/*
项目名称：JeriBlog
文件名称：stats.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import type { SiteStats, ArchiveStats } from '@@/types/stats'
import { createApi } from './createApi'

const statsApi = createApi<SiteStats>('')

/** 获取网站统计信息 */
export const getSiteStats = async () => {
  return statsApi.get<SiteStats>('/stats/site')
}

/** 获取归档统计信息 */
export const getArchiveStats = async () => {
  return statsApi.get<ArchiveStats>('/stats/archives')
}
