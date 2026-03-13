/**
 * 索引笔记相关 API
 */

import { fetchAPI } from './index'
import type { IndexNote } from '@/types/note'

// API 返回的原始数据类型
export interface IndexListResponse {
  title: string
  id: string
  subType?: string
}

export interface IndexDetailResponse {
  title: string
  url?: string
  type?: string
  content?: string
  path: string
}

/**
 * 获取索引列表
 */
export async function getIndexList(): Promise<IndexListResponse[]> {
  return fetchAPI<IndexListResponse[]>('/index')
}

/**
 * 获取索引详情
 */
export async function getIndexDetail(path: string): Promise<IndexDetailResponse> {
  return fetchAPI<IndexDetailResponse>(`/index/detail?path=${encodeURIComponent(path)}`)
}

/**
 * 将列表 API 响应转换为前端使用的 IndexNote 类型
 */
export function transformIndexItem(item: IndexListResponse): IndexNote {
  return {
    title: item.title,
    type: item.subType || '',
    content: '',
    path: item.id,
  }
}

/**
 * 将详情 API 响应转换为前端使用的 IndexNote 类型
 */
export function transformIndexDetail(data: IndexDetailResponse): IndexNote {
  return {
    title: data.title,
    url: data.url || '',
    type: data.type || '',
    content: data.content || '',
    path: data.path,
  }
}
