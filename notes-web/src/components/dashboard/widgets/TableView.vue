<template>
  <div class="table-view">
    <!-- 工具栏：筛选器 + 操作按钮 -->
    <div class="toolbar">
      <div class="filters">
        <template v-for="filter in filters" :key="filter.key">
          <!-- 搜索框 -->
          <el-input
            v-if="filter.type === 'search'"
            v-model="filterValues[filter.key]"
            :placeholder="filter.placeholder || '搜索...'"
            clearable
            size="small"
            :style="{ width: filter.width || '200px' }"
            @input="handleFilterChange"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          
          <!-- 单选按钮组 -->
          <el-radio-group
            v-else-if="filter.type === 'radio'"
            v-model="filterValues[filter.key]"
            size="small"
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
          
          <!-- 下拉选择 -->
          <el-select
            v-else-if="filter.type === 'select'"
            v-model="filterValues[filter.key]"
            :placeholder="filter.placeholder"
            clearable
            size="small"
            :style="{ width: filter.width || '120px' }"
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
      
      <!-- 操作按钮插槽 -->
      <div class="actions">
        <slot name="actions" />
      </div>
    </div>

    <!-- 表格 -->
    <el-table
      :data="pagedData"
      :row-key="rowKey"
      :height="height"
      :max-height="maxHeight"
      :stripe="stripe"
      :border="border"
      :highlight-current-row="highlightCurrentRow"
      :default-sort="defaultSort"
      v-loading="loading"
      @row-click="handleRowClick"
      @sort-change="handleSortChange"
      @selection-change="handleSelectionChange"
    >
      <!-- 选择列 -->
      <el-table-column v-if="selectable" type="selection" width="50" />
      
      <!-- 序号列 -->
      <el-table-column v-if="showIndex" type="index" label="#" width="60" />
      
      <!-- 数据列 -->
      <el-table-column
        v-for="col in columns"
        :key="col.prop"
        :prop="col.prop"
        :label="col.label"
        :width="col.width"
        :min-width="col.minWidth"
        :sortable="col.sortable"
        :align="col.align"
        :fixed="col.fixed"
      >
        <template #default="{ row }">
          <!-- 自定义列渲染 -->
          <template v-if="col.slot">
            <slot :name="`col-${col.prop}`" :row="row" :value="row[col.prop]">
              {{ row[col.prop] }}
            </slot>
          </template>
          
          <!-- 标签类型 -->
          <template v-else-if="col.type === 'tag'">
            <el-tag
              :type="getTagType(row, col)"
              :size="col.tagSize || 'small'"
            >
              {{ formatValue(row, col) }}
            </el-tag>
          </template>
          
          <!-- 链接类型 -->
          <template v-else-if="col.type === 'link'">
            <el-link
              :type="col.linkType || 'primary'"
              @click.stop="col.onClick?.(row)"
            >
              {{ formatValue(row, col) }}
            </el-link>
          </template>
          
          <!-- 进度条类型 -->
          <template v-else-if="col.type === 'progress'">
            <el-progress
              :percentage="row[col.prop] || 0"
              :status="col.getProgressStatus?.(row)"
              :stroke-width="6"
            />
          </template>
          
          <!-- 日期类型 -->
          <template v-else-if="col.type === 'date'">
            {{ formatDate(row[col.prop], col.dateFormat) }}
          </template>
          
          <!-- 图片类型 -->
          <template v-else-if="col.type === 'image'">
            <el-image
              :src="row[col.prop]"
              :style="{ width: col.imageWidth || '50px', height: col.imageHeight || '50px' }"
              fit="cover"
            />
          </template>
          
          <!-- 默认文本 -->
          <template v-else>
            {{ formatValue(row, col) }}
          </template>
        </template>
      </el-table-column>
      
      <!-- 操作列 -->
      <el-table-column
        v-if="$slots.operation"
        label="操作"
        :width="operationWidth"
        :fixed="operationFixed"
      >
        <template #default="{ row }">
          <slot name="operation" :row="row" />
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div v-if="pagination" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="pageSizes"
        :total="filteredData.length"
        :layout="paginationLayout"
        @size-change="handlePageSizeChange"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Search } from '@element-plus/icons-vue'

// 筛选器配置
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

// 列配置
export interface TableColumnConfig {
  prop: string
  label: string
  width?: number | string
  minWidth?: number | string
  sortable?: boolean | 'custom'
  align?: 'left' | 'center' | 'right'
  fixed?: boolean | 'left' | 'right'
  // 列类型
  type?: 'tag' | 'link' | 'progress' | 'date' | 'image'
  // 自定义插槽
  slot?: boolean
  // tag 类型
  tagSize?: 'large' | 'default' | 'small'
  getTagType?: (row: any) => 'success' | 'warning' | 'danger' | 'info' | 'primary'
  // link 类型
  linkType?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'default'
  onClick?: (row: any) => void
  // progress 类型
  getProgressStatus?: (row: any) => '' | 'success' | 'warning' | 'exception'
  // date 类型
  dateFormat?: string
  // image 类型
  imageWidth?: string
  imageHeight?: string
  // 格式化函数
  formatter?: (value: any, row: any) => string
}

