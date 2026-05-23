<template>
  <Dashboard
    :title="dashboardTitle"
    :loading="notesStore.loading"
    :rows="dashboardRows"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useNotesStore } from '@/stores/notes'
import { ICONS, GRADIENTS } from '@/composables/useDashboard'
import Dashboard from '@/components/dashboard/Dashboard.vue'

const notesStore = useNotesStore()

const dashboardTitle = '欢迎来到笔记库'
const s = computed(() => notesStore.stats)

const dashboardRows = computed(() => [
  {
    gutter: 16,
    widgets: [
      { type: 'stats' as const, icon: ICONS.reading, value: s.value.book || 0, label: '书籍', gradient: GRADIENTS.blue, route: '/book' },
      { type: 'stats' as const, icon: ICONS.videoPlay, value: s.value.video || 0, label: '视频', gradient: GRADIENTS.green, route: '/video' },
      { type: 'stats' as const, icon: ICONS.collection, value: s.value.knowledge || 0, label: '知识库', gradient: GRADIENTS.gray, route: '/knowledge' },
      { type: 'stats' as const, icon: ICONS.tools, value: s.value.executor_skill || 0, label: '技能', gradient: GRADIENTS.purple, route: '/executor_skill' },
      { type: 'stats' as const, icon: ICONS.questionFilled, value: s.value.project || 0, label: '项目', gradient: GRADIENTS.orange, route: '/project' },
      { type: 'stats' as const, icon: ICONS.star, value: s.value.github_repo || 0, label: 'GitHub', gradient: GRADIENTS.red, route: '/github_repo' },
    ],
  },
  {
    gutter: 20,
    widgets: [
      { type: 'quickLink' as const, icon: ICONS.reading, title: '书籍学习', description: `${s.value.book || 0} 本书籍`, color: '#409eff', route: '/book' },
      { type: 'quickLink' as const, icon: ICONS.videoPlay, title: '视频学习', description: `${s.value.video || 0} 个视频`, color: '#67c23a', route: '/video' },
      { type: 'quickLink' as const, icon: ICONS.star, title: 'GitHub 收藏', description: `${s.value.github_repo || 0} 个项目`, color: '#e6a23c', route: '/github_repo' },
    ],
  },
  {
    gutter: 16,
    widgets: [
      { type: 'stats' as const, icon: ICONS.videoPlay, value: s.value.anime || 0, label: '动漫', gradient: GRADIENTS.blue, route: '/anime' },
      { type: 'stats' as const, icon: ICONS.videoPlay, value: s.value.movie || 0, label: '电影', gradient: GRADIENTS.green, route: '/movie' },
      { type: 'stats' as const, icon: ICONS.tools, value: s.value.game || 0, label: '游戏', gradient: GRADIENTS.purple, route: '/game' },
      { type: 'stats' as const, icon: ICONS.reading, value: s.value.quote || 0, label: '金句', gradient: GRADIENTS.orange, route: '/quote' },
      { type: 'stats' as const, icon: ICONS.collection, value: s.value.project_learning || 0, label: '项目学习', gradient: GRADIENTS.gray, route: '/project_learning' },
    ],
  },
])
</script>
