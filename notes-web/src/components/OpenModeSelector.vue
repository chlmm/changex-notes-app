<template>
  <el-dropdown @command="handleChange">
    <el-button type="primary" link>
      <el-icon><View /></el-icon>
      <span>{{ currentLabel }}</span>
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item 
          v-for="option in options" 
          :key="option.value" 
          :command="option.value"
          :class="{ 'is-active': modelValue === option.value }"
        >
          <el-icon><component :is="option.icon" /></el-icon>
          <span>{{ option.label }}</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { View, Monitor, Tickets, Grid } from '@element-plus/icons-vue'
import type { OpenMode } from '@/composables/useNoteViewSettings'

interface Props {
  modelValue: OpenMode
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: OpenMode]
}>()

const options = [
  { value: 'drawer' as OpenMode, label: '侧边抽屉', icon: Tickets },
  { value: 'dialog' as OpenMode, label: '弹窗窗口', icon: Grid },
  { value: 'page' as OpenMode, label: '新页面', icon: Monitor },
]

const currentLabel = computed(() => {
  const option = options.find(o => o.value === props.modelValue)
  return option?.label || '侧边抽屉'
})

function handleChange(value: OpenMode) {
  emit('update:modelValue', value)
}
</script>

<style scoped>
.el-dropdown-menu__item.is-active {
  color: var(--el-color-primary);
  background-color: var(--el-dropdown-menuItem-hover-fill);
}
</style>
