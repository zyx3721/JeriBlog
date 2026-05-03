/*
项目名称：JeriBlog
文件名称：rssfeed.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - rssfeed
*/

import request from '@/utils/request';
import type { RssArticleQuery, RssArticleListData } from '@/types/rssfeed';

/**
 * 全部标记已读响应
 */
export interface MarkAllReadResponse {
  affected: number;
}

/**
 * 获取RSS文章列表
 * @param params 查询参数
 * @returns Promise<RssArticleListData>
 */
export const getRssArticles = async (params?: RssArticleQuery): Promise<RssArticleListData> => {
  return request.get('/admin/rssfeed', { params });
};

/**
 * 标记文章已读
 * @param id 文章ID
 * @returns Promise<void>
 */
export const markRssArticleRead = async (id: number): Promise<void> => {
  await request.put(`/admin/rssfeed/${id}/read`);
};

/**
 * 全部标记已读
 * @returns Promise<MarkAllReadResponse>
 */
export const markAllRssArticlesRead = async (): Promise<MarkAllReadResponse> => {
  return request.put('/admin/rssfeed/read-all');
};

/**
 * 立即刷新所有RSS订阅源
 */
export const refreshAllRssFeeds = async (): Promise<{ message: string }> => {
  return request.post('/admin/rssfeed/refresh');
};
