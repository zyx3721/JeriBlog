/*
项目名称：JeriBlog
文件名称：stats.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：类型定义 - stats类型
*/

// Dashboard 统计数据
export interface DashboardStats {
  total_articles: number
  total_friends: number
  total_moments: number
  total_views: number
  total_visitors: number
  total_comments: number
  total_users: number
  today_views: number
  today_visitors: number
  today_comments: number
  today_users: number
  views_growth: number
  visitors_growth: number
  comments_growth: number
  users_growth: number
}

// 趋势数据项
export interface TrendDataItem {
  date: string
  pv_count: number
  uv_count: number
}

// 趋势查询参数
export interface TrendQuery {
  start_date: string
  end_date: string
  type: 'daily' | 'monthly'
}

// 标签统计数据
export interface TagStats {
  name: string
  count: number
}

// 文章贡献数据
export interface ArticleContribution {
  date: string
  count: number
}

// 文章贡献查询参数
export interface ContributionQuery {
  year?: number  // 年份（可选）
  month?: number // 月份 1-12（可选）
}

// 分类统计数据
export interface CategoryStats {
  name: string
  count: number
}

// 访问日志实体
export interface Visit {
  id: number
  visitor_id: string
  ip: string
  page_url: string
  user_agent: string
  location: string
  browser: string
  os: string
  referer: string
  created_at: string
}

// 访问日志列表数据
export interface VisitListData {
  list: Visit[]
  total: number
  page: number
  page_size: number
}

// 访问日志查询参数
export interface VisitQuery {
  page: number
  page_size: number
  keyword?: string
  start_date?: string
  end_date?: string
}
