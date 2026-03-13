# 新增笔记类型开发指南

本文档描述如何为 Notes Web 项目添加新的笔记类型。以 GitHub 项目收藏为例，完整记录从后端到前端的实现步骤。

## 概述

添加新笔记类型需要修改以下文件：

| 层级 | 文件 | 说明 |
|------|------|------|
| 后端数据 | `models/note.go` | 数据结构定义 |
| 后端服务 | `services/note_service.go` | 加载和查询逻辑 |
| 后端处理 | `handlers/note_handler.go` | HTTP Handler |
| 后端路由 | `router/router.go` | API 路由 |
| 前端类型 | `types/note.ts` | TypeScript 类型 |
| 前端API | `api/<type>.ts` | API 调用层 |
| 前端状态 | `stores/notes.ts` | Pinia Store |
| 前端页面 | `views/<Type>.vue` | 页面组件 |
| 前端路由 | `router/index.ts` | 路由配置 |
| 前端菜单 | `components/Sidebar.vue` | 侧边栏菜单 |
| 文档 | `SKILL.md` | 功能文档 |

## 第一步：后端数据结构

### 1.1 添加 NoteType 常量

文件：`notes-server/models/note.go`

```go
const (
    // ... 现有类型
    NoteTypeGitHub NoteType = "github"  // 新增
)
```

### 1.2 添加数据结构

```go
// GitHubRepoNote GitHub 项目收藏
type GitHubRepoNote struct {
    URL         string   `json:"url"`              // 项目链接
    Name        string   `json:"name"`             // 名称
    Title       string   `json:"title,omitempty"`  // 标题
    Description string   `json:"desc,omitempty"`   // 描述
    Stars       int      `json:"stars,omitempty"`  // Star 数
    Forks       int      `json:"forks,omitempty"`  // Fork 数
    Language    string   `json:"lang,omitempty"`   // 语言
    Topics      []string `json:"topics,omitempty"` // Topics
    Tags        []string `json:"tags,omitempty"`   // 标签
    Comment     string   `json:"comment,omitempty"` // 备注
    Path        string   `json:"path"`             // 文件路径
    ID          string   `json:"id"`               // 唯一标识
}
```

### 1.3 添加 ToSummary 方法

```go
func (n *GitHubRepoNote) ToSummary() NoteSummary {
    return NoteSummary{
        ID:      n.ID,
        Title:   n.Title,
        Type:    NoteTypeGitHub,
        SubType: n.Language,
        Tags:    n.Tags,
    }
}
```

### 1.4 更新 Frontmatter 结构

```go
type Frontmatter struct {
    // ... 现有字段
    // 新增字段
    Name        string   `yaml:"name"`
    Description string   `yaml:"description"`
    Stars       int      `yaml:"stars"`
    Forks       int      `yaml:"forks"`
    Language    string   `yaml:"language"`
    Topics      []string `yaml:"topics"`
}
```

### 1.5 更新 NoteStats

```go
type NoteStats struct {
    // ... 现有字段
    TotalGitHub int `json:"totalGitHub"` // 新增
}
```

## 第二步：后端服务层

### 2.1 添加数据存储

文件：`notes-server/services/note_service.go`

```go
type NoteService struct {
    // ... 现有字段
    gitHubRepos []models.GitHubRepoNote  // 新增
}
```

### 2.2 在 LoadNotes 中调用加载

```go
func (s *NoteService) LoadNotes() error {
    // ... 现有加载逻辑
    
    // 加载 GitHub 项目
    gitHubPath := filepath.Join(s.cfg.NotesPath, "..", "Index", "GitHub")
    if err := s.loadGitHubRepos(gitHubPath); err != nil {
        fmt.Printf("Warning: failed to load GitHub repos: %v\n", err)
    }
    
    return nil
}
```

### 2.3 添加加载函数

```go
// loadGitHubRepos 加载 GitHub 项目收藏(支持单文件多条记录)
func (s *NoteService) loadGitHubRepos(rootPath string) error {
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
            if fm.URL == "" {  // 使用 URL 作为判断条件
                continue
            }

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
```

### 2.4 更新 parseMultipleFrontmatter

在 switch 语句中添加新字段解析：

```go
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
```

更新判断条件：

```go
// 支持 Quote 和 GitHub URL
if currentFm.Quote != "" || currentFm.Title != "" || currentFm.URL != "" {
    records = append(records, currentFm)
}
```

### 2.5 添加查询方法

```go
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
```

### 2.6 更新 GetStats

