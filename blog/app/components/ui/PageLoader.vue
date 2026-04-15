<!--
项目名称：JeriBlog
文件名称：PageLoader.vue
创建时间：2026-04-16 02:03:34

系统用户：jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面加载动画组件，星空背景+毛玻璃质感+流畅进度条动画
-->
<template>
  <Transition name="loader-fade">
    <div v-if="isLoading" class="page-loader">
      <!-- 星空背景 -->
      <div class="stars-bg">
        <div class="stars"></div>
        <div class="stars2"></div>
        <div class="stars3"></div>
      </div>

      <!-- 毛玻璃加载框 -->
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
          <div class="progress-bar" :style="{ width: `${progress}%` }">
            <div class="progress-glow"></div>
          </div>
        </div>

        <!-- 加载提示文字 -->
        <p class="loader-text">{{ loadingText }}</p>
        <p class="loader-percent">{{ progress }}%</p>
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
  background: radial-gradient(ellipse at bottom, #1b2735 0%, #090a0f 100%);
  overflow: hidden;
}

/* 星空背景 */
.stars-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.stars,
.stars2,
.stars3 {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: transparent;
}

/* 第一层星星 - 小而快 */
.stars {
  background-image:
    radial-gradient(2px 2px at 20px 30px, #eee, rgba(0,0,0,0)),
    radial-gradient(2px 2px at 60px 70px, #fff, rgba(0,0,0,0)),
    radial-gradient(1px 1px at 50px 50px, #fff, rgba(0,0,0,0)),
    radial-gradient(1px 1px at 130px 80px, #fff, rgba(0,0,0,0)),
    radial-gradient(2px 2px at 90px 10px, #fff, rgba(0,0,0,0));
  background-repeat: repeat;
  background-size: 200px 200px;
  animation: stars-move 50s linear infinite;
  opacity: 0.6;
}

/* 第二层星星 - 中等大小和速度 */
.stars2 {
  background-image:
    radial-gradient(1px 1px at 100px 120px, #fff, rgba(0,0,0,0)),
    radial-gradient(2px 2px at 40px 70px, #fff, rgba(0,0,0,0)),
    radial-gradient(1px 1px at 150px 60px, #fff, rgba(0,0,0,0)),
    radial-gradient(1px 1px at 90px 140px, #fff, rgba(0,0,0,0));
  background-repeat: repeat;
  background-size: 250px 250px;
  animation: stars-move 100s linear infinite;
  opacity: 0.4;
}

/* 第三层星星 - 大而慢，带闪烁 */
.stars3 {
  background-image:
    radial-gradient(3px 3px at 75px 125px, #fffacd, rgba(0,0,0,0)),
    radial-gradient(2px 2px at 180px 80px, #f0e68c, rgba(0,0,0,0)),
    radial-gradient(2px 2px at 120px 160px, #fafad2, rgba(0,0,0,0));
  background-repeat: repeat;
  background-size: 300px 300px;
  animation: stars-move 150s linear infinite, twinkle 3s ease-in-out infinite;
  opacity: 0.5;
}

@keyframes stars-move {
  from {
    transform: translateY(0);
  }
  to {
    transform: translateY(-200px);
  }
}

@keyframes twinkle {
  0%, 100% {
    opacity: 0.5;
  }
  50% {
    opacity: 0.8;
  }
}

/* 毛玻璃加载框 */
.loader-content {
  position: relative;
  z-index: 1;
  text-align: center;
  max-width: 500px;
  padding: 60px 50px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow:
    0 8px 32px 0 rgba(0, 0, 0, 0.37),
    inset 0 1px 0 0 rgba(255, 255, 255, 0.1);
}

[data-theme='dark'] .loader-content {
  background: rgba(15, 23, 42, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.loader-logo {
  margin-bottom: 40px;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-15px);
  }
}

.logo-circle {
  width: 100px;
  height: 100px;
  margin: 0 auto 24px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  animation: pulse 2s ease-in-out infinite;
  box-shadow:
    0 0 0 0 rgba(255, 255, 255, 0.4),
    0 0 30px rgba(255, 255, 255, 0.2);
}

@keyframes pulse {
  0%, 100% {
    box-shadow:
      0 0 0 0 rgba(255, 255, 255, 0.4),
      0 0 30px rgba(255, 255, 255, 0.2);
  }
  50% {
    box-shadow:
      0 0 0 20px rgba(255, 255, 255, 0),
      0 0 50px rgba(255, 255, 255, 0.4);
  }
}

.logo-inner {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 4px solid rgba(255, 255, 255, 0.8);
  border-top-color: transparent;
  border-right-color: transparent;
  animation: spin 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.loader-title {
  margin: 0;
  font-size: 32px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 3px;
  text-shadow:
    0 0 10px rgba(255, 255, 255, 0.5),
    0 0 20px rgba(255, 255, 255, 0.3),
    0 2px 10px rgba(0, 0, 0, 0.3);
  animation: glow 2s ease-in-out infinite;
}

@keyframes glow {
  0%, 100% {
    text-shadow:
      0 0 10px rgba(255, 255, 255, 0.5),
      0 0 20px rgba(255, 255, 255, 0.3),
      0 2px 10px rgba(0, 0, 0, 0.3);
  }
  50% {
    text-shadow:
      0 0 20px rgba(255, 255, 255, 0.8),
      0 0 30px rgba(255, 255, 255, 0.5),
      0 2px 10px rgba(0, 0, 0, 0.3);
  }
}

.loader-progress {
  width: 100%;
  height: 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  overflow: hidden;
  margin-bottom: 20px;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.3);
  position: relative;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(
    90deg,
    rgba(56, 189, 248, 0.8) 0%,
    rgba(59, 130, 246, 0.9) 50%,
    rgba(99, 102, 241, 1) 100%
  );
  border-radius: 10px;
  transition: width 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow:
    0 0 20px rgba(59, 130, 246, 0.6),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
  position: relative;
  overflow: hidden;
}

.progress-bar::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.4),
    transparent
  );
  animation: shimmer 2s infinite;
}

@keyframes shimmer {
  to {
    left: 200%;
  }
}

.progress-glow {
  position: absolute;
  top: 50%;
  right: 0;
  width: 20px;
  height: 20px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.8) 0%, transparent 70%);
  transform: translate(50%, -50%);
  filter: blur(4px);
  animation: glow-pulse 1s ease-in-out infinite;
}

@keyframes glow-pulse {
  0%, 100% {
    opacity: 0.6;
    transform: translate(50%, -50%) scale(1);
  }
  50% {
    opacity: 1;
    transform: translate(50%, -50%) scale(1.2);
  }
}

.loader-text {
  margin: 0 0 8px 0;
  font-size: 15px;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
  letter-spacing: 1.5px;
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.loader-percent {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.95);
  font-family: 'Courier New', monospace;
  letter-spacing: 2px;
  text-shadow: 0 0 10px rgba(59, 130, 246, 0.5);
}

/* 淡入淡出过渡效果 */
.loader-fade-enter-active {
  transition: opacity 0.3s ease;
}

.loader-fade-leave-active {
  transition: opacity 0.8s ease;
}

.loader-fade-enter-from,
.loader-fade-leave-to {
  opacity: 0;
}
</style>
