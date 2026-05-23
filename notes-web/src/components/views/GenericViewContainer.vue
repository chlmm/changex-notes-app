<template>
  <div class="generic-view-container">
    <!-- 头部 -->
    <div class="view-header">
      <div class="header-left">
        <h2>
          <el-icon v-if="type?.Icon" :size="22"><component :is="type.Icon" /></el-icon>
          {{ type?.Label || typeId }}
        </h2>
        <span class="header-count">{{ flatNotes.length }} 条</span>
      </div>
      <div class="header-right">
        <ViewSwitcher v-model="currentView" :available-views="availableViews" />
      </div>
    </div>

    <!-- 加载中 -->
    <el-skeleton v-if="loading" :rows="8" animated />

    <!-- 表格视图 -->
    <TableView
      v-else-if="currentView === 'table'"
      :columns="tableColumns"
      :data="flatNotes"
      :pagination="true"
      :page-size="20"
      @row-click="onItemClick"
    />

    <!-- 网格视图（卡片） -->
    <GridView
      v-else-if="currentView === 'grid'"
      :items="flatNotes"
      @item-click="onItemClick"
    >
      <template #item="{ item }">
        <GenericNoteCard
          :note="(item as FlatNote)._raw"
          :card-fields="gridCardFields"
        />
      </template>
    </GridView>

    <!-- 列表视图 -->
    <ListView
      v-else-if="currentView === 'list'"
      :items="flatNotes"
      :title-field="'title'"
      :desc-field="listDescField"
      :tag-field="listTagField"
      @item-click="onItemClick"
    />

    <!-- 看板视图 -->
    <KanbanView
      v-else-if="currentView === 'kanban'"
      :columns="kanbanColumns"
      :items="flatNotes"
      :status-field="kanbanGroupField"
      :title-field="'title'"
      :tag-field="kanbanGroupField"
      @card-click="onItemClick"
    />

    <!-- 时间线视图 -->
    <TimelineView
      v-else-if="currentView === 'timeline'"
      :items="flatNotes"
      :time-field="timelineField"
      :title-field="'title'"
      :desc-field="listDescField"
      :type-field="'typeId'"
      @item-click="onItemClick"
    />

    <!-- 图表视图 -->
    <ChartView
      v-else-if="currentView === 'chart'"
      :data="chartData"
      :default-type="'bar'"
      :title="type?.Label + ' 分布统计'"
      @chart-click="() => {}"
    />

    <!-- 树形视图 -->
    <TreeView
      v-else-if="currentView === 'tree'"
      :data="treeData"
      :title-field="'title'"
      :desc-field="listDescField"
      @node-click="onNodeClick"
    />

    <!-- 空状态 -->
    <el-empty v-else-if="!flatNotes.length" description="暂无数据" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import ViewSwitcher, { type ViewType } from '@/components/dashboard/widgets/ViewSwitcher.vue'
import TableView, { type TableColumnConfig } from '@/components/dashboard/widgets/TableView.vue'
import GridView from '@/components/dashboard/widgets/GridView.vue'
import ListView from '@/components/dashboard/widgets/ListView.vue'
import KanbanView, { type KanbanColumn } from '@/components/dashboard/widgets/KanbanView.vue'
import TimelineView from '@/components/dashboard/widgets/TimelineView.vue'
import ChartView, { type ChartDataItem } from '@/components/dashboard/widgets/ChartView.vue'
import TreeView, { type TreeNodeConfig } from '@/components/dashboard/widgets/TreeView.vue'
import GenericNoteCard from './GenericNoteCard.vue'
import { useSchema } from '@/composables/useSchema'
import type { GenericNote } from '@/types/note'
import type { TypeDef } from '@/types/schema'

// ========================
// Flat adapter: GenericNote → old view compatible flat object
// ========================
interface FlatNote {
  _id: string
  _typeId: string
  _path: string
  _updatedAt?: string
  _raw: GenericNote
  [key: string]: unknown
}

// 从扁平化字段中自动检测"标题"字段
// 多条目笔记的标题通常隐藏在 source.title / source.name 中
function detectTitle(fields: Record<string, unknown>): string {
  // 优先直接用 title
  if (fields['title']) return String(fields['title'])
  // 多条目类型：source.title
  if (fields['source.title']) return String(fields['source.title'])
  // GitHub 类型：source.name
  if (fields['source.name']) return String(fields['source.name'])
  // quote 类型：取前 30 字
  if (fields['quote']) {
    const q = String(fields['quote'])
    return q.length > 30 ? q.slice(0, 30) + '...' : q
  }
  // 兜底
  return ''
}

