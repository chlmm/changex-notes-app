package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"notes-server/models"
	"notes-server/services"
)

type NoteHandler struct {
	service *services.NoteService
}

func NewNoteHandler(service *services.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

// GetStats 获取统计信息
func (h *NoteHandler) GetStats(c *gin.Context) {
	stats := h.service.GetStats()
	c.JSON(http.StatusOK, stats)
}

// GetNotes 获取笔记列表
func (h *NoteHandler) GetNotes(c *gin.Context) {
	noteType := c.Query("type")
	status := c.Query("status")

	notes := h.service.GetAllNotes(noteType, status)
	c.JSON(http.StatusOK, notes)
}

// GetBooks 获取书籍列表
func (h *NoteHandler) GetBooks(c *gin.Context) {
	subType := c.Query("subType") // novel 或 non-fiction
	status := c.Query("status")

	books := h.service.GetBooks(subType, status)
	c.JSON(http.StatusOK, books)
}

// GetBookDetail 获取书籍详情
func (h *NoteHandler) GetBookDetail(c *gin.Context) {
	bookType := c.Param("type") // novel 或 non-fiction
	title := c.Param("title")

	// URL 解码
	decodedTitle, err := url.PathUnescape(title)
	if err != nil {
		decodedTitle = title
	}

	book, err := h.service.GetBookDetail(bookType, decodedTitle)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// GetVideos 获取视频列表
func (h *NoteHandler) GetVideos(c *gin.Context) {
	status := c.Query("status")
	videos := h.service.GetVideos(status)
	c.JSON(http.StatusOK, videos)
}

// GetVideoDetail 获取视频详情
func (h *NoteHandler) GetVideoDetail(c *gin.Context) {
	videoID := c.Param("id")

	video, err := h.service.GetVideoDetail(videoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
}

// Search 搜索笔记
func (h *NoteHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	results := h.service.Search(query)
	c.JSON(http.StatusOK, results)
}

// UpdateStatus 更新笔记状态
func (h *NoteHandler) UpdateStatus(c *gin.Context) {
	noteType := c.Param("type")
	id := c.Param("id")

	// URL 解码
	decodedID, err := url.PathUnescape(id)
	if err != nil {
		decodedID = id
	}

	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateNoteStatus(noteType, decodedID, req.Status); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

// GetRecentNotes 获取最近更新的笔记
func (h *NoteHandler) GetRecentNotes(c *gin.Context) {
	notes := h.service.GetRecentNotes(5)
	c.JSON(http.StatusOK, notes)
}

// GetKnowledge 获取知识库列表
func (h *NoteHandler) GetKnowledge(c *gin.Context) {
	category := c.Query("category")
	topic := c.Query("topic")

	knowledge := h.service.GetKnowledge(category, topic)
	c.JSON(http.StatusOK, knowledge)
}

// GetKnowledgeTree 获取知识库树形结构
func (h *NoteHandler) GetKnowledgeTree(c *gin.Context) {
	tree := h.service.GetKnowledgeTree()
	c.JSON(http.StatusOK, tree)
}

// GetKnowledgeDetail 获取知识库详情
func (h *NoteHandler) GetKnowledgeDetail(c *gin.Context) {
	path := c.Query("path")
	category := c.Query("category")
	topic := c.Query("topic")
	title := c.Query("title")

	var kn *models.KnowledgeNote
	var err error

	// 优先使用 path 查询
	if path != "" {
		decodedPath, err := url.PathUnescape(path)
		if err != nil {
			decodedPath = path
		}
		kn, err = h.service.GetKnowledgeDetail(decodedPath)
	} else if category != "" && topic != "" && title != "" {
		// 使用 category/topic/title 查询
		kn, err = h.service.GetKnowledgeDetailByQuery(category, topic, title)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path or category/topic/title parameters are required"})
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kn)
}

// GetSkills 获取技能列表
func (h *NoteHandler) GetSkills(c *gin.Context) {
	category := c.Query("category")
	skillType := c.Query("type")

	skills := h.service.GetSkills(category, skillType)
	c.JSON(http.StatusOK, skills)
}

// GetSkillDetail 获取技能详情
func (h *NoteHandler) GetSkillDetail(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter is required"})
		return
	}

	decodedPath, err := url.PathUnescape(path)
	if err != nil {
		decodedPath = path
	}

	skill, err := h.service.GetSkillDetail(decodedPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, skill)
}

// GetProblems 获取问题解决列表
func (h *NoteHandler) GetProblems(c *gin.Context) {
	problemType := c.Query("type")

	problems := h.service.GetProblems(problemType)
	c.JSON(http.StatusOK, problems)
}

// GetProblemDetail 获取问题解决详情
func (h *NoteHandler) GetProblemDetail(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter is required"})
		return
	}

	decodedPath, err := url.PathUnescape(path)
	if err != nil {
		decodedPath = path
	}

	problem, err := h.service.GetProblemDetail(decodedPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, problem)
}

// GetIndexNotes 获取索引列表
func (h *NoteHandler) GetIndexNotes(c *gin.Context) {
	notes := h.service.GetIndexNotes()
	c.JSON(http.StatusOK, notes)
}

// GetIndexDetail 获取索引详情
func (h *NoteHandler) GetIndexDetail(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path parameter is required"})
		return
	}

	decodedPath, err := url.PathUnescape(path)
	if err != nil {
		decodedPath = path
	}

	idx, err := h.service.GetIndexDetail(decodedPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, idx)
}

// GetQuotes 获取金句列表
func (h *NoteHandler) GetQuotes(c *gin.Context) {
	author := c.Query("author")
	tag := c.Query("tag")

	quotes := h.service.GetQuotes(author, tag)
	c.JSON(http.StatusOK, quotes)
}

// GetQuoteDetail 获取金句详情
func (h *NoteHandler) GetQuoteDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	decodedID, err := url.PathUnescape(id)
	if err != nil {
		decodedID = id
	}

	quote, err := h.service.GetQuoteDetail(decodedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quote)
}

