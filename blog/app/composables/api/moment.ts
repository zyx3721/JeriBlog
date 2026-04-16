/*
项目名称：JeriBlog
文件名称：moment.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import type { Moment } from '@@/types/moment'
import type { PaginationData, PaginationQuery } from '@@/types/request'
import { createApi } from './createApi'

const momentApi = createApi<Moment>('/moments')

/** 获取动态列表 */
export const getMoments = async (params: PaginationQuery = {}) => {
  return momentApi.getList(params)
}
