package cmd

import (
	"fmt"
	"time"

	"certctl/internal/cert"
	"certctl/internal/config"
	"certctl/internal/ui"

	"github.com/spf13/cobra"
)

var listOutputDir string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "查看已申请的证书",
	Long:  "列出已申请的所有 SSL 证书及其有效期",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listOutputDir, "output", "o", "", "证书目录")
}

func runList(cmd *cobra.Command, args []string) error {
	fmt.Println()

	// 如果没有指定目录，使用配置中的证书目录
	if listOutputDir == "" {
		listOutputDir = config.Get().CertsDir
	}

	// 使用共享的证书扫描函数
	certs, err := cert.ListCertificates(listOutputDir)
	if err != nil {
		ui.Error(fmt.Sprintf("读取证书失败: %v", err))
		return nil
	}

	if len(certs) == 0 {
		ui.Info("暂无已申请的证书")
		return nil
	}

	// 显示证书列表
	ui.Title("已申请的证书:")
	fmt.Println()

	for _, c := range certs {
		status := "✔"
		daysLeft := int(time.Until(c.NotAfter).Hours() / 24)

		if daysLeft < 0 {
			status = "✖ 已过期"
		} else if daysLeft < 30 {
			status = fmt.Sprintf("⚠ %d 天后过期", daysLeft)
		} else {
			status = fmt.Sprintf("✔ %d 天后过期", daysLeft)
		}

		fmt.Printf("  %s\n", c.Domain)
		fmt.Printf("    状态: %s\n", status)
		fmt.Printf("    有效期至: %s\n", c.NotAfter.Format("2006-01-02"))
		fmt.Printf("    证书: %s\n", c.CertPath)
		fmt.Printf("    私钥: %s\n", c.KeyPath)
		fmt.Println()
	}

	return nil
}

