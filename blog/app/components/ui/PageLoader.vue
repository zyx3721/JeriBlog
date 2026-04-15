<!--
项目名称：JeriBlog
文件名称：PageLoader.vue
创建时间：2026-04-16 02:03:34

系统用户：jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面加载动画组件，显示加载进度条和动画效果，支持淡入淡出过渡
-->
<template>
  <Transition name="loader-fade">
    <div v-if="isLoading" class="page-loader">
      <div class="loader-content">
        <!-- Logo 或标题 -->
        <div class="loader-logo">
          <div class="logo-circle">
            <div class="logo-inner"></div>
          </div>
          <h2 class="loader-title">{{ title }}</h2>
        </div>

        <!-- 加载进度条 -->
        <div class="loader-progress">
          <div class="progress-bar" :style="{ width: `${progress}%` }"></div>
        </div>

        <!-- 加载提示文字 -->
        <p class="loader-text">{{ loadingText }}</p>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
interface Props {
  isLoading?: boolean
  title?: string
  progress?: number
  loadingText?: string
}

const props = withDefaults(defineProps<Props>(), {
  isLoading: true,
  title: '加载中',
  progress: 0,
  loadingText: '正在加载资源...'
})
</script>

<style scoped>
.page-loader {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  backdrop-filter: blur(10px);
}

[data-theme='dark'] .page-loader {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
}

.loader-content {
  text-align: center;
  max-width: 400px;
  padding: 40px;
}

.loader-logo {
  margin-bottom: 32px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.logo-circle {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(255, 255, 255, 0.4);
  }
  50% {
    box-shadow: 0 0 0 20px rgba(255, 255, 255, 0);
  }
}

.logo-inner {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 3px solid #fff;
  border-top-color: transparent;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loader-title {
  margin: 0;
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 2px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.loader-progress {
  width: 100%;
  height: 4px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 16px;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #fff 0%, rgba(255, 255, 255, 0.8) 100%);
  border-radius: 2px;
  transition: width 0.3s ease;
  box-shadow: 0 0 10px rgba(255, 255, 255, 0.5);
}

.loader-text {
  margin: 0;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  font-weight: 500;
  letter-spacing: 1px;
}

/* 淡入淡出过渡效果 */
.loader-fade-enter-active,
.loader-fade-leave-active {
  transition: opacity 0.5s ease;
}

.loader-fade-enter-from,
.loader-fade-leave-to {
  opacity: 0;
}
</style>
