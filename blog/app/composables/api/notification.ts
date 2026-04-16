/*
项目名称：JeriBlog
文件名称：notification.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import type { NotificationListResponse, GetNotificationsParams } from '@@/types/notification'
import { createApi } from './createApi'

const notificationApi = createApi<NotificationListResponse>('')

/** 获取通知列表 */
export const getNotifications = async (params: GetNotificationsParams) => {
  return notificationApi.get<NotificationListResponse>('/notifications', params)
}

/** 标记单条通知已读 */
export const markAsRead = async (id: number) => {
  return notificationApi.put<void>(`/notifications/${id}/read`)
}

/** 标记全部通知已读 */
export const markAllAsRead = async () => {
  return notificationApi.put<void>('/notifications/read-all')
}
