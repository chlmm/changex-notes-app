<template>
  <div class="movies-page">
    <!-- 筛选区域 -->
    <div class="filter-section">
      <el-select
        v-model="selectedGenre"
        placeholder="按类型筛选"
        clearable
        @change="handleFilterChange"
      >
        <el-option
          v-for="genre in genres"
          :key="genre"
          :label="genre"
          :value="genre"
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
      title="电影收藏"
      :items="filteredMovies"
      :loading="notesStore.loading"
      empty-text="暂无收藏的电影"
      :available-views="['grid', 'list', 'table', 'chart', 'kanban']"
      default-view="grid"
      :view-config="viewConfig"
      :load-detail="loadMovieDetail"
      @item-click="handleItemClick"
    >
      <!-- GridView 自定义卡片 -->
      <template #gridItem="{ item }">
        <div class="movie-card" @click="openUrl(item.url)">
          <div class="movie-header">
            <span class="movie-title">{{ item.title }}</span>
            <span v-if="item.year" class="movie-year">{{ item.year }}</span>
          </div>
          <div class="movie-info">
            <span v-if="item.director">导演: {{ item.director }}</span>
          </div>
          <div v-if="item.rating" class="movie-rating">
            <el-rate :model-value="item.rating / 2" disabled show-score text-color="#ff9900" />
          </div>
          <div v-if="item.genre?.length" class="movie-genre">
            <el-tag
              v-for="g in item.genre.slice(0, 2)"
              :key="g"
              size="small"
              type="primary"
              class="genre-tag"
            >
              {{ g }}
            </el-tag>
          </div>
          <div class="movie-tags">
            <el-tag
              v-for="tag in (item.tags || []).slice(0, 2)"
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
        <div class="movie-list-item" @click="openUrl(item.url)">
          <div class="movie-content">
            <div class="movie-header">
              <span class="movie-title">{{ item.title }}</span>
              <span v-if="item.year" class="movie-year">{{ item.year }}</span>
            </div>
            <div class="movie-info">
              <span v-if="item.director">导演: {{ item.director }}</span>
              <span v-if="item.rating">评分: {{ item.rating }}</span>
            </div>
          </div>
          <div class="movie-genre">
            <el-tag
              v-for="g in (item.genre || []).slice(0, 2)"
              :key="g"
              size="small"
              type="primary"
            >
              {{ g }}
            </el-tag>
          </div>
        </div>
      </template>

      <!-- 详情内容 -->
      <template #detail="{ item, loading }">
        <div class="movie-detail-content">
          <template v-if="loading">
            <el-skeleton :rows="6" animated />
          </template>
          <template v-else-if="item">
            <div class="detail-header">
              <h3 class="detail-title">{{ item.title }}</h3>
              <el-button v-if="item.url" type="primary" @click="openUrl(item.url)">
                <el-icon><Link /></el-icon>
                访问链接
              </el-button>
            </div>
            
            <div class="detail-info">
              <div v-if="item.year" class="info-item">
                <span class="label">年份:</span>
                <span>{{ item.year }}</span>
              </div>
              <div v-if="item.director" class="info-item">
                <span class="label">导演:</span>
                <span>{{ item.director }}</span>
              </div>
              <div v-if="item.rating" class="info-item">
                <span class="label">评分:</span>
                <el-rate :model-value="item.rating / 2" disabled show-score text-color="#ff9900" />
              </div>
            </div>

            <div v-if="item.genre?.length" class="detail-genre">
              <h4>类型</h4>
              <el-tag
                v-for="g in item.genre"
                :key="g"
                size="small"
                type="primary"
                style="margin-right: 8px; margin-bottom: 8px;"
              >
                {{ g }}
              </el-tag>
            </div>

            <el-divider v-if="item.comment" />
            <div v-if="item.comment" class="detail-comment">
              <h4>观后感</h4>
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
import { Link } from '@element-plus/icons-vue'

const notesStore = useNotesStore()

