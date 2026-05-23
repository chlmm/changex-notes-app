<template>
  <div :class="['generic-note-detail', { overlay: overlay }]">
    <!-- 加载中 -->
    <el-skeleton v-if="loading" :rows="8" animated />

    <!-- 错误 -->
    <el-result
      v-else-if="error"
      icon="error"
      :title="error"
      :sub-title="`类型: ${typeId}, ID: ${noteId}`"
    >
      <template #extra>
        <el-button v-if="overlay" @click="handleClose">关闭</el-button>
        <el-button v-else @click="$router.back()">返回</el-button>
      </template>
    </el-result>

    <!-- 内容 -->
    <div v-else-if="note" class="detail-content">
      <!-- 返回/关闭按钮 -->
      <div class="detail-back">
        <el-button v-if="overlay" text @click="handleClose">
          <el-icon><Close /></el-icon>
          关闭
        </el-button>
        <el-button v-else text @click="$router.back()">
          <el-icon><ArrowLeft /></el-icon>
          返回列表
        </el-button>
        <div class="detail-actions">
          <el-button size="small" @click="handleEdit">
            <el-icon><Edit /></el-icon>
            编辑
          </el-button>
          <el-popconfirm title="确定要删除这条笔记吗？" @confirm="handleDelete">
            <template #reference>
              <el-button size="small" type="danger">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>

      <!-- 标题 -->
      <h1 class="detail-title">
        {{ type?.Label || typeId }}详情
      </h1>

      <!-- 编辑器 -->
      <NoteEditor ref="editorRef" :type-id="typeId" @saved="onNoteEdited" />

      <!-- 字段信息 -->
      <div class="detail-fields">
        <div
          v-for="field in fields"
          :key="field.key"
          class="detail-field-row"
        >
          <span class="field-label">{{ field.label }}</span>
          <span class="field-value">
            <FieldRenderer
              :value="note.fields[field.key]"
              :hint="field.hint"
            />
          </span>
        </div>
      </div>

      <!-- 子文件 -->
      <div v-if="subFiles.length" class="detail-subfiles">
        <el-divider />
        <h3>关联文件</h3>
        <el-tabs v-model="activeSubFile" type="card">
          <el-tab-pane
            v-for="sf in subFiles"
            :key="sf.key"
            :label="sf.label"
            :name="sf.key"
          >
          <div class="subfile-content">
            <MdRenderer :content="sf.content" />
          </div>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 正文 -->
      <div v-if="note.content" class="detail-body">
        <el-divider />
        <h3>正文</h3>
        <div class="body-content">
          <MdRenderer :content="note.content" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ArrowLeft, Close, Edit, Delete } from '@element-plus/icons-vue'
import { useSchema, getFieldLabel } from '@/composables/useSchema'
import { getNoteDetail, deleteNote } from '@/api/notesV2'
import { useRouter } from 'vue-router'
import FieldRenderer from '@/components/fields/FieldRenderer.vue'
import MdRenderer from '@/components/MdRenderer.vue'
import NoteEditor from '@/components/NoteEditor.vue'
import { ElMessage } from 'element-plus'
import type { GenericNote } from '@/types/note'
import type { TypeDef } from '@/types/schema'

const props = defineProps<{
  typeId: string
  noteId: string
  overlay?: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const router = useRouter()
const handleClose = () => emit('close')

const { getType, loadSchema } = useSchema()
const editorRef = ref<InstanceType<typeof NoteEditor>>()

const loading = ref(true)
const error = ref('')
const note = ref<GenericNote | null>(null)
const activeSubFile = ref('')

const type = computed<TypeDef | undefined>(() => getType(props.typeId))

interface FieldRow {
  key: string
  label: string
  hint: any | null
}

const fields = computed<FieldRow[]>(() => {
  if (!type.value) return []
  const hints = type.value.FieldHints || {}
  // 展示所有有值的 hint 字段
  return Object.entries(hints).map(([key, hint]) => ({
    key,
    label: getFieldLabel(key, hint),
    hint,
  }))
})

const subFiles = computed(() => {
  if (!note.value?.subFiles || !type.value?.SubFiles) return []
  const sf = note.value.subFiles
  return type.value.SubFiles
    .filter(s => sf[s.Name])
    .map(s => ({
      key: s.Name,
      label: s.Label,
      content: sf[s.Name] || '',
    }))
})

async function loadDetail() {
  loading.value = true
  error.value = ''
  note.value = null
  try {
    note.value = await getNoteDetail(props.typeId, props.noteId)
    if (subFiles.value.length) {
      activeSubFile.value = subFiles.value[0].key
    }
  } catch (e) {
    error.value = '笔记加载失败'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadSchema()
  await loadDetail()
})

// 同路由参数变化（如 /book/detail?id=1 → /game/detail?id=2），Vue Router 复用组件不 remount
watch(
  () => [props.typeId, props.noteId],
  async () => { await loadDetail() },
)

function handleEdit() {
  if (note.value) {
    editorRef.value?.open(note.value)
  }
}

async function handleDelete() {
  try {
    await deleteNote(props.typeId, props.noteId)
    ElMessage.success('笔记已删除')
    if (props.overlay) {
      emit('close')
    } else {
      router.push({ name: 'GenericList', params: { typeId: props.typeId } })
    }
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

async function onNoteEdited() {
  ElMessage.success('笔记已更新')
  await loadDetail()
}
</script>

<style scoped>
.generic-note-detail {
  padding: 16px;
  max-width: 900px;
}
.generic-note-detail.overlay {
  padding: 0;
  max-width: none;
}
.detail-back {
  margin-bottom: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.detail-actions {
  display: flex;
  gap: 8px;
}
.detail-title { margin-bottom: 20px; }
.detail-fields {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0;
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
  overflow: hidden;
}
.detail-field-row {
  display: flex;
  border-bottom: 1px solid var(--el-border-color-lighter);
}
.detail-field-row:nth-last-child(-n+2):nth-child(odd),
.detail-field-row:last-child {
  border-bottom: none;
}
.field-label {
  width: 120px;
  min-width: 120px;
  padding: 10px 16px;
  background: var(--el-fill-color-light);
  font-weight: 500;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
.field-value {
  flex: 1;
  padding: 10px 16px;
  display: flex;
  align-items: center;
}
.detail-subfiles { margin-top: 20px; }
.subfile-content {
  max-height: 400px;
  overflow-y: auto;
  padding: 12px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.8;
}
.detail-body { margin-top: 20px; }
.body-content {
  padding: 16px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.8;
}
.body-content :deep(h2) { margin-top: 0; }
.body-content :deep(h3) { font-size: 16px; }
.body-content :deep(li) { margin-left: 20px; }
.body-content :deep(code) {
  background: var(--el-color-primary-light-9);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
}
</style>
