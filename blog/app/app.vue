<!--
项目名称：JeriBlog
文件名称：app.vue
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：Vue 组件
-->

<script setup lang="ts">
import { getMenus } from '@/composables/api/menu'
import { getCategories } from '@/composables/api/category'
import { getTags } from '@/composables/api/tag'
import { getSiteStats } from '@/composables/api/stats'
import { getSettingGroup } from '@/composables/api/sysconfig'

const { toasts } = useToast()
const { showLoginModal } = useLoginModal()
const { showBindEmailModal, triggerGlobal, onBindSuccess } = useBindEmail()

// 全局加载状态
const { isLoading, hasInitialized, progress, loadingText, setLoading, setInitialized, setProgress, setLoadingText } = useAppLoading()

// 全局数据
const { blogConfig, basicConfig, oauthConfig, uploadConfig } = useSysConfig()
const { menus } = useMenus()
const { categories, total: categoriesTotal } = useCategories()
const { tags, total: tagsTotal } = useTags()
const { siteStats } = useStats()

// 使用SSR获取全局数据
const { data: globalData } = await useAsyncData('global-data', async () => {
  const [basicConfigData, blogConfigData, oauthConfigData, uploadConfigData, menusData, categoriesData, tagsData, statsData] = await Promise.all([
    getSettingGroup('basic'),
    getSettingGroup('blog'),
    getSettingGroup('oauth'),
    getSettingGroup('upload'),
    getMenus(),
    getCategories(),
    getTags(),
    getSiteStats()
  ])

  // 处理配置数据
  const processConfig = (config: any, prefix: string) => {
    const processed: Record<string, string> = {}
    Object.entries(config).forEach(([key, value]) => {
      if (key.startsWith(`${prefix}.`)) {
        processed[key.substring(prefix.length + 1)] = value as string
      }
    })
    return processed
  }

  return {
    basicConfig: processConfig(basicConfigData, 'basic'),
    blogConfig: processConfig(blogConfigData, 'blog'),
    oauthConfig: processConfig(oauthConfigData, 'oauth'),
    uploadConfig: processConfig(uploadConfigData, 'upload'),
    menus: menusData || [],
    categories: categoriesData.list,
    categoriesTotal: categoriesData.total,
    tags: tagsData.list,
    tagsTotal: tagsData.total,
    stats: statsData
  }
})

// 初始化全局数据
if (globalData.value) {
  basicConfig.value = globalData.value.basicConfig
  blogConfig.value = globalData.value.blogConfig
  oauthConfig.value = globalData.value.oauthConfig
  uploadConfig.value = globalData.value.uploadConfig
  menus.value = globalData.value.menus
  categories.value = globalData.value.categories
  tags.value = globalData.value.tags
  siteStats.value = globalData.value.stats
  if (globalData.value.categoriesTotal !== undefined) {
    categoriesTotal.value = globalData.value.categoriesTotal
  }
  if (globalData.value.tagsTotal !== undefined) {
    tagsTotal.value = globalData.value.tagsTotal
  }
}

// 全局路由切换时触发邮箱绑定提示
const router = useRouter()
router.afterEach(() => {
  triggerGlobal()
})

// 背景图片
const bgImage = computed(() => blogConfig.value.background_image || '/bg.webp')

// 刷新时恢复滚动位置
onMounted(() => {
  const key = 'scroll-y'
  const nav = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming
  if (nav?.type === 'reload') {
    const y = +(sessionStorage.getItem(key) || 0)
    if (y > 0) setTimeout(() => window.scrollTo(0, y), 100)
  }
  let t: ReturnType<typeof setTimeout>
  const save = () => sessionStorage.setItem(key, '' + window.scrollY)
  window.addEventListener('scroll', () => { clearTimeout(t); t = setTimeout(save, 200) }, { passive: true })
  window.addEventListener('pagehide', save)

  // 首次加载：模拟加载进度
  if (!hasInitialized.value) {
    const loadingSteps = [
      { progress: 20, text: '正在加载配置...', delay: 200 },
      { progress: 40, text: '正在加载菜单...', delay: 200 },
      { progress: 60, text: '正在加载分类和标签...', delay: 200 },
      { progress: 80, text: '正在加载统计数据...', delay: 200 },
      { progress: 100, text: '加载完成！', delay: 200 }
    ]

    const updateProgress = async (index: number) => {
      if (index >= loadingSteps.length) {
        setTimeout(() => {
          setLoading(false)
          setInitialized(true)
        }, 300)
        return
      }

      const step = loadingSteps[index]
      setProgress(step.progress)
      setLoadingText(step.text)

      await new Promise(resolve => setTimeout(resolve, step.delay))
      updateProgress(index + 1)
    }

    updateProgress(0)
  } else {
    // 已初始化过，直接关闭加载动画
    setLoading(false)
  }
})

