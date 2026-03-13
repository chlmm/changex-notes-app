<template>
  <NoteViewContainer
    title="索引与收藏"
    :items="notesStore.indexNotes"
    :loading="notesStore.loading"
    empty-text="暂无索引"
    :available-views="['grid', 'list', 'table', 'tree', 'chart', 'kanban']"
    default-view="grid"
    :view-config="viewConfig"
    :load-detail="loadIndexDetail"
    @item-click="handleItemClick"
  >
    <!-- GridView 自定义卡片 -->
    <template #gridItem="{ item }">
      <div class="index-card">
        <div class="card-icon">
          <el-icon size="32"><Collection /></el-icon>
        </div>
        <div class="card-content">
          <h3>{{ item.title }}</h3>
          <p class="card-stats">{{ getPreview(item.content) }}</p>
        </div>
      </div>
    </template>

    <!-- ListView 自定义项 -->
    <template #listItem="{ item }">
      <div class="index-list-item">
        <el-icon size="20"><Collection /></el-icon>
        <div class="list-info">
          <h3>{{ item.title }}</h3>
          <p>{{ getPreview(item.content) }}</p>
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
      <div class="index-tree-content">
        <div class="content-header">
          <h2>{{ selectedDetail?.title || item.title || item.label }}</h2>
        </div>
        <el-divider />
        <div class="content-body">
          <template v-if="loadingDetail">
            <el-skeleton :rows="6" animated />
          </template>
          <template v-else-if="selectedDetail?.content">
            <div class="index-links">
              <div v-for="(line, i) in parseContent(selectedDetail.content)" :key="i" class="link-line">
                <template v-if="line.isLink">
                  <a :href="line.url" target="_blank" class="link-item">
                    <el-icon><Link /></el-icon>
                    <span>{{ line.text }}</span>
                  </a>
                </template>
                <template v-else>
                  <span class="text-item">{{ line.text }}</span>
                </template>
              </div>
            </div>
          </template>
          <template v-else>
            <p class="empty-desc">暂无内容</p>
          </template>
        </div>
      </div>
    </template>

    <!-- 详情内容（抽屉/对话框） -->
    <template #detail="{ item, loading }">
      <div class="index-detail-content">
        <template v-if="loading">
          <el-skeleton :rows="10" animated />
        </template>
        <template v-else-if="item">
          <div class="index-links">
            <div v-for="(line, i) in parseContent(item.content || '')" :key="i" class="link-line">
              <template v-if="line.isLink">
                <a :href="line.url" target="_blank" class="link-item">
                  <el-icon><Link /></el-icon>
                  <span>{{ line.text }}</span>
                </a>
              </template>
              <template v-else>
                <span class="text-item">{{ line.text }}</span>
              </template>
            </div>
          </div>
        </template>
      </div>
    </template>
  </NoteViewContainer>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useNotesStore } from '@/stores/notes'
import NoteViewContainer from '@/components/dashboard/widgets/NoteViewContainer.vue'
import { Collection, Link } from '@element-plus/icons-vue'
import type { IndexNote } from '@/types/note'

const notesStore = useNotesStore()
const selectedDetail = ref<IndexNote | null>(null)
const loadingDetail = ref(false)

interface ContentLine {
  text: string
  isLink: boolean
  url?: string
}

// 树形数据：按首字母或标题分组
const treeData = computed(() => {
  const indices = notesStore.indexNotes
  const groupMap = new Map<string, IndexNote[]>()
  
  indices.forEach(idx => {
    const firstChar = idx.title?.charAt(0)?.toUpperCase() || '#'
    const group = /[A-Z]/.test(firstChar) ? firstChar : '#'
    if (!groupMap.has(group)) {
      groupMap.set(group, [])
    }
    groupMap.get(group)!.push(idx)
  })
  
  return Array.from(groupMap.entries())
    .sort(([a], [b]) => a.localeCompare(b))
    .map(([group, items]) => ({
      id: group,
      label: group === '#' ? '其他' : group,
      isCategory: true,
      children: items.map(idx => ({
        id: idx.path || idx.title,
        label: idx.title,
        isLeaf: true,
        ...idx
      }))
    }))
})

// 视图配置
const viewConfig = computed(() => ({
  titleField: 'title',
  descField: 'content',
  tagField: 'category',
  timeField: 'updatedAt',
  cardType: 'book' as const,
  gridColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
  columns: [
    { prop: 'title', label: '标题', minWidth: 200, sortable: true },
    { prop: 'path', label: '路径', width: 200 },
    { prop: 'updatedAt', label: '更新时间', width: 160, type: 'date' as const },
  ],
  treeData: treeData.value,
  chartType: 'bar' as const,
  chartGroupBy: 'category',
  kanbanColumns: getKanbanColumns(),
}))

// 动态生成看板列（按首字母分组）
function getKanbanColumns() {
  const groups = [...new Set(notesStore.indexNotes.map(idx => {
    const firstChar = idx.title?.charAt(0)?.toUpperCase() || '#'
    return /[A-Z]/.test(firstChar) ? firstChar : '#'
  }))].sort()
  
  return groups.map(g => ({
    id: g,
    title: g === '#' ? '其他' : g,
    statusValue: g,
  }))
}

function getPreview(content?: string): string {
  if (!content) return '暂无内容'
  const lines = content.split('\n').filter(l => l.trim())
  return lines.length > 0 ? `${lines.length} 条记录` : '暂无内容'
}

function parseContent(content: string): ContentLine[] {
  return content.split('\n')
    .filter(l => l.trim())
    .map(line => {
      const urlMatch = line.match(/(https?:\/\/[^\s]+)/)
      if (urlMatch) {
        return {
          text: line.replace(urlMatch[0], '').trim() || urlMatch[0],
          isLink: true,
          url: urlMatch[0],
        }
      }
      return {
        text: line,
        isLink: false,
      }
    })
}

// 加载详情的函数
async function loadIndexDetail(item: any) {
  if (item.path) {
    return await notesStore.getIndexDetail(item.path)
  }
  return item
}

// 处理树状视图点击事件
async function handleItemClick(item: any) {
  if (item.path && item.isLeaf) {
    loadingDetail.value = true
    const detail = await notesStore.getIndexDetail(item.path)
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
.index-card {
  padding: 16px;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #e4e7ed;
  cursor: pointer;
  transition: all 0.3s;
}

.index-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.1);
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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

.card-stats {
  margin: 0;
  font-size: 12px;
  color: #909399;
}

/* ListView 项 */
.index-list-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  width: 100%;
  cursor: pointer;
  color: #409eff;
}

.index-list-item:hover {
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

/* TreeView 右侧内容 */
.index-tree-content {
  padding: 0 16px;
}

.content-header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.content-body {
  line-height: 2;
}

.index-links {
  font-size: 14px;
}

.link-line {
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.link-line:last-child {
  border-bottom: none;
}

.link-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #409eff;
  text-decoration: none;
}

.link-item:hover {
  text-decoration: underline;
}

.text-item {
  color: #606266;
}

.empty-desc {
  color: #909399;
}

/* 详情内容 */
.index-detail-content {
  max-height: 70vh;
  overflow: auto;
  padding: 0 16px;
}
</style>
