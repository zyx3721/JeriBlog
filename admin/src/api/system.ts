/*
项目名称：JeriBlog
文件名称：system.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - system
*/

import request from "@/utils/request"
import type { SystemStatic, SystemDynamic } from "@/types/system"

/**
 * 获取系统静态信息
 */
export function getSystemStatic(): Promise<SystemStatic> {
  return request.get("/admin/system/static")
}

/**
 * 获取系统动态信息
 */
export function getSystemDynamic(): Promise<SystemDynamic> {
  return request.get("/admin/system/dynamic")
}
