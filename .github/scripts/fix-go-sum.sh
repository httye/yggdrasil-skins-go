#!/bin/bash
# 修复go.sum文件

# 删除现有的go.sum文件（如果存在）
if [ -f go.sum ]; then
    rm go.sum
    echo "已删除现有的go.sum文件"
fi

# 重新下载依赖并生成go.sum
go mod download
echo "已下载依赖"

# 验证模块
go mod verify
echo "模块验证完成"

# 生成go.sum文件
go mod tidy
echo "go.sum文件已修复"
