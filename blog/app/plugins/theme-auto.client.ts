/*
项目名称：JeriBlog
文件名称：theme-auto.client.ts
创建时间：2026-05-18 20:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

export default defineNuxtPlugin(() => {
  const { blogConfig } = useSysConfig();

  const tryInit = () => {
    initAutoSwitch({
      lightStart: blogConfig.value.theme_light_start || '06:00',
      darkStart: blogConfig.value.theme_dark_start || '18:00',
    });
  };

  tryInit();
  watch(blogConfig, tryInit);
});
