package domain

import (
	"errors"
	"strings"
)

var ErrInvalidDomain = errors.New("无效的域名格式")

// 常见的二级域名后缀
var secondLevelTLDs = map[string]bool{
	"com.cn": true, "net.cn": true, "org.cn": true, "gov.cn": true,
	"co.uk": true, "org.uk": true, "ac.uk": true,
	"co.jp": true, "ne.jp": true, "or.jp": true,
	"com.au": true, "net.au": true, "org.au": true,
	"co.kr": true, "or.kr": true,
	"com.hk": true, "org.hk": true,
	"com.tw": true, "org.tw": true,
}

// Parse 解析域名，返回根域名
func Parse(domain string) (string, error) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "*.")
	domain = strings.TrimSuffix(domain, "/")

	// 移除路径部分
	if idx := strings.Index(domain, "/"); idx != -1 {
		domain = domain[:idx]
	}

	// 移除端口
	if idx := strings.Index(domain, ":"); idx != -1 {
		domain = domain[:idx]
	}

	if domain == "" {
		return "", ErrInvalidDomain
	}

	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return "", ErrInvalidDomain
	}

	// 检查是否是二级域名后缀
	if len(parts) >= 3 {
		possibleTLD := parts[len(parts)-2] + "." + parts[len(parts)-1]
		if secondLevelTLDs[possibleTLD] {
			if len(parts) < 3 {
				return "", ErrInvalidDomain
			}
			return parts[len(parts)-3] + "." + possibleTLD, nil
		}
	}

	// 普通域名，返回最后两部分
	return parts[len(parts)-2] + "." + parts[len(parts)-1], nil
}

// GenerateWildcard 生成通配符域名列表
func GenerateWildcard(domain string) ([]string, error) {
	root, err := Parse(domain)
	if err != nil {
		return nil, err
	}

	return []string{root, "*." + root}, nil
}