const selectedGenre = ref('')
const selectedTag = ref('')
const genres = ref<string[]>([])
const tags = ref<string[]>([])

const filteredMovies = computed(() => {
  let result = notesStore.movies
  
  if (selectedGenre.value) {
    result = result.filter(m => m.genre?.includes(selectedGenre.value))
  }
  
  if (selectedTag.value) {
    result = result.filter(m => m.tags?.includes(selectedTag.value))
  }
  
  return result
})

const viewConfig = computed(() => ({
  titleField: 'title',
  descField: 'comment',
  tagField: 'tags',
  cardType: 'movie' as const,
  gridColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
  columns: [
    { prop: 'title', label: '标题', minWidth: 200, sortable: true },
    { prop: 'year', label: '年份', width: 80 },
    { prop: 'director', label: '导演', width: 120 },
    { prop: 'rating', label: '评分', width: 100 },
    { prop: 'tags', label: '标签', width: 200 },
  ],
  chartType: 'bar' as const,
  chartGroupBy: 'genre',
  kanbanColumns: getKanbanColumns(),
}))

function getKanbanColumns() {
  const allGenres = notesStore.movies.flatMap(m => m.genre || [])
  const uniqueGenres = [...new Set(allGenres)]
  return uniqueGenres.slice(0, 10).map(genre => ({
    id: genre,
    title: genre,
    statusValue: genre,
  }))
}

function openUrl(url?: string) {
  if (url) {
    window.open(url, '_blank')
  }
}

async function loadMovieDetail(item: any) {
  if (item.id) {
    return await notesStore.getMovieDetail(item.id)
  }
  return item
}

async function handleItemClick(item: any) {
  console.log('Movie clicked:', item)
}

async function handleFilterChange() {
  // 筛选逻辑已通过 computed 实现
}

async function loadFilterOptions() {
  genres.value = await notesStore.getMovieGenres()
  tags.value = await notesStore.getMovieTags()
}

onMounted(async () => {
  await notesStore.loadNotes()
  await loadFilterOptions()
})
</script>

<style scoped>
.movies-page {
  padding: 0;
}

.filter-section {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
}

.movie-card {
  padding: 20px;
  border-radius: 12px;
  background: linear-gradient(135deg, #f4f4f5 0%, #fff 100%);
  border: 1px solid #d3d4d6;
  cursor: pointer;
  transition: all 0.3s;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.movie-card:hover {
  border-color: #909399;
  box-shadow: 0 4px 16px rgba(144, 147, 153, 0.15);
  transform: translateY(-2px);
}

.movie-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 12px;
}

.movie-title {
  font-weight: 600;
  font-size: 16px;
  color: #303133;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.movie-year {
  font-size: 13px;
  color: #909399;
}

.movie-info {
  display: flex;
  gap: 12px;
  font-size: 13px;
  color: #909399;
  margin-bottom: 8px;
}

.movie-rating {
  margin: 8px 0;
}

.movie-genre {
  margin-bottom: 8px;
  display: flex;
  gap: 6px;
}

.movie-tags {
  margin-top: auto;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.movie-tags .tag {
  margin: 0;
}

.movie-list-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px;
  width: 100%;
  cursor: pointer;
}

.movie-list-item:hover {
  background: #f5f7fa;
}

.movie-list-item .movie-content {
  flex: 1;
  min-width: 0;
}

.movie-list-item .movie-header {
  margin-bottom: 8px;
}

.movie-list-item .movie-title {
  font-size: 15px;
  margin-right: 8px;
}

.movie-list-item .movie-genre {
  margin-top: 0;
  margin-left: 16px;
  flex-shrink: 0;
}

.movie-detail-content {
  max-height: 70vh;
  overflow: auto;
}

.detail-header {
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detail-title {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.detail-info {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-item .label {
  color: #909399;
  font-size: 14px;
}

.detail-genre h4,
.detail-comment h4,
.detail-tags h4 {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #909399;
}

.detail-comment p {
  margin: 0;
  font-size: 14px;
  line-height: 1.8;
  color: #606266;
}
</style>
