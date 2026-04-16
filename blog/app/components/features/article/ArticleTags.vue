<!--
项目名称：JeriBlog
文件名称：ArticleTags.vue
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：Vue 组件
-->

<script lang="ts" setup>
defineProps<{ article: { tags: Array<{ id: number; name: string; url: string }> } }>()
</script>

<template>
    <div class="article-tags" v-if="article.tags?.length">
        <a v-for="tag in article.tags" :key="tag.id" :href="tag.url">{{ tag.name }}</a>
    </div>
</template>

<style lang="scss" scoped>
.article-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;

    a {
        padding: 1px 8px;
        border: 1px solid var(--theme-color);
        border-radius: 8px;
        color: var(--theme-color);
        text-decoration: none;
        font-size: 0.85rem;

        &:hover {
            background: var(--theme-color);
            color: #fff;
        }
    }
}
</style>
