# Notes Web - AI 开发指南

## 项目概述

Notes Web 是一个基于 Vue 3 的个人笔记管理系统，支持 8 种笔记类型：书籍学习、视频学习、知识库、技能库、问题与解决、索引收藏、金句收藏、GitHub 项目收藏。每种笔记类型都支持 6 种视图模式。

## 技术栈

- **框架**: Vue 3.4 + TypeScript 5.3
- **构建工具**: Vite 5.0
- **状态管理**: Pinia 2.1
- **路由**: Vue Router 4.3
- **UI 组件库**: Element Plus 2.5
- **图表**: ECharts 6.0
- **Markdown 渲染**: markdown-it + highlight.js + KaTeX
- **自动导入**: unplugin-auto-import + unplugin-vue-components

## 项目结构

```
src/
├── api/                    # API 请求层
│   ├── index.ts           # 基础请求封装 (fetchAPI, putAPI, postAPI)
│   ├── books.ts           # 书籍 API
│   ├── videos.ts          # 视频 API
│   ├── knowledge.ts       # 知识库 API
│   ├── skills.ts          # 技能 API
│   ├── problems.ts        # 问题 API
│   ├── indexNotes.ts      # 索引 API
│   ├── quotes.ts          # 金句 API
│   ├── github.ts          # GitHub 项目 API
│   └── search.ts          # 搜索 API
│
├── components/            # 公共组件
│   ├── MdRenderer.vue     # Markdown 渲染器 (支持代码高亮、KaTeX 公式)
│   ├── NoteCard.vue       # 笔记卡片
│   ├── NoteDetailDisplay.vue  # 统一详情显示 (抽屉/对话框)
│   ├── OpenModeSelector.vue   # 打开方式选择器
│   ├── ProgressBar.vue    # 阅读进度条
│   ├── Sidebar.vue        # 左侧导航栏
│   ├── StatusTag.vue      # 状态标签
│   └── dashboard/         # Dashboard 相关组件
│       ├── Dashboard.vue  # Dashboard 布局
│       └── widgets/       # 可复用 Widget 组件
│           ├── NoteViewContainer.vue  # ⭐ 核心容器组件
│           ├── GridView.vue           # 网格视图
│           ├── ListView.vue           # 列表视图
│           ├── TableView.vue          # 表格视图
│           ├── TreeView.vue           # 树状视图
│           ├── KanbanView.vue         # 看板视图
│           ├── ChartView.vue          # 图表视图
│           ├── TimelineView.vue       # 时间线视图
│           ├── ViewSwitcher.vue       # 视图切换器
│           ├── StatsCard.vue          # 统计卡片
│           ├── ProgressCard.vue       # 进度卡片
│           ├── QuickEntryCard.vue     # 快捷入口卡片
│           └── QuickLinkCard.vue      # 快捷链接卡片
│
├── composables/           # 组合式函数
│   ├── useAsync.ts        # 异步状态管理
│   ├── useDashboard.ts    # Dashboard 配置
│   ├── useNoteViewSettings.ts  # 笔记视图设置 (打开模式)
│   └── useTree.ts         # 树形数据处理
│
├── stores/                # Pinia 状态管理
│   └── notes.ts           # ⭐ 笔记数据 Store
│
├── types/                 # TypeScript 类型定义
│   └── note.ts            # 笔记相关类型
│
├── views/                 # 页面视图
│   ├── Home.vue           # 首页 Dashboard
│   ├── Books.vue          # 书籍学习
│   ├── Videos.vue         # 视频学习
│   ├── Knowledge.vue      # 知识库
│   ├── Skills.vue         # 技能库
│   ├── Problems.vue       # 问题与解决
│   ├── Index.vue          # 索引收藏
│   ├── Quotes.vue         # 金句收藏
│   ├── BookDetail.vue     # 书籍详情页
│   └── VideoDetail.vue    # 视频详情页
│
├── router/index.ts        # 路由配置
├── App.vue                # 根组件
└── main.ts                # 入口文件
```

## 核心数据类型

```typescript
// 笔记类型
type NoteType = 'book' | 'video' | 'knowledge' | 'skill' | 'problem' | 'index' | 'quote' | 'github'

// 基础笔记
interface BaseNote {
  title: string
  tags?: string[]
  status?: NoteStatus
  created?: string
  content?: string
  path: string
}

// 各类型笔记继承 BaseNote
interface BookNote extends BaseNote {
  author: string
  category: string
  type: 'novel' | 'non-fiction'
  quotes?: string
  reviews?: string
}

interface VideoNote extends BaseNote {
  platform: string
  videoId: string
  url: string
}

interface KnowledgeNote extends BaseNote {
  category: string  // Math, Physics, etc.
  topic: string     // Algebra, Calculus, etc.
}

interface SkillNote extends BaseNote {
  name: string
  description: string
  category: string
  type: string  // SKILL, USER-SKILL
}

interface ProblemNote extends BaseNote {
  type: string  // "问题与解决" or "需求与解决"
}

interface IndexNote extends BaseNote {
  url?: string
  type?: string
}

// 金句笔记 (支持一个文件多条记录)
interface QuoteNote {
  id: string           // 唯一标识 (path-index 格式)
  path: string         // 文件路径
  quote: string        // 金句内容
  author?: string      // 作者
  source?: string      // 出处
  tags?: string[]      // 标签
  comment?: string     // 个人感悟/评论
}
```

