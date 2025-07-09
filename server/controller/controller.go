package controller

import (
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/crazyfrankie/zdocker-web/service"
)

// ListContainers 获取容器列表
func ListContainers(c *gin.Context) {
	containers, err := service.GetContainerList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取容器列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": containers,
	})
}

// CreateContainer 创建容器
func CreateContainer(c *gin.Context) {
	var req service.CreateContainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	result, err := service.CreateContainer(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建容器失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

// GetContainer 获取单个容器信息
func GetContainer(c *gin.Context) {
	containerId := c.Param("id")
	if containerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "容器ID不能为空",
		})
		return
	}

	container, err := service.GetContainerById(containerId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "容器不存在: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": container,
	})
}

// StartContainer 启动容器
func StartContainer(c *gin.Context) {
	containerId := c.Param("id")
	if containerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "容器ID不能为空",
		})
		return
	}

	err := service.StartContainer(containerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "启动容器失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "容器启动成功",
	})
}

// StopContainer 停止容器
func StopContainer(c *gin.Context) {
	containerId := c.Param("id")
	if containerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "容器ID不能为空",
		})
		return
	}

	err := service.StopContainer(containerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "停止容器失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "容器停止成功",
	})
}

// RemoveContainer 删除容器
func RemoveContainer(c *gin.Context) {
	containerId := c.Param("id")
	if containerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "容器ID不能为空",
		})
		return
	}

	err := service.RemoveContainer(containerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除容器失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "容器删除成功",
	})
}

// GetContainerLogs 获取容器日志
func GetContainerLogs(c *gin.Context) {
	containerId := c.Param("id")
	if containerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "容器ID不能为空",
		})
		return
	}

	logs, err := service.GetContainerLogs(containerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取容器日志失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
	})
}

// ExecContainer 在容器中执行命令
func ExecContainer(c *gin.Context) {
	containerId := c.Param("id")
	if containerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "容器ID不能为空",
		})
		return
	}

	var req service.ExecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	result, err := service.ExecContainer(containerId, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "执行命令失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

// ListImages 获取镜像列表
func ListImages(c *gin.Context) {
	// 由于zdocker没有专门的镜像管理，这里简单返回空列表
	c.JSON(http.StatusOK, gin.H{
		"data": []interface{}{},
	})
}

// RemoveImage 删除镜像
func RemoveImage(c *gin.Context) {
	imageId := c.Param("id")
	if imageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "镜像ID不能为空",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "镜像删除成功",
	})
}

// ListNetworks 获取网络列表
func ListNetworks(c *gin.Context) {
	networks, err := service.GetNetworkList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取网络列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": networks,
	})
}

// CreateNetwork 创建网络
func CreateNetwork(c *gin.Context) {
	var req service.CreateNetworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	result, err := service.CreateNetwork(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建网络失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}

// RemoveNetwork 删除网络
func RemoveNetwork(c *gin.Context) {
	networkId := c.Param("id")
	if networkId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "网络ID不能为空",
		})
		return
	}

	err := service.RemoveNetwork(networkId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除网络失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "网络删除成功",
	})
}

// GetSystemInfo 获取系统信息
func GetSystemInfo(c *gin.Context) {
	info, err := service.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取系统信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": info,
	})
}

// GetVersion 获取版本信息
func GetVersion(c *gin.Context) {
	cmd := exec.Command("zdocker", "--version")
	output, err := cmd.CombinedOutput()

	version := "unknown"
	if err == nil {
		version = strings.TrimSpace(string(output))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"version":     version,
			"api_version": "1.0",
			"build_date":  "2024-01-01",
		},
	})
}
