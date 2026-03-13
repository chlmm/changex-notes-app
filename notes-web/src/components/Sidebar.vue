<template>
  <div class="sidebar">
    <div class="logo">
      <el-icon size="24"><Notebook /></el-icon>
      <span>笔记库</span>
    </div>
    
    <el-menu
      :default-active="activeMenu"
      router
      class="sidebar-menu"
    >
      <el-menu-item index="/">
        <el-icon><HomeFilled /></el-icon>
        <span>首页</span>
      </el-menu-item>
      
      <el-menu-item index="/books">
        <el-icon><Reading /></el-icon>
        <span>书籍学习</span>
      </el-menu-item>
      
      <el-menu-item index="/videos">
        <el-icon><VideoPlay /></el-icon>
        <span>视频学习</span>
      </el-menu-item>

      <el-menu-item index="/knowledge">
        <el-icon><Collection /></el-icon>
        <span>知识库</span>
      </el-menu-item>

      <el-menu-item index="/skills">
        <el-icon><Tools /></el-icon>
        <span>技能库</span>
      </el-menu-item>

      <el-menu-item index="/problems">
        <el-icon><QuestionFilled /></el-icon>
        <span>问题与解决</span>
      </el-menu-item>

      <el-menu-item index="/index">
        <el-icon><Star /></el-icon>
        <span>索引收藏</span>
      </el-menu-item>

      <el-menu-item index="/quotes">
        <el-icon><Edit /></el-icon>
        <span>金句收藏</span>
      </el-menu-item>

      <el-menu-item index="/github">
        <el-icon><Link /></el-icon>
        <span>GitHub 收藏</span>
      </el-menu-item>

      <el-menu-item index="/anime">
        <el-icon><Film /></el-icon>
        <span>动漫收藏</span>
      </el-menu-item>

      <el-menu-item index="/movies">
        <el-icon><VideoCamera /></el-icon>
        <span>电影收藏</span>
      </el-menu-item>

      <el-menu-item index="/games">
        <el-icon><Monitor /></el-icon>
        <span>游戏收藏</span>
      </el-menu-item>
    </el-menu>

    <div class="sidebar-footer">
      <el-tag size="small" type="info">共 {{ totalNotes }} 条笔记</el-tag>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useNotesStore } from '@/stores/notes'
import { Notebook, HomeFilled, Reading, VideoPlay, Collection, Tools, QuestionFilled, Star, Edit, Link, Film, VideoCamera, Monitor } from '@element-plus/icons-vue'

const route = useRoute()
const notesStore = useNotesStore()

const activeMenu = computed(() => {
  if (route.path.startsWith('/books')) return '/books'
  return route.path
})

const totalNotes = computed(() => 
  notesStore.stats.totalBooks + 
  notesStore.stats.totalVideos + 
  notesStore.stats.totalKnowledge + 
  notesStore.stats.totalSkills + 
  notesStore.stats.totalProblems + 
  notesStore.stats.totalIndex +
  notesStore.stats.totalQuotes +
  notesStore.stats.totalGitHub +
  notesStore.stats.totalAnime +
  notesStore.stats.totalMovies +
  notesStore.stats.totalGames
)
</script>

<style scoped>
.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  border-bottom: 1px solid #e4e7ed;
}

.sidebar-menu {
  flex: 1;
  border-right: none;
}

.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid #e4e7ed;
}
</style>
