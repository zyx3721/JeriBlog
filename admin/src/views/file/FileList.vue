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
  <div class="file-list-page">
    <transition name="filter-slide">
      <file-filter
        v-if="showFilter"
        v-model="query"
        @close="showFilter = false"
        @search="loadList"
      />
    </transition>

    <common-list
      title="文件管理"
      :data="fileList"
      :loading="loading"
      :total="total"
      :show-create="false"
      :filter-active="showFilter"
      :filter-count="activeFilterCount"
      v-model:page="query.page"
      v-model:page-size="query.page_size"
      @refresh="loadList"
      @filter="toggleFilter"
      @update:page="loadList"
      @update:pageSize="loadList"
    >
      <template #toolbar-before>
        <template v-if="!showFilter">
          <el-select
            v-model="quickFilters.upload_type"
            placeholder="用途"
            clearable
            style="width: 140px"
            @change="handleQuickFilterChange"
          >
            <el-option-group label="用户相关">
              <el-option label="用户头像" value="用户头像" />
            </el-option-group>
            <el-option-group label="文章相关">
              <el-option label="文章封面" value="文章封面" />
              <el-option label="文章配图" value="文章配图" />
              <el-option label="文章视频" value="文章视频" />
              <el-option label="文章音频" value="文章音频" />
              <el-option label="文章附件" value="文章附件" />
            </el-option-group>
            <el-option-group label="动态相关">
              <el-option label="动态配图" value="动态配图" />
              <el-option label="动态视频" value="动态视频" />
              <el-option label="动态音频" value="动态音频" />
            </el-option-group>
            <el-option-group label="评论相关">
              <el-option label="评论贴图" value="评论贴图" />
            </el-option-group>
            <el-option-group label="友链相关">
              <el-option label="友情链接A" value="友情链接A" />
              <el-option label="友情链接S" value="友情链接S" />
            </el-option-group>
            <el-option-group label="系统设置">
              <el-option label="站长头像" value="站长头像" />
              <el-option label="站长形象" value="站长形象" />
              <el-option label="博客图标" value="博客图标" />
              <el-option label="博客背景" value="博客背景" />
              <el-option label="博客截图" value="博客截图" />
              <el-option label="展览图片" value="展览图片" />
              <el-option label="微信收款码" value="微信收款码" />
              <el-option label="支付宝收款码" value="支付宝收款码" />
            </el-option-group>
            <el-option-group label="其他">
              <el-option label="菜单图标" value="菜单图标" />
              <el-option label="反馈投诉" value="反馈投诉" />
            </el-option-group>
          </el-select>
          <el-select
            v-model="quickFilters.file_type"
            placeholder="全部类型"
            clearable
            class="quick-filter-769"
            style="width: 120px"
            @change="handleQuickFilterChange"
          >
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
            <el-option label="音频" value="audio" />
            <el-option label="文档" value="application" />
          </el-select>
          <el-select
            v-model="quickFilters.status"
            placeholder="使用状态"
            clearable
            class="quick-filter-769"
            style="width: 100px"
            @change="handleQuickFilterChange"
          >
            <el-option label="使用中" :value="1" />
            <el-option label="未使用" :value="0" />
          </el-select>
        </template>
      </template>

      <el-table-column label="预览" width="80" align="center">
        <template #default="{ row }">
          <el-image
            v-if="isImage(row)"
            :src="row.file_url"
            fit="cover"
            style="width: 50px; height: 50px; border-radius: 4px"
          />
        </template>
      </el-table-column>

      <el-table-column label="文件名" min-width="180" align="center">
        <template #default="{ row }">
          <div style="display: flex; flex-direction: column; align-items: center; gap: 4px">
            <span style="margin-right: 8px; font-weight: 500">{{ row.file_name }}</span>
            <span style="font-size: 12px; color: #909399">{{ formatFileSize(row.file_size) }}</span>
          </div>
        </template>
      </el-table-column>

      <el-table-column
        prop="original_name"
        label="原始文件名"
        align="center"
        min-width="200"
        show-overflow-tooltip
      />

      <el-table-column prop="file_type" label="类型" width="120" align="center" />

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
            @click="handleShowReferences(row)"
            :disabled="!row.reference_count || row.reference_count === 0"
          >
            {{ row.reference_count || 0 }}
          </el-link>
        </template>
      </el-table-column>

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
        <!-- 文件引用详情对话框 -->
        <el-dialog
          v-model="referencesDialogVisible"
          title="文件引用详情"
          width="90%"
          style="max-width: 600px"
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
                  <!-- 评论类型时显示所属内容类型 -->
                  <el-tag
                    v-if="ref.type === 'comment' && ref.target_type"
                    :type="getCommentTargetTypeTag(ref.target_type, ref.target_key)"
                    size="small"
                  >
                    {{ getCommentTargetTypeName(ref.target_type, ref.target_key) }}
                  </el-tag>
                </div>

                <div class="reference-body">
                  <div class="reference-title">{{ getReferenceTitle(ref) }}</div>
                  <el-link type="primary" class="reference-link" @click="handleReferenceClick(ref)">
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useRouter } from 'vue-router';
import CommonList from '@/components/common/CommonList.vue';
import FileFilter from './components/FileFilter.vue';
import { getFileList, deleteFile, getFileReferences, type FileReference } from '@/api/file';
import type { FileInfo, FileListQuery } from '@/types/file';
import { formatDateTime } from '@/utils/date';

