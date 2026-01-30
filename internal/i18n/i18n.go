package i18n

import "os"

// Lang 当前语言
var Lang = "zh"

// 初始化语言
func init() {
	// 检查环境变量
	if lang := os.Getenv("LANG"); lang != "" {
		if len(lang) >= 2 && lang[:2] == "en" {
			Lang = "en"
		}
	}
	if lang := os.Getenv("CERTCTL_LANG"); lang != "" {
		Lang = lang
	}
}

// SetLang 设置语言
func SetLang(lang string) {
	if lang == "en" || lang == "zh" {
		Lang = lang
	}
}

// T 翻译函数
func T(key string) string {
	if Lang == "en" {
		if msg, ok := enMessages[key]; ok {
			return msg
		}
	}
	if msg, ok := zhMessages[key]; ok {
		return msg
	}
	return key
}

// 中文消息
var zhMessages = map[string]string{
	// 步骤
	"step.prepare":      "准备申请信息",
	"step.init":         "初始化ACME客户端",
	"step.dns":          "配置DNS验证",
	"step.verify":       "验证DNS记录",
	"step.apply":        "申请证书",
	"step.save":         "保存证书",

	// 提示
	"prompt.domain":     "域名",
	"prompt.email":      "邮箱",
	"prompt.ali_key":    "阿里云 AccessKey ID",
	"prompt.ali_secret": "阿里云 AccessKey Secret",
	"prompt.dns_added":  "已添加 DNS记录?",

	// 详情
	"detail.domain":       "主域名",
	"detail.wildcard":     "通配符",
	"detail.email":        "邮箱",
	"detail.config_dir":   "配置目录",
	"detail.dns_mode":     "DNS模式",
	"detail.dns_aliyun":   "阿里云自动验证",
	"detail.dns_manual":   "手动验证",
	"detail.env_staging":  "Let's Encrypt 测试环境",
	"detail.env_prod":     "Let's Encrypt 生产环境",
	"detail.output":       "输出目录",

	// 进度
	"progress.params_done":   "参数收集完成",
	"progress.account_ok":    "账户加载成功",
	"progress.aliyun_ready":     "阿里云 DNS 就绪",
	"progress.tencentcloud_ready": "腾讯云 DNS 就绪",
	"progress.manual_ready":  "手动验证模式就绪",
	"progress.client_ready":  "ACME 客户端就绪",
	"progress.dns_ok":        "DNS 记录已生效",
	"progress.cert_ok":       "证书申请成功",
	"progress.saved":         "证书已保存",
	"progress.checking_dns":  "正在检查 DNS 记录传播...",
	"progress.applying":      "正在自动添加 DNS 记录并申请证书...",

	// 干跑
	"dryrun.warning":  "干跑模式：不会实际申请证书",
	"dryrun.skip_apply": "干跑模式：跳过证书申请",
	"dryrun.skip_save":  "干跑模式：跳过证书保存",
	"dryrun.done":       "干跑完成！参数验证通过，可以正式申请",

	// 错误
	"error.domain_empty":     "域名不能为空",
	"error.domain_invalid":   "域名解析失败",
	"error.domain_parse":     "解析域名失败: %v",
	"error.email_empty":      "邮箱不能为空",
	"error.ali_key_empty":    "阿里云 AccessKey 不完整",
	"error.account_fail":     "账户初始化失败",
	"error.aliyun_fail":      "阿里云 DNS 客户端创建失败",
	"error.aliyun_create":    "创建阿里云DNS客户端失败: %v",
	"error.tencentcloud_fail":   "腾讯云 DNS 客户端创建失败",
	"error.tencentcloud_create": "创建腾讯云DNS客户端失败: %v",
	"error.client_fail":      "ACME 客户端创建失败",
	"error.client_create":    "创建ACME客户端失败: %v",
	"error.dns_provider":     "设置DNS提供者失败: %v",
	"error.register_fail":    "账户注册失败",
	"error.register":         "注册账户失败: %v",
	"error.dns_fail":         "DNS 验证失败",
	"error.dns_add":          "添加DNS记录失败: %v",
	"error.dns_update":       "更新DNS记录失败: %v",
	"error.dns_delete":       "删除DNS记录失败: %v",
	"error.dns_query":        "查询DNS记录失败: %v",
	"error.dns_query_fail":   "DNS查询失败: %s",
	"error.dns_timeout":      "DNS记录验证超时",
	"error.cert_fail":        "证书申请失败",
	"error.cert_obtain":      "申请证书失败: %v",
	"error.cert_parse":       "解析证书失败",
	"error.cert_parse_err":   "解析证书失败: %v",
	"error.save_fail":        "证书保存失败",
	"error.user_cancel":      "用户取消操作",
	"error.ai_request":       "AI请求失败: %v",
	"error.ai_parse":         "AI响应解析失败: %v",
	"error.ai_error":         "AI错误: %s",
	"error.ai_no_response":   "AI无响应",
	"error.ai_connect":       "连接失败: %v",
	"error.api_error":        "API错误: %s",
	"error.no_response":      "无响应",
	"error.ai_disabled":      "AI未启用",

	// 提示
	"hint.domain_usage":      "请使用 -d 参数指定域名，如: certctl apply -d example.com",
	"hint.domain_format":     "请检查域名格式是否正确（如 example.com）",
	"hint.email_usage":       "邮箱用于 Let's Encrypt 账户注册",
	"hint.email_reminder":    "证书到期前会收到续期提醒邮件",
	"hint.ali_key_param":     "请通过参数 --ali-key 和 --ali-secret 指定",
	"hint.ali_key_env":       "或设置环境变量: ALICLOUD_ACCESS_KEY, ALICLOUD_SECRET_KEY",
	"hint.ali_key_url":       "获取方式: https://ram.console.aliyun.com/manage/ak",
	"hint.check_network":     "请检查网络连接",
	"hint.china_blocked":     "如果在国内，Let's Encrypt 服务器可能被阻断",
	"hint.dns_check":         "请检查 DNS 记录是否正确添加",
	"hint.dns_wait":          "DNS 传播可能需要几分钟，请稍后重试",
	"hint.rate_limit":        "触发了 Let's Encrypt 速率限制，请等待 1 小时后重试",

	// UI 菜单
	"ui.select_operation":    "请选择操作:",
	"ui.apply_cert":          "申请证书",
	"ui.renew_cert":          "证书续期",
	"ui.view_certs":          "查看证书",
	"ui.settings":            "设置",
	"ui.exit":                "✖ 退出",
	"ui.goodbye":             "再见!",

	// 设置
	"ui.settings_title":      "设置",
	"ui.language":            "语言",
	"ui.certs_dir":           "证书目录",
	"ui.dns_config":          "DNS配置",
	"ui.config_dir":          "配置目录",
	"ui.verbose":             "详细模式",
	"ui.on":                  "开启",
	"ui.off":                 "关闭",
	"ui.not_configured":      "未配置",
	"ui.configs_count":       "%d 个配置",
	"ui.select_setting":      "选择要修改的设置:",
	"ui.language_with":       "语言 (%s)",
	"ui.certs_dir_with":      "证书目录(%s)",
	"ui.verbose_with":        "详细模式(%s)",
	"ui.dns_key_config":      "DNS密钥配置",
	"ui.language_set":        "语言",
	"ui.certs_dir_set":       "证书目录",
	"ui.select_certs_folder": "选择证书存放目录",
	"ui.verbose_set":         "详细模式",
	"ui.dns_vendor_config":   "厂商DNS配置",
	"ui.back_main":           "返回主菜单",
	"ui.select_language":     "选择语言:",
	"ui.chinese":             "中文",
	"ui.english":             "English",
	"ui.lang_set_to":         "语言已设置为: %s",
	"ui.verbose_set_to":      "详细模式已%s",

	// DNS 配置
	"ui.dns_config_title":    "厂商DNS配置",
	"ui.no_dns_config":       "暂无已保存的 DNS 配置",
	"ui.select_action":       "选择操作:",
	"ui.add_aliyun":          "添加阿里云 DNS 配置",
	"ui.delete_config":       "删除已保存的配置",
	"ui.back":                "返回",
	"ui.config_name":         "配置名称 (如: 公司账号)",
	"ui.name_empty":          "配置名称不能为空",
	"ui.ali_key_id":          "阿里云 AccessKey ID",
	"ui.ali_key_empty":       "AccessKey ID 不能为空",
	"ui.ali_secret":          "阿里云 AccessKey Secret",
	"ui.ali_secret_empty":    "AccessKey Secret 不能为空",
	"ui.config_saved":        "阿里云DNS配置「%s」已保存",
	"ui.no_config_delete":    "没有可删除的配置",
	"ui.select_delete":       "选择要删除的配置:",
	"ui.cancel":              "取消",
	"ui.deleted":             "已删除「%s」配置",
	"ui.cancelled":           "已取消选择",
	"ui.dir_set_to":          "证书目录已设置为: %s",

	// AI 增强模式
	"ui.ai_config":            "AI增强",
	"ui.ai_config_title":      "AI增强模式",
	"ui.ai_enabled":           "已启用",
	"ui.ai_disabled":          "未启用",
	"ui.toggle_ai":            "启用/禁用 AI",
	"ui.set_api_key":          "设置API Key",
	"ui.select_model":         "选择模型",
	"ui.test_connection":      "测试连接",
	"ui.ai_status":            "当前状态",
	"ui.ai_model":             "模型",
	"ui.ai_key_prompt":        "智谱 AI API Key",
	"ui.ai_key_hint":          "获取地址: https://open.bigmodel.cn",
	"ui.ai_enabled_msg":       "AI增强模式已启用",
	"ui.ai_disabled_msg":      "AI增强模式已禁用",
	"ui.ai_key_saved":         "API Key 已保存",
	"ui.ai_model_saved":       "模型已设置为: %s",
	"ui.ai_testing":           "正在测试连接...",
	"ui.ai_test_ok":           "连接成功！AI 服务正常",
	"ui.ai_test_fail":         "连接失败: %s",
	"ui.ai_diagnosing":        "AI 诊断中...",

	// 申请证书
	"ui.apply_title":         "申请证书",
	"ui.domain":              "域名",
	"ui.domain_empty":        "域名不能为空",
	"ui.email":               "邮箱",
	"ui.email_empty":         "邮箱不能为空",
	"ui.dns_verify_method":   "DNS验证方式:",
	"ui.aliyun_auto":         "阿里云NS(自动)",
	"ui.manual_dns":          "手动添加DNS记录",
	"ui.select_aliyun_cfg":   "选择阿里云配置:",
	"ui.input_new_config":    "输入新的配置",
	"ui.using_config":        "使用配置「%s」",
	"ui.save_config_local":   "是否保存此配置到本地?",
	"ui.config_name_new":     "配置名称 (如: 新账号)",
	"ui.will_apply":          "将申请以下证书:",
	"ui.domain_info":         "域名: %s, *.%s",
	"ui.email_info":          "邮箱: %s",
	"ui.dns_aliyun_auto":     "DNS: 阿里云自动验证",
	"ui.dns_manual":          "DNS: 手动验证",
	"ui.confirm_apply":       "确认申请?",
	"ui.cancelled_op":        "已取消",

	// 续期
	"ui.renew_title":         "证书续期",
	"ui.select_renew_cert":   "选择要续期的证书:",
	"ui.retry_question":      "是否重试?",
	"ui.yes":                 "是",
	"ui.no":                  "否",
	"ui.days":                "天",
	"ui.expired":             "已过期",

	// 查看证书
	"ui.view_title":          "查看证书",
	"ui.no_certs":            "暂无已申请的证书",

	// 其他
	"ui.press_enter":         "按 Enter 键返回主菜单...",
}

