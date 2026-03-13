<template>
  <div class="kanban-view">
    <!-- 工具栏 -->
    <div v-if="showToolbar" class="kanban-toolbar">
      <el-input
        v-if="searchable"
        v-model="searchText"
        placeholder="搜索..."
        clearable
        size="small"
        style="width: 200px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <slot name="toolbar" />
    </div>

    <!-- 加载状态 -->
    <el-skeleton v-if="loading" :rows="5" animated />

    <!-- 看板列 -->
    <div v-else class="kanban-columns" :style="{ height: height }">
      <div
        v-for="column in columns"
        :key="column.id"
        class="kanban-column"
        :style="{ width: columnWidth }"
      >
        <!-- 列头 -->
        <div class="column-header" :style="{ background: column.color || defaultColors[columns.indexOf(column) % defaultColors.length] }">
          <span class="column-title">{{ column.title }}</span>
          <el-badge :value="getColumnItems(column).length" :max="99" />
        </div>

        <!-- 列内容 -->
        <div
          class="column-body"
          @dragover.prevent
          @drop="handleDrop($event, column)"
        >
          <!-- 空状态 -->
          <div v-if="getColumnItems(column).length === 0" class="empty-slot">
            拖拽卡片到此处
          </div>

          <!-- 卡片列表 -->
          <div
            v-for="item in getColumnItems(column)"
            :key="getItemKey(item)"
            class="kanban-card"
            draggable="true"
            @dragstart="handleDragStart($event, item, column)"
            @dragend="handleDragEnd"
            @click="handleCardClick(item)"
          >
            <slot name="card" :item="item" :column="column">
              <!-- 默认卡片内容 -->
              <div class="card-header">
                <span class="card-title">{{ getItemTitle(item) }}</span>
                <el-tag v-if="getItemTag(item)" size="small">
                  {{ getItemTag(item) }}
                </el-tag>
              </div>
              <div v-if="getItemDescription(item)" class="card-desc">
                {{ getItemDescription(item) }}
              </div>
              <div v-if="showMeta" class="card-meta">
                <span v-if="getItemDate(item)">
                  <el-icon><Clock /></el-icon>
                  {{ getItemDate(item) }}
                </span>
              </div>
            </slot>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Search, Clock } from '@element-plus/icons-vue'

// 看板列配置
export interface KanbanColumn {
  id: string
  title: string
  color?: string
  statusField?: string
  statusValue?: string
}

interface Props {
  columns: KanbanColumn[]
  items: any[]
  loading?: boolean
  height?: string
  columnWidth?: string
  
  // 字段映射
  itemKey?: string
  titleField?: string
  descField?: string
  tagField?: string
  dateField?: string
  statusField?: string
  
  // 功能开关
  searchable?: boolean
  showToolbar?: boolean
  showMeta?: boolean
  draggable?: boolean
  
  // 颜色
  defaultColors?: string[]
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  height: 'calc(100vh - 220px)',
  columnWidth: '300px',
  itemKey: 'id',
  titleField: 'title',
  descField: 'description',
  tagField: 'category',
  dateField: 'updatedAt',
  statusField: 'status',
  searchable: true,
  showToolbar: true,
  showMeta: true,
  draggable: true,
  defaultColors: () => [
    'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
  ],
})

const emit = defineEmits<{
  cardClick: [item: any]
  cardMove: [item: any, fromColumn: KanbanColumn, toColumn: KanbanColumn]
}>()

const searchText = ref('')
const draggedItem = ref<any>(null)
const draggedFromColumn = ref<KanbanColumn | null>(null)

// 获取列中的项目
function getColumnItems(column: KanbanColumn) {
  return props.items.filter(item => {
    // 状态匹配
    const status = item[props.statusField] || item[column.statusField || 'status']
    const matchesStatus = status === (column.statusValue || column.id)
    
    // 搜索过滤
    const matchesSearch = !searchText.value ||
      String(item[props.titleField] || '').toLowerCase().includes(searchText.value.toLowerCase()) ||
      String(item[props.descField] || '').toLowerCase().includes(searchText.value.toLowerCase())
    
    return matchesStatus && matchesSearch
  })
}

// 获取项目 key
function getItemKey(item: any): string {
  return item[props.itemKey] || JSON.stringify(item)
}

// 获取项目标题
function getItemTitle(item: any): string {
  return item[props.titleField] || item.name || '未命名'
}

// 获取项目描述
function getItemDescription(item: any): string {
  const desc = item[props.descField]
  if (!desc) return ''
  return desc.length > 50 ? desc.slice(0, 50) + '...' : desc
}

// 获取项目标签
function getItemTag(item: any): string | null {
  return item[props.tagField] || null
}

// 获取项目日期
function getItemDate(item: any): string | null {
  const date = item[props.dateField]
  if (!date) return null
  return new Date(date).toLocaleDateString()
}

// 拖拽开始
function handleDragStart(event: DragEvent, item: any, column: KanbanColumn) {
  if (!props.draggable) return
  
  draggedItem.value = item
  draggedFromColumn.value = column
  ;(event.target as HTMLElement).classList.add('dragging')
}

// 拖拽结束
function handleDragEnd(event: DragEvent) {
  ;(event.target as HTMLElement).classList.remove('dragging')
  draggedItem.value = null
  draggedFromColumn.value = null
}

// 放置
function handleDrop(_event: DragEvent, toColumn: KanbanColumn) {
  if (!props.draggable || !draggedItem.value || !draggedFromColumn.value) return
  
  // 不允许同一列
  if (draggedFromColumn.value.id === toColumn.id) return
  
  emit('cardMove', draggedItem.value, draggedFromColumn.value, toColumn)
}

// 点击卡片
function handleCardClick(item: any) {
  emit('cardClick', item)
}

// 暴露方法
defineExpose({
  searchText,
  getColumnItems,
})
</script>

<style scoped>
.kanban-view {
  width: 100%;
}

.kanban-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.kanban-columns {
  display: flex;
  gap: 16px;
  overflow-x: auto;
  padding-bottom: 8px;
}

.kanban-column {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
  border-radius: 8px;
  overflow: hidden;
}

.column-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  color: white;
}

.column-title {
  font-size: 14px;
  font-weight: 600;
}

.column-body {
  flex: 1;
  padding: 12px;
  overflow-y: auto;
  min-height: 100px;
}

.empty-slot {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 80px;
  color: #909399;
  font-size: 12px;
  border: 1px dashed #dcdfe6;
  border-radius: 4px;
}

.kanban-card {
  background: white;
  border-radius: 6px;
  padding: 12px;
  margin-bottom: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.2s;
}

.kanban-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.kanban-card.dragging {
  opacity: 0.5;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.card-title {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-desc {
  font-size: 12px;
  color: #606266;
  line-height: 1.5;
  margin-bottom: 8px;
}

.card-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #909399;
}

.card-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
