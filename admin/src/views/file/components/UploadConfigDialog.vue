<!--
项目名称：JeriBlog
文件名称：UploadConfigDialog.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面组件 - UploadConfigDialog页面
-->

<template>
  <el-dialog v-model="visible" :title="'上传配置'" width="600px" :close-on-click-modal="false" @close="handleClose">
    <el-form ref="formRef" :model="form" label-width="100px">
      <el-form-item label="存储类型">
        <el-select v-model="form.storage_type" placeholder="选择存储类型" style="width: 220px">
          <el-option label="本地存储" value="local" />
          <el-option label="亚马逊 S3" value="s3" />
          <el-option label="阿里云 OSS" value="oss" />
          <el-option label="腾讯云 COS" value="cos" />
          <el-option label="七牛云 Kodo" value="kodo" />
          <el-option label="Cloudflare R2" value="r2" />
          <el-option label="MinIO" value="minio" />
        </el-select>
      </el-form-item>
      <el-form-item label="最大文件大小">
        <el-input-number v-model="form.max_file_size" :min="0" :step="1" />
        <span style="margin-left:8px;color:#909399">MB</span>
      </el-form-item>
      <el-form-item label="文件命名">
        <el-input v-model="form.path_pattern" placeholder="{timestamp}_{random}{ext}" />
      </el-form-item>
      <template v-if="form.storage_type !== 'local'">
        <el-form-item :label="accessLabel">
          <el-input v-model="currentConfig.access_key" :placeholder="accessPlaceholder" clearable />
        </el-form-item>
        <el-form-item :label="secretLabel">
          <el-input v-model="currentConfig.secret_key" type="password" show-password :placeholder="secretPlaceholder"
            clearable />
        </el-form-item>
        <el-form-item v-if="showRegion" label="地域">
          <el-input v-model="currentConfig.region" :placeholder="regionPlaceholder" clearable />
        </el-form-item>
        <el-form-item label="存储桶">
          <el-input v-model="currentConfig.bucket" :placeholder="bucketPlaceholder" clearable />
        </el-form-item>
        <el-form-item v-if="showEndpoint" label="服务端点">
          <el-input v-model="currentConfig.endpoint" :placeholder="endpointPlaceholder" clearable />
        </el-form-item>
        <el-form-item label="自定义域名">
          <el-input v-model="currentConfig.domain" :placeholder="domainPlaceholder" clearable />
        </el-form-item>
        <el-form-item v-if="showUseSSL" label="启用 HTTPS">
          <el-switch v-model="currentConfig.use_ssl" :active-value="true" :inactive-value="false" />
        </el-form-item>
      </template>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSubmit">保存</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getSettingGroup, updateSettingGroup } from '@/api/sysconfig'
const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{ 'update:modelValue': [value: boolean], 'success': [] }>()

const visible = computed({
  get: () => props.modelValue,
  set: (val: boolean) => emit('update:modelValue', val)
})

const formRef = ref()
const saving = ref(false)

const form = ref({
  storage_type: 'local',
  max_file_size: 10,
  path_pattern: '{timestamp}_{random}{ext}',
  // Local 配置
  local: {
    enabled: true
  },
  // S3 配置
  s3: {
    access_key: '',
    secret_key: '',
    region: '',
    bucket: '',
    endpoint: '',
    domain: ''
  },
  // OSS 配置
  oss: {
    access_key: '',
    secret_key: '',
    region: '',
    bucket: '',
    domain: ''
  },
  // COS 配置
  cos: {
    secret_id: '',
    secret_key: '',
    region: '',
    bucket: '',
    domain: ''
  },
  // Kodo 配置
  kodo: {
    access_key: '',
    secret_key: '',
    region: '',
    bucket: '',
    domain: ''
  },
  // R2 配置
  r2: {
    access_key: '',
    secret_key: '',
    bucket: '',
    endpoint: '',
    domain: '',
    use_ssl: true
  },
  // MinIO 配置
  minio: {
    access_key: '',
    secret_key: '',
    region: '',
    bucket: '',
    endpoint: '',
    domain: '',
    use_ssl: true
  }
})

