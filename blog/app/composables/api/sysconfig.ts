/*
项目名称：JeriBlog
文件名称：sysconfig.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import type { SettingGroupType } from '@@/types/sysconfig'
import { createApi } from './createApi'

const settingApi = createApi<Record<string, string>>('')

/** 获取指定分组的配置 */
export const getSettingGroup = async (group: SettingGroupType) => {
  return settingApi.get<Record<string, string>>(`/settings/${group}`)
}
