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
import HomeHeader from './HomeHeader.vue'
import PostHeader from './PostHeader.vue'

const route = useRoute()
const type = computed(() => route.meta.typeHeader as string | undefined)
</script>

<template>
    <HomeHeader v-if="type === 'home'" />
    <PostHeader v-else-if="type === 'post'" />
</template>