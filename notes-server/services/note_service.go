package services

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"notes-server/config"
	"notes-server/models"
)

type NoteService struct {
	cfg         *config.Config
	books       []models.BookNote
	videos      []models.VideoNote
	knowledge   []models.KnowledgeNote
	skills      []models.SkillNote
	problems    []models.ProblemNote
	indexNotes  []models.IndexNote
	quotes      []models.QuoteNote
	gitHubRepos []models.GitHubRepoNote
	anime       []models.AnimeNote
	movies      []models.MovieNote
	games       []models.GameNote
	indexPath   string
}

func NewNoteService(cfg *config.Config) *NoteService {
	return &NoteService{
		cfg:       cfg,
		indexPath: filepath.Join(cfg.NotesPath, "notes-web", "public", "notes-index.json"),
	}
}

// LoadNotes 加载所有笔记
func (s *NoteService) LoadNotes() error {
	// NotesPath 指向 changex-notes 根目录
	booksPath := filepath.Join(s.cfg.NotesPath, "Resource", "书籍学习")
	videosPath := filepath.Join(s.cfg.NotesPath, "Resource", "视频学习")
	knowledgePath := filepath.Join(s.cfg.NotesPath, "Knowledge")
	codePath := filepath.Join(s.cfg.NotesPath, "Code")
	projectPath := filepath.Join(s.cfg.NotesPath, "Project")
	indexPath := filepath.Join(s.cfg.NotesPath, "Index")

	// 加载书籍
	if err := s.loadBooks(booksPath); err != nil {
		return fmt.Errorf("failed to load books: %w", err)
	}

	// 加载视频
	if err := s.loadVideos(videosPath); err != nil {
		return fmt.Errorf("failed to load videos: %w", err)
	}

	// 加载知识库
	if err := s.loadKnowledge(knowledgePath); err != nil {
		fmt.Printf("Warning: failed to load knowledge: %v\n", err)
	}

	// 加载技能
	if err := s.loadSkills(codePath); err != nil {
		fmt.Printf("Warning: failed to load skills: %v\n", err)
	}

	// 加载问题解决
	if err := s.loadProblems(projectPath); err != nil {
		fmt.Printf("Warning: failed to load problems: %v\n", err)
	}

	// 加载索引
	if err := s.loadIndex(indexPath); err != nil {
		fmt.Printf("Warning: failed to load index: %v\n", err)
	}

	// 加载金句
	quotesPath := filepath.Join(s.cfg.NotesPath, "Index", "Quotes")
	if err := s.loadQuotes(quotesPath); err != nil {
		fmt.Printf("Warning: failed to load quotes: %v\n", err)
	}

	// 加载 GitHub 项目收藏
	gitHubPath := filepath.Join(s.cfg.NotesPath, "Index", "GitHub")
	if err := s.loadGitHubRepos(gitHubPath); err != nil {
		fmt.Printf("Warning: failed to load GitHub repos: %v\n", err)
	}

	// 加载动漫收藏
	animePath := filepath.Join(s.cfg.NotesPath, "Index", "Anime")
	if err := s.loadAnime(animePath); err != nil {
		fmt.Printf("Warning: failed to load anime: %v\n", err)
	}

	// 加载电影收藏
	moviesPath := filepath.Join(s.cfg.NotesPath, "Index", "Movies")
	if err := s.loadMovies(moviesPath); err != nil {
		fmt.Printf("Warning: failed to load movies: %v\n", err)
	}

	// 加载游戏收藏
	gamesPath := filepath.Join(s.cfg.NotesPath, "Index", "Games")
	if err := s.loadGames(gamesPath); err != nil {
		fmt.Printf("Warning: failed to load games: %v\n", err)
	}

	return nil
}

// loadBooks 加载书籍笔记
func (s *NoteService) loadBooks(rootPath string) error {
	subDirs, err := os.ReadDir(rootPath)
	if err != nil {
		return err
	}

	for _, subDir := range subDirs {
		if !subDir.IsDir() {
			continue
		}

		subPath := filepath.Join(rootPath, subDir.Name())
		bookDirs, err := os.ReadDir(subPath)
		if err != nil {
			continue
		}

		for _, bookDir := range bookDirs {
			if !bookDir.IsDir() {
				continue
			}

			bookPath := filepath.Join(subPath, bookDir.Name(), "index.md")
			content, err := os.ReadFile(bookPath)
			if err != nil {
				continue
			}

			fm, body := parseFrontmatter(string(content))
			bookType := models.BookSubTypeNonFiction
			if subDir.Name() == "小说" {
				bookType = models.BookSubTypeNovel
			}

			book := models.BookNote{
				Title:    fm.Title,
				Author:   fm.Author,
				Category: fm.Category,
				Tags:     fm.Tags,
				Status:   models.NoteStatus(fm.Status),
				Created:  fm.Created,
				Content:  body,
				Type:     bookType,
				Path:     bookPath,
			}

			// 加载 quotes 和 reviews
			basePath := filepath.Join(subPath, bookDir.Name())
			if quotes, err := os.ReadFile(filepath.Join(basePath, "quotes.md")); err == nil {
				_, body := parseFrontmatter(string(quotes))
				book.Quotes = body
			}
			if reviews, err := os.ReadFile(filepath.Join(basePath, "reviews.md")); err == nil {
				_, body := parseFrontmatter(string(reviews))
				book.Reviews = body
			}

			s.books = append(s.books, book)
		}
	}

	return nil
}

