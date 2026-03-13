/**
 * 知识库相关 API
 */

import { fetchAPI } from './index'
import type { KnowledgeNote } from '@/types/note'

// API 返回的原始数据类型
export interface KnowledgeListResponse {
  title: string
  id: string
  subType?: string
}

export interface KnowledgeDetailResponse {
  title: string
  category?: string
  topic?: string
  content?: string
  path: string
}

// 知识树结构类型
export type KnowledgeTree = Record<string, Record<string, string[]>>

/**
 * 获取知识库列表
 */
export async function getKnowledgeList(): Promise<KnowledgeListResponse[]> {
  return fetchAPI<KnowledgeListResponse[]>('/knowledge')
}

/**
 * 获取知识树结构
 */
export async function getKnowledgeTree(): Promise<KnowledgeTree> {
  return fetchAPI<KnowledgeTree>('/knowledge/tree')
}

/**
 * 通过路径获取知识点详情
 */
export async function getKnowledgeDetailByPath(path: string): Promise<KnowledgeDetailResponse> {
  return fetchAPI<KnowledgeDetailResponse>(`/knowledge/detail?path=${encodeURIComponent(path)}`)
}

/**
 * 通过分类/主题/标题获取知识点详情
 */
export async function getKnowledgeDetailByQuery(
  category: string,
  topic: string,
  title: string
): Promise<KnowledgeDetailResponse> {
  const params = new URLSearchParams({
    category,
    topic,
    title,
  })
  return fetchAPI<KnowledgeDetailResponse>(`/knowledge/detail?${params}`)
}

/**
 * 将列表 API 响应转换为前端使用的 KnowledgeNote 类型
 */
export function transformKnowledgeItem(item: KnowledgeListResponse): KnowledgeNote {
  return {
    title: item.title,
    category: item.subType?.split('/')[0] || '',
    topic: item.subType?.split('/')[1] || '',
    content: '',
    path: item.id,
  }
}

/**
 * 将详情 API 响应转换为前端使用的 KnowledgeNote 类型
 */
export function transformKnowledgeDetail(data: KnowledgeDetailResponse): KnowledgeNote {
  return {
    title: data.title,
    category: data.category || '',
    topic: data.topic || '',
    content: data.content || '',
    path: data.path,
  }
}
