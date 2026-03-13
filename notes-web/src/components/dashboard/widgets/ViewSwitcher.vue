<template>
  <div class="view-switcher">
    <el-radio-group
      :model-value="modelValue"
      size="small"
      @change="handleChange"
    >
      <el-tooltip
        v-for="view in availableViewsList"
        :key="view.value"
        :content="view.label"
        placement="top"
      >
        <el-radio-button :value="view.value">
          <el-icon>
            <component :is="view.icon" />
          </el-icon>
        </el-radio-button>
      </el-tooltip>
    </el-radio-group>
  </div>
</template>

<script lang="ts">
import {
  Grid,
  List,
  Calendar,
  TrendCharts,
  Memo,
  Clock,
  Histogram,
} from '@element-plus/icons-vue'

// 视图类型定义
export type ViewType = 'grid' | 'list' | 'table' | 'tree' | 'chart' | 'kanban' | 'timeline' | 'calendar'

export interface ViewOption {
  value: ViewType
  label: string
  icon: any
  description: string
}

// 所有视图选项
export const ALL_VIEWS: ViewOption[] = [
  { value: 'grid', label: '网格视图', icon: Grid, description: '卡片式网格布局' },
  { value: 'list', label: '列表视图', icon: List, description: '单列列表布局' },
  { value: 'table', label: '表格视图', icon: Histogram, description: '结构化表格数据' },
  { value: 'tree', label: '树形视图', icon: Memo, description: '层级树形结构' },
  { value: 'chart', label: '图表视图', icon: TrendCharts, description: '数据可视化图表' },
  { value: 'kanban', label: '看板视图', icon: Calendar, description: '拖拽式看板' },
  { value: 'timeline', label: '时间线', icon: Clock, description: '时间顺序展示' },
]
</script>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  modelValue: ViewType
  availableViews?: ViewType[]
}

const props = withDefaults(defineProps<Props>(), {
  availableViews: () => ['grid', 'list', 'table'],
})

const emit = defineEmits<{
  'update:modelValue': [value: ViewType]
}>()

// 可用的视图选项
const availableViewsList = computed(() => {
  return ALL_VIEWS.filter(view => props.availableViews.includes(view.value))
})

function handleChange(value: string | number | boolean | undefined) {
  if (typeof value === 'string') {
    emit('update:modelValue', value as ViewType)
    
    // 保存到 localStorage
    if (props.availableViews.length > 0) {
      const key = `view-preference-${window.location.pathname}`
      localStorage.setItem(key, value)
    }
  }
}
</script>

<style scoped>
.view-switcher {
  display: inline-flex;
}

.el-radio-group {
  gap: 0;
}

:deep(.el-radio-button__inner) {
  padding: 8px 12px;
}
</style>