// loadVideos 加载视频笔记
func (s *NoteService) loadVideos(rootPath string) error {
	videoDirs, err := os.ReadDir(rootPath)
	if err != nil {
		return err
	}

	for _, videoDir := range videoDirs {
		if !videoDir.IsDir() {
			continue
		}

		videoPath := filepath.Join(rootPath, videoDir.Name(), "index.md")
		content, err := os.ReadFile(videoPath)
		if err != nil {
			continue
		}

		fm, body := parseFrontmatter(string(content))

		video := models.VideoNote{
			Title:    fm.Title,
			Platform: fm.Platform,
			VideoID:  fm.VideoID,
			URL:      fm.URL,
			Tags:     fm.Tags,
			Status:   models.NoteStatus(fm.Status),
			Created:  fm.Created,
			Content:  body,
			Path:     videoPath,
		}

		s.videos = append(s.videos, video)
	}

	return nil
}

// loadKnowledge 加载知识库笔记
func (s *NoteService) loadKnowledge(rootPath string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// 从路径提取分类和主题
		relPath, _ := filepath.Rel(rootPath, path)
		parts := strings.Split(relPath, string(filepath.Separator))
		category := ""
		topic := ""
		if len(parts) > 0 {
			category = parts[0]
		}
		if len(parts) > 1 {
			topic = parts[1]
		}

		// 从文件名提取标题
		title := strings.TrimSuffix(filepath.Base(path), ".md")

		s.knowledge = append(s.knowledge, models.KnowledgeNote{
			Title:    title,
			Category: category,
			Topic:    topic,
			Content:  string(content),
			Path:     path,
		})
		return nil
	})
}

// loadSkills 加载技能笔记
func (s *NoteService) loadSkills(rootPath string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		filename := filepath.Base(path)
		if filename != "SKILL.md" && filename != "USER-SKILL.md" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// 解析 frontmatter
		fm, body := parseFrontmatter(string(content))

		// 从路径提取分类
		relPath, _ := filepath.Rel(rootPath, path)
		parts := strings.Split(relPath, string(filepath.Separator))
		category := ""
		if len(parts) > 1 {
			category = parts[0] + "/" + parts[1]
		} else if len(parts) > 0 {
			category = parts[0]
		}

		skillType := "SKILL"
		if filename == "USER-SKILL.md" {
			skillType = "USER-SKILL"
		}

		name := fm.Title
		if name == "" {
			// 从路径提取名称
			if len(parts) > 2 {
				name = parts[2]
			}
		}

		s.skills = append(s.skills, models.SkillNote{
			Name:        name,
			Description: fm.Category, // 复用 category 字段存储 description
			Category:    category,
			Type:        skillType,
			Content:     body,
			Path:        path,
		})
		return nil
	})
}

// loadProblems 加载问题解决笔记
func (s *NoteService) loadProblems(rootPath string) error {
	subDirs, err := os.ReadDir(rootPath)
	if err != nil {
		return err
	}

	for _, subDir := range subDirs {
		if !subDir.IsDir() {
			continue
		}

		subPath := filepath.Join(rootPath, subDir.Name())
		files, err := os.ReadDir(subPath)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
				continue
			}

			filePath := filepath.Join(subPath, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				continue
			}

			title := strings.TrimSuffix(file.Name(), ".md")

			s.problems = append(s.problems, models.ProblemNote{
				Title:   title,
				Type:    subDir.Name(),
				Content: string(content),
				Path:    filePath,
			})
		}
	}

	return nil
}

// loadIndex 加载索引笔记
func (s *NoteService) loadIndex(rootPath string) error {
	files, err := os.ReadDir(rootPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		filePath := filepath.Join(rootPath, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		title := strings.TrimSuffix(file.Name(), ".md")

		s.indexNotes = append(s.indexNotes, models.IndexNote{
			Title:   title,
			Content: string(content),
			Path:    filePath,
		})
	}

	return nil
}

// loadQuotes 加载金句笔记(支持单文件多条记录)
func (s *NoteService) loadQuotes(rootPath string) error {
	// 检查目录是否存在
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// 解析多条 frontmatter 记录
		records := parseMultipleFrontmatter(string(content))
		for i, fm := range records {
			if fm.Quote == "" {
				continue
			}

			// 生成唯一 ID: 路径 + 索引
			id := fmt.Sprintf("%s#%d", path, i)

			s.quotes = append(s.quotes, models.QuoteNote{
				Quote:   fm.Quote,
				Author:  fm.Author,
				Source:  fm.Source,
				Tags:    fm.Tags,
				Comment: fm.Comment,
				Path:    path,
				ID:      id,
			})
		}
		return nil
	})
}

