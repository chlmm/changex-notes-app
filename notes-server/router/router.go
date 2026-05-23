package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"notes-server/config"
	"notes-server/handlers"
	"notes-server/schema"
	"notes-server/services"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS 中间件
	r.Use(corsMiddleware())

	// 加载 Schema
	noteSchema, err := schema.Load(cfg.SchemaPath)
	if err != nil {
		log.Fatalf("Failed to load schema: %v", err)
	}

	// 初始化通用服务
	genericService := services.NewGenericNoteService(cfg, noteSchema)
	if err := genericService.LoadAll(); err != nil {
		log.Printf("Warning: failed to load notes: %v", err)
	}
	genericHandler := handlers.NewGenericHandler(genericService, noteSchema)

	// API 路由
	api := r.Group("/api")
	{
		// Schema
		api.GET("/schema", genericHandler.GetSchema)

		// 通用 API
		api.GET("/v2/stats", genericHandler.GetStats)
		api.GET("/v2/notes/:typeId", genericHandler.GetNoteList)
		api.GET("/v2/notes/:typeId/filters/:field", genericHandler.GetFilters)
		api.GET("/v2/notes/:typeId/detail/*id", genericHandler.GetNoteDetail)
		api.GET("/v2/search", genericHandler.Search)
	}

	// 静态文件服务（前端）
	r.Static("/assets", cfg.StaticPath+"/assets")
	r.StaticFile("/", cfg.StaticPath+"/index.html")
	r.StaticFile("/favicon.ico", cfg.StaticPath+"/favicon.ico")

	// SPA 回退
	r.NoRoute(func(c *gin.Context) {
		c.File(cfg.StaticPath + "/index.html")
	})

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