const query = ref<FileListQuery>({ page: 1, page_size: 20 });
const fileList = ref<FileInfo[]>([]);
const router = useRouter();
const total = ref(0);
const loading = ref(false);
const showFilter = ref(false);

// 引用详情相关
const referencesDialogVisible = ref(false);
const referencesLoading = ref(false);
const references = ref<FileReference[]>([]);
const currentFile = ref<FileInfo | null>(null);

const quickFilters = reactive({
  file_type: undefined as string | undefined,
  status: undefined as number | undefined,
  upload_type: undefined as string | undefined,
});

/**
 * 计算当前激活的筛选项数量
 */
const activeFilterCount = computed(() => {
  let count = 0;
  if (query.value.keyword) count++;
  if (query.value.file_type !== undefined) count++;
  if (query.value.status !== undefined) count++;
  if (query.value.upload_type) count++;
  if (query.value.min_size || query.value.max_size) count++;
  if (query.value.start_time || query.value.end_time) count++;
  return count;
});

/**
 * 切换筛选面板显示状态
 */
const toggleFilter = () => {
  showFilter.value = !showFilter.value;
  if (!showFilter.value) {
    syncQuickFiltersFromQuery();
  }
};

/**
 * 从 query 同步筛选条件到快速筛选
 */
const syncQuickFiltersFromQuery = () => {
  quickFilters.file_type = query.value.file_type;
  quickFilters.status = query.value.status;
  quickFilters.upload_type = query.value.upload_type;
};

/**
 * 处理快速筛选变化
 */
const handleQuickFilterChange = () => {
  query.value.file_type = quickFilters.file_type;
  query.value.status = quickFilters.status;
  query.value.upload_type = quickFilters.upload_type;
  query.value.page = 1;
  loadList();
};

const loadList = async () => {
  loading.value = true;
  try {
    const [data] = await Promise.all([
      getFileList(query.value),
      new Promise(resolve => setTimeout(resolve, 300)),
    ]);
    fileList.value = data.list;
    total.value = data.total;
  } catch {
    ElMessage.error('加载失败');
  } finally {
    loading.value = false;
  }
};

