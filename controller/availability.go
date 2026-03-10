package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var availabilityDir = "availability" // 模型可用性检测结果的存放目录

func init() {
	if dir := os.Getenv("AVAILABILITY_DIR"); dir != "" {
		availabilityDir = dir // 允许通过环境变量自定义目录
	}
}

// ServeAvailability 提供模型可用性检测报告的静态文件服务
func ServeAvailability(c *gin.Context) {
	relativePath := c.Param("path") // 获取请求的相对路径

	// --- 处理默认路径 ---
	if relativePath == "" || relativePath == "/" {
		relativePath = "/index.html" // 默认返回首页
	}

	relativePath = strings.TrimPrefix(relativePath, "/") // 移除前导斜杠
	filePath := filepath.Join(availabilityDir, relativePath) // 拼接完整文件路径

	// --- 读取文件内容 ---
	content, err := os.ReadFile(filePath)
	if err != nil {
		c.Status(http.StatusNotFound) // 文件不存在返回404
		return
	}

	// --- 根据扩展名设置内容类型 ---
	contentType := "application/octet-stream" // 默认二进制类型
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".html":
		contentType = "text/html; charset=utf-8" // HTML文档
	case ".css":
		contentType = "text/css; charset=utf-8" // 样式表
	case ".js":
		contentType = "application/javascript; charset=utf-8" // 脚本
	case ".json":
		contentType = "application/json; charset=utf-8" // JSON数据
	case ".png":
		contentType = "image/png" // PNG图片
	case ".jpg", ".jpeg":
		contentType = "image/jpeg" // JPEG图片
	case ".svg":
		contentType = "image/svg+xml" // SVG矢量图
	case ".ico":
		contentType = "image/x-icon" // 网站图标
	}

	c.Header("Cache-Control", "no-cache") // 禁用缓存以获取最新结果
	c.Data(http.StatusOK, contentType, content) // 返回文件内容
}