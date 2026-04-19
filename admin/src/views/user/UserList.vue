<!--
项目名称：JeriBlog
文件名称：UserList.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面组件 - UserList页面
-->

<template>
  <common-list title="用户列表" :data="userList" :loading="loading" :total="total" v-model:page="queryParams.page"
    v-model:page-size="queryParams.page_size" create-text="新增用户" @create="handleCreate" @refresh="fetchUsers"
    @update:page="fetchUsers" @update:pageSize="fetchUsers">

    <!-- 搜索表单 -->
    <template #toolbar-before>
      <div class="search-form">
        <el-input
          v-model="queryParams.keyword"
          placeholder="搜索昵称、邮箱、网站..."
          clearable
          style="width: 240px"
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        />
        <el-select
          v-model="queryParams.role"
          placeholder="角色"
          clearable
          style="width: 140px"
          @change="handleSearch"
        >
          <el-option label="超级管理员" value="super_admin" />
          <el-option label="管理员" value="admin" />
          <el-option label="普通用户" value="user" />
          <el-option label="访客" value="guest" />
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
    </template>

    <!-- 表格列 -->
    <el-table-column label="头像" width="100" align="center">
      <template #default="{ row }">
        <el-avatar v-if="row.avatar" :src="row.avatar" :size="40" />
        <el-avatar v-else :size="40">
          <el-icon>
            <User />
          </el-icon>
        </el-avatar>
      </template>
    </el-table-column>

    <el-table-column label="昵称" min-width="130" align="center">
      <template #default="{ row }">
        <div style="display: flex; align-items: center; justify-content: center; gap: 8px">
          <span>{{ row.nickname }}</span>
          <el-tag v-if="row.badge" type="info" effect="plain" size="small">{{ row.badge }}</el-tag>
          <el-tag v-if="row.deleted_at" type="danger" size="small">已删除</el-tag>
        </div>
      </template>
    </el-table-column>

    <el-table-column label="邮箱" min-width="150" align="center">
      <template #default="{ row }">
        <div style="display: flex; align-items: center; justify-content: center; gap: 8px">
          <span v-if="row.email">{{ row.email }}</span>
          <span v-else style="color: #999">-</span>
          <el-tag v-if="!row.is_enabled" type="danger" size="small">禁用</el-tag>
        </div>
      </template>
    </el-table-column>

    <el-table-column label="网站地址" min-width="150" align="center">
      <template #default="{ row }">
        <span v-if="row.website">{{ row.website }}</span>
        <span v-else style="color: #999">-</span>
      </template>
    </el-table-column>

    <el-table-column label="角色" width="120" align="center">
      <template #default="{ row }">
        <el-tag v-if="row.role === 'super_admin'" type="danger" size="small">超级管理员</el-tag>
        <el-tag v-else-if="row.role === 'admin'" type="warning" size="small">管理员</el-tag>
        <el-tag v-else-if="row.role === 'user'" type="success" size="small">普通用户</el-tag>
        <el-tag v-else type="info" size="small">访客</el-tag>
      </template>
    </el-table-column>

    <el-table-column label="登录方式" width="150" align="center">
      <template #default="{ row }">
        <div class="login-methods">
          <template v-if="row.has_password">
            <i class="ri-lock-password-fill"></i>
          </template>

          <template v-if="row.github_id">
            <i class="ri-github-fill"></i>
          </template>

          <template v-if="row.google_id">
            <i class="ri-google-fill"></i>
          </template>

          <template v-if="row.qq_id">
            <i class="ri-qq-fill"></i>
          </template>

          <template v-if="row.microsoft_id">
            <i class="ri-microsoft-fill"></i>
          </template>
        </div>
      </template>
    </el-table-column>

    <el-table-column label="最后登录" width="180" align="center">
      <template #default="{ row }">
        {{ formatDateTime(row.last_login) }}
      </template>
    </el-table-column>

    <el-table-column label="操作" width="180" align="center" fixed="right">
      <template #default="{ row }">
        <template v-if="!row.deleted_at">
          <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
          <el-tooltip
            v-if="row.id === 1 && row.role === 'super_admin'"
            content="禁止删除默认超级管理员"
            placement="top"
          >
            <el-button type="danger" link size="small" disabled>删除</el-button>
          </el-tooltip>
          <el-button
            v-else
            type="danger"
            link
            size="small"
            @click="handleDelete(row.id)"
          >删除</el-button>
        </template>
      </template>
    </el-table-column>

    <!-- 额外内容 -->
    <template #extra>
      <!-- 用户表单对话框 -->
      <user-form-dialog v-model="dialogVisible" :edit-user="currentUser" @success="fetchUsers" />
    </template>
  </common-list>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import CommonList from '@/components/common/CommonList.vue'
import type { User as UserType, UserQuery } from '@/types/user'
import { getUsers, deleteUser } from '@/api/user'
import UserFormDialog from './components/UserFormDialog.vue'
import { formatDateTime } from '@/utils/date'

const loading = ref(false)
const userList = ref<UserType[]>([])
const total = ref(0)
const queryParams = ref<UserQuery>({
  page: 1,
  page_size: 20,
  keyword: undefined,
  role: undefined
})

// 对话框相关
const dialogVisible = ref(false)
const currentUser = ref<UserType | null>(null)

const fetchUsers = async () => {
  loading.value = true
  try {
    const [result] = await Promise.all([
      getUsers(queryParams.value),
      new Promise(resolve => setTimeout(resolve, 300))
    ])
    userList.value = result.list
    total.value = result.total
  } catch {
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  queryParams.value.page = 1
  fetchUsers()
}

const handleReset = () => {
  queryParams.value.keyword = undefined
  queryParams.value.role = undefined
  queryParams.value.page = 1
  fetchUsers()
}

const handleCreate = () => {
  currentUser.value = null
  dialogVisible.value = true
}

const handleEdit = (user: UserType) => {
  currentUser.value = user
  dialogVisible.value = true
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个用户吗？', '提示', { type: 'warning' })
    await deleteUser(id)
    ElMessage.success('删除成功')
    currentUser.value = null
    dialogVisible.value = false
    fetchUsers()
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message)
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.search-form {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.login-methods {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.login-methods i {
  font-size: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
