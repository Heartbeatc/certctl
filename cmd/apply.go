package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"certctl/internal/acme"
	"certctl/internal/ai"
	"certctl/internal/cert"
	"certctl/internal/config"
	"certctl/internal/dns"
	"certctl/internal/i18n"
	"certctl/internal/ui"
	"certctl/pkg/domain"

	"github.com/briandowns/spinner"
	"github.com/go-acme/lego/v4/challenge"
	legolog "github.com/go-acme/lego/v4/log"
	"github.com/spf13/cobra"
)

var (
	flagDomain    string
	flagEmail     string
	flagOutput    string
	flagStaging   bool
	flagDryRun    bool
	flagLang      string
	flagDNS       string
	flagAliKey    string
	flagAliSecret string
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "申请 SSL 证书",
	Long:  "申请通配符 SSL 证书，支持 DNS-01 验证",
	RunE:  runApply,
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringVarP(&flagDomain, "domain", "d", "", "要申请证书的域名")
	applyCmd.Flags().StringVarP(&flagEmail, "email", "e", "", "Let's Encrypt 账户邮箱")
	applyCmd.Flags().StringVarP(&flagOutput, "output", "o", "", "证书输出目录")
	applyCmd.Flags().BoolVar(&flagStaging, "staging", false, "使用 Let's Encrypt 测试环境")
	applyCmd.Flags().BoolVar(&flagDryRun, "dry-run", false, "干跑模式，模拟流程不实际申请")
	applyCmd.Flags().StringVar(&flagLang, "lang", "", "语言 (zh/en)")
	applyCmd.Flags().StringVar(&flagDNS, "dns", "", "DNS 提供商 (aliyun)")
	applyCmd.Flags().StringVar(&flagAliKey, "ali-key", "", "阿里云 AccessKey ID")
	applyCmd.Flags().StringVar(&flagAliSecret, "ali-secret", "", "阿里云 AccessKey Secret")
}

