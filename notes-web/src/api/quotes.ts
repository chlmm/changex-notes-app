/**
 * 金句笔记相关 API
 */

import { fetchAPI } from './index'
import type { QuoteNote } from '@/types/note'

// API 返回的原始数据类型
export interface QuoteListResponse {
  id: string
  title: string
  author?: string
  tags?: string[]
}

export interface QuoteDetailResponse {
  quote: string
  author?: string
  source?: string
  tags?: string[]
  comment?: string
  path: string
  id: string
}

/**
 * 获取金句列表
 */
export async function getQuotes(author?: string, tag?: string): Promise<QuoteListResponse[]> {
  const params = new URLSearchParams()
  if (author) params.append('author', author)
  if (tag) params.append('tag', tag)
  const query = params.toString() ? `?${params.toString()}` : ''
  return fetchAPI<QuoteListResponse[]>(`/quotes${query}`)
}

/**
 * 获取金句详情
 */
export async function getQuoteDetail(id: string): Promise<QuoteDetailResponse> {
  return fetchAPI<QuoteDetailResponse>(`/quotes/detail?id=${encodeURIComponent(id)}`)
}

/**
 * 获取所有作者列表
 */
export async function getQuoteAuthors(): Promise<string[]> {
  return fetchAPI<string[]>('/quotes/authors')
}

/**
 * 获取所有标签列表
 */
export async function getQuoteTags(): Promise<string[]> {
  return fetchAPI<string[]>('/quotes/tags')
}

/**
 * 将列表 API 响应转换为前端使用的 QuoteNote 类型
 */
export function transformQuoteItem(item: QuoteListResponse): QuoteNote {
  return {
    quote: item.title,
    author: item.author || '',
    tags: item.tags || [],
    path: item.id,
    id: item.id,
  }
}

/**
 * 将详情 API 响应转换为前端使用的 QuoteNote 类型
 */
export function transformQuoteDetail(data: QuoteDetailResponse): QuoteNote {
  return {
    quote: data.quote,
    author: data.author || '',
    source: data.source || '',
    tags: data.tags || [],
    comment: data.comment || '',
    path: data.path,
    id: data.id,
  }
}
