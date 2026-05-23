package services

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"notes-server/config"
	"notes-server/models"
	"notes-server/schema"
)

// GenericNoteService schema-driven 通用笔记服务
type GenericNoteService struct {
	cfg    *config.Config
	schema *schema.NoteSchema
	notes  map[string][]models.GenericNote // key: typeId
}

func NewGenericNoteService(cfg *config.Config, s *schema.NoteSchema) *GenericNoteService {
	return &GenericNoteService{
		cfg:    cfg,
		schema: s,
		notes:  make(map[string][]models.GenericNote),
	}
}

// LoadAll 按 schema 加载所有类型的笔记
func (s *GenericNoteService) LoadAll() error {
	for _, t := range s.schema.Types {
		notes, err := s.loadType(&t)
		if err != nil {
			log.Printf("GenericNoteService: failed to load %s: %v", t.ID, err)
			continue
		}
		s.notes[t.ID] = notes
		log.Printf("GenericNoteService: loaded %d notes for type %s", len(notes), t.ID)
	}
	return nil
}

// GetList 获取某类型笔记列表（支持过滤）
func (s *GenericNoteService) GetList(typeID string) ([]models.GenericNote, error) {
	notes, ok := s.notes[typeID]
	if !ok {
		return nil, fmt.Errorf("unknown note type: %s", typeID)
	}
	return notes, nil
}

// GetDetail 获取单条笔记详情
func (s *GenericNoteService) GetDetail(typeID, id string) (*models.GenericNote, error) {
	notes, ok := s.notes[typeID]
	if !ok {
		return nil, fmt.Errorf("unknown note type: %s", typeID)
	}
	for _, n := range notes {
		if n.ID == id {
			return &n, nil
		}
	}
	return nil, fmt.Errorf("note not found: %s/%s", typeID, id)
}

// GetFilterValues 获取某类型某个字段的所有去重值
func (s *GenericNoteService) GetFilterValues(typeID, field string) []string {
	notes := s.notes[typeID]
	seen := make(map[string]bool)
	var values []string
	for _, n := range notes {
		if v, ok := n.Fields[field]; ok {
			str := fmt.Sprintf("%v", v)
			if !seen[str] {
				seen[str] = true
				values = append(values, str)
			}
		}
	}
	sort.Strings(values)
	return values
}

// GetStats 通用统计（按 schema 类型汇总）
func (s *GenericNoteService) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	for _, t := range s.schema.Types {
		stats[t.ID] = len(s.notes[t.ID])
	}
	return stats
}

// Search 通用搜索
func (s *GenericNoteService) Search(query string) []models.GenericNote {
	query = strings.ToLower(query)
	var results []models.GenericNote
	for _, notes := range s.notes {
		for _, n := range notes {
			if matchGenericNote(n, query) {
				results = append(results, n)
			}
		}
	}
	return results
}

// =============================================================
// 内部加载逻辑
// =============================================================

// loadType 按 TypeDef 加载笔记
func (s *GenericNoteService) loadType(td *schema.TypeDef) ([]models.GenericNote, error) {
	files, err := s.scanFiles(td)
	if err != nil {
		return nil, fmt.Errorf("scan files: %w", err)
	}

	var notes []models.GenericNote
	for _, f := range files {
		loaded, err := s.loadFile(td, f)
		if err != nil {
			log.Printf("GenericNoteService: skip %s: %v", f, err)
			continue
		}
		notes = append(notes, loaded...)
	}
	return notes, nil
}

// scanFiles 按 TypeDef.Path 扫描匹配的文件
func (s *GenericNoteService) scanFiles(td *schema.TypeDef) ([]string, error) {
	scan := td.Path.Scan
	baseDir := s.cfg.NotesPath

	// 处理 ** 递归模式
	if strings.Contains(scan, "**") {
		return s.scanRecursive(baseDir, scan, td.Path.Exclude)
	}

	// 普通 glob 模式
	pattern := filepath.Join(baseDir, scan)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// scanRecursive 递归扫描（处理 ** 模式）
func (s *GenericNoteService) scanRecursive(baseDir, scan string, excludes []string) ([]string, error) {
	// 将 scan 转换为路径：把 ** 替换为空，获取扫描根目录
	parts := strings.SplitN(scan, "**", 2)
	scanRoot := filepath.Join(baseDir, parts[0])
	suffix := ""
	if len(parts) > 1 {
		suffix = strings.TrimPrefix(parts[1], "/")
	}

	var matches []string
	err := filepath.Walk(scanRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		// 检查 exclude
		rel, _ := filepath.Rel(baseDir, path)
		for _, ex := range excludes {
			if matched, _ := filepath.Match(ex, rel); matched {
				return nil
			}
			if matched, _ := filepath.Match(ex, filepath.Base(path)); matched {
				return nil
			}
		}
		// 匹配后缀（文件名模式）
		if suffix != "" {
			if matched, _ := filepath.Match(suffix, filepath.Base(path)); !matched {
				return nil
			}
		}
		matches = append(matches, path)
		return nil
	})
	return matches, err
}

// loadFile 加载单个文件，返回一个或多个 GenericNote
func (s *GenericNoteService) loadFile(td *schema.TypeDef, filePath string) ([]models.GenericNote, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fileInfo, _ := os.Stat(filePath)
	modTime := ""
	if fileInfo != nil {
		modTime = fileInfo.ModTime().Format(time.RFC3339)
	}

	if td.MultiEntry {
		return s.loadMultiEntry(td, filePath, string(content), modTime)
	}
	return s.loadSingleEntry(td, filePath, string(content), modTime)
}

