<!--
项目名称：JeriBlog
文件名称：VisitList.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：页面组件 - VisitList页面
-->

<template>
    <common-list title="访问日志" :data="visitList" :loading="loading" :total="total" :show-create="false"
        v-model:page="queryParams.page" v-model:page-size="queryParams.page_size" @refresh="fetchVisits"
        @update:page="fetchVisits" @update:pageSize="fetchVisits">
        <!-- 搜索表单 -->
        <template #toolbar-before>
            <div class="search-form">
                <el-input
                    v-model="queryParams.keyword"
                    placeholder="搜索访客ID、IP、页面URL、地理位置、浏览器、操作系统、来源..."
                    clearable
                    style="width: 420px"
                    @keyup.enter="handleSearch"
                    @clear="handleSearch"
                />
                <el-date-picker
                    v-model="dateRange"
                    type="datetimerange"
                    range-separator="至"
                    start-placeholder="开始时间"
                    end-placeholder="结束时间"
                    style="width: 360px"
                    @change="handleDateChange"
                    clearable
                    value-format="YYYY-MM-DD HH:mm:ss"
                />
                <el-button type="primary" @click="handleSearch">搜索</el-button>
                <el-button @click="handleReset">重置</el-button>
            </div>
        </template>

        <!-- 表格列 -->
        <el-table-column label="访客ID" width="150" align="center">
            <template #default="{ row }">
                <el-tooltip :content="row.visitor_id" placement="top">
                    <div style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap">
                        {{ row.visitor_id.substring(0, 12) }}...
                    </div>
                </el-tooltip>
            </template>
        </el-table-column>

        <el-table-column label="IP地址" width="140" align="center" prop="ip" />

        <el-table-column label="访问页面" min-width="250" align="center">
            <template #default="{ row }">
                <div style="display: flex; align-items: center; gap: 8px; width: 100%;">
                    <el-tooltip :content="row.page_url" placement="top">
                        <div style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1;">
                            {{ row.page_url }}
                        </div>
                    </el-tooltip>
                </div>
            </template>
        </el-table-column>

        <el-table-column label="地理位置" width="150" align="center" prop="location" />

        <el-table-column label="浏览器" width="140" align="center" prop="browser" />

        <el-table-column label="操作系统" width="120" align="center" prop="os" />

        <el-table-column label="来源" width="250" align="center">
            <template #default="{ row }">
                <el-tooltip v-if="row.referer" :content="row.referer" placement="top">
                    <div style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap">
                        {{ row.referer }}
                    </div>
                </el-tooltip>
                <span v-else style="color: #999">-</span>
            </template>
        </el-table-column>

        <el-table-column label="访问时间" width="180" align="center">
            <template #default="{ row }">
                {{ formatDateTime(row.created_at) }}
            </template>
        </el-table-column>
    </common-list>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import CommonList from '@/components/common/CommonList.vue'
import type { Visit, VisitQuery } from '@/types/stats'
import { getVisits } from '@/api/stats'
import { formatDateTime } from '@/utils/date'

const loading = ref(false)
const visitList = ref<Visit[]>([])
const total = ref(0)
const queryParams = ref<VisitQuery>({
    page: 1,
    page_size: 20,
    keyword: undefined,
    start_date: undefined,
    end_date: undefined
})
const dateRange = ref<[string, string] | null>(null)

const fetchVisits = async () => {
    loading.value = true
    try {
        const [result] = await Promise.all([
            getVisits(queryParams.value),
            new Promise(resolve => setTimeout(resolve, 300))
        ])
        visitList.value = result.list
        total.value = result.total
    } catch {
        ElMessage.error('获取访问日志失败')
    } finally {
        loading.value = false
    }
}

const handleDateChange = (value: [string, string] | null) => {
    if (value) {
        queryParams.value.start_date = value[0]
        queryParams.value.end_date = value[1]
    } else {
        queryParams.value.start_date = undefined
        queryParams.value.end_date = undefined
    }
}

const handleSearch = () => {
    queryParams.value.page = 1
    fetchVisits()
}

const handleReset = () => {
    queryParams.value.keyword = undefined
    queryParams.value.start_date = undefined
    queryParams.value.end_date = undefined
    queryParams.value.page = 1
    dateRange.value = null
    fetchVisits()
}

onMounted(fetchVisits)
</script>

<style scoped>
.search-form {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;
}
</style>