// SEO Meta
useSeoMeta({
  description: () => blogConfig.value.description,
  keywords: () => blogConfig.value.keywords,
  author: () => basicConfig.value.author,
  // Open Graph
  ogTitle: () => blogConfig.value.title,
  ogDescription: () => blogConfig.value.description,
  ogImage: () => blogConfig.value.favicon,
  ogType: 'website',
  ogSiteName: () => blogConfig.value.title,
  // Twitter Card
  twitterCard: 'summary_large_image',
  twitterTitle: () => blogConfig.value.title,
  twitterDescription: () => blogConfig.value.description,
  twitterImage: () => blogConfig.value.favicon,
})

// 页面标题模板和 favicon
const route = useRoute()
const siteTitle = computed(() => blogConfig.value.title)

useHead({
  titleTemplate: (title): string | null => {
    // 首页特殊处理：显示"网站标题 - 网站副标题"
    if (route.path === '/') {
      const subtitle = blogConfig.value.subtitle
      return subtitle ? `${siteTitle.value} - ${subtitle}` : siteTitle.value || null
    }

    // 其他页面：显示"页面标题 | 网站标题"
    const pageTitle = title || (route.meta.title as string)
    if (pageTitle) return `${pageTitle} | ${siteTitle.value}`
    return siteTitle.value || null
  },
  link: [
    { rel: 'icon', type: 'image/x-icon', href: blogConfig.value.favicon || '/favicon.ico?v=2' },
    // PWA Manifest
    { rel: 'manifest', href: '/manifest.json' },
    // RSS/Atom 订阅
    {
      rel: 'alternate',
      type: 'application/rss+xml',
      title: `${blogConfig.value.title} - RSS 2.0 Feed`,
      href: '/rss.xml'
    },
    {
      rel: 'alternate',
      type: 'application/atom+xml',
      title: `${blogConfig.value.title} - Atom Feed`,
      href: '/atom.xml'
    }
  ],
  meta: computed(() => [
    { name: 'description', content: blogConfig.value.description },
    { name: 'keywords', content: blogConfig.value.keywords },
    { name: 'author', content: blogConfig.value.author },
    // PWA 主题色
    { name: 'theme-color', content: '#f7f7f7' },
    { name: 'apple-mobile-web-app-capable', content: 'yes' },
    { name: 'apple-mobile-web-app-status-bar-style', content: 'default' }
  ]),
  script: [
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'WebSite',
        name: blogConfig.value.title,
        description: blogConfig.value.description,
      })
    }
  ]
})
</script>

<template>
  <!-- 页面加载动画 -->
  <UiPageLoader
    :is-loading="isLoading"
    :title="blogConfig.title || '加载中'"
    :progress="progress"
    :loading-text="loadingText"
  />

  <!-- Canvas 背景动画 -->
  <UiCanvasBackground />

  <!-- 背景图片 -->
  <div class="web_bg" :style="{ backgroundImage: `url(${bgImage})` }"></div>

  <!-- Nuxt 布局和页面系统 -->
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>

  <!-- Toast 消息提示 -->
  <UiToast v-for="toast in toasts" :key="toast.id" v-bind="toast" />

  <!-- 登录弹窗 -->
  <FeaturesModalsLoginModal v-model="showLoginModal" />

  <!-- 邮箱绑定弹窗 -->
  <FeaturesModalsBindEmailModal v-model="showBindEmailModal" @success="onBindSuccess" />

  <!-- 右键菜单 -->
  <UiContextMenu />
</template>

<style scoped>
.web_bg {
  position: fixed;
  width: 100%;
  height: 100%;
  z-index: -50;
  background-position: center;
  background-size: cover;
  background-repeat: no-repeat;
}

[data-theme='dark'] .web_bg::before {
  position: absolute;
  width: 100%;
  height: 100%;
  background-color: #121212b0;
  content: '';
}
</style>