// GetQuoteAuthors 获取金句作者列表
func (h *NoteHandler) GetQuoteAuthors(c *gin.Context) {
	authors := h.service.GetQuoteAuthors()
	c.JSON(http.StatusOK, authors)
}

// GetQuoteTags 获取金句标签列表
func (h *NoteHandler) GetQuoteTags(c *gin.Context) {
	tags := h.service.GetQuoteTags()
	c.JSON(http.StatusOK, tags)
}

// GetGitHubRepos 获取 GitHub 项目列表
func (h *NoteHandler) GetGitHubRepos(c *gin.Context) {
	language := c.Query("language")
	tag := c.Query("tag")

	repos := h.service.GetGitHubRepos(language, tag)
	c.JSON(http.StatusOK, repos)
}

// GetGitHubRepoDetail 获取 GitHub 项目详情
func (h *NoteHandler) GetGitHubRepoDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	decodedID, err := url.PathUnescape(id)
	if err != nil {
		decodedID = id
	}

	repo, err := h.service.GetGitHubRepoDetail(decodedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, repo)
}

// GetGitHubLanguages 获取 GitHub 项目语言列表
func (h *NoteHandler) GetGitHubLanguages(c *gin.Context) {
	languages := h.service.GetGitHubLanguages()
	c.JSON(http.StatusOK, languages)
}

// GetGitHubTags 获取 GitHub 项目标签列表
func (h *NoteHandler) GetGitHubTags(c *gin.Context) {
	tags := h.service.GetGitHubTags()
	c.JSON(http.StatusOK, tags)
}

// ==================== 动漫相关 Handler ====================

// GetAnime 获取动漫列表
func (h *NoteHandler) GetAnime(c *gin.Context) {
	status := c.Query("status")
	tag := c.Query("tag")

	anime := h.service.GetAnime(status, tag)
	c.JSON(http.StatusOK, anime)
}

// GetAnimeDetail 获取动漫详情
func (h *NoteHandler) GetAnimeDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	decodedID, err := url.PathUnescape(id)
	if err != nil {
		decodedID = id
	}

	anime, err := h.service.GetAnimeDetail(decodedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, anime)
}

// GetAnimeStatuses 获取动漫状态列表
func (h *NoteHandler) GetAnimeStatuses(c *gin.Context) {
	statuses := h.service.GetAnimeStatuses()
	c.JSON(http.StatusOK, statuses)
}

// GetAnimeTags 获取动漫标签列表
func (h *NoteHandler) GetAnimeTags(c *gin.Context) {
	tags := h.service.GetAnimeTags()
	c.JSON(http.StatusOK, tags)
}

// ==================== 电影相关 Handler ====================

// GetMovies 获取电影列表
func (h *NoteHandler) GetMovies(c *gin.Context) {
	genre := c.Query("genre")
	tag := c.Query("tag")

	movies := h.service.GetMovies(genre, tag)
	c.JSON(http.StatusOK, movies)
}

// GetMovieDetail 获取电影详情
func (h *NoteHandler) GetMovieDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	decodedID, err := url.PathUnescape(id)
	if err != nil {
		decodedID = id
	}

	movie, err := h.service.GetMovieDetail(decodedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// GetMovieGenres 获取电影类型列表
func (h *NoteHandler) GetMovieGenres(c *gin.Context) {
	genres := h.service.GetMovieGenres()
	c.JSON(http.StatusOK, genres)
}

// GetMovieTags 获取电影标签列表
func (h *NoteHandler) GetMovieTags(c *gin.Context) {
	tags := h.service.GetMovieTags()
	c.JSON(http.StatusOK, tags)
}

// ==================== 游戏相关 Handler ====================

// GetGames 获取游戏列表
func (h *NoteHandler) GetGames(c *gin.Context) {
	platform := c.Query("platform")
	status := c.Query("status")
	tag := c.Query("tag")

	games := h.service.GetGames(platform, status, tag)
	c.JSON(http.StatusOK, games)
}

// GetGameDetail 获取游戏详情
func (h *NoteHandler) GetGameDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	decodedID, err := url.PathUnescape(id)
	if err != nil {
		decodedID = id
	}

	game, err := h.service.GetGameDetail(decodedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, game)
}

// GetGamePlatforms 获取游戏平台列表
func (h *NoteHandler) GetGamePlatforms(c *gin.Context) {
	platforms := h.service.GetGamePlatforms()
	c.JSON(http.StatusOK, platforms)
}

// GetGameStatuses 获取游戏状态列表
func (h *NoteHandler) GetGameStatuses(c *gin.Context) {
	statuses := h.service.GetGameStatuses()
	c.JSON(http.StatusOK, statuses)
}

// GetGameTags 获取游戏标签列表
func (h *NoteHandler) GetGameTags(c *gin.Context) {
	tags := h.service.GetGameTags()
	c.JSON(http.StatusOK, tags)
}
