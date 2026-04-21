<!--
项目名称：JeriBlog
文件名称：FriendList.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面组件 - FriendList页面
-->

<template>
  <common-list title="友链管理" :data="friendList" :loading="loading" :total="total" v-model:page="queryParams.page"
    v-model:page-size="queryParams.page_size" create-text="新增友链" @create="handleCreate" @refresh="fetchFriends"
    @update:page="fetchFriends" @update:pageSize="fetchFriends">
    <!-- 搜索表单 -->
    <template #toolbar-before>
      <div class="search-form friend-search">
        <el-input
          v-model="queryParams.keyword"
          placeholder="搜索名称、链接、描述..."
          clearable
          style="width: 240px"
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        />
        <el-select
          v-model="queryParams.type_id"
          placeholder="类型"
          clearable
          style="width: 150px"
          @change="handleSearch"
        >
          <el-option
            v-for="type in typeList"
            :key="type.id"
            :label="type.name"
            :value="type.id"
          />
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
    </template>

    <!-- 额外按钮 -->
    <template #toolbar-after>
      <el-button @click="handleTypeManage">
        类型管理
      </el-button>
    </template>

    <!-- 表格列 -->
    <el-table-column label="头像" width="80" align="center">
      <template #default="{ row }">
        <el-avatar v-if="row.avatar" :src="row.avatar" :size="40" />
        <el-avatar v-else :size="40">
          <el-icon>
            <Link />
          </el-icon>
        </el-avatar>
      </template>
    </el-table-column>

    <el-table-column label="友链名称" min-width="130" align="center">
      <template #default="{ row }">
        <span>{{ row.name }}</span>
        <el-tag v-if="row.is_invalid" type="warning" size="small" style="margin-left: 8px">失效</el-tag>
        <el-tag v-else-if="row.accessible > 0" type="danger" size="small" style="margin-left: 8px">异常({{ row.accessible }})</el-tag>
        <el-tag v-else-if="row.is_pending" type="info" size="small" style="margin-left: 8px">待审核</el-tag>
      </template>
    </el-table-column>

    <el-table-column label="链接地址" min-width="160" align="center">
      <template #default="{ row }">
        <span>{{ row.url }}</span>
      </template>
    </el-table-column>

    <el-table-column label="描述" min-width="180" align="center">
      <template #default="{ row }">
        <span v-if="row.description">{{ row.description }}</span>
        <span v-else style="color: #999">-</span>
      </template>
    </el-table-column>

    <el-table-column label="类型" width="120" align="center">
      <template #default="{ row }">
        <span v-if="row.type_name">{{ row.type_name }}</span>
        <span v-else style="color: #999">-</span>
      </template>
    </el-table-column>

    <el-table-column prop="sort" label="排序" width="100" align="center" />

    <el-table-column label="最新文章" width="180" align="center">
      <template #default="{ row }">
        <span v-if="row.rss_latime" :class="getRSSTimeClass(row.rss_latime)">{{ formatDateTime(row.rss_latime) }}</span>
        <span v-else style="color: #999">-</span>
      </template>
    </el-table-column>

    <el-table-column label="操作" width="180" align="center" fixed="right">
      <template #default="{ row }">
        <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
        <el-button type="danger" link size="small" @click="handleDelete(row.id)">删除</el-button>
      </template>
    </el-table-column>

    <!-- 额外内容 -->
    <template #extra>
      <friend-form-dialog v-model="dialogVisible" :edit-friend="currentFriend" @success="handleFriendSuccess" />
      <friend-type-manager ref="typeManagerRef" v-model="typeManagerVisible" @success="handleTypeSuccess" />
    </template>
  </common-list>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Link } from '@element-plus/icons-vue'
import CommonList from '@/components/common/CommonList.vue'
import type { Friend, FriendType } from '@/types/friend'
import type { PaginationQuery } from '@/types/request'
import { getFriends, deleteFriend, getFriendTypes } from '@/api/friend'
import { formatDateTime } from '@/utils/date'
import FriendFormDialog from './components/FriendFormDialog.vue'
import FriendTypeManager from './components/FriendTypeManager.vue'

const loading = ref(false)
const friendList = ref<Friend[]>([])
const typeList = ref<FriendType[]>([])
const total = ref(0)
const queryParams = ref<PaginationQuery & { keyword?: string; type_id?: number }>({
  page: 1,
  page_size: 20,
  keyword: '',
  type_id: undefined
})

// 对话框相关
const dialogVisible = ref(false)
const currentFriend = ref<Friend | null>(null)

// 类型管理对话框
const typeManagerVisible = ref(false)
const typeManagerRef = ref<InstanceType<typeof FriendTypeManager>>()

// 获取类型列表
const fetchTypes = async () => {
  try {
    const result = await getFriendTypes()
    typeList.value = result.list
  } catch {
    ElMessage.error('获取类型列表失败')
  }
}

const fetchFriends = async () => {
  loading.value = true
  try {
    const [result] = await Promise.all([
      getFriends(queryParams.value),
      new Promise(resolve => setTimeout(resolve, 300))
    ])
    friendList.value = result.list
    total.value = result.total
  } catch {
    ElMessage.error('获取友链列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  queryParams.value.page = 1
  fetchFriends()
}

// 重置
const handleReset = () => {
  queryParams.value = {
    page: 1,
    page_size: queryParams.value.page_size,
    keyword: '',
    type_id: undefined
  }
  fetchFriends()
}

const handleCreate = () => {
  currentFriend.value = null
  dialogVisible.value = true
}

const handleEdit = (friend: Friend) => {
  currentFriend.value = friend
  dialogVisible.value = true
}

const handleTypeManage = () => {
  typeManagerVisible.value = true
}

const handleTypeSuccess = () => {
  fetchTypes()
  fetchFriends()
}

const handleFriendSuccess = () => {
  fetchFriends()
  // 同时刷新类型管理器的数据（更新友链数量）
  typeManagerRef.value?.refreshData()
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个友链吗？', '提示', { type: 'warning' })
    await deleteFriend(id)
    ElMessage.success('删除成功')
    fetchFriends()
    // 同时刷新类型管理器的数据（更新友链数量）
    typeManagerRef.value?.refreshData()
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message)
  }
}

const getRSSTimeClass = (rssLatime?: string): string => {
  if (!rssLatime) return ''
  const months = (Date.now() - new Date(rssLatime).getTime()) / (1000 * 60 * 60 * 24 * 30)
  if (months > 6) return 'rss-danger'
  if (months > 3) return 'rss-warning'
  return ''
}

onMounted(() => {
  fetchTypes()
  fetchFriends()
})
</script>

<style scoped>
/* 搜索表单样式已移至全局样式 main.scss */

.rss-warning {
  color: #e6a23c;
}

.rss-danger {
  color: #f56c6c;
}
</style>
