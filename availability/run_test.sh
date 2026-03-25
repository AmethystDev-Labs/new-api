#!/bin/bash
# 模型可用性检测定时任务脚本

cd "$(dirname "$0")" # 切换到脚本目录

# 从环境变量或配置文件读取 API 密钥
if [ -z "$API_KEY" ]; then
    if [ -f "api_key.txt" ]; then
        API_KEY=$(cat api_key.txt) # 从文件读取密钥
    else
        echo "错误: 请设置 API_KEY 环境变量或在 api_key.txt 中保存密钥"
        exit 1
    fi
fi

# 执行测试
python3 test_models.py -k "$API_KEY" -o results.json # 运行测试并保存结果

echo "检测完成: $(date)"