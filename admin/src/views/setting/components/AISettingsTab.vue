<!--
项目名称：JeriBlog
文件名称：AISettingsTab.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面组件 - AISettingsTab页面
-->

<template>
  <el-form :model="form" label-width="120px" class="setting-form">
    <el-divider content-position="left">基础配置</el-divider>

    <el-form-item label="API 端点">
      <el-input v-model="form.base_url" placeholder="例如 https://api.deepseek.com" :disabled="loading" />
    </el-form-item>

    <el-form-item label="API 密钥">
      <el-input v-model="form.api_key" type="password" show-password placeholder="输入 API Key" :disabled="loading"
        autocomplete="off" />
    </el-form-item>

    <el-form-item label="模型名称">
      <el-input v-model="form.model" placeholder="例如 deepseek-chat" :disabled="loading" />
    </el-form-item>

    <el-form-item label=" ">
      <el-button :loading="testing" @click="handleTest">测试连接</el-button>
    </el-form-item>

    <el-divider content-position="left">提示词配置</el-divider>

    <el-form-item label="文章摘要提示词">
      <el-input v-model="form.summary_prompt" type="textarea" :rows="5" placeholder="用于生成文章摘要，留空时使用系统默认提示词"
        :disabled="loading" />
    </el-form-item>

    <el-form-item label="AI 总结提示词">
      <el-input v-model="form.ai_summary_prompt" type="textarea" :rows="5" placeholder="用于生成 AI 总结，留空时使用系统默认提示词"
        :disabled="loading" />
    </el-form-item>

    <el-form-item label="标题提示词">
      <el-input v-model="form.title_prompt" type="textarea" :rows="5" placeholder="用于生成标题，留空时使用系统默认提示词"
        :disabled="loading" />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { testAIConfig } from '@/api/ai'

interface AIForm {
  base_url: string
  api_key: string
  model: string
  summary_prompt: string
  ai_summary_prompt: string
  title_prompt: string
}

const form = defineModel<AIForm>('form', { required: true })

defineProps<{
  loading?: boolean
}>()

const testing = ref(false)

async function handleTest() {
  if (!form.value.base_url || !form.value.api_key || !form.value.model) {
    ElMessage.warning('请先填写完整的 API 端点、密钥和模型名称')
    return
  }
  testing.value = true
  try {
    await testAIConfig({
      base_url: form.value.base_url,
      api_key: form.value.api_key,
      model: form.value.model,
    })
    ElMessage.success('连接成功，配置可用')
  } catch (e: any) {
    ElMessage.error(e?.message || '连接失败，请检查配置')
  } finally {
    testing.value = false
  }
}
</script>

<style lang="scss" scoped>
// 移动端适配
@media (max-width: 768px) {
  :deep(.el-form-item__label) {
    width: 100px !important;
    font-size: 13px;
  }
}
</style>
