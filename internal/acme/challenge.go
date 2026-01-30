package acme

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/dns01"
)

// Challenge DNS 验证挑战信息
type Challenge struct {
	Domain     string // 原始域名
	FQDN       string // 完整记录名 _acme-challenge.example.com
	RecordName string // 主机记录 _acme-challenge
	Value      string // TXT 记录值
}

// ManualDNSProvider 手动 DNS 验证提供者
type ManualDNSProvider struct {
	challenges map[string]*Challenge
	mu         sync.Mutex
	onPresent  func(*Challenge) error // 阻塞回调：显示信息、等待用户确认、检查 DNS
	onCleanup  func(*Challenge) error
}

// NewManualDNSProvider 创建手动 DNS 提供者
// onPresent 回调应该阻塞直到 DNS 记录验证通过
func NewManualDNSProvider(onPresent, onCleanup func(*Challenge) error) *ManualDNSProvider {
	return &ManualDNSProvider{
		challenges: make(map[string]*Challenge),
		onPresent:  onPresent,
		onCleanup:  onCleanup,
	}
}

func (p *ManualDNSProvider) Present(domain, token, keyAuth string) error {
	// 计算 TXT 记录值
	hash := sha256.Sum256([]byte(keyAuth))
	txtValue := base64.RawURLEncoding.EncodeToString(hash[:])

	fqdn := fmt.Sprintf("_acme-challenge.%s", domain)

	challenge := &Challenge{
		Domain:     domain,
		FQDN:       fqdn,
		RecordName: "_acme-challenge",
		Value:      txtValue,
	}

	p.mu.Lock()
	p.challenges[domain] = challenge
	p.mu.Unlock()

	// onPresent 回调会阻塞直到用户确认且 DNS 记录验证通过
	if p.onPresent != nil {
		return p.onPresent(challenge)
	}

	return nil
}

func (p *ManualDNSProvider) CleanUp(domain, token, keyAuth string) error {
	p.mu.Lock()
	challenge := p.challenges[domain]
	delete(p.challenges, domain)
	p.mu.Unlock()

	if p.onCleanup != nil && challenge != nil {
		return p.onCleanup(challenge)
	}

	return nil
}

// GetChallenge 获取挑战信息
func (p *ManualDNSProvider) GetChallenge(domain string) *Challenge {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.challenges[domain]
}

// 确保实现了接口
var _ challenge.Provider = (*ManualDNSProvider)(nil)

// GetChallengeInfo 从 lego 获取挑战信息（用于调试）
func GetChallengeInfo(domain, keyAuth string) (fqdn, value string) {
	fqdn, value = dns01.GetRecord(domain, keyAuth)
	return
}
