<template>
  <div class="anime-page">
    <!-- 筛选区域 -->
    <div class="filter-section">
      <el-select
        v-model="selectedStatus"
        placeholder="按状态筛选"
        clearable
        @change="handleFilterChange"
      >
        <el-option
          v-for="status in statuses"
          :key="status"
          :label="status"
          :value="status"
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
      title="动漫收藏"
      :items="filteredAnime"
      :loading="notesStore.loading"
      empty-text="暂无收藏的动漫"
      :available-views="['grid', 'list', 'table', 'chart', 'kanban']"
      default-view="grid"
      :view-config="viewConfig"
      :load-detail="loadAnimeDetail"
      @item-click="handleItemClick"
    >
      <!-- GridView 自定义卡片 -->
      <template #gridItem="{ item }">
        <div class="anime-card" @click="openUrl(item.url)">
          <div class="anime-header">
            <span class="anime-title">{{ item.title }}</span>
            <el-tag v-if="item.status" size="small" :type="getStatusType(item.status)">{{ item.status }}</el-tag>
          </div>
          <div class="anime-info">
            <span v-if="item.year">{{ item.year }}</span>
            <span v-if="item.studio">{{ item.studio }}</span>
            <span v-if="item.episodes">{{ item.episodes }}集</span>
          </div>
          <div v-if="item.rating" class="anime-rating">
            <el-rate :model-value="item.rating / 2" disabled show-score text-color="#ff9900" />
          </div>
          <div class="anime-tags">
            <el-tag
              v-for="tag in (item.tags || []).slice(0, 3)"
              :key="tag"
              size="small"
              type="info"
              class="tag"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>
      </template>

      <!-- ListView 自定义项 -->
      <template #listItem="{ item }">
        <div class="anime-list-item" @click="openUrl(item.url)">
          <div class="anime-content">
            <div class="anime-header">
              <span class="anime-title">{{ item.title }}</span>
              <el-tag v-if="item.status" size="small" :type="getStatusType(item.status)">{{ item.status }}</el-tag>
            </div>
            <div class="anime-info">
              <span v-if="item.year">{{ item.year }}</span>
              <span v-if="item.studio">{{ item.studio }}</span>
              <span v-if="item.episodes">{{ item.episodes }}集</span>
              <span v-if="item.rating">评分: {{ item.rating }}</span>
            </div>
          </div>
          <div class="anime-tags">
            <el-tag
              v-for="tag in (item.tags || []).slice(0, 2)"
              :key="tag"
              size="small"
              type="info"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>
      </template>

      <!-- 详情内容 -->
      <template #detail="{ item, loading }">
        <div class="anime-detail-content">
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
              <div v-if="item.studio" class="info-item">
                <span class="label">制作公司:</span>
                <span>{{ item.studio }}</span>
              </div>
              <div v-if="item.episodes" class="info-item">
                <span class="label">集数:</span>
                <span>{{ item.episodes }}</span>
              </div>
              <div v-if="item.status" class="info-item">
                <span class="label">状态:</span>
                <el-tag size="small" :type="getStatusType(item.status)">{{ item.status }}</el-tag>
              </div>
              <div v-if="item.rating" class="info-item">
                <span class="label">评分:</span>
                <el-rate :model-value="item.rating / 2" disabled show-score text-color="#ff9900" />
              </div>
            </div>

            <el-divider v-if="item.comment" />
            <div v-if="item.comment" class="detail-comment">
              <h4>评论</h4>
              <p>{{ item.comment }}</p>
            </div>

            <div class="detail-tags">
              <h4>标签</h4>
              <el-tag
                v-for="tag in item.tags"
                :key="tag"
                size="small"
                type="info"
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

const selectedStatus = ref('')
const selectedTag = ref('')
const statuses = ref<string[]>([])
const tags = ref<string[]>([])

const filteredAnime = computed(() => {
  let result = notesStore.anime
  
  if (selectedStatus.value) {
    result = result.filter(a => a.status === selectedStatus.value)
  }
  
  if (selectedTag.value) {
    result = result.filter(a => a.tags?.includes(selectedTag.value))
  }
  
  return result
})

const viewConfig = computed(() => ({
  titleField: 'title',
  descField: 'comment',
  tagField: 'tags',
  cardType: 'anime' as const,
  gridColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
  columns: [
    { prop: 'title', label: '标题', minWidth: 200, sortable: true },
    { prop: 'year', label: '年份', width: 80 },
    { prop: 'studio', label: '制作公司', width: 120 },
    { prop: 'status', label: '状态', width: 80 },
    { prop: 'rating', label: '评分', width: 100 },
  ],
  chartType: 'bar' as const,
  chartGroupBy: 'status',
  kanbanColumns: getKanbanColumns(),
}))

function getKanbanColumns() {
  const uniqueStatuses = [...new Set(notesStore.anime.map(a => a.status || '未知'))]
  return uniqueStatuses.slice(0, 10).map(status => ({
    id: status,
    title: status,
    statusValue: status,
  }))
}

function getStatusType(status: string): 'success' | 'warning' | 'danger' | 'info' | 'primary' {
  const typeMap: Record<string, 'success' | 'warning' | 'danger' | 'info' | 'primary'> = {
    '完结': 'success',
    '连载': 'warning',
    '待看': 'info',
    '已看完': 'success',
    '追番中': 'warning',
    '想看': 'info',
  }
  return typeMap[status] || 'info'
}

function openUrl(url?: string) {
  if (url) {
    window.open(url, '_blank')
  }
}

async function loadAnimeDetail(item: any) {
  if (item.id) {
    return await notesStore.getAnimeDetail(item.id)
  }
  return item
}

async function handleItemClick(item: any) {
  console.log('Anime clicked:', item)
}

async function handleFilterChange() {
  // 筛选逻辑已通过 computed 实现
}

async function loadFilterOptions() {
  statuses.value = await notesStore.getAnimeStatuses()
  tags.value = await notesStore.getAnimeTags()
}

onMounted(async () => {
  await notesStore.loadNotes()
  await loadFilterOptions()
})
</script>

<style scoped>
.anime-page {
  padding: 0;
}

.filter-section {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
}

.anime-card {
  padding: 20px;
  border-radius: 12px;
  background: linear-gradient(135deg, #fef0f0 0%, #fff 100%);
  border: 1px solid #fbc4c4;
  cursor: pointer;
  transition: all 0.3s;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.anime-card:hover {
  border-color: #f78989;
  box-shadow: 0 4px 16px rgba(247, 137, 137, 0.15);
  transform: translateY(-2px);
}

.anime-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 12px;
}

.anime-title {
  font-weight: 600;
  font-size: 16px;
  color: #303133;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.anime-info {
  display: flex;
  gap: 12px;
  font-size: 13px;
  color: #909399;
  margin-bottom: 8px;
}

.anime-rating {
  margin: 8px 0;
}

.anime-tags {
  margin-top: auto;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.anime-tags .tag {
  margin: 0;
}

.anime-list-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px;
  width: 100%;
  cursor: pointer;
}

.anime-list-item:hover {
  background: #f5f7fa;
}

.anime-list-item .anime-content {
  flex: 1;
  min-width: 0;
}

.anime-list-item .anime-header {
  margin-bottom: 8px;
}

.anime-list-item .anime-title {
  font-size: 15px;
  margin-right: 8px;
}

.anime-list-item .anime-tags {
  margin-top: 0;
  margin-left: 16px;
  flex-shrink: 0;
}

.anime-detail-content {
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
