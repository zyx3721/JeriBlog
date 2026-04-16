/*
项目名称：JeriBlog
文件名称：category.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：类型定义 - category类型
*/

// 分类实体
export interface Category {
  id: number
  name: string
  slug: string
  description: string
  count: number
  sort: number
}

// 分页数据
export interface CategoryListData {
  list: Category[]
  total: number
  page: number
  page_size: number
}