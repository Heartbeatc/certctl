package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"certctl/internal/ai"
	"certctl/internal/cert"
	"certctl/internal/config"
	"certctl/internal/i18n"
	"certctl/internal/ui"

	"github.com/spf13/cobra"
)

// 全局证书目录（从配置加载）
var certsDir string

func init() {
	// 加载配置
	cfg := config.Load()
	certsDir = cfg.CertsDir
	i18n.SetLang(cfg.Language)
}

var rootCmd = &cobra.Command{
	Use:   "certctl",
	Short: "轻量级 SSL 证书申请工具",
	Long: `certctl - 一个 CLI 风格的 SSL 证书申请工具

支持通配符证书申请，使用 Let's Encrypt 作为 CA，
通过 DNS-01 验证方式完成域名所有权验证。

输出 Nginx 格式的证书文件。`,
	Run: runMainMenu,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// runMainMenu 主菜单
func runMainMenu(cmd *cobra.Command, args []string) {
	showMainScreen()

	for {
		idx, _, err := ui.SelectWithExit(i18n.T("ui.select_operation"), []string{
			i18n.T("ui.apply_cert"),
			i18n.T("ui.renew_cert"),
			i18n.T("ui.settings"),
		}, i18n.T("ui.exit"))

		if err != nil || idx == -1 {
			ui.ClearScreen()
			ui.Info(i18n.T("ui.goodbye"))
			return
		}

		// 清屏进入功能页面
		ui.ClearScreen()

		needPressKey := true // 是否需要按键等待

		switch idx {
		case 0:
			runApplyInteractive()
		case 1:
			runRenewInteractive()
			needPressKey = false // 续期有自己的菜单选择
		case 2:
			showSettings()
			needPressKey = false // 设置菜单有自己的循环，不需要等待
		}

		// 只在需要时等待用户确认
		if needPressKey {
			ui.PressAnyKey()
		}

		// 清屏并重新显示主菜单
		ui.ClearScreen()
		showMainScreen()
	}
}

// showMainScreen 显示主界面
func showMainScreen() {
	ui.Logo()
}

// showCertsList 显示证书列表
func showCertsList() {
	ui.Header("查看证书")
	runList(nil, nil)
}


// switchLanguage 切换语言
func switchLanguage() {
	idx, _, err := ui.Select("选择语言 / Select Language:", []string{
		"中文",
		"English",
	})

	if err != nil {
		return
	}

	if idx == 0 {
		i18n.SetLang("zh")
	} else {
		i18n.SetLang("en")
	}
	ui.Success(fmt.Sprintf("语言已切换为: %s", i18n.Lang))
}

// setCertsDir 设置证书目录
func setCertsDir() {
	var newDir string
	var err error

	// 根据系统选择不同的方式
	switch runtime.GOOS {
	case "darwin":
		// Mac 使用 osascript 调用原生对话框
		newDir, err = selectFolderMac()
		if err != nil {
			ui.Info(i18n.T("ui.cancelled"))
			return
		}
	case "windows":
		// Windows 使用 dialog 库
		newDir, err = selectFolderWindows(i18n.T("ui.select_certs_folder"))
		if err != nil {
			ui.Info(i18n.T("ui.cancelled"))
			return
		}
	default:
		// Linux 等其他系统手动输入
		newDir, err = ui.Input(i18n.T("ui.certs_dir"), certsDir)
		if err != nil {
			return
		}
	}

	if newDir != "" {
		certsDir = newDir
		flagOutput = certsDir
		ui.Success(fmt.Sprintf(i18n.T("ui.dir_set_to"), certsDir))
	}
}

// selectFolderMac 使用 osascript 打开 Mac 文件夹选择器
func selectFolderMac() (string, error) {
	prompt := i18n.T("ui.select_certs_folder")
	script := fmt.Sprintf(`set folderPath to POSIX path of (choose folder with prompt "%s")`, prompt)
	cmd := exec.Command("osascript", "-e", `tell application "System Events" to activate`, "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// 去掉末尾换行符
	path := strings.TrimSpace(string(output))
	// 去掉末尾的斜杠
	path = strings.TrimSuffix(path, "/")
	return path, nil
}

// showSettings 设置子菜单
func showSettings() {
	for {
		// 每次循环开始时清屏
		ui.ClearScreen()
		ui.Header(i18n.T("ui.settings_title"))
		fmt.Println()
		ui.StatusLine(i18n.T("ui.language"), i18n.Lang)
		ui.StatusLine(i18n.T("ui.certs_dir"), certsDir)

		// 显示详细模式状态
		verboseStatus := i18n.T("ui.off")
		if config.GetVerbose() {
			verboseStatus = i18n.T("ui.on")
		}
		ui.StatusLine(i18n.T("ui.verbose"), verboseStatus)

		// 显示已配置的 DNS，显示配置名称
		dnsConfigs := config.GetDNSConfigs()
		if len(dnsConfigs) > 0 {
			// 显示所有配置名称
			names := []string{}
			for _, cfg := range dnsConfigs {
				names = append(names, cfg.Name)
			}
			ui.StatusLine(i18n.T("ui.dns_config"), strings.Join(names, ", "))
		} else {
			ui.StatusLine(i18n.T("ui.dns_config"), i18n.T("ui.not_configured"))
		}

		// 显示 AI 增强模式状态
		aiStatus := i18n.T("ui.ai_disabled")
		if config.IsAIEnabled() {
			aiStatus = i18n.T("ui.ai_enabled")
		}
		ui.StatusLine(i18n.T("ui.ai_config"), aiStatus)

		ui.StatusLine(i18n.T("ui.config_dir"), config.GetConfigDir())
		fmt.Println()

		// 构建菜单选项，带状态（复用上面已计算的verboseStatus）
		aiMenuStatus := i18n.T("ui.off")
		if config.IsAIEnabled() {
			aiMenuStatus = i18n.T("ui.on")
		}

		idx, _, err := ui.Select(i18n.T("ui.select_setting"), []string{
			i18n.T("ui.language_set"),
			i18n.T("ui.certs_dir_set"),
			fmt.Sprintf("%s (%s)", i18n.T("ui.verbose_set"), verboseStatus),
			i18n.T("ui.dns_vendor_config"),
			fmt.Sprintf("%s (%s)", i18n.T("ui.ai_config"), aiMenuStatus),
			i18n.T("ui.back_main"),
		})

		if err != nil {
			return
		}

		switch idx {
		case 0:
			// 切换语言
			langIdx, _, _ := ui.Select(i18n.T("ui.select_language"), []string{i18n.T("ui.chinese"), i18n.T("ui.english")})
			if langIdx == 0 {
				i18n.SetLang("zh")
				config.SetLanguage("zh")
			} else {
				i18n.SetLang("en")
				config.SetLanguage("en")
			}
			ui.Success(fmt.Sprintf(i18n.T("ui.lang_set_to"), i18n.Lang))
			return // 返回主菜单以刷新界面
		case 1:
			// 设置证书目录
			setCertsDirInner()
		case 2:
			// 切换详细模式
			newVerbose := !config.GetVerbose()
			config.SetVerbose(newVerbose)
			status := i18n.T("ui.off")
			if newVerbose {
				status = i18n.T("ui.on")
			}
			ui.Success(fmt.Sprintf(i18n.T("ui.verbose_set_to"), status))
		case 3:
			// DNS 配置
			manageDNSConfig()
		case 4:
			// AI 增强模式
			manageAIConfig()
		case 5:
			// 返回
			return
		}
	}
}

// manageDNSConfig 管理 DNS 配置
func manageDNSConfig() {
	for {
		// 清屏
		ui.ClearScreen()
		ui.Header(i18n.T("ui.dns_config_title"))
		fmt.Println()

		dnsConfigs := config.GetDNSConfigs()
		if len(dnsConfigs) > 0 {
			for _, dns := range dnsConfigs {
				maskedKey := dns.AccessKeyID
				if len(maskedKey) > 8 {
					maskedKey = maskedKey[:4] + "****" + maskedKey[len(maskedKey)-4:]
				}
				ui.StatusLine(dns.Name, fmt.Sprintf("%s (%s)", dns.Provider, maskedKey))
			}
		} else {
			ui.Info(i18n.T("ui.no_dns_config"))
		}
		fmt.Println()

		idx, _, err := ui.Select(i18n.T("ui.select_action"), []string{
			i18n.T("ui.add_aliyun"),
			i18n.T("ui.add_tencentcloud"),
			i18n.T("ui.delete_config"),
			i18n.T("ui.back"),
		})

		if err != nil {
			return
		}

		switch idx {
		case 0:
			// 添加阿里云配置
			addAliyunDNSConfig()
		case 1:
			// 添加腾讯云配置
			addTencentCloudDNSConfig()
		case 2:
			// 删除配置
			deletesDNSConfig()
		case 3:
			return
		}
	}
}

// addAliyunDNSConfig 添加阿里云 DNS 配置（支持命名）
func addAliyunDNSConfig() {
	// 输入配置名称
	name, err := ui.Input(i18n.T("ui.config_name"), "")
	if err != nil || name == "" {
		ui.Error(i18n.T("ui.name_empty"))
		return
	}

	accessKeyID, err := ui.Input(i18n.T("ui.ali_key_id"), "")
	if err != nil || accessKeyID == "" {
		ui.Error(i18n.T("ui.ali_key_empty"))
		return
	}

	accessKeySecret, err := ui.InputSecret(i18n.T("ui.ali_secret"))
	if err != nil || accessKeySecret == "" {
		ui.Error(i18n.T("ui.ali_secret_empty"))
		return
	}

	config.AddDNSConfig(name, "aliyun", accessKeyID, accessKeySecret)
	ui.Success(fmt.Sprintf(i18n.T("ui.config_saved"), name))
}

// deletesDNSConfig 删除 DNS 配置
func deletesDNSConfig() {
	dnsConfigs := config.GetDNSConfigs()
	if len(dnsConfigs) == 0 {
		ui.Info(i18n.T("ui.no_config_delete"))
		return
	}

	names := []string{}
	for _, dns := range dnsConfigs {
		names = append(names, fmt.Sprintf("%s (%s)", dns.Name, dns.Provider))
	}
	names = append(names, i18n.T("ui.cancel"))

	idx, result, err := ui.Select(i18n.T("ui.select_delete"), names)
	if err != nil || result == i18n.T("ui.cancel") {
		return
	}

	configName := dnsConfigs[idx].Name
	config.DeleteDNSConfig(configName)
	ui.Success(fmt.Sprintf(i18n.T("ui.deleted"), configName))
}

// addTencentCloudDNSConfig 添加腾讯云 DNS 配置
func addTencentCloudDNSConfig() {
	// 输入配置名称
	name, err := ui.Input(i18n.T("ui.config_name"), "")
	if err != nil || name == "" {
		ui.Error(i18n.T("ui.name_empty"))
		return
	}

	secretId, err := ui.Input("腾讯云 SecretId", "")
	if err != nil || secretId == "" {
		ui.Error("SecretId 不能为空")
		return
	}

	secretKey, err := ui.InputSecret("腾讯云 SecretKey")
	if err != nil || secretKey == "" {
		ui.Error("SecretKey 不能为空")
		return
	}

	config.AddDNSConfig(name, "tencentcloud", secretId, secretKey)
	ui.Success(fmt.Sprintf(i18n.T("ui.config_saved"), name))
}

// setCertsDirInner 内部设置证书目录（不按键返回）
func setCertsDirInner() {
	var newDir string
	var err error

	switch runtime.GOOS {
	case "darwin":
		newDir, err = selectFolderMac()
		if err != nil {
			ui.Info(i18n.T("ui.cancelled"))
			return
		}
	case "windows":
		newDir, err = selectFolderWindows(i18n.T("ui.certs_dir"))
		if err != nil {
			ui.Info(i18n.T("ui.cancelled"))
			return
		}
	default:
		newDir, err = ui.Input(i18n.T("ui.certs_dir"), certsDir)
		if err != nil {
			return
		}
	}

	if newDir != "" {
		certsDir = newDir
		flagOutput = certsDir
		config.SetCertsDir(certsDir)
		ui.Success(fmt.Sprintf(i18n.T("ui.dir_set_to"), certsDir))
	}
}

// manageAIConfig 管理 AI 增强模式配置
func manageAIConfig() {
	for {
		ui.ClearScreen()
		ui.Header(i18n.T("ui.ai_config_title"))
		fmt.Println()

		aiCfg := config.GetAIConfig()
		status := i18n.T("ui.ai_disabled")
		if aiCfg.Enabled {
			status = i18n.T("ui.ai_enabled")
		}
		ui.StatusLine(i18n.T("ui.ai_status"), status)

		if aiCfg.APIKey != "" {
			// 显示掩码后的 API Key
			key := aiCfg.APIKey
			maskedKey := key
			if len(key) > 8 {
				maskedKey = key[:4] + "****" + key[len(key)-4:]
			}
			ui.StatusLine("API Key", maskedKey)
		} else {
			ui.StatusLine("API Key", i18n.T("ui.not_configured"))
		}
		ui.StatusLine(i18n.T("ui.ai_model"), aiCfg.Model)
		fmt.Println()
		ui.Info(i18n.T("ui.ai_key_hint"))
		fmt.Println()

		idx, _, err := ui.Select(i18n.T("ui.select_action"), []string{
			i18n.T("ui.toggle_ai"),
			i18n.T("ui.set_api_key"),
			i18n.T("ui.select_model"),
			i18n.T("ui.test_connection"),
			i18n.T("ui.back"),
		})

		if err != nil {
			return
		}

		switch idx {
		case 0:
			// 启用/禁用 AI
			newEnabled := !aiCfg.Enabled
			config.SetAIEnabled(newEnabled)
			if newEnabled {
				ui.Success(i18n.T("ui.ai_enabled_msg"))
			} else {
				ui.Info(i18n.T("ui.ai_disabled_msg"))
			}
		case 1:
			// 设置 API Key
			key, err := ui.Input(i18n.T("ui.ai_key_prompt"), "")
			if err != nil || key == "" {
				continue
			}
			config.SetAIAPIKey(key)
			ui.Success(i18n.T("ui.ai_key_saved"))
		case 2:
			// 选择模型
			models := []string{"glm-4-flash (推荐)", "glm-4.7", "glm-4-plus"}
			modelValues := []string{"glm-4-flash", "glm-4.7", "glm-4-plus"}
			mIdx, _, err := ui.Select(i18n.T("ui.select_model"), models)
			if err != nil {
				continue
			}
			config.SetAIModel(modelValues[mIdx])
			ui.Success(fmt.Sprintf(i18n.T("ui.ai_model_saved"), modelValues[mIdx]))
		case 3:
			// 测试连接
			if aiCfg.APIKey == "" {
				ui.Error(i18n.T("ui.not_configured"))
				ui.PressAnyKey()
				continue
			}
			spin := ui.NewSpinner(i18n.T("ui.ai_testing"))
			spin.Start()
			client := ai.NewZhipuClient()
			err := client.TestConnection()
			spin.Stop()
			if err != nil {
				ui.Error(fmt.Sprintf(i18n.T("ui.ai_test_fail"), err))
			} else {
				ui.Success(i18n.T("ui.ai_test_ok"))
			}
			ui.PressAnyKey()
		case 4:
			return
		}
	}
}

// runApplyInteractive 交互式申请证书
func runApplyInteractive() {
	ui.Header(i18n.T("ui.apply_title"))

	// 1. 输入域名
	domain, err := ui.Input(i18n.T("ui.domain"), "")
	if err != nil || domain == "" {
		ui.Error(i18n.T("ui.domain_empty"))
		return
	}

	// 2. 输入邮箱
	email, err := ui.Input(i18n.T("ui.email"), "")
	if err != nil || email == "" {
		ui.Error(i18n.T("ui.email_empty"))
		return
	}

	// 3. 选择 DNS 验证方式
	dnsIdx, _, err := ui.Select(i18n.T("ui.dns_verify_method"), []string{
		i18n.T("ui.aliyun_auto"),
		"腾讯云NS(自动)",
		i18n.T("ui.manual_dns"),
	})
	if err != nil {
		return
	}

	var dnsProvider string
	if dnsIdx == 0 {
		dnsProvider = "aliyun"
	} else if dnsIdx == 1 {
		dnsProvider = "tencentcloud"
	}

	// 4. 如果是阿里云，检查已保存的配置
	var aliKey, aliSecret string
	var tencentId, tencentSecret string
	if dnsProvider == "aliyun" {
		// 获取所有阿里云配置
		aliyunConfigs := config.GetDNSConfigsByProvider("aliyun")

		if len(aliyunConfigs) > 0 {
			// 有已保存的配置，让用户选择
			options := []string{}
			for _, cfg := range aliyunConfigs {
				maskedKey := cfg.AccessKeyID
				if len(maskedKey) > 8 {
					maskedKey = maskedKey[:4] + "****" + maskedKey[len(maskedKey)-4:]
				}
				options = append(options, fmt.Sprintf("%s (%s)", cfg.Name, maskedKey))
			}
			options = append(options, i18n.T("ui.input_new_config"))

			cfgIdx, _, _ := ui.Select(i18n.T("ui.select_aliyun_cfg"), options)

			if cfgIdx < len(aliyunConfigs) {
				// 使用已保存的配置
				selectedConfig := aliyunConfigs[cfgIdx]
				aliKey = selectedConfig.AccessKeyID
				aliSecret = selectedConfig.AccessKeySecret
				ui.Success(fmt.Sprintf(i18n.T("ui.using_config"), selectedConfig.Name))
			} else {
				// 输入新配置
				aliKey, aliSecret = inputAliyunCredentials()
				if aliKey == "" || aliSecret == "" {
					return
				}
				// 询问是否保存
				if ui.ConfirmPrompt(i18n.T("ui.save_config_local")) {
					name, _ := ui.Input(i18n.T("ui.config_name_new"), "")
					if name != "" {
						config.AddDNSConfig(name, "aliyun", aliKey, aliSecret)
						ui.Success(fmt.Sprintf(i18n.T("ui.config_saved"), name))
					}
				}
			}
		} else {
			// 没有保存的配置，需要输入
			aliKey, aliSecret = inputAliyunCredentials()
			if aliKey == "" || aliSecret == "" {
				return
			}
			// 询问是否保存
			if ui.ConfirmPrompt(i18n.T("ui.save_config_local")) {
				name, _ := ui.Input(i18n.T("ui.config_name"), "")
				if name != "" {
					config.AddDNSConfig(name, "aliyun", aliKey, aliSecret)
					ui.Success(fmt.Sprintf(i18n.T("ui.config_saved"), name))
				}
			}
		}
	} else if dnsProvider == "tencentcloud" {
		// 获取所有腾讯云配置
		tencentConfigs := config.GetDNSConfigsByProvider("tencentcloud")

		if len(tencentConfigs) > 0 {
			// 有已保存的配置，让用户选择
			options := []string{}
			for _, cfg := range tencentConfigs {
				maskedKey := cfg.AccessKeyID
				if len(maskedKey) > 8 {
					maskedKey = maskedKey[:4] + "****" + maskedKey[len(maskedKey)-4:]
				}
				options = append(options, fmt.Sprintf("%s (%s)", cfg.Name, maskedKey))
			}
			options = append(options, i18n.T("ui.input_new_config"))

			cfgIdx, _, _ := ui.Select("选择腾讯云配置:", options)

			if cfgIdx < len(tencentConfigs) {
				// 使用已保存的配置
				selectedConfig := tencentConfigs[cfgIdx]
				tencentId = selectedConfig.AccessKeyID
				tencentSecret = selectedConfig.AccessKeySecret
				ui.Success(fmt.Sprintf(i18n.T("ui.using_config"), selectedConfig.Name))
			} else {
				// 输入新配置
				tencentId, _ = ui.Input("腾讯云 SecretId", os.Getenv("TENCENTCLOUD_SECRET_ID"))
				tencentSecret, _ = ui.InputSecret("腾讯云 SecretKey")
				if tencentId == "" || tencentSecret == "" {
					return
				}
				// 询问是否保存
				if ui.ConfirmPrompt(i18n.T("ui.save_config_local")) {
					name, _ := ui.Input(i18n.T("ui.config_name_new"), "")
					if name != "" {
						config.AddDNSConfig(name, "tencentcloud", tencentId, tencentSecret)
						ui.Success(fmt.Sprintf(i18n.T("ui.config_saved"), name))
					}
				}
			}
		} else {
			// 没有保存的配置，需要输入
			tencentId, _ = ui.Input("腾讯云 SecretId", os.Getenv("TENCENTCLOUD_SECRET_ID"))
			tencentSecret, _ = ui.InputSecret("腾讯云 SecretKey")
			if tencentId == "" || tencentSecret == "" {
				return
			}
			// 询问是否保存
			if ui.ConfirmPrompt(i18n.T("ui.save_config_local")) {
				name, _ := ui.Input(i18n.T("ui.config_name"), "")
				if name != "" {
					config.AddDNSConfig(name, "tencentcloud", tencentId, tencentSecret)
					ui.Success(fmt.Sprintf(i18n.T("ui.config_saved"), name))
				}
			}
		}
	}

	// 5. 确认
	fmt.Println()
	ui.Info(i18n.T("ui.will_apply"))
	ui.Detail(fmt.Sprintf(i18n.T("ui.domain_info"), domain, domain))
	ui.Detail(fmt.Sprintf(i18n.T("ui.email_info"), email))
	if dnsProvider == "aliyun" {
		ui.Detail(i18n.T("ui.dns_aliyun_auto"))
	} else if dnsProvider == "tencentcloud" {
		ui.Detail("DNS: 腾讯云自动验证")
	} else {
		ui.Detail(i18n.T("ui.dns_manual"))
	}
	fmt.Println()

	if !ui.ConfirmPrompt(i18n.T("ui.confirm_apply")) {
		ui.Info(i18n.T("ui.cancelled_op"))
		return
	}

	// 6. 设置参数并执行原有逻辑
	flagDomain = domain
	flagEmail = email
	flagDNS = dnsProvider
	flagAliKey = aliKey
	flagAliSecret = aliSecret
	flagTencentId = tencentId
	flagTencentSecret = tencentSecret
	flagDryRun = false

	runApply(nil, nil)
}

// runRenewInteractive 交互式续期证书
func runRenewInteractive() {
	ui.Header(i18n.T("ui.renew_title"))

	// 获取证书目录
	certsDir := config.Get().CertsDir
	if certsDir == "" {
		certsDir = "./certs"
	}

	// 扫描已有证书
	certs, err := cert.ListCertificates(certsDir)
	if err != nil {
		ui.Error(fmt.Sprintf("扫描证书失败: %v", err))
		return
	}

	if len(certs) == 0 {
		ui.Info(i18n.T("ui.no_certs"))
		fmt.Println()
		ui.Info("请先使用「申请证书」功能申请证书")
		return
	}

	// 构建选择列表 - 先计算最大域名长度用于对齐
	maxDomainLen := 0
	for _, c := range certs {
		if len(c.Domain) > maxDomainLen {
			maxDomainLen = len(c.Domain)
		}
	}

	options := []string{}
	for _, c := range certs {
		status := ""
		if c.DaysLeft <= 0 {
			status = fmt.Sprintf("%-10s", i18n.T("ui.expired"))
		} else if c.DaysLeft <= 30 {
			status = fmt.Sprintf("%3d %s ⚠️", c.DaysLeft, i18n.T("ui.days"))
		} else {
			status = fmt.Sprintf("%3d %s   ", c.DaysLeft, i18n.T("ui.days"))
		}
		// 只显示目录，不显示文件名
		certDir := filepath.Dir(c.CertPath)
		// 使用固定宽度格式化域名
		domainPadded := fmt.Sprintf("%-*s", maxDomainLen, c.Domain)
		options = append(options, fmt.Sprintf("%s  %s  %s", domainPadded, status, certDir))
	}
	options = append(options, i18n.T("ui.back"))

	// 让用户选择
	idx, _, err := ui.Select(i18n.T("ui.select_renew_cert"), options)
	if err != nil || idx == len(certs) {
		return
	}

	// 执行续期
	selectedCert := certs[idx]
	renewDomain = selectedCert.Domain
	renewOutput = certsDir
	runRenew(nil, nil)
}

// inputAliyunCredentials 输入阿里云凭证
func inputAliyunCredentials() (string, string) {
	aliKey, err := ui.Input(i18n.T("ui.ali_key_id"), os.Getenv("ALICLOUD_ACCESS_KEY"))
	if err != nil || aliKey == "" {
		ui.Error(i18n.T("ui.ali_key_empty"))
		return "", ""
	}

	aliSecret, err := ui.InputSecret(i18n.T("ui.ali_secret"))
	if err != nil || aliSecret == "" {
		ui.Error(i18n.T("ui.ali_secret_empty"))
		return "", ""
	}

	return aliKey, aliSecret
}

// noopLogger 实现 lego 的 StdLogger 接口，丢弃所有日志
type noopLogger struct{}

func (n *noopLogger) Fatal(args ...interface{})                 {}
func (n *noopLogger) Fatalln(args ...interface{})               {}
func (n *noopLogger) Fatalf(format string, args ...interface{}) {}
func (n *noopLogger) Print(args ...interface{})                 {}
func (n *noopLogger) Println(args ...interface{})               {}
func (n *noopLogger) Printf(format string, args ...interface{}) {}
