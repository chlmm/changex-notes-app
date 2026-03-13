package models

import "time"

type NoteType string

const (
	NoteTypeBook      NoteType = "book"
	NoteTypeVideo     NoteType = "video"
	NoteTypeKnowledge NoteType = "knowledge"
	NoteTypeSkill     NoteType = "skill"
	NoteTypeProblem   NoteType = "problem"
	NoteTypeIndex     NoteType = "index"
	NoteTypeQuote     NoteType = "quote"
	NoteTypeGitHub    NoteType = "github"
	NoteTypeAnime     NoteType = "anime"
	NoteTypeMovie     NoteType = "movie"
	NoteTypeGame      NoteType = "game"
)

type BookSubType string

const (
	BookSubTypeNovel      BookSubType = "novel"
	BookSubTypeNonFiction BookSubType = "non-fiction"
)

type NoteStatus string

const (
	StatusNotStarted  NoteStatus = "未开始"
	StatusReading     NoteStatus = "阅读中"
	StatusCompleted   NoteStatus = "已完成"
	StatusInterpreted NoteStatus = "已解读"
)

// BookNote 书籍笔记
type BookNote struct {
	Title    string     `json:"title"`
	Author   string     `json:"author"`
	Category string     `json:"category"`
	Tags     []string   `json:"tags"`
	Status   NoteStatus `json:"status"`
	Created  string     `json:"created"`
	Content  string     `json:"content,omitempty"`
	Type     BookSubType `json:"type"`
	Path     string     `json:"path"`
	Quotes   string     `json:"quotes,omitempty"`
	Reviews  string     `json:"reviews,omitempty"`
}

// VideoNote 视频笔记
type VideoNote struct {
	Title    string     `json:"title"`
	Platform string     `json:"platform"`
	VideoID  string     `json:"videoId"`
	URL      string     `json:"url"`
	Tags     []string   `json:"tags"`
	Status   NoteStatus `json:"status"`
	Created  string     `json:"created"`
	Content  string     `json:"content,omitempty"`
	Path     string     `json:"path"`
}

// KnowledgeNote 知识库笔记
type KnowledgeNote struct {
	Title    string   `json:"title"`
	Category string   `json:"category"` // Math, Physics, etc.
	Topic    string   `json:"topic"`    // Algebra, Calculus, etc.
	Content  string   `json:"content,omitempty"`
	Path     string   `json:"path"`
}