// loadGitHubRepos 加载 GitHub 项目收藏(支持单文件多条记录)
func (s *NoteService) loadGitHubRepos(rootPath string) error {
	// 检查目录是否存在
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// 解析多条 frontmatter 记录
		records := parseMultipleFrontmatter(string(content))
		for i, fm := range records {
			if fm.URL == "" {
				continue
			}

			// 生成唯一 ID: 路径 + 索引
			id := fmt.Sprintf("%s#%d", path, i)

			s.gitHubRepos = append(s.gitHubRepos, models.GitHubRepoNote{
				URL:         fm.URL,
				Name:        fm.Name,
				Title:       fm.Title,
				Description: fm.Description,
				Stars:       fm.Stars,
				Forks:       fm.Forks,
				Language:    fm.Language,
				Topics:      fm.Topics,
				Tags:        fm.Tags,
				Comment:     fm.Comment,
				Path:        path,
				ID:          id,
			})
		}
		return nil
	})
}

// loadAnime 加载动漫收藏(支持单文件多条记录)
func (s *NoteService) loadAnime(rootPath string) error {
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		records := parseMultipleFrontmatter(string(content))
		for i, fm := range records {
			if fm.Title == "" {
				continue
			}

			id := fmt.Sprintf("%s#%d", path, i)

			s.anime = append(s.anime, models.AnimeNote{
				URL:      fm.URL,
				Title:    fm.Title,
				Year:     fm.Year,
				Studio:   fm.Studio,
				Episodes: fm.Episodes,
				Status:   fm.Status,
				Rating:   fm.Rating,
				Tags:     fm.Tags,
				Comment:  fm.Comment,
				Path:     path,
				ID:       id,
			})
		}
		return nil
	})
}

// loadMovies 加载电影收藏(支持单文件多条记录)
func (s *NoteService) loadMovies(rootPath string) error {
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		records := parseMultipleFrontmatter(string(content))
		for i, fm := range records {
			if fm.Title == "" {
				continue
			}

			id := fmt.Sprintf("%s#%d", path, i)

			s.movies = append(s.movies, models.MovieNote{
				URL:      fm.URL,
				Title:    fm.Title,
				Year:     fm.Year,
				Director: fm.Director,
				Rating:   fm.Rating,
				Genre:    fm.Genre,
				Tags:     fm.Tags,
				Comment:  fm.Comment,
				Path:     path,
				ID:       id,
			})
		}
		return nil
	})
}

// loadGames 加载游戏收藏(支持单文件多条记录)
func (s *NoteService) loadGames(rootPath string) error {
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		records := parseMultipleFrontmatter(string(content))
		for i, fm := range records {
			if fm.Title == "" {
				continue
			}

			id := fmt.Sprintf("%s#%d", path, i)

			s.games = append(s.games, models.GameNote{
				URL:       fm.URL,
				Title:     fm.Title,
				Platform:  fm.Platform,
				Developer: fm.Developer,
				Year:      fm.Year,
				Status:    fm.Status,
				Rating:    fm.Rating,
				Tags:      fm.Tags,
				Comment:   fm.Comment,
				Path:      path,
				ID:        id,
			})
		}
		return nil
	})
}

