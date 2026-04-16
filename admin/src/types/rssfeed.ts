/*
项目名称：JeriBlog
文件名称：rssfeed.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：类型定义 - rssfeed类型
*/

// RSS文章实体
export interface RssArticle {
  id: number
  friend_id: number
  friend_name: string
  friend_url: string
  title: string
  link: string
  description: string
  published_at?: string
  is_read: boolean
  created_at: string
}

// RSS文章列表查询参数
export interface RssArticleQuery {
  page?: number
  page_size?: number
}

// RSS文章列表响应数据
export interface RssArticleListData {
  list: RssArticle[]
  total: number
  page: number
  page_size: number
  unread_count: number
}
