<!--
项目名称：JeriBlog
文件名称：CommentEmpty.vue
创建时间：2026-04-16 15:10:34

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：Vue 组件
-->

<template>
  <div class="comments-empty">
    <i class="ri-chat-off-line"></i>
    <p>暂无评论，快来发表第一条评论吧！</p>
  </div>
</template>

<style lang="scss" scoped>
.comments-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--flec-lightText);
  text-align: center;

  i {
    font-size: 4rem;
    margin-bottom: 15px;
    opacity: 0.5;
  }

  p {
    font-size: 1rem;
    margin: 0;
  }
}

// 响应式设�?
@media screen and (max-width: 768px) {
  .comments-empty {
    padding: 40px 20px;

    i {
      font-size: 3rem;
    }
  }
}
</style>


