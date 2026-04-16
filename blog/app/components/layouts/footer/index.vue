<!--
项目名称：JeriBlog
文件名称：index.vue
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：Vue 组件
-->

<script lang="ts" setup>
</script>

<template>
  <footer id="footer">
    <LayoutsFooterContact />
    <LayoutsFooterNavigation />
    <LayoutsFooterCopyright />
  </footer>
</template>

<style lang="scss" scoped>
#footer {
  position: relative;
  padding-top: 1rem;
  display: flex;
  flex-direction: column;
  gap: 32px;
  align-items: center;
  width: 100%;
  background: linear-gradient(to top, var(--flec-card-bg) 65%, var(--flec-page-bg) 90%);
}
</style>