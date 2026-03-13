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
import { GRADIENTS, ICONS } from '@/composables/useDashboard'
import Dashboard from '@/components/dashboard/Dashboard.vue'
import type { ChartDataItem } from '@/components/dashboard/widgets/ChartView.vue'

const notesStore = useNotesStore()

// Dashboard 标题
const dashboardTitle = '欢迎来到笔记库'

// 计算阅读进度百分比
const readingPercent = computed(() => {
  const total = notesStore.stats.totalBooks
  if (total === 0) return 0
  return Math.round((notesStore.stats.readingBooks / total) * 100)
})

const completedPercent = computed(() => {
  const total = notesStore.stats.totalBooks
  if (total === 0) return 0
  return Math.round((notesStore.stats.completedBooks / total) * 100)
})

// 1. 笔记分布图表数据
const chartData = computed<ChartDataItem[]>(() => {
  const total = notesStore.stats.totalBooks +
    notesStore.stats.totalVideos +
    notesStore.stats.totalKnowledge +
    notesStore.stats.totalSkills +
    notesStore.stats.totalProblems +
    notesStore.stats.totalIndex

  const data = [
    { name: '书籍', value: notesStore.stats.totalBooks },
    { name: '视频', value: notesStore.stats.totalVideos },
    { name: '知识点', value: notesStore.stats.totalKnowledge },
    { name: '技能', value: notesStore.stats.totalSkills },
    { name: '问题', value: notesStore.stats.totalProblems },
    { name: '收藏', value: notesStore.stats.totalIndex },
  ]

  return data.map(item => ({
    ...item,
    percentage: total > 0 ? ((item.value / total) * 100).toFixed(1) : '0',
  }))
})

// 2. 最近更新时间线数据
const recentNotes = computed(() => {
  const allNotes: Array<{ title: string; type: string; updatedAt?: string; created?: string; description?: string }> = []

  // 收集所有笔记
  notesStore.books.forEach(b => {
    allNotes.push({
      title: b.title,
      type: 'book',
      updatedAt: b.updatedAt || b.created,
      created: b.created,
      description: b.author ? `${b.author} - ${b.category || ''}` : b.category,
    })
  })

  notesStore.videos.forEach(v => {
    allNotes.push({
      title: v.title,
      type: 'video',
      updatedAt: v.updatedAt || v.created,
      created: v.created,
      description: v.platform,
    })
  })

  notesStore.knowledge.forEach(k => {
    allNotes.push({
      title: k.title,
      type: 'knowledge',
      updatedAt: k.updatedAt || k.created,
      created: k.created,
      description: `${k.category} - ${k.topic}`,
    })
  })

  notesStore.skills.forEach(s => {
    allNotes.push({
      title: s.name || s.title,
      type: 'skill',
      updatedAt: s.updatedAt || s.created,
      created: s.created,
      description: s.category,
    })
  })

  notesStore.problems.forEach(p => {
    allNotes.push({
      title: p.title,
      type: 'problem',
      updatedAt: p.updatedAt || p.created,
      created: p.created,
      description: p.type,
    })
  })

  // 按时间排序，取前10条
  return allNotes
    .sort((a, b) => {
      const timeA = new Date(a.updatedAt || a.created || 0).getTime()
      const timeB = new Date(b.updatedAt || b.created || 0).getTime()
      return timeB - timeA
    })
    .slice(0, 10)
})

// 3. 正在阅读的书籍
const readingBooks = computed(() => {
  return notesStore.books
    .filter(b => b.status === '阅读中')
    .slice(0, 6)
})

// 4. 待解决问题
const pendingProblems = computed(() => {
  return notesStore.problems.slice(0, 5)
})

// 5. 本周活动统计
const weeklyActivity = computed(() => {
  const now = new Date()
  const weekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)

  let newBooks = 0
  let newVideos = 0
  let newKnowledge = 0
  let newSkills = 0
  let newProblems = 0

  notesStore.books.forEach(b => {
    const date = new Date(b.created || 0)
    if (date >= weekAgo) newBooks++
  })

  notesStore.videos.forEach(v => {
    const date = new Date(v.created || 0)
    if (date >= weekAgo) newVideos++
  })

  notesStore.knowledge.forEach(k => {
    const date = new Date(k.created || 0)
    if (date >= weekAgo) newKnowledge++
  })

  notesStore.skills.forEach(s => {
    const date = new Date(s.created || 0)
    if (date >= weekAgo) newSkills++
  })

  notesStore.problems.forEach(p => {
    const date = new Date(p.created || 0)
    if (date >= weekAgo) newProblems++
  })

  return {
    total: newBooks + newVideos + newKnowledge + newSkills + newProblems,
    newBooks,
    newVideos,
    newKnowledge,
    newSkills,
    newProblems,
  }
})

