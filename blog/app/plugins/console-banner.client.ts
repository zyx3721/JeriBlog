/*
项目名称：JeriBlog
文件名称：console-banner.client.ts
创建时间：2026-05-28 21:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：TypeScript 模块
*/

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig();
  const version = config.public.appVersion;

  console.log(
    `%c FlecBlog %c v${version} %c github.com/zyx3721/JeriBlog `,
    'background: #49b1f5; color: #fff; padding: 4px 6px; border-radius: 4px 0 0 4px; font-weight: bold;',
    'background: #3a8fd4; color: #fff; padding: 4px 6px;',
    'background: #2d7ab8; color: #fff; padding: 4px 6px; border-radius: 0 4px 4px 0;'
  );
});
