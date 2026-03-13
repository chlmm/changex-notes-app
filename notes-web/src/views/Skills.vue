<template>
  <NoteViewContainer
    title="技能库"
    :items="notesStore.skills"
    :loading="notesStore.loading"
    empty-text="暂无技能"
    :available-views="['grid', 'list', 'table', 'tree', 'chart', 'kanban']"
    default-view="grid"
    :view-config="viewConfig"
    :load-detail="loadSkillDetail"
    @item-click="handleItemClick"
  >
    <!-- GridView 自定义卡片 -->
    <template #gridItem="{ item }">
      <div class="skill-card">
        <div class="skill-icon">
          <el-icon size="32"><Tools /></el-icon>
        </div>
        <div class="skill-content">
          <h3>{{ item.name }}</h3>
          <p class="skill-category">{{ item.category }}</p>
          <el-tag size="small" :type="item.type === 'skill' ? 'primary' : 'success'">
            {{ item.type }}
          </el-tag>
        </div>
      </div>
    </template>

    <!-- ListView 自定义项 -->
    <template #listItem="{ item }">
      <div class="skill-list-item">
        <div class="skill-icon-small">
          <el-icon size="24"><Tools /></el-icon>
        </div>
        <div class="skill-info">
          <h3>{{ item.name }}</h3>
          <div class="skill-meta">
            <el-tag size="small">{{ item.category }}</el-tag>
            <el-tag size="small" :type="item.type === 'skill' ? 'primary' : 'success'">
              {{ item.type }}
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
      <div class="skill-tree-content">
        <div class="content-header">
          <h2>{{ selectedDetail?.name || item.name || item.label }}</h2>
          <div class="content-tags" v-if="selectedDetail">
            <el-tag>{{ selectedDetail.category }}</el-tag>
            <el-tag :type="selectedDetail.type === 'skill' ? 'primary' : 'success'">
              {{ selectedDetail.type }}
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
            <p class="empty-desc">暂无内容</p>
          </template>
        </div>
      </div>
    </template>

    <!-- 详情内容（抽屉/对话框） -->
    <template #detail="{ item, loading }">
      <div class="skill-detail-content">
        <template v-if="loading">
          <el-skeleton :rows="10" animated />
        </template>
        <template v-else-if="item">
          <el-descriptions :column="1" border>
            <el-descriptions-item label="名称">{{ item.name }}</el-descriptions-item>
            <el-descriptions-item label="分类">{{ item.category }}</el-descriptions-item>
            <el-descriptions-item label="类型">{{ item.type }}</el-descriptions-item>
          </el-descriptions>
          <el-divider content-position="left">技能内容</el-divider>
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
import { Tools } from '@element-plus/icons-vue'
import type { SkillNote } from '@/types/note'

const notesStore = useNotesStore()
const selectedDetail = ref<SkillNote | null>(null)
const loadingDetail = ref(false)

// 树形数据：按分类分组
const treeData = computed(() => {
  const skills = notesStore.skills
  const categoryMap = new Map<string, SkillNote[]>()
  
  skills.forEach(skill => {
    const cat = skill.category || '未分类'
    if (!categoryMap.has(cat)) {
      categoryMap.set(cat, [])
    }
    categoryMap.get(cat)!.push(skill)
  })
  
  return Array.from(categoryMap.entries()).map(([category, items]) => ({
    id: category,
    label: category,
    isCategory: true,
    children: items.map(skill => ({
      id: skill.path || skill.name,
      label: skill.name,
      isLeaf: true,
      ...skill
    }))
  }))
})

// 视图配置
const viewConfig = computed(() => ({
  titleField: 'name',
  descField: 'description',
  tagField: 'category',
  statusField: 'type',
  timeField: 'updatedAt',
  cardType: 'book' as const,
  gridColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
  columns: [
    { prop: 'name', label: '技能名称', minWidth: 200, sortable: true },
    { prop: 'category', label: '分类', width: 120 },
    { prop: 'type', label: '类型', width: 100 },
    { prop: 'updatedAt', label: '更新时间', width: 160, type: 'date' as const },
  ],
  treeData: treeData.value,
  chartType: 'pie' as const,
  chartGroupBy: 'category',
  kanbanColumns: getKanbanColumns(),
}))

// 动态生成看板列
function getKanbanColumns() {
  const categories = [...new Set(notesStore.skills.map(s => s.category || '未分类'))]
  return categories.map(cat => ({
    id: cat,
    title: cat,
    statusValue: cat,
  }))
}

// 加载详情的函数
async function loadSkillDetail(item: any) {
  if (item.path) {
    return await notesStore.getSkillDetail(item.path)
  }
  return item
}

// 处理树状视图点击事件
async function handleItemClick(item: any) {
  if (item.path && item.isLeaf) {
    loadingDetail.value = true
    const detail = await notesStore.getSkillDetail(item.path)
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
.skill-card {
  padding: 16px;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #e4e7ed;
  cursor: pointer;
  transition: all 0.3s;
}

.skill-card:hover {
  border-color: #409eff;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.1);
}

.skill-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-bottom: 12px;
}

.skill-content h3 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #303133;
}

.skill-category {
  margin: 0 0 8px 0;
  font-size: 12px;
  color: #909399;
}

/* ListView 项 */
.skill-list-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  width: 100%;
  cursor: pointer;
}

.skill-list-item:hover {
  background: #f5f7fa;
}

.skill-icon-small {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.skill-info {
  flex: 1;
  min-width: 0;
}

.skill-info h3 {
  margin: 0 0 8px 0;
  font-size: 14px;
  color: #303133;
}

.skill-meta {
  display: flex;
  gap: 8px;
}

/* TreeView 右侧内容 */
.skill-tree-content {
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

.empty-desc {
  color: #909399;
}

/* 详情内容 */
.skill-detail-content {
  padding: 0 16px;
}
</style>