func runApply(cmd *cobra.Command, args []string) error {
	// 设置语言
	if flagLang != "" {
		i18n.SetLang(flagLang)
	}

	// 如果没有指定输出目录，使用配置中的证书目录
	if flagOutput == "" {
		flagOutput = config.Get().CertsDir
		if flagOutput == "" {
			flagOutput = "./certs"
		}
	}

	// 获取详细模式设置
	verbose := config.GetVerbose()

	// 只在非详细模式下禁用 lego 日志
	if !verbose {
		legolog.Logger = &noopLogger{}
	}
	// 详细模式下使用标准日志输出，可以看到 ACME 请求和 DNS 操作的详细信息

	totalSteps := 5
	if flagDNS == "" {
		totalSteps = 6
	}

	var progress *ui.StepProgress
	if verbose {
		progress = ui.NewStepProgress(totalSteps)
	}

	if flagDryRun {
		ui.Warning(i18n.T("dryrun.warning"))
	}

	// Step 1: 收集参数
	if verbose && progress != nil {
		progress.Next(i18n.T("step.prepare"))
	}

	inputDomain := flagDomain
	if inputDomain == "" {
		inputDomain = ui.Prompt(i18n.T("prompt.domain"))
	}

	if inputDomain == "" {
		ui.ErrorWithHint(i18n.T("error.domain_empty"), []string{
			i18n.T("hint.domain_usage"),
		})
		return nil
	}

	domains, err := domain.GenerateWildcard(inputDomain)
	if err != nil {
		ui.ErrorWithHint(i18n.T("error.domain_invalid"), []string{
			fmt.Sprintf("Input: %s", inputDomain),
			i18n.T("hint.domain_format"),
		})
		return nil
	}
	rootDomain := domains[0]
	ui.Detail(fmt.Sprintf("%s: %s", i18n.T("detail.domain"), rootDomain))
	ui.Detail(fmt.Sprintf("%s: *.%s", i18n.T("detail.wildcard"), rootDomain))

	email := flagEmail
	if email == "" {
		email = ui.Prompt(i18n.T("prompt.email"))
	}
	if email == "" {
		ui.ErrorWithHint(i18n.T("error.email_empty"), []string{
			i18n.T("hint.email_usage"),
			i18n.T("hint.email_reminder"),
		})
		return nil
	}
	ui.Detail(fmt.Sprintf("%s: %s", i18n.T("detail.email"), email))
	ui.ProgressDone(i18n.T("progress.params_done"))

	// Step 2: 初始化客户端
	if verbose && progress != nil {
		progress.Next(i18n.T("step.init"))
	}

	configDir := getConfigDir()
	if verbose {
		ui.Detail(fmt.Sprintf("%s: %s", i18n.T("detail.config_dir"), configDir))
		ui.Detail(fmt.Sprintf("  ACME 服务: Let's Encrypt"))
		ui.Detail(fmt.Sprintf("  账户邮箱: %s", email))
	}

	account, err := acme.LoadOrCreateAccount(configDir, email)
	if err != nil {
		ui.ErrorWithHint(i18n.T("error.account_fail"), []string{
			fmt.Sprintf("Error: %v", err),
		})
		return nil
	}
	if verbose {
		ui.ProgressDone(i18n.T("progress.account_ok"))
	}

	// Step 3: 配置 DNS 验证
	if verbose && progress != nil {
		progress.Next(i18n.T("step.dns"))
	}

	var provider challenge.Provider

	if flagDNS == "aliyun" {
		aliKey := getAliKey()
		aliSecret := getAliSecret()

		if aliKey == "" || aliSecret == "" {
			ui.ErrorWithHint(i18n.T("error.ali_key_empty"), []string{
				i18n.T("hint.ali_key_param"),
				i18n.T("hint.ali_key_env"),
				i18n.T("hint.ali_key_url"),
			})
			return nil
		}

		if verbose {
			ui.Detail(fmt.Sprintf("%s: %s", i18n.T("detail.dns_mode"), i18n.T("detail.dns_aliyun")))
			ui.Detail(fmt.Sprintf("  AccessKey ID: %s****", aliKey[:4]))
			ui.Detail(fmt.Sprintf("  DNS API: dns.aliyuncs.com"))
		}
		aliyunProvider, err := acme.NewAliyunDNSProvider(aliKey, aliSecret, "")
		if err != nil {
			ui.ErrorWithHint(i18n.T("error.aliyun_fail"), []string{
				fmt.Sprintf("Error: %v", err),
			})
			return nil
		}
		provider = aliyunProvider
		if verbose {
			ui.ProgressDone(i18n.T("progress.aliyun_ready"))
		}
	} else {
		ui.Detail(fmt.Sprintf("%s: %s", i18n.T("detail.dns_mode"), i18n.T("detail.dns_manual")))
		provider = acme.NewManualDNSProvider(
			func(c *acme.Challenge) error {
				fmt.Println()
				ui.DNSRecord(c.RecordName, "TXT", c.Value, c.FQDN)

				fmt.Println()
				if !ui.Confirm(i18n.T("prompt.dns_added")) {
					return fmt.Errorf(i18n.T("error.user_cancel"))
				}

				if verbose && progress != nil {
					progress.Next(i18n.T("step.verify"))
				}
				spin := ui.NewSpinner(i18n.T("progress.checking_dns"))
				spin.Start()

				err := dns.WaitForRecord(c.FQDN, c.Value, 5*time.Minute, func(attempt int) {
					spin.Suffix = fmt.Sprintf(" (%d)", attempt)
				})
				spin.Stop()

				if err != nil {
					ui.ErrorWithHint(i18n.T("error.dns_fail"), []string{
						i18n.T("hint.dns_check"),
						i18n.T("hint.dns_wait"),
					})
					return err
				}
				ui.ProgressDone(i18n.T("progress.dns_ok"))
				return nil
			},
			nil,
		)
		ui.ProgressDone(i18n.T("progress.manual_ready"))
	}

	client, err := acme.NewClient(account, flagStaging, provider)
	if err != nil {
		ui.ErrorWithHint(i18n.T("error.client_fail"), []string{
			fmt.Sprintf("Error: %v", err),
		})
		return nil
	}

	if err := client.Register(); err != nil {
		errMsg := err.Error()
		hints := []string{fmt.Sprintf("Error: %v", err)}
		if strings.Contains(errMsg, "dial") || strings.Contains(errMsg, "timeout") {
			hints = append(hints, i18n.T("hint.check_network"))
			hints = append(hints, i18n.T("hint.china_blocked"))
		}
		ui.ErrorWithHint(i18n.T("error.register_fail"), hints)
		return nil
	}

	acme.SaveAccount(configDir, account)
	ui.ProgressDone(i18n.T("progress.client_ready"))

	// Step 4/5: 申请证书
	if verbose && progress != nil {
		progress.Next(i18n.T("step.apply"))
	}

	if verbose {
		if flagStaging {
			ui.Detail(i18n.T("detail.env_staging"))
		} else {
			ui.Detail(i18n.T("detail.env_prod"))
		}
		ui.Detail(fmt.Sprintf("  主域名: %s", rootDomain))
		ui.Detail(fmt.Sprintf("  通配符: *.%s", rootDomain))
		ui.Detail(fmt.Sprintf("  证书类型: 通配符证书 (Wildcard)"))
		ui.Detail(fmt.Sprintf("  有效期: 90 天"))
	}

	// 干跑模式
	if flagDryRun {
		ui.ProgressDone(i18n.T("dryrun.skip_apply"))

		if verbose && progress != nil {
			progress.Next(i18n.T("step.save"))
		}
		ui.ProgressDone(i18n.T("dryrun.skip_save"))

		if verbose && progress != nil {
			progress.Done(i18n.T("dryrun.done"))
		}
		ui.Detail(fmt.Sprintf("%s: %s, *.%s", i18n.T("detail.domain"), rootDomain, rootDomain))
		ui.Detail(fmt.Sprintf("%s: %s", i18n.T("detail.email"), email))
		ui.Detail(fmt.Sprintf("%s: %s", i18n.T("detail.output"), flagOutput))
		fmt.Println()
		return nil
	}

	var spin *spinner.Spinner
	if flagDNS == "aliyun" {
		spin = ui.NewSpinner(i18n.T("progress.applying"))
		spin.Start()
	}

	certificate, err := client.ObtainCertificate(domains)

	if spin != nil {
		spin.Stop()
	}

	if err != nil {
		errMsg := err.Error()

		// 如果 AI 启用，只显示错误并调用 AI 诊断（不显示内置提示以免干扰 AI）
		if config.IsAIEnabled() {
			ui.Error(i18n.T("error.cert_fail"))
			// 只有详细模式才显示原始错误信息
			if verbose {
				fmt.Println()
				ui.Detail(fmt.Sprintf("Error: %v", err))
			}
			fmt.Println()
			spin := ui.NewSpinner(i18n.T("ui.ai_diagnosing"))
			spin.Start()
			diagnosis, aiErr := ai.DiagnoseError(errMsg, rootDomain, flagDNS)
			spin.Stop()
			if aiErr != nil {
				ui.Info(fmt.Sprintf("AI 诊断失败: %v", aiErr))
			} else if diagnosis != "" {
				ui.AIBox(diagnosis)
			} else {
				ui.Info("AI 返回空结果")
			}
			// AI 诊断后直接返回，不询问重试
			fmt.Println()
			return nil
		} else {
			// AI 未启用，使用内置提示
			hints := []string{fmt.Sprintf("Error: %v", err)}
			if strings.Contains(errMsg, "DNS") || strings.Contains(errMsg, "TXT") {
				hints = append(hints, i18n.T("hint.dns_check"))
			}
			if strings.Contains(errMsg, "rate limit") {
				hints = append(hints, i18n.T("hint.rate_limit"))
			}
			ui.ErrorWithHint(i18n.T("error.cert_fail"), hints)
			fmt.Println()
			return nil
		}
	}

	ui.ProgressDone(i18n.T("progress.cert_ok"))

	// Step 5/6: 保存证书
	if verbose && progress != nil {
		progress.Next(i18n.T("step.save"))
	}

	certPath, keyPath, err := cert.Save(flagOutput, rootDomain, certificate.Certificate, certificate.PrivateKey)
	if err != nil {
		ui.ErrorWithHint(i18n.T("error.save_fail"), []string{
			fmt.Sprintf("Error: %v", err),
		})
		return nil
	}

	absOut, _ := filepath.Abs(flagOutput)
	certPath = filepath.Join(absOut, rootDomain, rootDomain+".pem")
	keyPath = filepath.Join(absOut, rootDomain, rootDomain+".key")

	ui.ProgressDone(i18n.T("progress.saved"))

	// 完成
	if verbose && progress != nil {
		progress.Done(i18n.T("progress.cert_ok"))
	}
	ui.CertResult(certPath, keyPath, certificate.NotAfter.Format("2006-01-02"))
	fmt.Println()

	return nil
}

func getAliKey() string {
	if flagAliKey != "" {
		return flagAliKey
	}
	if key := os.Getenv("ALICLOUD_ACCESS_KEY"); key != "" {
		return key
	}
	return ui.Prompt(i18n.T("prompt.ali_key"))
}

func getAliSecret() string {
	if flagAliSecret != "" {
		return flagAliSecret
	}
	if secret := os.Getenv("ALICLOUD_SECRET_KEY"); secret != "" {
		return secret
	}
	return ui.PromptSecret(i18n.T("prompt.ali_secret"))
}

func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".certctl"
	}
	return filepath.Join(home, ".certctl")
}
