<template>
  <div class="timeline-view">
    <!-- 工具栏 -->
    <div v-if="showToolbar" class="timeline-toolbar">
      <el-radio-group v-model="localViewMode" size="small">
        <el-radio-button value="vertical">垂直</el-radio-button>
        <el-radio-button value="horizontal">水平</el-radio-button>
      </el-radio-group>
      
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

    <!-- 空状态 -->
    <el-empty v-else-if="filteredItems.length === 0" :description="emptyText" />

    <!-- 时间线 -->
    <div v-else :class="['timeline-container', localViewMode]">
      <el-timeline :reverse="reverse">
        <el-timeline-item
          v-for="(item, index) in filteredItems"
          :key="getItemKey(item)"
          :timestamp="getItemTime(item)"
          :placement="placement"
          :type="getItemType(item)"
          :size="size"
          :hollow="hollow"
          :color="getItemColor(item)"
        >
          <template #dot>
            <slot name="dot" :item="item" :index="index">
              <el-icon v-if="getItemIcon(item)">
                <component :is="getItemIcon(item)" />
              </el-icon>
            </slot>
          </template>
          
          <div class="timeline-card" @click="handleItemClick(item)">
            <slot name="item" :item="item" :index="index">
              <div class="card-header">
                <span class="card-title">{{ getItemTitle(item) }}</span>
                <el-tag v-if="getItemTag(item)" size="small" :type="getItemTagType(item)">
                  {{ getItemTag(item) }}
                </el-tag>
              </div>
              <div v-if="getItemDescription(item)" class="card-desc">
                {{ getItemDescription(item) }}
              </div>
              <div v-if="showMeta" class="card-meta">
                <span v-if="getItemDate(item)">
                  <el-icon><Calendar /></el-icon>
                  {{ getItemDate(item) }}
                </span>
              </div>
            </slot>
          </div>
        </el-timeline-item>
      </el-timeline>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Search, Calendar, Document, Star, Clock, Location } from '@element-plus/icons-vue'

interface Props {
  items: any[]
  loading?: boolean
  emptyText?: string
  
  // 字段映射
  itemKey?: string
  titleField?: string
  descField?: string
  tagField?: string
  timeField?: string
  dateField?: string
  
  // 显示配置
  viewMode?: 'vertical' | 'horizontal'
  placement?: 'top' | 'bottom'
  size?: 'normal' | 'large'
  hollow?: boolean
  reverse?: boolean
  
  // 功能开关
  showToolbar?: boolean
  searchable?: boolean
  showMeta?: boolean
  
  // 类型/颜色映射
  typeField?: string
  typeColors?: Record<string, string>
  typeIcons?: Record<string, any>
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyText: '暂无数据',
  itemKey: 'id',
  titleField: 'title',
  descField: 'description',
  tagField: 'category',
  timeField: 'time',
  dateField: 'date',
  viewMode: 'vertical',
  placement: 'top',
  size: 'normal',
  hollow: false,
  reverse: false,
  showToolbar: true,
  searchable: true,
  showMeta: true,
  typeField: 'type',
  typeColors: () => ({
    book: '#409eff',
    video: '#67c23a',
    knowledge: '#e6a23c',
    skill: '#f56c6c',
    problem: '#909399',
  }),
  typeIcons: () => ({
    book: Document,
    video: Star,
    knowledge: Star,
    skill: Clock,
    problem: Location,
  }),
})

const emit = defineEmits<{
  itemClick: [item: any]
}>()

const searchText = ref('')
const localViewMode = ref<'vertical' | 'horizontal'>(props.viewMode)

// 过滤后的项目
const filteredItems = computed(() => {
  if (!searchText.value) return props.items
  
  const keyword = searchText.value.toLowerCase()
  return props.items.filter(item => {
    const title = String(item[props.titleField] || '').toLowerCase()
    const desc = String(item[props.descField] || '').toLowerCase()
    return title.includes(keyword) || desc.includes(keyword)
  })
})

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
  return desc.length > 100 ? desc.slice(0, 100) + '...' : desc
}

// 获取项目标签
function getItemTag(item: any): string | null {
  return item[props.tagField] || null
}

// 获取标签类型
function getItemTagType(item: any): 'success' | 'warning' | 'danger' | 'info' | 'primary' {
  const tag = item[props.tagField]
  const typeMap: Record<string, 'success' | 'warning' | 'danger' | 'info' | 'primary'> = {
    done: 'success',
    pending: 'warning',
    error: 'danger',
  }
  return typeMap[tag] || 'primary'
}

// 获取时间
function getItemTime(item: any): string {
  const time = item[props.timeField] || item[props.dateField] || item.createdAt || item.updatedAt
  if (!time) return ''
  return new Date(time).toLocaleString()
}

// 获取日期
function getItemDate(item: any): string | null {
  const date = item[props.dateField] || item[props.timeField] || item.createdAt || item.updatedAt
  if (!date) return null
  return new Date(date).toLocaleDateString()
}

// 获取类型
function getItemType(item: any): 'primary' | 'success' | 'warning' | 'danger' | 'info' {
  const type = item[props.typeField] || 'primary'
  const typeMap: Record<string, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
    primary: 'primary',
    success: 'success',
    warning: 'warning',
    danger: 'danger',
    info: 'info',
    done: 'success',
    pending: 'warning',
    error: 'danger',
  }
  return typeMap[type] || 'primary'
}

// 获取颜色
function getItemColor(item: any): string | undefined {
  const type = item[props.typeField]
  return props.typeColors[type]
}

// 获取图标
function getItemIcon(item: any): any {
  const type = item[props.typeField]
  return props.typeIcons[type] || Document
}

// 点击事件
function handleItemClick(item: any) {
  emit('itemClick', item)
}

// 暴露方法
defineExpose({
  searchText,
  filteredItems,
})
</script>

<style scoped>
.timeline-view {
  width: 100%;
}

.timeline-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.timeline-container {
  padding: 16px;
}

.timeline-container.horizontal {
  overflow-x: auto;
}

.timeline-card {
  background: white;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.2s;
}

.timeline-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transform: translateX(4px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
}

.card-title {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
  flex: 1;
}

.card-desc {
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
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