// parseMultipleFrontmatter 解析单文件多条 frontmatter 记录
func parseMultipleFrontmatter(content string) []models.Frontmatter {
	var records []models.Frontmatter

	// 按 --- 分割
	parts := strings.Split(content, "---")

	var currentFm models.Frontmatter

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// 检查是否是 YAML 块
		if isYamlBlock(part) {
			// 如果之前有记录，保存它 (支持 Quote 和 GitHub URL)
			if currentFm.Quote != "" || currentFm.Title != "" || currentFm.URL != "" {
				records = append(records, currentFm)
			}
			currentFm = models.Frontmatter{}

			// 解析 YAML
			scanner := bufio.NewScanner(strings.NewReader(part))
			for scanner.Scan() {
				line := scanner.Text()
				colonIdx := strings.Index(line, ":")
				if colonIdx == -1 {
					continue
				}

				key := strings.TrimSpace(line[:colonIdx])
				value := strings.TrimSpace(line[colonIdx+1:])

				switch key {
				case "title":
					currentFm.Title = strings.Trim(value, `"`)
				case "author":
					currentFm.Author = strings.Trim(value, `"`)
				case "category":
					currentFm.Category = strings.Trim(value, `"`)
				case "status":
					currentFm.Status = strings.Trim(value, `"`)
				case "created":
					currentFm.Created = strings.Trim(value, `"`)
				case "platform":
					currentFm.Platform = strings.Trim(value, `"`)
				case "video_id":
					currentFm.VideoID = strings.Trim(value, `"`)
				case "url":
					currentFm.URL = strings.Trim(value, `"`)
				case "quote":
					currentFm.Quote = strings.Trim(value, `"`)
				case "source":
					currentFm.Source = strings.Trim(value, `"`)
				case "comment":
					currentFm.Comment = strings.Trim(value, `"`)
				case "tags":
					if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
						tagsStr := value[1 : len(value)-1]
						tags := strings.Split(tagsStr, ",")
						for i, t := range tags {
							tags[i] = strings.TrimSpace(strings.Trim(t, `"`))
						}
						currentFm.Tags = tags
					}
				// GitHub 项目相关字段
				case "name":
					currentFm.Name = strings.Trim(value, `"`)
				case "description":
					currentFm.Description = strings.Trim(value, `"`)
				case "stars":
					if v, err := strconv.Atoi(value); err == nil {
						currentFm.Stars = v
					}
				case "forks":
					if v, err := strconv.Atoi(value); err == nil {
						currentFm.Forks = v
					}
				case "language":
					currentFm.Language = strings.Trim(value, `"`)
				case "topics":
					if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
						topicsStr := value[1 : len(value)-1]
						topics := strings.Split(topicsStr, ",")
						for i, t := range topics {
							topics[i] = strings.TrimSpace(strings.Trim(t, `"`))
						}
						currentFm.Topics = topics
					}
				// 动漫、电影、游戏相关字段
				case "year":
					if v, err := strconv.Atoi(value); err == nil {
						currentFm.Year = v
					}
				case "studio":
					currentFm.Studio = strings.Trim(value, `"`)
				case "episodes":
					if v, err := strconv.Atoi(value); err == nil {
						currentFm.Episodes = v
					}
				case "rating":
					if v, err := strconv.ParseFloat(value, 64); err == nil {
						currentFm.Rating = v
					}
				case "director":
					currentFm.Director = strings.Trim(value, `"`)
				case "genre":
					if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
						genreStr := value[1 : len(value)-1]
						genre := strings.Split(genreStr, ",")
						for i, g := range genre {
							genre[i] = strings.TrimSpace(strings.Trim(g, `"`))
						}
						currentFm.Genre = genre
					}
				case "developer":
					currentFm.Developer = strings.Trim(value, `"`)
				}
			}
		}
	}

	// 添加最后一条记录 (支持 Quote 和 GitHub URL)
	if currentFm.Quote != "" || currentFm.Title != "" || currentFm.URL != "" {
		records = append(records, currentFm)
	}

	return records
}

// isYamlBlock 检查是否是 YAML 块(包含 key: value 格式)
func isYamlBlock(content string) bool {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// 检查是否有 key: value 格式
		if strings.Contains(line, ":") {
			return true
		}
	}
	return false
}

// parseFrontmatter 解析 YAML frontmatter
func parseFrontmatter(content string) (models.Frontmatter, string) {
	var fm models.Frontmatter

	// 检查是否有 frontmatter
	if !strings.HasPrefix(content, "---\n") {
		return fm, content
	}

	// 找到第二个 ---
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return fm, content
	}

	fmStr := parts[1]
	body := strings.TrimSpace(parts[2])

	// 简单解析 YAML
	scanner := bufio.NewScanner(strings.NewReader(fmStr))
	for scanner.Scan() {
		line := scanner.Text()
		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			continue
		}

		key := strings.TrimSpace(line[:colonIdx])
		value := strings.TrimSpace(line[colonIdx+1:])

		switch key {
		case "title":
			fm.Title = strings.Trim(value, `"`)
		case "author":
			fm.Author = strings.Trim(value, `"`)
		case "category":
			fm.Category = strings.Trim(value, `"`)
		case "status":
			fm.Status = strings.Trim(value, `"`)
		case "created":
			fm.Created = strings.Trim(value, `"`)
		case "platform":
			fm.Platform = strings.Trim(value, `"`)
		case "video_id":
			fm.VideoID = strings.Trim(value, `"`)
		case "url":
			fm.URL = strings.Trim(value, `"`)
		case "quote":
			fm.Quote = strings.Trim(value, `"`)
		case "source":
			fm.Source = strings.Trim(value, `"`)
		case "comment":
			fm.Comment = strings.Trim(value, `"`)
		case "tags":
			// 解析数组 [tag1, tag2]
			if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
				tagsStr := value[1 : len(value)-1]
				tags := strings.Split(tagsStr, ",")
				for i, t := range tags {
					tags[i] = strings.TrimSpace(strings.Trim(t, `"`))
				}
				fm.Tags = tags
			}
		}
	}

	return fm, body
}

// GetAllNotes 获取所有笔记摘要
func (s *NoteService) GetAllNotes(noteType string, status string) []models.NoteSummary {
	var summaries []models.NoteSummary

	// 过滤书籍
	if noteType == "" || noteType == "book" {
		for _, book := range s.books {
			if status == "" || string(book.Status) == status {
				summaries = append(summaries, book.ToSummary())
			}
		}
	}

	// 过滤视频
	if noteType == "" || noteType == "video" {
		for _, video := range s.videos {
			if status == "" || string(video.Status) == status {
				summaries = append(summaries, video.ToSummary())
			}
		}
	}

	// 按创建时间排序
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Created > summaries[j].Created
	})

	return summaries
}

