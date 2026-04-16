<!--
项目名称：JeriBlog
文件名称：Logo.vue
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：Vue 组件
-->

<script setup lang="ts">
const { blogConfig } = useSysConfig()
</script>

<template>
  <div class="logo-container">
    <a href="/" class="logo brighten" aria-label="返回首页">
      {{ blogConfig.subtitle }}
    </a>
  </div>
</template>

<style lang="scss" scoped>
.logo-container {
  .logo {
    font-weight: bold;
    font-size: 1.25rem;
  }
}

@media screen and (max-width: 768px) {
  .logo-container {
    display: none !important;
  }
}
</style>