// SkillNote 技能笔记
type SkillNote struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`    // Mine/Others, skills/code-skills etc.
	Type        string   `json:"type"`        // SKILL, USER-SKILL
	Content     string   `json:"content,omitempty"`
	Path        string   `json:"path"`
}

// ProblemNote 问题解决笔记
type ProblemNote struct {
	Title      string `json:"title"`
	Type       string `json:"type"` // "问题与解决" or "需求与解决"
	Problem    string `json:"problem,omitempty"`
	Solution   string `json:"solution,omitempty"`
	Content    string `json:"content,omitempty"`
	Path       string `json:"path"`
}

// IndexNote 索引笔记
type IndexNote struct {
	Title   string `json:"title"`
	URL     string `json:"url,omitempty"`
	Type    string `json:"type,omitempty"` // 小说, 动漫, etc.
	Content string `json:"content,omitempty"`
	Path    string `json:"path"`
}

// QuoteNote 金句笔记
type QuoteNote struct {
	Quote   string   `json:"quote"`            // 金句内容
	Author  string   `json:"author,omitempty"` // 出处人
	Source  string   `json:"source,omitempty"` // 出处处
	Tags    []string `json:"tags,omitempty"`
	Comment string   `json:"comment,omitempty"` // 解读/感想
	Path    string   `json:"path"`              // 文件路径
	ID      string   `json:"id"`                // 唯一标识(路径+索引)
}

// GitHubRepoNote GitHub 项目收藏
type GitHubRepoNote struct {
	URL         string   `json:"url"`              // GitHub 项目链接
	Name        string   `json:"name"`             // owner/repo 格式
	Title       string   `json:"title,omitempty"`  // 项目标题
	Description string   `json:"desc,omitempty"`   // 项目描述
	Stars       int      `json:"stars,omitempty"`  // Star 数
	Forks       int      `json:"forks,omitempty"`  // Fork 数
	Language    string   `json:"lang,omitempty"`   // 主要语言
	Topics      []string `json:"topics,omitempty"` // GitHub Topics
	Tags        []string `json:"tags,omitempty"`   // 个人标签
	Comment     string   `json:"comment,omitempty"` // 个人备注
	Path        string   `json:"path"`             // 文件路径
	ID          string   `json:"id"`               // 唯一标识(路径+索引)
}

// AnimeNote 动漫收藏
type AnimeNote struct {
	URL       string   `json:"url"`               // 链接
	Title     string   `json:"title"`             // 标题
	Year      int      `json:"year,omitempty"`    // 年份
	Studio    string   `json:"studio,omitempty"`  // 制作公司
	Episodes  int      `json:"episodes,omitempty"` // 集数
	Status    string   `json:"status,omitempty"`  // 状态：连载/完结
	Rating    float64  `json:"rating,omitempty"`  // 评分
	Tags      []string `json:"tags,omitempty"`    // 标签
	Comment   string   `json:"comment,omitempty"` // 评论
	Path      string   `json:"path"`              // 文件路径
	ID        string   `json:"id"`                // 唯一标识(路径+索引)
}

// MovieNote 电影收藏
type MovieNote struct {
	URL       string   `json:"url"`               // 链接
	Title     string   `json:"title"`             // 标题
	Year      int      `json:"year,omitempty"`    // 年份
	Director  string   `json:"director,omitempty"` // 导演
	Rating    float64  `json:"rating,omitempty"`  // 评分
	Genre     []string `json:"genre,omitempty"`   // 类型
	Tags      []string `json:"tags,omitempty"`    // 标签
	Comment   string   `json:"comment,omitempty"` // 观后感
	Path      string   `json:"path"`              // 文件路径
	ID        string   `json:"id"`                // 唯一标识(路径+索引)
}

// GameNote 游戏收藏
type GameNote struct {
	URL       string   `json:"url"`               // 链接
	Title     string   `json:"title"`             // 标题
	Platform  string   `json:"platform,omitempty"` // 平台
	Developer string   `json:"developer,omitempty"` // 开发商
	Year      int      `json:"year,omitempty"`    // 年份
	Status    string   `json:"status,omitempty"`  // 状态：想玩/在玩/通关
	Rating    float64  `json:"rating,omitempty"`  // 评分
	Tags      []string `json:"tags,omitempty"`    // 标签
	Comment   string   `json:"comment,omitempty"` // 评论
	Path      string   `json:"path"`              // 文件路径
	ID        string   `json:"id"`                // 唯一标识(路径+索引)
}

// NoteSummary 笔记摘要（列表用）
type NoteSummary struct {
	ID       string     `json:"id"`
	Title    string     `json:"title"`
	Type     NoteType   `json:"type"`
	SubType  string     `json:"subType,omitempty"`
	Status   NoteStatus `json:"status"`
	Created  string     `json:"created"`
	Author   string     `json:"author,omitempty"`
	Tags     []string   `json:"tags"`
}

// NoteStats 统计信息
type NoteStats struct {
	TotalBooks      int `json:"totalBooks"`
	TotalVideos     int `json:"totalVideos"`
	ReadingBooks    int `json:"readingBooks"`
	CompletedBooks  int `json:"completedBooks"`
	TotalKnowledge  int `json:"totalKnowledge"`
	TotalSkills     int `json:"totalSkills"`
	TotalProblems   int `json:"totalProblems"`
	TotalIndex      int `json:"totalIndex"`
	TotalQuotes     int `json:"totalQuotes"`
	TotalGitHub     int `json:"totalGitHub"`
	TotalAnime      int `json:"totalAnime"`
	TotalMovies     int `json:"totalMovies"`
	TotalGames      int `json:"totalGames"`
}

// SearchResult 搜索结果
type SearchResult struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Type    NoteType `json:"type"`
	Snippet string   `json:"snippet"`
}

// Frontmatter 解析后的元数据
type Frontmatter struct {
	Title       string   `yaml:"title"`
	Author      string   `yaml:"author"`
	Category    string   `yaml:"category"`
	Tags        []string `yaml:"tags"`
	Status      string   `yaml:"status"`
	Created     string   `yaml:"created"`
	Platform    string   `yaml:"platform"`
	VideoID     string   `yaml:"video_id"`
	URL         string   `yaml:"url"`
	Quote       string   `yaml:"quote"`
	Source      string   `yaml:"source"`
	Comment     string   `yaml:"comment"`
	// GitHub 项目相关字段
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Stars       int      `yaml:"stars"`
	Forks       int      `yaml:"forks"`
	Language    string   `yaml:"language"`
	Topics      []string `yaml:"topics"`
	// 动漫相关字段
	Year      int     `yaml:"year"`
	Studio    string  `yaml:"studio"`
	Episodes  int     `yaml:"episodes"`
	Rating    float64 `yaml:"rating"`
	// 电影相关字段
	Director string   `yaml:"director"`
	Genre    []string `yaml:"genre"`
	// 游戏相关字段
	Developer string `yaml:"developer"`
}

// NoteIndexEntry 笔记索引条目
type NoteIndexEntry struct {
	Path    string `json:"path"`
	Type    string `json:"type"`
	SubType string `json:"subType"`
	Name    string `json:"name"`
}

// UpdateStatusRequest 更新状态请求
type UpdateStatusRequest struct {
	Status NoteStatus `json:"status"`
}

func (n *BookNote) ToSummary() NoteSummary {
	return NoteSummary{
		ID:      n.Title,
		Title:   n.Title,
		Type:    NoteTypeBook,
		SubType: string(n.Type),
		Status:  n.Status,
		Created: n.Created,
		Author:  n.Author,
		Tags:    n.Tags,
	}
}

func (n *VideoNote) ToSummary() NoteSummary {
	return NoteSummary{
		ID:      n.VideoID,
		Title:   n.Title,
		Type:    NoteTypeVideo,
		Status:  n.Status,
		Created: n.Created,
		Tags:    n.Tags,
	}
}

// GetCreatedTime 获取创建时间
func (n *BookNote) GetCreatedTime() time.Time {
	t, _ := time.Parse("2006-01-02", n.Created)
	return t
}

func (n *VideoNote) GetCreatedTime() time.Time {
	t, _ := time.Parse("2006-01-02", n.Created)
	return t
}

func (n *KnowledgeNote) ToSummary() NoteSummary {
	return NoteSummary{
		ID:      n.Path,
		Title:   n.Title,
		Type:    NoteTypeKnowledge,
		SubType: n.Topic,
	}
}

func (n *SkillNote) ToSummary() NoteSummary {
	return NoteSummary{
		ID:      n.Path,
		Title:   n.Name,
		Type:    NoteTypeSkill,
		SubType: n.Category,
	}
}

func (n *ProblemNote) ToSummary() NoteSummary {
	return NoteSummary{
		ID:      n.Path,
		Title:   n.Title,
		Type:    NoteTypeProblem,
		SubType: n.Type,
	}
}

func (n *IndexNote) ToSummary() NoteSummary {
	return NoteSummary{
		ID:      n.Path,
		Title:   n.Title,
		Type:    NoteTypeIndex,
		SubType: n.Type,
	}
}

func (n *QuoteNote) ToSummary() NoteSummary {
	return NoteSummary{
		ID:      n.ID,
		Title:   n.Quote,
		Type:    NoteTypeQuote,
		Author:  n.Author,
		Tags:    n.Tags,
	}
}

func (n *GitHubRepoNote) ToSummary() NoteSummary {
	return NoteSummary{
		ID:      n.ID,
		Title:   n.Title,
		Type:    NoteTypeGitHub,
		SubType: n.Language,
		Tags:    n.Tags,
	}
}

func (n *AnimeNote) ToSummary() NoteSummary {
	return NoteSummary{
		ID:      n.ID,
		Title:   n.Title,
		Type:    NoteTypeAnime,
		SubType: n.Status,
		Tags:    n.Tags,
	}
}

func (n *MovieNote) ToSummary() NoteSummary {
	return NoteSummary{
		ID:      n.ID,
		Title:   n.Title,
		Type:    NoteTypeMovie,
		SubType: n.Director,
		Tags:    n.Tags,
	}
}

func (n *GameNote) ToSummary() NoteSummary {
	return NoteSummary{
		ID:      n.ID,
		Title:   n.Title,
		Type:    NoteTypeGame,
		SubType: n.Platform,
		Tags:    n.Tags,
	}
}
