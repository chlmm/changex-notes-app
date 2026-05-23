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

// GetDetail 获取单条笔记详情（返回指向 slice 元素的指针，可安全修改缓存）
func (s *GenericNoteService) GetDetail(typeID, id string) (*models.GenericNote, error) {
	notes, ok := s.notes[typeID]
	if !ok {
		return nil, fmt.Errorf("unknown note type: %s", typeID)
	}
	for i := range notes {
		if notes[i].ID == id {
			return &notes[i], nil
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

// UpdateField 更新笔记的单个字段，写回 .md 文件
func (s *GenericNoteService) UpdateField(typeID, id, field, value string) error {
	note, err := s.GetDetail(typeID, id)
	if err != nil {
		return err
	}

	filePath := filepath.Join(s.cfg.NotesPath, note.Path)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	// 判断是否为多条目
	isMulti := strings.Contains(id, "#")
	idx := -1
	if isMulti {
		parts := strings.SplitN(id, "#", 2)
		if len(parts) == 2 {
			fmt.Sscanf(parts[1], "%d", &idx)
		}
	}

	newContent, err := s.updateFieldInContent(typeID, string(content), idx, field, value)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	// 更新内存缓存
	note.Fields[field] = value
	return nil
}

// updateFieldInContent 在不破坏 Markdown 内容的前提下修改 YAML frontmatter 中的字段
func (s *GenericNoteService) updateFieldInContent(typeID, content string, idx int, field, value string) (string, error) {
	if idx >= 0 {
		return s.updateMultiEntryField(content, idx, field, value)
	}
	return s.updateSingleEntryField(content, field, value)
}

func (s *GenericNoteService) updateSingleEntryField(content, field, value string) (string, error) {
	content = strings.TrimSpace(content)
	if !strings.HasPrefix(content, "---") {
		return "", fmt.Errorf("no frontmatter delimiter")
	}

	rest := content[3:]
	endIdx := strings.Index(rest, "\n---")
	if endIdx == -1 {
		return "", fmt.Errorf("no closing frontmatter delimiter")
	}

	yamlStr := rest[:endIdx]
	body := strings.TrimSpace(rest[endIdx+4:])

	var fm map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &fm); err != nil {
		return "", fmt.Errorf("yaml parse: %w", err)
	}

	// 更新字段（支持 nested: parent.child）
	if err := setNestedFieldWithUnnesting(fm, field, value); err != nil {
		return "", err
	}

	// 序列化回 YAML
	newYAML, err := yaml.Marshal(fm)
	if err != nil {
		return "", fmt.Errorf("yaml marshal: %w", err)
	}

	sb := strings.Builder{}
	sb.WriteString("---\n")
	sb.WriteString(strings.TrimSpace(string(newYAML)))
	sb.WriteString("\n---")
	if body != "" {
		sb.WriteString("\n\n")
		sb.WriteString(body)
	}
	sb.WriteString("\n")
	return sb.String(), nil
}

func (s *GenericNoteService) updateMultiEntryField(content string, idx int, field, value string) (string, error) {
	// 解析全部 entries（和 loadMultiEntry 用同样的方法），过滤掉 body block
	type entryData struct {
		fm   map[string]interface{}
		body string
	}
	var entries []entryData

	blocks := splitFrontmatterBlocks(content)
	for _, block := range blocks {
		fm, body, err := parseYAMLFrontmatter(block)
		if err != nil || isEmptyEntry(fm) {
			continue
		}
		entries = append(entries, entryData{fm, body})
	}

	if idx < 0 || idx >= len(entries) {
		return "", fmt.Errorf("entry index %d out of range (total %d entries)", idx, len(entries))
	}

	// 更新目标 entry 的字段
	if err := setNestedFieldWithUnnesting(entries[idx].fm, field, value); err != nil {
		return "", err
	}

	// 全部重新序列化（跳过 body block —— 恢复干净的文件结构）
	var sb strings.Builder
	for i, e := range entries {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		unnested := unnestMap(e.fm)
		yamlBytes, yerr := yaml.Marshal(unnested)
		if yerr != nil {
			return "", fmt.Errorf("entry %d marshal: %w", i, yerr)
		}
		sb.WriteString("---\n")
		sb.WriteString(strings.TrimSpace(string(yamlBytes)))
		sb.WriteString("\n---")
		if e.body != "" {
			sb.WriteString("\n")
			sb.WriteString(e.body)
		}
	}
	sb.WriteString("\n")
	return sb.String(), nil
}

// setNestedFieldWithUnnesting 设置嵌套字段（支持 "parent.child" 表示法）
// 同时处理"反向嵌套"：如果存在扁平化后的 key（如 "source.title"），
// 则还原为嵌套结构写入
func setNestedFieldWithUnnesting(fm map[string]interface{}, field, value string) error {
	// 先尝试作为已经扁平化的 key 还原
	if strings.Contains(field, ".") {
		parts := strings.SplitN(field, ".", 2)
		parent := parts[0]
		child := parts[1]
		if nested, ok := fm[parent]; ok {
			if nestedMap, ok := nested.(map[string]interface{}); ok {
				nestedMap[child] = value
				return nil
			}
		}
		// 如果父级不存在，创建嵌套结构
		newNested := map[string]interface{}{child: value}
		fm[parent] = newNested
		return nil
	}
	fm[field] = value
	return nil
}

// =============================================================
// CRUD 操作
// =============================================================

// CreateNote 创建新笔记
func (s *GenericNoteService) CreateNote(typeID string, fields map[string]interface{}, content string) (*models.GenericNote, error) {
	td := s.schema.GetType(typeID)
	if td == nil {
		return nil, fmt.Errorf("unknown note type: %s", typeID)
	}

	// 确定文件路径：优先从 field 推断，否则自动生成
	filePath := s.resolveCreatePath(td, fields)
	fullPath := filepath.Join(s.cfg.NotesPath, filePath)

	// 检查文件是否已存在
	if _, err := os.Stat(fullPath); err == nil {
		if td.MultiEntry {
			// 多条目：追加到现有文件
			return s.appendMultiEntry(td, fullPath, filePath, fields, content)
		}
		return nil, fmt.Errorf("文件已存在: %s", filePath)
	}

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, fmt.Errorf("create dir: %w", err)
	}

	// 将展开的扁平 key 还原为嵌套 YAML
	yamlStr := s.marshalFields(fields)
	sb := strings.Builder{}
	sb.WriteString("---\n")
	sb.WriteString(yamlStr)
	sb.WriteString("---\n")
	if content != "" {
		sb.WriteString("\n")
		sb.WriteString(content)
		sb.WriteString("\n")
	}

	if err := os.WriteFile(fullPath, []byte(sb.String()), 0644); err != nil {
		return nil, fmt.Errorf("write file: %w", err)
	}

	note := &models.GenericNote{
		TypeID:  typeID,
		ID:      filePath,
		Fields:  fields,
		Content: content,
		Path:    filePath,
	}

	// 追加到内存缓存
	s.notes[typeID] = append(s.notes[typeID], *note)
	return note, nil
}

// UpdateNote 更新笔记
func (s *GenericNoteService) UpdateNote(typeID, id string, fields map[string]interface{}, content string) error {
	note, err := s.GetDetail(typeID, id)
	if err != nil {
		return err
	}

	filePath := filepath.Join(s.cfg.NotesPath, note.Path)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	isMulti := strings.Contains(id, "#")
	idx := -1
	if isMulti {
		if parts := strings.SplitN(id, "#", 2); len(parts) == 2 {
			fmt.Sscanf(parts[1], "%d", &idx)
		}
	}

	newContent, err := s.replaceEntryContent(string(fileContent), idx, fields, content)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	// 更新内存
	for k, v := range fields {
		note.Fields[k] = v
	}
	note.Content = content
	return nil
}

// DeleteNote 删除笔记
func (s *GenericNoteService) DeleteNote(typeID, id string) error {
	note, err := s.GetDetail(typeID, id)
	if err != nil {
		return err
	}

	filePath := filepath.Join(s.cfg.NotesPath, note.Path)

	isMulti := strings.Contains(id, "#")
	if !isMulti {
		// 单条目：直接删除文件
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("remove file: %w", err)
		}
	} else {
		// 多条目：删除指定 entry
		parts := strings.SplitN(id, "#", 2)
		idx := 0
		fmt.Sscanf(parts[1], "%d", &idx)

		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		newContent, err := s.removeMultiEntry(string(fileContent), idx)
		if err != nil {
			return err
		}

		if strings.TrimSpace(newContent) == "" {
			// 全部删完了，删除文件
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("remove empty file: %w", err)
			}
		} else {
			if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
				return fmt.Errorf("write file: %w", err)
			}
		}
	}

	// 从内存缓存中移除
	notes := s.notes[typeID]
	for i, n := range notes {
		if n.ID == id {
			s.notes[typeID] = append(notes[:i], notes[i+1:]...)
			break
		}
	}
	return nil
}

