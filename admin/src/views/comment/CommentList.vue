<!--
项目名称：JeriBlog
文件名称：CommentList.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面组件 - CommentList页面
-->

<template>
  <common-list title="评论管理" :data="commentList" :loading="loading" :total="total" :show-create="false"
    v-model:page="queryParams.page" v-model:page-size="queryParams.page_size" @refresh="fetchComments"
    @update:page="fetchComments" @update:pageSize="fetchComments">

    <!-- 搜索表单 -->
    <template #toolbar-before>
      <div class="search-form comment-search">
        <el-input
          v-model="queryParams.keyword"
          placeholder="搜索内容、用户、来源..."
          clearable
          style="width: 240px"
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        />
        <el-select
          v-model="queryParams.status"
          placeholder="状态"
          clearable
          style="width: 120px"
          @change="handleSearch"
        >
          <el-option label="显示" :value="1" />
          <el-option label="隐藏" :value="0" />
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
    </template>

    <!-- 表格列 -->
    <el-table-column label="用户信息" width="180" align="center">
      <template #default="{ row }">
        <div style="display: flex; align-items: center; gap: 8px">
          <el-avatar :size="40" :src="row.user.avatar" style="flex-shrink: 0">
            <el-icon>
              <User />
            </el-icon>
          </el-avatar>
          <div style="flex: 1; min-width: 0; overflow: hidden; text-align: left">
            <div style="font-weight: 500; overflow: hidden; text-overflow: ellipsis; white-space: nowrap">{{
              row.user.nickname }}</div>
            <el-tooltip :content="row.user.email" placement="top">
              <div style="font-size: 12px; color: #999; overflow: hidden; text-overflow: ellipsis; white-space: nowrap">{{
                row.user.email }}</div>
            </el-tooltip>
          </div>
        </div>
      </template>
    </el-table-column>

    <el-table-column label="评论内容" min-width="300">
      <template #default="{ row }">
        <div style="line-height: 1.6; display: flex; align-items: center; gap: 8px">
          <span>{{ row.content }}</span>
          <el-tag v-if="row.deleted_at" type="danger" size="small">已删除</el-tag>
          <el-tag v-if="row.parent_id" type="info" size="small">子评论</el-tag>
        </div>
      </template>
    </el-table-column>

    <el-table-column label="评论来源" width="220" align="center">
      <template #default="{ row }">
        <div style="display: flex; align-items: center; gap: 8px">
          <el-tag v-if="row.target.type !== 'article'" type="success" size="small">
            {{ getTargetTypeText(row.target.type) }}
          </el-tag>
          <el-tooltip :content="row.target.title" placement="top">
            <div style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 12px; flex: 1">
              {{ row.target.title }}
            </div>
          </el-tooltip>
        </div>
      </template>
    </el-table-column>

    <el-table-column label="评论时间" width="180" align="center">
      <template #default="{ row }">
        {{ formatDateTime(row.created_at) }}
      </template>
    </el-table-column>

    <el-table-column label="状态" width="100" align="center">
      <template #default="{ row }">
        <el-switch v-model="row.status" :active-value="1" :inactive-value="0" inline-prompt active-text="显示"
          inactive-text="隐藏" @change="handleStatusChange(row)" />
      </template>
    </el-table-column>

    <el-table-column label="操作" width="220" align="center" fixed="right">
      <template #default="{ row }">
        <el-button v-if="!row.deleted_at" type="primary" link size="small" @click="openReplyDialog(row)">
          回复
        </el-button>
        <el-button v-if="row.deleted_at" type="success" link size="small" @click="handleRestore(row.id)">
          恢复
        </el-button>
        <el-button v-else type="danger" link size="small" @click="handleDelete(row.id)">
          删除
        </el-button>
      </template>
    </el-table-column>
  </common-list>

  <!-- 回复对话框 -->
  <el-dialog v-model="replyDialogVisible" title="回复评论" width="500px" destroy-on-close>
    <div v-if="replyingComment" class="reply-info">
      <div class="info-row">
        <span class="label">评论来源：</span>
        <span class="value">
          <el-tag v-if="replyingComment.target.type !== 'article'" type="success" size="small">
            {{ getTargetTypeText(replyingComment.target.type) }}
          </el-tag>
          {{ replyingComment.target.title }}
        </span>
      </div>
      <div class="info-row">
        <span class="label">评论时间：</span>
        <span class="value">{{ formatDateTime(replyingComment.created_at) }}</span>
      </div>
      <el-divider style="margin: 12px 0" />
      <div class="reply-to">回复 <span class="nickname">{{ replyingComment.user.nickname }}</span>：</div>
      <div class="original-content">{{ replyingComment.content }}</div>
    </div>
    <el-form :model="replyForm" label-width="80px" style="margin-top: 16px">
      <el-form-item label="回复内容">
        <el-input v-model="replyForm.content" type="textarea" :rows="4" placeholder="请输入回复内容..." show-word-limit
          maxlength="500" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="replyDialogVisible = false">取消</el-button>
      <el-button type="primary" :loading="replying" @click="handleReply">提交回复</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import CommonList from '@/components/common/CommonList.vue'
