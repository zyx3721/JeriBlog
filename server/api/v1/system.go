/*
项目名称：JeriBlog
文件名称：system.go
创建时间：2026-04-16 15:02:06

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TODO
*/

package v1

import (
	"net/http"

	"jeri_blog/internal/service"
	"jeri_blog/pkg/response"

	"github.com/gin-gonic/gin"
)

// SystemController 系统信息控制器。
type SystemController struct {
	systemService *service.SystemService
}

// NewSystemController 创建系统信息控制器。
func NewSystemController(systemService *service.SystemService) *SystemController {
	return &SystemController{systemService: systemService}
}

// GetSystemStatic 获取系统静态信息。
func (h *SystemController) GetSystemStatic(c *gin.Context) {
	response.Success(c, h.systemService.GetStaticInfo())
}

// GetSystemDynamic 获取系统动态信息。
func (h *SystemController) GetSystemDynamic(c *gin.Context) {
	response.Success(c, h.systemService.GetDynamicInfo())
}

// Health 健康检查接口，检查服务及数据库连接状态。
func (h *SystemController) Health(c *gin.Context) {
	dbStatus := h.systemService.CheckHealth()
	if dbStatus != "正常" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "db": dbStatus})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "db": dbStatus})
}
