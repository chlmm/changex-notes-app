import { ref, watch } from 'vue'
import { STORAGE_KEYS } from '@/constants'

export type ThemeMode = 'light' | 'dark' | 'system'

// 全局单例状态
const theme = ref<ThemeMode>(loadTheme())
const isDark = ref(false)

function loadTheme(): ThemeMode {
  try {
    const saved = localStorage.getItem(STORAGE_KEYS.THEME)
    if (saved === 'dark' || saved === 'light' || saved === 'system') return saved
  } catch { /* ignore */ }
  return 'light'
}

function applyTheme(mode: ThemeMode) {
  const prefersDark =
    mode === 'dark' || (mode === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)
  isDark.value = prefersDark
  document.documentElement.classList.toggle('dark', prefersDark)
}

// 初始化
applyTheme(theme.value)

// 持久化
watch(theme, (mode) => {
  localStorage.setItem(STORAGE_KEYS.THEME, mode)
  applyTheme(mode)
})

// 监听系统主题变化
if (window.matchMedia) {
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    if (theme.value === 'system') applyTheme('system')
  })
}

const themeLabel: Record<ThemeMode, string> = {
  light: '浅色',
  dark: '深色',
  system: '跟随系统',
}

export function useTheme() {
  function setTheme(mode: ThemeMode) {
    theme.value = mode
  }

  function toggleTheme() {
    const cycle: ThemeMode[] = ['light', 'dark', 'system']
    const idx = cycle.indexOf(theme.value)
    theme.value = cycle[(idx + 1) % cycle.length]
  }

  return {
    theme,
    isDark,
    themeLabel,
    setTheme,
    toggleTheme,
  }
}
