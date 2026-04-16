/*
项目名称：JeriBlog
文件名称：article.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import type { Article, ArticleQuery } from '@@/types/article'
import type { PaginationData } from '@@/types/request'
import { createApi } from './createApi'

const articleApi = createApi<Article>('/articles')

/** 获取文章列表 */
export const getArticlesForWeb = async (params: ArticleQuery = {}) => {
  return articleApi.getList(params)
}

/** 获取文章详情 */
export const getArticleBySlug = async (slug: string) => {
  return articleApi.getOne(slug)
}

/** 搜索文章 */
export const searchArticles = async (keyword: string, params: Partial<ArticleQuery> = {}) => {
  return articleApi.get<PaginationData<Article>>('/search', { keyword, ...params })
}
