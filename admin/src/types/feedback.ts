/*
项目名称：JeriBlog
文件名称：feedback.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：类型定义 - feedback类型
*/

// 举报类型
export type ReportType = 'copyright' | 'inappropriate' | 'summary' | 'suggestion'

// 反馈状态
export type FeedbackStatus = 'pending' | 'resolved' | 'closed'

// 反馈内容结构
export interface FeedbackContent {
  description: string
  reason?: string
  attachmentFiles?: string[]
}

// 反馈对象
export interface Feedback {
  id: number
  ticket_no: string
  report_url: string
  report_type: ReportType
  form_content: FeedbackContent
  email: string
  status: FeedbackStatus
  admin_reply: string
  reply_time?: string
  user_agent: string
  ip: string
  feedback_time: string
}

// 反馈列表数据
export interface FeedbackListData {
  list: Feedback[]
  total: number
  page: number
  page_size: number
}

// 更新反馈请求
export interface FeedbackUpdateRequest {
  status: FeedbackStatus
  admin_reply?: string
}