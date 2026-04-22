<!--
项目名称：JeriBlog
文件名称：FileList.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面组件 - FileList页面
-->

<template>
  <common-list title="文件管理" :data="fileList" :loading="loading" :total="total" :show-create="false"
    v-model:page="query.page" v-model:page-size="query.page_size" @refresh="loadList" @update:page="loadList"
    @update:pageSize="loadList" row-key="id">
    <!-- 搜索表单 -->
    <template #toolbar-before>
      <div class="search-form file-search">
        <el-input
          v-model="query.keyword"
          placeholder="搜索文件名、原始文件名..."
          clearable
          style="width: 240px"
          @keyup.enter="handleSearch"
          @clear="handleSearch"
        />
        <el-select
          v-model="query.status"
          placeholder="状态"
          clearable
          style="width: 140px"
          @change="handleSearch"
        >
          <el-option label="未使用" :value="0" />
          <el-option label="使用中" :value="1" />
        </el-select>
        <el-select
          v-model="query.upload_type"
          placeholder="用途"
          clearable
          style="width: 140px"
          @change="handleSearch"
        >
          <el-option label="用户头像" value="用户头像" />
          <el-option label="文章封面" value="文章封面" />
          <el-option label="站长头像" value="站长头像" />
          <el-option label="站长形象" value="站长形象" />
          <el-option label="博客图标" value="博客图标" />
          <el-option label="博客背景" value="博客背景" />
          <el-option label="博客截图" value="博客截图" />
          <el-option label="展览图片" value="展览图片" />
          <el-option label="友情链接A" value="友情链接A" />
          <el-option label="友情链接S" value="友情链接S" />
          <el-option label="微信收款码" value="微信收款码" />
          <el-option label="支付宝收款码" value="支付宝收款码" />
          <el-option label="动态配图" value="动态配图" />
          <el-option label="动态视频" value="动态视频" />
          <el-option label="反馈投诉" value="反馈投诉" />
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
    </template>

    <!-- 右上角工具栏 -->
    <template #toolbar-after>
      <el-button type="primary" @click="uploadDialogVisible = true">上传配置</el-button>
    </template>

    <!-- 表格列 -->
    <el-table-column label="预览" width="80" align="center">
      <template #default="{ row }">
        <el-image v-if="isImage(row)" :src="row.file_url" fit="cover"
          style="width: 50px; height: 50px; border-radius: 4px" />
      </template>
    </el-table-column>

    <el-table-column label="文件名" min-width="180" align="center">
      <template #default="{ row }">
        <div style="display: flex; flex-direction: column; align-items: center; gap: 4px;">
          <span style="font-weight: 500">{{ row.file_name }}</span>
          <span style="font-size: 12px; color: #909399">{{ formatFileSize(row.file_size) }}</span>
        </div>
      </template>
    </el-table-column>

    <el-table-column prop="original_name" label="原始文件名" min-width="200" align="center" show-overflow-tooltip />

    <el-table-column prop="file_type" label="类型" width="100" align="center" />

    <el-table-column label="状态" width="100" align="center">
      <template #default="{ row }">
        <el-tag :type="getStatusTagType(row.status)" size="small" effect="light">
          {{ getStatusText(row.status) }}
        </el-tag>
      </template>
    </el-table-column>

    <el-table-column label="引用数" width="100" align="center">
      <template #default="{ row }">
        <el-link
          type="primary"
          :underline="false"
          @click="handleShowReferences(row)"
          :disabled="!row.reference_count || row.reference_count === 0"
        >
          {{ row.reference_count || 0 }}
        </el-link>
      </template>
    </el-table-column>

    <el-table-column prop="upload_type" label="用途" width="100" align="center" />

    <el-table-column label="上传时间" width="180" align="center">
      <template #default="{ row }">
        {{ formatDateTime(row.upload_time) }}
      </template>
    </el-table-column>

    <el-table-column label="操作" width="180" align="center" fixed="right">
      <template #default="{ row }">
        <el-button link type="primary" size="small" @click="copyUrl(row)">复制链接</el-button>
        <el-button link type="danger" size="small" @click="handleDelete(row.id)">删除</el-button>
      </template>
    </el-table-column>
    <!-- 额外挂载区域 -->
    <template #extra>
      <upload-config-dialog v-model="uploadDialogVisible" />

      <!-- 文件引用详情对话框 -->
      <el-dialog
        v-model="referencesDialogVisible"
        title="文件引用详情"
        width="600px"
        :close-on-click-modal="false"
      >
        <div v-loading="referencesLoading" class="references-content">
          <el-empty v-if="!referencesLoading && references.length === 0" description="暂无引用" />

          <div v-else class="reference-list">
            <div v-for="(ref, index) in references" :key="index" class="reference-item">
              <div class="reference-header">
                <el-tag :type="getReferenceTypeTag(ref.type)" size="small">
                  {{ getReferenceTypeName(ref.type) }}
                </el-tag>
                <span class="reference-field">{{ ref.field }}</span>
              </div>

              <div class="reference-body">
                <div class="reference-title">{{ ref.title }}</div>
                <el-link
                  v-if="ref.url"
                  :href="ref.url"
                  target="_blank"
                  type="primary"
                  :underline="false"
                  class="reference-link"
                >
                  <i class="ri-external-link-line"></i>
                  查看详情
                </el-link>
              </div>
            </div>
          </div>
        </div>

        <template #footer>
          <el-button @click="referencesDialogVisible = false">关闭</el-button>
        </template>
      </el-dialog>
    </template>
  </common-list>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import CommonList from '@/components/common/CommonList.vue'
