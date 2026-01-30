package acme

import (
	"certctl/internal/dns/tencentcloud"
	"certctl/internal/i18n"
	"certctl/pkg/domain"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/go-acme/lego/v4/challenge"
)

// TencentCloudDNSProvider 腾讯云 DNS 验证提供者
type TencentCloudDNSProvider struct {
	client *tencentcloud.DNSClient
}

// NewTencentCloudDNSProvider 创建腾讯云 DNS 提供者
func NewTencentCloudDNSProvider(secretId, secretKey, region string) (*TencentCloudDNSProvider, error) {
	client, err := tencentcloud.NewDNSClient(secretId, secretKey, region)
	if err != nil {
		return nil, err
	}

	return &TencentCloudDNSProvider{client: client}, nil
}

func (p *TencentCloudDNSProvider) Present(domainName, token, keyAuth string) error {
	// 计算 TXT 记录值
	hash := sha256.Sum256([]byte(keyAuth))
	txtValue := base64.RawURLEncoding.EncodeToString(hash[:])

	// 解析根域名
	rootDomain, err := domain.Parse(domainName)
	if err != nil {
		return fmt.Errorf(i18n.T("error.domain_parse"), err)
	}

	// 添加 TXT 记录
	rr := "_acme-challenge"
	if domainName != rootDomain {
		// 子域名情况
		rr = "_acme-challenge." + domainName[:len(domainName)-len(rootDomain)-1]
	}

	return p.client.AddTXTRecord(rootDomain, rr, txtValue)
}

func (p *TencentCloudDNSProvider) CleanUp(domainName, token, keyAuth string) error {
	// 解析根域名
	rootDomain, err := domain.Parse(domainName)
	if err != nil {
		return nil // 清理时忽略错误
	}

	rr := "_acme-challenge"
	if domainName != rootDomain {
		rr = "_acme-challenge." + domainName[:len(domainName)-len(rootDomain)-1]
	}

	return p.client.DeleteTXTRecord(rootDomain, rr)
}

// 确保实现了接口
var _ challenge.Provider = (*TencentCloudDNSProvider)(nil)
