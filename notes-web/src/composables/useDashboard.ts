/**
 * Dashboard 配置管理 composable
 */

import { ref, computed, type Component } from 'vue'
import {
  Reading,
  VideoPlay,
  Collection,
  Tools,
  QuestionFilled,
  Star,
} from '@element-plus/icons-vue'

// Widget 类型定义
export interface StatsWidgetConfig {
  type: 'stats'
  icon: Component
  value: number
  label: string
  gradient: string
  route: string
}

export interface ProgressWidgetItem {
  label: string
  percentage: number
  count: string
  status?: '' | 'success' | 'warning' | 'exception'
}

export interface ProgressWidgetConfig {
  type: 'progress'
  title: string
  items: ProgressWidgetItem[]
  showMore?: boolean
  moreRoute?: string
}

export interface QuickEntryItem {
  icon: Component
  label: string
  color: string
  route: string
}

export interface QuickEntryWidgetConfig {
  type: 'quickEntry'
  title: string
  items: QuickEntryItem[]
}

export interface QuickLinkWidgetConfig {
  type: 'quickLink'
  icon: Component
  title: string
  description: string
  color: string
  route: string
}

// GridView 筛选器配置
export interface GridFilterOption {
  label: string
  value: string
}

export interface GridFilterConfig {
  key: string
  type: 'radio' | 'select'
  options: GridFilterOption[]
  placeholder?: string
  width?: string
  defaultValue?: string
}

// GridView widget 配置
export interface GridViewWidgetConfig {
  type: 'gridView'
  items: any[]
  loading?: boolean
  emptyText?: string
  cardType?: 'book' | 'video'
  filters?: GridFilterConfig[]
  gridColumns?: string
  gridGap?: string
  itemKey?: string
  fullSpan?: boolean  // 是否占满整行
}

// ListView 筛选器配置
export interface ListFilterOption {
  label: string
  value: string
}

export interface ListFilterConfig {
  key: string
  type: 'search' | 'radio' | 'select'
  options?: ListFilterOption[]
  placeholder?: string
  width?: string
  defaultValue?: string
  searchFields?: string[]
}

// ListView widget 配置
export interface ListViewWidgetConfig {
  type: 'listView'
  items: any[]
  loading?: boolean
  emptyText?: string
  filters?: ListFilterConfig[]
  itemKey?: string
  titleField?: string
  descField?: string
  tagField?: string
}

// TableView 筛选器配置
export interface TableFilterOption {
  label: string
  value: string
}

export interface TableFilterConfig {
  key: string
  type: 'search' | 'radio' | 'select'
  options?: TableFilterOption[]
  placeholder?: string
  width?: string
  defaultValue?: string
  searchFields?: string[]
}

// TableView 列配置
export interface TableColumnConfig {
  prop: string
  label: string
  width?: number | string
  minWidth?: number | string
  sortable?: boolean | 'custom'
  align?: 'left' | 'center' | 'right'
  fixed?: boolean | 'left' | 'right'
  type?: 'tag' | 'link' | 'progress' | 'date' | 'image'
  slot?: boolean
  tagSize?: 'large' | 'default' | 'small'
  getTagType?: (row: any) => 'success' | 'warning' | 'danger' | 'info' | 'primary'
  linkType?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'default'
  onClick?: (row: any) => void
  getProgressStatus?: (row: any) => '' | 'success' | 'warning' | 'exception'
  dateFormat?: string
  imageWidth?: string
  imageHeight?: string
  formatter?: (value: any, row: any) => string
}

// TableView widget 配置
export interface TableViewWidgetConfig {
  type: 'tableView'
  columns: TableColumnConfig[]
  data: any[]
  loading?: boolean
  rowKey?: string
  height?: number | string
  maxHeight?: number | string
  stripe?: boolean
  border?: boolean
  highlightCurrentRow?: boolean
  defaultSort?: { prop: string; order: 'ascending' | 'descending' }
  filters?: TableFilterConfig[]
  pagination?: boolean
  pageSize?: number
  pageSizes?: number[]
  paginationLayout?: string
  selectable?: boolean
  showIndex?: boolean
  operationWidth?: number | string
  operationFixed?: boolean | 'left' | 'right'
}

// TreeView 节点配置
export interface TreeNodeConfig {
  id?: string
  label: string
  children?: TreeNodeConfig[]
  isLeaf?: boolean
  isCategory?: boolean
  [key: string]: any
}

