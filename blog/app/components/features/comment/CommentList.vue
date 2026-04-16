<!--
项目名称：JeriBlog
文件名称：CommentList.vue
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：Vue 组件
-->

<script setup lang="ts">
import type { Comment } from '@@/types/comment'
import CommentItem from './CommentItem.vue'

interface FlatComment {
  comment: Comment
  depth: number
}

interface Props {
  comments: FlatComment[]
}

const props = defineProps<Props>()

// 将扁平化的评论列表转换为分组结构
const groupedComments = computed(() => {
  const groups: Array<{
    parent: FlatComment
    replies: FlatComment[]
  }> = []

  let currentGroup: { parent: FlatComment; replies: FlatComment[] } | null = null

  props.comments.forEach(item => {
    if (item.depth === 0) {
      // 顶级评论，创建新�?
      currentGroup = {
        parent: item,
        replies: []
      }
      groups.push(currentGroup)
    } else {
      // 回复评论，添加到当前�?
      if (currentGroup) {
        currentGroup.replies.push(item)
      }
    }
  })

  return groups
})
</script>

<template>
  <div class="comments-list">
    <div 
      v-for="group in groupedComments"
      :key="group.parent.comment.id"
      class="comment-card"
    >
      <!-- 顶级评论 -->
      <CommentItem 
        :comment="group.parent.comment"
        :depth="group.parent.depth"
      >
        <!-- 子评论列�?-->
        <template v-if="group.replies.length > 0" #replies>
          <div class="replies-list">
            <CommentItem 
              v-for="reply in group.replies"
              :key="reply.comment.id"
              :comment="reply.comment"
              :depth="reply.depth"
            />
          </div>
        </template>
      </CommentItem>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.comments-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.comment-card {
  background: var(--flec-card-bg);
  border-radius: 8px;
  padding: 16px;
}

@media screen and (max-width: 768px) {
  .comment-card {
    padding: 12px;
  }
}
</style>


