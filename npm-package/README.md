# certctl-cli

轻量级 SSL 证书申请工具 / Lightweight SSL Certificate CLI Tool

支持通配符证书申请 & 阿里云 DNS 自动验证。

## 安装

```bash
npm install -g certctl-cli
```

## 使用

```bash
# 交互式菜单
certctl

# 手动 DNS 验证
certctl apply -d example.com -e admin@example.com

# 阿里云 DNS 自动验证
certctl apply -d example.com --dns aliyun --ali-key YOUR_KEY --ali-secret YOUR_SECRET

# 查看证书
certctl list

# 续期证书
certctl renew
```

## 文档

完整文档: https://github.com/Heartbeatc/certctl

## License

MIT
