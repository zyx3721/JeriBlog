/*
项目名称：JeriBlog
文件名称：file.ts
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：API 接口定义 - file
*/

import request from "@/utils/request";
import type { FileInfo, FileListData, FileListQuery, FileQuery } from "@/types/file";

/**
 * 上传文件响应接口
 */
export interface UploadResponse {
  file_url: string
  file_name: string
  file_size: number
}

/**
 * 上传文件
 * @param {File} file - 要上传的文件
 * @param {string} [type='image'] - 文件类型（默认为'image'）
 * @returns {Promise<UploadResponse>} 上传结果
 */
export async function uploadFile(file: File, type = 'image'): Promise<UploadResponse> {
  const formData = new FormData();
  formData.append("file", file);
  formData.append("type", type);
  try {
    return await request.post("/admin/files", formData, {
      headers: { "Content-Type": "multipart/form-data" },
      timeout: 300000 // 5分钟超时，支持大文件上传
    });
  } catch (error: any) {
    // 尝试从响应中提取详细错误信息
    if (error.response?.data?.message) {
      throw new Error(error.response.data.message);
    }
    throw error;
  }
}

/**
 * 获取文件列表
 * @param {FileQuery} params - 查询参数
 * @returns {Promise<FileListData>} 文件列表
 */
export function getFileList(params: FileQuery): Promise<FileListData> {
  return request.get("/admin/files", { params });
}

/**
 * 删除文件
 * @param {number} id - 文件ID
 * @returns {Promise<void>}
 */
export function deleteFile(id: number): Promise<void> {
  return request.delete(`/admin/files/${id}`);
}

/**
 * 文件引用信息
 */
export interface FileReference {
  type: string // 引用类型：article/user/friend/setting
  id: number // 引用对象ID
  title: string // 引用对象标题
  field: string // 引用字段
  url?: string // 引用对象链接
}

/**
 * 获取文件引用详情
 * @param {number} id - 文件ID
 * @returns {Promise<FileReference[]>} 引用列表
 */
export function getFileReferences(id: number): Promise<FileReference[]> {
  return request.get(`/admin/files/${id}/references`);
}
