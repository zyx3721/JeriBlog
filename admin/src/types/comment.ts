/*
项目名称：JeriBlog
文件名称：comment.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：类型定义 - comment类型
*/

// 评论实体
export interface Comment {
    id: number
    content: string
    status: number  // 0: 隐藏, 1: 显示
    parent_id?: number  // 父评论ID，用于回复
    created_at: string
    deleted_at?: string
    target: {
        type: string  // 评论目标类型：article(文章)、page(页面)等
        key: string   // 目标标识符
        title: string // 目标标题
    }
    user: {
        id: number
        nickname: string
        email: string
        avatar: string
    }
}

// 分页数据
export interface CommentListData {
    list: Comment[]
    total: number
    page: number
    page_size: number
}

// 评论导入相关类型
export interface ImportCommentError {
    index: number
    content: string
    error: string
}

export interface ImportCommentsResult {
    total: number
    success: number
    failed: number
    user_created: number
    errors?: ImportCommentError[]
}

