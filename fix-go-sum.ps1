# 修复go.sum文件

# 删除现有的go.sum文件
if (Test-Path go.sum) {
    Remove-Item go.sum
    Write-Host "已删除现有的go.sum文件"
}

# 重新下载依赖并生成go.sum
go mod download
Write-Host "已下载依赖"

# 验证模块
go mod verify
Write-Host "模块验证完成"

# 生成go.sum文件
go mod tidy
Write-Host "go.sum文件已修复"
