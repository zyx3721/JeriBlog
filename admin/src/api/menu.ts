/*
项目名称：JeriBlog
文件名称：menu.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - menu
*/

import request from '@/utils/request';
import type {
  MenuTreeNode,
  CreateMenuRequest,
  UpdateMenuRequest,
  MenuResponse,
} from '@/types/menu';

/**
 * 获取菜单树
 * @param type 菜单类型
 * @returns Promise<MenuTreeNode[]>
 */
export function getMenuTree(type?: string): Promise<MenuTreeNode[]> {
  return request.get('/admin/menus', {
    params: { type },
  });
}

/**
 * 创建菜单
 * @param data 菜单数据
 * @returns Promise<MenuResponse>
 */
export function createMenu(data: CreateMenuRequest): Promise<MenuResponse> {
  return request.post('/admin/menus', data);
}

/**
 * 获取菜单详情
 * @param id 菜单ID
 * @returns Promise<MenuResponse>
 */
export function getMenuDetail(id: number): Promise<MenuResponse> {
  return request.get(`/admin/menus/${id}`);
}

/**
 * 更新菜单
 * @param id 菜单ID
 * @param data 更新数据
 * @returns Promise<MenuResponse>
 */
export function updateMenu(id: number, data: UpdateMenuRequest): Promise<MenuResponse> {
  return request.put(`/admin/menus/${id}`, data);
}

/**
 * 删除菜单
 * @param id 菜单ID
 * @param data 子菜单处理方式
 * @returns Promise<void>
 */
export function deleteMenu(
  id: number,
  data?: { children_action?: 'delete' | 'upgrade' }
): Promise<void> {
  return request.delete(`/admin/menus/${id}`, { data });
}