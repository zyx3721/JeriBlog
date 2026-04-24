<!--
项目名称：JeriBlog
文件名称：UploadSettingsTab.vue
创建时间：2026-04-24 16:30:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：系统设置 - 上传配置标签页
-->

<template>
  <div class="upload-settings-tab">
    <el-form :model="currentForm" label-width="120px" class="setting-form">
      <!-- 基础配置 -->
      <el-divider content-position="left">基础配置</el-divider>
      <el-form-item label="最大文件大小">
        <el-input-number v-model="currentForm.max_file_size" :min="0" :step="1" />
        <span style="margin-left:8px;color:#909399">MB</span>
      </el-form-item>
      <el-form-item label="文件命名">
        <el-input v-model="currentForm.path_pattern" placeholder="{timestamp}_{random}{ext}" />
        <div style="margin-top: 4px; font-size: 12px; color: #909399;">
          支持变量: {timestamp} {random} {filename} {ext} {type} {userid} YYYY MM DD HH mm ss
        </div>
      </el-form-item>

      <!-- 存储类型选择 -->
      <el-divider content-position="left">存储类型</el-divider>
      <el-form-item label="存储类型">
        <el-select v-model="currentStorageType" placeholder="选择存储类型" style="width: 220px" @change="handleStorageTypeChange">
          <el-option label="本地存储" value="local" />
          <el-option label="亚马逊 S3" value="s3" />
          <el-option label="阿里云 OSS" value="oss" />
          <el-option label="腾讯云 COS" value="cos" />
          <el-option label="七牛云 Kodo" value="kodo" />
          <el-option label="Cloudflare R2" value="r2" />
          <el-option label="MinIO" value="minio" />
        </el-select>
      </el-form-item>

      <!-- 云存储配置 -->
      <template v-if="currentStorageType !== 'local'">
        <el-divider content-position="left">{{ storageTypeLabel }} 配置</el-divider>
        <el-form-item :label="accessLabel">
          <el-input v-model="currentForm.access_key" :placeholder="accessPlaceholder" clearable />
        </el-form-item>
        <el-form-item :label="secretLabel">
          <el-input v-model="currentForm.secret_key" type="password" show-password :placeholder="secretPlaceholder" clearable />
        </el-form-item>
        <el-form-item v-if="showRegion" label="地域">
          <el-input v-model="currentForm.region" :placeholder="regionPlaceholder" clearable />
        </el-form-item>
        <el-form-item label="存储桶">
          <el-input v-model="currentForm.bucket" :placeholder="bucketPlaceholder" clearable />
        </el-form-item>
        <el-form-item v-if="showEndpoint" label="服务端点">
          <el-input v-model="currentForm.endpoint" :placeholder="endpointPlaceholder" clearable />
        </el-form-item>
        <el-form-item label="自定义域名">
          <el-input v-model="currentForm.domain" :placeholder="domainPlaceholder" clearable />
          <div style="margin-top: 4px; font-size: 12px; color: #909399;">
            支持 http:// 或 https:// 开头,如: http://oss.example.com
          </div>
        </el-form-item>
        <el-form-item v-if="showUseSSL" label="启用 HTTPS">
          <el-switch v-model="currentForm.use_ssl" :active-value="true" :inactive-value="false" />
        </el-form-item>
      </template>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'

// 定义存储类型配置接口
interface StorageConfig {
  max_file_size: number
  path_pattern: string
  access_key: string
  secret_key: string
  region: string
  bucket: string
  endpoint: string
  domain: string
  use_ssl: boolean
}

// 定义存储类型键
type StorageType = 'local' | 's3' | 'oss' | 'cos' | 'kodo' | 'r2' | 'minio'

// 定义表单接口
export interface UploadForm {
  storage_type: string
  local: StorageConfig
  s3: StorageConfig
  oss: StorageConfig
  cos: StorageConfig
  kodo: StorageConfig
  r2: StorageConfig
  minio: StorageConfig
}

const props = defineProps<{
  form: UploadForm
  loading?: boolean
}>()

const emit = defineEmits<{
  'update:form': [value: UploadForm]
}>()

// 当前选中的存储类型
const currentStorageType = ref(props.form.storage_type)

// 当前存储类型的配置
const currentForm = computed({
  get: () => {
    const type = currentStorageType.value as StorageType
    return props.form[type]
  },
  set: (val: StorageConfig) => {
    const newForm = { ...props.form }
    const type = currentStorageType.value as StorageType
    newForm[type] = val
    emit('update:form', newForm)
  }
})

