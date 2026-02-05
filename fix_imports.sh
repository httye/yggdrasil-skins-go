#!/bin/bash
# 批量修复所有Go文件中的导入路径

echo "开始修复所有Go文件中的导入路径..."

# 查找所有包含旧导入路径的.go文件
find src -name "*.go" -type f -exec grep -l "yggdrasil-api-go/" {} \; | while read file; do
    echo "修复文件: $file"
    
    # 使用sed替换所有旧导入路径为新导入路径
    sed -i 's|yggdrasil-api-go/|github.com/httye/yggdrasil-skins-go/|g' "$file"
    
    echo "已修复: $file"
done

echo "所有文件导入路径修复完成！"

# 验证修复结果
echo "验证修复结果..."
remaining_files=$(find src -name "*.go" -type f -exec grep -l "yggdrasil-api-go/" {} \; | wc -l)
if [ "$remaining_files" -eq 0 ]; then
    echo "✅ 所有导入路径已成功修复！"
else
    echo "⚠️  还有 $remaining_files 个文件需要修复"
    find src -name "*.go" -type f -exec grep -l "yggdrasil-api-go/" {} \;
fi
