package acme

import "time"

// Certificate 证书
type Certificate struct {
	Domain      string
	Certificate []byte
	PrivateKey  []byte
	NotAfter    time.Time
}
