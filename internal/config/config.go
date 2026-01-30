package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// DNSConfig DNS 提供商配置（支持命名的多个配置）
type DNSConfig struct {
	Name            string `json:"name"`            // 配置名称，如 "阿里云-公司账号"
	Provider        string `json:"provider"`        // 提供商类型：aliyun, tencent, cloudflare
	AccessKeyID     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
}

// AIConfig AI 增强模式配置
type AIConfig struct {
	Enabled bool   `json:"enabled"` // 是否启用
	APIKey  string `json:"apiKey"`  // 智谱 AI API Key
	Model   string `json:"model"`   // 模型，默认 glm-4.7
}

// Config 应用配置
type Config struct {
	Language string      `json:"language"`
	CertsDir string      `json:"certsDir"`
	Verbose  bool        `json:"verbose"` // 详细模式
	DNS      []DNSConfig `json:"dns"`     // 改为数组，支持多个配置
	AI       AIConfig    `json:"ai"`      // AI 增强模式
}

var (
	configDir  string
	configFile string
	current    *Config
)

func init() {
	home, _ := os.UserHomeDir()
	configDir = filepath.Join(home, ".certctl")
	configFile = filepath.Join(configDir, "config.json")
}

// GetConfigDir 获取配置目录
func GetConfigDir() string {
	return configDir
}

// Load 加载配置
func Load() *Config {
	if current != nil {
		return current
	}

	// 默认证书目录使用用户主目录下的 certs
	home, _ := os.UserHomeDir()
	defaultCertsDir := filepath.Join(home, "certs")

	current = &Config{
		Language: "zh",
		CertsDir: defaultCertsDir,
		DNS:      []DNSConfig{},
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return current
	}

	json.Unmarshal(data, current)
	if current.DNS == nil {
		current.DNS = []DNSConfig{}
	}
	return current
}

// Save 保存配置
func Save() error {
	// 确保目录存在
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件，权限 600（仅用户可读写）
	return os.WriteFile(configFile, data, 0600)
}

// Get 获取当前配置
func Get() *Config {
	if current == nil {
		Load()
	}
	return current
}

// SetLanguage 设置语言
func SetLanguage(lang string) {
	Get().Language = lang
	Save()
}

// SetCertsDir 设置证书目录
func SetCertsDir(dir string) {
	Get().CertsDir = dir
	Save()
}

// SetVerbose 设置详细模式
func SetVerbose(verbose bool) {
	Get().Verbose = verbose
	Save()
}

// GetVerbose 获取详细模式
func GetVerbose() bool {
	return Get().Verbose
}

// GetDNSConfigs 获取所有 DNS 配置
func GetDNSConfigs() []DNSConfig {
	return Get().DNS
}

// GetDNSConfigsByProvider 获取指定提供商的所有配置
func GetDNSConfigsByProvider(provider string) []DNSConfig {
	var configs []DNSConfig
	for _, dns := range Get().DNS {
		if dns.Provider == provider {
			configs = append(configs, dns)
		}
	}
	return configs
}

// GetDNSConfigByName 根据名称获取配置
func GetDNSConfigByName(name string) (DNSConfig, bool) {
	for _, dns := range Get().DNS {
		if dns.Name == name {
			return dns, true
		}
	}
	return DNSConfig{}, false
}

// AddDNSConfig 添加 DNS 配置
func AddDNSConfig(name, provider, accessKeyID, accessKeySecret string) {
	cfg := Get()
	// 如果已存在同名配置，先删除
	DeleteDNSConfig(name)
	cfg.DNS = append(cfg.DNS, DNSConfig{
		Name:            name,
		Provider:        provider,
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
	})
	Save()
}

// DeleteDNSConfig 根据名称删除 DNS 配置
func DeleteDNSConfig(name string) {
	cfg := Get()
	newDNS := []DNSConfig{}
	for _, dns := range cfg.DNS {
		if dns.Name != name {
			newDNS = append(newDNS, dns)
		}
	}
	cfg.DNS = newDNS
	Save()
}

// HasDNSConfigs 检查是否有任何 DNS 配置
func HasDNSConfigs() bool {
	return len(Get().DNS) > 0
}

// HasProviderConfigs 检查是否有指定提供商的配置
func HasProviderConfigs(provider string) bool {
	return len(GetDNSConfigsByProvider(provider)) > 0
}

// GetAIConfig 获取 AI 配置
func GetAIConfig() AIConfig {
	cfg := Get()
	if cfg.AI.Model == "" {
		cfg.AI.Model = "glm-4-flash" // 默认模型（有免费额度）
	}
	return cfg.AI
}

// SetAIEnabled 设置 AI 启用状态
func SetAIEnabled(enabled bool) {
	cfg := Get()
	cfg.AI.Enabled = enabled
	Save()
}

// SetAIAPIKey 设置 AI API Key
func SetAIAPIKey(apiKey string) {
	cfg := Get()
	cfg.AI.APIKey = apiKey
	Save()
}

// SetAIModel 设置 AI 模型
func SetAIModel(model string) {
	cfg := Get()
	cfg.AI.Model = model
	Save()
}

// IsAIEnabled 检查 AI 是否启用且有效
func IsAIEnabled() bool {
	cfg := GetAIConfig()
	return cfg.Enabled && cfg.APIKey != ""
}
