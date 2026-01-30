package dns

import (
	"fmt"
	"net"
	"strings"
	"time"

	"certctl/internal/i18n"

	"github.com/miekg/dns"
)

var defaultResolvers = []string{
	"8.8.8.8:53",
	"1.1.1.1:53",
	"223.5.5.5:53",
}

// CheckTXTRecord 检查 TXT 记录是否存在
func CheckTXTRecord(fqdn, expectedValue string) (bool, error) {
	if !strings.HasSuffix(fqdn, ".") {
		fqdn = fqdn + "."
	}

	for _, resolver := range defaultResolvers {
		values, err := queryTXT(fqdn, resolver)
		if err != nil {
			continue
		}

		for _, v := range values {
			if v == expectedValue {
				return true, nil
			}
		}
	}

	return false, nil
}

func queryTXT(fqdn, resolver string) ([]string, error) {
	c := new(dns.Client)
	c.Timeout = 5 * time.Second

	m := new(dns.Msg)
	m.SetQuestion(fqdn, dns.TypeTXT)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, resolver)
	if err != nil {
		return nil, err
	}

	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf(i18n.T("error.dns_query_fail"), dns.RcodeToString[r.Rcode])
	}

	var values []string
	for _, ans := range r.Answer {
		if txt, ok := ans.(*dns.TXT); ok {
			values = append(values, strings.Join(txt.Txt, ""))
		}
	}

	return values, nil
}

// WaitForRecord 等待 DNS 记录生效
func WaitForRecord(fqdn, expectedValue string, timeout time.Duration, onCheck func(attempt int)) error {
	deadline := time.Now().Add(timeout)
	attempt := 0

	for time.Now().Before(deadline) {
		attempt++
		if onCheck != nil {
			onCheck(attempt)
		}

		found, err := CheckTXTRecord(fqdn, expectedValue)
		if err == nil && found {
			return nil
		}

		time.Sleep(10 * time.Second)
	}

	return fmt.Errorf(i18n.T("error.dns_timeout"))
}

// GetLocalIP 获取本机出口 IP
func GetLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
