<template>
  <div class="quotes-page">
    <!-- 筛选区域 -->
    <div class="filter-section">
      <el-select
        v-model="selectedAuthor"
        placeholder="按作者筛选"
        clearable
        @change="handleFilterChange"
      >
        <el-option
          v-for="author in authors"
          :key="author"
          :label="author"
          :value="author"
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
      title="金句收藏"
      :items="filteredQuotes"
      :loading="notesStore.loading"
      empty-text="暂无金句"
      :available-views="['grid', 'list', 'table', 'chart', 'kanban']"
      default-view="grid"
      :view-config="viewConfig"
      :load-detail="loadQuoteDetail"
      @item-click="handleItemClick"
    >
      <!-- GridView 自定义卡片 -->
      <template #gridItem="{ item }">
        <div class="quote-card">
          <div class="quote-content">
            <el-icon class="quote-icon"><ChatDotSquare /></el-icon>
            <p class="quote-text">{{ item.quote }}</p>
          </div>
          <div class="quote-meta">
            <span v-if="item.author" class="author">{{ item.author }}</span>
            <span v-if="item.source" class="source">《{{ item.source }}》</span>
          </div>
          <div class="quote-tags">
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
        <div class="quote-list-item">
          <div class="quote-content">
            <p class="quote-text">{{ item.quote }}</p>
            <div class="quote-meta">
              <span v-if="item.author" class="author">—— {{ item.author }}</span>
              <span v-if="item.source" class="source">《{{ item.source }}》</span>
            </div>
          </div>
          <div class="quote-tags">
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

      <!-- Table 操作列 -->
      <template #tableOperation>
        <el-button size="small" type="primary" link>
          查看
        </el-button>
      </template>

      <!-- 详情内容（抽屉/对话框） -->
      <template #detail="{ item, loading }">
        <div class="quote-detail-content">
          <template v-if="loading">
            <el-skeleton :rows="6" animated />
          </template>
          <template v-else-if="item">
            <div class="detail-quote">
              <el-icon class="quote-icon"><ChatDotSquare /></el-icon>
              <p class="quote-text">{{ item.quote }}</p>
              <div class="quote-meta">
                <span v-if="item.author" class="author">{{ item.author }}</span>
                <span v-if="item.source" class="source">《{{ item.source }}》</span>
              </div>
            </div>
            <el-divider v-if="item.comment" />
            <div v-if="item.comment" class="detail-comment">
              <h4>解读</h4>
              <p>{{ item.comment }}</p>
            </div>
            <div class="detail-tags">
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
import { ChatDotSquare } from '@element-plus/icons-vue'

const notesStore = useNotesStore()

// 筛选状态
const selectedAuthor = ref('')
const selectedTag = ref('')
const authors = ref<string[]>([])
const tags = ref<string[]>([])

// 筛选后的金句列表
const filteredQuotes = computed(() => {
  let result = notesStore.quotes
  
  if (selectedAuthor.value) {
    result = result.filter(q => q.author === selectedAuthor.value)
  }
  
  if (selectedTag.value) {
    result = result.filter(q => q.tags?.includes(selectedTag.value))
  }
  
  return result
})

// 视图配置
const viewConfig = computed(() => ({
  titleField: 'quote',
  descField: 'comment',
  tagField: 'tags',
  timeField: 'createdAt',
  cardType: 'quote' as const,
  gridColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
  columns: [
    { prop: 'quote', label: '金句', minWidth: 300, sortable: true },
    { prop: 'author', label: '作者', width: 120 },
    { prop: 'source', label: '出处', width: 150 },
    { prop: 'tags', label: '标签', width: 200 },
  ],
  chartType: 'bar' as const,
  chartGroupBy: 'author',
  kanbanColumns: getKanbanColumns(),
}))

// 动态生成看板列（按作者分组）
function getKanbanColumns() {
  const uniqueAuthors = [...new Set(notesStore.quotes.map(q => q.author || '未知'))]
  
  return uniqueAuthors.slice(0, 10).map(author => ({
    id: author,
    title: author,
    statusValue: author,
  }))
}

// 加载详情的函数
async function loadQuoteDetail(item: any) {
  if (item.id) {
    return await notesStore.getQuoteDetail(item.id)
  }
  return item
}

// 处理点击事件
async function handleItemClick(item: any) {
  console.log('Quote clicked:', item)
}

// 筛选变化时重新加载
async function handleFilterChange() {
  // 筛选逻辑已通过 computed 实现
}

// 加载作者和标签列表
async function loadFilterOptions() {
  authors.value = await notesStore.getQuoteAuthors()
  tags.value = await notesStore.getQuoteTags()
}

onMounted(async () => {
  await notesStore.loadNotes()
  await loadFilterOptions()
})
</script>

<style scoped>
.quotes-page {
  padding: 0;
}

.filter-section {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
}

/* GridView 卡片 */
.quote-card {
  padding: 20px;
  border-radius: 12px;
  background: linear-gradient(135deg, #f5f7fa 0%, #fff 100%);
  border: 1px solid #e4e7ed;
  cursor: pointer;
  transition: all 0.3s;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.quote-card:hover {
  border-color: #409eff;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.15);
  transform: translateY(-2px);
}

.quote-content {
  flex: 1;
}

.quote-icon {
  color: #c0c4cc;
  font-size: 24px;
  margin-bottom: 8px;
}

.quote-text {
  margin: 0;
  font-size: 15px;
  line-height: 1.8;
  color: #303133;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.quote-meta {
  margin-top: 16px;
  font-size: 13px;
  color: #909399;
}

.quote-meta .author {
  font-weight: 500;
  color: #606266;
}

.quote-meta .source {
  margin-left: 8px;
}

.quote-tags {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.quote-tags .tag {
  margin: 0;
}

/* ListView 项 */
.quote-list-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px;
  width: 100%;
  cursor: pointer;
}

.quote-list-item:hover {
  background: #f5f7fa;
}

.quote-list-item .quote-content {
  flex: 1;
  min-width: 0;
}

.quote-list-item .quote-text {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #303133;
  line-height: 1.6;
}

.quote-list-item .quote-meta {
  font-size: 12px;
  color: #909399;
  margin-top: 0;
}

.quote-list-item .quote-tags {
  margin-top: 0;
  margin-left: 16px;
  flex-shrink: 0;
}

/* 详情内容 */
.quote-detail-content {
  max-height: 70vh;
  overflow: auto;
}

.detail-quote {
  text-align: center;
  padding: 20px 0;
}

.detail-quote .quote-icon {
  font-size: 32px;
  color: #c0c4cc;
  margin-bottom: 16px;
}

.detail-quote .quote-text {
  font-size: 18px;
  line-height: 1.8;
  color: #303133;
  margin: 0;
}

.detail-quote .quote-meta {
  margin-top: 16px;
  font-size: 14px;
}

.detail-quote .quote-meta .author {
  font-weight: 500;
  color: #409eff;
}

.detail-comment {
  margin-top: 16px;
}

.detail-comment h4 {
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

.detail-tags {
  margin-top: 16px;
}
</style>
