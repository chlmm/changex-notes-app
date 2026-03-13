<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>{{ title }}</span>
        <el-button v-if="showMore" text type="primary" @click="handleMore">
          查看全部
        </el-button>
      </div>
    </template>
    <div class="progress-stats">
      <div v-for="item in items" :key="item.label" class="progress-item">
        <span class="progress-label">{{ item.label }}</span>
        <el-progress :percentage="item.percentage" :stroke-width="12" :status="item.status" />
        <span class="progress-count">{{ item.count }}</span>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
interface ProgressItem {
  label: string
  percentage: number
  count: string
  status?: '' | 'success' | 'warning' | 'exception'
}

interface Props {
  title: string
  items: ProgressItem[]
  showMore?: boolean
}

withDefaults(defineProps<Props>(), {
  showMore: false,
})

const emit = defineEmits<{
  more: []
}>()

function handleMore() {
  emit('more')
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.progress-stats {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.progress-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.progress-label {
  width: 60px;
  font-size: 14px;
  color: #606266;
}

.progress-item :deep(.el-progress) {
  flex: 1;
}

.progress-count {
  width: 50px;
  text-align: right;
  font-size: 14px;
  color: #909399;
}
</style>
