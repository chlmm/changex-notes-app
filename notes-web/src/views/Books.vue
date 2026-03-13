<template>
  <NoteViewContainer
    title="书籍学习"
    :items="notesStore.books"
    :loading="notesStore.loading"
    empty-text="暂无书籍"
    :available-views="['grid', 'list', 'table', 'tree', 'chart', 'kanban']"
    default-view="grid"
    :view-config="viewConfig"
    :load-detail="loadBookDetail"
    :detail-route="'/books/:type/:title'"
    @item-click="handleItemClick"
  >
    <!-- GridView 自定义卡片 -->
    <template #gridItem="{ item }">
      <NoteCard
        :note="item"
        type="book"
      />
    </template>
    
    <!-- ListView 自定义项 -->
    <template #listItem="{ item }">
      <div class="book-list-item">
        <div class="book-cover">
          <el-icon size="32"><Reading /></el-icon>
        </div>
        <div class="book-info">
          <h3>{{ item.title }}</h3>
          <div class="book-meta">
            <el-tag size="small">{{ item.type === 'novel' ? '小说' : '非虚构' }}</el-tag>
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
      <div class="book-tree-content">
        <div class="content-header">
          <h2>{{ selectedDetail?.title || item.title }}</h2>
          <div class="content-tags">
            <el-tag>{{ selectedDetail?.type === 'novel' ? '小说' : '非虚构' }}</el-tag>
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
            <p class="book-desc">暂无内容</p>
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
      <div class="book-detail-content">
        <template v-if="loading">
          <el-skeleton :rows="10" animated />
        </template>
        <template v-else-if="item">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="书名">{{ item.title }}</el-descriptions-item>
            <el-descriptions-item label="类型">{{ item.type === 'novel' ? '小说' : '非虚构' }}</el-descriptions-item>
            <el-descriptions-item label="状态">{{ item.status || '未开始' }}</el-descriptions-item>
            <el-descriptions-item label="作者" v-if="item.author">{{ item.author }}</el-descriptions-item>
          </el-descriptions>
          <el-divider content-position="left">笔记内容</el-divider>
          <MdRenderer :content="item.content || ''" />
        </template>
      </div>
    </template>
  </NoteViewContainer>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useNotesStore } from '@/stores/notes'
import NoteViewContainer from '@/components/dashboard/widgets/NoteViewContainer.vue'
import NoteCard from '@/components/NoteCard.vue'
import MdRenderer from '@/components/MdRenderer.vue'
import { Reading } from '@element-plus/icons-vue'
import type { BookNote } from '@/types/note'

const router = useRouter()
const notesStore = useNotesStore()
const selectedDetail = ref<BookNote | null>(null)
const loadingDetail = ref(false)

// 树形数据：按类型分组
const treeData = computed(() => {
  const books = notesStore.books
  const novelBooks = books.filter(b => b.type === 'novel')
  const nonFictionBooks = books.filter(b => b.type === 'non-fiction')
  
  return [
    {
      id: 'novel',
      label: '小说',
      isCategory: true,
      children: novelBooks.map(book => ({
        id: book.path || book.title,
        label: book.title,
        isLeaf: true,
        ...book
      }))
    },
    {
      id: 'non-fiction',
      label: '非虚构',
      isCategory: true,
      children: nonFictionBooks.map(book => ({
        id: book.path || book.title,
        label: book.title,
        isLeaf: true,
        ...book
      }))
    }
  ]
})

// 视图配置
const viewConfig = computed(() => ({
  titleField: 'title',
  descField: 'description',
  tagField: 'type',
  statusField: 'status',
  timeField: 'updatedAt',
  typeField: 'type',
  cardType: 'book' as const,
  gridColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
  columns: [
    { prop: 'title', label: '书名', minWidth: 200, sortable: true },
    { 
      prop: 'type', 
      label: '类型', 
      width: 100,
      formatter: (val: string) => val === 'novel' ? '小说' : '非虚构'
    },
    { 
      prop: 'status', 
      label: '状态', 
      width: 100, 
      type: 'tag' as const,
      getTagType: (row: any): 'success' | 'warning' | 'danger' | 'info' | 'primary' => {
        const map: Record<string, 'success' | 'warning' | 'danger' | 'info' | 'primary'> = {
          '未开始': 'info',
          '阅读中': 'warning',
          '已完成': 'success',
          '已解读': 'success',
        }
        return map[row.status] || 'info'
      }
    },
    { prop: 'progress', label: '进度', width: 120, type: 'progress' as const },
    { prop: 'updatedAt', label: '更新时间', width: 160, type: 'date' as const },
  ],
  kanbanColumns: [
    { id: 'notStarted', title: '未开始', statusValue: '未开始' },
    { id: 'reading', title: '阅读中', statusValue: '阅读中' },
    { id: 'finished', title: '已完成', statusValue: '已完成' },
    { id: 'reviewed', title: '已解读', statusValue: '已解读' },
  ],
  chartType: 'bar' as const,
  chartGroupBy: 'type',
  treeData: treeData.value,
}))

// 获取状态类型
function getStatusType(status?: string): 'success' | 'warning' | 'danger' | 'info' | 'primary' {
  const map: Record<string, 'success' | 'warning' | 'danger' | 'info' | 'primary'> = {
    '未开始': 'info',
    '阅读中': 'warning',
    '已完成': 'success',
    '已解读': 'success',
  }
  return map[status || ''] || 'info'
}

// 加载详情的函数
async function loadBookDetail(item: any) {
  if (item.type && item.title) {
    return await notesStore.getBookDetail(item.type, item.title)
  }
  return item
}

// 跳转到详情页
function goToDetail(book: BookNote) {
  router.push(`/books/${book.type}/${encodeURIComponent(book.title)}`)
}

// 处理树状视图点击事件
async function handleItemClick(item: any) {
  if (item.type && item.title && item.isLeaf) {
    loadingDetail.value = true
    const detail = await notesStore.getBookDetail(item.type, item.title)
    if (detail) {
      selectedDetail.value = detail
    }
    loadingDetail.value = false
  }
}

</script>

<style scoped>
.book-list-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  width: 100%;
}

.book-cover {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.book-info {
  flex: 1;
  min-width: 0;
}

.book-info h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.book-meta {
  display: flex;
  gap: 8px;
}

/* TreeView 右侧内容样式 */
.book-tree-content {
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

.book-desc {
  margin: 0 0 16px 0;
}

.content-footer {
  margin-top: 24px;
}

/* 详情内容 */
.book-detail-content {
  padding: 0 16px;
}
</style>
