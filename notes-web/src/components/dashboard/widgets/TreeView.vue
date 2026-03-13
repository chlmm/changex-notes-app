<template>
  <div class="tree-view">
    <el-row :gutter="20">
      <!-- 左侧树形导航 -->
      <el-col :span="leftSpan">
        <el-card class="tree-card" :style="{ height: height }">
          <!-- 搜索框 -->
          <el-input
            v-if="searchable"
            v-model="filterText"
            placeholder="搜索..."
            clearable
            size="small"
            class="search-input"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          
          <!-- 树形结构 -->
          <el-tree
            ref="treeRef"
            :data="data"
            :props="treeProps"
            :node-key="nodeKey"
            :default-expand-all="defaultExpandAll"
            :expand-on-click-node="expandOnClickNode"
            :highlight-current="highlightCurrent"
            :filter-node-method="filterNode"
            @node-click="handleNodeClick"
          >
            <template #default="{ node, data }">
              <slot name="node" :node="node" :data="data">
                <span class="tree-node">
                  <el-icon v-if="getNodeIcon(data, node)">
                    <component :is="getNodeIcon(data, node)" />
                  </el-icon>
                  <span class="node-label">{{ node.label }}</span>
                  <el-tag
                    v-if="getNodeTag(data)"
                    size="small"
                    class="node-tag"
                  >
                    {{ getNodeTag(data) }}
                  </el-tag>
                </span>
              </slot>
            </template>
          </el-tree>
        </el-card>
      </el-col>

      <!-- 右侧内容区 -->
      <el-col :span="24 - leftSpan">
        <el-card class="content-card" :style="{ height: height }">
          <template v-if="selectedItem">
            <slot name="content" :item="selectedItem" :data="selectedItem">
              <!-- 默认内容展示 -->
              <div class="default-content">
                <div class="content-header">
                  <h2>{{ getItemTitle(selectedItem) }}</h2>
                  <el-tag v-if="getItemTag(selectedItem)">
                    {{ getItemTag(selectedItem) }}
                  </el-tag>
                </div>
                <el-divider />
                <div class="content-body">
                  {{ getItemDescription(selectedItem) || '暂无内容' }}
                </div>
              </div>
            </slot>
          </template>
          <template v-else>
            <el-empty :description="emptyText" />
          </template>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, shallowRef } from 'vue'
import { Search, Folder, FolderOpened, Document } from '@element-plus/icons-vue'

// 树节点配置
export interface TreeNodeConfig {
  id?: string
  label: string
  children?: TreeNodeConfig[]
  isLeaf?: boolean
  isCategory?: boolean
  [key: string]: any
}

interface Props {
  data: TreeNodeConfig[]
  loading?: boolean
  height?: string
  leftSpan?: number
  nodeKey?: string
  labelField?: string
  childrenField?: string
  defaultExpandAll?: boolean
  expandOnClickNode?: boolean
  highlightCurrent?: boolean
  searchable?: boolean
  emptyText?: string
  // 字段映射
  titleField?: string
  descField?: string
  tagField?: string
  // 图标配置
  categoryIcon?: any
  leafIcon?: any
}

const props = withDefaults(defineProps<Props>(), {
  height: 'calc(100vh - 180px)',
  leftSpan: 6,
  nodeKey: 'id',
  labelField: 'label',
  childrenField: 'children',
  defaultExpandAll: true,
  expandOnClickNode: false,
  highlightCurrent: true,
  searchable: true,
  emptyText: '请从左侧选择',
  titleField: 'title',
  descField: 'description',
  tagField: 'category',
  categoryIcon: Folder,
  leafIcon: Document,
})

const emit = defineEmits<{
  nodeClick: [data: any, node: any]
}>()

const treeRef = shallowRef()
const filterText = ref('')
const selectedItem = ref<any>(null)

// 树形属性配置
const treeProps = {
  children: props.childrenField,
  label: props.labelField,
}

// 过滤节点
function filterNode(value: string, data: any) {
  if (!value) return true
  return data[props.labelField]?.toLowerCase().includes(value.toLowerCase())
}

// 监听搜索框
watch(filterText, (val) => {
  treeRef.value?.filter(val)
})

// 获取节点图标
function getNodeIcon(data: any, node: any) {
  if (data.isCategory) return props.categoryIcon
  if (node.isLeaf || data.isLeaf) return props.leafIcon
  if (node.expanded) return FolderOpened
  return Folder
}

// 获取节点标签
function getNodeTag(data: any) {
  if (props.tagField && data[props.tagField]) {
    return data[props.tagField]
  }
  return null
}

// 获取项目标题
function getItemTitle(item: any) {
  return item[props.titleField] || item[props.labelField] || '未命名'
}

// 获取项目描述
function getItemDescription(item: any) {
  return item[props.descField] || item.content || ''
}

// 获取项目标签
function getItemTag(item: any) {
  return item[props.tagField] || null
}

// 处理节点点击
function handleNodeClick(data: any, node: any) {
  // 只有叶子节点才选中
  if (node.isLeaf || data.isLeaf || (!data[props.childrenField]?.length)) {
    selectedItem.value = data
    emit('nodeClick', data, node)
  }
}

// 暴露方法
defineExpose({
  treeRef,
  selectedItem,
  setSelectedItem: (item: any) => {
    selectedItem.value = item
  },
  clearSelection: () => {
    selectedItem.value = null
  },
})
</script>

<style scoped>
.tree-view {
  width: 100%;
}

.tree-card {
  overflow: auto;
}

.search-input {
  margin-bottom: 12px;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 6px;
  width: 100%;
}

.node-label {
  flex: 1;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.node-tag {
  margin-left: auto;
}

.content-card {
  overflow: auto;
}

.default-content {
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
</style>
