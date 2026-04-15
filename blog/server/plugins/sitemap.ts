// 动态生成 sitemap 路由
export default defineNitroPlugin(async (nitroApp) => {
  nitroApp.hooks.hook('sitemap:resolved', async (ctx) => {
    try {
      const config = useRuntimeConfig()
      const apiUrl = config.public.apiUrl

      // 从后端 API 获取文章列表
      const articlesRes = await $fetch<any>(`${apiUrl}/articles`).catch(() => null)

      // 添加文章路由到 sitemap
      const articles = articlesRes?.data?.list || []
      articles.forEach((article: any) => {
        ctx.urls.push({
          loc: article.url,
          lastmod: article.update_time || article.publish_time,
          changefreq: 'weekly' as const,
          priority: 0.8
        } as any)
      })

      // 添加分类路由
      const categoriesRes = await $fetch<any>(`${apiUrl}/categories`).catch(() => null)
      const categories = categoriesRes?.data?.list || []
      categories.forEach((category: any) => {
        ctx.urls.push({
          loc: category.url,
          changefreq: 'weekly' as const,
          priority: 0.6
        } as any)
      })

      // 添加标签路由
      const tagsRes = await $fetch<any>(`${apiUrl}/tags`).catch(() => null)
      const tags = tagsRes?.data?.list || []
      tags.forEach((tag: any) => {
        ctx.urls.push({
          loc: tag.url,
          changefreq: 'weekly' as const,
          priority: 0.5
        } as any)
      })

    } catch (error) {
    }
  })
})