// 监听存储类型变化
const handleStorageTypeChange = (newType: string) => {
  const newForm = { ...props.form }
  newForm.storage_type = newType
  emit('update:form', newForm)
}

// 存储类型标签
const storageTypeLabel = computed(() => {
  const labels: Record<string, string> = {
    local: '本地存储',
    s3: '亚马逊 S3',
    oss: '阿里云 OSS',
    cos: '腾讯云 COS',
    kodo: '七牛云 Kodo',
    r2: 'Cloudflare R2',
    minio: 'MinIO'
  }
  return labels[currentStorageType.value] || '云存储'
})

// Access Key 标签
const accessLabel = computed(() => {
  switch (currentStorageType.value) {
    case 'cos': return 'SecretId'
    case 'oss': return 'AccessKeyId'
    case 'kodo': return 'AccessKey'
    case 'r2': return 'Access Key'
    case 'minio': return 'Access Key'
    default: return 'Access Key'
  }
})

// Secret Key 标签
const secretLabel = computed(() => {
  switch (currentStorageType.value) {
    case 'cos': return 'SecretKey'
    case 'oss': return 'AccessKeySecret'
    case 'kodo': return 'SecretKey'
    case 'r2': return 'Secret Key'
    case 'minio': return 'Secret Key'
    default: return 'Secret Key'
  }
})

// Access Key 占位符
const accessPlaceholder = computed(() => {
  switch (currentStorageType.value) {
    case 'cos': return '例如 AKIDxxxxxxxxxxxxxxxxxxxx'
    case 'oss': return '例如 LTAIxxxxxxxxxxxxxxxx'
    default: return ''
  }
})

// Secret Key 占位符
const secretPlaceholder = computed(() => {
  switch (currentStorageType.value) {
    case 'cos': return 'COS 的 SecretKey'
    case 'oss': return 'OSS 的 AccessKeySecret'
    default: return ''
  }
})

// 地域占位符
const regionPlaceholder = computed(() => {
  switch (currentStorageType.value) {
    case 's3': return '例如 us-east-1, ap-southeast-1'
    case 'cos': return '例如 ap-guangzhou, ap-beijing'
    case 'oss': return '例如 cn-hangzhou, cn-beijing'
    case 'kodo': return '例如 cn-east-1, cn-north-1, cn-south-1'
    case 'minio': return '例如 us-east-1, cn-east-1'
    default: return ''
  }
})

// 服务端点占位符
const endpointPlaceholder = computed(() => {
  switch (currentStorageType.value) {
    case 's3': return '可选，例如 s3.us-east-1.amazonaws.com'
    case 'r2': return '例如 <account-id>.r2.cloudflarestorage.com'
    case 'minio': return '例如 localhost:9000 或 minio.example.com'
    default: return ''
  }
})

// 存储桶占位符
const bucketPlaceholder = computed(() => {
  switch (currentStorageType.value) {
    case 's3': return '例如 my-bucket'
    case 'cos': return '例如 my-bucket-1234567890'
    case 'oss': return '例如 my-bucket'
    case 'kodo': return '例如 my-bucket'
    case 'r2': return '例如 my-bucket'
    case 'minio': return '例如 my-bucket'
    default: return '例如 my-bucket'
  }
})

// 自定义域名占位符
const domainPlaceholder = computed(() => {
  switch (currentStorageType.value) {
    case 'kodo': return '必需，例如 https://cdn.example.com (七牛云CDN域名)'
    case 'cos': return '可选，例如 https://cdn.example.com (腾讯云CDN域名)'
    case 'oss': return '可选，例如 http://oss.example.com (阿里云CDN域名)'
    default: return '可选，例如 https://cdn.example.com'
  }
})

// 是否显示地域配置
const showRegion = computed(() => {
  const t = currentStorageType.value
  return t === 's3' || t === 'cos' || t === 'oss' || t === 'kodo' || t === 'minio'
})

// 是否显示服务端点配置
const showEndpoint = computed(() => {
  const t = currentStorageType.value
  return t === 's3' || t === 'r2' || t === 'minio'
})

// 是否显示 SSL 配置
const showUseSSL = computed(() => {
  const t = currentStorageType.value
  return t === 'r2' || t === 'minio'
})

// 监听 props.form 变化,同步 currentStorageType
watch(() => props.form.storage_type, (newType) => {
  currentStorageType.value = newType
})
</script>

<style scoped lang="scss">
.upload-settings-tab {
  .el-divider {
    margin: 24px 0 16px 0;
  }
}
</style>
