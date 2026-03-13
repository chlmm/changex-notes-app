/**
 * GitHub 项目收藏相关 API
 */

import { fetchAPI } from './index'
import type { GitHubRepoNote } from '@/types/note'

// API 返回的原始数据类型
export interface GitHubRepoListResponse {
  id: string
  title: string
  subType?: string  // language
  tags?: string[]
}

export interface GitHubRepoDetailResponse {
  url: string
  name: string
  title?: string
  desc?: string
  stars?: number
  forks?: number
  lang?: string
  topics?: string[]
  tags?: string[]
  comment?: string
  path: string
  id: string
}

/**
 * 获取 GitHub 项目列表
 */
export async function getGitHubRepos(language?: string, tag?: string): Promise<GitHubRepoListResponse[]> {
  const params = new URLSearchParams()
  if (language) params.append('language', language)
  if (tag) params.append('tag', tag)
  const query = params.toString() ? `?${params.toString()}` : ''
  return fetchAPI<GitHubRepoListResponse[]>(`/github${query}`)
}

/**
 * 获取 GitHub 项目详情
 */
export async function getGitHubRepoDetail(id: string): Promise<GitHubRepoDetailResponse> {
  return fetchAPI<GitHubRepoDetailResponse>(`/github/detail?id=${encodeURIComponent(id)}`)
}

/**
 * 获取所有语言列表
 */
export async function getGitHubLanguages(): Promise<string[]> {
  return fetchAPI<string[]>('/github/languages')
}

/**
 * 获取所有标签列表
 */
export async function getGitHubTags(): Promise<string[]> {
  return fetchAPI<string[]>('/github/tags')
}

/**
 * 将列表 API 响应转换为前端使用的 GitHubRepoNote 类型
 */
export function transformGitHubRepoItem(item: GitHubRepoListResponse): GitHubRepoNote {
  return {
    url: '',
    name: item.title,
    title: item.title,
    lang: item.subType,
    tags: item.tags || [],
    path: item.id,
    id: item.id,
  }
}

/**
 * 将详情 API 响应转换为前端使用的 GitHubRepoNote 类型
 */
export function transformGitHubRepoDetail(data: GitHubRepoDetailResponse): GitHubRepoNote {
  return {
    url: data.url,
    name: data.name,
    title: data.title || '',
    desc: data.desc || '',
    stars: data.stars,
    forks: data.forks,
    lang: data.lang || '',
    topics: data.topics || [],
    tags: data.tags || [],
    comment: data.comment || '',
    path: data.path,
    id: data.id,
  }
}
