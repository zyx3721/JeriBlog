<!--
项目名称：JeriBlog
文件名称：AISummary.vue
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：Vue 组件
-->

<script setup lang="ts">
const props = withDefaults(defineProps<{
  summary: string
  modelName?: string
  chatTitle?: string
}>(), { modelName: 'AI', chatTitle: 'AI 摘要' })

const displayText = ref('正在生成摘要...')

function typeText(text: string) {
  let i = 0
  const type = () => {
    if (i <= text.length) {
      displayText.value = text.slice(0, i++)
      setTimeout(type, 50)
    }
  }

  let len = displayText.value.length
  const erase = () => {
    if (len > 0) {
      displayText.value = displayText.value.slice(0, --len)
      setTimeout(erase, 80)
    } else {
      type()
    }
  }
  erase()
}

onMounted(() => {
  if (props.summary) setTimeout(() => typeText(props.summary), 2000)
})
</script>

<template>
  <div v-if="summary" class="ai-summary">
    <div class="ai-title">
      <div class="ai-title-left">
        <i class="ri-sparkling-line"></i>
        <div class="ai-title-text">{{ chatTitle }}</div>
      </div>
      <div class="ai-tag">{{ modelName }}</div>
    </div>
    <div class="ai-explanation typing">{{ displayText }}</div>
  </div>
</template>

<style lang="scss" scoped>
.ai-summary {
  font-size: .9rem;
  background: var(--flec-card-bg);
  border-radius: 12px;
  padding: 8px 8px 12px;
  margin-bottom: 16px;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.ai-explanation {
  padding: 8px 12px;
  line-height: 1.4;
  color: var(--font-color);
  text-align: justify;

  &.typing::after {
    content: '';
    display: inline-block;
    width: 8px;
    height: 2px;
    margin-left: 2px;
    background: var(--font-color);
    vertical-align: bottom;
    animation: blink-underline 1s ease-in-out infinite;
    position: relative;
    bottom: 3px;
  }
}

.ai-title {
  display: flex;
  align-items: center;
  padding: 0 12px;
  user-select: none;

  .ai-title-left {
    display: flex;
    align-items: center;
    color: var(--theme-color);

    i {
      margin-right: 3px;
    }

    .ai-title-text {
      font-weight: 500;
    }
  }

  .ai-tag {
    color: var(--font-color);
    font-weight: 300;
    margin-left: auto;
  }
}

@keyframes blink-underline {

  0%,
  100% {
    opacity: 1;
  }

  50% {
    opacity: 0;
  }
}

@media (max-width: 768px) {
  .ai-summary {
    padding: 8px;
  }

  .ai-explanation {
    padding: 6px 10px;
  }

  .ai-title {
    padding: 0 10px;
  }
}
</style>
