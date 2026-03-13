<template>
  <div class="progress-bar">
    <el-progress 
      :percentage="percentage" 
      :stroke-width="6"
      :show-text="false"
      :color="color"
    />
    <span class="progress-text">{{ completed }}/{{ total }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  content: string
}>()

// 解析 checkbox 计算进度
const { completed, total, percentage } = computed(() => {
  const checkboxRegex = /- \[([ xX])\]/g
  const matches = props.content.match(checkboxRegex) || []
  const total = matches.length
  const completed = matches.filter(m => m.includes('x') || m.includes('X')).length
  const percentage = total > 0 ? Math.round((completed / total) * 100) : 0
  return { completed, total, percentage }
}).value

const color = computed(() => {
  if (percentage >= 100) return '#67c23a'
  if (percentage >= 50) return '#e6a23c'
  return '#409eff'
})
</script>

<style scoped>
.progress-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100px;
}

.progress-text {
  font-size: 12px;
  color: #909399;
  white-space: nowrap;
}
</style>