const copyUrl = async (file: FileInfo) => {
  try {
    // 优先使用现代 Clipboard API
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(file.file_url);
      ElMessage.success('已复制');
      return;
    }

    // 降级方案：使用传统的 document.execCommand
    const textArea = document.createElement('textarea');
    textArea.value = file.file_url;
    textArea.style.position = 'fixed';
    textArea.style.left = '-999999px';
    textArea.style.top = '-999999px';
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    const successful = document.execCommand('copy');
    document.body.removeChild(textArea);

    if (successful) {
      ElMessage.success('已复制');
    } else {
      ElMessage.error('复制失败');
    }
  } catch (error) {
    ElMessage.error('复制失败');
  }
};

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个文件吗？', '提示', {
      type: 'warning',
    });
    await deleteFile(id);
    ElMessage.success('删除成功');
    loadList();
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message);
  }
};

// 显示文件引用详情
const handleShowReferences = async (file: FileInfo) => {
  if (!file.reference_count || file.reference_count === 0) {
    return;
  }

  currentFile.value = file;
  referencesDialogVisible.value = true;
  referencesLoading.value = true;
  references.value = [];

  try {
    references.value = await getFileReferences(file.id);
  } catch (error: unknown) {
    ElMessage.error((error as Error).message);
  } finally {
    referencesLoading.value = false;
  }
};

// 处理引用详情点击跳转
const handleReferenceClick = (ref: FileReference) => {
  // 根据引用类型进行不同的跳转处理
  switch (ref.type) {
    case 'user':
      // 跳转到用户管理页面，通过 state 传递搜索关键词
      referencesDialogVisible.value = false; // 当前窗口跳转，关闭对话框
      router.push({
        path: '/users',
        state: { keyword: ref.title },
      });
      break;

    case 'article':
      // 跳转到文章管理页面，通过 state 传递搜索关键词
      referencesDialogVisible.value = false; // 当前窗口跳转，关闭对话框
      router.push({
        path: '/articles',
        state: { keyword: ref.title },
      });
      break;

    case 'moment':
      // 跳转到动态管理页面，通过 state 传递搜索关键词
      referencesDialogVisible.value = false; // 当前窗口跳转，关闭对话框
      router.push({
        path: '/moments',
        state: { keyword: ref.title },
      });
      break;

    case 'comment':
      // 跳转到评论所属文章/页面的前端访问链接
      if (ref.url) {
        window.open(ref.url, '_blank'); // 新窗口打开，不关闭对话框
      } else {
        ElMessage.warning('该评论所属内容已被删除');
      }
      break;

    case 'friend':
      // 跳转到友链管理页面，通过 state 传递搜索关键词
      referencesDialogVisible.value = false; // 当前窗口跳转，关闭对话框
      router.push({
        path: '/friends',
        state: { keyword: ref.title },
      });
      break;

    case 'setting':
      // 跳转到系统设置页面
      referencesDialogVisible.value = false; // 当前窗口跳转，关闭对话框
      router.push('/settings');
      break;

    case 'menu':
      // 跳转到菜单管理页面
      referencesDialogVisible.value = false; // 当前窗口跳转，关闭对话框
      router.push('/menus');
      break;

    case 'feedback':
      // 跳转到反馈投诉页面，通过 state 传递搜索关键词
      referencesDialogVisible.value = false; // 当前窗口跳转，关闭对话框
      router.push({
        path: '/feedback',
        state: { keyword: ref.title },
      });
      break;

    default:
      ElMessage.warning('未知的引用类型');
  }
};

// 获取引用类型标签颜色
const getReferenceTypeTag = (
  type: string
): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const typeMap: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    article: 'primary',
    user: 'success',
    friend: 'warning',
    setting: 'info',
    moment: 'primary',
    comment: 'success',
    menu: 'info',
    feedback: 'danger',
  };
  return typeMap[type] || 'info';
};

// 获取引用类型名称
const getReferenceTypeName = (type: string) => {
  const nameMap: Record<string, string> = {
    article: '文章',
    user: '用户',
    friend: '友链',
    setting: '系统设置',
    moment: '动态',
    comment: '评论',
    menu: '菜单',
    feedback: '反馈',
  };
  return nameMap[type] || type;
};

