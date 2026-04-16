/*
项目名称：JeriBlog
文件名称：rssfeed.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - rssfeed
*/

import request from '@/utils/request'
import type { RssArticleQuery, RssArticleListData } from '@/types/rssfeed'

/**
 * 获取RSS文章列表
 */
export const getRssArticles = async (params?: RssArticleQuery): Promise<RssArticleListData> => {
  return request.get('/admin/rssfeed', { params })
}

/**
 * 标记文章已读
 */
export const markRssArticleRead = async (id: number): Promise<void> => {
  await request.put(`/admin/rssfeed/${id}/read`)
}

/**
 * 全部标记已读
 */
export const markAllRssArticlesRead = async (): Promise<{ affected: number }> => {
  return request.put('/admin/rssfeed/read-all')
}
