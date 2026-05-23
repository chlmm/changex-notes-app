<template>
  <el-card class="note-card" shadow="hover" @click="$emit('click')">
    <div class="card-header">
      <span class="title">{{ title }}</span>
      <StatusTag :status="(note.fields['status'] as string) || ''" />
    </div>
    
    <div class="card-meta">
      <span v-if="note.fields['author']" class="meta-item">
        <el-icon><User /></el-icon>
        {{ note.fields['author'] }}
      </span>
      <span v-if="note.fields['category']" class="meta-item">
        <el-icon><Folder /></el-icon>
        {{ note.fields['category'] }}
      </span>
      <span v-if="note.fields['platform']" class="meta-item">
        <el-icon><VideoPlay /></el-icon>
        {{ note.fields['platform'] }}
      </span>
    </div>
    
    <div class="card-tags">
      <el-tag 
        v-for="tag in tags" 
        :key="tag" 
        size="small" 
        type="info"
        class="tag"
      >
        {{ tag }}
      </el-tag>
    </div>
    
    <div class="card-footer">
      <span class="date">{{ note.fields['created'] || note.updatedAt }}</span>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { GenericNote } from '@/types/note'
import StatusTag from './StatusTag.vue'

const props = defineProps<{
  note: GenericNote
  type: string
}>()

defineEmits(['click'])

const title = computed(() =>
  props.note.fields['title'] || props.note.fields['source.title'] || props.note.fields['name'] || props.note.id
)

const tags = computed(() => {
  const t = props.note.fields['tags']
  if (Array.isArray(t)) return t.slice(0, 3)
  return []
})
</script>

<style scoped>
.note-card { cursor: pointer; margin-bottom: 12px; }
.card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.title { font-size: 14px; font-weight: 500; color: #303133; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1; margin-right: 8px; }
.card-meta { display: flex; gap: 12px; font-size: 12px; color: #909399; margin-bottom: 8px; flex-wrap: wrap; }
.meta-item { display: flex; align-items: center; gap: 4px; }
.card-tags { display: flex; gap: 4px; flex-wrap: wrap; margin-bottom: 8px; }
.card-footer { display: flex; justify-content: space-between; align-items: center; font-size: 12px; color: #c0c4cc; }
</style>
