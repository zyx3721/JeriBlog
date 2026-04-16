/*
项目名称：JeriBlog
文件名称：ai.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：类型定义 - ai类型
*/

// AI功能请求类型
export interface AISummaryRequest {
    content: string
}

export interface AIAISummaryRequest {
    content: string
}

export interface AITitleRequest {
    content: string
}

// AI功能响应类型
export interface AISummaryResponse {
    summary: string
}

export interface AIAISummaryResponse {
    summary: string
}

export interface AITitleResponse {
    title: string
}