/*
项目名称：JeriBlog
文件名称：tools.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：类型定义 - tools类型
*/

// 工具API类型定义

// 获取链接信息请求
export interface FetchLinkRequest {
  url: string
}

// 链接信息响应
export interface LinkInfo {
  url: string
  title: string
  description: string
  favicon: string
}

// 解析视频请求
export interface ParseVideoRequest {
  url: string
}

// 视频信息响应
export interface VideoInfo {
  platform: string
  video_id: string
}

// 下载图片请求
export interface DownloadImageRequest {
  url: string
}

// 下载图片响应
export interface DownloadImageResponse {
  content_type: string
  content_length: number
  data: string
}