```go
func (s *NoteService) GetStats() models.NoteStats {
    stats := models.NoteStats{
        // ... 现有字段
        TotalGitHub: len(s.gitHubRepos),
    }
    // ...
}
```

## 第三步：后端 Handler

文件：`notes-server/handlers/note_handler.go`

```go
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
    decodedID, _ := url.PathUnescape(id)
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
```

## 第四步：后端路由

文件：`notes-server/router/router.go`

```go
// GitHub 项目收藏
api.GET("/github", noteHandler.GetGitHubRepos)
api.GET("/github/detail", noteHandler.GetGitHubRepoDetail)
api.GET("/github/languages", noteHandler.GetGitHubLanguages)
api.GET("/github/tags", noteHandler.GetGitHubTags)
```

## 第五步：前端类型定义

文件：`notes-web/src/types/note.ts`

```typescript
// 更新 NoteType
export type NoteType = 'book' | 'video' | 'knowledge' | 'skill' | 'problem' | 'index' | 'quote' | 'github'

// 添加接口
export interface GitHubRepoNote {
  url: string
  name: string
  title?: string
  desc?: string
  stars?: number
  forks?: number
  lang?: string
  topics?: string[]
  tags?: string[]
  comment?: string
  path: string
  id: string
}

// 更新 NoteStats
export interface NoteStats {
  // ... 现有字段
  totalGitHub: number
}
```

## 第六步：前端 API 层

文件：`notes-web/src/api/github.ts`

```typescript
import { fetchAPI } from './index'
import type { GitHubRepoNote } from '@/types/note'

export interface GitHubRepoListResponse {
  id: string
  title: string
  subType?: string
  tags?: string[]
}

export interface GitHubRepoDetailResponse {
  url: string
  name: string
  title?: string
  desc?: string
  stars?: number
  forks?: number
  lang?: string
  topics?: string[]
  tags?: string[]
  comment?: string
  path: string
  id: string
}

export async function getGitHubRepos(language?: string, tag?: string): Promise<GitHubRepoListResponse[]> {
  const params = new URLSearchParams()
  if (language) params.append('language', language)
  if (tag) params.append('tag', tag)
  const query = params.toString() ? `?${params.toString()}` : ''
  return fetchAPI<GitHubRepoListResponse[]>(`/github${query}`)
}

export async function getGitHubRepoDetail(id: string): Promise<GitHubRepoDetailResponse> {
  return fetchAPI<GitHubRepoDetailResponse>(`/github/detail?id=${encodeURIComponent(id)}`)
}

export async function getGitHubLanguages(): Promise<string[]> {
  return fetchAPI<string[]>('/github/languages')
}

export async function getGitHubTags(): Promise<string[]> {
  return fetchAPI<string[]>('/github/tags')
}

export function transformGitHubRepoItem(item: GitHubRepoListResponse): GitHubRepoNote {
  return {
    url: '',
    name: item.title,
    title: item.title,
    lang: item.subType,
    tags: item.tags || [],
    path: item.id,
    id: item.id,
  }
}

export function transformGitHubRepoDetail(data: GitHubRepoDetailResponse): GitHubRepoNote {
  return {
    url: data.url,
    name: data.name,
    title: data.title || '',
    desc: data.desc || '',
    stars: data.stars,
    forks: data.forks,
    lang: data.lang || '',
    topics: data.topics || [],
    tags: data.tags || [],
    comment: data.comment || '',
    path: data.path,
    id: data.id,
  }
}
```

## 第七步：前端 Store

文件：`notes-web/src/stores/notes.ts`

```typescript
// 导入
import type { GitHubRepoNote } from '@/types/note'
import * as githubApi from '@/api/github'

// 状态
const githubRepos = ref<GitHubRepoNote[]>([])

// stats 计算属性
const stats = computed<NoteStats>(() => ({
  // ... 现有字段
  totalGitHub: githubRepos.value.length,
}))

// loadNotes 中加载
const [/* ... */, githubData] = await Promise.all([
  // ... 现有调用
  githubApi.getGitHubRepos(),
])
githubRepos.value = githubData.map(githubApi.transformGitHubRepoItem)

// 添加方法
async function getGitHubRepoDetail(id: string): Promise<GitHubRepoNote | null> {
  try {
    const data = await githubApi.getGitHubRepoDetail(id)
    return githubApi.transformGitHubRepoDetail(data)
  } catch (e) {
    console.error('Failed to load github repo detail:', e)
    return null
  }
}

async function getGitHubLanguages() {
  return githubApi.getGitHubLanguages()
}

async function getGitHubTags() {
  return githubApi.getGitHubTags()
}

// 导出
return {
  // ... 现有导出
  githubRepos,
  getGitHubRepoDetail,
  getGitHubLanguages,
  getGitHubTags,
}
```