import { getFileList, deleteFile, getFileReferences, type FileReference } from '@/api/file'
import type { FileInfo, FileQuery } from '@/types/file'
import { formatDateTime } from '@/utils/date'
import UploadConfigDialog from '@/views/file/components/UploadConfigDialog.vue'

const query = reactive<FileQuery>({
  page: 1,
  page_size: 20,
  keyword: undefined,
  status: undefined,
  upload_type: undefined
})
const fileList = ref<FileInfo[]>([])
const total = ref(0)
const loading = ref(false)
const uploadDialogVisible = ref(false)

// 引用详情相关
const referencesDialogVisible = ref(false)
const referencesLoading = ref(false)
const references = ref<FileReference[]>([])
const currentFile = ref<FileInfo | null>(null)

const loadList = async () => {
  loading.value = true
  try {
    const [data] = await Promise.all([
      getFileList(query),
      new Promise(resolve => setTimeout(resolve, 300))
    ])
    fileList.value = data.list
    total.value = data.total
  } catch {
    ElMessage.error('加载失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  query.page = 1
  loadList()
}

const handleReset = () => {
  query.keyword = undefined
  query.status = undefined
  query.upload_type = undefined
  query.page = 1
  loadList()
}

const copyUrl = async (file: FileInfo) => {
  try {
    await navigator.clipboard.writeText(file.file_url)
    ElMessage.success('已复制')
  } catch {
    ElMessage.error('复制失败')
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个文件吗？', '提示', { type: 'warning' })
    await deleteFile(id)
    ElMessage.success('删除成功')
    loadList()
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message)
  }
}

// 显示文件引用详情
const handleShowReferences = async (file: FileInfo) => {
  if (!file.reference_count || file.reference_count === 0) {
    return
  }

  currentFile.value = file
  referencesDialogVisible.value = true
  referencesLoading.value = true
  references.value = []

  try {
    references.value = await getFileReferences(file.id)
  } catch (error: any) {
    console.error('加载引用详情失败:', error)
    ElMessage.error('加载引用详情失败')
  } finally {
    referencesLoading.value = false
  }
}

// 获取引用类型标签颜色
const getReferenceTypeTag = (type: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const typeMap: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    article: 'primary',
    user: 'success',
    friend: 'warning',
    setting: 'info'
  }
  return typeMap[type] || 'info'
}

// 获取引用类型名称
const getReferenceTypeName = (type: string) => {
  const nameMap: Record<string, string> = {
    article: '文章',
    user: '用户',
    friend: '友链',
    setting: '系统设置'
  }
  return nameMap[type] || type
}

const isImage = (file: FileInfo) => file.file_type?.startsWith('image/')

const formatFileSize = (size: number) => {
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  return (size / (1024 * 1024)).toFixed(1) + ' MB'
}

const getStatusTagType = (status: number) => {
  return status === 1 ? 'success' : 'info'
}

const getStatusText = (status: number) => {
  return status === 1 ? '使用中' : '未使用'
}

onMounted(loadList)
</script>

<style scoped lang="scss">
/* 搜索表单样式已移至全局样式 main.scss */

.references-content {
  min-height: 150px;
  max-height: 500px;
  overflow-y: auto;
  padding: 4px;
}

.reference-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.reference-item {
  padding: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #fafafa;
  transition: all 0.3s;

  &:hover {
    border-color: #409eff;
    box-shadow: 0 2px 12px rgba(64, 158, 255, 0.15);
    transform: translateY(-2px);
  }
}

.reference-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;

  .reference-field {
    font-size: 13px;
    color: #909399;
    padding: 2px 8px;
    background: #fff;
    border-radius: 4px;
    border: 1px solid #e4e7ed;
  }
}

.reference-body {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;

  .reference-title {
    flex: 1;
    font-size: 14px;
    color: #303133;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .reference-link {
    flex-shrink: 0;
    font-size: 13px;

    i {
      margin-right: 4px;
    }
  }
}
</style>
