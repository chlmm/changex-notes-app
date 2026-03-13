/**
 * 技能相关 API
 */

import { fetchAPI } from './index'
import type { SkillNote } from '@/types/note'

// API 返回的原始数据类型
export interface SkillListResponse {
  title: string
  id: string
  subType?: string
}

export interface SkillDetailResponse {
  name: string
  description?: string
  category?: string
  type?: string
  content?: string
  path: string
}

/**
 * 获取技能列表
 */
export async function getSkills(): Promise<SkillListResponse[]> {
  return fetchAPI<SkillListResponse[]>('/skills')
}

/**
 * 获取技能详情
 */
export async function getSkillDetail(path: string): Promise<SkillDetailResponse> {
  return fetchAPI<SkillDetailResponse>(`/skills/detail?path=${encodeURIComponent(path)}`)
}

/**
 * 将列表 API 响应转换为前端使用的 SkillNote 类型
 */
export function transformSkillItem(item: SkillListResponse): SkillNote {
  return {
    title: item.title,
    name: item.title,
    description: '',
    category: item.subType || '',
    type: 'SKILL',
    content: '',
    path: item.id,
  }
}

/**
 * 将详情 API 响应转换为前端使用的 SkillNote 类型
 */
export function transformSkillDetail(data: SkillDetailResponse): SkillNote {
  return {
    title: data.name || '',
    name: data.name,
    description: data.description || '',
    category: data.category || '',
    type: data.type || 'SKILL',
    content: data.content || '',
    path: data.path,
  }
}
