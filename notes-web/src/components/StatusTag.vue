<template>
  <el-tag :type="tagType" size="small">
    {{ status }}
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NoteStatus } from '@/types/note'

const props = defineProps<{
  status: NoteStatus | string
}>()

const tagType = computed(() => {
  switch (props.status) {
    case NoteStatus.NotStarted:
    case '未开始':
      return 'info'
    case NoteStatus.Reading:
    case '阅读中':
      return 'warning'
    case NoteStatus.Completed:
    case '已完成':
    case NoteStatus.Interpreted:
    case '已解读':
      return 'success'
    default:
      return 'info'
  }
})
</script>