## 核心组件使用

### NoteViewContainer (核心容器)

所有笔记列表页面的核心容器组件，统一处理 6 种视图模式。

```vue
<NoteViewContainer
  title="书籍学习"
  :items="notesStore.books"
  :loading="notesStore.loading"
  empty-text="暂无书籍"
  :available-views="['grid', 'list', 'table', 'tree', 'chart', 'kanban']"
  default-view="grid"
  :view-config="viewConfig"
  :load-detail="loadDetailFn"
  :detail-route="'/books/:type/:title'"
  @item-click="handleItemClick"
>
  <!-- 自定义 Grid 卡片 -->
  <template #gridItem="{ item }">
    <NoteCard :note="item" type="book" />
  </template>

  <!-- 自定义 Tree 右侧内容 -->
  <template #treeContent="{ item }">
    <MdRenderer :content="item.content" />
  </template>

  <!-- 自定义详情内容 (抽屉/对话框) -->
  <template #detail="{ item, loading }">
    <MdRenderer :content="item?.content" />
  </template>
</NoteViewContainer>
```

### viewConfig 配置

```typescript
const viewConfig = computed(() => ({
  titleField: 'title',
  descField: 'description',
  tagField: 'category',
  statusField: 'status',
  timeField: 'updatedAt',
  cardType: 'book' as const,
  gridColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
  
  // Table 列配置
  columns: [
    { prop: 'title', label: '标题', minWidth: 200, sortable: true },
    { prop: 'status', label: '状态', width: 120 },
    { prop: 'updatedAt', label: '更新时间', width: 160, type: 'date' },
  ],
  
  // Tree 数据
  treeData: treeData.value,
  
  // Chart 配置
  chartType: 'bar' as const,
  chartGroupBy: 'category',
  
  // Kanban 列配置
  kanbanColumns: [
    { id: 'notStarted', title: '未开始', statusValue: '未开始' },
    { id: 'reading', title: '阅读中', statusValue: '阅读中' },
    { id: 'finished', title: '已完成', statusValue: '已完成' },
  ],
}))
```

## 视图模式说明

| 视图 | 说明 | 点击行为 |
|------|------|----------|
| Grid | 网格卡片布局 | 使用打开模式 |
| List | 列表布局 | 使用打开模式 |
| Table | 表格布局 | 使用打开模式 |
| Tree | 树状+右侧内容 | 仅右侧显示，不走打开模式 |
| Kanban | 看板拖拽 | 使用打开模式 |
| Chart | 图表统计 | 使用打开模式 |

## 打开模式

用户可在界面选择打开方式：

| 模式 | 说明 |
|------|------|
| `drawer` | 右侧抽屉 |
| `dialog` | 居中对话框 |
| `page` | 跳转详情页 |

设置保存在 localStorage，key: `note-view-settings`

## API 约定

所有 API 请求基于 `/api` 前缀，由 Vite proxy 代理到后端。

```typescript
// 基础请求
fetchAPI<T>(endpoint: string): Promise<T>
putAPI<T>(endpoint: string, body: unknown): Promise<T>
postAPI<T>(endpoint: string, body: unknown): Promise<T>

// 示例
GET  /api/books           -> BookNote[]
GET  /api/books/:type/:name -> BookDetail
GET  /api/videos          -> VideoNote[]
GET  /api/knowledge/tree  -> KnowledgeTree[]
GET  /api/skills          -> SkillNote[]
GET  /api/problems        -> ProblemNote[]
GET  /api/index           -> IndexNote[]
GET  /api/quotes          -> QuoteNote[]
GET  /api/quotes/detail   -> QuoteNote (单条详情)
GET  /api/quotes/authors  -> string[] (作者列表)
GET  /api/quotes/tags     -> string[] (标签列表)
GET  /api/github          -> GitHubRepoNote[]
GET  /api/github/detail   -> GitHubRepoNote (单条详情)
GET  /api/github/languages -> string[] (语言列表)
GET  /api/github/tags     -> string[] (标签列表)
GET  /api/search?q=xxx    -> SearchResult[]
```