// =============================================================
// CRUD 辅助函数
// =============================================================

// resolveCreatePath 解析创建路径
func (s *GenericNoteService) resolveCreatePath(td *schema.TypeDef, fields map[string]interface{}) string {
	// 尝试从 field 中获取目录/文件名
	scan := td.Path.Scan

	// 如果有 title，用 title 作为文件名
	title := ""
	if t, ok := fields["title"]; ok {
		if ts, ok := t.(string); ok {
			title = ts
		}
	}
	if title == "" {
		if t, ok := fields["source.title"]; ok {
			title = fmt.Sprintf("%v", t)
		}
	}
	if title == "" {
		title = "untitled"
	}
	fileName := sanitizeFileName(title) + ".md"

	// 从 scan 模式中取目录
	dir := ""
	if strings.Contains(scan, "**") {
		parts := strings.SplitN(scan, "**", 2)
		dir = parts[0]
	} else {
		dir = filepath.Dir(scan)
	}

	return filepath.Join(dir, fileName)
}

func sanitizeFileName(name string) string {
	replacer := strings.NewReplacer(
		"/", "-", "\\", "-", ":", "-", "*", "-", "?", "-",
		"\"", "-", "<", "-", ">", "-", "|", "-",
	)
	return replacer.Replace(name)
}

