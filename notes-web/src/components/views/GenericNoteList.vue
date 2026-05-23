<template>
  <GenericViewContainer
    :type-id="typeId"
    :notes="notes"
    :loading="loading"
    @item-click="onItemClick"
  />

  <!-- 抽屉/弹窗模式的详情浮层 -->
  <NoteDetailDisplay
    ref="detailDisplayRef"
    :title="detailTitle"
    :load-detail="loadDetailData"
  >
    <template #default="{ detail, loading: detailLoading }">
      <GenericNoteDetail
        v-if="detail"
        :type-id="typeId"
        :note-id="detail.id"
        :overlay="true"
        @close="detailDisplayRef?.close()"
      />
      <el-skeleton v-else-if="detailLoading" :rows="8" animated />
    </template>
  </NoteDetailDisplay>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useSchema } from '@/composables/useSchema'
import { useNoteViewSettings } from '@/composables/useNoteViewSettings'
import { getNoteList } from '@/api/notesV2'
import GenericViewContainer from './GenericViewContainer.vue'
import NoteDetailDisplay from '@/components/NoteDetailDisplay.vue'
import GenericNoteDetail from './GenericNoteDetail.vue'
import type { GenericNote } from '@/types/note'

const props = defineProps<{
  typeId: string
}>()

const { loadSchema } = useSchema()
const { settings } = useNoteViewSettings()
const router = useRouter()

const notes = ref<GenericNote[]>([])
const loading = ref(true)
const detailDisplayRef = ref<InstanceType<typeof NoteDetailDisplay> | null>(null)
const selectedNote = ref<GenericNote | null>(null)

const detailTitle = computed(() => {
  if (!selectedNote.value) return '详情'
  const fields = selectedNote.value.fields
  return String(fields['title'] || fields['source.title'] || fields['source.name'] || fields['name'] || '详情')
})

async function loadData() {
  loading.value = true
  notes.value = []
  try {
    notes.value = await getNoteList(props.typeId)
  } catch (e) {
    console.error('GenericNoteList: load failed', e)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadSchema()
  await loadData()
})

// 同路由参数变化（如 /book → /game），Vue Router 复用组件不 remount，需要 watch 重新加载
watch(() => props.typeId, async () => {
  await loadData()
})

async function loadDetailData(item: unknown): Promise<GenericNote> {
  // NoteDetailDisplay 传过来的是 open 时的 item
  return item as GenericNote
}

function onItemClick(note: GenericNote) {
  const mode = settings.value.openMode
  if (mode === 'page') {
    router.push({
      name: 'GenericDetail',
      params: { typeId: props.typeId },
      query: { id: note.id },
    })
  } else {
    selectedNote.value = note
    detailDisplayRef.value?.open(note)
  }
}
</script>
