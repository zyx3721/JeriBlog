<!--
项目名称：JeriBlog
文件名称：Sidebar.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：公共组件 - Sidebar组件
-->

<template>
  <div class="sidebar" :class="{ 'is-collapse': isCollapse }">
    <div class="logo">
      <span v-show="!isCollapse">Jeri 管理系统</span>
    </div>
    <el-menu :default-active="route.path" :collapse="isCollapse" background-color="#304156" text-color="#bfcbd9"
      active-text-color="#409eff" router @select="handleMenuSelect">
      <el-menu-item index="/" @click="handleItemClick('/')">
        <i class="ri-dashboard-2-line ri-lg"></i>
        <template #title><span>仪表盘</span></template>
      </el-menu-item>

      <el-sub-menu index="content">
        <template #title>
          <i class="ri-layout-2-line ri-lg"></i>
          <span>内容管理</span>
        </template>
        <el-menu-item index="/articles" @click="handleItemClick('/articles')">
          <i class="ri-article-line ri-lg"></i>
          <template #title>文章管理</template>
        </el-menu-item>
        <el-menu-item index="/moments" @click="handleItemClick('/moments')">
          <i class="ri-chat-3-line ri-lg"></i>
          <template #title>动态管理</template>
        </el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="interaction">
        <template #title>
          <i class="ri-chat-2-line ri-lg"></i>
          <span>互动管理</span>
        </template>
        <el-menu-item index="/friends" @click="handleItemClick('/friends')">
          <i class="ri-links-line ri-lg"></i>
          <template #title>友链管理</template>
        </el-menu-item>
        <el-menu-item index="/comments" @click="handleItemClick('/comments')">
          <i class="ri-message-3-line ri-lg"></i>
          <template #title>评论管理</template>
        </el-menu-item>
        <el-menu-item index="/rssfeed" @click="handleItemClick('/rssfeed')">
          <i class="ri-rss-line ri-lg"></i>
          <template #title>RSS订阅</template>
        </el-menu-item>
        <el-menu-item index="/feedback" @click="handleItemClick('/feedback')">
          <i class="ri-feedback-line ri-lg"></i>
          <template #title>反馈投诉</template>
        </el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="management">
        <template #title>
          <i class="ri-admin-line ri-lg"></i>
          <span>系统管理</span>
        </template>
        <el-menu-item index="/users" @click="handleItemClick('/users')">
          <i class="ri-team-line ri-lg"></i>
          <template #title>用户管理</template>
        </el-menu-item>
        <el-menu-item index="/files" @click="handleItemClick('/files')">
          <i class="ri-folder-image-line ri-lg"></i>
          <template #title>文件管理</template>
        </el-menu-item>
        <el-menu-item index="/menus" @click="handleItemClick('/menus')">
          <i class="ri-menu-line ri-lg"></i>
          <template #title>菜单管理</template>
        </el-menu-item>
        <el-menu-item index="/visits" @click="handleItemClick('/visits')">
          <i class="ri-file-list-3-line ri-lg"></i>
          <template #title>访问日志</template>
        </el-menu-item>
        <el-menu-item index="/systems" @click="handleItemClick('/systems')">
          <i class="ri-information-line ri-lg"></i>
          <template #title>系统信息</template>
        </el-menu-item>
        <el-menu-item index="/settings" @click="handleItemClick('/settings')">
          <i class="ri-settings-3-line ri-lg"></i>
          <template #title>系统设置</template>
        </el-menu-item>
      </el-sub-menu>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

defineProps<{
  isCollapse: boolean
}>()

const emit = defineEmits(['menu-click'])

// 菜单选择事件处理
const handleMenuSelect = () => {
  emit('menu-click')
}

// 菜单项点击事件处理（支持重复点击刷新）
const handleItemClick = (path: string) => {
  // 如果点击的是当前路由，则刷新页面
  if (route.path === path) {
    // 通过改变路由的 query 参数来触发组件重新加载
    router.replace({
      path: path,
      query: { _t: Date.now() }
    }).then(() => {
      // 立即移除 query 参数，保持 URL 干净
      router.replace({ path: path })
    })
  }
}
</script>

<style scoped lang="scss">
.sidebar {
  height: 100%;
  background: linear-gradient(180deg, #1f2937 0%, #304156 50%, #1a2332 100%);
  overflow-y: auto;
  overflow-x: hidden;

  // 隐藏滚动条但保持滚动功能
  &::-webkit-scrollbar {
    width: 0;
    height: 0;
  }

  scrollbar-width: none; // Firefox
  -ms-overflow-style: none; // IE/Edge

  .logo {
    height: 60px;
    padding: 10px 0;
    margin-bottom: 10px;
    display: flex;
    align-items: center;
    justify-content: center;

    span {
      color: #fff;
      font-size: 16px;
      font-weight: 600;
    }
  }

  &.is-collapse {
    .logo {
      padding: 10px;

      span {
        display: none;
      }
    }
  }

  :deep(.el-menu) {
    border-right: none;
    background: transparent;

    .el-menu-item,
    .el-sub-menu__title {
      font-size: 16px;

      i {
        margin-right: 10px;
      }
    }
  }
}
</style>