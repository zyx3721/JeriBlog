/*
项目名称：JeriBlog
文件名称：auth.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

import { ref, computed } from 'vue'

// 博客端专用的 token key，避免与管理后台冲突
const BLOG_ACCESS_TOKEN_KEY = 'blog_access_token'
const BLOG_REFRESH_TOKEN_KEY = 'blog_refresh_token'

// Token 状态
export const accessToken = ref<string | null>(null)
export const refreshToken = ref<string | null>(null)

// 响应式登录状态
export const isLoggedIn = computed(() => !!accessToken.value && accessToken.value !== '')

// 从 localStorage 初始化 token（仅客户端）
if (process.client) {
  accessToken.value = localStorage.getItem(BLOG_ACCESS_TOKEN_KEY)
  refreshToken.value = localStorage.getItem(BLOG_REFRESH_TOKEN_KEY)
}

/**
 * 设置双token
 */
export const setTokens = (access: string, refresh: string): void => {
  accessToken.value = access
  refreshToken.value = refresh

  // 同步到 localStorage（仅客户端）
  if (process.client) {
    localStorage.setItem(BLOG_ACCESS_TOKEN_KEY, access)
    localStorage.setItem(BLOG_REFRESH_TOKEN_KEY, refresh)
  }
}

/**
 * 设置access token（用于token刷新）
 */
export const setAccessToken = (access: string): void => {
  accessToken.value = access

  // 同步到 localStorage（仅客户端）
  if (process.client) {
    localStorage.setItem(BLOG_ACCESS_TOKEN_KEY, access)
  }
}

/**
 * 登出操作（清除双token）
 */
export const logout = (): void => {
  accessToken.value = null
  refreshToken.value = null
  if (process.client) {
    localStorage.removeItem(BLOG_ACCESS_TOKEN_KEY)
    localStorage.removeItem(BLOG_REFRESH_TOKEN_KEY)
  }
}

/**
 * 获取响应式的登录状态
 */
export const useAuth = () => isLoggedIn