## 开发规范

### 1. 组件命名

- 页面组件：PascalCase (如 `Books.vue`)
- 公共组件：PascalCase (如 `NoteCard.vue`)
- Widget 组件：PascalCase + View/Card 后缀

### 2. 状态管理

- 使用 Pinia Store 集中管理数据
- 数据在 `App.vue` 的 `onMounted` 中预加载
- 各视图组件不再重复调用 `loadNotes()`

### 3. 类型安全

- 所有 API 返回数据需有类型定义
- 组件 props 使用 TypeScript 接口
- 避免 `any`，使用具体类型或泛型

### 4. 样式规范

- 使用 scoped 样式
- 颜色使用 Element Plus 变量
- 响应式布局使用 CSS Grid/Flex

### 5. 新增笔记类型步骤

1. 在 `types/note.ts` 添加类型定义
2. 在 `api/` 添加 API 模块
3. 在 `stores/notes.ts` 添加状态和方法
4. 在 `views/` 添加页面组件
5. 在 `router/index.ts` 添加路由
6. 在 `Sidebar.vue` 添加菜单项

## 常见问题

### Q: 页面切换后数据不显示？

确保 `App.vue` 中有预加载数据：
```typescript
onMounted(async () => {
  await notesStore.loadNotes()
})
```

### Q: GridView 点击不响应？

检查是否正确传递了 slot，并确保 `@item-click` 事件已绑定。

### Q: Markdown 渲染问题？

`MdRenderer` 组件需要 `content: string`，如果内容可能为空，使用：
```vue
<MdRenderer :content="item.content || ''" />
```

## 构建命令

```bash
npm run dev      # 开发服务器
npm run build    # 生产构建 (包含类型检查)
npm run preview  # 预览生产构建
```

## 金句收藏 (Quotes)

金句收藏是一种特殊的笔记类型，支持**一个文件包含多条记录**，适合批量管理大量金句。

### 文件格式

每个金句文件由多个 YAML frontmatter 块组成，每个块代表一条金句：

```markdown
---
quote: 人生·工作的结果 = 思维方式 × 热情 × 能力
author: 稻盛和夫
source: 活法
tags:
  - 人生哲学
  - 工作态度
comment: 这是我非常喜欢的一句话，它告诉我们要有正确的思维方式
---

---
quote: 习惯若不是最好的仆人，便就是最差的主人
author: 爱默生
source:
tags:
  - 习惯
  - 自律
comment:
---
```

### 字段说明

| 字段 | 必填 | 说明 |
|------|------|------|
| `quote` | 是 | 金句内容 |
| `author` | 否 | 作者 |
| `source` | 否 | 出处（书名、演讲等） |
| `tags` | 否 | 标签数组 |
| `comment` | 否 | 个人感悟/评论 |

### 存储位置

金句文件存放在 `changex-notes/Index/Quotes/` 目录下，按主题分文件：

```
changex-notes/Index/Quotes/
├── 人生哲学.md
├── 工作态度.md
├── 学习成长.md
└── ...
```

### 前端特性

- 支持按作者筛选
- 支持按标签筛选
- 点击金句展开显示完整内容和评论
- 统计显示金句总数

## GitHub 项目收藏 (GitHub)

GitHub 项目收藏是一种特殊的笔记类型，支持**一个文件包含多条记录**，适合批量管理 GitHub 项目。

### 文件格式

每个 GitHub 项目文件由多个 YAML frontmatter 块组成，每个块代表一个项目：

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
url: https://github.com/raoooool/flomo-reminder
name: raoooool/flomo-reminder
title: flomo-reminder
stars: 56
language: Go
topics: [flomo, reminder]
tags: [工具]
comment:
---
```

### 字段说明

| 字段 | 必填 | 说明 |
|------|------|------|
| `url` | 是 | GitHub 项目链接 |
| `name` | 是 | owner/repo 格式 |
| `title` | 否 | 项目标题 |
| `description` | 否 | 项目描述 |
| `stars` | 否 | Star 数 |
| `forks` | 否 | Fork 数 |
| `language` | 否 | 主要语言 |
| `topics` | 否 | GitHub Topics |
| `tags` | 否 | 个人标签 |
| `comment` | 否 | 个人备注 |

### 存储位置

GitHub 项目文件存放在 `changex-notes/Index/GitHub/` 目录下：

```
changex-notes/Index/GitHub/
├── AI工具.md
├── 前端框架.md
├── Go项目.md
└── ...
```

### 前端特性

- 支持按语言筛选
- 支持按标签筛选
- 点击卡片可直接跳转到 GitHub
- 展示 Stars、Forks 等统计信息
- 支持多种视图模式（网格/列表/表格/图表/看板）
