<template>
  <div class="generic-note-list">
    <!-- 头部 -->
    <div class="list-header">
      <h2>
        <el-icon v-if="type?.Icon" :size="22"><component :is="type.Icon" /></el-icon>
        {{ type?.Label || typeId }}
      </h2>
      <span class="list-count">{{ notes.length }} 条</span>
    </div>

    <!-- 加载中 -->
    <el-skeleton v-if="loading" :rows="6" animated />

    <!-- 空状态 -->
    <el-empty v-else-if="!notes.length" description="暂无数据" />

    <!-- 数据表格 -->
    <el-table
      v-else
      :data="notes"
      stripe
      style="width: 100%"
      @row-click="onRowClick"
      row-class-name="clickable-row"
    >
      <el-table-column
        v-for="col in columns"
        :key="col.key"
        :label="col.label"
        :min-width="col.minWidth"
      >
        <template #default="{ row }">
          <FieldRenderer
            :value="row.fields[col.key]"
            :hint="col.hint"
          />
        </template>
      </el-table-column>

      <!-- 操作列 -->
      <el-table-column label="" width="60" align="center">
        <template #default="{ row }">
          <el-button text type="primary" size="small" @click.stop="onRowClick(row)">
            详情
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div v-if="notes.length > pageSize" class="list-pagination">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="notes.length"
        layout="total, prev, pager, next"
        background
        small
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSchema } from '@/composables/useSchema'
import { getNoteList } from '@/api/notesV2'
import { getFieldLabel } from '@/composables/useSchema'
import FieldRenderer from '@/components/fields/FieldRenderer.vue'
import type { GenericNote } from '@/types/note'
import type { TypeDef, FieldHint } from '@/types/schema'

const props = defineProps<{
  typeId: string
}>()

const { getType, loadSchema } = useSchema()
const router = useRouter()

const notes = ref<GenericNote[]>([])
const loading = ref(true)
const currentPage = ref(1)
const pageSize = 20

const type = computed<TypeDef | undefined>(() => getType(props.typeId))

interface ColumnDef {
  key: string
  label: string
  hint: FieldHint | null
  minWidth: number
}

const columns = computed<ColumnDef[]>(() => {
  if (!type.value) return []
  const cols = type.value.ListColumns
  const hints = type.value.FieldHints || {}
  return cols.map(key => ({
    key,
    label: getFieldLabel(key, hints[key]),
    hint: hints[key] || null,
    minWidth: key.includes('.') ? 140 : 100,
  }))
})

onMounted(async () => {
  await loadSchema()
  try {
    notes.value = await getNoteList(props.typeId)
  } catch (e) {
    console.error('GenericNoteList: load failed', e)
  } finally {
    loading.value = false
  }
})

function onRowClick(row: GenericNote) {
  router.push({
    name: 'GenericDetail',
    params: { typeId: props.typeId },
    query: { id: row.id },
  })
}
</script>

<style scoped>
.generic-note-list {
  padding: 16px;
}
.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.list-header h2 {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}
.list-count { color: var(--el-text-color-secondary); font-size: 14px; }
.clickable-row { cursor: pointer; }
.list-pagination { margin-top: 16px; display: flex; justify-content: center; }
</style>
