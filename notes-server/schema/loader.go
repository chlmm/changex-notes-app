package schema

import (
	"os"

	"gopkg.in/yaml.v3"
)

// NoteSchema 完整笔记 schema 定义（映射 note-schema.yaml）
type NoteSchema struct {
	Meta       MetaInfo              `yaml:"meta"`
	Categories []CategoryDef         `yaml:"categories"`
	Types      []TypeDef             `yaml:"types"`
	UITypes    map[string]UITypeDef  `yaml:"ui_types"`
}

// MetaInfo schema 元信息
type MetaInfo struct {
	Version     string `yaml:"version"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// CategoryDef 顶层分类
type CategoryDef struct {
	ID    string `yaml:"id"`
	Label string `yaml:"label"`
	Icon  string `yaml:"icon"`
	Desc  string `yaml:"desc"`
}

// TypeDef 笔记类型定义
type TypeDef struct {
	ID                  string               `yaml:"id"`
	Label               string               `yaml:"label"`
	Icon                string               `yaml:"icon"`
	Category            string               `yaml:"category"`
	Path                PathDef              `yaml:"path"`
	SubFiles            []SubFileDef         `yaml:"sub_files"`
	FieldHints          map[string]FieldHint `yaml:"field_hints"`
	ListColumns         []string             `yaml:"list_columns"`
	MultiEntry          bool                 `yaml:"multi_entry"`
	FrontmatterOptional bool                 `yaml:"frontmatter_optional"`
	ResourceDirs        []string             `yaml:"resource_dirs"`
	Views               *ViewConfig          `yaml:"views"`
}

// ViewConfig 视图配置
type ViewConfig struct {
	Enabled       []string `yaml:"enabled"`
	KanbanField   string   `yaml:"kanban_field"`
	TimelineField string   `yaml:"timeline_field"`
	ChartGroupBy  string   `yaml:"chart_group_by"`
}

// PathDef 扫描路径定义
type PathDef struct {
	Scan    string   `yaml:"scan"`
	Depth   int      `yaml:"depth"`
	Exclude []string `yaml:"exclude"`
}

// SubFileDef 附属文件定义
type SubFileDef struct {
	Name    string `yaml:"name"`
	Label   string `yaml:"label"`
	Pattern string `yaml:"pattern"`
}

// FieldHint 字段渲染提示
type FieldHint struct {
	Label   string   `yaml:"label"`
	UI      string   `yaml:"ui"`
	Options []string `yaml:"options"`
	Min     int      `yaml:"min"`
	Max     int      `yaml:"max"`
}

// UITypeDef UI 组件类型定义
type UITypeDef struct {
	Component string                 `yaml:"component"`
	Desc      string                 `yaml:"desc"`
	Props     map[string]interface{} `yaml:"props,omitempty"`
}

// Load 加载并解析 note-schema.yaml
func Load(path string) (*NoteSchema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var s NoteSchema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	return &s, nil
}

// GetType 按 id 获取类型定义
func (s *NoteSchema) GetType(typeID string) *TypeDef {
	for i := range s.Types {
		if s.Types[i].ID == typeID {
			return &s.Types[i]
		}
	}
	return nil
}

// GetCategory 按 id 获取分类定义
func (s *NoteSchema) GetCategory(categoryID string) *CategoryDef {
	for i := range s.Categories {
		if s.Categories[i].ID == categoryID {
			return &s.Categories[i]
		}
	}
	return nil
}

// GetTypesByCategory 获取某分类下的所有类型
func (s *NoteSchema) GetTypesByCategory(categoryID string) []TypeDef {
	var result []TypeDef
	for _, t := range s.Types {
		if t.Category == categoryID {
			result = append(result, t)
		}
	}
	return result
}