function flattenNote(note: GenericNote): FlatNote {
  const flat: FlatNote = {
    _id: note.id,
    _typeId: note.typeId,
    _path: note.path,
    _updatedAt: note.updatedAt,
    _raw: note,
    // 注入 title: 让 list/kanban/tree 等视图能找到名称
    title: detectTitle(note.fields) || note.id,
  }
  // Merge all fields to top level (safe: field keys don't start with _)
  for (const [key, value] of Object.entries(note.fields)) {
    flat[key] = value
  }
  return flat
}

// ========================
// Props & Emits
// ========================
const props = defineProps<{
  typeId: string
  notes: GenericNote[]
  loading?: boolean
}>()

const emit = defineEmits<{
  itemClick: [note: GenericNote]
  nodeClick: [node: TreeNodeConfig]
}>()

const route = useRoute()
const { getType } = useSchema()

// ========================
// State
// ========================
const currentView = ref<ViewType>('table')

// ========================
// Computed: type & flat notes
// ========================
const type = computed<TypeDef | undefined>(() => getType(props.typeId))

const flatNotes = computed<FlatNote[]>(() =>
  props.notes.map(flattenNote)
)

// ========================
// Available views (explicit config > auto-detect)
// ========================
const availableViews = computed<ViewType[]>(() => {
  if (!type.value) return ['table', 'grid', 'list']

  // Explicit view config takes priority
  const explicit = type.value.Views?.Enabled
  if (explicit && explicit.length > 0) {
    return explicit.filter(v =>
      (['table', 'grid', 'list', 'kanban', 'timeline', 'chart', 'tree'] as string[]).includes(v)
    ) as ViewType[]
  }

  // Auto-detect from field hints
  const views: ViewType[] = ['table', 'grid', 'list']
  const hints = type.value.FieldHints || {}
  const hintValues = Object.values(hints)

  if (hintValues.some(h => h.UI === 'select')) views.push('kanban')
  if (hintValues.some(h => h.UI === 'date') ||
      Object.keys(hints).some(k => /date|time|created|updated|timestamp/i.test(k)))
    views.push('timeline')
  views.push('chart')
  if (hintValues.some(h => h.UI === 'select') ||
      Object.keys(hints).some(k => /category|parent|group/i.test(k)))
    views.push('tree')

  return views
})

// ========================
// Table config
// ========================
const tableColumns = computed<TableColumnConfig[]>(() => {
  if (!type.value) return []

  const hints = type.value.FieldHints || {}
  const columns = type.value.ListColumns.map((key): TableColumnConfig => {
    const hint = hints[key]
    const ui = hint?.UI || 'text'

    // Map schema UI types to TableView column types
    let colType: TableColumnConfig['type'] = undefined
    if (ui === 'date') colType = 'date'
    if (ui === 'tags') colType = 'tag'
    if (ui === 'select') colType = 'tag'

    // Pick a sensible label
    const parts = key.split('.')
    const label = hint?.Label || parts[parts.length - 1].replace(/_/g, ' ')

    return {
      prop: key,       // becomes top-level after flatten
      label,
      minWidth: key.includes('.') ? 120 : 80,
      sortable: true,
      type: colType,
    }
  })

  return columns
})

// ========================
// Grid card config
// ========================
const gridCardFields = computed(() => {
  if (!type.value) return []
  const cols = type.value.ListColumns
  const hints = type.value.FieldHints || {}
  return cols.slice(0, 4).map(key => {
    const parts = key.split('.')
    return {
      key,
      label: hints[key]?.Label || parts[parts.length - 1].replace(/_/g, ' '),
      hint: hints[key] || null,
    }
  })
})

// ========================
// List view config
// ========================
const listDescField = computed(() => {
  if (!type.value) return 'description'
  const keys = Object.keys(type.value.FieldHints || {})
  return keys.find(k => /desc|content|summary|note/i.test(k)) || keys[1] || 'description'
})

const listTagField = computed(() => {
  if (!type.value) return ''
  const hints = type.value.FieldHints || {}
  const selectField = Object.entries(hints).find(([, h]) => h.UI === 'select')
  return selectField?.[0] || ''
})

// ========================
// Kanban config
// ========================
const kanbanGroupField = computed(() => {
  if (!type.value) return 'status'
  if (type.value.Views?.KanbanField) return type.value.Views.KanbanField
  const hints = type.value.FieldHints || {}
  const selectField = Object.entries(hints).find(([, h]) => h.UI === 'select')
  return selectField?.[0] || 'status'
})

