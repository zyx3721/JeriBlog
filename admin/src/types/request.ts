/*
项目名称：JeriBlog
文件名称：request.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：类型定义 - request类型
*/

// API响应数据结构
export interface ApiResponse<T = any> {
    code: number
    message: string
    data: T
}

// 分页查询参数
export interface PaginationQuery {
    page: number
    page_size: number
}
