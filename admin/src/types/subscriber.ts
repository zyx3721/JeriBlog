/*
项目名称：JeriBlog
文件名称：subscriber.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：类型定义 - subscriber类型
*/

import type { PaginationQuery } from './request'

// 订阅者实体
export interface Subscriber {
    id: number
    email: string
    active: boolean
    created_at: string
    updated_at: string
}

// 订阅者查询参数
export interface SubscriberQuery extends PaginationQuery {}