## 第八步：前端页面组件

文件：`notes-web/src/views/GitHub.vue`

参考 `views/Quotes.vue` 创建页面，关键点：

1. 使用 `NoteViewContainer` 组件
2. 实现筛选功能（语言/标签）
3. 自定义 Grid/List 卡片样式
4. 实现 `loadDetail` 函数

## 第九步：前端路由

文件：`notes-web/src/router/index.ts`

```typescript
{
  path: '/github',
  name: 'GitHub',
  component: () => import('@/views/GitHub.vue'),
  meta: { title: 'GitHub 项目收藏' },
},
```

## 第十步：前端菜单

文件：`notes-web/src/components/Sidebar.vue`

```vue
<el-menu-item index="/github">
  <el-icon><Link /></el-icon>
  <span>GitHub 收藏</span>
</el-menu-item>
```

更新统计：

```typescript
const totalNotes = computed(() => 
  // ... 现有字段
  notesStore.stats.totalGitHub
)
```

## 第十一步：创建数据目录和示例文件

```bash
mkdir -p changex-notes/Index/GitHub
```

创建示例文件 `changex-notes/Index/GitHub/AI工具.md`：

```markdown
---
url: https://github.com/jina-ai/MCP
name: jina-ai/MCP
title: Jina AI Remote MCP Server
description: Model Context Protocol server implementation
stars: 1234
language: Python
topics: [mcp, ai]
tags: [AI工具, 值得关注]
comment: 可以用于扩展 AI 能力的 MCP 服务器
---

---
url: https://github.com/example/project
name: example/project
title: Example Project
stars: 56
language: Go
topics: [tool]
tags: [工具]
comment:
---
```

## 第十二步：更新文档

文件：`notes-web/SKILL.md`

1. 更新项目概述中的笔记类型数量
2. 更新项目结构中的 API 文件列表
3. 更新核心数据类型
4. 更新 API 约定
5. 添加新类型的完整文档

## 扩展：添加电影/动漫等资源

按照上述流程，只需替换以下内容：

### 电影收藏示例

**数据目录**: `changex-notes/Index/Movies/`

**数据结构**:
```go
type MovieNote struct {
    URL       string   `json:"url"`        // 链接
    Title     string   `json:"title"`      // 标题
    Year      int      `json:"year"`       // 年份
    Director  string   `json:"director"`   // 导演
    Rating    float64  `json:"rating"`     // 评分
    Genre     []string `json:"genre"`      // 类型
    Tags      []string `json:"tags"`       // 个人标签
    Comment   string   `json:"comment"`    // 观后感
    Path      string   `json:"path"`
    ID        string   `json:"id"`
}
```

**文件格式**:
```markdown
---
url: https://movie.douban.com/subject/123456
title: 肖申克的救赎
year: 1994
director: 弗兰克·德拉邦特
rating: 9.7
genre: [剧情, 犯罪]
tags: [经典, 治愈]
comment: 希望让人自由
---
```

### 动漫收藏示例

**数据目录**: `changex-notes/Index/Anime/`

**数据结构**:
```go
type AnimeNote struct {
    URL       string   `json:"url"`
    Title     string   `json:"title"`
    Year      int      `json:"year"`
    Studio    string   `json:"studio"`     // 制作公司
    Episodes  int      `json:"episodes"`   // 集数
    Status    string   `json:"status"`     // 状态：连载/完结
    Rating    float64  `json:"rating"`
    Tags      []string `json:"tags"`
    Comment   string   `json:"comment"`
    Path      string   `json:"path"`
    ID        string   `json:"id"`
}
```

## 检查清单

- [ ] 后端: NoteType 常量
- [ ] 后端: 数据结构体
- [ ] 后端: ToSummary 方法
- [ ] 后端: Frontmatter 字段
- [ ] 后端: NoteStats 字段
- [ ] 后端: Service 数据存储
- [ ] 后端: Service 加载函数
- [ ] 后端: Service 查询方法
- [ ] 后端: Handler 函数
- [ ] 后端: 路由配置
- [ ] 后端: 编译通过
- [ ] 前端: TypeScript 类型
- [ ] 前端: API 层
- [ ] 前端: Store 状态和方法
- [ ] 前端: 页面组件
- [ ] 前端: 路由配置
- [ ] 前端: 菜单项
- [ ] 数据: 创建目录
- [ ] 数据: 创建示例文件
- [ ] 文档: 更新 SKILL.md
