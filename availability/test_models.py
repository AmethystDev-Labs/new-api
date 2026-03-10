#!/usr/bin/env python3
"""
API 模型可用性检测脚本
测试 API 端点上所有模型并生成报告
"""

import json
import time
import requests
import argparse
import os
from datetime import datetime

scriptDir = os.path.dirname(os.path.abspath(__file__)) # 获取脚本所在目录
defaultBaseUrl = "https://api.amethyst.ltd/v1" # 默认 API 地址
sampleImageBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8DwHwAFBQIAX8jx0gAAAABJRU5ErkJggg==" # 测试用1x1像素图片
timeout = 120 # 请求超时时间
delayBetweenRequests = 2.0 # 请求间隔时间
retryOn429 = 3 # 429错误重试次数


def getModelsList(baseUrl: str, apiKey: str) -> list:
    headers = {"Authorization": f"Bearer {apiKey}"} # 设置认证头
    try:
        response = requests.get(f"{baseUrl}/models", headers=headers, timeout=30) # 请求模型列表
        response.raise_for_status() # 检查响应状态
        data = response.json() # 解析响应数据
        return [m.get("id") for m in data.get("data", []) if m.get("id")] # 提取模型ID列表
    except Exception as e:
        print(f"获取模型列表失败: {e}")
        return []


def testBasic(baseUrl: str, apiKey: str, model: str) -> dict:
    headers = {"Authorization": f"Bearer {apiKey}", "Content-Type": "application/json"} # 请求头
    payload = {
        "model": model,
        "messages": [{"role": "user", "content": "回复OK"}], # 简单测试消息
        "max_tokens": 30, # 限制输出长度
        "temperature": 0.1 # 降低随机性
    }
    
    result = {"available": False, "error": None, "responseTime": None, "isThinking": False} # 初始化结果
    
    # --- 重试循环 ---
    for attempt in range(retryOn429):
        try:
            startTime = time.time() # 记录开始时间
            response = requests.post(f"{baseUrl}/chat/completions", headers=headers, json=payload, timeout=timeout) # 发送测试请求
            result["responseTime"] = time.time() - startTime # 计算响应时间
            
            if response.status_code == 200: # 请求成功
                data = response.json()
                choices = data.get("choices", [])
                if choices:
                    msg = choices[0].get("message", {})
                    content = msg.get("content", "")
                    result["available"] = True # 标记模型可用
                    
                    # --- 检测思考模型 ---
                    reasoning = msg.get("reasoning_content") or msg.get("thinking")
                    if reasoning or (content and ("<tool_call>" in content or "thinking" in content.lower())):
                        result["isThinking"] = True # 标记为思考模型
                return result
            elif response.status_code == 429: # 触发限流
                waitTime = 5 * (attempt + 1) # 递增等待时间
                print(f"  429限流，等待{waitTime}秒...")
                time.sleep(waitTime)
                continue
            else:
                result["error"] = f"HTTP {response.status_code}" # 记录错误码
                return result
        except requests.exceptions.Timeout:
            result["error"] = "请求超时"
            return result
        except Exception as e:
            result["error"] = str(e)[:50] # 截断错误信息
            return result
    
    result["error"] = "多次429重试失败"
    return result


def testVision(baseUrl: str, apiKey: str, model: str) -> bool:
    headers = {"Authorization": f"Bearer {apiKey}", "Content-Type": "application/json"}
    payload = {
        "model": model,
        "messages": [{
            "role": "user",
            "content": [
                {"type": "text", "text": "什么颜色?"}, # 询问图片颜色
                {"type": "image_url", "image_url": {"url": f"data:image/png;base64,{sampleImageBase64}"}}
            ]
        }],
        "max_tokens": 20,
        "temperature": 0.1
    }
    
    try:
        response = requests.post(f"{baseUrl}/chat/completions", headers=headers, json=payload, timeout=timeout)
        if response.status_code == 200:
            data = response.json()
            return bool(data.get("choices", [{}])[0].get("message", {}).get("content")) # 有返回内容则支持视觉
    except:
        pass
    return False


def testAllModels(baseUrl: str, apiKey: str, delay: float):
    print(f"\nAPI地址: {baseUrl}")
    print(f"请求间隔: {delay}秒\n")
    
    # --- 获取模型列表 ---
    print("获取模型列表...")
    models = getModelsList(baseUrl, apiKey)
    if not models:
        print("未找到模型")
        return []
    print(f"找到 {len(models)} 个模型\n")
    
    results = [] # 所有模型结果
    thinkingModels = [] # 思考模型列表
    visionModels = [] # 多模态模型列表
    
    # --- 逐个测试模型 ---
    for i, model in enumerate(models, 1):
        print(f"[{i}/{len(models)}] {model[:45]}")
        
        basic = testBasic(baseUrl, apiKey, model) # 测试基本能力
        
        if basic["available"]:
            isThinking = " [思考模型]" if basic["isThinking"] else ""
            print(f"  可用{isThinking} ({basic['responseTime']:.1f}秒)")
            
            if basic["isThinking"]:
                thinkingModels.append(model)
            
            time.sleep(delay)
            if testVision(baseUrl, apiKey, model): # 测试多模态能力
                print(f"  [支持视觉]")
                visionModels.append(model)
        else:
            print(f"  失败: {basic['error']}")
        
        results.append({"model": model, **basic, "supportsVision": model in visionModels})
        
        if i < len(models):
            time.sleep(delay) # 请求间隔
    
    # --- 打印摘要 ---
    print(f"\n总计: {len(results)} | 可用: {len([r for r in results if r['available']])}")
    print(f"思考模型: {len(thinkingModels)}")
    print(f"多模态模型: {len(visionModels)}")
    
    return results


def saveResults(results: list, outputFile: str):
    available = [r for r in results if r.get("available")] # 筛选可用模型
    output = {
        "timestamp": datetime.now().isoformat(), # 时间戳
        "total": len(results), # 总数
        "available": len(available), # 可用数
        "thinkingModels": [r["model"] for r in available if r.get("isThinking")], # 思考模型
        "visionModels": [r["model"] for r in available if r.get("supportsVision")], # 多模态模型
        "models": results # 详细结果
    }
    with open(outputFile, "w", encoding="utf-8") as f:
        json.dump(output, f, ensure_ascii=False, indent=2) # 保存JSON文件
    print(f"\n结果已保存: {outputFile}")


def main():
    parser = argparse.ArgumentParser(description="测试API模型可用性")
    parser.add_argument("--key", "-k", required=True, help="API密钥")
    parser.add_argument("--url", default=defaultBaseUrl, help="API地址")
    parser.add_argument("--output", "-o", default=os.path.join(scriptDir, "results.json")) # 默认输出到同目录
    parser.add_argument("--delay", "-d", type=float, default=2.0, help="请求间隔秒数")
    args = parser.parse_args()
    
    results = testAllModels(args.url, args.key, args.delay)
    if results:
        saveResults(results, args.output)


if __name__ == "__main__":
    main()