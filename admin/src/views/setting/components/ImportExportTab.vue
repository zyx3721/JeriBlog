<!--
项目名称：JeriBlog
文件名称：ImportExportTab.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面组件 - ImportExportTab页面
-->

<template>
  <el-form label-width="100px" class="setting-form">
    <el-form-item label="文章数据">
      <el-button type="primary" :disabled="readonly" @click="articleImportVisible = true"
        >导入文章</el-button
      >
    </el-form-item>

    <el-form-item label="评论数据">
      <el-button type="primary" :disabled="readonly" @click="commentImportVisible = true"
        >导入评论</el-button
      >
    </el-form-item>
  </el-form>

  <!-- 文章导入对话框 -->
  <el-dialog
    v-model="articleImportVisible"
    title="导入文章"
    width="90%"
    style="max-width: 600px"
    :close-on-click-modal="false"
    align-center
    class="import-dialog"
  >
    <div class="dialog-scroll-content">
      <el-form label-width="100px">
        <el-form-item label="数据来源">
          <el-select
            v-model="articleSourceType"
            placeholder="请选择数据来源"
            style="width: 100%"
            :disabled="readonly"
          >
            <el-option label="Markdown 格式" value="markdown" />
            <el-option label="Hexo 格式" value="hexo" />
          </el-select>
          <div class="form-tip">Hexo 格式需要包含 Front Matter，Markdown 格式 仅需 Markdown 内容</div>
        </el-form-item>

        <el-form-item label="上传文件">
          <el-upload
            :auto-upload="false"
            :file-list="articleFileList"
            :on-change="handleArticleFileChange"
            :on-remove="handleArticleFileRemove"
            accept=".md,.markdown"
            :limit="100"
            multiple
            drag
            :disabled="readonly"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">拖拽或点击选择文件</div>
            <template #tip>
              <div class="el-upload__tip">最多添加 100 个文件，如遇上传失败请减少数量后重试</div>
            </template>
          </el-upload>
        </el-form-item>

        <el-form-item label="图片处理" class="image-process-item">
          <div class="switch-wrapper">
            <el-switch v-model="articleUploadImages" :disabled="readonly" />
            <div class="form-tip switch-tip">开启后将自动下载并上传文章中的图片</div>
          </div>
        </el-form-item>

        <el-form-item label="图片代理" v-if="articleUploadImages">
          <el-input
            v-model="articleImageProxy"
            placeholder="https://gh.llkk.cc/"
            clearable
            style="width: 100%"
            class="proxy-input"
          >
            <template #prepend>
              <span class="prepend-text">代理地址</span>
            </template>
          </el-input>
          <div class="form-tip">
            用于加速 GitHub 等国外图片下载，留空则使用默认代理：https://gh.llkk.cc/
          </div>
        </el-form-item>
      </el-form>

      <el-alert
        v-if="articleImportResult"
        :type="articleImportResult.failed > 0 ? 'warning' : 'success'"
        :closable="false"
        style="margin-top: 16px"
      >
        <div>成功 {{ articleImportResult.success }} 篇，失败 {{ articleImportResult.failed }} 篇</div>
        <div
          v-if="articleImportResult.errors?.length"
          style="margin-top: 8px; font-size: 12px; color: #909399"
        >
          <div v-for="(err, i) in articleImportResult.errors" :key="i">
            {{ err.filename }}: {{ err.error }}
          </div>
        </div>
      </el-alert>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="articleImportVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="articleUploading"
          :disabled="readonly || articleFileList.length === 0"
          @click="handleArticleImport"
        >
          {{ articleUploading ? '导入中...' : '开始导入' }}
        </el-button>
      </span>
    </template>
  </el-dialog>

  <!-- 评论导入对话框 -->
  <el-dialog
    v-model="commentImportVisible"
    title="导入评论"
    width="90%"
    style="max-width: 600px"
    :close-on-click-modal="false"
    align-center
    class="import-dialog"
  >
    <div class="dialog-scroll-content">
      <el-form label-width="100px">
        <el-form-item label="数据来源">
          <el-select
            v-model="commentSourceType"
            placeholder="请选择数据来源"
            style="width: 100%"
            :disabled="readonly"
          >
            <el-option label="Artalk" value="artalk" />
          </el-select>
          <div class="form-tip">选择评论数据的来源系统，目前支持 Artalk 评论系统</div>
        </el-form-item>

        <el-form-item label="上传文件">
          <el-upload
            :auto-upload="false"
            :file-list="commentFileList"
            :on-change="handleCommentFileChange"
            :on-remove="handleCommentFileRemove"
            accept=".json,.artrans"
            :limit="1"
            drag
            :disabled="readonly"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">拖拽或点击选择文件</div>
            <template #tip>
              <div class="el-upload__tip">支持 JSON 或 Artrans 格式文件，单个文件最大 10MB</div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>

      <el-alert
        v-if="commentImportResult"
        :type="commentImportResult.failed > 0 ? 'warning' : 'success'"
        :closable="false"
        style="margin-top: 16px"
      >
        <div>
          <strong>导入完成</strong>
        </div>
        <div style="margin-top: 8px">
          总计 {{ commentImportResult.total }} 条，成功 {{ commentImportResult.success }} 条，失败
          {{ commentImportResult.failed }}
          条
        </div>
        <div
          v-if="commentImportResult.user_created > 0"
          style="margin-top: 4px; font-size: 12px; color: #909399"
        >
          自动创建了 {{ commentImportResult.user_created }} 个游客用户账号
        </div>
        <div
          v-if="commentImportResult.errors?.length"
          style="
            margin-top: 12px;
            font-size: 12px;
            color: #909399;
            max-height: 200px;
            overflow-y: auto;
          "
        >
          <div><strong>失败详情：</strong></div>
          <div v-for="(err, i) in commentImportResult.errors" :key="i" style="margin-top: 4px">
            第 {{ err.index + 1 }} 条: {{ err.error }}
          </div>
        </div>
      </el-alert>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="commentImportVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="commentUploading"
          :disabled="readonly || commentFileList.length === 0 || !commentSourceType"
          @click="handleCommentImport"
        >
          {{ commentUploading ? '导入中...' : '开始导入' }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { UploadFilled } from '@element-plus/icons-vue';
import type { UploadUserFile, UploadFile } from 'element-plus';
import { importArticles } from '@/api/article';
import { importComments } from '@/api/comment';
import type { ImportArticlesResult } from '@/types/article';
import type { ImportCommentsResult } from '@/types/comment';

const props = withDefaults(defineProps<{ readonly?: boolean }>(), {
  readonly: false,
});

const emit = defineEmits<{
  'import-success': [];
}>();

// 文章导入相关
const articleImportVisible = ref(false);
const articleFileList = ref<UploadUserFile[]>([]);
const articleUploading = ref(false);
const articleImportResult = ref<ImportArticlesResult | undefined>();
const articleSourceType = ref<string>('markdown');
const articleUploadImages = ref(false);
const articleImageProxy = ref<string>('');

const handleArticleFileChange = (file: UploadFile, files: UploadUserFile[]) => {
  articleFileList.value = files;
};

const handleArticleFileRemove = (file: UploadFile, files: UploadUserFile[]) => {
  articleFileList.value = files;
};

const handleArticleImport = async () => {
  if (props.readonly) {
    ElMessage.warning('仅超级管理员可执行导入操作');
    return;
  }

  if (articleFileList.value.length === 0) return;

  try {
    articleUploading.value = true;
    articleImportResult.value = undefined;

    const formData = new FormData();
    formData.append('source_type', articleSourceType.value);
    formData.append('upload_images', String(articleUploadImages.value));
    // 传递图片代理地址，留空则后端使用默认值
    if (articleUploadImages.value && articleImageProxy.value) {
      formData.append('image_proxy', articleImageProxy.value.trim());
    }

    articleFileList.value.forEach(file => {
      if (file.raw) formData.append('files', file.raw);
    });

    const result = await importArticles(formData);
    articleImportResult.value = result;

    if (result.failed === 0) {
      ElMessage.success(`成功导入 ${result.success} 篇文章`);
      emit('import-success');
    } else if (result.success > 0) {
      ElMessage.warning(`导入完成：成功 ${result.success} 篇，失败 ${result.failed} 篇`);
      emit('import-success');
    } else {
      ElMessage.error('导入失败');
    }
  } catch (error: unknown) {
    ElMessage.error((error as Error)?.message || '导入失败');
  } finally {
    articleUploading.value = false;
  }
};

watch(articleImportVisible, val => {
  if (!val) {
    setTimeout(() => {
      articleFileList.value = [];
      articleImportResult.value = undefined;
      articleSourceType.value = 'markdown';
      articleUploadImages.value = false;
    }, 300);
  }
});

// 评论导入相关
const commentImportVisible = ref(false);
const commentFileList = ref<UploadUserFile[]>([]);
const commentUploading = ref(false);
const commentImportResult = ref<ImportCommentsResult | undefined>();
const commentSourceType = ref<string>('artalk');

const handleCommentFileChange = (file: UploadFile, files: UploadUserFile[]) => {
  commentFileList.value = files;
};

const handleCommentFileRemove = (file: UploadFile, files: UploadUserFile[]) => {
  commentFileList.value = files;
};

const handleCommentImport = async () => {
  if (props.readonly) {
    ElMessage.warning('仅超级管理员可执行导入操作');
    return;
  }

  if (commentFileList.value.length === 0) {
    ElMessage.warning('请选择要导入的文件');
    return;
  }

  if (!commentSourceType.value) {
    ElMessage.warning('请选择数据来源');
    return;
  }

  try {
    commentUploading.value = true;
    commentImportResult.value = undefined;

    const formData = new FormData();
    const rawFile = commentFileList.value[0]?.raw;
    if (!rawFile) {
      ElMessage.error('文件读取失败');
      return;
    }
    formData.append('file', rawFile);
    formData.append('source_type', commentSourceType.value);

    const result = await importComments(formData);
    commentImportResult.value = result;

    if (result.failed === 0) {
      ElMessage.success(`成功导入 ${result.success} 条评论`);
      emit('import-success');
    } else if (result.success > 0) {
      ElMessage.warning(`导入完成：成功 ${result.success} 条，失败 ${result.failed} 条`);
      emit('import-success');
    } else {
      ElMessage.error('导入失败，请检查文件格式');
    }
  } catch (error: unknown) {
    ElMessage.error((error as Error)?.message || '导入失败');
  } finally {
    commentUploading.value = false;
  }
};

watch(commentImportVisible, val => {
  if (!val) {
    setTimeout(() => {
      commentFileList.value = [];
      commentImportResult.value = undefined;
      commentSourceType.value = 'artalk';
    }, 300);
  }
});
</script>

<style lang="scss" scoped>
:deep(.el-icon--upload) {
  font-size: 40px;
  color: #409eff;
  margin-bottom: 12px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  margin-top: 8px;
}

// 图片处理开关布局
.switch-wrapper {
  display: flex;
  align-items: center;

  .switch-tip {
    margin-top: 0;
    margin-left: 12px;
  }
}

// 导入对话框样式
.import-dialog {
  :deep(.el-dialog) {
    max-height: 85vh;
  }

  .dialog-scroll-content {
    max-height: calc(85vh - 140px);
    overflow-y: auto;
    overflow-x: hidden;
    padding-right: 8px;
    -webkit-overflow-scrolling: touch;
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
}

// 移动端适配
@media (max-width: 768px) {
  .import-dialog {
    :deep(.el-dialog) {
      max-height: 90vh;
    }

    .dialog-scroll-content {
      max-height: calc(90vh - 140px);
      padding-right: 4px;
    }

    // 对话框标题
    :deep(.el-dialog__header) {
      padding: 16px;

      .el-dialog__title {
        font-size: 16px;
      }
    }

    // 对话框内容
    :deep(.el-dialog__body) {
      padding: 16px;
    }

    // 对话框底部
    :deep(.el-dialog__footer) {
      padding: 12px 16px;
      border-top: 1px solid #f0f0f0;
    }

    .dialog-footer {
      width: 100%;
      gap: 10px;

      .el-button {
        padding: 10px 20px;
        font-size: 14px;
      }
    }

    // 表单标签宽度
    :deep(.el-form-item) {
      margin-bottom: 18px;
    }

    :deep(.el-form-item__label) {
      width: 80px !important;
      font-size: 14px;
      line-height: 32px;
    }

    // 表单内容
    :deep(.el-form-item__content) {
      margin-left: 80px !important;
    }

    // 下拉选择框
    :deep(.el-select) {
      width: 100%;
    }

    // 输入框
    :deep(.el-input) {
      font-size: 14px;

      .el-input__inner {
        font-size: 14px;
        height: 36px;
        line-height: 36px;
      }

      .el-input-group__prepend {
        font-size: 13px;
        padding: 0 12px;
      }
    }

    // 移动端隐藏图片代理输入框的前置文本
    .proxy-input {
      :deep(.el-input-group__prepend) {
        display: none;
      }

      :deep(.el-input__wrapper) {
        border-radius: 4px;
      }
    }

    // 上传组件
    :deep(.el-upload) {
      width: 100%;

      .el-upload-dragger {
        width: 100%;
        padding: 24px 12px;
      }
    }

    :deep(.el-icon--upload) {
      font-size: 36px;
      margin-bottom: 10px;
    }

    :deep(.el-upload__text) {
      font-size: 13px;
      line-height: 1.5;
    }

    :deep(.el-upload__tip) {
      font-size: 11px;
      margin-top: 6px;
    }

    // 文件列表
    :deep(.el-upload-list) {
      margin-top: 10px;

      .el-upload-list__item {
        font-size: 13px;
        padding: 8px 10px;
        line-height: 1.5;

        .el-upload-list__item-name {
          max-width: calc(100% - 70px);
        }
      }
    }

    // 开关组件
    :deep(.el-switch) {
      margin-right: 10px;
    }

    // 图片处理项移动端布局调整
    .image-process-item {
      .switch-wrapper {
        flex-direction: column;
        align-items: flex-start;

        .switch-tip {
          margin-left: 0;
        }
      }
    }

    // 提示文本
    .form-tip {
      font-size: 11px;
      margin-top: 8px;
      line-height: 1.5;
    }

    // 结果提示框
    :deep(.el-alert) {
      padding: 12px;
      font-size: 13px;
      line-height: 1.6;

      .el-alert__content {
        font-size: 13px;
      }

      .el-alert__description {
        font-size: 12px;
        margin-top: 8px;
      }
    }

    // 错误详情滚动区域
    :deep(.el-alert) {
      div[style*='max-height: 200px'] {
        max-height: 160px !important;
        font-size: 11px;
        line-height: 1.5;
      }
    }
  }
}

// 小屏幕进一步优化
@media (max-width: 480px) {
  .import-dialog {
    :deep(.el-dialog) {
      max-height: 92vh;
    }

    .dialog-scroll-content {
      max-height: calc(92vh - 130px);
    }

    :deep(.el-dialog__header) {
      padding: 14px;

      .el-dialog__title {
        font-size: 15px;
      }
    }

    :deep(.el-dialog__body) {
      padding: 14px;
    }

    :deep(.el-dialog__footer) {
      padding: 10px 14px;
    }

    .dialog-footer {
      gap: 8px;

      .el-button {
        flex: 1;
        padding: 9px 16px;
        font-size: 13px;
      }
    }

    // 表单标签宽度
    :deep(.el-form-item) {
      margin-bottom: 16px;
    }

    :deep(.el-form-item__label) {
      width: 70px !important;
      font-size: 13px;
      line-height: 32px;
    }

    // 表单内容
    :deep(.el-form-item__content) {
      margin-left: 70px !important;
    }

    :deep(.el-input) {
      .el-input__inner {
        height: 34px;
        line-height: 34px;
        font-size: 13px;
      }

      .el-input-group__prepend {
        font-size: 12px;
        padding: 0 10px;
      }
    }

    // 移动端隐藏图片代理输入框的前置文本
    .proxy-input {
      :deep(.el-input-group__prepend) {
        display: none;
      }

      :deep(.el-input__wrapper) {
        border-radius: 4px;
      }
    }

    :deep(.el-upload-dragger) {
      padding: 20px 10px;
    }

    :deep(.el-icon--upload) {
      font-size: 32px;
      margin-bottom: 8px;
    }

    :deep(.el-upload__text) {
      font-size: 12px;
    }

    :deep(.el-upload__tip) {
      font-size: 10px;
      margin-top: 4px;
    }

    :deep(.el-upload-list__item) {
      font-size: 12px;
      padding: 6px 8px;
    }

    .form-tip {
      font-size: 10px;
      margin-top: 6px;
    }

    :deep(.el-alert) {
      padding: 10px;
      font-size: 12px;

      div[style*='max-height: 200px'] {
        max-height: 140px !important;
        font-size: 10px;
      }
    }
  }
}
</style>
