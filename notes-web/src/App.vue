<template>
  <div class="app-container">
    <el-container>
      <el-aside width="var(--sidebar-width)" class="sidebar">
        <Sidebar />
      </el-aside>
      <el-container>
        <el-header class="header">
          <div class="header-title">
            <el-icon><Notebook /></el-icon>
            <span>笔记库</span>
          </div>
          <div class="header-search">
            <el-input
              v-model="searchQuery"
              placeholder="搜索笔记..."
              clearable
              @keyup.enter="handleSearch"
              @focus="showSearchPanel = true"
            >
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
            <div v-if="showSearchPanel && searchResults.length > 0" class="search-panel">
              <div
                v-for="result in searchResults"
                :key="result.id"
                class="search-item"
                @click="goToResult(result)"
              >
                <div class="result-header">
                  <el-tag size="small" type="info">{{ result.typeId }}</el-tag>
                  <span class="result-title">{{ getTitle(result) }}</span>
                </div>
                <p class="result-snippet">{{ getSnippet(result) }}</p>
              </div>
            </div>
          </div>
        </el-header>
        <el-main class="main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '@/components/Sidebar.vue'
import { useNotesStore } from '@/stores/notes'
import type { GenericNote } from '@/types/note'

const router = useRouter()
const notesStore = useNotesStore()

const searchQuery = ref('')
const searchResults = ref<GenericNote[]>([])
const showSearchPanel = ref(false)

onMounted(async () => {
  await notesStore.loadNotes()
})

async function handleSearch() {
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    return
  }
  searchResults.value = await notesStore.searchNotes(searchQuery.value)
  showSearchPanel.value = true
}

function getTitle(note: GenericNote): string {
  const f = note.fields
  return String(f['title'] || f['source.title'] || f['quote'] || f['name'] || note.id)
}

function getSnippet(note: GenericNote): string {
  return (note.content || '').slice(0, 200) + (note.content && note.content.length > 200 ? '...' : '')
}

function goToResult(note: GenericNote) {
  showSearchPanel.value = false
  searchQuery.value = ''
  router.push({
    name: 'GenericDetail',
    params: { typeId: note.typeId },
    query: { id: note.id },
  })
}

document.addEventListener('click', (e) => {
  const target = e.target as HTMLElement
  if (!target.closest('.header-search')) {
    showSearchPanel.value = false
  }
})
</script>

<style scoped>
.app-container { height: 100%; }
.sidebar {
  background: #fff;
  border-right: 1px solid #e4e7ed;
  height: 100vh;
  overflow-y: auto;
}
.header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: var(--header-height);
}
.header-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}
.header-search { position: relative; width: 300px; }
.search-panel {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  max-height: 400px;
  overflow-y: auto;
  z-index: 1000;
}
.search-item {
  padding: 12px 16px;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
}
.search-item:last-child { border-bottom: none; }
.search-item:hover { background: #f5f7fa; }
.result-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}
.result-title { font-size: 14px; font-weight: 500; color: #303133; }
.result-snippet {
  margin: 0;
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}
.main {
  background: var(--bg-color);
  padding: 20px;
  overflow-y: auto;
  height: calc(100vh - var(--header-height));
}
</style>
