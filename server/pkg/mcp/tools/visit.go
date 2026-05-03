/*
项目名称：JeriBlog
文件名称：visit.go
创建时间：2026-05-03 20:47:47

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：访问日志管理 MCP 工具
*/

package tools

import (
	"context"
	"fmt"

	"jeri_blog/internal/dto"
	"jeri_blog/internal/service"

	"github.com/google/jsonschema-go/jsonschema"
	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	visitActionList        = "list"
	visitActionDelete      = "delete"
	visitActionBatchDelete = "batch_delete"
)

// ============ MCP 类型定义============

// VisitLogItem 访问日志项
type VisitLogItem struct {
	ID        uint   `json:"id"`         // 记录ID
	VisitorID string `json:"visitor_id"` // 访客唯一标识
	IP        string `json:"ip"`         // 访客IP
	PageURL   string `json:"page_url"`   // 访问页面URL
	UserAgent string `json:"user_agent"` // 浏览器UA
	Location  string `json:"location"`   // 地理位置
	Browser   string `json:"browser"`    // 浏览器
	OS        string `json:"os"`         // 操作系统
	Referer   string `json:"referer"`    // 来源页面
	CreatedAt string `json:"created_at"` // 创建时间
}

// ============ 聚合 Tool 输入/输出类型============

// VisitManageInput visit_manage 聚合 tool 输入
type VisitManageInput struct {
	Action  string              `json:"action"` // list|delete|batch_delete
	Payload VisitManagePayload  `json:"payload"`
}

// VisitManagePayload visit_manage 载荷
type VisitManagePayload struct {
	// 用于 list
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Keyword    string `json:"keyword"`     // 搜索关键词（页面URL）
	VisitorID  string `json:"visitor_id"`  // 访客ID
	IP         string `json:"ip"`          // IP地址
	ExcludeIPs string `json:"exclude_ips"` // 排除的IP地址，多个用逗号分隔
	Location   string `json:"location"`    // 地理位置
	Browser    string `json:"browser"`     // 浏览器
	OS         string `json:"os"`          // 操作系统
	StartTime  string `json:"start_time"`  // 开始时间（格式：2006-01-02）
	EndTime    string `json:"end_time"`    // 结束时间（格式：2006-01-02）

	// 用于 delete
	ID uint `json:"id"`

	// 用于 batch_delete
	IDs []uint `json:"ids"`
}

// VisitManageOutput visit_manage 聚合 tool 输出
type VisitManageOutput struct {
	// list 结果
	List     []VisitLogItem `json:"list,omitempty"`
	Total    int64          `json:"total,omitempty"`
	Page     int            `json:"page,omitempty"`
	PageSize int            `json:"page_size,omitempty"`

	// delete 结果
	DeleteSuccess *bool `json:"delete_success,omitempty"`
	ID            *uint `json:"id,omitempty"`

	// batch_delete 结果
	BatchDeleteSuccess *bool `json:"batch_delete_success,omitempty"`
	DeletedCount       *int  `json:"deleted_count,omitempty"`

	// 错误信息（如果有）
	Error string `json:"error,omitempty"`
}

// ============ 服务包装器============

// VisitWrapper 访问日志服务包装器
type VisitWrapper struct {
	statsService *service.StatsService
}

// NewVisitWrapper 创建访问日志服务包装器
func NewVisitWrapper(statsService *service.StatsService) *VisitWrapper {
	return &VisitWrapper{statsService: statsService}
}

// ============ 聚合 Tool Handler============

// ManageVisit 访问日志管理聚合入口
func (w *VisitWrapper) ManageVisit(
	_ context.Context,
	_ *sdkmcp.CallToolRequest,
	input VisitManageInput,
) (*sdkmcp.CallToolResult, VisitManageOutput, error) {
	switch input.Action {
	case visitActionList:
		return w.listVisits(input.Payload)
	case visitActionDelete:
		return w.deleteVisit(input.Payload)
	case visitActionBatchDelete:
		return w.batchDeleteVisits(input.Payload)
	default:
		return nil, VisitManageOutput{}, fmt.Errorf("不支持的操作: %s", input.Action)
	}
}

