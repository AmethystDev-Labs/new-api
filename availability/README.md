# Availability 模块 - 模型可用性检测

## 目录结构

```
availability/
├── test_models.py    # 测试脚本
├── run_test.sh       # 定时任务运行脚本
├── index.html        # 可视化报告页面
├── results.json      # 测试结果（运行后生成）
└── api_key.txt       # API密钥（需要创建）
```

## 使用方法

### 1. 配置 API 密钥

```bash
echo "your-api-key" > availability/api_key.txt
chmod 600 availability/api_key.txt
```

### 2. 手动运行测试

```bash
cd availability
python3 test_models.py -k YOUR_API_KEY
```

或使用便捷脚本：

```bash
cd availability
./run_test.sh
```

### 3. 配置定时任务

编辑 crontab：

```bash
crontab -e
```

添加以下内容（每天凌晨3点运行）：

```
0 3 * * * cd /path/to/new-api/availability && ./run_test.sh >> /var/log/availability.log 2>&1
```

### 4. 访问报告

启动 new-api 服务后，访问：

```
http://your-domain/availability/
```

## 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `AVAILABILITY_DIR` | 报告文件目录 | `availability` |
| `API_KEY` | API 密钥 | 从 api_key.txt 读取 |

## API 端点

| 端点 | 说明 |
|------|------|
| `GET /availability/` | 可视化报告首页 |
| `GET /availability/results.json` | 原始测试结果 |