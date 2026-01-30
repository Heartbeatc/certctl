package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"certctl/internal/config"
	"certctl/internal/i18n"
)

const ZhipuAPIURL = "https://open.bigmodel.cn/api/paas/v4/chat/completions"

// ZhipuClient 智谱 AI 客户端
type ZhipuClient struct {
	APIKey string
	Model  string
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ZhipuRequest 请求结构
type ZhipuRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

// ZhipuResponse 响应结构
type ZhipuResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"error"`
}

// NewZhipuClient 创建智谱 AI 客户端
func NewZhipuClient() *ZhipuClient {
	cfg := config.GetAIConfig()
	model := cfg.Model
	if model == "" {
		model = "glm-4-flash"
	}
	return &ZhipuClient{
		APIKey: cfg.APIKey,
		Model:  model,
	}
}

// Diagnose 诊断证书申请错误
func (c *ZhipuClient) Diagnose(errorMsg, domain, dnsProvider string) (string, error) {
	prompt := buildDiagnosisPrompt(errorMsg, domain, dnsProvider)

	req := ZhipuRequest{
		Model: c.Model,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		MaxTokens: 500,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequest("POST", ZhipuAPIURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf(i18n.T("error.ai_request"), err)
	}
	defer resp.Body.Close()

	var result ZhipuResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf(i18n.T("error.ai_parse"), err)
	}

	if result.Error != nil {
		return "", fmt.Errorf(i18n.T("error.ai_error"), result.Error.Message)
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf(i18n.T("error.ai_no_response"))
}

// TestConnection 测试连接
func (c *ZhipuClient) TestConnection() error {
	req := ZhipuRequest{
		Model: c.Model,
		Messages: []Message{
			{Role: "user", Content: "你好，请回复 OK"},
		},
		MaxTokens: 10,
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", ZhipuAPIURL, bytes.NewBuffer(body))
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf(i18n.T("error.ai_connect"), err)
	}
	defer resp.Body.Close()

	var result ZhipuResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Error != nil {
		return fmt.Errorf(i18n.T("error.api_error"), result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return fmt.Errorf(i18n.T("error.no_response"))
	}

	return nil
}

// buildDiagnosisPrompt 构建诊断提示词
func buildDiagnosisPrompt(errorMsg, domain, dnsProvider string) string {
	lang := i18n.Lang
	
	if lang == "en" {
		return fmt.Sprintf(`You are an SSL certificate expert. Please analyze the following Let's Encrypt certificate error and provide solutions.

Error message:
%s

Domain: %s
DNS Provider: %s

Please answer concisely in English with the following format:

🔍 Problem: (one sentence describing the issue)

✅ Solutions:
1. xxx
2. xxx

💡 Retry recommended: Yes/No`, errorMsg, domain, dnsProvider)
	}
	
	// 默认中文
	return fmt.Sprintf(`你是一个 SSL 证书申请专家。请分析以下 Let's Encrypt 证书申请错误并给出解决方案。

错误信息:
%s

域名: %s
DNS 提供商: %s

请用简洁的中文回答，格式如下：

🔍 问题原因：（一句话描述问题）

✅ 解决方案：
1. xxx
2. xxx
3. xxx
..........

💡 是否建议重试：是/否`, errorMsg, domain, dnsProvider)
}

// DiagnoseError 便捷函数：诊断错误
func DiagnoseError(errorMsg, domain, dnsProvider string) (string, error) {
	if !config.IsAIEnabled() {
		return "", fmt.Errorf(i18n.T("error.ai_disabled"))
	}
	client := NewZhipuClient()
	return client.Diagnose(errorMsg, domain, dnsProvider)
}
