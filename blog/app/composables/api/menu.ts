/*
项目名称：JeriBlog
文件名称：menu.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import type { Menu } from '@@/types/menu'
import { createApi } from './createApi'

const menuApi = createApi<Menu[]>('')

/** 获取菜单列表 */
export const getMenus = async () => {
  return menuApi.get<Menu[]>('/menus')
}
