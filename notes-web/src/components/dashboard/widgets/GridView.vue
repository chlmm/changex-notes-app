<template>
  <div class="grid-view">
    <!-- 筛选器 -->
    <div v-if="filters && filters.length > 0" class="filters">
      <template v-for="filter in filters" :key="filter.key">
        <el-radio-group
          v-if="filter.type === 'radio'"
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

    <!-- 加载状态 -->
    <el-skeleton v-if="loading" :rows="5" animated />

    <!-- 空状态 -->
    <el-empty v-else-if="filteredItems.length === 0" :description="emptyText" />

    <!-- 网格内容 -->
    <div v-else class="grid-container" :style="gridStyle">
      <div
        v-for="item in filteredItems"
        :key="item[itemKey] || JSON.stringify(item)"
        class="grid-item"
        @click="handleItemClick(item)"
      >
        <slot name="item" :item="item">
          <!-- 默认使用 NoteCard (仅支持 book/video) -->
          <NoteCard
            v-if="cardType === 'book' || cardType === 'video'"
            :note="item"
            :type="cardType"
          />
          <!-- 其他类型显示简单卡片 -->
          <el-card v-else class="simple-card" shadow="hover">
            <div class="card-title">{{ item.title || item.name }}</div>
            <div v-if="item.tags?.length" class="card-tags">
              <el-tag v-for="tag in item.tags.slice(0, 3)" :key="tag" size="small" type="info">
                {{ tag }}
              </el-tag>
            </div>
          </el-card>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import NoteCard from '@/components/NoteCard.vue'

export interface FilterOption {
  label: string
  value: string
}

export interface FilterConfig {
  key: string
  type: 'radio' | 'select'
  options: FilterOption[]
  placeholder?: string
  width?: string
  defaultValue?: string
}

interface Props {
  title?: string
  items: any[]
  loading?: boolean
  emptyText?: string
  filters?: FilterConfig[]
  cardType?: 'book' | 'video' | 'anime' | 'movie' | 'game' | 'quote' | 'github'
  gridColumns?: string
  gridGap?: string
  itemKey?: string
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyText: '暂无数据',
  cardType: 'book',
  gridColumns: 'repeat(auto-fill, minmax(320px, 1fr))',
  gridGap: '20px',
  itemKey: 'path',
})

const emit = defineEmits<{
  itemClick: [item: any]
  filterChange: [filters: Record<string, string>]
}>()

const route = useRoute()

// 筛选值
const filterValues = ref<Record<string, string>>({})

// 网格样式
const gridStyle = computed(() => ({
  gridTemplateColumns: props.gridColumns,
  gap: props.gridGap,
}))

// 筛选后的数据
const filteredItems = computed(() => {
  let result = props.items

  // 应用所有筛选器
  for (const [key, value] of Object.entries(filterValues.value)) {
    if (value && value !== 'all') {
      result = result.filter(item => {
        // 支持嵌套属性，如 'type' 或 'status'
        return item[key] === value
      })
    }
  }

  return result
})

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
  // 从 props 初始化默认值
  if (props.filters) {
    for (const filter of props.filters) {
      filterValues.value[filter.key] = filter.defaultValue || ''
    }
  }

  // 从 URL 参数读取
  if (props.filters) {
    for (const filter of props.filters) {
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
.grid-view {
  width: 100%;
}

.filters {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.grid-container {
  display: grid;
}

.grid-item {
  cursor: pointer;
  transition: transform 0.2s;
}

.grid-item:hover {
  transform: translateY(-2px);
}
</style>
