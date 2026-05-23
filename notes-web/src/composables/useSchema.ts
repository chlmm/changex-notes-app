import { ref, computed } from 'vue'
import { fetchAPI } from '@/api/index'
import type { NoteSchema, TypeDef, CategoryDef, FieldHint } from '@/types/schema'

// 全局单例
const schema = ref<NoteSchema | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)
let loaded = false

export function useSchema() {
  async function loadSchema(force = false): Promise<NoteSchema | null> {
    if (loaded && !force) return schema.value

    loading.value = true
    error.value = null

    try {
      const data = await fetchAPI<NoteSchema>('/schema')
      schema.value = data
      loaded = true
      return data
    } catch (e) {
      error.value = 'Failed to load schema'
      console.error('useSchema: load failed', e)
      return null
    } finally {
      loading.value = false
    }
  }

  // 按 ID 获取类型
  function getType(typeId: string): TypeDef | undefined {
    return schema.value?.Types.find(t => t.ID === typeId)
  }

  // 按分类获取类型列表
  function getTypesByCategory(categoryId: string): TypeDef[] {
    return schema.value?.Types.filter(t => t.Category === categoryId) ?? []
  }

  // 所有分类
  const categories = computed<CategoryDef[]>(() => schema.value?.Categories ?? [])

  // 获取类型字段列表（用于渲染）
  function getTypeFields(typeId: string): { key: string; hint: FieldHint }[] {
    const type = getType(typeId)
    if (!type?.FieldHints) return []

    return Object.entries(type.FieldHints).map(([key, hint]) => ({
      key,
      hint,
    }))
  }

  return {
    schema,
    loading,
    error,
    categories,
    loadSchema,
    getType,
    getTypesByCategory,
    getTypeFields,
  }
}

// 导出类型辅助函数（不依赖 Vue 响应式）
export function getFieldUI(hint: FieldHint | undefined): string {
  return hint?.UI ?? 'text'
}

export function getFieldLabel(key: string, hint?: FieldHint): string {
  if (hint?.Label) return hint.Label
  // 自动转换: camelCase → 自然语言, source.title → Title
  const parts = key.split('.')
  const last = parts[parts.length - 1]
  return last
    .replace(/([A-Z])/g, ' $1')
    .replace(/^./, s => s.toUpperCase())
    .replace(/_/g, ' ')
}

export function getFieldOptions(hint?: FieldHint): string[] {
  return hint?.Options ?? []
}
