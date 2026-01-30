package aliyun

import (
	"fmt"
	"strings"

	"certctl/internal/i18n"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

// DNSClient 阿里云 DNS 客户端
type DNSClient struct {
	client *alidns.Client
}

// NewDNSClient 创建阿里云 DNS 客户端
func NewDNSClient(accessKey, accessSecret, region string) (*DNSClient, error) {
	if region == "" {
		region = "cn-hangzhou"
	}

	client, err := alidns.NewClientWithAccessKey(region, accessKey, accessSecret)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.aliyun_create"), err)
	}

	return &DNSClient{client: client}, nil
}

// AddTXTRecord 添加 TXT 记录
func (c *DNSClient) AddTXTRecord(domain, rr, value string) error {
	request := alidns.CreateAddDomainRecordRequest()
	request.Scheme = "https"
	request.DomainName = domain
	request.RR = rr
	request.Type = "TXT"
	request.Value = value
	request.TTL = requests.NewInteger(600)

	_, err := c.client.AddDomainRecord(request)
	if err != nil {
		// 如果记录已存在，尝试更新
		if strings.Contains(err.Error(), "DomainRecordDuplicate") {
			return c.UpdateTXTRecord(domain, rr, value)
		}
		return fmt.Errorf(i18n.T("error.dns_add"), err)
	}

	return nil
}

// UpdateTXTRecord 更新 TXT 记录
func (c *DNSClient) UpdateTXTRecord(domain, rr, value string) error {
	// 先查询记录 ID
	recordID, err := c.getRecordID(domain, rr, "TXT")
	if err != nil {
		return err
	}

	if recordID == "" {
		// 记录不存在，添加新记录
		return c.AddTXTRecord(domain, rr, value)
	}

	// 更新记录
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = recordID
	request.RR = rr
	request.Type = "TXT"
	request.Value = value

	_, err = c.client.UpdateDomainRecord(request)
	if err != nil {
		return fmt.Errorf(i18n.T("error.dns_update"), err)
	}

	return nil
}

// DeleteTXTRecord 删除 TXT 记录
func (c *DNSClient) DeleteTXTRecord(domain, rr string) error {
	recordID, err := c.getRecordID(domain, rr, "TXT")
	if err != nil {
		return err
	}

	if recordID == "" {
		return nil // 记录不存在，无需删除
	}

	request := alidns.CreateDeleteDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = recordID

	_, err = c.client.DeleteDomainRecord(request)
	if err != nil {
		return fmt.Errorf(i18n.T("error.dns_delete"), err)
	}

	return nil
}

// getRecordID 获取记录 ID
func (c *DNSClient) getRecordID(domain, rr, recordType string) (string, error) {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.DomainName = domain
	request.RRKeyWord = rr
	request.TypeKeyWord = recordType

	response, err := c.client.DescribeDomainRecords(request)
	if err != nil {
		return "", fmt.Errorf(i18n.T("error.dns_query"), err)
	}

	for _, record := range response.DomainRecords.Record {
		if record.RR == rr && record.Type == recordType {
			return record.RecordId, nil
		}
	}

	return "", nil
}