// listVisits 获取访问日志列表
func (w *VisitWrapper) listVisits(payload VisitManagePayload) (*sdkmcp.CallToolResult, VisitManageOutput, error) {
	page, pageSize := NormalizePage(payload.Page, payload.PageSize)

	req := &dto.GetVisitLogsRequest{
		Page:       page,
		PageSize:   pageSize,
		Keyword:    payload.Keyword,
		VisitorID:  payload.VisitorID,
		IP:         payload.IP,
		ExcludeIPs: payload.ExcludeIPs,
		Location:   payload.Location,
		Browser:    payload.Browser,
		OS:         payload.OS,
		StartTime:  payload.StartTime,
		EndTime:    payload.EndTime,
	}

	visits, total, returnedPage, returnedPageSize, err := w.statsService.GetVisitLogs(req)
	if err != nil {
		return nil, VisitManageOutput{Error: fmt.Sprintf("获取访问日志列表失败: %v", err)}, nil
	}

	list := make([]VisitLogItem, len(visits))
	for i, visit := range visits {
		list[i] = VisitLogItem{
			ID:        visit.ID,
			VisitorID: visit.VisitorID,
			IP:        visit.IP,
			PageURL:   visit.PageURL,
			UserAgent: visit.UserAgent,
			Location:  visit.Location,
			Browser:   visit.Browser,
			OS:        visit.OS,
			Referer:   visit.Referer,
			CreatedAt: visit.CreatedAt,
		}
	}

	return nil, VisitManageOutput{
		List:     list,
		Total:    total,
		Page:     returnedPage,
		PageSize: returnedPageSize,
	}, nil
}

// deleteVisit 删除访问日志
func (w *VisitWrapper) deleteVisit(payload VisitManagePayload) (*sdkmcp.CallToolResult, VisitManageOutput, error) {
	if payload.ID == 0 {
		return nil, VisitManageOutput{Error: "访问日志 ID 不能为空"}, nil
	}

	err := w.statsService.DeleteVisitLog(payload.ID)
	if err != nil {
		return nil, VisitManageOutput{Error: fmt.Sprintf("删除访问日志失败: %v", err)}, nil
	}

	success := true
	return nil, VisitManageOutput{DeleteSuccess: &success, ID: &payload.ID}, nil
}

// batchDeleteVisits 批量删除访问日志
func (w *VisitWrapper) batchDeleteVisits(payload VisitManagePayload) (*sdkmcp.CallToolResult, VisitManageOutput, error) {
	if len(payload.IDs) == 0 {
		return nil, VisitManageOutput{Error: "访问日志 ID 列表不能为空"}, nil
	}

	err := w.statsService.BatchDeleteVisitLogs(payload.IDs)
	if err != nil {
		return nil, VisitManageOutput{Error: fmt.Sprintf("批量删除访问日志失败: %v", err)}, nil
	}

	success := true
	count := len(payload.IDs)
	return nil, VisitManageOutput{BatchDeleteSuccess: &success, DeletedCount: &count}, nil
}

// VisitManageInputSchema 返回 visit_manage 的自定义输入 schema
func VisitManageInputSchema() *jsonschema.Schema {
	listPayload := BuildPayloadSchema(map[string]*jsonschema.Schema{
		"page":        {Type: "integer"},
		"page_size":   PageSizeSchema(),
		"keyword":     {Type: "string", Description: "搜索关键词（页面URL）"},
		"visitor_id":  {Type: "string", Description: "访客ID"},
		"ip":          {Type: "string", Description: "IP地址"},
		"exclude_ips": {Type: "string", Description: "排除的IP地址，多个用逗号分隔"},
		"location":    {Type: "string", Description: "地理位置"},
		"browser":     {Type: "string", Description: "浏览器"},
		"os":          {Type: "string", Description: "操作系统"},
		"start_time":  {Type: "string", Description: "开始时间（格式：2006-01-02）"},
		"end_time":    {Type: "string", Description: "结束时间（格式：2006-01-02）"},
	})

	deletePayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"id": {Type: "integer"},
		},
		"id",
	)

	batchDeletePayload := BuildPayloadSchema(
		map[string]*jsonschema.Schema{
			"ids": {
				Type:        "array",
				Items:       &jsonschema.Schema{Type: "integer"},
				Description: "要删除的访问日志ID列表",
			},
		},
		"ids",
	)

	return &jsonschema.Schema{
		Type: "object",
		Properties: map[string]*jsonschema.Schema{
			"action": {
				Type: "string",
				Enum: []any{
					visitActionList,
					visitActionDelete,
					visitActionBatchDelete,
				},
			},
			"payload": {Type: "object"},
		},
		Required: []string{"action", "payload"},
		OneOf: []*jsonschema.Schema{
			BuildActionSchema(visitActionList, "获取访问日志列表", listPayload),
			BuildActionSchema(visitActionDelete, "删除单条访问日志。风险操作，谨慎使用，不可恢复", deletePayload),
			BuildActionSchema(visitActionBatchDelete, "批量删除访问日志。风险操作，谨慎使用，不可恢复", batchDeletePayload),
		},
	}
}
