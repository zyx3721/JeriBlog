/*
项目名称：JeriBlog
文件名称：friend.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import type { FriendGroupedResponse, FriendQueryParams, FriendApplyRequest } from '@@/types/friend'
import { createApi } from './createApi'

const friendApi = createApi<FriendGroupedResponse>('')

/** 获取友链列表 */
export const getFriends = async (params?: FriendQueryParams) => {
  return friendApi.get('/friends', params)
}

/** 申请友链 */
export const applyFriend = async (data: FriendApplyRequest) => {
  return friendApi.post('/friends/apply', data)
}