import type { Comment, CommentQuery } from '@/types/comment'
import { getComments, deleteComment, restoreComment, toggleCommentStatus, createComment } from '@/api/comment'
import { formatDateTime } from '@/utils/date'

const loading = ref(false)
const commentList = ref<Comment[]>([])
const total = ref(0)
const queryParams = ref<CommentQuery>({
  page: 1,
  page_size: 20,
  keyword: '',
  status: undefined
})

// 回复相关状态
const replyDialogVisible = ref(false)
const replying = ref(false)
const replyingComment = ref<Comment | null>(null)
const replyForm = ref({
  content: ''
})

const fetchComments = async () => {
  loading.value = true
  try {
    const [result] = await Promise.all([
      getComments(queryParams.value),
      new Promise(resolve => setTimeout(resolve, 300))
    ])
    commentList.value = result.list
    total.value = result.total
  } catch {
    ElMessage.error('获取评论列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  queryParams.value.page = 1
  fetchComments()
}

// 重置
const handleReset = () => {
  queryParams.value = {
    page: 1,
    page_size: queryParams.value.page_size,
    keyword: '',
    status: undefined
  }
  fetchComments()
}

const handleStatusChange = async (comment: Comment) => {
  const statusText = comment.status === 1 ? '显示' : '隐藏'
  try {
    await toggleCommentStatus(comment.id)
    ElMessage.success(`已设置为${statusText}`)
  } catch (error) {
    comment.status = comment.status === 1 ? 0 : 1
    if (error instanceof Error) {
      ElMessage.error(error.message)
    } else {
      ElMessage.error('状态切换失败')
    }
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', { type: 'warning' })
    await deleteComment(id)
    ElMessage.success('删除成功')
    fetchComments()
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message)
  }
}

const handleRestore = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要恢复这条评论吗？', '提示', { type: 'info' })
    await restoreComment(id)
    ElMessage.success('恢复成功')
    fetchComments()
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message)
  }
}

const openReplyDialog = (comment: Comment) => {
  replyingComment.value = comment
  replyForm.value.content = ''
  replyDialogVisible.value = true
}

const handleReply = async () => {
  if (!replyForm.value.content.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  if (!replyingComment.value) {
    ElMessage.error('评论信息错误')
    return
  }

  replying.value = true
  try {
    await createComment({
      content: replyForm.value.content,
      target_type: replyingComment.value.target.type,
      target_key: replyingComment.value.target.key,
      parent_id: replyingComment.value.id
    })
    ElMessage.success('回复成功')
    replyDialogVisible.value = false
    fetchComments()
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error(error.message)
    } else {
      ElMessage.error('回复失败')
    }
  } finally {
    replying.value = false
  }
}

// 获取目标类型显示文本
const getTargetTypeText = (type: string) => {
  const typeMap: Record<string, string> = {
    article: '文章',
    moment: '动态',
    guestbook: '留言',
    page: '页面'
  }
  return typeMap[type] || type
}

onMounted(fetchComments)
</script>

<style scoped lang="scss">
/* 搜索表单样式已移至全局样式 main.scss */

.reply-info {
  padding: 12px;
  background-color: #f5f7fa;
  border-radius: 4px;
  border: 1px solid #dcdfe6;

  .info-row {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    font-size: 13px;

    &:last-child {
      margin-bottom: 0;
    }

    .label {
      color: #909399;
      flex-shrink: 0;
    }

    .value {
      color: #606266;
      display: flex;
      align-items: center;
      gap: 6px;
    }
  }

  .reply-to {
    font-size: 14px;
    color: #606266;
    margin-bottom: 8px;

    .nickname {
      color: #409eff;
      font-weight: 500;
    }
  }

  .original-content {
    font-size: 13px;
    color: #909399;
    line-height: 1.6;
    white-space: pre-wrap;
    word-wrap: break-word;
  }
}
</style>
