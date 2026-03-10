package router

import (
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// SetAvailabilityRouter 设置模型可用性检测报告的路由
func SetAvailabilityRouter(router *gin.Engine) {
	availabilityRouter := router.Group("/availability") // 可用性报告路径前缀
	availabilityRouter.Use(middleware.RouteTag("availability")) // 标记路由类型
	availabilityRouter.Use(gzip.Gzip(gzip.DefaultCompression)) // 启用gzip压缩
	{
		availabilityRouter.GET("/*path", controller.ServeAvailability) // 所有子路径由控制器处理
	}
}