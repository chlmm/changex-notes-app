/**
 * 动漫收藏相关 API
 */

import { fetchAPI } from './index'
import type { AnimeNote } from '@/types/note'

// API 返回的原始数据类型
export interface AnimeListResponse {
  id: string
  title: string
  subType?: string  // status
  tags?: string[]
}

export interface AnimeDetailResponse {
  url?: string
  title: string
  year?: number
  studio?: string
  episodes?: number
  status?: string
  rating?: number
  tags?: string[]
  comment?: string
  path: string
  id: string
}

/**
 * 获取动漫列表
 */
export async function getAnime(status?: string, tag?: string): Promise<AnimeListResponse[]> {
  const params = new URLSearchParams()
  if (status) params.append('status', status)
  if (tag) params.append('tag', tag)
  const query = params.toString() ? `?${params.toString()}` : ''
  return fetchAPI<AnimeListResponse[]>(`/anime${query}`)
}

/**
 * 获取动漫详情
 */
export async function getAnimeDetail(id: string): Promise<AnimeDetailResponse> {
  return fetchAPI<AnimeDetailResponse>(`/anime/detail?id=${encodeURIComponent(id)}`)
}

/**
 * 获取所有状态列表
 */
export async function getAnimeStatuses(): Promise<string[]> {
  return fetchAPI<string[]>('/anime/statuses')
}

/**
 * 获取所有标签列表
 */
export async function getAnimeTags(): Promise<string[]> {
  return fetchAPI<string[]>('/anime/tags')
}

/**
 * 将列表 API 响应转换为前端使用的 AnimeNote 类型
 */
export function transformAnimeItem(item: AnimeListResponse): AnimeNote {
  return {
    title: item.title,
    status: item.subType,
    tags: item.tags || [],
    path: item.id,
    id: item.id,
  }
}

/**
 * 将详情 API 响应转换为前端使用的 AnimeNote 类型
 */
export function transformAnimeDetail(data: AnimeDetailResponse): AnimeNote {
  return {
    url: data.url || '',
    title: data.title,
    year: data.year,
    studio: data.studio || '',
    episodes: data.episodes,
    status: data.status || '',
    rating: data.rating,
    tags: data.tags || [],
    comment: data.comment || '',
    path: data.path,
    id: data.id,
  }
}
