import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getStats, getNoteList, searchNotes as searchV2 } from '@/api/notesV2'
import type { GenericNote } from '@/types/note'

export const useNotesStore = defineStore('notes', () => {
  const notesMap = ref<Record<string, GenericNote[]>>({})
  const stats = ref<Record<string, number>>({})
  const loading = ref(false)
  const loaded = ref(false)

  const totalNotes = computed(() =>
    Object.values(stats.value).reduce((a, b) => a + b, 0)
  )

  async function loadStats() {
    try {
      stats.value = await getStats()
    } catch (e) {
      console.error('loadStats failed:', e)
    }
  }

  async function loadNotes(force = false) {
    if (loaded.value && !force) return
    loading.value = true
    await loadStats()
    loaded.value = true
    loading.value = false
  }

  async function loadTypeNotes(typeId: string): Promise<GenericNote[]> {
    try {
      const notes = await getNoteList(typeId)
      notesMap.value[typeId] = notes
      return notes
    } catch (e) {
      console.error('loadTypeNotes failed:', typeId, e)
      return []
    }
  }

  async function searchNotes(query: string): Promise<GenericNote[]> {
    try {
      return await searchV2(query)
    } catch (e) {
      console.error('search failed:', e)
      return []
    }
  }

  return {
    notesMap,
    stats,
    loading,
    totalNotes,
    loadNotes,
    loadStats,
    loadTypeNotes,
    searchNotes,
  }
})
