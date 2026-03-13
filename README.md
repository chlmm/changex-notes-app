# Notes App

个人笔记管理系统，用于管理和展示多类型笔记内容。

## 技术栈

| 模块 | 技术 |
|------|------|
| 前端 | Vue 3 + Vite + TypeScript + Element Plus |
| 后端 | Go + Gin |
| 数据 | Markdown 文件（独立仓库） |

## 目录结构

```
changex-notes-app/
├── notes-server/          # Go 后端服务
│   ├── config/           # 配置
│   ├── handlers/         # 请求处理
│   ├── models/           # 数据模型
│   ├── router/           # 路由
│   └── services/         # 业务逻辑
├── notes-web/            # Vue 前端
│   └── src/
│       ├── api/          # API 请求
│       ├── components/   # 组件
│       ├── views/        # 页面
│       └── stores/       # 状态管理
├── .env                  # 配置文件
└── start.sh              # 启动脚本
```

## 快速开始

### 1. 克隆仓库

```bash
git clone https://github.com/chlmm/changex-notes-app.git
cd changex-notes-app
```

### 2. 配置数据路径

编辑 `.env` 文件，设置笔记数据路径：

```env
NOTES_PATH=/path/to/changex-notes
STATIC_PATH=./notes-web/dist
```

> 笔记数据仓库：[changex-notes](https://github.com/chlmm/changex-notes)

### 3. 安装依赖

```bash
# 前端依赖
cd notes-web && npm install

# 后端依赖（Go 会自动下载）
cd ../notes-server
```

### 4. 启动服务

```bash
./start.sh
```

- 前端：http://localhost:3000
- 后端：http://localhost:8080

### 其他命令

```bash
./start.sh stop     # 停止服务
./start.sh restart  # 重启服务
./start.sh logs     # 查看日志
```

## 功能特性

- 📚 **书籍管理** - 书籍笔记、读书进度、金句摘录
- 🎬 **视频笔记** - 视频学习记录
- 📖 **知识库** - 分类知识整理
- 💻 **技能树** - 技能学习追踪
- 🔧 **问题解决** - 问题记录与解决方案
- 💬 **金句收藏** - 名言警句收集
- ⭐ **GitHub 收藏** - 开源项目收藏
- 🎌 **动漫收藏** - 动漫观看记录
- 🎥 **电影收藏** - 电影观看记录
- 🎮 **游戏收藏** - 游戏进度追踪

## 配置说明

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `NOTES_PATH` | 笔记数据目录 | `../changex-notes` |
| `STATIC_PATH` | 前端静态文件目录 | `./notes-web/dist` |

## 开发

### 前端开发

```bash
cd notes-web
npm run dev
```

### 后端开发

```bash
cd notes-server
go run main.go
```

### 构建生产版本

```bash
cd notes-web
npm run build
```

## License

MIT
