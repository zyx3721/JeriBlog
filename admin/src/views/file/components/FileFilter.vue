<!--
项目名称：JeriBlog
文件名称：FileFilter.vue
创建时间：2026-04-28 15:30:00

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：文件筛选面板组件
-->

<template>
  <filter-panel v-model="filterForm" title="筛选条件" @reset="handleReset" @close="$emit('close')">
    <el-col :span="5">
      <el-form-item label="关键词">
        <el-input v-model="filterForm.keyword" placeholder="搜索文件名" clearable>
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="文件类型">
        <el-select
          v-model="filterForm.file_type"
          placeholder="选择文件类型"
          clearable
          style="width: 100%"
        >
          <el-option label="图片" value="image" />
          <el-option label="视频" value="video" />
          <el-option label="音频" value="audio" />
          <el-option label="文档" value="application" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="使用状态">
        <el-select
          v-model="filterForm.status"
          placeholder="选择使用状态"
          clearable
          style="width: 100%"
        >
          <el-option label="使用中" :value="1" />
          <el-option label="未使用" :value="0" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="上传用途">
        <el-select
          v-model="filterForm.upload_type"
          placeholder="选择上传用途"
          clearable
          style="width: 100%"
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
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="文件大小">
        <el-select
          v-model="sizeRange"
          placeholder="选择文件大小"
          clearable
          style="width: 100%"
          @change="handleSizeChange"
        >
          <el-option label="< 100KB" value="0-102400" />
          <el-option label="100KB - 1MB" value="102400-1048576" />
          <el-option label="1MB - 10MB" value="1048576-10485760" />
          <el-option label="> 10MB" value="10485760-0" />
        </el-select>
      </el-form-item>
    </el-col>

    <!-- PC 端：日期范围选择器 -->
    <el-col :span="7" class="date-range-pc">
      <el-form-item label="上传时间">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 100%"
          popper-class="date-picker-left"
          @change="handleDateChange"
        />
      </el-form-item>
    </el-col>

    <!-- 移动端：两个独立的单日期选择器 -->
    <el-col :span="12" class="date-range-mobile">
      <el-form-item label="开始日期">
        <el-date-picker
          v-model="startDate"
          type="date"
          placeholder="选择开始日期"
          value-format="YYYY-MM-DD"
          style="width: 100%"
          popper-class="mobile-date-picker"
          @change="handleMobileDateChange"
        />
      </el-form-item>
    </el-col>

    <el-col :span="12" class="date-range-mobile">
      <el-form-item label="结束日期">
        <el-date-picker
          v-model="endDate"
          type="date"
          placeholder="选择结束日期"
          value-format="YYYY-MM-DD"
          style="width: 100%"
          popper-class="mobile-date-picker"
          @change="handleMobileDateChange"
        />
      </el-form-item>
    </el-col>
  </filter-panel>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { Search } from '@element-plus/icons-vue';
import FilterPanel from '@/components/common/FilterPanel.vue';
import type { FileListQuery } from '@/types/file';

/**
 * 防抖函数
 * @param fn 要执行的函数
 * @param delay 延迟时间（毫秒）
 * @returns 防抖后的函数
 */
function debounce<T extends (...args: unknown[]) => unknown>(fn: T, delay: number) {
  let timer: ReturnType<typeof setTimeout> | null = null;
  return function (this: ThisParameterType<T>, ...args: Parameters<T>) {
    if (timer) clearTimeout(timer);
    timer = setTimeout(() => {
      fn.apply(this, args);
    }, delay);
  };
}

/**
 * 组件属性定义
 */
const props = defineProps<{
  modelValue: FileListQuery;
}>();

/**
 * 组件事件定义
 */
const emit = defineEmits<{
  'update:modelValue': [value: FileListQuery];
  search: [];
  close: [];
}>();

const filterForm = ref<FileListQuery>({ ...props.modelValue });
const dateRange = ref<[string, string] | null>(null);
const sizeRange = ref<string | null>(null);

let isExternalUpdate = false;
let isResetting = false;

watch(
  () => props.modelValue,
  newVal => {
    isExternalUpdate = true;
    filterForm.value = { ...newVal };
    if (newVal.start_time && newVal.end_time) {
      dateRange.value = [newVal.start_time, newVal.end_time];
      startDate.value = newVal.start_time;
      endDate.value = newVal.end_time;
    } else {
      dateRange.value = null;
      startDate.value = '';
      endDate.value = '';
    }
    setTimeout(() => {
      isExternalUpdate = false;
    }, 0);
  },
  { deep: true }
);

