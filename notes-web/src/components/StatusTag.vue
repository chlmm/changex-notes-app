<template>
  <el-tag :type="tagType" size="small">
    {{ status }}
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  status: string
}>()

const tagType = computed(() => {
  switch (props.status) {
    case '未开始':
    case '想看':
    case '想玩':
    case '未开始':
      return 'info'
    case '阅读中':
    case '追番中':
    case '游玩中':
    case '观看中':
    case '学习中':
      return 'warning'
    case '已完成':
    case '已解读':
    case '已看完':
    case '已通关':
    case '已看':
    case '已完成':
      return 'success'
    default:
      return 'info'
  }
})
</script>
