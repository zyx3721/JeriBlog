/*
项目名称：JeriBlog
文件名称：stats.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - stats
*/

import request from "@/utils/request";
import type { DashboardStats, TrendDataItem, TrendQuery, CategoryStats, TagStats, ArticleContribution, VisitListData, ContributionQuery } from "@/types/stats";
import type { PaginationQuery } from "@/types/request";

/**
 * 获取仪表板统计数据
 * @returns Promise<DashboardStats>
 */
export function getDashboardStats(): Promise<DashboardStats> {
  return request.get("/admin/stats/dashboard");
}

/**
 * 获取趋势数据
 * @param params 查询参数
 * @returns Promise<TrendDataItem[]>
 */
export function getTrendData(params: TrendQuery): Promise<TrendDataItem[]> {
  return request.get("/admin/stats/trend", { params });
}

/**
 * 获取分类统计数据
 * @returns Promise<CategoryStats[]>
 */
export function getCategoryStats(): Promise<CategoryStats[]> {
  return request.get("/admin/stats/category");
}

/**
 * 获取标签统计数据
 * @returns Promise<TagStats[]>
 */
export function getTagStats(): Promise<TagStats[]> {
  return request.get("/admin/stats/tag");
}

/**
 * 获取文章贡献数据
 * @param params 查询参数
 * @returns Promise<ArticleContribution[]>
 */
export function getArticleContribution(params?: ContributionQuery): Promise<ArticleContribution[]> {
  return request.get("/admin/stats/contribution", { params });
}

/**
 * 获取访问日志列表
 * @param params 查询参数
 * @returns Promise<VisitListData>
 */
export function getVisits(params: PaginationQuery): Promise<VisitListData> {
  return request.get("/admin/stats/visits", { params });
}
