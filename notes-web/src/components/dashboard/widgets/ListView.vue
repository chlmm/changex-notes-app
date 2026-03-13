<template>
  <div class="list-view">
    <!-- 筛选栏 -->
    <div v-if="showFilters" class="filter-bar">
      <template v-for="filter in filters" :key="filter.key">
        <!-- 搜索框 -->
        <el-input
          v-if="filter.type === 'search'"
          v-model="filterValues[filter.key]"
          :placeholder="filter.placeholder || '搜索...'"
          clearable
          :style="{ width: filter.width || '300px' }"
          @input="handleFilterChange"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>

        <!-- Radio 筛选 -->
        <el-radio-group
          v-else-if="filter.type === 'radio'"
          v-model="filterValues[filter.key]"
          @change="handleFilterChange"
        >
          <el-radio-button
            v-for="option in filter.options"
            :key="option.value"
            :value="option.value"
          >
            {{ option.label }}
          </el-radio-button>
        </el-radio-group>

        <!-- Select 筛选 -->
        <el-select
          v-else-if="filter.type === 'select'"
          v-model="filterValues[filter.key]"
          :placeholder="filter.placeholder"
          clearable
          :style="{ width: filter.width || '200px' }"
          @change="handleFilterChange"
        >
          <el-option
            v-for="option in filter.options"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          />
        </el-select>
      </template>
    </div>

    <!-- 加载状态 -->
    <el-skeleton v-if="loading" :rows="5" animated />

    <!-- 空状态 -->
    <el-empty v-else-if="filteredItems.length === 0" :description="emptyText" />

    <!-- 列表内容 -->
    <template v-else>
      <el-card
        v-for="item in filteredItems"
        :key="getItemKey(item)"
        class="list-card"
        @click="handleItemClick(item)"
      >
        <slot name="item" :item="item">
          <!-- 默认列表项渲染 -->
          <div class="list-item-default">
            <div class="item-header">
              <el-tag v-if="getItemTag(item)" size="small">{{ getItemTag(item) }}</el-tag>
              <h3>{{ getItemTitle(item) }}</h3>
            </div>
            <p class="item-desc">{{ getItemDesc(item) }}</p>
          </div>
        </slot>
      </el-card>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Search } from '@element-plus/icons-vue'

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
  // search 专用：搜索的字段列表
  searchFields?: string[]
}

interface Props {
  items: any[]
  loading?: boolean
  emptyText?: string
  filters?: ListFilterConfig[]
  itemKey?: string
  titleField?: string
  descField?: string
  tagField?: string
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyText: '暂无数据',
  itemKey: 'path',
  titleField: 'title',
  descField: 'description',
  tagField: '',
})

const emit = defineEmits<{
  itemClick: [item: any]
  filterChange: [filters: Record<string, string>]
}>()

const route = useRoute()

// 筛选值
const filterValues = ref<Record<string, string>>({})

// 是否显示筛选栏
const showFilters = computed(() => props.filters && props.filters.length > 0)

// 筛选后的数据
const filteredItems = computed(() => {
  let result = props.items

  for (const filter of (props.filters || [])) {
    const value = filterValues.value[filter.key]
    if (!value || value === 'all') continue

    if (filter.type === 'search') {
      // 搜索：在指定字段中查找
      const query = value.toLowerCase()
      const fields = filter.searchFields || [props.titleField, props.descField]
      result = result.filter(item => 
        fields.some(field => {
          const fieldValue = item[field]
          return fieldValue && String(fieldValue).toLowerCase().includes(query)
        })
      )
    } else {
      // 普通筛选：精确匹配
      result = result.filter(item => item[filter.key] === value)
    }
  }

  return result
})

// 获取项目 key
function getItemKey(item: any): string {
  return item[props.itemKey] || JSON.stringify(item)
}

// 获取标题
function getItemTitle(item: any): string {
  return item[props.titleField] || item.name || item.title || ''
}

// 获取描述
function getItemDesc(item: any): string {
  const desc = item[props.descField] || item.description || item.content || ''
  return desc.length > 100 ? desc.slice(0, 100) + '...' : desc
}

// 获取标签
function getItemTag(item: any): string {
  return props.tagField ? item[props.tagField] : ''
}

// 处理项目点击
function handleItemClick(item: any) {
  emit('itemClick', item)
}

// 处理筛选变化
function handleFilterChange() {
  emit('filterChange', { ...filterValues.value })
}

// 初始化筛选值
onMounted(() => {
  if (props.filters) {
    for (const filter of props.filters) {
      // 从 props 获取默认值
      filterValues.value[filter.key] = filter.defaultValue || ''
      
      // 从 URL 参数读取
      const queryValue = route.query[filter.key] as string
      if (queryValue) {
        filterValues.value[filter.key] = queryValue
      }
    }
  }
})

// 暴露筛选值供外部使用
defineExpose({
  filterValues,
  filteredItems,
})
</script>

<style scoped>
.list-view {
  width: 100%;
}

.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 20px;
  align-items: center;
}

.list-card {
  cursor: pointer;
  margin-bottom: 16px;
  transition: all 0.3s;
}

.list-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.list-item-default {
  width: 100%;
}

.item-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.item-header h3 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.item-desc {
  margin: 0;
  color: #909399;
  font-size: 14px;
  line-height: 1.5;
}
</style>