const debouncedSearch = debounce(() => {
  emit('update:modelValue', { ...filterForm.value });
  emit('search');
}, 500);

watch(
  filterForm,
  () => {
    if (!isExternalUpdate && !isResetting) {
      debouncedSearch();
    }
  },
  { deep: true }
);

/**
 * 处理日期范围变化
 * @param val 日期范围值
 */
const handleDateChange = (val: [string, string] | null) => {
  if (val) {
    filterForm.value.start_time = val[0];
    filterForm.value.end_time = val[1];
    // 同步到移动端
    startDate.value = val[0];
    endDate.value = val[1];
  } else {
    filterForm.value.start_time = undefined;
    filterForm.value.end_time = undefined;
    startDate.value = '';
    endDate.value = '';
  }
};

/**
 * 处理移动端日期变化
 */
const handleMobileDateChange = () => {
  // 情况1：两个日期都清空
  if (!startDate.value && !endDate.value) {
    filterForm.value.start_time = undefined;
    filterForm.value.end_time = undefined;
    dateRange.value = null;
    return
  }

  // 情况2：只选择了开始日期或结束日期，不触发筛选
  if (!startDate.value || !endDate.value) {
    return;
  }

  // 情况3：两个日期都已选择，进行合法性校验
  if (startDate.value && endDate.value) {
    // 时间合法性校验：开始时间不能大于结束时间
    if (startDate.value > endDate.value) {
      ElMessage.error('开始时间不能大于结束时间，请重新选择');
      return;
    }

    // 校验通过，更新筛选条件
    filterForm.value.start_time = startDate.value;
    filterForm.value.end_time = endDate.value;
    // 同步到 PC 端
    dateRange.value = [startDate.value, endDate.value];
  }
};

/**
 * 处理文件大小范围变化
 * @param val 大小范围值
 */
const handleSizeChange = (val: string | null) => {
  if (val) {
    const [min, max] = val.split('-').map(Number);
    filterForm.value.min_size = min;
    filterForm.value.max_size = max || undefined;
  } else {
    filterForm.value.min_size = undefined;
    filterForm.value.max_size = undefined;
  }
};

/**
 * 处理重置
 */
const handleReset = () => {
  isResetting = true;
  dateRange.value = null;
  sizeRange.value = null;
  startDate.value = '';
  endDate.value = '';

  const page = filterForm.value.page;
  const pageSize = filterForm.value.page_size;
  filterForm.value = { page, page_size: pageSize };

  emit('update:modelValue', { ...filterForm.value });
  emit('search');

  setTimeout(() => {
    isResetting = false;
  }, 100);
};

onMounted(() => {
  // 组件挂载时初始化日期范围
  if (filterForm.value.start_time && filterForm.value.end_time) {
    dateRange.value = [filterForm.value.start_time, filterForm.value.end_time];
    startDate.value = filterForm.value.start_time;
    endDate.value = filterForm.value.end_time;
  }
});
</script>

<style scoped>
/* 时间选择器左对齐 */
:deep(.date-picker-left) {
  left: 0 !important;
}

/* 默认显示 PC 端日期范围选择器 */
.date-range-pc {
  display: block;
}

.date-range-mobile {
  display: none;
}

/* 移动端优化 */
@media (max-width: 768px) {
  /* 隐藏 PC 端日期范围选择器 */
  .date-range-pc {
    display: none !important;
  }

  /* 显示移动端单日期选择器 */
  .date-range-mobile {
    display: block !important;
  }

  /* 移动端日期选择器弹出层优化 */
  :deep(.mobile-date-picker.el-popper) {
    /* 弹出层定位优化,防止被输入法遮挡 */
    position: fixed !important;
    transform: translateY(-50%) !important;
    top: 50% !important;
    left: 50% !important;
    margin-left: -47.5vw !important;
    max-width: 95vw !important;
    z-index: 3000 !important;
  }

  /* 日期面板单列紧凑布局 */
  :deep(.mobile-date-picker .el-date-picker) {
    width: 100% !important;
  }

  /* 日期表格优化 */
  :deep(.mobile-date-picker .el-date-table) {
    width: 100% !important;
  }

  /* 月份切换按钮优化 */
  :deep(.mobile-date-picker .el-date-picker__header) {
    display: flex !important;
    justify-content: space-between !important;
    align-items: center !important;
    padding: 8px 12px !important;
  }

  /* 确保切换按钮可点击 */
  :deep(.mobile-date-picker .el-picker-panel__icon-btn) {
    padding: 8px !important;
    min-width: 32px !important;
    min-height: 32px !important;
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
  }
}
</style>

