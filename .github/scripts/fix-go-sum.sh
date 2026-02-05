#!/bin/bash
# 修复go.sum文件

echo "开始修复go.sum文件..."

# 删除现有的go.sum文件（如果存在）
if [ -f go.sum ]; then
    rm go.sum
    echo "已删除现有的go.sum文件"
fi

# 清理Go模块缓存
go clean -modcache
echo "已清理Go模块缓存"

# 重新下载依赖并生成go.sum
go mod download
echo "已下载依赖"

# 验证模块
go mod verify
echo "模块验证完成"

# 生成go.sum文件
go mod tidy
echo "go.sum文件已修复"

# 验证生成的go.sum
if [ -f go.sum ]; then
    echo "go.sum文件生成成功"
    echo "Go模块路径: $(go list -m)"
else
    echo "错误: go.sum文件生成失败"
    exit 1
fi
