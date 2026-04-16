/*
项目名称：JeriBlog
文件名称：tag.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import type { Tag } from '@@/types/tag'
import { createApi } from './createApi'

const tagApi = createApi<Tag>('/tags')

/** 获取标签列表 */
export const getTags = async () => {
  return tagApi.getList()
}

/** 根据ID获取标签详情 */
export const getTagById = async (id: number) => {
  return tagApi.getOne(id)
}

/** 根据Slug获取标签详情 */
export const getTagBySlug = async (slug: string) => {
  return tagApi.getOne(slug)
}
