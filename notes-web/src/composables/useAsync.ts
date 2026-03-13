/**
 * 异步状态管理 composable
 */

import { ref, type Ref } from 'vue'

export interface AsyncState<T> {
  data: Ref<T | null>
  loading: Ref<boolean>
  error: Ref<string | null>
}

export function useAsync<T>(initialValue: T | null = null) {
  const data = ref<T | null>(initialValue) as Ref<T | null>
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function execute(asyncFn: () => Promise<T>): Promise<T | null> {
    loading.value = true
    error.value = null
    
    try {
      const result = await asyncFn()
      data.value = result
      return result
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Unknown error'
      return null
    } finally {
      loading.value = false
    }
  }

  function reset() {
    data.value = initialValue
    loading.value = false
    error.value = null
  }

  return {
    data,
    loading,
    error,
    execute,
    reset,
  }
}
