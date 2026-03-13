<template>
  <div class="games-page">
    <!-- 筛选区域 -->
    <div class="filter-section">
      <el-select
        v-model="selectedPlatform"
        placeholder="按平台筛选"
        clearable
        @change="handleFilterChange"
      >
        <el-option
          v-for="platform in platforms"
          :key="platform"
          :label="platform"
          :value="platform"
        />
      </el-select>
      <el-select
        v-model="selectedStatus"
        placeholder="按状态筛选"
        clearable
        style="margin-left: 12px;"
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
      title="游戏收藏"
      :items="filteredGames"
      :loading="notesStore.loading"
      empty-text="暂无收藏的游戏"
      :available-views="['grid', 'list', 'table', 'chart', 'kanban']"
      default-view="grid"
      :view-config="viewConfig"
      :load-detail="loadGameDetail"
      @item-click="handleItemClick"
    >
      <!-- GridView 自定义卡片 -->
      <template #gridItem="{ item }">
        <div class="game-card" @click="openUrl(item.url)">
          <div class="game-header">
            <span class="game-title">{{ item.title }}</span>
            <el-tag v-if="item.status" size="small" :type="getStatusType(item.status)">{{ item.status }}</el-tag>
          </div>
          <div class="game-info">
            <span v-if="item.platform">{{ item.platform }}</span>
            <span v-if="item.developer">{{ item.developer }}</span>
            <span v-if="item.year">{{ item.year }}</span>
          </div>
          <div v-if="item.rating" class="game-rating">
            <el-rate :model-value="item.rating / 2" disabled show-score text-color="#ff9900" />
          </div>
          <div class="game-tags">
            <el-tag
              v-for="tag in (item.tags || []).slice(0, 3)"
              :key="tag"
              size="small"
              type="success"
              class="tag"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>
      </template>

      <!-- ListView 自定义项 -->
      <template #listItem="{ item }">
        <div class="game-list-item" @click="openUrl(item.url)">
          <div class="game-content">
            <div class="game-header">
              <span class="game-title">{{ item.title }}</span>
              <el-tag v-if="item.status" size="small" :type="getStatusType(item.status)">{{ item.status }}</el-tag>
            </div>
            <div class="game-info">
              <span v-if="item.platform">{{ item.platform }}</span>
              <span v-if="item.developer">{{ item.developer }}</span>
              <span v-if="item.rating">评分: {{ item.rating }}</span>
            </div>
          </div>
          <div class="game-tags">
            <el-tag
              v-for="tag in (item.tags || []).slice(0, 2)"
              :key="tag"
              size="small"
              type="success"
            >
              {{ tag }}
            </el-tag>
          </div>
        </div>
      </template>

      <!-- 详情内容 -->
      <template #detail="{ item, loading }">
        <div class="game-detail-content">
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
              <div v-if="item.platform" class="info-item">
                <span class="label">平台:</span>
                <span>{{ item.platform }}</span>
              </div>
              <div v-if="item.developer" class="info-item">
                <span class="label">开发商:</span>
                <span>{{ item.developer }}</span>
              </div>
              <div v-if="item.year" class="info-item">
                <span class="label">年份:</span>
                <span>{{ item.year }}</span>
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
                type="success"
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

const selectedPlatform = ref('')
const selectedStatus = ref('')
const selectedTag = ref('')
const platforms = ref<string[]>([])
const statuses = ref<string[]>([])
const tags = ref<string[]>([])

const filteredGames = computed(() => {
  let result = notesStore.games
  
  if (selectedPlatform.value) {
    result = result.filter(g => g.platform === selectedPlatform.value)
  }
  
  if (selectedStatus.value) {
    result = result.filter(g => g.status === selectedStatus.value)
  }
  
  if (selectedTag.value) {
    result = result.filter(g => g.tags?.includes(selectedTag.value))
  }
  
  return result
})

const viewConfig = computed(() => ({
  titleField: 'title',
  descField: 'comment',
  tagField: 'tags',
  cardType: 'game' as const,
  gridColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
  columns: [
    { prop: 'title', label: '游戏名称', minWidth: 200, sortable: true },
    { prop: 'platform', label: '平台', width: 100 },
    { prop: 'developer', label: '开发商', width: 120 },
    { prop: 'status', label: '状态', width: 80 },
    { prop: 'rating', label: '评分', width: 100 },
  ],
  chartType: 'bar' as const,
  chartGroupBy: 'platform',
  kanbanColumns: getKanbanColumns(),
}))

function getKanbanColumns() {
  const uniqueStatuses = [...new Set(notesStore.games.map(g => g.status || '未知'))]
  return uniqueStatuses.slice(0, 10).map(status => ({
    id: status,
    title: status,
    statusValue: status,
  }))
}

function getStatusType(status: string): 'success' | 'warning' | 'danger' | 'info' | 'primary' {
  const typeMap: Record<string, 'success' | 'warning' | 'danger' | 'info' | 'primary'> = {
    '通关': 'success',
    '在玩': 'warning',
    '想玩': 'info',
    '已通关': 'success',
    '游玩中': 'warning',
  }
  return typeMap[status] || 'info'
}

function openUrl(url?: string) {
  if (url) {
    window.open(url, '_blank')
  }
}

async function loadGameDetail(item: any) {
  if (item.id) {
    return await notesStore.getGameDetail(item.id)
  }
  return item
}

async function handleItemClick(item: any) {
  console.log('Game clicked:', item)
}

async function handleFilterChange() {
  // 筛选逻辑已通过 computed 实现
}

async function loadFilterOptions() {
  platforms.value = await notesStore.getGamePlatforms()
  statuses.value = await notesStore.getGameStatuses()
  tags.value = await notesStore.getGameTags()
}

onMounted(async () => {
  await notesStore.loadNotes()
  await loadFilterOptions()
})
</script>

<style scoped>
.games-page {
  padding: 0;
}

.filter-section {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
}

.game-card {
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

.game-card:hover {
  border-color: #67c23a;
  box-shadow: 0 4px 16px rgba(103, 194, 58, 0.15);
  transform: translateY(-2px);
}

.game-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 12px;
}

.game-title {
  font-weight: 600;
  font-size: 16px;
  color: #303133;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.game-info {
  display: flex;
  gap: 12px;
  font-size: 13px;
  color: #909399;
  margin-bottom: 8px;
}

.game-rating {
  margin: 8px 0;
}

.game-tags {
  margin-top: auto;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.game-tags .tag {
  margin: 0;
}

.game-list-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px;
  width: 100%;
  cursor: pointer;
}

.game-list-item:hover {
  background: #f5f7fa;
}

.game-list-item .game-content {
  flex: 1;
  min-width: 0;
}

.game-list-item .game-header {
  margin-bottom: 8px;
}

.game-list-item .game-title {
  font-size: 15px;
  margin-right: 8px;
}

.game-list-item .game-tags {
  margin-top: 0;
  margin-left: 16px;
  flex-shrink: 0;
}

.game-detail-content {
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
