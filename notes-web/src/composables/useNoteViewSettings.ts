import { ref, watch } from 'vue'

export type OpenMode = 'page' | 'drawer' | 'dialog'

export interface NoteViewSettings {
  openMode: OpenMode
  drawerSize?: string
  dialogWidth?: string
}

const STORAGE_KEY = 'note-view-settings'

const defaultSettings: NoteViewSettings = {
  openMode: 'drawer',
  drawerSize: '50%',
  dialogWidth: '60%',
}

// 全局设置状态
const settings = ref<NoteViewSettings>({ ...defaultSettings })

// 从 localStorage 加载
function loadSettings() {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const parsed = JSON.parse(saved)
      settings.value = { ...defaultSettings, ...parsed }
    }
  } catch (e) {
    console.error('Failed to load note view settings:', e)
  }
}

// 保存到 localStorage
function saveSettings() {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings.value))
  } catch (e) {
    console.error('Failed to save note view settings:', e)
  }
}

// 监听变化自动保存
watch(settings, saveSettings, { deep: true })

// 初始化加载
loadSettings()

export function useNoteViewSettings() {
  function setOpenMode(mode: OpenMode) {
    settings.value.openMode = mode
  }

  function setDrawerSize(size: string) {
    settings.value.drawerSize = size
  }

  function setDialogWidth(width: string) {
    settings.value.dialogWidth = width
  }

  function resetSettings() {
    settings.value = { ...defaultSettings }
  }

  return {
    settings,
    setOpenMode,
    setDrawerSize,
    setDialogWidth,
    resetSettings,
  }
}
