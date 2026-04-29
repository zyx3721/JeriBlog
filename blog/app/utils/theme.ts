/*
项目名称：JeriBlog
文件名称：theme.ts
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

export const isDark = ref(false);

// 客户端初始化
if (import.meta.client) {
  // 从 DOM 读取已设置的主题（由 nuxt.config.ts 内联脚本提前设置）
  const currentTheme = document.documentElement.getAttribute('data-theme');
  isDark.value = currentTheme === 'dark';

  // 监听主题变化，自动更新 DOM 和 localStorage
  watch(isDark, dark => {
    document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
    localStorage.setItem('theme', dark ? 'dark' : 'light');
  });

  // 监听系统主题变化（仅在用户未手动设置主题时生效）
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
  mediaQuery.addEventListener('change', e => {
    // 如果用户没有手动设置过主题，则跟随系统
    if (!localStorage.getItem('theme')) {
      isDark.value = e.matches;
    }
  });
}

/**
 * 切换暗黑模式
 */
export const toggleTheme = (): void => {
  isDark.value = !isDark.value;
};
