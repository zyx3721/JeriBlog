/*
项目名称：JeriBlog
文件名称：tools.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - tools
*/

import request from '@/utils/request'
import type { FetchLinkRequest, LinkInfo, ParseVideoRequest, VideoInfo } from '@/types/tools'

/**
 * 工具API模块 - 用于视频解析、链接元数据获取等通用工具功能
 */

/**
 * 根据URL获取链接信息
 * @param data 链接URL
 * @returns Promise<LinkInfo>
 */
export function fetchLinkInfo(data: FetchLinkRequest): Promise<LinkInfo> {
  return request.post("/admin/tools/fetch-linkmeta", data)
}

/**
 * 解析视频URL
 * @param data 视频URL
 * @returns Promise<VideoInfo>
 */
export function parseVideo(data: ParseVideoRequest): Promise<VideoInfo> {
  return request.post("/admin/tools/parse-video", data)
}

/**
 * 下载图片
 * @param data 图片URL
 * @returns Promise<{ content_type: string, content_length: number, data: string }>
 */
export function downloadImage(data: { url: string }): Promise<{ content_type: string, content_length: number, data: string }> {
  return request.post("/admin/tools/download-image", data)
}

/**
 * 解析音乐信息（代理Meting API）
 * @param params 音乐参数 { server: string, type: string, id: string }
 * @returns Promise<any[]>
 */
export function parseMusic(params: { server: string, type: string, id: string }): Promise<any[]> {
  return request.get("/admin/tools/parse-music/", { params })
}