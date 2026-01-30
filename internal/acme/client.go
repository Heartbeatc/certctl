package acme

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"certctl/internal/i18n"

	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

const (
	// LetsEncryptProduction 生产环境
	LetsEncryptProduction = "https://acme-v02.api.letsencrypt.org/directory"
	// LetsEncryptStaging 测试环境
	LetsEncryptStaging = "https://acme-staging-v02.api.letsencrypt.org/directory"
)

// Client ACME 客户端
type Client struct {
	account  *Account
	client   *lego.Client
	provider challenge.Provider
}

// NewClient 创建 ACME 客户端
func NewClient(account *Account, staging bool, provider challenge.Provider) (*Client, error) {
	config := lego.NewConfig(account)

	if staging {
		config.CADirURL = LetsEncryptStaging
	} else {
		config.CADirURL = LetsEncryptProduction
	}

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.client_create"), err)
	}

	// 设置 DNS 验证，禁用 lego 自带的传播检查（我们自己检查）
	if err := client.Challenge.SetDNS01Provider(
		provider,
		dns01.AddRecursiveNameservers([]string{"8.8.8.8:53", "1.1.1.1:53"}),
		dns01.DisableCompletePropagationRequirement(),
	); err != nil {
		return nil, fmt.Errorf(i18n.T("error.dns_provider"), err)
	}

	return &Client{
		account:  account,
		client:   client,
		provider: provider,
	}, nil
}

// Register 注册账户
func (c *Client) Register() error {
	if c.account.Registration != nil {
		return nil
	}

	reg, err := c.client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return fmt.Errorf(i18n.T("error.register"), err)
	}

	c.account.Registration = reg
	return nil
}

// ObtainCertificate 申请证书
func (c *Client) ObtainCertificate(domains []string) (*Certificate, error) {
	request := certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	}

	certificates, err := c.client.Certificate.Obtain(request)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.cert_obtain"), err)
	}

	// 解析证书获取过期时间
	block, _ := pem.Decode(certificates.Certificate)
	if block == nil {
		return nil, fmt.Errorf(i18n.T("error.cert_parse"))
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.cert_parse_err"), err)
	}

	return &Certificate{
		Domain:      certificates.Domain,
		Certificate: certificates.Certificate,
		PrivateKey:  certificates.PrivateKey,
		NotAfter:    cert.NotAfter,
	}, nil
}

// GetProvider 获取 DNS 提供者
func (c *Client) GetProvider() challenge.Provider {
	return c.provider
}
