package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"certctl/internal/acme"
	"certctl/internal/cert"
	"certctl/internal/dns"
	"certctl/internal/ui"
	"certctl/pkg/domain"

	legolog "github.com/go-acme/lego/v4/log"
	"github.com/spf13/cobra"
)

var (
	renewDomain  string
	renewEmail   string
	renewOutput  string
	renewStaging bool
)

var renewCmd = &cobra.Command{
	Use:   "renew",
	Short: "续期 SSL 证书",
	Long:  "续期已申请的 SSL 证书，需要重新进行 DNS 验证",
	RunE:  runRenew,
}

func init() {
	rootCmd.AddCommand(renewCmd)

	renewCmd.Flags().StringVarP(&renewDomain, "domain", "d", "", "要续期的域名")
	renewCmd.Flags().StringVarP(&renewEmail, "email", "e", "", "Let's Encrypt 账户邮箱（可选，使用已保存的账户）")
	renewCmd.Flags().StringVarP(&renewOutput, "output", "o", "./certs", "证书输出目录")
	renewCmd.Flags().BoolVar(&renewStaging, "staging", false, "使用 Let's Encrypt 测试环境")
}

func runRenew(cmd *cobra.Command, args []string) error {
	// 禁用 lego 库的日志输出
	legolog.Logger = &noopLogger{}

	fmt.Println()

	// 1. 获取域名
	inputDomain := renewDomain
	if inputDomain == "" {
		inputDomain = ui.Prompt("请输入要续期的域名:")
		if inputDomain == "" {
			ui.Error("域名不能为空")
			return nil
		}
	}

	// 2. 生成通配符域名
	domains, err := domain.GenerateWildcard(inputDomain)
	if err != nil {
		ui.Error(fmt.Sprintf("域名解析失败: %v", err))
		return nil
	}

	rootDomain := domains[0]

	// 3. 检查证书是否存在
	certPath := filepath.Join(renewOutput, rootDomain, rootDomain+".pem")
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		ui.Warning(fmt.Sprintf("未找到域名 %s 的证书，将进行首次申请", rootDomain))
	} else {
		// 显示当前证书信息
		notAfter, err := cert.ParseCertExpiry(certPath)
		if err == nil {
			daysLeft := int(time.Until(notAfter).Hours() / 24)
			if daysLeft > 30 {
				ui.Info(fmt.Sprintf("当前证书还有 %d 天过期 (%s)", daysLeft, notAfter.Format("2006-01-02")))
				if !ui.Confirm("确定要续期吗?") {
					ui.Info("已取消")
					return nil
				}
			} else {
				ui.Warning(fmt.Sprintf("当前证书将于 %d 天后过期 (%s)", daysLeft, notAfter.Format("2006-01-02")))
			}
		}
	}

	// 4. 确认续期
	ui.Title("将为以下域名续期证书:")
	fmt.Println()
	ui.DomainList(domains)
	fmt.Println()

	if !ui.Confirm("继续?") {
		ui.Info("已取消")
		return nil
	}

	fmt.Println()

	// 5. 加载账户
	configDir := getConfigDir()

	spin := ui.NewSpinner("正在加载账户信息...")
	spin.Start()

	// 尝试加载已有账户
	var account *acme.Account
	email := renewEmail

	// 如果没有指定邮箱，尝试从配置中加载
	if email == "" {
		existingAccount, err := acme.LoadOrCreateAccount(configDir, "")
		if err == nil && existingAccount.Email != "" {
			account = existingAccount
			email = existingAccount.Email
		}
	}

	if account == nil {
		spin.Stop()
		if email == "" {
			email = ui.Prompt("请输入邮箱 (用于 Let's Encrypt 账户):")
			if email == "" {
				ui.Error("邮箱不能为空")
				return nil
			}
		}
		spin = ui.NewSpinner("正在初始化 ACME 客户端...")
		spin.Start()

		account, err = acme.LoadOrCreateAccount(configDir, email)
		if err != nil {
			spin.Stop()
			ui.Error(fmt.Sprintf("加载账户失败: %v", err))
			return nil
		}
	}

	// 创建 DNS Provider，onPresent 回调会阻塞直到 DNS 验证通过
	provider := acme.NewManualDNSProvider(
		func(c *acme.Challenge) error {
			// 显示 DNS 记录信息
			fmt.Println()
			ui.DNSRecord(
				c.RecordName,
				"TXT",
				c.Value,
				c.FQDN,
			)
			fmt.Println()

			ui.Info("💡 如果之前已添加过 TXT 记录，请更新记录值")
			fmt.Println()

			if !ui.Confirm("已添加/更新 DNS 记录?") {
				return fmt.Errorf("用户取消")
			}

			fmt.Println()

			// 检查 DNS 记录
			checkSpin := ui.NewSpinner("检查 DNS 记录是否生效...")
			checkSpin.Start()

			err := dns.WaitForRecord(c.FQDN, c.Value, 5*time.Minute, func(attempt int) {
				checkSpin.Suffix = fmt.Sprintf(" 检查 DNS 记录... (第 %d 次)", attempt)
			})

			checkSpin.Stop()

			if err != nil {
				ui.Error("DNS 记录验证超时，请确认记录已正确添加")
				return err
			}

			ui.Success("DNS 记录已生效")
			fmt.Println()

			return nil
		},
		nil,
	)

	client, err := acme.NewClient(account, renewStaging, provider)
	if err != nil {
		spin.Stop()
		ui.Error(fmt.Sprintf("创建客户端失败: %v", err))
		return nil
	}

	// 注册账户
	if err := client.Register(); err != nil {
		spin.Stop()
		ui.Error(fmt.Sprintf("注册账户失败: %v", err))
		return nil
	}

	// 保存账户
	if err := acme.SaveAccount(configDir, account); err != nil {
		spin.Stop()
		ui.Warning(fmt.Sprintf("保存账户失败: %v", err))
	}

	spin.Stop()
	ui.Success("ACME 客户端就绪")
	fmt.Println()

	// 6. 申请证书（onPresent 会阻塞等待 DNS 验证通过）
	ui.Info("正在与 Let's Encrypt 通信...")
	fmt.Println()

	certificate, err := client.ObtainCertificate(domains)

	if err != nil {
		ui.Error(fmt.Sprintf("证书续期失败: %v", err))
		return nil
	}

	ui.Success("证书续期成功!")
	fmt.Println()

	// 7. 保存证书
	certPathNew, keyPath, err := cert.Save(renewOutput, rootDomain, certificate.Certificate, certificate.PrivateKey)
	if err != nil {
		ui.Error(fmt.Sprintf("保存证书失败: %v", err))
		return nil
	}

	// 转换为绝对路径
	absOut, _ := filepath.Abs(renewOutput)
	certPathNew = filepath.Join(absOut, rootDomain, rootDomain+".pem")
	keyPath = filepath.Join(absOut, rootDomain, rootDomain+".key")

	// 8. 显示结果
	ui.CertResult(certPathNew, keyPath, certificate.NotAfter.Format("2006-01-02"))
	fmt.Println()

	return nil
}
