<!--
项目名称：JeriBlog
文件名称：CommonList.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：公共组件 - CommonList组件
-->

<template>
    <div class="common-list">
        <el-card>
            <!-- 工具栏 -->
            <div class="toolbar">
                <h2>{{ title }}</h2>
                <div class="actions">
                    <!-- 前工具栏 -->
                    <slot name="toolbar-before" />
                    <el-button v-if="showCreate" type="primary" @click="$emit('create')">
                        {{ createText }}
                    </el-button>
                    <!-- 后工具栏 -->
                    <slot name="toolbar-after" />
                    <el-button class="refresh-btn" @click="$emit('refresh')">
                        <el-icon>
                            <Refresh />
                        </el-icon>
                    </el-button>
                </div>
            </div>

            <!-- 额外内容 -->
            <slot name="extra" />

            <!-- 表格区域 -->
            <div class="table-wrapper">
                <!-- 加载状态 -->
                <div v-if="loading" class="common-list-loading">
                    <el-skeleton :rows="5" animated />
                </div>

                <!-- 表格 - 完全由外部控制 -->
                <el-table v-else :data="data" border style="width: 100%; height: 100%" v-bind="$attrs" :row-key="rowKey">
                    <slot />
                </el-table>
            </div>

            <!-- 分页 -->
            <div v-if="showPagination" class="pagination">
                <el-pagination :current-page="page" :page-size="pageSize" :page-sizes="[10, 20, 50, 100]" :total="total"
                    layout="total, sizes, prev, pager, next" @current-change="$emit('update:page', $event)"
                    @size-change="$emit('update:pageSize', $event)" />
            </div>
        </el-card>
    </div>
</template>

<script setup lang="ts">
import { Refresh } from '@element-plus/icons-vue'

withDefaults(defineProps<{
    title: string
    data: any[]
    loading?: boolean
    total?: number
    page?: number
    pageSize?: number
    showPagination?: boolean
    showCreate?: boolean
    createText?: string
    rowKey?: string
}>(), {
    loading: false,
    total: 0,
    page: 1,
    pageSize: 10,
    showPagination: true,
    showCreate: true,
    createText: '新增',
    rowKey: 'id'
})

defineEmits<{
    create: []
    refresh: []
    'update:page': [page: number]
    'update:pageSize': [size: number]
}>()
</script>

<style scoped lang="scss">
.common-list {
    height: 100%;

    :deep(.el-card) {
        height: 100%;
        display: flex;
        flex-direction: column;

        .el-card__body {
            flex: 1;
            display: flex;
            flex-direction: column;
            overflow: hidden;
        }
    }

    .toolbar {
        margin-bottom: 12px;
        display: flex;
        justify-content: space-between;
        align-items: center;

        h2 {
            margin: 0;
            font-size: 20px;
            font-weight: 500;
        }

        .actions {
            display: flex;
            gap: 12px;

            :deep(.el-button + .el-button) {
                margin-left: 0;
            }
        }

        @media (max-width: 767px) {
            flex-direction: column;
            align-items: flex-start;
            gap: 12px;

            h2 {
                font-size: 18px;
            }

            .actions {
                width: 100%;
                flex-wrap: wrap;

                .refresh-btn {
                    display: none;
                }
            }
        }

        @media (max-width: 480px) {
            .actions {
                /* 用户管理：搜索表单 + 新增用户按钮三等分 */
                &:has(.user-search) {
                    .user-search {
                        width: 100%;
                        display: flex;
                        flex-wrap: wrap;
                        gap: 12px;
                    }
                    /* 新增用户按钮与搜索/重置按钮同行三等分 */
                    > .el-button {
                        flex: 1 1 calc(33.333% - 8px);
                        min-width: 0;
                    }
                }

                /* 文件管理：搜索表单 + 上传配置按钮三等分 */
                &:has(.file-search) {
                    .file-search {
                        width: 100%;
                        display: flex;
                        flex-wrap: wrap;
                        gap: 12px;
                    }
                    /* 上传配置按钮与搜索/重置按钮同行三等分 */
                    > .el-button {
                        flex: 1 1 calc(33.333% - 8px);
                        min-width: 0;
                    }
                }

                /* 友链管理：新增友链和类型管理按钮同行平分 */
                &:has(.friend-search) {
                    .friend-search {
                        width: 100%;
                    }
                    /* 新增友链和类型管理按钮同行平分 */
                    > .el-button {
                        flex: 1 1 calc(50% - 6px);
                        min-width: 0;
                    }
                }

                /* 文章管理：新增文章、分类管理、标签管理按钮同行三等分 */
                &:has(.article-search) {
                    .article-search {
                        width: 100%;
                    }
                    /* 新增文章、分类管理、标签管理按钮同行三等分 */
                    > .el-button {
                        flex: 1 1 calc(33.333% - 8px);
                        min-width: 0;
                    }
                }

                /* RSS订阅：多个按钮自适应 */
                &:has(.rss-search) {
                    .rss-search {
                        width: 100%;
                    }
                    /* 本站订阅、立即抓取RSS、全部已读按钮自适应 */
                    > .el-button,
                    > .unread-badge .el-button {
                        flex: 1 1 auto;
                        min-width: 0;
                    }
                }
            }
        }
    }

    .table-wrapper {
        flex: 1;
        overflow: auto;

        :deep(.el-table__header th .cell) {
            text-align: center;
        }
    }

    .pagination {
        display: flex;
        justify-content: flex-end;
        padding-top: 12px;

        @media (max-width: 767px) {
            justify-content: center;
        }
    }
}
</style>