// marshalFields 将扁平字段还原为嵌套 YAML 结构
func (s *GenericNoteService) marshalFields(fields map[string]interface{}) string {
	unnested := unnestMap(fields)
	yamlBytes, _ := yaml.Marshal(unnested)
	return strings.TrimSpace(string(yamlBytes))
}

// unnestMap 将 "parent.child" key 还原为嵌套 map
func unnestMap(flat map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range flat {
		if strings.Contains(k, ".") {
			parts := strings.SplitN(k, ".", 2)
			parent := parts[0]
			child := parts[1]
			if existing, ok := result[parent]; ok {
				if existingMap, ok := existing.(map[string]interface{}); ok {
					existingMap[child] = v
				}
			} else {
				result[parent] = map[string]interface{}{child: v}
			}
		} else {
			result[k] = v
		}
	}
	return result
}

// appendMultiEntry 追加条目到多条目文件
func (s *GenericNoteService) appendMultiEntry(td *schema.TypeDef, fullPath, relPath string, fields map[string]interface{}, content string) (*models.GenericNote, error) {
	existing, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	existingStr := strings.TrimSpace(string(existing))

	// 计数现有条目，确定新索引
	blocks := splitFrontmatterBlocks(existingStr)
	newIdx := len(blocks)

	yamlStr := s.marshalFields(fields)

	sb := strings.Builder{}
	sb.WriteString(existingStr)
	sb.WriteString("\n\n---\n")
	sb.WriteString(yamlStr)
	sb.WriteString("\n---")
	if content != "" {
		sb.WriteString("\n")
		sb.WriteString(content)
	}
	sb.WriteString("\n")

	if err := os.WriteFile(fullPath, []byte(sb.String()), 0644); err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s#%d", relPath, newIdx)
	note := &models.GenericNote{
		TypeID:  td.ID,
		ID:      id,
		Fields:  fields,
		Content: content,
		Path:    relPath,
	}
	s.notes[td.ID] = append(s.notes[td.ID], *note)
	return note, nil
}

// replaceEntryContent 替换条目的完整 frontmatter + content
func (s *GenericNoteService) replaceEntryContent(content string, idx int, fields map[string]interface{}, body string) (string, error) {
	if idx >= 0 {
		return s.replaceMultiEntryContent(content, idx, fields, body)
	}
	return s.replaceSingleEntryContent(content, fields, body)
}

func (s *GenericNoteService) replaceSingleEntryContent(content string, fields map[string]interface{}, body string) (string, error) {
	yamlStr := s.marshalFields(fields)

	sb := strings.Builder{}
	sb.WriteString("---\n")
	sb.WriteString(yamlStr)
	sb.WriteString("\n---")
	if body != "" {
		sb.WriteString("\n\n")
		sb.WriteString(body)
	}
	sb.WriteString("\n")
	return sb.String(), nil
}

func (s *GenericNoteService) replaceMultiEntryContent(content string, idx int, fields map[string]interface{}, body string) (string, error) {
	blocks := splitFrontmatterBlocks(content)
	if idx < 0 || idx >= len(blocks) {
		return "", fmt.Errorf("entry index %d out of range", idx)
	}

	newBlock, err := s.replaceSingleEntryContent(blocks[idx], fields, body)
	if err != nil {
		return "", err
	}
	blocks[idx] = strings.TrimSpace(newBlock)
	return strings.Join(blocks, "\n---\n") + "\n", nil
}

// removeMultiEntry 删除多条目文件中的指定条目
func (s *GenericNoteService) removeMultiEntry(content string, idx int) (string, error) {
	blocks := splitFrontmatterBlocks(content)
	if idx < 0 || idx >= len(blocks) {
		return "", fmt.Errorf("entry index %d out of range", idx)
	}
	blocks = append(blocks[:idx], blocks[idx+1:]...)
	if len(blocks) == 0 {
		return "", nil
	}
	return strings.Join(blocks, "\n---\n") + "\n", nil
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
