package cert

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"time"
)

// Save 保存证书到文件
func Save(outputDir, domain string, certPEM, keyPEM []byte) (certPath, keyPath string, err error) {
	domainDir := filepath.Join(outputDir, domain)
	if err = os.MkdirAll(domainDir, 0755); err != nil {
		return
	}

	certPath = filepath.Join(domainDir, domain+".pem")
	keyPath = filepath.Join(domainDir, domain+".key")

	if err = os.WriteFile(certPath, certPEM, 0644); err != nil {
		return
	}

	if err = os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return
	}

	return certPath, keyPath, nil
}

// Certificate 证书信息
type Certificate struct {
	CertPath string
	KeyPath  string
	Domain   string
	NotAfter time.Time
	DaysLeft int
}

// ListCertificates 扫描目录下的所有证书
func ListCertificates(certsDir string) ([]Certificate, error) {
	var certs []Certificate

	// 检查目录是否存在
	if _, err := os.Stat(certsDir); os.IsNotExist(err) {
		return certs, nil
	}

	// 遍历证书目录
	entries, err := os.ReadDir(certsDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		domain := entry.Name()
		domainDir := filepath.Join(certsDir, domain)

		// 支持多种证书命名格式
		certNames := []string{domain + ".pem", "fullchain.pem", "cert.pem", "certificate.pem"}
		keyNames := []string{domain + ".key", "privkey.pem", "key.pem", "private.key"}

		var certPath, keyPath string

		// 查找证书文件
		for _, name := range certNames {
			path := filepath.Join(domainDir, name)
			if _, err := os.Stat(path); err == nil {
				certPath = path
				break
			}
		}

		// 查找私钥文件
		for _, name := range keyNames {
			path := filepath.Join(domainDir, name)
			if _, err := os.Stat(path); err == nil {
				keyPath = path
				break
			}
		}

		// 如果没找到证书，跳过
		if certPath == "" {
			continue
		}

		// 解析证书获取有效期
		notAfter, err := ParseCertExpiry(certPath)
		if err != nil {
			continue
		}

		daysLeft := int(time.Until(notAfter).Hours() / 24)

		certs = append(certs, Certificate{
			CertPath: certPath,
			KeyPath:  keyPath,
			Domain:   domain,
			NotAfter: notAfter,
			DaysLeft: daysLeft,
		})
	}

	return certs, nil
}

// ParseCertExpiry 解析证书有效期
func ParseCertExpiry(certPath string) (time.Time, error) {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return time.Time{}, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return time.Time{}, err
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Time{}, err
	}

	return cert.NotAfter, nil
}
