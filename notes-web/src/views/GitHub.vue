<template>
  <div class="github-page">
    <!-- 筛选区域 -->
    <div class="filter-section">
      <el-select
        v-model="selectedLanguage"
        placeholder="按语言筛选"
        clearable
        @change="handleFilterChange"
      >
        <el-option
          v-for="lang in languages"
          :key="lang"
          :label="lang"
          :value="lang"
        />
      </el-select>
      <el-select
        v-model="selectedTag"
        placeholder="按标签筛选"
        clearable
        style="margin-left: 12px;"
        @change="handleFilterChange"
      >
        <el-option
          v-for="tag in tags"
          :key="tag"
          :label="tag"
          :value="tag"
        />
      </el-select>
    </div>

    <NoteViewContainer
      title="GitHub 项目收藏"
      :items="filteredRepos"
      :loading="notesStore.loading"
      empty-text="暂无收藏的项目"
      :available-views="['grid', 'list', 'table', 'chart', 'kanban']"
      default-view="grid"
      :view-config="viewConfig"
      :load-detail="loadRepoDetail"
      @item-click="handleItemClick"
    >
      <!-- GridView 自定义卡片 -->
      <template #gridItem="{ item }">
        <div class="repo-card" @click="openGitHub(item.url)">
          <div class="repo-header">
            <el-icon class="repo-icon"><Link /></el-icon>
            <span class="repo-name">{{ item.name }}</span>
            <el-tag v-if="item.lang" size="small" type="info" class="lang-tag">{{ item.lang }}</el-tag>
          </div>
          <p class="repo-title">{{ item.title || item.name }}</p>
          <div class="repo-stats">
            <span v-if="item.stars" class="stat">
              <el-icon><Star /></el-icon>
              {{ formatNumber(item.stars) }}
            </span>
            <span v-if="item.forks" class="stat">
              <el-icon><Fork /></el-icon>
              {{ formatNumber(item.forks) }}
            </span>
          </div>
          <div class="repo-tags">
            <el-tag
              v-for="tag in (item.tags || []).slice(0, 3)"
              :key="tag"
              size="small"
              type="warning"
              class="tag"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>
      </template>

      <!-- ListView 自定义项 -->
      <template #listItem="{ item }">
        <div class="repo-list-item" @click="openGitHub(item.url)">
          <div class="repo-content">
            <div class="repo-header">
              <span class="repo-name">{{ item.name }}</span>
              <el-tag v-if="item.lang" size="small" type="info">{{ item.lang }}</el-tag>
            </div>
            <p class="repo-desc">{{ item.title || item.desc || '暂无描述' }}</p>
            <div class="repo-stats">
              <span v-if="item.stars" class="stat">
                <el-icon><Star /></el-icon>
                {{ formatNumber(item.stars) }}
              </span>
              <span v-if="item.forks" class="stat">
                <el-icon><Fork /></el-icon>
                {{ formatNumber(item.forks) }}
              </span>
            </div>
          </div>
          <div class="repo-tags">
            <el-tag
              v-for="tag in (item.tags || []).slice(0, 2)"
              :key="tag"
              size="small"
              type="warning"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>
      </template>

      <!-- Table 操作列 -->
      <template #tableOperation="{ row }">
        <el-button size="small" type="primary" link @click="openGitHub(row.url)">
          打开
        </el-button>
      </template>

      <!-- 详情内容（抽屉/对话框） -->
      <template #detail="{ item, loading }">
        <div class="repo-detail-content">
          <template v-if="loading">
            <el-skeleton :rows="6" animated />
          </template>
          <template v-else-if="item">
            <div class="detail-header">
              <h3 class="detail-name">{{ item.name }}</h3>
              <p v-if="item.title" class="detail-title">{{ item.title }}</p>
              <el-button type="primary" @click="openGitHub(item.url)">
                <el-icon><Link /></el-icon>
                访问 GitHub
              </el-button>
            </div>
            
            <div class="detail-stats">
              <div v-if="item.stars" class="stat-item">
                <el-icon><Star /></el-icon>
                <span>{{ formatNumber(item.stars) }} Stars</span>
              </div>
              <div v-if="item.forks" class="stat-item">
                <el-icon><Fork /></el-icon>
                <span>{{ formatNumber(item.forks) }} Forks</span>
              </div>
              <div v-if="item.lang" class="stat-item">
                <el-icon><Document /></el-icon>
                <span>{{ item.lang }}</span>
              </div>
            </div>

            <el-divider v-if="item.desc" />
            <div v-if="item.desc" class="detail-desc">
              <h4>项目描述</h4>
              <p>{{ item.desc }}</p>
            </div>

            <div v-if="item.topics?.length" class="detail-topics">
              <h4>Topics</h4>
              <div class="topics-list">
                <el-tag
                  v-for="topic in item.topics"
                  :key="topic"
                  size="small"
                  type="info"
                  style="margin-right: 8px; margin-bottom: 8px;"
                >
                  {{ topic }}
                </el-tag>
              </div>
            </div>

            <el-divider v-if="item.comment" />
            <div v-if="item.comment" class="detail-comment">
              <h4>个人备注</h4>
              <p>{{ item.comment }}</p>
            </div>

            <div class="detail-tags">
              <h4>标签</h4>
              <el-tag
                v-for="tag in item.tags"
                :key="tag"
                size="small"
                type="warning"
                style="margin-right: 8px;"
              >
                {{ tag }}
              </el-tag>
            </div>
          </template>
        </div>
      </template>
    </NoteViewContainer>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useNotesStore } from '@/stores/notes'
