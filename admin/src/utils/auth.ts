import router from '@/router'

const ACCESS_TOKEN_KEY = 'access_token'
const REFRESH_TOKEN_KEY = 'refresh_token'

/**
 * 获取本地存储中的access token
 * @returns {string | null} access token字符串或null
 */
export const getAccessToken = (): string | null => {
  return localStorage.getItem(ACCESS_TOKEN_KEY)
}

/**
 * 获取本地存储中的refresh token
 * @returns {string | null} refresh token字符串或null
 */
export const getRefreshToken = (): string | null => {
  return localStorage.getItem(REFRESH_TOKEN_KEY)
}

/**
 * 将双token保存到本地存储
 * @param {string} accessToken access token字符串
 * @param {string} refreshToken refresh token字符串
 */
export const setTokens = (accessToken: string, refreshToken: string): void => {
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken)
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
}

/**
 * 设置access token（用于token刷新）
 * @param {string} accessToken access token字符串
 */
export const setAccessToken = (accessToken: string): void => {
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken)
}

/**
 * 从本地存储中移除双token
 */
export const removeTokens = (): void => {
  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
}

/**
 * 检查用户是否已经登录（是否有access token）
 * @returns {boolean} true表示已登录，false表示未登录
 */
export const checkAuth = (): boolean => {
  const token = getAccessToken()
  return token !== null && token !== ''
}

/**
 * 注销用户，清除token和用户信息并跳转到登录页
 */
export const logout = (): void => {
  removeTokens()
  localStorage.removeItem('userInfo')
  router.push('/login')
}