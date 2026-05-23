package handlers

import (
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
