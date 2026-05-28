<!--
项目名称：JeriBlog
文件名称：JsonListEditor.vue
创建时间：2026-04-16 15:08:10

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：公共组件 - JsonListEditor组件
-->

<template>
  <div :class="['json-list-editor', `mobile-layout-${mobileLayout}`]">
    <div
      v-for="(item, index) in internalValue"
      :key="index"
      :class="['editor-item', { 'has-controls': !hideControls }]"
    >
      <!-- 排序按钮 -->
      <template v-if="!hideControls">
        <el-button
          :icon="ArrowUp"
          circle
          size="small"
          @click="moveUp(index)"
          :disabled="disabled || index === 0"
        />
        <el-button
          :icon="ArrowDown"
          circle
          size="small"
          @click="moveDown(index)"
          :disabled="disabled || index === internalValue.length - 1"
          style="margin-left: 0"
        />
      </template>

      <!-- 动态字段 -->
      <template v-for="field in fields" :key="field.key">
        <!-- 文本输入 -->
        <el-input
          v-if="field.type === 'text'"
          v-model="item[field.key]"
          :class="['editor-field', `field-${field.key}`]"
          :placeholder="field.placeholder"
          :style="field.style"
          :disabled="disabled"
          @input="emitUpdate"
        />

        <!-- 下拉选择 -->
        <el-select
          v-else-if="field.type === 'select'"
          v-model="item[field.key]"
          :class="['editor-field', `field-${field.key}`]"
          :placeholder="field.placeholder"
          :style="field.style"
          :disabled="disabled"
          :filterable="field.filterable"
          :allow-create="field.allowCreate"
          @change="emitUpdate"
        >
          <template v-if="field.prefix" #prefix>{{ field.prefix }}</template>
          <el-option
            v-for="option in field.options"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          >
            <template v-if="option.icon">
              <i :class="option.icon" style="margin-right: 8px; font-size: 16px"></i>
              {{ option.label }}
            </template>
          </el-option>
        </el-select>

        <!-- 颜色选择器 -->
        <el-color-picker
          v-else-if="field.type === 'color'"
          v-model="item[field.key]"
          :class="['editor-field', `field-${field.key}`]"
          :disabled="disabled"
          @change="emitUpdate"
        />
      </template>

      <!-- 删除按钮 -->
      <el-button
        v-if="!hideControls"
        type="danger"
        :icon="Delete"
        circle
        size="small"
        @click="removeItem(index)"
        :disabled="disabled"
      />
    </div>

    <!-- 添加按钮行 -->
    <div v-if="!hideControls" class="editor-item add-row">
      <el-button
        type="primary"
        :icon="Plus"
        circle
        size="small"
        @click="addItem"
        :disabled="disabled"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { Delete, Plus, ArrowUp, ArrowDown } from '@element-plus/icons-vue';

export type MobileLayout =
  | 'none'
  | 'name-icon'
  | 'footer-social'
  | 'name-url'
  | 'profile-color';

export interface FieldConfig {
  key: string;
  type: 'text' | 'select' | 'color';
  placeholder?: string;
  style?: string;
  prefix?: string;
  filterable?: boolean;
  allowCreate?: boolean;
  options?: Array<{ label: string; value: string; icon?: string }>;
}

export interface JsonListEditorProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  modelValue: any[];
  fields: FieldConfig[];
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  defaultItem?: Record<string, any>;
  disabled?: boolean;
  hideControls?: boolean;
  mobileLayout?: MobileLayout;
}

const props = withDefaults(defineProps<JsonListEditorProps>(), {
  disabled: false,
  defaultItem: () => ({}),
  hideControls: false,
  mobileLayout: 'none',
});

const emit = defineEmits<{
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  'update:modelValue': [value: any[]];
}>();

// 内部值（深拷贝避免直接修改 prop）
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const internalValue = ref<any[]>([]);

// 监听 modelValue 变化
watch(
  () => props.modelValue,
  newVal => {
    internalValue.value = JSON.parse(JSON.stringify(newVal || []));
  },
  { immediate: true, deep: true }
);

// 发送更新
const emitUpdate = () => {
  emit('update:modelValue', JSON.parse(JSON.stringify(internalValue.value)));
};

// 上移
const moveUp = (index: number) => {
  if (index <= 0) return;
  [internalValue.value[index], internalValue.value[index - 1]] = [
    internalValue.value[index - 1],
    internalValue.value[index],
  ];
  emitUpdate();
};

