/**
 * v2 通用 API 模块（schema-driven）
 */
import { fetchAPI, putAPI, postAPI, deleteAPI } from './index'
import type { GenericNote } from '@/types/note'

// 获取统计
export function getStats(): Promise<Record<string, number>> {
  return fetchAPI('/v2/stats')
}

// 获取某类型笔记列表
export function getNoteList(typeId: string): Promise<GenericNote[]> {
  return fetchAPI(`/v2/notes/${typeId}`)
}

// 获取笔记详情（id 可能含 / 和 #，encodeURIComponent 防止被浏览器当 fragment 截断）
export function getNoteDetail(typeId: string, id: string): Promise<GenericNote> {
  return fetchAPI(`/v2/notes/${typeId}/detail/${encodeURIComponent(id)}`)
}

// 获取字段筛选值
export function getFilters(typeId: string, field: string): Promise<{ field: string; values: string[] }> {
  return fetchAPI(`/v2/notes/${typeId}/filters/${field}`)
}

// 搜索
export function searchNotes(query: string): Promise<GenericNote[]> {
  return fetchAPI(`/v2/search?q=${encodeURIComponent(query)}`)
}

// 更新笔记单个字段（看板拖拽等场景）
export function updateNoteField(
  typeId: string,
  id: string,
  field: string,
  value: string,
): Promise<{ ok: boolean }> {
  return putAPI(`/v2/field/${typeId}`, { id, field, value })
}

// 创建笔记
export function createNote(
  typeId: string,
  fields: Record<string, unknown>,
  content?: string,
): Promise<GenericNote> {
  return postAPI(`/v2/notes/${typeId}`, { fields, content })
}

// 更新笔记
export function updateNote(
  typeId: string,
  id: string,
  fields: Record<string, unknown>,
  content?: string,
): Promise<{ ok: boolean }> {
  return putAPI(`/v2/notes/${typeId}/${encodeURIComponent(id)}`, { fields, content })
}

// 删除笔记
export function deleteNote(
  typeId: string,
  id: string,
): Promise<{ ok: boolean }> {
  return deleteAPI(`/v2/notes/${typeId}/${encodeURIComponent(id)}`)
}
