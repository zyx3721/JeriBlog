/*
项目名称：JeriBlog
文件名称：notification.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - notification
*/

import request from "@/utils/request";
import type {
  NotificationListData,
  NotificationQueryParams
} from "@/types/notification";

/**
 * 获取管理员通知列表
 * @param params 查询参数
 * @returns Promise<NotificationListData>
 */
export function getNotifications(params: NotificationQueryParams): Promise<NotificationListData> {
  return request.get("/admin/notifications", { params });
}

/**
 * 标记单条已读
 * @param id 通知ID
 * @returns Promise<void>
 */
export function markAsRead(id: number): Promise<void> {
  return request.put(`/admin/notifications/${id}/read`);
}

/**
 * 全部标记已读
 * @returns Promise<void>
 */
export function markAllAsRead(): Promise<void> {
  return request.put("/admin/notifications/read-all");
}