// 获取评论所属内容类型的标签颜色
const getCommentTargetTypeTag = (targetType: string, targetKey?: string) => {
  if (targetType === 'article') return 'success';
  if (targetType === 'page') {
    if (targetKey === 'moment') return 'warning';
    if (targetKey === 'message') return 'info';
    if (targetKey === 'friend') return 'danger';
  }
  return undefined;
};

// 获取评论所属内容类型的名称
const getCommentTargetTypeName = (targetType: string, targetKey?: string) => {
  if (targetType === 'article') return '文章';
  if (targetType === 'page') {
    if (targetKey === 'moment') return '动态';
    if (targetKey === 'message') return '留言';
    if (targetKey === 'friend') return '友链';
    return '页面';
  }
  return targetType;
};

// 获取引用标题显示内容
const getReferenceTitle = (ref: FileReference) => {
  // 如果是系统设置类型，根据字段名区分基本配置和博客配置
  if (ref.type === 'setting' && ref.field) {
    // 基本配置字段
    const basicFields = ['站长头像', '站长形象'];
    // 博客配置字段
    const blogFields = [
      '博客图标',
      '博客背景',
      '博客截图',
      '展览图片',
      '微信收款码',
      '支付宝收款码',
    ];

    if (basicFields.includes(ref.field)) {
      return '基本配置';
    } else if (blogFields.includes(ref.field)) {
      return '博客配置';
    }
  }

  // 其他类型返回原始 title
  return ref.title;
};

const isImage = (file: FileInfo) => file.file_type?.startsWith('image/');

const formatFileSize = (size: number) => {
  if (size < 1024) return size + ' B';
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB';
  return (size / (1024 * 1024)).toFixed(1) + ' MB';
};

const getStatusTagType = (status: number) => {
  return status === 1 ? 'success' : 'info';
};

const getStatusText = (status: number) => {
  return status === 1 ? '使用中' : '未使用';
};

onMounted(() => {
  syncQuickFiltersFromQuery();
  loadList();
});
</script>

<style scoped lang="scss">
  /* 搜索表单样式已移至全局样式 main.scss */

  .file-list-page {
    height: 100%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .filter-slide-enter-active,
  .filter-slide-leave-active {
    transition: all 0.1s linear;
  }

  .filter-slide-enter-from,
  .filter-slide-leave-to {
    opacity: 0;
    transform: translateY(-4px);
  }

  .filter-slide-enter-to,
  .filter-slide-leave-from {
    opacity: 1;
    transform: translateY(0);
  }

  .file-list-page > :deep(.filter-panel) {
    flex-shrink: 0;
  }

  .file-list-page > :deep(.common-list) {
    flex: 1;
    min-height: 0;
  }

  .references-content {
    min-height: 150px;
    padding: 4px;
  }

  .reference-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    max-height: calc(
      3 * (16px + 12px + 12px + 16px + 24px)
    ); /* 3行的高度: padding + gap + header + body + margin */
    overflow-y: auto;

    /* 美化滚动条 */
    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-track {
      background: #f1f1f1;
      border-radius: 3px;
    }

    &::-webkit-scrollbar-thumb {
      background: #c1c1c1;
      border-radius: 3px;

      &:hover {
        background: #a8a8a8;
      }
    }
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

  /* 移动端优化 */
  @media (max-width: 768px) {
    .reference-item {
      padding: 12px;
    }

    .reference-header {
      flex-wrap: wrap;
      gap: 8px;
      margin-bottom: 10px;

      .reference-field {
        font-size: 12px;
        padding: 2px 6px;
      }
    }

    .reference-body {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;

      .reference-title {
        white-space: normal;
        word-break: break-all;
        line-height: 1.5;
      }

      .reference-link {
        align-self: flex-end;
      }
    }
  }
</style>
