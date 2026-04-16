<!--
项目名称：JeriBlog
文件名称：callback.vue
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：Vue 组件
-->

<script setup lang="ts">
const router = useRouter()
const route = useRoute()
const { fetchUserInfo } = useUser()

onMounted(async () => {
  // 只在客户端执行
  if (!process.client) return

  const token = route.query.token as string
  const refreshToken = route.query.refresh_token as string
  const redirect = route.query.redirect as string

  if (token && refreshToken) {
    // 保存 Token
    setTokens(token, refreshToken)

    // 获取用户信息
    await fetchUserInfo()

    // 跳转回原页面，如果有 redirect 参数则使用，否则回首页
    if (redirect) {
      router.push(decodeURIComponent(redirect))
    } else {
      router.push('/')
    }
  } else {
    // 失败处理：跳转首页
    router.push('/')
  }
})
</script>

<template>
  <div class="callback-page">
    <div class="loading-content">
      <i class="ri-loader-4-line spin"></i>
      <p>正在登录中，请稍候...</p>
    </div>
  </div>
</template>

<style scoped>
.callback-page {
  height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-content {
  text-align: center;
  color: var(--text-secondary);
}

.spin {
  font-size: 2rem;
  animation: spin 1s linear infinite;
  display: block;
  margin: 0 auto 1rem;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
