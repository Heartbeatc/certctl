package tencentcloud

import (
	"fmt"
	"strings"

	"certctl/internal/i18n"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
)

// DNSClient 腾讯云 DNS 客户端
type DNSClient struct {
	client *dnspod.Client
}

// NewDNSClient 创建腾讯云 DNS 客户端
func NewDNSClient(secretId, secretKey, region string) (*DNSClient, error) {
	if region == "" {
		region = "ap-guangzhou"
	}

	credential := common.NewCredential(secretId, secretKey)
	cpf := profile.NewClientProfile()
	client, err := dnspod.NewClient(credential, region, cpf)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.tencentcloud_create"), err)
	}

	return &DNSClient{client: client}, nil
}

// AddTXTRecord 添加 TXT 记录
func (c *DNSClient) AddTXTRecord(domain, rr, value string) error {
	request := dnspod.NewCreateRecordRequest()
	request.Domain = common.StringPtr(domain)
	request.SubDomain = common.StringPtr(rr)
	request.RecordType = common.StringPtr("TXT")
	request.RecordLine = common.StringPtr("默认")
	request.Value = common.StringPtr(value)

	_, err := c.client.CreateRecord(request)
	if err != nil {
		// 如果记录已存在，尝试更新
		if strings.Contains(err.Error(), "RecordAlreadyExists") || strings.Contains(err.Error(), "记录已存在") {
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

	if recordID == 0 {
		// 记录不存在，添加新记录
		return c.AddTXTRecord(domain, rr, value)
	}

	// 更新记录
	request := dnspod.NewModifyRecordRequest()
	request.Domain = common.StringPtr(domain)
	request.RecordId = common.Uint64Ptr(recordID)
	request.SubDomain = common.StringPtr(rr)
	request.RecordType = common.StringPtr("TXT")
	request.RecordLine = common.StringPtr("默认")
	request.Value = common.StringPtr(value)

	_, err = c.client.ModifyRecord(request)
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

	if recordID == 0 {
		return nil // 记录不存在，无需删除
	}

	request := dnspod.NewDeleteRecordRequest()
	request.Domain = common.StringPtr(domain)
	request.RecordId = common.Uint64Ptr(recordID)

	_, err = c.client.DeleteRecord(request)
	if err != nil {
		return fmt.Errorf(i18n.T("error.dns_delete"), err)
	}

	return nil
}

// getRecordID 获取记录 ID
func (c *DNSClient) getRecordID(domain, rr, recordType string) (uint64, error) {
	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = common.StringPtr(domain)
	request.Subdomain = common.StringPtr(rr)
	request.RecordType = common.StringPtr(recordType)

	response, err := c.client.DescribeRecordList(request)
	if err != nil {
		return 0, fmt.Errorf(i18n.T("error.dns_query"), err)
	}

	if response.Response.RecordCountInfo != nil && *response.Response.RecordCountInfo.TotalCount > 0 {
		if len(response.Response.RecordList) > 0 {
			return *response.Response.RecordList[0].RecordId, nil
		}
	}

	return 0, nil
}
