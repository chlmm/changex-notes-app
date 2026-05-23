/**
 * Schema 类型定义（对应 Go backend note-schema.yaml 解析后的 JSON）
 * Go struct 字段首字母大写，JSON key 也是大写首字母
 */

// Schema 分类
export interface CategoryDef {
  ID: string
  Label: string
  Icon: string
  Desc: string
}

// 扫描路径定义
export interface PathDef {
  Scan: string
  Depth: number
  Exclude: string[] | null
}

// 子文件定义
export interface SubFileDef {
  Name: string
  Label: string
  Pattern: string
}

// 字段渲染提示
export interface FieldHint {
  Label: string
  UI: string          // text | number | link | star | select | tags | badges | date | image | quote | textarea
  Options: string[] | null
  Min: number
  Max: number
}

// UI 组件类型定义
export interface UITypeDef {
  Component: string
  Desc: string
  Props: Record<string, unknown> | null
}

// 笔记类型定义
export interface TypeDef {
  ID: string
  Label: string
  Icon: string
  Category: string
  Path: PathDef
  SubFiles: SubFileDef[] | null
  FieldHints: Record<string, FieldHint> | null
  ListColumns: string[]
  MultiEntry: boolean
  FrontmatterOptional: boolean
  ResourceDirs: string[] | null
}

// Schema 元信息
export interface MetaInfo {
  Version: string
  Name: string
  Description: string
}

// 完整 Schema
export interface NoteSchema {
  Meta: MetaInfo
  Categories: CategoryDef[]
  Types: TypeDef[]
  UITypes: Record<string, UITypeDef>
}
