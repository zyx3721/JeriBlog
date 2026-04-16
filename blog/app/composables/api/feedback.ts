/*
项目名称：JeriBlog
文件名称：feedback.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import type { Feedback, SubmitFeedbackParams } from '@@/types/feedback'
import { createApi } from './createApi'

const feedbackApi = createApi<Feedback>('')

/** 提交反馈 */
export const submitFeedback = async (data: SubmitFeedbackParams) => {
  return feedbackApi.post('/feedback', data)
}

/** 查询反馈状态 */
export const getFeedbackByTicketNo = async (ticketNo: string) => {
  return feedbackApi.get(`/feedback/ticket/${ticketNo}`)
}
