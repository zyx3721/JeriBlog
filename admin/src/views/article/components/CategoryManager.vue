<!--
项目名称：JeriBlog
文件名称：CategoryManager.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面组件 - CategoryManager页面
-->

<template>
  <el-dialog v-model="visible" title="分类管理" width="800px" :align-center="true">
    <el-table :data="list" style="margin: 20px 0" max-height="350">
      <el-table-column prop="name" label="分类名称" />
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column prop="count" label="文章数" width="100" align="center" />
      <el-table-column prop="sort" label="排序" width="100" align="center" />
      <el-table-column label="操作" width="120" align="center">
        <template #header>
          <el-button type="primary" plain size="small" @click="openForm()">新增</el-button>
        </template>
        <template #default="{ row }">
          <el-button type="primary" link @click="openForm(row)">编辑</el-button>
          <el-button type="danger" link @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="formVisible" :title="current.id ? '编辑' : '新增'" width="400px" append-to-body>
      <el-form :model="current" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="current.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="current.description" type="textarea" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="current.sort" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" @click="save">确定</el-button>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import { getCategories, createCategory, updateCategory, deleteCategory } from '@/api/category'
import type { Category } from '@/types/category'
const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits(['update:modelValue', 'success'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const list = ref<Category[]>([])

const formVisible = ref(false)
const current = ref<Partial<Category>>({ id: 0, name: '', description: '', sort: 0 })

// 初始化加载数据
onMounted(() => {
  loadData()
})

// 加载分类列表
async function loadData() {
  const loading = ElLoading.service()
  try {
    const res = await getCategories()
    list.value = res.list
  } catch (err) {
    ElMessage.error('加载分类列表失败')
  } finally {
    loading.close()
  }
}

// 打开表单
function openForm(row?: Category) {
  if (row) {
    current.value = { ...row }
  } else {
    current.value = { id: 0, name: '', description: '', sort: 0 }
  }
  formVisible.value = true
}

async function remove(row: Category) {
  try {
    await ElMessageBox.confirm('确定要删除这个分类吗？')
    await deleteCategory(row.id)
    await loadData()
    emit('success')
    ElMessage.success('删除成功')
  } catch { }
}

async function save() {
  if (!current.value.name?.trim()) {
    return ElMessage.warning('请输入分类名称')
  }

  const loading = ElLoading.service()
  try {
    if (current.value.id) {
      await updateCategory(current.value.id, current.value)
    } else {
      await createCategory(current.value)
    }
    await loadData()
    formVisible.value = false
    emit('success')
    ElMessage.success('保存成功')
  } catch (err) {
    ElMessage.error('保存失败')
  } finally {
    loading.close()
  }
}
</script>
