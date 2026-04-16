/*
项目名称：JeriBlog
文件名称：friend.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - friend
*/

import request from "@/utils/request";
import type {
  Friend,
  FriendListData,
  FriendQuery,
  CreateFriendRequest,
  UpdateFriendRequest,
  FriendType,
  FriendTypeListData,
  CreateFriendTypeRequest,
  UpdateFriendTypeRequest
} from "@/types/friend";

// ========== 友链类型相关API ==========

/**
 * 获取友链类型列表（管理端，包括隐藏的）
 * @returns Promise<FriendTypeListData>
 */
export function getFriendTypes(): Promise<FriendTypeListData> {
  return request.get("/admin/friends/types");
}

/**
 * 获取友链类型详情
 * @param id 类型ID
 * @returns Promise<FriendType>
 */
export function getFriendTypeDetail(id: number): Promise<FriendType> {
  return request.get(`/admin/friends/types/${id}`);
}

/**
 * 创建友链类型
 * @param data 类型数据
 * @returns Promise<FriendType>
 */
export function createFriendType(data: CreateFriendTypeRequest): Promise<FriendType> {
  return request.post("/admin/friends/types", data);
}

/**
 * 更新友链类型
 * @param id 类型ID
 * @param data 更新数据
 * @returns Promise<FriendType>
 */
export function updateFriendType(id: number, data: UpdateFriendTypeRequest): Promise<FriendType> {
  return request.put(`/admin/friends/types/${id}`, data);
}

/**
 * 删除友链类型
 * @param id 类型ID
 * @returns Promise<void>
 */
export function deleteFriendType(id: number): Promise<void> {
  return request.delete(`/admin/friends/types/${id}`);
}

// ========== 友链相关API ==========

/**
 * 获取友链列表
 * @param params 查询参数
 * @returns Promise<FriendListData>
 */
export function getFriends(params?: FriendQuery): Promise<FriendListData> {
  return request.get("/admin/friends", { params });
}

/**
 * 获取友链详情
 * @param id 友链ID
 * @returns Promise<Friend>
 */
export function getFriendDetail(id: number): Promise<Friend> {
  return request.get(`/admin/friends/${id}`);
}

/**
 * 创建友链
 * @param data 友链数据
 * @returns Promise<Friend>
 */
export function createFriend(data: CreateFriendRequest): Promise<Friend> {
  return request.post("/admin/friends", data);
}

/**
 * 更新友链
 * @param id 友链ID
 * @param data 更新数据
 * @returns Promise<Friend>
 */
export function updateFriend(id: number, data: UpdateFriendRequest): Promise<Friend> {
  return request.put(`/admin/friends/${id}`, data);
}

/**
 * 删除友链
 * @param id 友链ID
 * @returns Promise<void>
 */
export function deleteFriend(id: number): Promise<void> {
  return request.delete(`/admin/friends/${id}`);
}

