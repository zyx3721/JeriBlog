/*
项目名称：JeriBlog
文件名称：category.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import type { Category } from '@@/types/category'
import { createApi } from './createApi'

const categoryApi = createApi<Category>('/categories')

/** 获取分类列表 */
export const getCategories = async () => {
  return categoryApi.getList()
}

/** 获取分类详情（By ID） */
export const getCategoryById = async (id: number) => {
  return categoryApi.getOne(id)
}

/** 获取分类详情（By Slug） */
export const getCategoryBySlug = async (slug: string) => {
  return categoryApi.getOne(slug)
}
