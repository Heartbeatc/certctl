# Certctl

轻量级 SSL 证书申请 CLI 工具，支持通过 Let's Encrypt 自动申请**通配符证书**。

![Certctl 交互式菜单](docs/screenshot.png)

## ✨ 特性

- 🔐 支持通配符证书（*.example.com）
- 🤖 阿里云 DNS 自动验证
- 🌐 中英文双语界面
- 📋 证书管理（申请、续期、列表）
- 🎨 美观的交互式菜单

## 📦 安装

### 方式一：NPM 安装（推荐）

```bash
npm install -g certctl-cli
```

### 方式二：从源码编译

```bash
git clone https://github.com/cuijianzhong/certctl.git
cd certctl
go build -o certctl
```

## 🚀 快速开始

### 交互式菜单

直接运行 `certctl`，进入交互式菜单：

```bash
certctl
```

### 命令行使用

#### 1. 申请证书

**使用阿里云 DNS 自动验证**（推荐）：

```bash
certctl apply -d example.com \
  -e admin@example.com \
  --dns aliyun \
  --ali-key YOUR_ACCESS_KEY \
  --ali-secret YOUR_ACCESS_SECRET
```

或通过环境变量：

```bash
export ALICLOUD_ACCESS_KEY=YOUR_KEY
export ALICLOUD_SECRET_KEY=YOUR_SECRET
certctl apply -d example.com -e admin@example.com --dns aliyun
```

**手动 DNS 验证**：

```bash
certctl apply -d example.com -e admin@example.com
# 按提示手动添加 DNS TXT 记录
```

#### 2. 查看证书

```bash
certctl list
```

#### 3. 续期证书

```bash
certctl renew
# 或指定域名
certctl renew -d example.com
```

## 📂 证书输出

证书以 Nginx 格式保存到 `~/.certctl/certs/` 目录：

```
~/.certctl/certs/
└── example.com/
    ├── example.com.pem  # 证书链（公钥）
    └── example.com.key  # 私钥
```

### Nginx 配置示例

```nginx
server {
    listen 443 ssl;
    server_name example.com;

    ssl_certificate     /root/.certctl/certs/example.com/example.com.pem;
    ssl_certificate_key /root/.certctl/certs/example.com/example.com.key;

    # 其他配置...
}
```

## ⚙️ 配置

### 语言设置

```bash
# 在交互式菜单中选择"设置" -> "语言"
certctl
```

配置文件位置：`~/.certctl/config.json`

## 🔑 获取阿里云 AccessKey

1. 访问 https://ram.console.aliyun.com/manage/ak
2. 创建 AccessKey
3. 赋予 DNS 管理权限

## 🌍 环境选择

测试环境（不计入速率限制）：

```bash
certctl apply -d example.com --staging
```

生产环境（默认）：

```bash
certctl apply -d example.com
```

## 📋 常见问题

### 1. 证书到期了怎么办？

使用 `certctl renew` 命令续期，或设置 cron 定时任务：

```bash
# 每月 1 号凌晨 3 点检查续期
0 3 1 * * certctl renew -d example.com
```

### 2. 支持哪些 DNS 提供商？

目前支持：
- 阿里云 DNS（自动验证）
- 手动验证（所有 DNS 提供商）

### 3. Windows 上安装后找不到命令？

请确保 npm 全局安装目录在系统 PATH 中：

```powershell
npm bin -g
# 将输出的路径添加到系统环境变量 Path
```

## 📝 License

MIT

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📧 联系

- GitHub: https://github.com/Heartbeatc/certctl
- NPM: https://www.npmjs.com/package/certctl-cli
