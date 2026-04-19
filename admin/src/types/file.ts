/*
项目名称：JeriBlog
文件名称：file.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：类型定义 - file类型
*/

// 文件信息
export interface FileInfo {
  id: number
  file_name: string
  original_name: string
  file_url: string
  file_type: string
  file_size: number
  upload_type: string
  upload_time: string
  status: number // 0:未使用 1:使用中
}

// 文件列表查询
export interface FileListQuery {
  page?: number
  page_size?: number
  type?: string
}

// 文件查询参数
export interface FileQuery {
  page: number
  page_size: number
  keyword?: string
  status?: number
  upload_type?: string
}

// 文件列表响应
export interface FileListData {
  list: FileInfo[]
  total: number
  page: number
  page_size: number
}
