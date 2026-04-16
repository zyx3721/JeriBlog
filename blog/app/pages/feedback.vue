<!--
项目名称：JeriBlog
文件名称：feedback.vue
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：Vue 组件
-->

<script setup lang="ts">

definePageMeta({
  showSidebar: false
})

useSeoMeta({
  title: '反馈投诉',
  description: '如果您有任何问题或建议，欢迎通过此表单向我们反馈'
})

const handleSuccess = () => {
  // 提交成功后的处理
  console.log('Feedback submitted successfully')
}
</script>

<template>
  <div id="page">
    <h1 class="page-title">反馈投诉</h1>

    <div class="feedback-content">
      <FeaturesFeedbackForm @success="handleSuccess" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#page {
  @extend .cardHover;
  align-self: flex-start;
  padding: 40px;

  .page-title {
    margin: 0 0 30px;
    font-weight: bold;
    font-size: 2rem;
    color: var(--font-color);
  }
}

.feedback-content {
  max-width: 100%;
}

@media (max-width: 768px) {
  #page {
    padding: 20px;

    .page-title {
      font-size: 1.5rem;
      margin-bottom: 20px;
    }
  }
}
</style>
