<template>
  <NoteViewContainer
    title="知识库"
    :items="notesStore.knowledge"
    :loading="notesStore.loading"
    empty-text="暂无知识"
    :available-views="['tree', 'grid', 'list', 'table', 'chart', 'kanban']"
    default-view="tree"
    :view-config="viewConfig"
    :load-detail="loadKnowledgeDetail"
    @item-click="handleItemClick"
  >
    <!-- TreeView 右侧内容 -->
    <template #treeContent="{ item }">
      <div class="knowledge-tree-content">
        <div class="content-header">
          <h2>{{ selectedDetail?.title || item.title || item.label }}</h2>
          <el-tag v-if="selectedDetail">{{ selectedDetail.category }} / {{ selectedDetail.topic }}</el-tag>
        </div>
        <el-divider />
        <div class="content-body" v-if="selectedDetail">
          <MdRenderer :content="selectedDetail.content || ''" />
        </div>
        <div class="content-body" v-else>
          <el-skeleton :rows="6" animated />
        </div>
      </div>
    </template>

    <!-- GridView 自定义卡片 -->
    <template #gridItem="{ item }">
      <div class="knowledge-card">
        <div class="card-icon">
          <el-icon size="32"><Collection /></el-icon>
        </div>
        <div class="card-content">
          <h3>{{ item.title }}</h3>
          <p class="card-category">{{ item.category }} / {{ item.topic }}</p>
        </div>
      </div>
    </template>

    <!-- ListView 自定义项 -->
    <template #listItem="{ item }">
      <div class="knowledge-list-item">
        <el-icon><Collection /></el-icon>
        <div class="list-info">
          <h3>{{ item.title }}</h3>
          <p>{{ item.category }} / {{ item.topic }}</p>
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
      <div class="knowledge-detail-content">
        <template v-if="loading">
          <el-skeleton :rows="10" animated />
        </template>
        <template v-else-if="item">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="标题">{{ item.title }}</el-descriptions-item>
            <el-descriptions-item label="分类">{{ item.category }}</el-descriptions-item>
            <el-descriptions-item label="主题">{{ item.topic }}</el-descriptions-item>
          </el-descriptions>
          <el-divider content-position="left">知识内容</el-divider>
          <MdRenderer :content="item.content || ''" />
        </template>
      </div>
    </template>
  </NoteViewContainer>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useNotesStore } from '@/stores/notes'
import { useTree } from '@/composables/useTree'
import NoteViewContainer from '@/components/dashboard/widgets/NoteViewContainer.vue'
import MdRenderer from '@/components/MdRenderer.vue'
import { Collection } from '@element-plus/icons-vue'
import type { KnowledgeNote } from '@/types/note'

const notesStore = useNotesStore()
const { treeData, setRawData } = useTree()
const selectedDetail = ref<KnowledgeNote | null>(null)

// 加载树形数据
async function loadTreeData() {
  const rawData = await notesStore.getKnowledgeTree()
  setRawData(rawData)
}

// 视图配置
const viewConfig = computed(() => ({
  titleField: 'title',
  descField: 'description',
  tagField: 'category',
  statusField: 'category',
  timeField: 'updatedAt',
  cardType: 'book' as const,
  gridColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
  columns: [
    { prop: 'title', label: '标题', minWidth: 200, sortable: true },
    { prop: 'category', label: '分类', width: 120 },
    { prop: 'topic', label: '主题', width: 120 },
    { prop: 'updatedAt', label: '更新时间', width: 160, type: 'date' as const },
  ],
  treeData: treeData.value,
  chartType: 'bar' as const,
  chartGroupBy: 'category',
  kanbanColumns: [
    { id: 'frontend', title: '前端', statusValue: '前端' },
    { id: 'backend', title: '后端', statusValue: '后端' },
    { id: 'devops', title: '运维', statusValue: '运维' },
    { id: 'other', title: '其他', statusValue: '其他' },
  ],
}))

// 加载详情的函数
async function loadKnowledgeDetail(item: any) {
  if (item.category && item.topic && (item.label || item.title)) {
    return await notesStore.getKnowledgeDetailByQuery(
      item.category,
      item.topic,
      item.label || item.title
    )
  }
  return item
}

// 处理点击（树状视图）
async function handleItemClick(item: any) {
  if (item.category && item.topic && (item.label || item.title)) {
    const detail = await notesStore.getKnowledgeDetailByQuery(
      item.category,
      item.topic,
      item.label || item.title
    )
    if (detail) {
      selectedDetail.value = detail
    }
  }
}

onMounted(async () => {
  await loadTreeData()
})
</script>

<style scoped>
/* TreeView 右侧内容 */
.knowledge-tree-content {
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

/* GridView 卡片 */
.knowledge-card {
  padding: 16px;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #e4e7ed;
  cursor: pointer;
  transition: all 0.3s;
}

.knowledge-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.1);
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: linear-gradient(135deg, #36d1dc 0%, #5b86e5 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-bottom: 12px;
}

.card-content h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #303133;
}

.card-category {
  margin: 0;
  font-size: 12px;
  color: #909399;
}

/* ListView 项 */
.knowledge-list-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  width: 100%;
  cursor: pointer;
}

.knowledge-list-item:hover {
  background: #f5f7fa;
}

.list-info {
  flex: 1;
  min-width: 0;
}

.list-info h3 {
  margin: 0 0 4px 0;
  font-size: 14px;
  color: #303133;
}

.list-info p {
  margin: 0;
  font-size: 12px;
  color: #909399;
}

/* 详情内容 */
.knowledge-detail-content {
  padding: 0 16px;
}
</style>
