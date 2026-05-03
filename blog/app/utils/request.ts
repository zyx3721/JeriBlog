/*
项目名称：JeriBlog
文件名称：request.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import { $fetch } from 'ofetch';
import type { FetchOptions } from 'ofetch';
import { accessToken, setAccessToken, logout } from './auth';
import type { ApiResponse } from '@@/types/request';

type HttpMethod =
  | 'GET'
  | 'HEAD'
  | 'PATCH'
  | 'POST'
  | 'PUT'
  | 'DELETE'
  | 'CONNECT'
  | 'OPTIONS'
  | 'TRACE';

// 获取 API baseURL
const getBaseURL = () => useRuntimeConfig().public.apiUrl as string;

// Token 刷新状态
let isRefreshing = false;
let refreshPromise: Promise<boolean> | null = null;

// 刷新 Token
const doRefreshToken = async (): Promise<boolean> => {
  try {
    const res = await $fetch<ApiResponse<{ access_token: string }>>('/auth/refresh', {
      method: 'POST',
      baseURL: getBaseURL(),
      credentials: 'include', // 发送 Cookie
    });
    if (res.code === 0 && res.data) {
      setAccessToken(res.data.access_token);
      return true;
    }
    return false;
  } catch {
    return false;
  }
};

// 封装请求，支持自动 Token 和 401 刷新
/* eslint-disable @typescript-eslint/no-explicit-any -- 通用请求封装，返回类型由调用方指定 */
export async function apiRequest<T = any>(
  url: string,
  options: Omit<FetchOptions, 'method'> & { method?: HttpMethod; _retry?: boolean } = {}
): Promise<T> {
  const config = useRuntimeConfig();
  const headers: Record<string, string> = {
    ...((options.headers as Record<string, string>) || {}),
  };

  if (accessToken.value && url !== '/auth/refresh') {
    headers['Authorization'] = `Bearer ${accessToken.value}`;
  }

  try {
    const response = await $fetch<T>(url, {
      ...options,
      baseURL: config.public.apiUrl,
      headers,
      credentials: 'include', // 发送 Cookie
    } as any);

    // 检查业务状态码
    if (response && typeof response === 'object' && 'code' in response) {
      const apiResponse = response as any;
      if (apiResponse.code !== 0) {
        // 业务错误，抛出异常
        const error = new Error(apiResponse.message || '请求失败');
        (error as any).response = { status: apiResponse.code, data: apiResponse };
        throw error;
      }
    }

    return response;
  } catch (error: any) {
    // 401 自动刷新 token
    if (error?.response?.status === 401 && !options._retry) {
      // 避免并发刷新
      if (!isRefreshing) {
        isRefreshing = true;
        refreshPromise = doRefreshToken().finally(() => {
          isRefreshing = false;
        });
      }

      const success = await refreshPromise;
      if (success) {
        return apiRequest<T>(url, { ...options, _retry: true });
      }
      logout();
    }
    throw error;
  }
}
/* eslint-enable @typescript-eslint/no-explicit-any */

/**
 * GET请求
 */
/* eslint-disable @typescript-eslint/no-explicit-any -- 通用请求封装，返回类型由调用方指定 */
export async function get<T = any>(
  url: string,
  options: Omit<FetchOptions, 'method'> = {}
): Promise<T> {
  return await apiRequest<T>(url, { ...options, method: 'GET' });
}
/* eslint-enable @typescript-eslint/no-explicit-any */

/**
 * POST请求
 */
/* eslint-disable @typescript-eslint/no-explicit-any -- 通用请求封装，返回类型和请求体由调用方指定 */
export async function post<T = any>(
  url: string,
  body?: any,
  options: Omit<FetchOptions, 'method'> = {}
): Promise<T> {
  return await apiRequest<T>(url, { ...options, method: 'POST', body });
}
/* eslint-enable @typescript-eslint/no-explicit-any */

/**
 * PUT请求
 */
/* eslint-disable @typescript-eslint/no-explicit-any -- 通用请求封装，返回类型和请求体由调用方指定 */
export async function put<T = any>(
  url: string,
  body?: any,
  options: Omit<FetchOptions, 'method'> = {}
): Promise<T> {
  return await apiRequest<T>(url, { ...options, method: 'PUT', body });
}
/* eslint-enable @typescript-eslint/no-explicit-any */

/**
 * PATCH请求
 */
/* eslint-disable @typescript-eslint/no-explicit-any -- 通用请求封装，返回类型和请求体由调用方指定 */
export async function patch<T = any>(
  url: string,
  body?: any,
  options: Omit<FetchOptions, 'method'> = {}
): Promise<T> {
  return await apiRequest<T>(url, { ...options, method: 'PATCH', body });
}
/* eslint-enable @typescript-eslint/no-explicit-any */

/**
 * DELETE请求
 */
/* eslint-disable @typescript-eslint/no-explicit-any -- 通用请求封装，返回类型由调用方指定 */
export async function del<T = any>(
  url: string,
  options: Omit<FetchOptions, 'method'> = {}
): Promise<T> {
  return await apiRequest<T>(url, { ...options, method: 'DELETE' });
}
/* eslint-enable @typescript-eslint/no-explicit-any */
