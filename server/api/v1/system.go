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
// @Summary 获取系统静态信息
// @Description 获取系统静态信息（操作系统、Go版本、架构等）
// @Tags 系统信息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/system/static [get]
func (h *SystemController) GetSystemStatic(c *gin.Context) {
	response.Success(c, h.systemService.GetStaticInfo())
}

// GetSystemDynamic 获取系统动态信息。
// @Summary 获取系统动态信息
// @Description 获取系统动态信息（CPU使用率、内存使用率、磁盘使用率等）
// @Tags 系统信息
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/system/dynamic [get]
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
