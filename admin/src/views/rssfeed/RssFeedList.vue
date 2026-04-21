<!--
项目名称：JeriBlog
文件名称：RssFeedList.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面组件 - RssFeedList页面
-->

<template>
  <common-list title="RSS订阅" :data="articleList" :loading="loading" :total="total" v-model:page="queryParams.page"
    v-model:page-size="queryParams.page_size" :show-create="false" @refresh="fetchArticles"
    @update:page="fetchArticles" @update:pageSize="fetchArticles">
    <!-- 搜索表单 -->
    <template #toolbar-before>
      <div class="search-form rss-search">
        <el-input
          v-model="queryParams.keyword"
          placeholder="搜索文章标题..."
          clearable
          style="width: 240px"
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        />
        <el-select
          v-model="statusFilter"
          placeholder="状态"
          clearable
          style="width: 120px"
          @change="handleStatusChange"
        >
          <el-option label="未读" value="unread" />
          <el-option label="已读" value="read" />
          <el-option label="已删除" value="deleted" />
        </el-select>
        <el-select
          v-model="queryParams.friend_id"
          placeholder="来源"
          clearable
          filterable
          style="width: 180px"
          @change="handleSearch"
        >
          <el-option
            v-for="friend in friendList"
            :key="friend.id"
            :label="friend.name"
            :value="friend.id"
          />
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
    </template>

    <!-- 额外按钮 -->
    <template #toolbar-after>
      <el-button type="primary" @click="openSubscriberDialog">
        本站订阅
      </el-button>
      <el-button type="warning" :loading="refreshing" @click="handleRefreshRss" v-if="isSuperAdmin">
        <el-icon v-if="!refreshing">
          <Download />
        </el-icon>
        {{ refreshing ? '抓取中...' : '立即抓取RSS' }}
      </el-button>
      <el-badge :value="unreadCount" :hidden="unreadCount === 0" :max="99" class="unread-badge">
        <el-button type="success" :disabled="unreadCount === 0" @click="handleMarkAllRead"
          v-if="isSuperAdmin">
          全部已读
        </el-button>
      </el-badge>
    </template>

    <!-- 表格列 -->
    <el-table-column label="状态" width="80" align="center">
      <template #default="{ row }">
        <el-tag v-if="row.is_deleted" type="danger" size="small">
          已删除
        </el-tag>
        <el-tag v-else :type="row.is_read ? 'info' : 'warning'" size="small">
          {{ row.is_read ? '已读' : '未读' }}
        </el-tag>
      </template>
    </el-table-column>

    <el-table-column label="文章标题" min-width="300" align="center">
      <template #default="{ row }">
        <div style="display: flex; align-items: center; gap: 8px; justify-content: center;">
          <a :href="row.link" target="_blank" class="article-link" :class="{ read: row.is_read }">
            {{ row.title }}
          </a>
          <el-tag v-if="row.update_type && row.update_type.includes('title')" type="success" size="small">
            标题已更新
          </el-tag>
          <el-tag v-if="row.update_type && row.update_type.includes('link')" type="warning" size="small">
            链接已更新
          </el-tag>
          <el-tag v-if="row.update_type && row.update_type.includes('published_at')" type="info" size="small">
            发布时间已更新
          </el-tag>
        </div>
      </template>
    </el-table-column>

    <el-table-column label="来源" width="180" align="center">
      <template #default="{ row }">
        <a :href="row.friend_url" target="_blank" class="friend-link">
          {{ row.friend_name }}
        </a>
      </template>
    </el-table-column>

    <el-table-column label="发布时间" width="180" align="center">
      <template #default="{ row }">
        <span v-if="row.published_at">{{ formatDateTime(row.published_at) }}</span>
        <span v-else style="color: #999">-</span>
      </template>
    </el-table-column>

    <el-table-column label="抓取时间" width="180" align="center">
      <template #default="{ row }">
        {{ formatDateTime(row.created_at) }}
      </template>
    </el-table-column>

    <el-table-column label="操作" width="120" align="center" fixed="right">
      <template #default="{ row }">
        <el-button v-if="!row.is_read && isSuperAdmin" type="primary" link size="small"
          @click="handleMarkRead(row)">
          标记已读
        </el-button>
        <span v-else style="color: #999">-</span>
      </template>
    </el-table-column>
  </common-list>

  <!-- 本站订阅弹窗 -->
  <el-dialog v-model="subscriberDialogVisible" title="本站订阅者" width="700px" destroy-on-close>
    <el-table :data="subscriberList" v-loading="subscriberLoading" border style="width: 100%">
      <el-table-column label="邮箱地址" min-width="240">
        <template #default="{ row }">
          <div style="display: flex; align-items: center; gap: 8px">
            <el-icon size="16" color="#409eff">
              <Message />
            </el-icon>
            <span>{{ row.email }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.active ? 'success' : 'info'" size="small">
            {{ row.active ? '活跃' : '已退订' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="订阅时间" width="170" align="center">
        <template #default="{ row }">
          {{ formatDateTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="80" align="center">
        <template #default="{ row }">
          <el-button type="danger" link size="small" @click="handleDeleteSubscriber(row.id)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <div style="margin-top: 16px; display: flex; justify-content: flex-end;">
      <el-pagination v-model:current-page="subscriberQuery.page" v-model:page-size="subscriberQuery.page_size"
        :total="subscriberTotal" :page-sizes="[10, 20, 50]" layout="total, sizes, prev, pager, next"
        @current-change="fetchSubscribers" @size-change="fetchSubscribers" />
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Message, Download } from '@element-plus/icons-vue'
import CommonList from '@/components/common/CommonList.vue'
import type { RssArticle, RssArticleQuery } from '@/types/rssfeed'
import type { User } from '@/types/user'
import type { Subscriber } from '@/types/subscriber'
import type { Friend } from '@/types/friend'
import { getRssArticles, markRssArticleRead, markAllRssArticlesRead, refreshAllRssFeeds } from '@/api/rssfeed'
import { getSubscribers, deleteSubscriber } from '@/api/subscriber'
import { getFriends } from '@/api/friend'
import { formatDateTime } from '@/utils/date'

const userInfo = computed(() => {
  const stored = localStorage.getItem('userInfo')
  if (stored) {
    return JSON.parse(stored) as User
  }
  return null
})

const isSuperAdmin = computed(() => userInfo.value?.role === 'super_admin')

const loading = ref(false)
const refreshing = ref(false)
const articleList = ref<RssArticle[]>([])
const total = ref(0)
const unreadCount = ref(0)
const queryParams = ref<RssArticleQuery>({ page: 1, page_size: 20 })
const friendList = ref<Friend[]>([])
const statusFilter = ref<string>('')

/**
 * 获取友链列表（用于来源筛选）
 */
const fetchFriends = async () => {
  try {
    const result = await getFriends({ page: 1, page_size: 1000 })
    // 只显示配置了RSS的友链
    friendList.value = result.list.filter(f => f.rss_url && f.rss_url.trim() !== '')
  } catch {
    // 静默失败，不影响主功能
  }
}

/**
 * 获取RSS文章列表
 */
const fetchArticles = async () => {
  loading.value = true
  try {
    const result = await getRssArticles(queryParams.value)
    articleList.value = result.list
    total.value = result.total
    unreadCount.value = result.unread_count
  } catch {
    ElMessage.error('获取RSS文章列表失败')
  } finally {
    loading.value = false
  }
}

/**
 * 状态筛选变化处理
 */
const handleStatusChange = () => {
  if (statusFilter.value === 'unread') {
    queryParams.value.is_read = false
    queryParams.value.is_deleted = false
  } else if (statusFilter.value === 'read') {
    queryParams.value.is_read = true
    queryParams.value.is_deleted = false
  } else if (statusFilter.value === 'deleted') {
    queryParams.value.is_read = undefined
    queryParams.value.is_deleted = true
  } else {
    queryParams.value.is_read = undefined
    queryParams.value.is_deleted = undefined
  }
  handleSearch()
}

/**
 * 搜索
 */
const handleSearch = () => {
  queryParams.value.page = 1
  fetchArticles()
}

/**
 * 重置搜索
 */
const handleReset = () => {
  queryParams.value = { page: 1, page_size: 20 }
  statusFilter.value = ''
  fetchArticles()
}

/**
 * 标记单篇文章已读
 */
const handleMarkRead = async (article: RssArticle) => {
  try {
    await markRssArticleRead(article.id)
    // 在列表中找到对应文章并更新
    const index = articleList.value.findIndex(a => a.id === article.id)
    if (index !== -1 && articleList.value[index]) {
      articleList.value[index].is_read = true
      articleList.value[index].update_type = ''
    }
    unreadCount.value = Math.max(0, unreadCount.value - 1)
    ElMessage.success('已标记为已读')
  } catch {
    ElMessage.error('操作失败')
  }
}

/**
 * 全部标记已读
 */
const handleMarkAllRead = async () => {
  try {
    await ElMessageBox.confirm('确定要将所有未读文章标记为已读吗？', '提示', { type: 'warning' })
    const result = await markAllRssArticlesRead()
    ElMessage.success(`已标记 ${result.affected} 篇文章为已读`)
    fetchArticles()
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) {
      ElMessage.error(error.message)
    }
  }
}

/**
 * 立即刷新RSS订阅源
 */
const handleRefreshRss = async () => {
  try {
    await ElMessageBox.confirm('确定要立即抓取所有友链的RSS订阅内容吗？', '提示', {
      type: 'info',
      confirmButtonText: '立即抓取',
      cancelButtonText: '取消'
    })
    refreshing.value = true
    await refreshAllRssFeeds()
    ElMessage.success('RSS订阅源刷新成功，正在重新加载列表...')
    // 延迟1秒后刷新列表，确保数据已入库
    setTimeout(() => {
      fetchArticles()
    }, 1000)
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) {
      ElMessage.error(error.message)
    }
  } finally {
    refreshing.value = false
  }
}

// 订阅者相关
const subscriberDialogVisible = ref(false)
const subscriberLoading = ref(false)
const subscriberList = ref<Subscriber[]>([])
const subscriberTotal = ref(0)
const subscriberQuery = reactive({ page: 1, page_size: 10 })

/**
 * 打开订阅者弹窗
 */
const openSubscriberDialog = () => {
  subscriberDialogVisible.value = true
  subscriberQuery.page = 1
  fetchSubscribers()
}

/**
 * 获取订阅者列表
 */
const fetchSubscribers = async () => {
  subscriberLoading.value = true
  try {
    const result = await getSubscribers(subscriberQuery)
    subscriberList.value = result.list
    subscriberTotal.value = result.total
  } catch {
    ElMessage.error('获取订阅者列表失败')
  } finally {
    subscriberLoading.value = false
  }
}

/**
 * 删除订阅者
 */
const handleDeleteSubscriber = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该订阅者吗？此操作不可恢复。', '提示', { type: 'warning' })
    await deleteSubscriber(id)
    ElMessage.success('删除成功')
    fetchSubscribers()
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message)
  }
}

onMounted(() => {
  fetchFriends()
  fetchArticles()
})
</script>

<style scoped>
/* 搜索表单样式已移至全局样式 main.scss */

.article-link {
  color: var(--el-color-primary);
  text-decoration: none;
  transition: color 0.2s;
}

.article-link:hover {
  text-decoration: underline;
}

.article-link.read {
  color: var(--el-text-color-regular);
}

.friend-link {
  color: var(--el-text-color-secondary);
  text-decoration: none;
}

.friend-link:hover {
  color: var(--el-color-primary);
  text-decoration: underline;
}
</style>