// loadSingleEntry 加载单条目文件
func (s *GenericNoteService) loadSingleEntry(td *schema.TypeDef, filePath, content, modTime string) ([]models.GenericNote, error) {
	fm, body, err := parseYAMLFrontmatter(content)
	if err != nil {
		if td.FrontmatterOptional {
			// 无 frontmatter，从文件路径/标题推断
			title := strings.TrimSuffix(filepath.Base(filePath), ".md")
			fm = map[string]interface{}{"title": title}
			body = content
		} else {
			return nil, fmt.Errorf("parse frontmatter: %w", err)
		}
	}

	relPath, _ := filepath.Rel(s.cfg.NotesPath, filePath)
	id := relPath

	note := models.GenericNote{
		TypeID:    td.ID,
		ID:        id,
		Fields:    fm,
		Content:   body,
		Path:      relPath,
		UpdatedAt: modTime,
	}

	// 加载子文件
	s.loadSubFiles(td, filePath, &note)

	return []models.GenericNote{note}, nil
}

// loadMultiEntry 加载多条目文件（--- 分隔的多个 YAML 块）
func (s *GenericNoteService) loadMultiEntry(td *schema.TypeDef, filePath, content, modTime string) ([]models.GenericNote, error) {
	relPath, _ := filepath.Rel(s.cfg.NotesPath, filePath)
	blocks := splitFrontmatterBlocks(content)

	var notes []models.GenericNote
	for i, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		fm, body, err := parseYAMLFrontmatter(block)
		if err != nil {
			log.Printf("GenericNoteService: skip entry %d in %s: %v", i, filePath, err)
			continue
		}

		// 跳过空条目
		if isEmptyEntry(fm) {
			continue
		}

		id := fmt.Sprintf("%s#%d", relPath, i)
		note := models.GenericNote{
			TypeID:    td.ID,
			ID:        id,
			Fields:    fm,
			Content:   body,
			Path:      relPath,
			UpdatedAt: modTime,
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// loadSubFiles 加载附属文件（quotes.md, reviews.md 等）
func (s *GenericNoteService) loadSubFiles(td *schema.TypeDef, indexPath string, note *models.GenericNote) {
	if len(td.SubFiles) == 0 {
		return
	}
	dir := filepath.Dir(indexPath)
	if note.SubFiles == nil {
		note.SubFiles = make(map[string]string)
	}
	for _, sf := range td.SubFiles {
		subPath := filepath.Join(dir, sf.Pattern)
		if subContent, err := os.ReadFile(subPath); err == nil {
			note.SubFiles[sf.Name] = string(subContent)
		}
	}
}

// =============================================================
// YAML Frontmatter 解析
// =============================================================

// parseYAMLFrontmatter 解析 YAML frontmatter + body
// 支持标准格式：
//
//	---
//	key: value
//	source:
//	  url: xxx
//	---
//	body content
func parseYAMLFrontmatter(content string) (fields map[string]interface{}, body string, err error) {
	content = strings.TrimSpace(content)

	if !strings.HasPrefix(content, "---") {
		return nil, content, fmt.Errorf("no frontmatter delimiter")
	}

	// 找到第二个 ---
	rest := content[3:] // skip first ---
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		// 可能是单行 ---key: value--- 或无结束分隔符
		return nil, content, fmt.Errorf("no closing frontmatter delimiter")
	}

	yamlStr := strings.TrimSpace(rest[:idx])
	body = strings.TrimSpace(rest[idx+4:]) // skip \n---

	// 解析 YAML
	var raw map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &raw); err != nil {
		return nil, body, fmt.Errorf("yaml parse: %w", err)
	}

	// 扁平化嵌套对象（source.url → "source.url"）
	fields = flattenMap(raw, "")
	return fields, body, nil
}

// flattenMap 将嵌套 map 扁平化为 "parent.child" 格式
func flattenMap(m map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			for nk, nv := range flattenMap(val, key) {
				result[nk] = nv
			}
		default:
			result[key] = v
		}
	}
	return result
}

// splitFrontmatterBlocks 按 --- 分割多条目文件
// 条目之间由 \n---\n 分隔，每个条目各自有 YAML frontmatter
func splitFrontmatterBlocks(content string) []string {
	content = strings.TrimSpace(content)
	sep := "\n---\n"
	parts := strings.Split(content, sep)

	var blocks []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// 确保被 --- 包裹
		if !strings.HasPrefix(part, "---") {
			part = "---\n" + part
		}
		if !strings.HasSuffix(part, "---") {
			part = part + "\n---"
		}
		blocks = append(blocks, part)
	}
	return blocks
}

// isEmptyEntry 检查是否是空条目（扁平化后字段已是平铺 key）
func isEmptyEntry(fm map[string]interface{}) bool {
	// 名言类型
	if q, ok := fm["quote"]; ok {
		if s, ok := q.(string); ok && s != "" {
			return false
		}
	}
	// GitHub / anime / game / movie 类型（扁平化后是 source.url / source.title 等）
	if url, ok := fm["source.url"]; ok {
		if s, ok := url.(string); ok && s != "" {
			return false
		}
	}
	if title, ok := fm["source.title"]; ok {
		if s, ok := title.(string); ok && s != "" {
			return false
		}
	}
	// 有 title 字段
	if t, ok := fm["title"]; ok {
		if s, ok := t.(string); ok && s != "" {
			return false
		}
	}
	return true
}

// matchGenericNote 搜索匹配
func matchGenericNote(note models.GenericNote, query string) bool {
	// 搜索所有字段值
	for _, v := range note.Fields {
		if strings.Contains(strings.ToLower(fmt.Sprintf("%v", v)), query) {
			return true
		}
	}
	// 搜索正文
	if strings.Contains(strings.ToLower(note.Content), query) {
		return true
	}
	return false
}
