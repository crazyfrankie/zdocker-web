package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/zdocker-web/controller"
	"github.com/crazyfrankie/zdocker-web/middleware"
)

func main() {
	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 获取端口号，默认8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 创建gin路由
	r := gin.Default()

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// 添加中间件
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// 注册路由
	setupRoutes(r)

	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}

func setupRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API路由组
	api := r.Group("/api/v1")
	{
		// 容器相关路由
		containers := api.Group("/containers")
		{
			containers.GET("", controller.ListContainers)
			containers.POST("", controller.CreateContainer)
			containers.GET("/logs/:name", controller.GetContainerLogs)
			containers.GET("/:id", controller.GetContainer)
			containers.POST("/:id/start", controller.StartContainer)
			containers.POST("/stop/:name", controller.StopContainer)
			containers.DELETE("/:name", controller.RemoveContainer)
			containers.POST("/:id/exec", controller.ExecContainer)
		}

		// 镜像相关路由
		images := api.Group("/images")
		{
			images.GET("", controller.ListImages)
			images.DELETE("/:id", controller.RemoveImage)
		}

		// 网络相关路由
		networks := api.Group("/networks")
		{
			networks.GET("", controller.ListNetworks)
			networks.POST("", controller.CreateNetwork)
			networks.DELETE("/:id", controller.RemoveNetwork)
		}

		// 系统信息
		api.GET("/system/info", controller.GetSystemInfo)
		api.GET("/system/version", controller.GetVersion)
	}
}
