/*
项目名称：JeriBlog
文件名称：subscriber.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - subscriber
*/

import request from "@/utils/request";
import type { Subscriber, SubscriberQuery } from "@/types/subscriber";

/**
 * 获取订阅者列表
 * @param params 查询参数
 * @returns Promise<SubscriberListResponse>
 */
export function getSubscribers(params: SubscriberQuery): Promise<any> {
    return request.get("/admin/subscribers", { params });
}

/**
 * 删除订阅者
 * @param id 订阅者ID
 * @returns Promise<void>
 */
export function deleteSubscriber(id: number): Promise<void> {
    return request.delete(`/admin/subscribers/${id}`);
}
