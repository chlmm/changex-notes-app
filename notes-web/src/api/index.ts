/**
 * API 基础请求封装
 * 统一处理请求、响应和错误
 */

const API_BASE = '/api'

export interface ApiResponse<T> {
  data: T | null
  error: string | null
}

/**
 * 基础 GET 请求
 */
export async function fetchAPI<T>(endpoint: string): Promise<T> {
  const res = await fetch(`${API_BASE}${endpoint}`)
  if (!res.ok) {
    throw new Error(`API error: ${res.status}`)
  }
  return res.json()
}

/**
 * 基础 PUT 请求
 */
export async function putAPI<T>(endpoint: string, body: unknown): Promise<T> {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    throw new Error(`API error: ${res.status}`)
  }
  return res.json()
}

/**
 * 基础 POST 请求
 */
export async function postAPI<T>(endpoint: string, body: unknown): Promise<T> {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!res.ok) {
    throw new Error(`API error: ${res.status}`)
  }
  return res.json()
}

/**
 * 检查响应是否成功（用于不需要返回值的请求）
 */
export async function checkSuccess(endpoint: string, options?: RequestInit): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}${endpoint}`, options)
    return res.ok
  } catch {
    return false
  }
}
