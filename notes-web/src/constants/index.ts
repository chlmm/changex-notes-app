/**
 * 全局常量
 */

// API 基础路径
export const API_BASE = '/api'

// 笔记类型
export const NOTE_TYPES = {
  BOOK: 'book',
  VIDEO: 'video',
  KNOWLEDGE: 'knowledge',
  SKILL: 'skill',
  PROBLEM: 'problem',
  INDEX: 'index',
} as const

// 存储键名
export const STORAGE_KEYS = {
  THEME: 'notes-theme',
  SIDEBAR_COLLAPSED: 'notes-sidebar-collapsed',
} as const
