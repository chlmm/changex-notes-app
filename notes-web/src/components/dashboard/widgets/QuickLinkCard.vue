<template>
  <el-card class="entry-card" @click="handleClick">
    <div class="entry-content">
      <el-icon size="48" :color="color">
        <component :is="icon" />
      </el-icon>
      <div class="entry-info">
        <div class="entry-title">{{ title }}</div>
        <div class="entry-desc">{{ description }}</div>
      </div>
    </div>
    <el-icon class="entry-arrow"><ArrowRight /></el-icon>
  </el-card>
</template>

<script setup lang="ts">
import { ArrowRight } from '@element-plus/icons-vue'
import type { Component } from 'vue'

interface Props {
  icon: Component
  title: string
  description: string
  color: string
  route?: string
}

defineProps<Props>()

const emit = defineEmits<{
  click: []
}>()

function handleClick() {
  emit('click')
}
</script>

<style scoped>
.entry-card {
  cursor: pointer;
  transition: all 0.3s;
}

.entry-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.entry-card :deep(.el-card__body) {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.entry-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.entry-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.entry-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.entry-desc {
  font-size: 14px;
  color: #909399;
}

.entry-arrow {
  color: #c0c4cc;
}
</style>
