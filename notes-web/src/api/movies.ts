/**
 * 电影收藏相关 API
 */

import { fetchAPI } from './index'
import type { MovieNote } from '@/types/note'

// API 返回的原始数据类型
export interface MovieListResponse {
  id: string
  title: string
  subType?: string  // director
  tags?: string[]
}

export interface MovieDetailResponse {
  url?: string
  title: string
  year?: number
  director?: string
  rating?: number
  genre?: string[]
  tags?: string[]
  comment?: string
  path: string
  id: string
}

/**
 * 获取电影列表
 */
export async function getMovies(genre?: string, tag?: string): Promise<MovieListResponse[]> {
  const params = new URLSearchParams()
  if (genre) params.append('genre', genre)
  if (tag) params.append('tag', tag)
  const query = params.toString() ? `?${params.toString()}` : ''
  return fetchAPI<MovieListResponse[]>(`/movies${query}`)
}

/**
 * 获取电影详情
 */
export async function getMovieDetail(id: string): Promise<MovieDetailResponse> {
  return fetchAPI<MovieDetailResponse>(`/movies/detail?id=${encodeURIComponent(id)}`)
}

/**
 * 获取所有类型列表
 */
export async function getMovieGenres(): Promise<string[]> {
  return fetchAPI<string[]>('/movies/genres')
}

/**
 * 获取所有标签列表
 */
export async function getMovieTags(): Promise<string[]> {
  return fetchAPI<string[]>('/movies/tags')
}

/**
 * 将列表 API 响应转换为前端使用的 MovieNote 类型
 */
export function transformMovieItem(item: MovieListResponse): MovieNote {
  return {
    title: item.title,
    director: item.subType,
    tags: item.tags || [],
    path: item.id,
    id: item.id,
  }
}

/**
 * 将详情 API 响应转换为前端使用的 MovieNote 类型
 */
export function transformMovieDetail(data: MovieDetailResponse): MovieNote {
  return {
    url: data.url || '',
    title: data.title,
    year: data.year,
    director: data.director || '',
    rating: data.rating,
    genre: data.genre || [],
    tags: data.tags || [],
    comment: data.comment || '',
    path: data.path,
    id: data.id,
  }
}
