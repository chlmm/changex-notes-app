<template>
  <el-card class="note-card" shadow="hover" @click="$emit('click')">
    <div class="card-header">
      <span class="title">{{ note.title }}</span>
      <StatusTag :status="note.status || ''" />
    </div>
    
    <div class="card-meta">
      <template v-if="type === 'book'">
        <span class="meta-item">
          <el-icon><User /></el-icon>
          {{ (note as BookNote).author || '未知作者' }}
        </span>
        <span class="meta-item">
          <el-icon><Folder /></el-icon>
          {{ (note as BookNote).category || '未分类' }}
        </span>
      </template>
      <template v-else>
        <span class="meta-item">
          <el-icon><VideoPlay /></el-icon>
          {{ (note as VideoNote).platform }}
        </span>
      </template>
    </div>
    
    <div class="card-tags">
      <el-tag 
        v-for="tag in (note.tags || []).slice(0, 3)" 
        :key="tag" 
        size="small" 
        type="info"
        class="tag"
      >
        {{ tag }}
      </el-tag>
      <el-tag v-if="(note.tags || []).length > 3" size="small" type="info">
        +{{ (note.tags || []).length - 3 }}
      </el-tag>
    </div>
    
    <div class="card-footer">
      <span class="date">{{ note.created }}</span>
      <ProgressBar v-if="type === 'book'" :content="note.content || ''" />
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { BookNote, VideoNote } from '@/types/note'
import StatusTag from './StatusTag.vue'
import ProgressBar from './ProgressBar.vue'

defineProps<{
  note: BookNote | VideoNote
  type: 'book' | 'video'
}>()

defineEmits(['click'])
</script>

<style scoped>
.note-card {
  cursor: pointer;
  transition: all 0.3s;
}

.note-card:hover {
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  flex: 1;
  margin-right: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #909399;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.tag {
  margin: 0;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
}

.date {
  font-size: 12px;
  color: #c0c4cc;
}
</style>
