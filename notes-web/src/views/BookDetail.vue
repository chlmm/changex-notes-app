<template>
  <div class="book-detail">
    <div class="back-btn">
      <el-button @click="$router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回列表
      </el-button>
    </div>

    <el-skeleton v-if="loading" :rows="10" animated />

    <template v-else-if="book">
      <!-- 头部信息 -->
      <el-card class="info-card">
        <div class="book-header">
          <div class="book-info">
            <h1 class="book-title">{{ book.title }}</h1>
            <div class="book-meta">
              <span class="meta-item">
                <el-icon><User /></el-icon>
                {{ book.author || '未知作者' }}
              </span>
              <span class="meta-item">
                <el-icon><Folder /></el-icon>
                {{ book.category || '未分类' }}
              </span>
              <span class="meta-item">
                <el-icon><Calendar /></el-icon>
                {{ book.created }}
              </span>
            </div>
            <div class="book-tags">
              <el-tag v-for="tag in book.tags" :key="tag" size="small" type="info">
                {{ tag }}
              </el-tag>
            </div>
          </div>
          <div class="book-status">
            <StatusTag :status="book.status || ''" />
            <ProgressBar v-if="hasProgress" :content="book.content || ''" />
          </div>
        </div>
      </el-card>

      <!-- 内容区域 -->
      <el-card class="content-card">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="学习笔记" name="note">
            <MdRenderer :content="book.content || ''" />
          </el-tab-pane>
          <el-tab-pane label="金句摘录" name="quotes" :disabled="!book.quotes">
            <MdRenderer v-if="book.quotes" :content="book.quotes" />
            <el-empty v-else description="暂无金句摘录" />
          </el-tab-pane>
          <el-tab-pane label="读后感" name="reviews" :disabled="!book.reviews">
            <MdRenderer v-if="book.reviews" :content="book.reviews" />
            <el-empty v-else description="暂无读后感" />
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </template>

    <el-empty v-else description="书籍不存在" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useNotesStore } from '@/stores/notes'
import { useAsync } from '@/composables/useAsync'
import MdRenderer from '@/components/MdRenderer.vue'
import StatusTag from '@/components/StatusTag.vue'
import ProgressBar from '@/components/ProgressBar.vue'
import type { BookNote } from '@/types/note'

const route = useRoute()
const notesStore = useNotesStore()

const { data: book, loading, execute } = useAsync<BookNote | null>()
const activeTab = ref('note')

// 检查是否有进度信息
const hasProgress = computed(() => {
  const content = book.value?.content || ''
  return content.includes('- [ ]') || content.includes('- [x]')
})

onMounted(async () => {
  const type = route.params.type as string
  const name = decodeURIComponent(route.params.name as string)
  
  await execute(() => notesStore.getBookDetail(type, name) as Promise<BookNote | null>)
})
</script>

<style scoped>
.book-detail {
  max-width: 1000px;
  margin: 0 auto;
}

.back-btn {
  margin-bottom: 16px;
}

.info-card {
  margin-bottom: 20px;
}

.book-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.book-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
}

.book-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #909399;
}

.book-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.book-status {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
}

.content-card {
  margin-bottom: 20px;
}

:deep(.el-tabs__header) {
  margin-bottom: 20px;
}
</style>
