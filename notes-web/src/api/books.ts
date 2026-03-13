/**
 * 书籍相关 API
 */

import { fetchAPI, putAPI } from './index'
import type { BookNote } from '@/types/note'
import { NoteStatus } from '@/types/note'

// API 返回的原始数据类型
export interface BookListResponse {
  title: string
  author?: string
  subType?: string
  tags?: string[]
  status?: NoteStatus
  created?: string
  id: string
}

export interface BookDetailResponse {
  title: string
  author?: string
  category?: string
  tags?: string[]
  status?: NoteStatus
  created?: string
  content?: string
  path: string
  type: 'novel' | 'non-fiction'
  quotes?: string
  reviews?: string
}

/**
 * 获取书籍列表
 */
export async function getBooks(): Promise<BookListResponse[]> {
  return fetchAPI<BookListResponse[]>('/books')
}

/**
 * 获取书籍详情
 */
export async function getBookDetail(type: string, name: string): Promise<BookDetailResponse> {
  return fetchAPI<BookDetailResponse>(`/books/${type}/${encodeURIComponent(name)}`)
}

/**
 * 更新书籍状态
 */
export async function updateBookStatus(type: string, title: string, status: string): Promise<boolean> {
  try {
    await putAPI(`/books/${type}/${encodeURIComponent(title)}/status`, { status })
    return true
  } catch {
    return false
  }
}

/**
 * 将 API 响应转换为前端使用的 BookNote 类型
 */
export function transformBookItem(item: BookListResponse): BookNote {
  return {
    title: item.title,
    author: item.author || '',
    category: item.subType === 'novel' ? '小说' : '非虚构',
    tags: item.tags || [],
    status: item.status || NoteStatus.NotStarted,
    created: item.created || '',
    content: '',
    path: item.id,
    type: item.subType === 'novel' ? 'novel' : 'non-fiction',
  }
}

/**
 * 将详情 API 响应转换为前端使用的 BookNote 类型
 */
export function transformBookDetail(data: BookDetailResponse): BookNote {
  return {
    title: data.title,
    author: data.author || '',
    category: data.category || '',
    tags: data.tags || [],
    status: data.status || NoteStatus.NotStarted,
    created: data.created || '',
    content: data.content || '',
    path: data.path,
    type: data.type,
    quotes: data.quotes,
    reviews: data.reviews,
  }
}
