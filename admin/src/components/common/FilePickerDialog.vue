<!--
项目名称：JeriBlog
文件名称：FilePickerDialog.vue
创建时间：2026-04-17 10:30:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：公共组件 - 文件选择对话框
-->

<template>
  <el-dialog v-model="dialogVisible" title="选择文件" :width="dialogWidth" :close-on-click-modal="false" @close="handleClose">
    <div class="file-picker-container">
      <!-- 搜索栏 -->
      <div class="filter-bar">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索文件名"
          clearable
          @input="handleSearch"
        >
          <template #prefix>
            <i class="ri-search-line"></i>
          </template>
        </el-input>
      </div>

      <!-- 文件列表 -->
      <div v-loading="loading" class="file-list">
        <div v-if="fileList.length === 0" class="empty-state">
          <el-empty description="暂无文件" />
        </div>
        <div v-else class="file-grid">
          <div
            v-for="file in fileList"
            :key="file.id"
            class="file-item"
            :class="{ selected: selectedFile?.id === file.id }"
            @click="handleSelectFile(file)"
          >
            <div class="file-preview">
              <img v-if="isImage(file.file_type)" :src="file.file_url" :alt="file.original_name" />
              <i v-else class="ri-file-line file-icon"></i>
            </div>
            <div class="file-info">
              <div class="file-name" :title="file.original_name">{{ file.original_name }}</div>
              <div class="file-meta">
                <span>{{ formatFileSize(file.file_size) }}</span>
              </div>
            </div>
            <div v-if="selectedFile?.id === file.id" class="selected-badge">
              <i class="ri-check-line"></i>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-bar">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[20, 50, 100, 200, 999999]"
          layout="total, sizes, prev, pager, next"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" :disabled="!selectedFile" @click="handleConfirm">确定</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getFileList } from '@/api/file'
import type { FileInfo } from '@/types/file'
import { useDebounceFn } from '@vueuse/core'

interface Props {
  modelValue: boolean
  fileType?: string // 限制文件类型，如 'image'
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm', file: FileInfo): void
}

const props = withDefaults(defineProps<Props>(), {
  fileType: ''
})

const emit = defineEmits<Emits>()

const dialogVisible = ref(false)
const loading = ref(false)
const fileList = ref<FileInfo[]>([])
const selectedFile = ref<FileInfo | null>(null)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const searchKeyword = ref('')

// 响应式对话框宽度
const dialogWidth = computed(() => {
  const width = window.innerWidth
  if (width <= 768) return '95%'
  if (width <= 1024) return '80%'
  if (width <= 1440) return '700px'
  return '800px'
})

// 监听 modelValue 变化
watch(
  () => props.modelValue,
  (val) => {
    dialogVisible.value = val
    if (val) {
      // 打开对话框时重置并加载数据
      selectedFile.value = null
      searchKeyword.value = ''
      currentPage.value = 1
      fetchFileList()
    }
  },
  { immediate: true }
)

// 监听 dialogVisible 变化
watch(dialogVisible, (val) => {
  emit('update:modelValue', val)
})

// 获取文件列表
const fetchFileList = async () => {
  try {
    loading.value = true
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value === 999999 ? 999999 : pageSize.value
    }

    // 如果有文件类型限制，添加到参数中
    if (props.fileType) {
      params.type = props.fileType
    }

    const response = await getFileList(params)

    // 如果有搜索关键词，在前端过滤
    if (searchKeyword.value.trim()) {
      const keyword = searchKeyword.value.trim().toLowerCase()
      fileList.value = response.list.filter(file =>
        file.original_name.toLowerCase().includes(keyword) ||
        file.filename.toLowerCase().includes(keyword)
      )
      total.value = fileList.value.length
    } else {
      fileList.value = response.list
      total.value = response.total
    }
  } catch (error) {
    ElMessage.error('获取文件列表失败')
  } finally {
    loading.value = false
  }
}

// 判断是否为图片
const isImage = (fileType: string) => {
  return fileType.startsWith('image/')
}

// 格式化文件大小
const formatFileSize = (size: number) => {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`
  return `${(size / (1024 * 1024)).toFixed(2)} MB`
}

// 选择文件
const handleSelectFile = (file: FileInfo) => {
  selectedFile.value = file
}

// 搜索（使用防抖）
const handleSearch = useDebounceFn(() => {
  currentPage.value = 1
  fetchFileList()
}, 500)

// 分页变化
const handlePageChange = () => {
  fetchFileList()
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchFileList()
}

// 确认选择
const handleConfirm = () => {
  if (selectedFile.value) {
    emit('confirm', selectedFile.value)
    handleClose()
  }
}

// 关闭对话框
const handleClose = () => {
  dialogVisible.value = false
  selectedFile.value = null
}
</script>

<style scoped lang="scss">
.file-picker-container {
  .filter-bar {
    margin-bottom: 16px;
  }

  .file-list {
    min-height: 300px;
    max-height: 450px;
    overflow-y: auto;
    margin-bottom: 16px;

    .empty-state {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 300px;
    }

    .file-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
      gap: 12px;

      @media (max-width: 768px) {
        grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
        gap: 10px;
      }

      .file-item {
        position: relative;
        border: 2px solid #e4e7ed;
        border-radius: 8px;
        padding: 8px;
        cursor: pointer;
        transition: all 0.3s;
        background: #fff;

        &:hover {
          border-color: #409eff;
          box-shadow: 0 2px 12px rgba(64, 158, 255, 0.2);
        }

        &.selected {
          border-color: #409eff;
          background: #ecf5ff;
        }

        .file-preview {
          width: 100%;
          height: 100px;
          display: flex;
          align-items: center;
          justify-content: center;
          background: #f5f7fa;
          border-radius: 4px;
          overflow: hidden;
          margin-bottom: 8px;

          @media (max-width: 768px) {
            height: 80px;
          }

          img {
            width: 100%;
            height: 100%;
            object-fit: cover;
          }

          .file-icon {
            font-size: 40px;
            color: #909399;

            @media (max-width: 768px) {
              font-size: 32px;
            }
          }
        }

        .file-info {
          .file-name {
            font-size: 12px;
            color: #303133;
            margin-bottom: 4px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .file-meta {
            font-size: 11px;
            color: #909399;
            display: flex;
            justify-content: space-between;
            gap: 4px;

            span {
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
            }
          }
        }

        .selected-badge {
          position: absolute;
          top: 8px;
          right: 8px;
          width: 22px;
          height: 22px;
          background: #409eff;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          color: #fff;
          font-size: 14px;
        }
      }
    }
  }

  .pagination-bar {
    display: flex;
    justify-content: center;

    :deep(.el-pagination) {
      @media (max-width: 768px) {
        .el-pagination__sizes {
          display: none;
        }
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
