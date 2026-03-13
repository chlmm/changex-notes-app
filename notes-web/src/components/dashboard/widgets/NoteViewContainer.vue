<template>
  <div class="note-view-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">{{ title }}</h1>
        <span v-if="showStats" class="page-stats">共 {{ total }} 条</span>
      </div>
      <div class="header-right">
        <OpenModeSelector v-model="openMode" />
        <ViewSwitcher
          v-model="currentView"
          :available-views="availableViews"
        />
        <slot name="actions" />
      </div>
    </div>

    <!-- 视图内容 -->
    <div class="view-content">
      <!-- GridView -->
      <GridView
        v-if="currentView === 'grid'"
        :items="items"
        :loading="loading"
        :empty-text="emptyText"
        :filters="gridFilters"
        :card-type="viewConfig.cardType"
        :grid-columns="viewConfig.gridColumns"
        @item-click="handleItemClick"
      >
        <template #item="{ item }">
          <slot name="gridItem" :item="item" />
        </template>
      </GridView>

      <!-- ListView -->
      <ListView
        v-else-if="currentView === 'list'"
        :items="items"
        :loading="loading"
        :empty-text="emptyText"
        :title-field="viewConfig.titleField"
        :desc-field="viewConfig.descField"
        :filters="viewConfig.filters"
        @item-click="handleItemClick"
      >
        <template #item="{ item }">
          <slot name="listItem" :item="item" />
        </template>
      </ListView>

      <!-- TableView -->
      <TableView
        v-else-if="currentView === 'table'"
        :columns="tableColumns"
        :data="items"
        :loading="loading"
        :filters="viewConfig.filters"
        :pagination="true"
        @row-click="handleItemClick"
      >
        <template v-if="$slots.tableOperation" #operation="{ row }">
          <slot name="tableOperation" :row="row" />
        </template>
      </TableView>

      <!-- TreeView -->
      <TreeView
        v-else-if="currentView === 'tree'"
        :data="treeData"
        :loading="loading"
        :height="viewConfig.height"
        :title-field="viewConfig.titleField"
        :desc-field="viewConfig.descField"
        @node-click="handleItemClick"
      >
        <template #content="{ item }">
          <slot name="treeContent" :item="item" />
        </template>
      </TreeView>

      <!-- KanbanView -->
      <KanbanView
        v-else-if="currentView === 'kanban'"
        :columns="kanbanColumns"
        :items="items"
        :loading="loading"
        :title-field="viewConfig.titleField"
        :desc-field="viewConfig.descField"
        :status-field="viewConfig.statusField"
        @card-click="handleItemClick"
        @card-move="handleCardMove"
      >
        <template #card="{ item, column }">
          <slot name="kanbanCard" :item="item" :column="column" />
        </template>
      </KanbanView>

      <!-- TimelineView -->
      <TimelineView
        v-else-if="currentView === 'timeline'"
        :items="sortedByTime"
        :loading="loading"
        :empty-text="emptyText"
        :title-field="viewConfig.titleField"
        :desc-field="viewConfig.descField"
        :time-field="viewConfig.timeField"
        :type-field="viewConfig.typeField"
        @item-click="handleItemClick"
      >
        <template #item="{ item, index }">
          <slot name="timelineItem" :item="item" :index="index" />
        </template>
      </TimelineView>

      <!-- ChartView -->
      <ChartView
        v-else-if="currentView === 'chart'"
        :data="chartData"
        :loading="loading"
        :height="viewConfig.height || '400px'"
        :default-type="viewConfig.chartType || 'bar'"
        @chart-click="handleChartClick"
      />
    </div>

    <!-- 统一的详情显示（抽屉/对话框） -->
    <NoteDetailDisplay
      ref="detailDisplay"
      :title="currentItem?.title || currentItem?.name"
      :load-detail="loadDetailFn"
    >
      <template #default="{ detail, loading }">
        <slot name="detail" :item="detail" :loading="loading" />
      </template>
    </NoteDetailDisplay>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ViewSwitcher, { type ViewType } from './ViewSwitcher.vue'
import GridView from './GridView.vue'
import ListView from './ListView.vue'
import TableView from './TableView.vue'
import TreeView from './TreeView.vue'
import KanbanView from './KanbanView.vue'
import TimelineView from './TimelineView.vue'
import ChartView from './ChartView.vue'
import OpenModeSelector from '@/components/OpenModeSelector.vue'
import NoteDetailDisplay from '@/components/NoteDetailDisplay.vue'
import { useNoteViewSettings, type OpenMode } from '@/composables/useNoteViewSettings'
import type { TableColumnConfig, TableFilterConfig } from './TableView.vue'
import type { KanbanColumn } from './KanbanView.vue'
import type { ChartDataItem } from './ChartView.vue'
import type { TreeNodeConfig } from './TreeView.vue'
import type { FilterConfig as GridFilterConfig } from './GridView.vue'

// 视图配置
export interface NoteViewConfig {
  // 基础字段映射
  titleField?: string
  descField?: string
  tagField?: string
  statusField?: string
  timeField?: string
  dateField?: string
  typeField?: string
  
  // GridView 配置
  cardType?: 'book' | 'video' | 'anime' | 'movie' | 'game' | 'quote' | 'github'
  gridColumns?: string
  
  // TableView 配置
  columns?: TableColumnConfig[]
  
  // TreeView 配置
  treeData?: TreeNodeConfig[]
  
  // KanbanView 配置
  kanbanColumns?: KanbanColumn[]
  
  // ChartView 配置
  chartType?: 'bar' | 'line' | 'pie' | 'radar'
  chartGroupBy?: string
  
  // 公共配置
  filters?: TableFilterConfig[]
  height?: string
}

