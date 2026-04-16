/*
项目名称：JeriBlog
文件名称：sysconfig.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - sysconfig
*/

import request from "@/utils/request";
import type { SettingGroupType } from "@/types/sysconfig";

// 获取指定分组的配置
export const getSettingGroup = (
  group: SettingGroupType
): Promise<Record<string, string>> => {
  return request.get(`/admin/settings/${group}`);
};

// 更新指定分组的配置
export const updateSettingGroup = (
  group: SettingGroupType,
  data: Record<string, string>
) => {
  return request.patch(`/admin/settings/${group}`, data);
};