const loadConfigs = async () => {
  try {
    const data = await getSettingGroup('upload')
    // 适配新的扁平化数据格式
    const configs: Record<string, string> = {}
    Object.entries(data).forEach(([key, value]) => {
      const key2 = key.replace('upload.', '') // 移除前缀
      configs[key2] = value
    })

    form.value.storage_type = configs.storage_type || 'local'
    form.value.max_file_size = Number(configs.max_file_size || 10)
    form.value.path_pattern = configs.path_pattern || '{timestamp}_{random}{ext}'

    // Local 配置
    form.value.local.enabled = (configs['local.enabled'] || 'true') === 'true'

    // S3 配置
    form.value.s3.access_key = configs['s3.access_key'] || ''
    form.value.s3.secret_key = configs['s3.secret_key'] || ''
    form.value.s3.region = configs['s3.region'] || ''
    form.value.s3.bucket = configs['s3.bucket'] || ''
    form.value.s3.endpoint = configs['s3.endpoint'] || ''
    form.value.s3.domain = configs['s3.domain'] || ''

    // OSS 配置
    form.value.oss.access_key = configs['oss.access_key'] || ''
    form.value.oss.secret_key = configs['oss.secret_key'] || ''
    form.value.oss.region = configs['oss.region'] || ''
    form.value.oss.bucket = configs['oss.bucket'] || ''
    form.value.oss.domain = configs['oss.domain'] || ''

    // COS 配置
    form.value.cos.secret_id = configs['cos.secret_id'] || ''
    form.value.cos.secret_key = configs['cos.secret_key'] || ''
    form.value.cos.region = configs['cos.region'] || ''
    form.value.cos.bucket = configs['cos.bucket'] || ''
    form.value.cos.domain = configs['cos.domain'] || ''

    // Kodo 配置
    form.value.kodo.access_key = configs['kodo.access_key'] || ''
    form.value.kodo.secret_key = configs['kodo.secret_key'] || ''
    form.value.kodo.region = configs['kodo.region'] || ''
    form.value.kodo.bucket = configs['kodo.bucket'] || ''
    form.value.kodo.domain = configs['kodo.domain'] || ''

    // R2 配置
    form.value.r2.access_key = configs['r2.access_key'] || ''
    form.value.r2.secret_key = configs['r2.secret_key'] || ''
    form.value.r2.bucket = configs['r2.bucket'] || ''
    form.value.r2.endpoint = configs['r2.endpoint'] || ''
    form.value.r2.domain = configs['r2.domain'] || ''
    form.value.r2.use_ssl = (configs['r2.use_ssl'] || 'true') === 'true'

    // MinIO 配置
    form.value.minio.access_key = configs['minio.access_key'] || ''
    form.value.minio.secret_key = configs['minio.secret_key'] || ''
    form.value.minio.region = configs['minio.region'] || ''
    form.value.minio.bucket = configs['minio.bucket'] || ''
    form.value.minio.endpoint = configs['minio.endpoint'] || ''
    form.value.minio.domain = configs['minio.domain'] || ''
    form.value.minio.use_ssl = (configs['minio.use_ssl'] || 'true') === 'true'
  } catch {
    ElMessage.error('获取配置失败')
  }
}

watch(() => visible.value, (val) => { if (val) loadConfigs() })

// 定义配置对象的类型
interface StorageConfig {
  access_key?: string
  secret_key?: string
  region?: string
  bucket?: string
  endpoint?: string
  domain?: string
  use_ssl?: boolean
}

// 根据当前存储类型获取对应的配置对象
const currentConfig = computed<StorageConfig>(() => {
  const storageType = form.value.storage_type
  switch (storageType) {
    case 's3':
      return form.value.s3
    case 'oss':
      return form.value.oss
    case 'cos':
      // COS 使用 secret_id 而不是 access_key
      return {
        access_key: form.value.cos.secret_id,
        secret_key: form.value.cos.secret_key,
        region: form.value.cos.region,
        bucket: form.value.cos.bucket,
        domain: form.value.cos.domain
      }
    case 'kodo':
      return form.value.kodo
    case 'r2':
      return form.value.r2
    case 'minio':
      return form.value.minio
    default:
      return {}
  }
})

// 监听 currentConfig 的变化，同步回原始配置
watch(() => currentConfig.value, (newVal) => {
  const storageType = form.value.storage_type
  if (storageType === 'cos' && newVal.access_key !== undefined) {
    // COS 特殊处理：将 access_key 同步到 secret_id
    form.value.cos.secret_id = newVal.access_key
  }
}, { deep: true })

