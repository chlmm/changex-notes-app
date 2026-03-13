<template>
  <div class="dashboard">
    <h1 class="page-title">{{ title }}</h1>

    <!-- 加载状态 -->
    <template v-if="loading">
      <el-skeleton :rows="5" animated />
    </template>

    <!-- 渲染各行 widgets -->
    <template v-else>
      <template v-for="(row, rowIndex) in rows" :key="rowIndex">
        <!-- GridView 特殊处理：占满整行 -->
        <template v-if="row.widgets.length === 1 && row.widgets[0].type === 'gridView'">
          <GridView
            v-bind="row.widgets[0]"
            @item-click="handleGridItemClick"
          />
        </template>
        
        <!-- ListView 特殊处理：占满整行 -->
        <template v-else-if="row.widgets.length === 1 && row.widgets[0].type === 'listView'">
          <ListView
            v-bind="row.widgets[0]"
            @item-click="handleListItemClick"
          />
        </template>
        
        <!-- TableView 特殊处理：占满整行 -->
        <template v-else-if="row.widgets.length === 1 && row.widgets[0].type === 'tableView'">
          <TableView
            v-bind="row.widgets[0]"
            @row-click="handleTableRowClick"
          />
        </template>
        
        <!-- TreeView 特殊处理：占满整行 -->
        <template v-else-if="row.widgets.length === 1 && row.widgets[0].type === 'treeView'">
          <TreeView
            v-bind="row.widgets[0]"
            @node-click="handleTreeNodeClick"
          />
        </template>
        
        <!-- KanbanView 特殊处理：占满整行 -->
        <template v-else-if="row.widgets.length === 1 && row.widgets[0].type === 'kanbanView'">
          <KanbanView
            v-bind="row.widgets[0]"
            @card-click="handleKanbanCardClick"
            @card-move="handleKanbanCardMove"
          />
        </template>
        
        <!-- TimelineView 特殊处理：占满整行 -->
        <template v-else-if="row.widgets.length === 1 && row.widgets[0].type === 'timelineView'">
          <TimelineView
            v-bind="row.widgets[0]"
            @item-click="handleTimelineItemClick"
          />
        </template>
        
        <!-- ChartView 特殊处理：占满整行 -->
        <template v-else-if="row.widgets.length === 1 && row.widgets[0].type === 'chartView'">
          <ChartView
            v-bind="row.widgets[0]"
            @chart-click="handleChartClick"
          />
        </template>
        
        <!-- NoteViewContainer 特殊处理：占满整行 -->
        <template v-else-if="row.widgets.length === 1 && row.widgets[0].type === 'noteView'">
          <NoteViewContainer
            v-bind="row.widgets[0]"
            @item-click="handleNoteItemClick"
            @card-move="handleNoteCardMove"
            @chart-click="handleNoteChartClick"
            @view-change="handleViewChange"
          />
        </template>
        
        <!-- 普通行布局 -->
        <el-row v-else :gutter="row.gutter || 20" class="dashboard-row">
          <template v-for="(widget, widgetIndex) in row.widgets" :key="widgetIndex">
            <!-- Stats widgets 占 4 列 -->
            <el-col v-if="widget.type === 'stats'" :span="4">
              <StatsCard
                :icon="widget.icon"
                :value="widget.value"
                :label="widget.label"
                :gradient="widget.gradient"
                @click="$router.push(widget.route)"
              />
            </el-col>

            <!-- Progress widget 占 12 列 -->
            <el-col v-else-if="widget.type === 'progress'" :span="12">
              <ProgressCard
                :title="widget.title"
                :items="widget.items"
                :show-more="widget.showMore"
                @more="widget.moreRoute && $router.push(widget.moreRoute)"
              />
            </el-col>

            <!-- QuickEntry widget 占 12 列 -->
            <el-col v-else-if="widget.type === 'quickEntry'" :span="12">
              <QuickEntryCard
                :title="widget.title"
                :items="widget.items"
                @click="(item) => item.route && $router.push(item.route)"
              />
            </el-col>

            <!-- QuickLink widget 占 12 列 -->
            <el-col v-else-if="widget.type === 'quickLink'" :span="12">
              <QuickLinkCard
                :icon="widget.icon"
                :title="widget.title"
                :description="widget.description"
                :color="widget.color"
                @click="widget.route && $router.push(widget.route)"
              />
            </el-col>

            <!-- GridView widget -->
            <el-col v-else-if="widget.type === 'gridView'" :span="24">
              <GridView
                v-bind="widget"
                @item-click="handleGridItemClick"
              />
            </el-col>
          </template>
        </el-row>
      </template>
    </template>
  </div>
</template>

<script setup lang="ts">
import type { DashboardConfig } from '@/composables/useDashboard'
import StatsCard from './widgets/StatsCard.vue'
import ProgressCard from './widgets/ProgressCard.vue'
import QuickEntryCard from './widgets/QuickEntryCard.vue'
import QuickLinkCard from './widgets/QuickLinkCard.vue'
import GridView from './widgets/GridView.vue'
import ListView from './widgets/ListView.vue'
import TableView from './widgets/TableView.vue'
import TreeView from './widgets/TreeView.vue'
import KanbanView from './widgets/KanbanView.vue'
import TimelineView from './widgets/TimelineView.vue'
import ChartView from './widgets/ChartView.vue'
import NoteViewContainer from './widgets/NoteViewContainer.vue'

interface Props {
  title: string
  loading?: boolean
  rows: DashboardConfig['rows']
}

withDefaults(defineProps<Props>(), {
  loading: false,
})

const emit = defineEmits<{
  gridItemClick: [item: any]
  listItemClick: [item: any]
  tableRowClick: [row: any]
  treeNodeClick: [data: any, node: any]
  kanbanCardClick: [item: any]
  kanbanCardMove: [item: any, from: any, to: any]
  timelineItemClick: [item: any]
  chartClick: [params: any]
  noteItemClick: [item: any]
  noteCardMove: [item: any, from: any, to: any]
  noteChartClick: [params: any]
  viewChange: [view: string]
}>()

function handleGridItemClick(item: any) {
  emit('gridItemClick', item)
}

function handleListItemClick(item: any) {
  emit('listItemClick', item)
}

function handleTableRowClick(row: any) {
  emit('tableRowClick', row)
}

function handleTreeNodeClick(data: any, node: any) {
  emit('treeNodeClick', data, node)
}

function handleKanbanCardClick(item: any) {
  emit('kanbanCardClick', item)
}

function handleKanbanCardMove(item: any, from: any, to: any) {
  emit('kanbanCardMove', item, from, to)
}

function handleTimelineItemClick(item: any) {
  emit('timelineItemClick', item)
}

function handleChartClick(params: any) {
  emit('chartClick', params)
}

function handleNoteItemClick(item: any) {
  emit('noteItemClick', item)
}

function handleNoteCardMove(item: any, from: any, to: any) {
  emit('noteCardMove', item, from, to)
}

function handleNoteChartClick(params: any) {
  emit('noteChartClick', params)
}

function handleViewChange(view: string) {
  emit('viewChange', view)
}
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
  margin: 0 auto;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 24px;
}

.dashboard-row {
  margin-bottom: 20px;
}
</style>
