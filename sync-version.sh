#!/bin/bash
# 同步版本号到 package.json

set -e

VERSION=$(cat VERSION)

# 使用 sed 更新 package.json 中的版本号
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/\"version\": \".*\"/\"version\": \"$VERSION\"/" npm-package/package.json
else
    # Linux
    sed -i "s/\"version\": \".*\"/\"version\": \"$VERSION\"/" npm-package/package.json
fi

echo "✓ Synced version $VERSION to package.json"
