/**
 * 通用数据转换工具函数
 */

/**
 * 安全获取嵌套属性
 */
export function safeGet<T>(obj: unknown, path: string, defaultValue: T): T {
  const keys = path.split('.')
  let result: unknown = obj
  
  for (const key of keys) {
    if (result && typeof result === 'object' && key in result) {
      result = (result as Record<string, unknown>)[key]
    } else {
      return defaultValue
    }
  }
  
  return result as T
}

/**
 * 解析分类路径 (e.g., "Math/Algebra" -> ["Math", "Algebra"])
 */
export function parseCategoryPath(path: string): [string, string] {
  const parts = path.split('/')
  return [parts[0] || '', parts[1] || '']
}

/**
 * 构建分类路径
 */
export function buildCategoryPath(category: string, topic: string): string {
  return `${category}/${topic}`
}

/**
 * URL 安全编码
 */
export function safeEncode(str: string): string {
  return encodeURIComponent(str)
}

/**
 * 批量转换数组
 */
export function transformArray<T, R>(
  items: T[],
  transformer: (item: T) => R
): R[] {
  return items.map(transformer)
}

/**
 * 过滤无效项目
 */
export function filterValidItems<T>(
  items: (T | null | undefined)[]
): T[] {
  return items.filter((item): item is T => item != null)
}
