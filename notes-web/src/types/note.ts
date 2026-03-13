// 笔记状态枚举
export enum NoteStatus {
  NotStarted = '未开始',
  Reading = '阅读中',
  Completed = '已完成',
  Interpreted = '已解读',
}

// 笔记类型
export type NoteType = 'book' | 'video' | 'knowledge' | 'skill' | 'problem' | 'index' | 'quote' | 'github' | 'anime' | 'movie' | 'game'

// 基础笔记接口
export interface BaseNote {
  title: string
  tags?: string[]
  status?: NoteStatus
  created?: string
  updatedAt?: string
  content?: string
  path: string
}

// 书籍笔记接口
export interface BookNote extends BaseNote {
  author: string
  category: string
  quotes?: string
  reviews?: string
  type: 'novel' | 'non-fiction'
}

// 视频笔记接口
export interface VideoNote extends BaseNote {
  platform: string
  videoId: string
  url: string
  subtitle?: string
}

// 知识库笔记接口
export interface KnowledgeNote extends BaseNote {
  category: string  // Math, Physics, etc.
  topic: string     // Algebra, Calculus, etc.
}

// 技能笔记接口
export interface SkillNote extends BaseNote {
  name: string
  description: string
  category: string  // Mine/Others, skills/code-skills etc.
  type: string      // SKILL, USER-SKILL
}

// 问题解决笔记接口
export interface ProblemNote extends BaseNote {
  type: string      // "问题与解决" or "需求与解决"
  problem?: string
  solution?: string
}

// 索引笔记接口
export interface IndexNote extends BaseNote {
  url?: string
  type?: string     // 小说, 动漫, etc.
}

// 金句笔记接口
export interface QuoteNote {
  quote: string           // 金句内容
  author?: string         // 出处人
  source?: string         // 出处处
  tags?: string[]
  comment?: string        // 解读/感想
  path: string            // 文件路径
  id: string              // 唯一标识
}

// GitHub 项目收藏接口
export interface GitHubRepoNote {
  url: string             // GitHub 项目链接
  name: string            // owner/repo 格式
  title?: string          // 项目标题
  desc?: string           // 项目描述
  stars?: number          // Star 数
  forks?: number          // Fork 数
  lang?: string           // 主要语言
  topics?: string[]       // GitHub Topics
  tags?: string[]         // 个人标签
  comment?: string        // 个人备注
  path: string            // 文件路径
  id: string              // 唯一标识
}

// 动漫收藏接口
export interface AnimeNote {
  url?: string            // 链接
  title: string           // 标题
  year?: number           // 年份
  studio?: string         // 制作公司
  episodes?: number       // 集数
  status?: string         // 状态：连载/完结
  rating?: number         // 评分
  tags?: string[]         // 标签
  comment?: string        // 评论
  path: string            // 文件路径
  id: string              // 唯一标识
}

// 电影收藏接口
export interface MovieNote {
  url?: string            // 链接
  title: string           // 标题
  year?: number           // 年份
  director?: string       // 导演
  rating?: number         // 评分
  genre?: string[]        // 类型
  tags?: string[]         // 标签
  comment?: string        // 观后感
  path: string            // 文件路径
  id: string              // 唯一标识
}

// 游戏收藏接口
export interface GameNote {
  url?: string            // 链接
  title: string           // 标题
  platform?: string       // 平台
  developer?: string      // 开发商
  year?: number           // 年份
  status?: string         // 状态：想玩/在玩/通关
  rating?: number         // 评分
  tags?: string[]         // 标签
  comment?: string        // 评论
  path: string            // 文件路径
  id: string              // 唯一标识
}

// 笔记摘要
export interface NoteSummary {
  id: string
  title: string
  type: NoteType
  subType?: string
  status?: NoteStatus
  created?: string
  author?: string
  tags?: string[]
}

// 笔记统计
export interface NoteStats {
  totalBooks: number
  totalVideos: number
  readingBooks: number
  completedBooks: number
  totalKnowledge: number
  totalSkills: number
  totalProblems: number
  totalIndex: number
  totalQuotes: number
  totalGitHub: number
  totalAnime: number
  totalMovies: number
  totalGames: number
  recentNotes: NoteSummary[]
}

// 搜索结果
export interface SearchResult {
  id: string
  title: string
  type: NoteType
  snippet: string
}
