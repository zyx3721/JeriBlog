<!--
项目名称：JeriBlog
文件名称：BasicSettingsTab.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面组件 - BasicSettingsTab页面
-->

<template>
  <el-form :model="form" label-width="120px" class="setting-form">
    <el-divider content-position="left">站长信息</el-divider>

    <el-form-item label="站长姓名">
      <el-input v-model="form.author" placeholder="站长姓名" :disabled="loading" />
    </el-form-item>

    <el-form-item label="站长邮箱">
      <el-input v-model="form.author_email" placeholder="站长联系邮箱" :disabled="loading" />
    </el-form-item>

    <el-form-item label="站长简介">
      <el-input v-model="form.author_desc" type="textarea" :rows="3" placeholder="站长个人简介" :disabled="loading" />
    </el-form-item>

    <div class="image-row">
      <el-form-item label="站长头像">
        <div class="image-upload-wrapper">
          <ImageUploader ref="authorAvatarUploaderRef" v-model="form.author_avatar" upload-type="站长头像" width="120px"
            height="120px" />
          <el-button class="select-file-btn" @click="openFilePicker('author_avatar')">
            <i class="ri-folder-image-line"></i>
            选择文件
          </el-button>
        </div>
      </el-form-item>

      <el-form-item label="站长形象">
        <div class="image-upload-wrapper">
          <ImageUploader ref="authorPhotoUploaderRef" v-model="form.author_photo" upload-type="站长形象" width="80px"
            height="120px" />
          <el-button class="select-file-btn" @click="openFilePicker('author_photo')">
            <i class="ri-folder-image-line"></i>
            选择文件
          </el-button>
        </div>
      </el-form-item>
    </div>

    <el-divider content-position="left">备案信息</el-divider>

    <el-form-item label="ICP备案号">
      <el-input v-model="form.icp" placeholder="ICP备案号" :disabled="loading" />
    </el-form-item>

    <el-form-item label="公安备案号">
      <el-input v-model="form.police_record" placeholder="公安备案号" :disabled="loading" />
    </el-form-item>

    <el-divider content-position="left">系统地址</el-divider>

    <el-form-item label="管理地址">
      <el-input v-model="form.admin_url" placeholder="例如 https://admin.your-site.com" :disabled="loading" />
    </el-form-item>

    <el-form-item label="博客地址">
      <el-input v-model="form.blog_url" placeholder="例如 https://blog.your-site.com" :disabled="loading" />
    </el-form-item>

    <el-form-item label="主页地址">
      <el-input v-model="form.home_url" placeholder="例如 https://your-site.com" :disabled="loading" />
    </el-form-item>
  </el-form>

  <!-- 文件选择对话框 -->
  <FilePickerDialog v-model="filePickerVisible" file-type="image" @confirm="handleFileSelect" />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import ImageUploader from '@/components/common/ImageUploader.vue'
import FilePickerDialog from '@/components/common/FilePickerDialog.vue'
import type { FileInfo } from '@/types/file'

interface BasicForm {
  author: string
  author_email: string
  author_desc: string
  author_avatar: string
  author_photo: string
  icp: string
  police_record: string
  admin_url: string
  blog_url: string
  home_url: string
}

const form = defineModel<BasicForm>('form', { required: true })

defineProps<{
  loading?: boolean
}>()

// 图片上传器引用
const authorAvatarUploaderRef = ref<InstanceType<typeof ImageUploader>>()
const authorPhotoUploaderRef = ref<InstanceType<typeof ImageUploader>>()

// 文件选择对话框
const filePickerVisible = ref(false)
const currentField = ref<'author_avatar' | 'author_photo'>('author_avatar')

// 打开文件选择对话框
const openFilePicker = (field: 'author_avatar' | 'author_photo') => {
  currentField.value = field
  filePickerVisible.value = true
}

// 处理文件选择
const handleFileSelect = (file: FileInfo) => {
  form.value[currentField.value] = file.file_url
  ElMessage.success('已选择文件')
}

// 暴露给父组件使用
defineExpose({
  authorAvatarUploaderRef,
  authorPhotoUploaderRef
})
</script>

<style lang="scss" scoped>
.setting-form {
  .image-row {
    display: flex;
    gap: 40px;

    .el-form-item {
      margin-bottom: 22px;
    }
  }

  .image-upload-wrapper {
    display: flex;
    flex-direction: column;
    gap: 8px;

    .select-file-btn {
      width: 100%;

      i {
        margin-right: 4px;
      }
    }
  }
}

// 移动端适配
@media (max-width: 768px) {
  .setting-form {
    .image-row {
      flex-direction: column;
      gap: 0;
    }
  }

  :deep(.el-form-item__label) {
    width: 100px !important;
    font-size: 13px;
  }
}
</style>
