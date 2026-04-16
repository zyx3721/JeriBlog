/*
项目名称：JeriBlog
文件名称：feedback.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - feedback
*/

import request from "@/utils/request"
import type {
  Feedback,
  FeedbackListData,
  FeedbackUpdateRequest
} from "@/types/feedback"
import type { PaginationQuery } from "@/types/request"

/**
 * 获取反馈列表
 */
export function getFeedbackList(params: PaginationQuery): Promise<FeedbackListData> {
  return request.get("/admin/feedback", { params })
}

/**
 * 获取反馈详情
 */
export function getFeedbackDetail(id: number): Promise<Feedback> {
  return request.get(`/admin/feedback/${id}`)
}

/**
 * 更新反馈
 */
export function updateFeedback(id: number, data: FeedbackUpdateRequest): Promise<void> {
  return request.put(`/admin/feedback/${id}`, data)
}

/**
 * 删除反馈
 */
export function deleteFeedback(id: number): Promise<void> {
  return request.delete(`/admin/feedback/${id}`)
}