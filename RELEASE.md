# 版本发布流程

## 版本管理

版本号统一在 `VERSION` 文件中管理，会自动同步到：
- 编译的二进制文件（通过 -ldflags 注入）
- 交互式菜单显示
- npm package.json

## 发布新版本

### 1. 更新版本号

编辑 `VERSION` 文件：
```bash
echo "1.0.6" > VERSION
```

### 2. 同步到 package.json

```bash
./sync-version.sh
```

### 3. 编译多平台二进制

```bash
./build.sh
```

这会自动：
- 读取 VERSION 文件
- 通过 `-ldflags` 将版本号注入到代码中
- 编译所有平台的二进制文件

### 4. 发布到 npm

```bash
cd npm-package
npm publish
```

### 5. 提交到 Git

```bash
git add VERSION npm-package/package.json npm-package/bin/
git commit -m "chore: bump version to x.x.x"
git push
git tag vx.x.x
git push --tags
```

## 自动化脚本（可选）

创建 `release.sh` 一键发布：

```bash
#!/bin/bash
set -e

VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: ./release.sh <version>"
    exit 1
fi

echo "Releasing version $VERSION..."

# 更新版本
echo "$VERSION" > VERSION

# 同步和编译
./sync-version.sh
./build.sh

# 发布
cd npm-package
npm publish
cd ..

# Git
git add VERSION npm-package/
git commit -m "chore: bump version to $VERSION"
git push
git tag "v$VERSION"
git push --tags

echo "✓ Released $VERSION successfully!"
```

使用：
```bash
chmod +x release.sh
./release.sh 1.0.6
```

## 版本号说明

遵循语义化版本 (Semantic Versioning)：

- **主版本号** (1.x.x): 不兼容的 API 变更
- **次版本号** (x.1.x): 向下兼容的功能性新增
- **修订号** (x.x.1): 向下兼容的问题修复

## 检查当前版本

```bash
cat VERSION
./certctl --version  # (如果实现了 version 命令)
```
