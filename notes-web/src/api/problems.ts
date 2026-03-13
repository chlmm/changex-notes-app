/**
 * 问题相关 API
 */

import { fetchAPI } from './index'
import type { ProblemNote } from '@/types/note'

// API 返回的原始数据类型
export interface ProblemListResponse {
  title: string
  id: string
  subType?: string
}

export interface ProblemDetailResponse {
  title: string
  type?: string
  content?: string
  path: string
}

/**
 * 获取问题列表
 */
export async function getProblems(): Promise<ProblemListResponse[]> {
  return fetchAPI<ProblemListResponse[]>('/problems')
}

/**
 * 获取问题详情
 */
export async function getProblemDetail(path: string): Promise<ProblemDetailResponse> {
  return fetchAPI<ProblemDetailResponse>(`/problems/detail?path=${encodeURIComponent(path)}`)
}

/**
 * 将列表 API 响应转换为前端使用的 ProblemNote 类型
 */
export function transformProblemItem(item: ProblemListResponse): ProblemNote {
  return {
    title: item.title,
    type: item.subType || '',
    content: '',
    path: item.id,
  }
}

/**
 * 将详情 API 响应转换为前端使用的 ProblemNote 类型
 */
export function transformProblemDetail(data: ProblemDetailResponse): ProblemNote {
  return {
    title: data.title,
    type: data.type || '',
    content: data.content || '',
    path: data.path,
  }
}