// TreeView widget 配置
export interface TreeViewWidgetConfig {
  type: 'treeView'
  data: TreeNodeConfig[]
  loading?: boolean
  height?: string
  leftSpan?: number
  nodeKey?: string
  titleField?: string
  descField?: string
  tagField?: string
}

// KanbanView 列配置
export interface KanbanColumn {
  id: string
  title: string
  color?: string
  statusField?: string
  statusValue?: string
}

// KanbanView widget 配置
export interface KanbanViewWidgetConfig {
  type: 'kanbanView'
  columns: KanbanColumn[]
  items: any[]
  loading?: boolean
  height?: string
  columnWidth?: string
  titleField?: string
  descField?: string
  tagField?: string
  statusField?: string
}

// TimelineView widget 配置
export interface TimelineViewWidgetConfig {
  type: 'timelineView'
  items: any[]
  loading?: boolean
  emptyText?: string
  titleField?: string
  descField?: string
  tagField?: string
  timeField?: string
  typeField?: string
  viewMode?: 'vertical' | 'horizontal'
}

// ChartView 数据项
export interface ChartDataItem {
  name: string
  value: number
  percentage?: string
}

// ChartView widget 配置
export interface ChartViewWidgetConfig {
  type: 'chartView'
  data?: ChartDataItem[]
  series?: { name: string; data: number[] }[]
  categories?: string[]
  loading?: boolean
  height?: string
  title?: string
  defaultType?: 'bar' | 'line' | 'pie' | 'radar'
  showTypeSelector?: boolean
  showLegend?: boolean
}

// NoteViewContainer 配置
export interface NoteViewContainerConfig {
  type: 'noteView'
  title: string
  items: any[]
  loading?: boolean
  emptyText?: string
  showStats?: boolean
  availableViews?: ('grid' | 'list' | 'table' | 'tree' | 'chart' | 'kanban' | 'timeline')[]
  defaultView?: 'grid' | 'list' | 'table' | 'tree' | 'chart' | 'kanban' | 'timeline'
  viewConfig?: {
    titleField?: string
    descField?: string
    tagField?: string
    statusField?: string
    timeField?: string
    typeField?: string
    cardType?: 'book' | 'video'
    gridColumns?: string
    columns?: TableColumnConfig[]
    kanbanColumns?: KanbanColumn[]
    chartType?: 'bar' | 'line' | 'pie' | 'radar'
    chartGroupBy?: string
    filters?: TableFilterConfig[]
    height?: string
  }
}

export type WidgetConfig =
  | StatsWidgetConfig
  | ProgressWidgetConfig
  | QuickEntryWidgetConfig
  | QuickLinkWidgetConfig
  | GridViewWidgetConfig
  | ListViewWidgetConfig
  | TableViewWidgetConfig
  | TreeViewWidgetConfig
  | KanbanViewWidgetConfig
  | TimelineViewWidgetConfig
  | ChartViewWidgetConfig
  | NoteViewContainerConfig

// 行配置
export interface RowConfig {
  gutter?: number
  widgets: WidgetConfig[]
}

// Dashboard 配置
export interface DashboardConfig {
  title: string
  loading?: boolean
  rows: RowConfig[]
}

// 预设渐变色
export const GRADIENTS = {
  blue: 'linear-gradient(135deg, #409eff 0%, #66b1ff 100%)',
  green: 'linear-gradient(135deg, #67c23a 0%, #85ce61 100%)',
  gray: 'linear-gradient(135deg, #909399 0%, #b1b3b8 100%)',
  purple: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
  orange: 'linear-gradient(135deg, #e6a23c 0%, #f0c78a 100%)',
  red: 'linear-gradient(135deg, #f56c6c 0%, #f89898 100%)',
}

// 图标映射
export const ICONS = {
  reading: Reading,
  videoPlay: VideoPlay,
  collection: Collection,
  tools: Tools,
  questionFilled: QuestionFilled,
  star: Star,
}

export function useDashboard() {
  const config = ref<DashboardConfig>({
    title: '',
    loading: false,
    rows: [],
  })

  function setConfig(newConfig: DashboardConfig) {
    config.value = newConfig
  }

  function setLoading(loading: boolean) {
    config.value.loading = loading
  }

  const title = computed(() => config.value.title)
  const loading = computed(() => config.value.loading)
  const rows = computed(() => config.value.rows)

  return {
    config,
    title,
    loading,
    rows,
    setConfig,
    setLoading,
    GRADIENTS,
    ICONS,
  }
}
