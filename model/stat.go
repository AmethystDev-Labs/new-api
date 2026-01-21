package model

import (
	"time"

	"github.com/QuantumNous/new-api/common"
)

type ModelStat struct {
	ModelName    string  `json:"model_name"`
	SuccessCount int64   `json:"success_count"`
	ErrorCount   int64   `json:"error_count"`
	AvgLatency   float64 `json:"avg_latency"`
}

var modelStats []ModelStat
var lastUpdate time.Time

func UpdateModelStats() {
	for {
		common.SysLog("正在更新模型成功率统计数据...")
		newStats, err := calculateModelStats()
		if err != nil {
			common.SysError("更新模型成功率统计失败: " + err.Error())
		} else {
			modelStats = newStats
			lastUpdate = time.Now()
			common.SysLog("模型成功率统计数据更新成功")
		}
		// 每 10 分钟更新一次
		time.Sleep(10 * time.Minute)
	}
}

func calculateModelStats() ([]ModelStat, error) {
	var stats []ModelStat
	// 统计过去 24 小时的数据
	startTime := time.Now().Add(-24 * time.Hour).Unix()

	// 使用原生 SQL 进行聚合统计
	// type = 2 为 LogTypeConsume (成功), type = 5 为 LogTypeError (失败)
	// 新增逻辑：即使 type 为 2，如果 completion_tokens <= 0 且为非流式请求（假设 stream 请求 latency 较低或有其他标识），也视为失败
	// 但根据用户要求：completion_tokens 在等于小于 0 的情况下也算失败
	err := LOG_DB.Table("logs").
		Select("model_name, "+
			"SUM(CASE WHEN type = 2 AND completion_tokens > 0 THEN 1 ELSE 0 END) as success_count, "+
			"SUM(CASE WHEN type = 5 OR (type = 2 AND completion_tokens <= 0) THEN 1 ELSE 0 END) as error_count, "+
			"AVG(use_time) as avg_latency").
		Where("created_at > ? AND model_name != ''", startTime).
		Group("model_name").
		Scan(&stats).Error

	return stats, err
}

func GetModelStats() []ModelStat {
	return modelStats
}
