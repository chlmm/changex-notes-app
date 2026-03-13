/**
 * 游戏收藏相关 API
 */

import { fetchAPI } from './index'
import type { GameNote } from '@/types/note'

// API 返回的原始数据类型
export interface GameListResponse {
  id: string
  title: string
  subType?: string  // platform
  status?: string
  tags?: string[]
}

export interface GameDetailResponse {
  url?: string
  title: string
  platform?: string
  developer?: string
  year?: number
  status?: string
  rating?: number
  tags?: string[]
  comment?: string
  path: string
  id: string
}

/**
 * 获取游戏列表
 */
export async function getGames(platform?: string, status?: string, tag?: string): Promise<GameListResponse[]> {
  const params = new URLSearchParams()
  if (platform) params.append('platform', platform)
  if (status) params.append('status', status)
  if (tag) params.append('tag', tag)
  const query = params.toString() ? `?${params.toString()}` : ''
  return fetchAPI<GameListResponse[]>(`/games${query}`)
}

/**
 * 获取游戏详情
 */
export async function getGameDetail(id: string): Promise<GameDetailResponse> {
  return fetchAPI<GameDetailResponse>(`/games/detail?id=${encodeURIComponent(id)}`)
}

/**
 * 获取所有平台列表
 */
export async function getGamePlatforms(): Promise<string[]> {
  return fetchAPI<string[]>('/games/platforms')
}

/**
 * 获取所有状态列表
 */
export async function getGameStatuses(): Promise<string[]> {
  return fetchAPI<string[]>('/games/statuses')
}

/**
 * 获取所有标签列表
 */
export async function getGameTags(): Promise<string[]> {
  return fetchAPI<string[]>('/games/tags')
}

/**
 * 将列表 API 响应转换为前端使用的 GameNote 类型
 */
export function transformGameItem(item: GameListResponse): GameNote {
  return {
    title: item.title,
    platform: item.subType,
    status: item.status,
    tags: item.tags || [],
    path: item.id,
    id: item.id,
  }
}

/**
 * 将详情 API 响应转换为前端使用的 GameNote 类型
 */
export function transformGameDetail(data: GameDetailResponse): GameNote {
  return {
    url: data.url || '',
    title: data.title,
    platform: data.platform || '',
    developer: data.developer || '',
    year: data.year,
    status: data.status || '',
    rating: data.rating,
    tags: data.tags || [],
    comment: data.comment || '',
    path: data.path,
    id: data.id,
  }
}
