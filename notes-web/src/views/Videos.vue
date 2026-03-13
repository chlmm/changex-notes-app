<template>
  <NoteViewContainer
    title="视频学习"
    :items="notesStore.videos"
    :loading="notesStore.loading"
    empty-text="暂无视频"
    :available-views="['grid', 'list', 'table', 'tree', 'chart', 'kanban']"
    default-view="grid"
    :view-config="viewConfig"
    :load-detail="loadVideoDetail"
    :detail-route="'/videos/:id'"
    @item-click="handleItemClick"
  >
    <!-- GridView 自定义卡片 -->
    <template #gridItem="{ item }">
      <NoteCard
        :note="item"
        type="video"
      />
    </template>

    <!-- ListView 自定义项 -->
    <template #listItem="{ item }">
      <div class="video-list-item">
        <div class="video-cover">
          <el-icon size="32"><VideoPlay /></el-icon>
        </div>
        <div class="video-info">
          <h3>{{ item.title }}</h3>
          <div class="video-meta">
            <el-tag size="small">{{ item.platform || 'bilibili' }}</el-tag>
            <el-tag size="small" :type="getStatusType(item.status)">
              {{ item.status || '未开始' }}
            </el-tag>
          </div>
        </div>
      </div>
    </template>

    <!-- Table 操作列 -->
    <template #tableOperation>
      <el-button size="small" type="primary" link>
        查看
      </el-button>
    </template>

    <!-- TreeView 右侧内容 -->
    <template #treeContent="{ item }">
      <div class="video-tree-content">
        <div class="content-header">
          <h2>{{ selectedDetail?.title || item.title }}</h2>
          <div class="content-tags">
            <el-tag>{{ selectedDetail?.platform || 'bilibili' }}</el-tag>
            <el-tag :type="getStatusType(selectedDetail?.status || '未开始')">
              {{ selectedDetail?.status || '未开始' }}
            </el-tag>
          </div>
        </div>
        <el-divider />
        <div class="content-body">
          <template v-if="loadingDetail">
            <el-skeleton :rows="6" animated />
          </template>
          <template v-else-if="selectedDetail?.content">
            <MdRenderer :content="selectedDetail.content || ''" />
          </template>
          <template v-else>
            <p class="video-desc">暂无内容</p>
          </template>
        </div>
        <div class="content-footer" v-if="selectedDetail">
          <el-button type="primary" @click="goToDetail(selectedDetail)">
            查看完整详情
          </el-button>
        </div>
      </div>
    </template>

    <!-- 详情内容（抽屉/对话框） -->
    <template #detail="{ item, loading }">
      <div class="video-detail-content">
        <template v-if="loading">
          <el-skeleton :rows="10" animated />
        </template>
        <template v-else-if="item">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="标题">{{ item.title }}</el-descriptions-item>
            <el-descriptions-item label="平台">{{ item.platform || 'bilibili' }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ item.status || '未开始' }}</el-descriptions-item>
          </el-descriptions>
          <el-divider content-position="left">笔记内容</el-divider>
          <MdRenderer :content="item.content || ''" />
        </template>
      </div>
    </template>
  </NoteViewContainer>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useNotesStore } from '@/stores/notes'
import NoteViewContainer from '@/components/dashboard/widgets/NoteViewContainer.vue'
import NoteCard from '@/components/NoteCard.vue'
import MdRenderer from '@/components/MdRenderer.vue'
import { VideoPlay } from '@element-plus/icons-vue'
import type { VideoNote } from '@/types/note'

const router = useRouter()
const notesStore = useNotesStore()
const selectedDetail = ref<VideoNote | null>(null)
const loadingDetail = ref(false)

// 树形数据：按状态分组
const treeData = computed(() => {
  const videos = notesStore.videos
  const statusGroups: Record<string, VideoNote[]> = {
    '未开始': [],
    '阅读中': [],
    '已解读': [],
  }
  
  videos.forEach(v => {
    const status = v.status || '未开始'
    if (statusGroups[status]) {
      statusGroups[status].push(v)
    }
  })
  
  return Object.entries(statusGroups)
    .filter(([, items]) => items.length > 0)
    .map(([status, items]) => ({
      id: status,
      label: `${status} (${items.length})`,
      isCategory: true,
      children: items.map(video => ({
        id: video.videoId,
        label: video.title,
        isLeaf: true,
        ...video
      }))
    }))
})

// 视图配置
const viewConfig = computed(() => ({
  titleField: 'title',
  descField: 'description',
  tagField: 'platform',
  statusField: 'status',
  timeField: 'created',
  cardType: 'video' as const,
  gridColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
  columns: [
    { prop: 'title', label: '标题', minWidth: 200, sortable: true },
    { prop: 'platform', label: '平台', width: 100 },
    { 
      prop: 'status', 
      label: '状态', 
      width: 100, 
      type: 'tag' as const,
      getTagType: (row: any): 'success' | 'warning' | 'danger' | 'info' | 'primary' => {
        const map: Record<string, 'success' | 'warning' | 'danger' | 'info' | 'primary'> = {
          '未开始': 'info',
          '阅读中': 'warning',
          '已解读': 'success',
        }
        return map[row.status] || 'info'
      }
    },
    { prop: 'created', label: '创建时间', width: 160, type: 'date' as const },
  ],
  kanbanColumns: [
    { id: 'notStarted', title: '未开始', statusValue: '未开始' },
    { id: 'watching', title: '观看中', statusValue: '阅读中' },
    { id: 'reviewed', title: '已解读', statusValue: '已解读' },
  ],
  chartType: 'bar' as const,
  chartGroupBy: 'platform',
  treeData: treeData.value,
}))

// 获取状态类型
function getStatusType(status?: string): 'success' | 'warning' | 'danger' | 'info' | 'primary' {
  const map: Record<string, 'success' | 'warning' | 'danger' | 'info' | 'primary'> = {
    '未开始': 'info',
    '阅读中': 'warning',
    '已解读': 'success',
  }
  return map[status || ''] || 'info'
}

// 加载详情的函数
async function loadVideoDetail(item: any) {
  if (item.videoId) {
    return await notesStore.getVideoDetail(item.videoId)
  }
  return item
}

// 跳转到详情页
function goToDetail(video: VideoNote) {
  router.push(`/videos/${video.videoId}`)
}

// 处理树状视图点击事件
async function handleItemClick(item: any) {
  if (item.videoId && item.isLeaf) {
    loadingDetail.value = true
    const detail = await notesStore.getVideoDetail(item.videoId)
    if (detail) {
      selectedDetail.value = detail
    }
    loadingDetail.value = false
  }
}

onMounted(() => {
  notesStore.loadNotes()
})
</script>

<style scoped>
/* ListView 项 */
.video-list-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  width: 100%;
}

.video-cover {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.video-info {
  flex: 1;
  min-width: 0;
}

.video-info h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.video-meta {
  display: flex;
  gap: 8px;
}

/* TreeView 右侧内容样式 */
.video-tree-content {
  padding: 0 16px;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.content-header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.content-tags {
  display: flex;
  gap: 8px;
}

.content-body {
  line-height: 1.8;
  color: #606266;
}

.video-desc {
  margin: 0 0 16px 0;
}

.content-footer {
  margin-top: 24px;
}

/* 详情内容 */
.video-detail-content {
  padding: 0 16px;
}
</style>
