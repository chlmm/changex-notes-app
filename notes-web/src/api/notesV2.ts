/**
 * v2 通用 API 模块（schema-driven）
 */
import { fetchAPI } from './index'
import type { GenericNote } from '@/types/note'

// 获取统计
export function getStats(): Promise<Record<string, number>> {
  return fetchAPI('/v2/stats')
}

// 获取某类型笔记列表
export function getNoteList(typeId: string): Promise<GenericNote[]> {
  return fetchAPI(`/v2/notes/${typeId}`)
}

// 获取笔记详情（id 可能含 /，后端 /*id 通配捕获）
export function getNoteDetail(typeId: string, id: string): Promise<GenericNote> {
  return fetchAPI(`/v2/notes/${typeId}/detail/${id}`)
}

// 获取字段筛选值
export function getFilters(typeId: string, field: string): Promise<{ field: string; values: string[] }> {
  return fetchAPI(`/v2/notes/${typeId}/filters/${field}`)
}

// 搜索
export function searchNotes(query: string): Promise<GenericNote[]> {
  return fetchAPI(`/v2/search?q=${encodeURIComponent(query)}`)
}
