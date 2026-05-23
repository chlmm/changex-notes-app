<template>
  <GenericViewContainer
    :type-id="typeId"
    :notes="notes"
    :loading="loading"
    @item-click="onItemClick"
  />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSchema } from '@/composables/useSchema'
import { getNoteList } from '@/api/notesV2'
import GenericViewContainer from './GenericViewContainer.vue'
import type { GenericNote } from '@/types/note'

const props = defineProps<{
  typeId: string
}>()

const { loadSchema } = useSchema()
const router = useRouter()

const notes = ref<GenericNote[]>([])
const loading = ref(true)

onMounted(async () => {
  await loadSchema()
  try {
    notes.value = await getNoteList(props.typeId)
  } catch (e) {
    console.error('GenericNoteList: load failed', e)
  } finally {
    loading.value = false
  }
})

function onItemClick(note: GenericNote) {
  router.push({
    name: 'GenericDetail',
    params: { typeId: props.typeId },
    query: { id: note.id },
  })
}
</script>
