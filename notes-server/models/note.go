package models

// GenericNote 通用笔记结构，字段由 schema 动态定义
type GenericNote struct {
	TypeID    string                 `json:"typeId"`
	ID        string                 `json:"id"`
	Fields    map[string]interface{} `json:"fields"`
	Content   string                 `json:"content,omitempty"`
	SubFiles  map[string]string      `json:"subFiles,omitempty"`
	Path      string                 `json:"path"`
	UpdatedAt string                 `json:"updatedAt,omitempty"`
}
