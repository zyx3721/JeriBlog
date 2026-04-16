/*
项目名称：JeriBlog
文件名称：theme.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

export const isDark = ref(false)

// 客户端初始化
if (process.client) {
  // 从 localStorage 读取，或使用系统偏好
  const stored = localStorage.getItem('theme')
  if (stored) {
    isDark.value = stored === 'dark'
  } else {
    isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  
  // 监听主题变化，自动更新 DOM
  watch(isDark, (dark) => {
    document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light')
    localStorage.setItem('theme', dark ? 'dark' : 'light')
  }, { immediate: true })
}

/**
 * 切换暗黑模式
 */
export const toggleTheme = (): void => {
  isDark.value = !isDark.value
}