interface Props {
  title: string
  items: any[]
  loading?: boolean
  emptyText?: string
  showStats?: boolean
  availableViews?: ViewType[]
  defaultView?: ViewType
  viewConfig?: NoteViewConfig
  // 详情加载函数
  loadDetail?: (item: any) => Promise<any>
  // 跳转路由模板（支持 :id 等占位符）
  detailRoute?: string
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyText: '暂无数据',
  showStats: true,
  availableViews: () => ['grid', 'list', 'table'],
  defaultView: 'grid',
  viewConfig: () => ({}),
})

const emit = defineEmits<{
  itemClick: [item: any]
  cardMove: [item: any, from: any, to: any]
  chartClick: [params: any]
  viewChange: [view: ViewType]
}>()

const route = useRoute()
const router = useRouter()
const { settings, setOpenMode } = useNoteViewSettings()
const currentView = ref<ViewType>(props.defaultView)
const openMode = computed({
  get: () => settings.value.openMode,
  set: (value: OpenMode) => setOpenMode(value),
})
const detailDisplay = ref<InstanceType<typeof NoteDetailDisplay> | null>(null)
const currentItem = ref<any>(null)

// 详情加载函数
const loadDetailFn = computed(() => props.loadDetail)

// 总数
const total = computed(() => props.items.length)

// GridView 过滤器（只支持 radio 和 select 类型）
const gridFilters = computed<GridFilterConfig[]>(() => {
  if (!props.viewConfig.filters) return []
  return props.viewConfig.filters
    .filter(f => f.type === 'radio' || f.type === 'select')
    .map(f => ({
      key: f.key,
      type: f.type as 'radio' | 'select',
      options: f.options || [],
      placeholder: f.placeholder,
      width: f.width,
      defaultValue: f.defaultValue,
    }))
})

// 表格列配置
const tableColumns = computed<TableColumnConfig[]>(() => {
  if (props.viewConfig.columns) {
    return props.viewConfig.columns
  }
  
  // 默认列配置
  return [
    { prop: props.viewConfig.titleField || 'title', label: '标题', minWidth: 200, sortable: true },
    { prop: props.viewConfig.tagField || 'category', label: '分类', width: 120 },
    { prop: props.viewConfig.statusField || 'status', label: '状态', width: 100, type: 'tag' },
    { prop: props.viewConfig.timeField || 'updatedAt', label: '更新时间', width: 160, type: 'date' },
  ]
})

// 树形数据
const treeData = computed(() => {
  return props.viewConfig.treeData || []
})

// 看板列
const kanbanColumns = computed<KanbanColumn[]>(() => {
  if (props.viewConfig.kanbanColumns) {
    return props.viewConfig.kanbanColumns
  }
  
  // 默认看板列
  return [
    { id: 'pending', title: '待处理', statusValue: 'pending' },
    { id: 'doing', title: '进行中', statusValue: 'doing' },
    { id: 'done', title: '已完成', statusValue: 'done' },
  ]
})

// 按时间排序的数据
const sortedByTime = computed(() => {
  const timeField = props.viewConfig.timeField || 'updatedAt'
  return [...props.items].sort((a, b) => {
    const timeA = new Date(a[timeField] || 0).getTime()
    const timeB = new Date(b[timeField] || 0).getTime()
    return timeB - timeA // 降序，最新的在前
  })
})

// 图表数据
const chartData = computed<ChartDataItem[]>(() => {
  const groupBy = props.viewConfig.chartGroupBy || 'category'
  
  // 分组统计
  const groups: Record<string, number> = {}
  props.items.forEach(item => {
    const key = item[groupBy] || '未知'
    groups[key] = (groups[key] || 0) + 1
  })
  
  return Object.entries(groups).map(([name, value]) => ({
    name,
    value,
    percentage: ((value / props.items.length) * 100).toFixed(1),
  }))
})

// 处理点击
function handleItemClick(item: any) {
  currentItem.value = item
  emit('itemClick', item)

  // 树状视图使用专用交互，只显示右侧内容，不走打开模式
  if (currentView.value === 'tree') {
    return
  }

  // 如果有 detailRoute 且是页面模式，跳转到详情页
  if (openMode.value === 'page' && props.detailRoute) {
    const path = props.detailRoute
      .replace(':id', item.id || item.path || '')
      .replace(':type', item.type || '')
      .replace(':title', encodeURIComponent(item.title || item.name || ''))
    router.push(path)
    return
  }

  // 抽屉或对话框模式
  if (detailDisplay.value && openMode.value !== 'page') {
    detailDisplay.value.open(item)
  }
}

// 处理看板移动
function handleCardMove(item: any, from: any, to: any) {
  emit('cardMove', item, from, to)
}

// 处理图表点击
function handleChartClick(params: any) {
  emit('chartClick', params)
}

// 初始化视图
function initView() {
  // 从 URL 参数读取
  const viewParam = route.query.view as ViewType
  if (viewParam && props.availableViews.includes(viewParam)) {
    currentView.value = viewParam
    return
  }
  
  // 从 localStorage 读取
  const savedView = localStorage.getItem(`view-${route.path}`)
  if (savedView && props.availableViews.includes(savedView as ViewType)) {
    currentView.value = savedView as ViewType
    return
  }
  
  currentView.value = props.defaultView
}

// 监听视图变化
watch(currentView, (newView) => {
  localStorage.setItem(`view-${route.path}`, newView)
  emit('viewChange', newView)
})

onMounted(initView)

// 暴露方法
defineExpose({
  currentView,
  tableColumns,
  kanbanColumns,
  chartData,
  sortedByTime,
})
</script>

<style scoped>
.note-view-container {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.page-stats {
  color: #909399;
  font-size: 14px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.view-content {
  width: 100%;
}
</style>
