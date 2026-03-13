package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"notes-server/config"
	"notes-server/handlers"
	"notes-server/services"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS 中间件
	r.Use(corsMiddleware())

	// 初始化服务
	noteService := services.NewNoteService(cfg)
	if err := noteService.LoadNotes(); err != nil {
		panic("Failed to load notes: " + err.Error())
	}

	// 初始化 handler
	noteHandler := handlers.NewNoteHandler(noteService)

	// API 路由
	api := r.Group("/api")
	{
		// 统计
		api.GET("/stats", noteHandler.GetStats)

		// 笔记列表
		api.GET("/notes", noteHandler.GetNotes)

		// 书籍
		api.GET("/books", noteHandler.GetBooks)
		api.GET("/books/:type/:title", noteHandler.GetBookDetail)
		api.PUT("/books/:type/:title/status", noteHandler.UpdateStatus)

		// 视频
		api.GET("/videos", noteHandler.GetVideos)
		api.GET("/videos/:id", noteHandler.GetVideoDetail)
		api.PUT("/videos/:id/status", noteHandler.UpdateStatus)

		// 知识库
		api.GET("/knowledge", noteHandler.GetKnowledge)
		api.GET("/knowledge/tree", noteHandler.GetKnowledgeTree)
		api.GET("/knowledge/detail", noteHandler.GetKnowledgeDetail)

		// 技能
		api.GET("/skills", noteHandler.GetSkills)
		api.GET("/skills/detail", noteHandler.GetSkillDetail)

		// 问题解决
		api.GET("/problems", noteHandler.GetProblems)
		api.GET("/problems/detail", noteHandler.GetProblemDetail)

		// 索引
		api.GET("/index", noteHandler.GetIndexNotes)
		api.GET("/index/detail", noteHandler.GetIndexDetail)

		// 金句
		api.GET("/quotes", noteHandler.GetQuotes)
		api.GET("/quotes/detail", noteHandler.GetQuoteDetail)
		api.GET("/quotes/authors", noteHandler.GetQuoteAuthors)
		api.GET("/quotes/tags", noteHandler.GetQuoteTags)

		// GitHub 项目收藏
		api.GET("/github", noteHandler.GetGitHubRepos)
		api.GET("/github/detail", noteHandler.GetGitHubRepoDetail)
		api.GET("/github/languages", noteHandler.GetGitHubLanguages)
		api.GET("/github/tags", noteHandler.GetGitHubTags)

		// 动漫收藏
		api.GET("/anime", noteHandler.GetAnime)
		api.GET("/anime/detail", noteHandler.GetAnimeDetail)
		api.GET("/anime/statuses", noteHandler.GetAnimeStatuses)
		api.GET("/anime/tags", noteHandler.GetAnimeTags)

		// 电影收藏
		api.GET("/movies", noteHandler.GetMovies)
		api.GET("/movies/detail", noteHandler.GetMovieDetail)
		api.GET("/movies/genres", noteHandler.GetMovieGenres)
		api.GET("/movies/tags", noteHandler.GetMovieTags)

		// 游戏收藏
		api.GET("/games", noteHandler.GetGames)
		api.GET("/games/detail", noteHandler.GetGameDetail)
		api.GET("/games/platforms", noteHandler.GetGamePlatforms)
		api.GET("/games/statuses", noteHandler.GetGameStatuses)
		api.GET("/games/tags", noteHandler.GetGameTags)

		// 搜索
		api.GET("/search", noteHandler.Search)

		// 最近更新
		api.GET("/recent", noteHandler.GetRecentNotes)
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

// corsMiddleware CORS 中间件
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