import NoteViewContainer from '@/components/dashboard/widgets/NoteViewContainer.vue'
import { Link, Star, Document } from '@element-plus/icons-vue'

// Fork 图标 (Element Plus 没有，用 Star 旋转替代)
const Fork = Star

const notesStore = useNotesStore()

// 筛选状态
const selectedLanguage = ref('')
const selectedTag = ref('')
const languages = ref<string[]>([])
const tags = ref<string[]>([])

// 筛选后的项目列表
const filteredRepos = computed(() => {
  let result = notesStore.githubRepos
  
  if (selectedLanguage.value) {
    result = result.filter(r => r.lang === selectedLanguage.value)
  }
  
  if (selectedTag.value) {
    result = result.filter(r => r.tags?.includes(selectedTag.value))
  }
  
  return result
})

// 视图配置
const viewConfig = computed(() => ({
  titleField: 'name',
  descField: 'desc',
  tagField: 'tags',
  timeField: 'createdAt',
  cardType: 'github' as const,
  gridColumns: 'repeat(auto-fill, minmax(320px, 1fr))',
  columns: [
    { prop: 'name', label: '项目名称', minWidth: 200, sortable: true },
    { prop: 'title', label: '描述', minWidth: 250 },
    { prop: 'lang', label: '语言', width: 100 },
    { prop: 'stars', label: 'Stars', width: 100, sortable: true },
    { prop: 'tags', label: '标签', width: 200 },
  ],
  chartType: 'bar' as const,
  chartGroupBy: 'lang',
  kanbanColumns: getKanbanColumns(),
}))

// 动态生成看板列（按语言分组）
function getKanbanColumns() {
  const uniqueLanguages = [...new Set(notesStore.githubRepos.map(r => r.lang || '未知'))]
  
  return uniqueLanguages.slice(0, 10).map(lang => ({
    id: lang,
    title: lang,
    statusValue: lang,
  }))
}

// 格式化数字
function formatNumber(num: number): string {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return String(num)
}

// 打开 GitHub 链接
function openGitHub(url: string) {
  if (url) {
    window.open(url, '_blank')
  }
}

// 加载详情的函数
async function loadRepoDetail(item: any) {
  if (item.id) {
    return await notesStore.getGitHubRepoDetail(item.id)
  }
  return item
}

// 处理点击事件
async function handleItemClick(item: any) {
  console.log('GitHub repo clicked:', item)
}

// 筛选变化时重新加载
async function handleFilterChange() {
  // 筛选逻辑已通过 computed 实现
}

// 加载语言和标签列表
async function loadFilterOptions() {
  languages.value = await notesStore.getGitHubLanguages()
  tags.value = await notesStore.getGitHubTags()
}

onMounted(async () => {
  await notesStore.loadNotes()
  await loadFilterOptions()
})
</script>

<style scoped>
.github-page {
  padding: 0;
}

.filter-section {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
}

/* GridView 卡片 */
.repo-card {
  padding: 20px;
  border-radius: 12px;
  background: linear-gradient(135deg, #f0f9eb 0%, #fff 100%);
  border: 1px solid #e1f3d8;
  cursor: pointer;
  transition: all 0.3s;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.repo-card:hover {
  border-color: #67c23a;
  box-shadow: 0 4px 16px rgba(103, 194, 58, 0.15);
  transform: translateY(-2px);
}

.repo-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.repo-icon {
  color: #909399;
  font-size: 18px;
}

.repo-name {
  font-weight: 600;
  font-size: 16px;
  color: #303133;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.lang-tag {
  flex-shrink: 0;
}

.repo-title {
  margin: 0;
  font-size: 13px;
  line-height: 1.6;
  color: #606266;
  flex: 1;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.repo-stats {
  display: flex;
  gap: 16px;
  margin-top: 12px;
}

.stat {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #909399;
}

.repo-tags {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.repo-tags .tag {
  margin: 0;
}

/* ListView 项 */
.repo-list-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px;
  width: 100%;
  cursor: pointer;
}

.repo-list-item:hover {
  background: #f5f7fa;
}

.repo-list-item .repo-content {
  flex: 1;
  min-width: 0;
}

.repo-list-item .repo-header {
  margin-bottom: 8px;
}

.repo-list-item .repo-name {
  font-size: 15px;
  margin-right: 8px;
}

.repo-list-item .repo-desc {
  margin: 0;
  font-size: 13px;
  color: #909399;
  line-height: 1.5;
}

.repo-list-item .repo-stats {
  margin-top: 8px;
}

.repo-list-item .repo-tags {
  margin-top: 0;
  margin-left: 16px;
  flex-shrink: 0;
}

/* 详情内容 */
.repo-detail-content {
  max-height: 70vh;
  overflow: auto;
}

.detail-header {
  margin-bottom: 20px;
}

.detail-name {
  margin: 0 0 8px 0;
  font-size: 20px;
  color: #303133;
}

.detail-title {
  margin: 0 0 16px 0;
  font-size: 14px;
  color: #909399;
}

.detail-stats {
  display: flex;
  gap: 24px;
  margin-bottom: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
}

.detail-desc h4,
.detail-comment h4,
.detail-tags h4,
.detail-topics h4 {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #909399;
}

.detail-desc p,
.detail-comment p {
  margin: 0;
  font-size: 14px;
  line-height: 1.8;
  color: #606266;
}

.detail-topics {
  margin-top: 16px;
}

.topics-list {
  display: flex;
  flex-wrap: wrap;
}

.detail-comment {
  margin-top: 16px;
}

.detail-tags {
  margin-top: 16px;
}
</style>
