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

### Windows 用户特别说明

#### 安装方式

**推荐：使用 PowerShell**

```powershell
# 1. 安装 Node.js (如未安装)
# 下载地址: https://nodejs.org/

# 2. 使用官方 npm 源安装
npm install -g certctl-cli --registry https://registry.npmjs.org

# 3. 验证安装
certctl --help
```

#### 常见问题

**Q: 提示"无法将 certctl 项识别为 cmdlet"**

A: 这是 PATH 环境变量问题，解决方法：

1. **方法一：重启 PowerShell**（推荐）
   ```powershell
   # 关闭当前 PowerShell 窗口，重新打开
   certctl
   ```

2. **方法二：使用 npx 运行**（不需要安装）
   ```powershell
   npx certctl-cli
   ```

3. **方法三：手动添加 PATH**
   ```powershell
   # 查看 npm 全局路径
   npm bin -g
   
   # 将输出的路径添加到系统环境变量 Path
   # 控制面板 → 系统 → 高级系统设置 → 环境变量
   ```

**Q: 使用淘宝镜像安装失败**

A: 淘宝镜像同步需要时间，使用官方源：
```powershell
npm config set registry https://registry.npmjs.org
npm install -g certctl-cli
```

### Linux/macOS 用户

#### 安装

```bash
npm install -g certctl-cli
```

### 完整命令行参数

#### `certctl apply` - 申请证书

```
Usage:
  certctl apply [flags]

参数说明:
  -d, --domain string       要申请证书的域名（必填）
  -e, --email string        Let's Encrypt 账户邮箱（必填）
  -o, --output string       证书输出目录（默认: ~/.certctl/certs）
  
  DNS 自动验证:
      --dns string          DNS 提供商 (目前支持: aliyun)
      --ali-key string      阿里云 AccessKey ID
      --ali-secret string   阿里云 AccessKey Secret
  
  其他选项:
      --staging             使用测试环境（不计入速率限制）
      --dry-run             干跑模式（模拟流程，不实际申请）
      --lang string         界面语言 (zh/en)
  -h, --help                显示帮助信息

示例:
  # 手动 DNS 验证
  certctl apply -d example.com -e admin@example.com
  
  # 阿里云 DNS 自动验证
  certctl apply -d example.com -e admin@example.com \
    --dns aliyun --ali-key YOUR_KEY --ali-secret YOUR_SECRET
  
  # 使用测试环境
  certctl apply -d example.com -e admin@example.com --staging
  
  # 指定输出目录
  certctl apply -d example.com -e admin@example.com -o /path/to/certs
```

#### `certctl renew` - 续期证书

```
Usage:
  certctl renew [flags]

参数说明:
  -d, --domain string       要续期的域名
  -e, --email string        Let's Encrypt 账户邮箱（可选，使用已保存账户）
  -o, --output string       证书输出目录（默认: ~/.certctl/certs）
      --staging             使用测试环境
  -h, --help                显示帮助信息

示例:
  # 交互式选择续期
  certctl renew
  
  # 指定域名续期
  certctl renew -d example.com
  
  # 指定证书目录
  certctl renew -d example.com -o /path/to/certs
```

#### `certctl list` - 查看证书

```
Usage:
  certctl list [flags]

参数说明:
  -o, --output string       证书目录（默认: ~/.certctl/certs）
  -h, --help                显示帮助信息

示例:
  # 查看所有证书
  certctl list
  
  # 指定证书目录
  certctl list -o /path/to/certs
```

### 环境变量

支持通过环境变量配置阿里云 AccessKey：

```bash
# Linux/macOS
export ALICLOUD_ACCESS_KEY=YOUR_KEY
export ALICLOUD_SECRET_KEY=YOUR_SECRET
certctl apply -d example.com -e admin@example.com --dns aliyun

# Windows PowerShell
$env:ALICLOUD_ACCESS_KEY="YOUR_KEY"
$env:ALICLOUD_SECRET_KEY="YOUR_SECRET"
certctl apply -d example.com -e admin@example.com --dns aliyun
```

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
