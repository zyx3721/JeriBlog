<template>
  <NuxtLayout>
    <div class="error-page">
      <div class="error-container">
        <div class="error-code">{{ error.statusCode }}</div>
        <div class="error-message">{{ errorMessage }}</div>
        <div class="error-description">{{ errorDescription }}</div>
        <div class="error-actions">
          <button class="btn-primary" @click="handleBack">
            <i class="ri-arrow-left-line"></i>
            返回上一页
          </button>
          <NuxtLink to="/" class="btn-secondary">
            <i class="ri-home-line"></i>
            返回首页
          </NuxtLink>
        </div>
      </div>
    </div>
  </NuxtLayout>
</template>

<script setup lang="ts">
import type { NuxtError } from '#app'

const props = defineProps({
  error: {
    type: Object as () => NuxtError,
    required: true
  }
})

const errorMessage = computed(() => {
  const statusCode = props.error.statusCode
  switch (statusCode) {
    case 404:
      return '页面未找到'
    case 403:
      return '访问被拒绝'
    case 500:
      return '服务器错误'
    default:
      return '出错了'
  }
})

const errorDescription = computed(() => {
  const statusCode = props.error.statusCode
  switch (statusCode) {
    case 404:
      return '抱歉，您访问的页面不存在或已被删除'
    case 403:
      return '抱歉，您没有权限访问此页面'
    case 500:
      return '抱歉，服务器遇到了一些问题，请稍后再试'
    default:
      return props.error.message || '抱歉，发生了一些错误'
  }
})

const handleBack = () => {
  if (window.history.length > 1) {
    window.history.back()
  } else {
    navigateTo('/')
  }
}
</script>

<style scoped lang="scss">
.error-page {
  min-height: calc(100vh - 200px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.error-container {
  text-align: center;
  max-width: 600px;
  width: 100%;
}

.error-code {
  font-size: 8rem;
  font-weight: 700;
  line-height: 1;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-hover) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 1rem;

  @media screen and (max-width: 768px) {
    font-size: 6rem;
  }
}

.error-message {
  font-size: 2rem;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 1rem;

  @media screen and (max-width: 768px) {
    font-size: 1.5rem;
  }
}

.error-description {
  font-size: 1rem;
  color: var(--text-secondary);
  margin-bottom: 2rem;
  line-height: 1.6;
}

.error-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;

  button,
  a {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1.5rem;
    border-radius: 0.5rem;
    font-size: 1rem;
    font-weight: 500;
    text-decoration: none;
    transition: all 0.3s ease;
    cursor: pointer;
    border: none;

    i {
      font-size: 1.2rem;
    }
  }

  .btn-primary {
    background: var(--primary-color);
    color: white;

    &:hover {
      background: var(--primary-hover);
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }
  }

  .btn-secondary {
    background: var(--bg-secondary);
    color: var(--text-color);

    &:hover {
      background: var(--bg-hover);
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }
  }
}
</style>
