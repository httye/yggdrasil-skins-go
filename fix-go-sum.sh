#!/bin/bash
# 修复go.sum文件

# 删除现有的go.sum文件
rm -f go.sum

# 重新下载依赖并生成go.sum
go mod download

# 验证模块
go mod verify

# 生成go.sum文件
go mod tidy

echo "go.sum 文件已修复"
