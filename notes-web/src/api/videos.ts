/**
 * 视频相关 API
 */

import { fetchAPI, putAPI } from './index'
import type { VideoNote } from '@/types/note'
import { NoteStatus } from '@/types/note'

// API 返回的原始数据类型
export interface VideoListResponse {
  title: string
  id: string
  tags?: string[]
  status?: NoteStatus
  created?: string
}

export interface VideoDetailResponse {
  title: string
  platform?: string
  videoId: string
  url?: string
  tags?: string[]
  status?: NoteStatus
  created?: string
  content?: string
  path: string
}

/**
 * 获取视频列表
 */
export async function getVideos(): Promise<VideoListResponse[]> {
  return fetchAPI<VideoListResponse[]>('/videos')
}

/**
 * 获取视频详情
 */
export async function getVideoDetail(id: string): Promise<VideoDetailResponse> {
  return fetchAPI<VideoDetailResponse>(`/videos/${id}`)
}

/**
 * 更新视频状态
 */
export async function updateVideoStatus(id: string, status: string): Promise<boolean> {
  try {
    await putAPI(`/videos/${id}/status`, { status })
    return true
  } catch {
    return false
  }
}

/**
 * 将列表 API 响应转换为前端使用的 VideoNote 类型
 */
export function transformVideoItem(item: VideoListResponse): VideoNote {
  return {
    title: item.title,
    platform: 'bilibili',
    videoId: item.id,
    url: '',
    tags: item.tags || [],
    status: item.status || NoteStatus.NotStarted,
    created: item.created || '',
    content: '',
    path: item.id,
  }
}

/**
 * 将详情 API 响应转换为前端使用的 VideoNote 类型
 */
export function transformVideoDetail(data: VideoDetailResponse): VideoNote {
  return {
    title: data.title,
    platform: data.platform || 'bilibili',
    videoId: data.videoId,
    url: data.url || '',
    tags: data.tags || [],
    status: data.status || NoteStatus.NotStarted,
    created: data.created || '',
    content: data.content || '',
    path: data.path,
  }
}