interface Props {
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
  // 筛选
  filters?: TableFilterConfig[]
  // 分页
  pagination?: boolean
  pageSize?: number
  pageSizes?: number[]
  paginationLayout?: string
  // 选择
  selectable?: boolean
  showIndex?: boolean
  // 操作列
  operationWidth?: number | string
  operationFixed?: boolean | 'left' | 'right'
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  rowKey: 'id',
  stripe: true,
  border: false,
  highlightCurrentRow: false,
  pagination: true,
  pageSize: 10,
  pageSizes: () => [10, 20, 50, 100],
  paginationLayout: 'total, sizes, prev, pager, next, jumper',
  selectable: false,
  showIndex: false,
  operationWidth: 150,
  operationFixed: 'right',
})

const emit = defineEmits<{
  rowClick: [row: any]
  sortChange: [sort: { prop: string; order: string | null }]
  selectionChange: [selection: any[]]
  filterChange: [filters: Record<string, string>]
  pageChange: [page: number]
  pageSizeChange: [size: number]
}>()

const route = useRoute()

// 筛选值
const filterValues = ref<Record<string, string>>({})

// 分页
const currentPage = ref(1)
const pageSize = ref(props.pageSize)

// 排序
const currentSort = ref<{ prop: string; order: string | null }>({
  prop: '',
  order: null,
})

// 筛选后的数据
const filteredData = computed(() => {
  let result = props.data

  // 应用筛选
  for (const [key, value] of Object.entries(filterValues.value)) {
    if (value && value !== 'all') {
      const filter = props.filters?.find(f => f.key === key)
      
      if (filter?.type === 'search' && filter.searchFields) {
        // 搜索筛选：多字段匹配
        const searchValue = value.toLowerCase()
        result = result.filter(item =>
          filter.searchFields!.some(field =>
            String(item[field] || '').toLowerCase().includes(searchValue)
          )
        )
      } else {
        // 精确匹配
        result = result.filter(item => item[key] === value)
      }
    }
  }

  return result
})

// 排序后的数据
const sortedData = computed(() => {
  if (!currentSort.value.prop || !currentSort.value.order) {
    return filteredData.value
  }

  const { prop, order } = currentSort.value
  return [...filteredData.value].sort((a, b) => {
    const aVal = a[prop]
    const bVal = b[prop]
    
    if (aVal === bVal) return 0
    
    let result = aVal > bVal ? 1 : -1
    if (order === 'descending') result = -result
    
    return result
  })
})

// 分页数据
const pagedData = computed(() => {
  if (!props.pagination) return sortedData.value
  
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  
  return sortedData.value.slice(start, end)
})

// 格式化值
function formatValue(row: any, col: TableColumnConfig): string {
  const value = row[col.prop]
  if (col.formatter) {
    return col.formatter(value, row)
  }
  return value ?? ''
}

// 格式化日期
function formatDate(value: any, format = 'YYYY-MM-DD'): string {
  if (!value) return ''
  const date = new Date(value)
  if (isNaN(date.getTime())) return value
  
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  
  return format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
}

// 获取 tag 类型
function getTagType(row: any, col: TableColumnConfig): 'success' | 'warning' | 'danger' | 'info' | 'primary' {
  if (col.getTagType) {
    const type = col.getTagType(row)
    return type || 'primary'
  }
  return 'primary'
}

// 事件处理
function handleRowClick(row: any) {
  emit('rowClick', row)
}

function handleSortChange(sort: { prop: string; order: string | null }) {
  currentSort.value = sort
  emit('sortChange', sort)
}

function handleSelectionChange(selection: any[]) {
  emit('selectionChange', selection)
}

function handleFilterChange() {
  currentPage.value = 1 // 重置分页
  emit('filterChange', { ...filterValues.value })
}

function handlePageChange(page: number) {
  emit('pageChange', page)
}

function handlePageSizeChange(size: number) {
  currentPage.value = 1
  emit('pageSizeChange', size)
}

// 初始化
onMounted(() => {
  // 初始化筛选默认值
  if (props.filters) {
    for (const filter of props.filters) {
      filterValues.value[filter.key] = filter.defaultValue || ''
    }
    
    // 从 URL 参数读取
    for (const filter of props.filters) {
      const queryValue = route.query[filter.key] as string
      if (queryValue) {
        filterValues.value[filter.key] = queryValue
      }
    }
  }
})

// 暴露方法
defineExpose({
  filterValues,
  currentPage,
  pageSize,
  filteredData,
  sortedData,
  pagedData,
  refresh: () => {
    currentPage.value = 1
  },
})
</script>

<style scoped>
.table-view {
  width: 100%;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.filters {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.actions {
  display: flex;
  gap: 8px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