// Dashboard 配置
const dashboardRows = computed(() => [
  // 第一行：统计卡片
  {
    gutter: 16,
    widgets: [
      {
        type: 'stats' as const,
        icon: ICONS.reading,
        value: notesStore.stats.totalBooks,
        label: '书籍',
        gradient: GRADIENTS.blue,
        route: '/books',
      },
      {
        type: 'stats' as const,
        icon: ICONS.videoPlay,
        value: notesStore.stats.totalVideos,
        label: '视频',
        gradient: GRADIENTS.green,
        route: '/videos',
      },
      {
        type: 'stats' as const,
        icon: ICONS.collection,
        value: notesStore.stats.totalKnowledge,
        label: '知识点',
        gradient: GRADIENTS.gray,
        route: '/knowledge',
      },
      {
        type: 'stats' as const,
        icon: ICONS.tools,
        value: notesStore.stats.totalSkills,
        label: '技能',
        gradient: GRADIENTS.purple,
        route: '/skills',
      },
      {
        type: 'stats' as const,
        icon: ICONS.questionFilled,
        value: notesStore.stats.totalProblems,
        label: '问题',
        gradient: GRADIENTS.orange,
        route: '/problems',
      },
      {
        type: 'stats' as const,
        icon: ICONS.star,
        value: notesStore.stats.totalIndex,
        label: '收藏',
        gradient: GRADIENTS.red,
        route: '/index',
      },
    ],
  },
  // 第二行：进度 + 快捷入口
  {
    gutter: 20,
    widgets: [
      {
        type: 'progress' as const,
        title: '阅读进度',
        showMore: true,
        moreRoute: '/books',
        items: [
          {
            label: '阅读中',
            percentage: readingPercent.value,
            count: `${notesStore.stats.readingBooks} 本`,
            status: 'warning' as const,
          },
          {
            label: '已完成',
            percentage: completedPercent.value,
            count: `${notesStore.stats.completedBooks} 本`,
            status: 'success' as const,
          },
        ],
      },
      {
        type: 'quickEntry' as const,
        title: '快速入口',
        items: [
          {
            icon: ICONS.collection,
            label: '知识库',
            color: '#909399',
            route: '/knowledge',
          },
          {
            icon: ICONS.tools,
            label: '技能库',
            color: '#667eea',
            route: '/skills',
          },
          {
            icon: ICONS.questionFilled,
            label: '问题解决',
            color: '#e6a23c',
            route: '/problems',
          },
          {
            icon: ICONS.star,
            label: '索引收藏',
            color: '#f56c6c',
            route: '/index',
          },
        ],
      },
    ],
  },
  // 第三行：快捷链接
  {
    gutter: 20,
    widgets: [
      {
        type: 'quickLink' as const,
        icon: ICONS.reading,
        title: '书籍学习',
        description: `${notesStore.stats.totalBooks} 本书籍`,
        color: '#409eff',
        route: '/books',
      },
      {
        type: 'quickLink' as const,
        icon: ICONS.videoPlay,
        title: '视频学习',
        description: `${notesStore.stats.totalVideos} 个视频`,
        color: '#67c23a',
        route: '/videos',
      },
    ],
  },
  // 第四行：笔记分布图表 + 最近更新时间线
  {
    gutter: 20,
    widgets: [
      {
        type: 'chartView' as const,
        title: '笔记分布',
        data: chartData.value,
        defaultType: 'pie' as const,
        height: '320px',
        showTypeSelector: true,
        showLegend: true,
      },
      {
        type: 'timelineView' as const,
        items: recentNotes.value,
        titleField: 'title',
        descField: 'description',
        timeField: 'updatedAt',
        typeField: 'type',
        viewMode: 'vertical' as const,
        showToolbar: false,
        searchable: false,
        showMeta: true,
      },
    ],
  },
  // 第五行：正在阅读 + 待解决问题
  {
    gutter: 20,
    widgets: [
      {
        type: 'gridView' as const,
        items: readingBooks.value,
        cardType: 'book' as const,
        gridColumns: 'repeat(auto-fill, minmax(200px, 1fr))',
        gridGap: '16px',
        emptyText: '暂无正在阅读的书籍',
      },
      {
        type: 'listView' as const,
        items: pendingProblems.value,
        titleField: 'title',
        descField: 'type',
        tagField: 'status',
        emptyText: '暂无待解决问题',
      },
    ],
  },
  // 第六行：本周活动
  {
    gutter: 16,
    widgets: [
      {
        type: 'stats' as const,
        icon: ICONS.reading,
        value: weeklyActivity.value.newBooks,
        label: '本周新增书籍',
        gradient: GRADIENTS.blue,
        route: '/books',
      },
      {
        type: 'stats' as const,
        icon: ICONS.videoPlay,
        value: weeklyActivity.value.newVideos,
        label: '本周新增视频',
        gradient: GRADIENTS.green,
        route: '/videos',
      },
      {
        type: 'stats' as const,
        icon: ICONS.collection,
        value: weeklyActivity.value.newKnowledge,
        label: '本周新增知识点',
        gradient: GRADIENTS.gray,
        route: '/knowledge',
      },
      {
        type: 'stats' as const,
        icon: ICONS.tools,
        value: weeklyActivity.value.newSkills,
        label: '本周新增技能',
        gradient: GRADIENTS.purple,
        route: '/skills',
      },
      {
        type: 'stats' as const,
        icon: ICONS.questionFilled,
        value: weeklyActivity.value.newProblems,
        label: '本周新增问题',
        gradient: GRADIENTS.orange,
        route: '/problems',
      },
      {
        type: 'stats' as const,
        icon: ICONS.star,
        value: weeklyActivity.value.total,
        label: '本周总计',
        gradient: GRADIENTS.red,
        route: '/',
      },
    ],
  },
])
</script>
