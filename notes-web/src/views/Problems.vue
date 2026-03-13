<template>
  <NoteViewContainer
    title="问题与解决"
    :items="notesStore.problems"
    :loading="notesStore.loading"
    empty-text="暂无问题"
    :available-views="['grid', 'list', 'table', 'tree', 'chart', 'kanban']"
    default-view="grid"
    :view-config="viewConfig"
    :load-detail="loadProblemDetail"
    @item-click="handleItemClick"
  >
    <!-- GridView 自定义卡片 -->
    <template #gridItem="{ item }">
      <div class="problem-card">
        <div class="problem-header">
          <el-tag :type="item.type === '问题与解决' ? 'warning' : 'success'" size="small">
            {{ item.type }}
          </el-tag>
        </div>
        <h3>{{ item.title }}</h3>
        <p class="problem-preview">{{ getPreview(item.content) }}</p>
      </div>
    </template>

    <!-- ListView 自定义项 -->
    <template #listItem="{ item }">
      <div class="problem-list-item">
        <el-tag :type="item.type === '问题与解决' ? 'warning' : 'success'" size="small">
          {{ item.type }}
        </el-tag>
        <div class="problem-info">
          <h3>{{ item.title }}</h3>
          <p class="problem-preview">{{ getPreview(item.content) }}</p>
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
      <div class="problem-tree-content">
        <div class="content-header">
          <h2>{{ selectedDetail?.title || item.title || item.label }}</h2>
          <el-tag v-if="selectedDetail" :type="selectedDetail.type === '问题与解决' ? 'warning' : 'success'">
            {{ selectedDetail.type }}
          </el-tag>
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
            <p class="empty-desc">暂无内容</p>
          </template>
        </div>
      </div>
    </template>

    <!-- 详情内容（抽屉/对话框） -->
    <template #detail="{ item, loading }">
      <div class="problem-detail-content">
        <template v-if="loading">
          <el-skeleton :rows="10" animated />
        </template>
        <template v-else-if="item">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="类型">{{ item.type }}</el-descriptions-item>
          </el-descriptions>
          <el-divider content-position="left">问题描述与解决方案</el-divider>
          <MdRenderer :content="item.content || ''" />
        </template>
      </div>
    </template>
  </NoteViewContainer>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useNotesStore } from '@/stores/notes'
import NoteViewContainer from '@/components/dashboard/widgets/NoteViewContainer.vue'
import MdRenderer from '@/components/MdRenderer.vue'
import type { ProblemNote } from '@/types/note'

const notesStore = useNotesStore()
const selectedDetail = ref<ProblemNote | null>(null)
const loadingDetail = ref(false)

// 树形数据：按类型分组
const treeData = computed(() => {
  const problems = notesStore.problems
  const problemList = problems.filter(p => p.type === '问题与解决')
  const requirementList = problems.filter(p => p.type === '需求与解决')
  
  return [
    {
      id: 'problem',
      label: '问题与解决',
      isCategory: true,
      children: problemList.map(p => ({
        id: p.path || p.title,
        label: p.title,
        isLeaf: true,
        ...p
      }))
    },
    {
      id: 'requirement',
      label: '需求与解决',
      isCategory: true,
      children: requirementList.map(p => ({
        id: p.path || p.title,
        label: p.title,
        isLeaf: true,
        ...p
      }))
    }
  ]
})

// 视图配置
const viewConfig = computed(() => ({
  titleField: 'title',
  descField: 'content',
  tagField: 'type',
  statusField: 'type',
  timeField: 'updatedAt',
  cardType: 'book' as const,
  gridColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
  columns: [
    { prop: 'title', label: '标题', minWidth: 200, sortable: true },
    { prop: 'type', label: '类型', width: 120 },
    { prop: 'updatedAt', label: '更新时间', width: 160, type: 'date' as const },
  ],
  treeData: treeData.value,
  chartType: 'pie' as const,
  chartGroupBy: 'type',
  kanbanColumns: [
    { id: 'problem', title: '问题与解决', statusValue: '问题与解决' },
    { id: 'requirement', title: '需求与解决', statusValue: '需求与解决' },
  ],
}))

function getPreview(content?: string): string {
  if (!content) return '暂无内容'
  const lines = content.split('\n').filter(l => l.trim())
  const preview = lines.find(l => !l.startsWith('#')) || lines[0] || ''
  return preview.length > 100 ? preview.slice(0, 100) + '...' : preview
}

// 加载详情的函数
async function loadProblemDetail(item: any) {
  if (item.path) {
    return await notesStore.getProblemDetail(item.path)
  }
  return item
}

// 处理树状视图点击事件
async function handleItemClick(item: any) {
  if (item.path && item.isLeaf) {
    loadingDetail.value = true
    const detail = await notesStore.getProblemDetail(item.path)
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
/* GridView 卡片 */
.problem-card {
  padding: 16px;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #e4e7ed;
  cursor: pointer;
  transition: all 0.3s;
}

.problem-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.1);
}

.problem-header {
  margin-bottom: 8px;
}

.problem-card h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.problem-preview {
  margin: 0;
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* ListView 项 */
.problem-list-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  width: 100%;
  cursor: pointer;
}

.problem-list-item:hover {
  background: #f5f7fa;
}

.problem-info {
  flex: 1;
  min-width: 0;
}

.problem-info h3 {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #303133;
}

.problem-info .problem-preview {
  font-size: 12px;
}

/* TreeView 右侧内容 */
.problem-tree-content {
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

.content-body {
  line-height: 1.8;
  color: #606266;
}

.empty-desc {
  color: #909399;
}

/* 详情内容 */
.problem-detail-content {
  padding: 0 16px;
  max-height: 70vh;
  overflow: auto;
}
</style>
