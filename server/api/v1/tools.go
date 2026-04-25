/*
项目名称：JeriBlog
文件名称：tools.go
创建时间：2026-04-16 15:02:06

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TODO
*/

package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"jeri_blog/pkg/linkparser"
	"jeri_blog/pkg/response"
	"jeri_blog/pkg/videoparser"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ToolsController 工具控制器
type ToolsController struct{}

// NewToolsController 创建工具控制器
func NewToolsController() *ToolsController {
	return &ToolsController{}
}

// ParseVideo 解析视频URL
// @Summary 解析视频URL
// @Description 解析视频URL，支持多个视频平台（抖音、快手、B站等）
// @Tags 管理工具
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{url=string} true "视频URL"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /admin/tools/parse-video [post]
func (c *ToolsController) ParseVideo(ctx *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	videoInfo := videoparser.ParseVideoURL(req.URL)
	if videoInfo == nil {
		response.Failed(ctx, "无法解析视频URL")
		return
	}

	if videoInfo.Platform == "" && videoInfo.VideoID == "" {
		response.Failed(ctx, "不支持的视频平台")
		return
	}

	response.Success(ctx, videoInfo)
}

// FetchLinkMetadata 获取链接元数据
// @Summary 获取链接元数据
// @Description 获取链接的标题、描述、图标等元数据信息
// @Tags 管理工具
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{url=string} true "链接URL"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /admin/tools/fetch-link-meta [post]
func (c *ToolsController) FetchLinkMetadata(ctx *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	metadata, err := linkparser.Parse(req.URL)
	if err != nil {
		response.Failed(ctx, err.Error())
		return
	}

	response.Success(ctx, metadata)
}

// DownloadImage 下载图片
// @Summary 下载图片到服务器
// @Description 从指定URL下载图片并保存到服务器
// @Tags 管理工具
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body object{url=string} true "图片URL"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /admin/tools/download-image [post]
func (c *ToolsController) DownloadImage(ctx *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidateFailed(ctx, err.Error())
		return
	}

	// 验证URL
	if !strings.HasPrefix(req.URL, "http") {
		response.Failed(ctx, "无效的URL")
		return
	}

	// 下载图片
	resp, err := http.Get(req.URL)
	if err != nil {
		response.Failed(ctx, fmt.Sprintf("下载图片失败: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		response.Failed(ctx, fmt.Sprintf("HTTP状态码: %d", resp.StatusCode))
		return
	}

	// 获取Content-Type，但不强制要求是图片类型
	contentType := resp.Header.Get("Content-Type")

	// 读取图片数据
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Failed(ctx, fmt.Sprintf("读取图片数据失败: %v", err))
		return
	}

	// 返回图片信息
	response.Success(ctx, gin.H{
		"content_type":   contentType,
		"content_length": len(data),
		"data":           data,
	})
}

// ParseMusic 解析音乐信息（代理Meting API）
// @Summary 解析音乐信息
// @Description 通过 Meting API 解析音乐平台的歌曲信息（网易云、QQ音乐等）
// @Tags 管理工具
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param server query string true "音乐平台（netease/tencent/kugou等）"
// @Param type query string true "类型（song/playlist/album等）"
// @Param id query string true "音乐ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /admin/tools/parse-music [get]
func (c *ToolsController) ParseMusic(ctx *gin.Context) {
	server := ctx.Query("server")
	musicType := ctx.Query("type")
	id := ctx.Query("id")

	// 参数验证
	if server == "" || musicType == "" || id == "" {
		response.ValidateFailed(ctx, "缺少必需参数: server, type, id")
		return
	}

	// 构建请求URL
	apiURL := fmt.Sprintf("https://api.injahow.cn/meting/?server=%s&type=%s&id=%s", server, musicType, id)

	// 发起请求
	resp, err := http.Get(apiURL)
	if err != nil {
		response.Failed(ctx, fmt.Sprintf("请求音乐API失败: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		response.Failed(ctx, fmt.Sprintf("音乐API返回错误状态码: %d", resp.StatusCode))
		return
	}

	// 读取响应数据
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Failed(ctx, fmt.Sprintf("读取响应数据失败: %v", err))
		return
	}

	// 解析JSON数据
	var musicData []interface{}
	if err := json.Unmarshal(data, &musicData); err != nil {
		response.Failed(ctx, fmt.Sprintf("解析音乐数据失败: %v", err))
		return
	}

	// 使用标准响应格式返回
	response.Success(ctx, musicData)
}
