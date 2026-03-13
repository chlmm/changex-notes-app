/**
 * 搜索相关 API
 */

import { fetchAPI } from './index'
import type { SearchResult } from '@/types/note'

/**
 * 搜索笔记
 */
export async function searchNotes(query: string): Promise<SearchResult[]> {
  return fetchAPI<SearchResult[]>(`/search?q=${encodeURIComponent(query)}`)
}

/**
 * 获取统计数据
 */
export async function getStats(): Promise<{
  totalBooks: number
  totalVideos: number
  readingBooks: number
  completedBooks: number
  totalKnowledge: number
  totalSkills: number
  totalProblems: number
  totalIndex: number
}> {
  return fetchAPI('/stats')
}