// 英文消息
var enMessages = map[string]string{
	// Steps
	"step.prepare":      "Collecting information",
	"step.init":         "Initializing ACME client",
	"step.dns":          "Configuring DNS validation",
	"step.verify":       "Verifying DNS records",
	"step.apply":        "Requesting certificate",
	"step.save":         "Saving certificate",

	// Prompts
	"prompt.domain":     "Domain",
	"prompt.email":      "Email",
	"prompt.ali_key":    "Aliyun AccessKey ID",
	"prompt.ali_secret": "Aliyun AccessKey Secret",
	"prompt.dns_added":  "DNS record added?",

	// Details
	"detail.domain":       "Domain",
	"detail.wildcard":     "Wildcard",
	"detail.email":        "Email",
	"detail.config_dir":   "Config dir",
	"detail.dns_mode":     "DNS mode",
	"detail.dns_aliyun":   "Aliyun auto",
	"detail.dns_manual":   "Manual",
	"detail.env_staging":  "Let's Encrypt Staging",
	"detail.env_prod":     "Let's Encrypt Production",
	"detail.output":       "Output",

	// Progress
	"progress.params_done":   "Parameters collected",
	"progress.account_ok":    "Account loaded",
	"progress.aliyun_ready":     "Aliyun DNS ready",
	"progress.tencentcloud_ready": "Tencent Cloud DNS ready",
	"progress.manual_ready":  "Manual mode ready",
	"progress.client_ready":  "ACME client ready",
	"progress.dns_ok":        "DNS record verified",
	"progress.cert_ok":       "Certificate issued",
	"progress.saved":         "Certificate saved",
	"progress.checking_dns":  "Checking DNS propagation...",
	"progress.applying":      "Adding DNS records and requesting certificate...",

	// Dry run
	"dryrun.warning":    "Dry run mode: no certificate will be issued",
	"dryrun.skip_apply": "Dry run: skipping certificate request",
	"dryrun.skip_save":  "Dry run: skipping certificate save",
	"dryrun.done":       "Dry run complete! All parameters validated",

	// Errors
	"error.domain_empty":     "Domain cannot be empty",
	"error.domain_invalid":   "Invalid domain format",
	"error.domain_parse":     "Failed to parse domain: %v",
	"error.email_empty":      "Email cannot be empty",
	"error.ali_key_empty":    "Aliyun AccessKey incomplete",
	"error.account_fail":     "Account initialization failed",
	"error.aliyun_fail":         "Aliyun DNS client creation failed",
	"error.aliyun_create":       "Failed to create Aliyun DNS client: %v",
	"error.tencentcloud_fail":   "Tencent Cloud DNS client creation failed",
	"error.tencentcloud_create": "Failed to create Tencent Cloud DNS client: %v",
	"error.client_fail":      "ACME client creation failed",
	"error.client_create":    "Failed to create ACME client: %v",
	"error.dns_provider":     "Failed to set DNS provider: %v",
	"error.register_fail":    "Account registration failed",
	"error.register":         "Failed to register account: %v",
	"error.dns_fail":         "DNS validation failed",
	"error.dns_add":          "Failed to add DNS record: %v",
	"error.dns_update":       "Failed to update DNS record: %v",
	"error.dns_delete":       "Failed to delete DNS record: %v",
	"error.dns_query":        "Failed to query DNS record: %v",
	"error.dns_query_fail":   "DNS query failed: %s",
	"error.dns_timeout":      "DNS record verification timeout",
	"error.cert_fail":        "Certificate request failed",
	"error.cert_obtain":      "Failed to obtain certificate: %v",
	"error.cert_parse":       "Failed to parse certificate",
	"error.cert_parse_err":   "Failed to parse certificate: %v",
	"error.save_fail":        "Certificate save failed",
	"error.user_cancel":      "User cancelled",
	"error.ai_request":       "AI request failed: %v",
	"error.ai_parse":         "AI response parse failed: %v",
	"error.ai_error":         "AI error: %s",
	"error.ai_no_response":   "AI no response",
	"error.ai_connect":       "Connection failed: %v",
	"error.api_error":        "API error: %s",
	"error.no_response":      "No response",
	"error.ai_disabled":      "AI not enabled",

	// Hints
	"hint.domain_usage":      "Use -d to specify domain, e.g.: certctl apply -d example.com",
	"hint.domain_format":     "Check domain format (e.g. example.com)",
	"hint.email_usage":       "Email is used for Let's Encrypt account registration",
	"hint.email_reminder":    "You will receive renewal reminders before expiry",
	"hint.ali_key_param":     "Use --ali-key and --ali-secret parameters",
	"hint.ali_key_env":       "Or set env: ALICLOUD_ACCESS_KEY, ALICLOUD_SECRET_KEY",
	"hint.ali_key_url":       "Get keys at: https://ram.console.aliyun.com/manage/ak",
	"hint.check_network":     "Check your network connection",
	"hint.china_blocked":     "If in China, Let's Encrypt servers may be blocked",
	"hint.dns_check":         "Verify DNS record is correctly added",
	"hint.dns_wait":          "DNS propagation may take a few minutes",
	"hint.rate_limit":        "Rate limit hit, please wait 1 hour and retry",

	// UI Menu
	"ui.select_operation":    "Select operation:",
	"ui.apply_cert":          "Apply Certificate",
	"ui.renew_cert":          "Certificate Renewal",
	"ui.view_certs":          "View Certificates",
	"ui.settings":            "Settings",
	"ui.exit":                "✖ Exit",
	"ui.goodbye":             "Goodbye!",

	// Settings
	"ui.settings_title":      "Settings",
	"ui.language":            "Language",
	"ui.certs_dir":           "Certs Dir",
	"ui.dns_config":          "DNS Config",
	"ui.config_dir":          "Config Dir",
	"ui.verbose":             "Verbose Mode",
	"ui.on":                  "On",
	"ui.off":                 "Off",
	"ui.not_configured":      "Not configured",
	"ui.configs_count":       "%d configs",
	"ui.select_setting":      "Select setting to modify:",
	"ui.language_with":       "Language (%s)",
	"ui.certs_dir_with":      "Certs Dir (%s)",
	"ui.verbose_with":        "Verbose Mode (%s)",
	"ui.dns_key_config":      "DNS Key Config",
	"ui.language_set":        "Language",
	"ui.certs_dir_set":       "Certs Directory",
	"ui.select_certs_folder": "Select certificate folder",
	"ui.verbose_set":         "Verbose Mode",
	"ui.dns_vendor_config":   "Vendor DNS Config",
	"ui.back_main":           "Back to Main Menu",
	"ui.select_language":     "Select language:",
	"ui.chinese":             "中文",
	"ui.english":             "English",
	"ui.lang_set_to":         "Language set to: %s",
	"ui.verbose_set_to":      "Verbose mode %s",

	// DNS Config
	"ui.dns_config_title":    "Vendor DNS Config",
	"ui.no_dns_config":       "No saved DNS configs",
	"ui.select_action":       "Select action:",
	"ui.add_aliyun":          "Add Aliyun DNS Config",
	"ui.delete_config":       "Delete Saved Config",
	"ui.back":                "Back",
	"ui.config_name":         "Config Name (e.g. Company)",
	"ui.name_empty":          "Config name cannot be empty",
	"ui.ali_key_id":          "Aliyun AccessKey ID",
	"ui.ali_key_empty":       "AccessKey ID cannot be empty",
	"ui.ali_secret":          "Aliyun AccessKey Secret",
	"ui.ali_secret_empty":    "AccessKey Secret cannot be empty",
	"ui.config_saved":        "Aliyun DNS config [%s] saved",
	"ui.no_config_delete":    "No config to delete",
	"ui.select_delete":       "Select config to delete:",
	"ui.cancel":              "Cancel",
	"ui.deleted":             "Deleted [%s] config",
	"ui.cancelled":           "Selection cancelled",
	"ui.dir_set_to":          "Certs dir set to: %s",

	// AI Enhancement
	"ui.ai_config":            "AI Enhancement",
	"ui.ai_config_title":      "AI Enhancement Mode",
	"ui.ai_enabled":           "Enabled",
	"ui.ai_disabled":          "Disabled",
	"ui.toggle_ai":            "Enable/Disable AI",
	"ui.set_api_key":          "Set API Key",
	"ui.select_model":         "Select Model",
	"ui.test_connection":      "Test Connection",
	"ui.ai_status":            "Status",
	"ui.ai_model":             "Model",
	"ui.ai_key_prompt":        "Zhipu AI API Key",
	"ui.ai_key_hint":          "Get at: https://open.bigmodel.cn",
	"ui.ai_enabled_msg":       "AI Enhancement enabled",
	"ui.ai_disabled_msg":      "AI Enhancement disabled",
	"ui.ai_key_saved":         "API Key saved",
	"ui.ai_model_saved":       "Model set to: %s",
	"ui.ai_testing":           "Testing connection...",
	"ui.ai_test_ok":           "Connection successful! AI service is working",
	"ui.ai_test_fail":         "Connection failed: %s",
	"ui.ai_diagnosing":        "AI diagnosing...",

	// Apply Certificate
	"ui.apply_title":         "Apply Certificate",
	"ui.domain":              "Domain",
	"ui.domain_empty":        "Domain cannot be empty",
	"ui.email":               "Email",
	"ui.email_empty":         "Email cannot be empty",
	"ui.dns_verify_method":   "DNS Verification Method:",
	"ui.aliyun_auto":         "Aliyun DNS (Auto)",
	"ui.manual_dns":          "Manual DNS Record",
	"ui.select_aliyun_cfg":   "Select Aliyun config:",
	"ui.input_new_config":    "Input new config",
	"ui.using_config":        "Using config [%s]",
	"ui.save_config_local":   "Save this config locally?",
	"ui.config_name_new":     "Config Name (e.g. New Account)",
	"ui.will_apply":          "Will apply for:",
	"ui.domain_info":         "Domain: %s, *.%s",
	"ui.email_info":          "Email: %s",
	"ui.dns_aliyun_auto":     "DNS: Aliyun Auto",
	"ui.dns_manual":          "DNS: Manual",
	"ui.confirm_apply":       "Confirm apply?",
	"ui.cancelled_op":        "Cancelled",

	// Renew
	"ui.renew_title":         "Certificate Renewal",
	"ui.select_renew_cert":   "Select certificate to renew:",
	"ui.retry_question":      "Retry?",
	"ui.yes":                 "Yes",
	"ui.no":                  "No",
	"ui.days":                "days",
	"ui.expired":             "Expired",

	// View Certificates
	"ui.view_title":          "View Certificates",
	"ui.no_certs":            "No certificates applied",

	// Other
	"ui.press_enter":         "Press Enter to return...",
}
