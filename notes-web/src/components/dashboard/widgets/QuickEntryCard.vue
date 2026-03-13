<template>
  <el-card>
    <template #header>
      <span>{{ title }}</span>
    </template>
    <div class="quick-entries">
      <div
        v-for="item in items"
        :key="item.label"
        class="entry-item"
        @click="handleClick(item)"
      >
        <el-icon size="32" :color="item.color">
          <component :is="item.icon" />
        </el-icon>
        <span>{{ item.label }}</span>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

interface QuickEntryItem {
  icon: Component
  label: string
  color: string
  route?: string
}

interface Props {
  title: string
  items: QuickEntryItem[]
}

defineProps<Props>()

const emit = defineEmits<{
  click: [item: QuickEntryItem]
}>()

function handleClick(item: QuickEntryItem) {
  emit('click', item)
}
</script>

<style scoped>
.quick-entries {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.entry-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.entry-item:hover {
  background: #f5f7fa;
}

.entry-item span {
  font-size: 12px;
  color: #606266;
}
</style>