// GetBooks 获取书籍列表
func (s *NoteService) GetBooks(subType string, status string) []models.NoteSummary {
	var summaries []models.NoteSummary

	for _, book := range s.books {
		if subType != "" && string(book.Type) != subType {
			continue
		}
		if status != "" && string(book.Status) != status {
			continue
		}
		summaries = append(summaries, book.ToSummary())
	}

	return summaries
}

// GetVideos 获取视频列表
func (s *NoteService) GetVideos(status string) []models.NoteSummary {
	var summaries []models.NoteSummary

	for _, video := range s.videos {
		if status != "" && string(video.Status) != status {
			continue
		}
		summaries = append(summaries, video.ToSummary())
	}

	return summaries
}

// GetBookDetail 获取书籍详情
func (s *NoteService) GetBookDetail(bookType string, title string) (*models.BookNote, error) {
	for _, book := range s.books {
		if string(book.Type) == bookType && book.Title == title {
			return &book, nil
		}
	}
	return nil, fmt.Errorf("book not found: %s/%s", bookType, title)
}

// GetVideoDetail 获取视频详情
func (s *NoteService) GetVideoDetail(videoID string) (*models.VideoNote, error) {
	for _, video := range s.videos {
		if video.VideoID == videoID {
			return &video, nil
		}
	}
	return nil, fmt.Errorf("video not found: %s", videoID)
}

// GetKnowledge 获取知识库列表
func (s *NoteService) GetKnowledge(category string, topic string) []models.NoteSummary {
	var summaries []models.NoteSummary

	for _, kn := range s.knowledge {
		if category != "" && kn.Category != category {
			continue
		}
		if topic != "" && kn.Topic != topic {
			continue
		}
		summaries = append(summaries, kn.ToSummary())
	}

	return summaries
}

// GetKnowledgeDetail 获取知识库详情
func (s *NoteService) GetKnowledgeDetail(path string) (*models.KnowledgeNote, error) {
	for _, kn := range s.knowledge {
		if kn.Path == path {
			return &kn, nil
		}
	}
	return nil, fmt.Errorf("knowledge not found: %s", path)
}

// GetKnowledgeDetailByQuery 按分类/主题/标题查询知识库详情
func (s *NoteService) GetKnowledgeDetailByQuery(category, topic, title string) (*models.KnowledgeNote, error) {
	for _, kn := range s.knowledge {
		if kn.Category == category && kn.Topic == topic && kn.Title == title {
			return &kn, nil
		}
	}
	return nil, fmt.Errorf("knowledge not found: %s/%s/%s", category, topic, title)
}

// GetKnowledgeTree 获取知识库树形结构
func (s *NoteService) GetKnowledgeTree() map[string]map[string][]string {
	tree := make(map[string]map[string][]string)
	for _, kn := range s.knowledge {
		if tree[kn.Category] == nil {
			tree[kn.Category] = make(map[string][]string)
		}
		tree[kn.Category][kn.Topic] = append(tree[kn.Category][kn.Topic], kn.Title)
	}
	return tree
}

// GetSkills 获取技能列表
func (s *NoteService) GetSkills(category string, skillType string) []models.NoteSummary {
	var summaries []models.NoteSummary

	for _, skill := range s.skills {
		if category != "" && skill.Category != category {
			continue
		}
		if skillType != "" && skill.Type != skillType {
			continue
		}
		summaries = append(summaries, skill.ToSummary())
	}

	return summaries
}

// GetSkillDetail 获取技能详情
func (s *NoteService) GetSkillDetail(path string) (*models.SkillNote, error) {
	for _, skill := range s.skills {
		if skill.Path == path {
			return &skill, nil
		}
	}
	return nil, fmt.Errorf("skill not found: %s", path)
}

// GetProblems 获取问题解决列表
func (s *NoteService) GetProblems(problemType string) []models.NoteSummary {
	var summaries []models.NoteSummary

	for _, prob := range s.problems {
		if problemType != "" && prob.Type != problemType {
			continue
		}
		summaries = append(summaries, prob.ToSummary())
	}

	return summaries
}

// GetProblemDetail 获取问题解决详情
func (s *NoteService) GetProblemDetail(path string) (*models.ProblemNote, error) {
	for _, prob := range s.problems {
		if prob.Path == path {
			return &prob, nil
		}
	}
	return nil, fmt.Errorf("problem not found: %s", path)
}

// GetIndexNotes 获取索引列表
func (s *NoteService) GetIndexNotes() []models.NoteSummary {
	var summaries []models.NoteSummary
	for _, idx := range s.indexNotes {
		summaries = append(summaries, idx.ToSummary())
	}
	return summaries
}