const kanbanColumns = computed<KanbanColumn[]>(() => {
  if (!type.value) return []
  const hints = type.value.FieldHints || {}
  const fieldHint = hints[kanbanGroupField.value]
  const options = fieldHint?.Options || []

  if (options.length > 0) {
    return options.map(opt => ({
      id: opt,
      title: opt,
      statusValue: opt,
    }))
  }

  // Auto-detect from data
  const values = new Set<string>()
  for (const note of props.notes) {
    const v = note.fields[kanbanGroupField.value]
    if (v && typeof v === 'string') values.add(v)
  }

  return Array.from(values).map(v => ({
    id: v,
    title: v,
    statusValue: v,
  }))
})

// ========================
// Timeline config
// ========================
const timelineField = computed(() => {
  if (!type.value) return 'updatedAt'
  if (type.value.Views?.TimelineField) return type.value.Views.TimelineField
  const hints = type.value.FieldHints || {}
  const dateField = Object.entries(hints).find(([, h]) => h.UI === 'date')
  if (dateField) return dateField[0]
  const timeField = Object.keys(hints).find(k => /created|updated|date|time/i.test(k))
  return timeField || 'updatedAt'
})

// ========================
// Chart data (aggregate by first select or category field)
// ========================
const chartData = computed<ChartDataItem[]>(() => {
  if (!type.value || !props.notes.length) return []

  const hints = type.value.FieldHints || {}

  // Pick a group-by field: explicit config > select field > tags field > category
  const groupField =
    type.value.Views?.ChartGroupBy ||
    Object.entries(hints).find(([, h]) => h.UI === 'select' || h.UI === 'tags')?.[0] ||
    Object.keys(hints).find(k => /category|status|type/i.test(k)) ||
    'status'

  const groups: Record<string, number> = {}
  for (const note of props.notes) {
    const raw = note.fields[groupField]
    if (!raw) continue

    if (Array.isArray(raw)) {
      for (const item of raw) {
        const key = String(item)
        groups[key] = (groups[key] || 0) + 1
      }
    } else {
      const key = String(raw)
      groups[key] = (groups[key] || 0) + 1
    }
  }

  const total = props.notes.length
  return Object.entries(groups)
    .map(([name, value]) => ({
      name,
      value,
      percentage: total > 0 ? ((value / total) * 100).toFixed(1) : '0',
    }))
    .sort((a, b) => b.value - a.value)
})

// ========================
// Tree data (group by select/category field)
// ========================
const treeData = computed<TreeNodeConfig[]>(() => {
  if (!type.value || !props.notes.length) return []

  const hints = type.value.FieldHints || {}
  const groupField =
    Object.entries(hints).find(([, h]) => h.UI === 'select')?.[0] ||
    Object.keys(hints).find(k => /category|group|parent/i.test(k))

  if (!groupField) {
    // Flat tree: just list all
    return flatNotes.value.map(n => ({
      id: n._id,
      label: String(n.title || n._id),
      isLeaf: true,
      _raw: n._raw,
    }))
  }

  const groups: Record<string, TreeNodeConfig[]> = {}
  for (const note of flatNotes.value) {
    const key = String(note[groupField] || '其他')
    if (!groups[key]) groups[key] = []
    groups[key].push({
      id: note._id,
      label: String(note.title || note._id),
      isLeaf: true,
      _raw: note._raw,
    } as TreeNodeConfig)
  }

  return Object.entries(groups).map(([name, children]) => ({
    id: name,
    label: `${name} (${children.length})`,
    isCategory: true,
    children,
  }))
})

// ========================
// Handlers
// ========================
function onItemClick(item: FlatNote) {
  emit('itemClick', item._raw)
}

function onNodeClick(data: TreeNodeConfig) {
  if (data._raw) {
    emit('itemClick', data._raw as unknown as GenericNote)
  }
  emit('nodeClick', data)
}

// ========================
// View preference persistence
// ========================
function initView() {
  const viewParam = route.query.view as ViewType | undefined
  const available = availableViews.value

  if (viewParam && available.includes(viewParam)) {
    currentView.value = viewParam
    return
  }

  const saved = localStorage.getItem(`view-${props.typeId}`)
  if (saved && available.includes(saved as ViewType)) {
    currentView.value = saved as ViewType
    return
  }

  currentView.value = 'table'
}

watch(currentView, (v) => {
  localStorage.setItem(`view-${props.typeId}`, v)
})

onMounted(initView)
</script>

<style scoped>
.generic-view-container {
  padding: 16px;
  max-width: 1400px;
  margin: 0 auto;
}

.view-header {
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

.header-left h2 {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 20px;
  color: var(--el-text-color-primary);
}

.header-count {
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.view-content {
  width: 100%;
}
</style>