// 下移
const moveDown = (index: number) => {
  if (index >= internalValue.value.length - 1) return;
  [internalValue.value[index], internalValue.value[index + 1]] = [
    internalValue.value[index + 1],
    internalValue.value[index],
  ];
  emitUpdate();
};

// 删除项
const removeItem = (index: number) => {
  internalValue.value.splice(index, 1);
  emitUpdate();
};

// 添加项
const addItem = () => {
  const newItem = { ...props.defaultItem };
  // 确保所有字段都有默认值
  props.fields.forEach(({ key }) => {
    if (!(key in newItem)) newItem[key] = '';
  });
  internalValue.value.push(newItem);
  emitUpdate();
};
</script>

<style scoped lang="scss">
.json-list-editor {
  width: 100%;

  .editor-item {
    display: flex;
    align-items: center;
    margin-bottom: 12px;
    gap: 8px;

    &.add-row {
      justify-content: flex-end;
    }
  }

  @media (max-width: 768px) {
    .editor-item {
      .editor-field {
        margin: 0 !important;

        :deep(.el-input),
        :deep(.el-select),
        :deep(.el-input__wrapper),
        :deep(.el-select__wrapper) {
          width: 100%;
        }
      }
    }

    &:not(.mobile-layout-none) {
      .editor-item {
        display: grid;
        grid-template-columns: repeat(6, minmax(0, 1fr));
        align-items: flex-start;
        gap: 8px;
        padding: 10px 0;
        border-bottom: 1px solid var(--el-border-color-lighter);

        > .el-button {
          grid-row: 1;
          width: 32px;
          height: 32px;
          min-width: 32px;
          min-height: 32px;
          margin-left: 0;
        }

        > .el-button:nth-of-type(1) {
          grid-column: 1;
        }

        > .el-button:nth-of-type(2) {
          grid-column: 1;
          margin-left: 40px !important;
        }

        > .el-button:last-child {
          grid-column: -2 / -1;
          justify-self: end;
        }

        .editor-field {
          min-width: 0;
        }

        &.add-row {
          justify-content: flex-end;
          padding-top: 8px;
          border-bottom: 0;
        }
      }
    }

    &.mobile-layout-name-icon {
      .field-name,
      .field-icon {
        grid-row: 1;
      }

      .field-name {
        grid-column: 1 / 4;
      }

      .field-icon {
        grid-column: 4 / 7;
      }

      .field-url {
        grid-column: 1 / 7;
        grid-row: 2;
      }

      .editor-item.has-controls {
        .field-name,
        .field-icon {
          grid-row: 2;
        }

        .field-url {
          grid-row: 3;
        }
      }
    }

    &.mobile-layout-footer-social {
      .editor-item {
        grid-template-columns: repeat(3, minmax(0, 1fr));
      }

      .field-name,
      .field-position {
        grid-row: 1;
      }

      .field-name {
        grid-column: 1 / 3;
      }

      .field-position {
        grid-column: 3 / 4;
      }

      .field-icon {
        grid-column: 1 / 4;
        grid-row: 2;
      }

      .field-url {
        grid-column: 1 / 4;
        grid-row: 3;
      }

      .editor-item.has-controls {
        .field-name,
        .field-position {
          grid-row: 2;
        }

        .field-icon {
          grid-row: 3;
        }

        .field-url {
          grid-row: 4;
        }
      }
    }

    &.mobile-layout-name-url {
      .editor-item {
        grid-template-columns: repeat(3, minmax(0, 1fr));
      }

      .field-name {
        grid-column: 1 / 2;
        grid-row: 1;
      }

      .field-url,
      .field-version {
        grid-column: 2 / 4;
        grid-row: 1;
      }

      .editor-item.has-controls {
        .field-name,
        .field-url,
        .field-version {
          grid-row: 2;
        }
      }
    }

    &.mobile-layout-profile-color {
      .editor-item {
        grid-template-columns: minmax(0, 1fr) minmax(0, 2fr) 40px;
      }

      .field-label {
        grid-column: 1 / 2;
        grid-row: 1;
      }

      .field-value {
        grid-column: 2 / 3;
        grid-row: 1;
      }

      .field-color {
        grid-column: 3 / 4;
        grid-row: 1;
      }

      .editor-item.has-controls {
        .field-label,
        .field-value,
        .field-color {
          grid-row: 2;
        }
      }
    }

    &.mobile-layout-name-icon,
    &.mobile-layout-footer-social,
    &.mobile-layout-name-url,
    &.mobile-layout-profile-color {
      .editor-field {
        width: 100% !important;
      }
    }
  }
}
</style>
