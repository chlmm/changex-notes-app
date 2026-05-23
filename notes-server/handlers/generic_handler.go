package handlers

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"notes-server/schema"
	"notes-server/services"
)

type GenericHandler struct {
	service *services.GenericNoteService
	schema  *schema.NoteSchema
}

func NewGenericHandler(service *services.GenericNoteService, s *schema.NoteSchema) *GenericHandler {
	return &GenericHandler{service: service, schema: s}
}

// GetSchema 获取完整 schema
func (h *GenericHandler) GetSchema(c *gin.Context) {
	c.JSON(http.StatusOK, h.schema)
}

// GetStats 通用统计
func (h *GenericHandler) GetStats(c *gin.Context) {
	stats := h.service.GetStats()
	c.JSON(http.StatusOK, stats)
}

// GetNoteList 获取某类型笔记列表
// GET /api/v2/notes/:typeId
func (h *GenericHandler) GetNoteList(c *gin.Context) {
	typeID := c.Param("typeId")
	notes, err := h.service.GetList(typeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notes)
}

// GetNoteDetail 获取笔记详情
// GET /api/v2/notes/:typeId/*id（*id 捕获包含 / 的路径）
func (h *GenericHandler) GetNoteDetail(c *gin.Context) {
	typeID := c.Param("typeId")
	id := c.Param("id")

	// *id 包含前导 /，需要去掉
	id = strings.TrimPrefix(id, "/")

	// URL-decode：前端用 encodeURIComponent 编码了 # 和 / 等特殊字符
	if decoded, err := url.QueryUnescape(id); err == nil {
		id = decoded
	}

	note, err := h.service.GetDetail(typeID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, note)
}

// GetFilters 获取某类型某字段的筛选选项
// GET /api/v2/notes/:typeId/filters/:field
func (h *GenericHandler) GetFilters(c *gin.Context) {
	typeID := c.Param("typeId")
	field := c.Param("field")

	values := h.service.GetFilterValues(typeID, field)
	c.JSON(http.StatusOK, gin.H{
		"field":  field,
		"values": values,
	})
}

// Search 搜索笔记
// GET /api/v2/search?q=xxx
func (h *GenericHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}
	results := h.service.Search(query)
	c.JSON(http.StatusOK, results)
}

// UpdateField 更新笔记单个字段
// PUT /api/v2/notes/:typeId/field
// Body: {"id": "xxx", "field": "status", "value": "done"}
func (h *GenericHandler) UpdateField(c *gin.Context) {
	typeID := c.Param("typeId")

	var req struct {
		ID    string `json:"id" binding:"required"`
		Field string `json:"field" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateField(typeID, req.ID, req.Field, req.Value); err != nil {
		log.Printf("UpdateField error: type=%s id=%s field=%s value=%s err=%v", typeID, req.ID, req.Field, req.Value, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// CreateNote 创建笔记
// POST /api/v2/notes/:typeId
func (h *GenericHandler) CreateNote(c *gin.Context) {
	typeID := c.Param("typeId")

	var req struct {
		Fields  map[string]interface{} `json:"fields" binding:"required"`
		Content string                 `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.service.CreateNote(typeID, req.Fields, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, note)
}

// UpdateNote 更新笔记
// PUT /api/v2/notes/:typeId/*id
func (h *GenericHandler) UpdateNote(c *gin.Context) {
	typeID := c.Param("typeId")
	id := c.Param("id")

	// *id 包含前导 /，需要去掉；URL-decode 处理 # 等特殊字符
	id = strings.TrimPrefix(id, "/")
	if decoded, err := url.QueryUnescape(id); err == nil {
		id = decoded
	}

	var req struct {
		Fields  map[string]interface{} `json:"fields"`
		Content string                 `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateNote(typeID, id, req.Fields, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// DeleteNote 删除笔记
// DELETE /api/v2/notes/:typeId/*id
func (h *GenericHandler) DeleteNote(c *gin.Context) {
	typeID := c.Param("typeId")
	id := c.Param("id")

	// *id 包含前导 /，需要去掉；URL-decode 处理 # 等特殊字符
	id = strings.TrimPrefix(id, "/")
	if decoded, err := url.QueryUnescape(id); err == nil {
		id = decoded
	}

	if err := h.service.DeleteNote(typeID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}