// GetIndexDetail 获取索引详情
func (s *NoteService) GetIndexDetail(path string) (*models.IndexNote, error) {
	for _, idx := range s.indexNotes {
		if idx.Path == path {
			return &idx, nil
		}
	}
	return nil, fmt.Errorf("index not found: %s", path)
}

// GetQuotes 获取金句列表
func (s *NoteService) GetQuotes(author string, tag string) []models.NoteSummary {
	var summaries []models.NoteSummary

	for _, quote := range s.quotes {
		if author != "" && quote.Author != author {
			continue
		}
		if tag != "" {
			hasTag := false
			for _, t := range quote.Tags {
				if t == tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		summaries = append(summaries, quote.ToSummary())
	}

	return summaries
}

// GetQuoteDetail 获取金句详情
func (s *NoteService) GetQuoteDetail(id string) (*models.QuoteNote, error) {
	for _, quote := range s.quotes {
		if quote.ID == id {
			return &quote, nil
		}
	}
	return nil, fmt.Errorf("quote not found: %s", id)
}

// GetQuoteAuthors 获取所有作者列表
func (s *NoteService) GetQuoteAuthors() []string {
	authorMap := make(map[string]bool)
	for _, quote := range s.quotes {
		if quote.Author != "" {
			authorMap[quote.Author] = true
		}
	}

	var authors []string
	for author := range authorMap {
		authors = append(authors, author)
	}
	sort.Strings(authors)
	return authors
}

// GetQuoteTags 获取所有标签列表
func (s *NoteService) GetQuoteTags() []string {
	tagMap := make(map[string]bool)
	for _, quote := range s.quotes {
		for _, tag := range quote.Tags {
			if tag != "" {
				tagMap[tag] = true
			}
		}
	}

	var tags []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// GetGitHubRepos 获取 GitHub 项目列表
func (s *NoteService) GetGitHubRepos(language string, tag string) []models.NoteSummary {
	var summaries []models.NoteSummary

	for _, repo := range s.gitHubRepos {
		if language != "" && repo.Language != language {
			continue
		}
		if tag != "" {
			hasTag := false
			for _, t := range repo.Tags {
				if t == tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		summaries = append(summaries, repo.ToSummary())
	}

	return summaries
}

// GetGitHubRepoDetail 获取 GitHub 项目详情
func (s *NoteService) GetGitHubRepoDetail(id string) (*models.GitHubRepoNote, error) {
	for _, repo := range s.gitHubRepos {
		if repo.ID == id {
			return &repo, nil
		}
	}
	return nil, fmt.Errorf("github repo not found: %s", id)
}

// GetGitHubLanguages 获取所有语言列表
func (s *NoteService) GetGitHubLanguages() []string {
	langMap := make(map[string]bool)
	for _, repo := range s.gitHubRepos {
		if repo.Language != "" {
			langMap[repo.Language] = true
		}
	}

	var languages []string
	for lang := range langMap {
		languages = append(languages, lang)
	}
	sort.Strings(languages)
	return languages
}

// GetGitHubTags 获取所有标签列表
func (s *NoteService) GetGitHubTags() []string {
	tagMap := make(map[string]bool)
	for _, repo := range s.gitHubRepos {
		for _, tag := range repo.Tags {
			if tag != "" {
				tagMap[tag] = true
			}
		}
	}

	var tags []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// ==================== 动漫相关方法 ====================

// GetAnime 获取动漫列表
func (s *NoteService) GetAnime(status string, tag string) []models.NoteSummary {
	var summaries []models.NoteSummary

	for _, a := range s.anime {
		if status != "" && a.Status != status {
			continue
		}
		if tag != "" {
			hasTag := false
			for _, t := range a.Tags {
				if t == tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		summaries = append(summaries, a.ToSummary())
	}

	return summaries
}

// GetAnimeDetail 获取动漫详情
func (s *NoteService) GetAnimeDetail(id string) (*models.AnimeNote, error) {
	for _, a := range s.anime {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, fmt.Errorf("anime not found: %s", id)
}

// GetAnimeStatuses 获取所有状态列表
func (s *NoteService) GetAnimeStatuses() []string {
	statusMap := make(map[string]bool)
	for _, a := range s.anime {
		if a.Status != "" {
			statusMap[a.Status] = true
		}
	}

	var statuses []string
	for status := range statusMap {
		statuses = append(statuses, status)
	}
	sort.Strings(statuses)
	return statuses
}

// GetAnimeTags 获取所有标签列表
func (s *NoteService) GetAnimeTags() []string {
	tagMap := make(map[string]bool)
	for _, a := range s.anime {
		for _, tag := range a.Tags {
			if tag != "" {
				tagMap[tag] = true
			}
		}
	}

	var tags []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// ==================== 电影相关方法 ====================

// GetMovies 获取电影列表
func (s *NoteService) GetMovies(genre string, tag string) []models.NoteSummary {
	var summaries []models.NoteSummary

	for _, m := range s.movies {
		if genre != "" {
			hasGenre := false
			for _, g := range m.Genre {
				if g == genre {
					hasGenre = true
					break
				}
			}
			if !hasGenre {
				continue
			}
		}
		if tag != "" {
			hasTag := false
			for _, t := range m.Tags {
				if t == tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		summaries = append(summaries, m.ToSummary())
	}

	return summaries
}

// GetMovieDetail 获取电影详情
func (s *NoteService) GetMovieDetail(id string) (*models.MovieNote, error) {
	for _, m := range s.movies {
		if m.ID == id {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("movie not found: %s", id)
}

// GetMovieGenres 获取所有类型列表
func (s *NoteService) GetMovieGenres() []string {
	genreMap := make(map[string]bool)
	for _, m := range s.movies {
		for _, g := range m.Genre {
			if g != "" {
				genreMap[g] = true
			}
		}
	}

	var genres []string
	for genre := range genreMap {
		genres = append(genres, genre)
	}
	sort.Strings(genres)
	return genres
}

// GetMovieTags 获取所有标签列表
func (s *NoteService) GetMovieTags() []string {
	tagMap := make(map[string]bool)
	for _, m := range s.movies {
		for _, tag := range m.Tags {
			if tag != "" {
				tagMap[tag] = true
			}
		}
	}

	var tags []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// ==================== 游戏相关方法 ====================

// GetGames 获取游戏列表
func (s *NoteService) GetGames(platform string, status string, tag string) []models.NoteSummary {
	var summaries []models.NoteSummary

	for _, g := range s.games {
		if platform != "" && g.Platform != platform {
			continue
		}
		if status != "" && g.Status != status {
			continue
		}
		if tag != "" {
			hasTag := false
			for _, t := range g.Tags {
				if t == tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		summaries = append(summaries, g.ToSummary())
	}

	return summaries
}

// GetGameDetail 获取游戏详情
func (s *NoteService) GetGameDetail(id string) (*models.GameNote, error) {
	for _, g := range s.games {
		if g.ID == id {
			return &g, nil
		}
	}
	return nil, fmt.Errorf("game not found: %s", id)
}

// GetGamePlatforms 获取所有平台列表
func (s *NoteService) GetGamePlatforms() []string {
	platformMap := make(map[string]bool)
	for _, g := range s.games {
		if g.Platform != "" {
			platformMap[g.Platform] = true
		}
	}

	var platforms []string
	for platform := range platformMap {
		platforms = append(platforms, platform)
	}
	sort.Strings(platforms)
	return platforms
}

// GetGameStatuses 获取所有状态列表
func (s *NoteService) GetGameStatuses() []string {
	statusMap := make(map[string]bool)
	for _, g := range s.games {
		if g.Status != "" {
			statusMap[g.Status] = true
		}
	}

	var statuses []string
	for status := range statusMap {
		statuses = append(statuses, status)
	}
	sort.Strings(statuses)
	return statuses
}

// GetGameTags 获取所有标签列表
func (s *NoteService) GetGameTags() []string {
	tagMap := make(map[string]bool)
	for _, g := range s.games {
		for _, tag := range g.Tags {
			if tag != "" {
				tagMap[tag] = true
			}
		}
	}

	var tags []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// GetStats 获取统计信息
func (s *NoteService) GetStats() models.NoteStats {
	stats := models.NoteStats{
		TotalBooks:     len(s.books),
		TotalVideos:    len(s.videos),
		TotalKnowledge: len(s.knowledge),
		TotalSkills:    len(s.skills),
		TotalProblems:  len(s.problems),
		TotalIndex:     len(s.indexNotes),
		TotalQuotes:    len(s.quotes),
		TotalGitHub:    len(s.gitHubRepos),
		TotalAnime:     len(s.anime),
		TotalMovies:    len(s.movies),
		TotalGames:     len(s.games),
	}

	for _, book := range s.books {
		if book.Status == models.StatusReading {
			stats.ReadingBooks++
		}
		if book.Status == models.StatusCompleted || book.Status == models.StatusInterpreted {
			stats.CompletedBooks++
		}
	}

	return stats
}

// Search 搜索笔记
func (s *NoteService) Search(query string) []models.SearchResult {
	var results []models.SearchResult
	query = strings.ToLower(query)

	// 搜索书籍
	for _, book := range s.books {
		if s.matchNote(book.Title, book.Content, query) {
			snippet := getSnippet(book.Content, query, 100)
			results = append(results, models.SearchResult{
				ID:      book.Title,
				Title:   book.Title,
				Type:    models.NoteTypeBook,
				Snippet: snippet,
			})
		}
	}

	// 搜索视频
	for _, video := range s.videos {
		if s.matchNote(video.Title, video.Content, query) {
			snippet := getSnippet(video.Content, query, 100)
			results = append(results, models.SearchResult{
				ID:      video.VideoID,
				Title:   video.Title,
				Type:    models.NoteTypeVideo,
				Snippet: snippet,
			})
		}
	}

	// 搜索知识库
	for _, kn := range s.knowledge {
		if s.matchNote(kn.Title, kn.Content, query) {
			snippet := getSnippet(kn.Content, query, 100)
			results = append(results, models.SearchResult{
				ID:      kn.Path,
				Title:   kn.Title,
				Type:    models.NoteTypeKnowledge,
				Snippet: snippet,
			})
		}
	}

	// 搜索技能
	for _, skill := range s.skills {
		searchText := skill.Name + " " + skill.Description + " " + skill.Content
		if s.matchNote(skill.Name, searchText, query) {
			snippet := getSnippet(skill.Content, query, 100)
			results = append(results, models.SearchResult{
				ID:      skill.Path,
				Title:   skill.Name,
				Type:    models.NoteTypeSkill,
				Snippet: snippet,
			})
		}
	}

	// 搜索问题解决
	for _, prob := range s.problems {
		if s.matchNote(prob.Title, prob.Content, query) {
			snippet := getSnippet(prob.Content, query, 100)
			results = append(results, models.SearchResult{
				ID:      prob.Path,
				Title:   prob.Title,
				Type:    models.NoteTypeProblem,
				Snippet: snippet,
			})
		}
	}

	// 搜索索引
	for _, idx := range s.indexNotes {
		if s.matchNote(idx.Title, idx.Content, query) {
			snippet := getSnippet(idx.Content, query, 100)
			results = append(results, models.SearchResult{
				ID:      idx.Path,
				Title:   idx.Title,
				Type:    models.NoteTypeIndex,
				Snippet: snippet,
			})
		}
	}

	// 搜索金句
	for _, quote := range s.quotes {
		searchText := quote.Quote + " " + quote.Author + " " + quote.Source + " " + quote.Comment
		if s.matchNote(quote.Quote, searchText, query) {
			snippet := getSnippet(quote.Quote+" "+quote.Comment, query, 100)
			results = append(results, models.SearchResult{
				ID:      quote.ID,
				Title:   quote.Quote,
				Type:    models.NoteTypeQuote,
				Snippet: snippet,
			})
		}
	}

	return results
}

// matchNote 检查是否匹配搜索词
func (s *NoteService) matchNote(title, content, query string) bool {
	title = strings.ToLower(title)
	content = strings.ToLower(content)

	return strings.Contains(title, query) || strings.Contains(content, query)
}

// getSnippet 获取搜索片段
func getSnippet(content, query string, maxLen int) string {
	content = strings.ToLower(content)
	query = strings.ToLower(query)

	idx := strings.Index(content, query)
	if idx == -1 {
		if len(content) > maxLen {
			return content[:maxLen] + "..."
		}
		return content
	}

	start := idx - 30
	if start < 0 {
		start = 0
	}
	end := idx + len(query) + 30
	if end > len(content) {
		end = len(content)
	}

	snippet := content[start:end]
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(content) {
		snippet = snippet + "..."
	}

	return snippet
}

// UpdateNoteStatus 更新笔记状态
func (s *NoteService) UpdateNoteStatus(noteType string, id string, status models.NoteStatus) error {
	if noteType == "book" {
		for i, book := range s.books {
			if book.Title == id {
				s.books[i].Status = status
				return s.updateFileStatus(book.Path, status)
			}
		}
	} else if noteType == "video" {
		for i, video := range s.videos {
			if video.VideoID == id {
				s.videos[i].Status = status
				return s.updateFileStatus(video.Path, status)
			}
		}
	}
	return fmt.Errorf("note not found")
}

// updateFileStatus 更新文件中的状态
func (s *NoteService) updateFileStatus(filePath string, status models.NoteStatus) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 替换 status 行
	lines := bytes.Split(content, []byte("\n"))
	for i, line := range lines {
		if bytes.HasPrefix(line, []byte("status:")) {
			lines[i] = []byte(fmt.Sprintf("status: %s", status))
			break
		}
	}

	newContent := bytes.Join(lines, []byte("\n"))
	return os.WriteFile(filePath, newContent, 0644)
}

// GetRecentNotes 获取最近更新的笔记
func (s *NoteService) GetRecentNotes(limit int) []models.NoteSummary {
	var summaries []models.NoteSummary

	for _, book := range s.books {
		summaries = append(summaries, book.ToSummary())
	}
	for _, video := range s.videos {
		summaries = append(summaries, video.ToSummary())
	}
	for _, kn := range s.knowledge {
		summaries = append(summaries, kn.ToSummary())
	}
	for _, skill := range s.skills {
		summaries = append(summaries, skill.ToSummary())
	}
	for _, prob := range s.problems {
		summaries = append(summaries, prob.ToSummary())
	}

	// 按创建时间排序
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Created > summaries[j].Created
	})

	if len(summaries) > limit {
		summaries = summaries[:limit]
	}

	return summaries
}