const handleSubmit = async () => {
  saving.value = true
  try {
    const payload: Record<string, string> = {
      'upload.storage_type': form.value.storage_type,
      'upload.max_file_size': String(form.value.max_file_size),
      'upload.path_pattern': form.value.path_pattern,
      // Local 配置
      'upload.local.enabled': form.value.local.enabled ? 'true' : 'false',
      // S3 配置
      'upload.s3.access_key': form.value.s3.access_key,
      'upload.s3.secret_key': form.value.s3.secret_key,
      'upload.s3.region': form.value.s3.region,
      'upload.s3.bucket': form.value.s3.bucket,
      'upload.s3.endpoint': form.value.s3.endpoint,
      'upload.s3.domain': form.value.s3.domain,
      // OSS 配置
      'upload.oss.access_key': form.value.oss.access_key,
      'upload.oss.secret_key': form.value.oss.secret_key,
      'upload.oss.region': form.value.oss.region,
      'upload.oss.bucket': form.value.oss.bucket,
      'upload.oss.domain': form.value.oss.domain,
      // COS 配置
      'upload.cos.secret_id': form.value.cos.secret_id,
      'upload.cos.secret_key': form.value.cos.secret_key,
      'upload.cos.region': form.value.cos.region,
      'upload.cos.bucket': form.value.cos.bucket,
      'upload.cos.domain': form.value.cos.domain,
      // Kodo 配置
      'upload.kodo.access_key': form.value.kodo.access_key,
      'upload.kodo.secret_key': form.value.kodo.secret_key,
      'upload.kodo.region': form.value.kodo.region,
      'upload.kodo.bucket': form.value.kodo.bucket,
      'upload.kodo.domain': form.value.kodo.domain,
      // R2 配置
      'upload.r2.access_key': form.value.r2.access_key,
      'upload.r2.secret_key': form.value.r2.secret_key,
      'upload.r2.bucket': form.value.r2.bucket,
      'upload.r2.endpoint': form.value.r2.endpoint,
      'upload.r2.domain': form.value.r2.domain,
      'upload.r2.use_ssl': form.value.r2.use_ssl ? 'true' : 'false',
      // MinIO 配置
      'upload.minio.access_key': form.value.minio.access_key,
      'upload.minio.secret_key': form.value.minio.secret_key,
      'upload.minio.region': form.value.minio.region,
      'upload.minio.bucket': form.value.minio.bucket,
      'upload.minio.endpoint': form.value.minio.endpoint,
      'upload.minio.domain': form.value.minio.domain,
      'upload.minio.use_ssl': form.value.minio.use_ssl ? 'true' : 'false'
    }
    await updateSettingGroup('upload', payload)
    ElMessage.success('保存成功（配置已自动热重载）')
    emit('success')
    visible.value = false
  } catch (e) {
    if (e instanceof Error) ElMessage.error(e.message)
    else ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const handleClose = () => { }

const accessLabel = computed(() => {
  switch (form.value.storage_type) {
    case 'cos': return 'SecretId'
    case 'oss': return 'AccessKeyId'
    case 'kodo': return 'AccessKey'
    case 'r2': return 'Access Key'
    case 'minio': return 'Access Key'
    default: return 'Access Key'
  }
})

const secretLabel = computed(() => {
  switch (form.value.storage_type) {
    case 'cos': return 'SecretKey'
    case 'oss': return 'AccessKeySecret'
    case 'kodo': return 'SecretKey'
    case 'r2': return 'Secret Key'
    case 'minio': return 'Secret Key'
    default: return 'Secret Key'
  }
})

const accessPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 'cos': return '例如 AKIDxxxxxxxxxxxxxxxxxxxx'
    case 'oss': return '例如 LTAIxxxxxxxxxxxxxxxx'
    default: return ''
  }
})

const secretPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 'cos': return 'COS 的 SecretKey'
    case 'oss': return 'OSS 的 AccessKeySecret'
    default: return ''
  }
})

const regionPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 's3': return '例如 us-east-1, ap-southeast-1'
    case 'cos': return '例如 ap-guangzhou, ap-beijing'
    case 'oss': return '例如 cn-hangzhou, cn-beijing'
    case 'kodo': return '例如 cn-east-1, cn-north-1, cn-south-1'
    case 'minio': return '例如 us-east-1, cn-east-1'
    default: return ''
  }
})

const endpointPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 's3': return '可选，例如 s3.us-east-1.amazonaws.com'
    case 'r2': return '例如 <account-id>.r2.cloudflarestorage.com'
    case 'minio': return '例如 localhost:9000 或 minio.example.com'
    default: return ''
  }
})

const showRegion = computed(() => {
  const t = form.value.storage_type
  return t === 's3' || t === 'cos' || t === 'oss' || t === 'kodo' || t === 'minio'
})

const showEndpoint = computed(() => {
  const t = form.value.storage_type
  return t === 's3' || t === 'r2' || t === 'minio'
})

const showUseSSL = computed(() => {
  const t = form.value.storage_type
  return t === 'r2' || t === 'minio'
})

const bucketPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 's3': return '例如 my-bucket'
    case 'cos': return '例如 my-bucket-1234567890'
    case 'oss': return '例如 my-bucket'
    case 'kodo': return '例如 my-bucket'
    case 'r2': return '例如 my-bucket'
    case 'minio': return '例如 my-bucket'
    default: return '例如 my-bucket'
  }
})

const domainPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 'kodo': return '必需，例如 https://cdn.example.com (七牛云CDN域名)'
    case 'cos': return '可选，例如 https://cdn.example.com (腾讯云CDN域名)'
    case 'oss': return '可选，例如 https://cdn.example.com (阿里云CDN域名)'
    default: return '可选，例如 https://cdn.example.com'
  }
})
</script>

<style scoped>
.dialog-footer {
  display: inline-flex;
  gap: 12px;
}
</style>