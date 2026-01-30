package acme

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/registration"
)

const accountFileName = "account.json"
const keyFileName = "account.key"

// Account ACME 账户
type Account struct {
	Email        string                 `json:"email"`
	Registration *registration.Resource `json:"registration"`
	KeyPath      string                 `json:"key_path"`
	key          crypto.PrivateKey
}

func (a *Account) GetEmail() string {
	return a.Email
}

func (a *Account) GetRegistration() *registration.Resource {
	return a.Registration
}

func (a *Account) GetPrivateKey() crypto.PrivateKey {
	return a.key
}

// LoadOrCreateAccount 加载或创建账户
func LoadOrCreateAccount(configDir, email string) (*Account, error) {
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, err
	}

	accountPath := filepath.Join(configDir, accountFileName)
	keyPath := filepath.Join(configDir, keyFileName)

	// 尝试加载已有账户
	if _, err := os.Stat(accountPath); err == nil {
		return loadAccount(accountPath, keyPath)
	}

	// 创建新账户
	return createAccount(configDir, email)
}

func loadAccount(accountPath, keyPath string) (*Account, error) {
	data, err := os.ReadFile(accountPath)
	if err != nil {
		return nil, err
	}

	var account Account
	if err := json.Unmarshal(data, &account); err != nil {
		return nil, err
	}

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	account.key = key
	return &account, nil
}

func createAccount(configDir, email string) (*Account, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	keyPath := filepath.Join(configDir, keyFileName)
	if err := savePrivateKey(keyPath, privateKey); err != nil {
		return nil, err
	}

	account := &Account{
		Email:   email,
		KeyPath: keyPath,
		key:     privateKey,
	}

	return account, nil
}

func savePrivateKey(path string, key *ecdsa.PrivateKey) error {
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	}

	return os.WriteFile(path, pem.EncodeToMemory(block), 0600)
}

// SaveAccount 保存账户信息
func SaveAccount(configDir string, account *Account) error {
	data, err := json.MarshalIndent(account, "", "  ")
	if err != nil {
		return err
	}

	accountPath := filepath.Join(configDir, accountFileName)
	return os.WriteFile(accountPath, data, 0600)
}
