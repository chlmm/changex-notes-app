/**
 * 通用 Note 类型（schema-driven）
 */

// 通用笔记（字段由 schema 动态定义）
export interface GenericNote {
  typeId: string
  id: string
  fields: Record<string, unknown>
  content?: string
  subFiles?: Record<string, string>
  path: string
  updatedAt?: string
}
