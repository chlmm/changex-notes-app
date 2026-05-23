<template>
  <el-dialog
    v-model="visible"
    :title="isEdit ? '编辑笔记' : '新建笔记'"
    width="700px"
    top="5vh"
    destroy-on-close
    @closed="resetForm"
  >
    <el-form ref="formRef" :model="form" label-width="100px" label-position="top">
      <!-- 动态字段 -->
      <template v-for="hint in fieldHints" :key="hint.key">
        <el-form-item :label="hint.label">
          <template v-if="hint.ui === 'text' || hint.ui === undefined">
            <el-input v-model="form.fields[hint.key]" />
          </template>
          <template v-else-if="hint.ui === 'textarea'">
            <el-input v-model="form.fields[hint.key]" type="textarea" :rows="3" />
          </template>
          <template v-else-if="hint.ui === 'select' && hint.options">
            <el-select v-model="form.fields[hint.key]" clearable style="width: 100%">
              <el-option
                v-for="opt in hint.options"
                :key="opt"
                :label="opt"
                :value="opt"
              />
            </el-select>
          </template>
          <template v-else-if="hint.ui === 'date'">
            <el-date-picker
              v-model="form.fields[hint.key]"
              type="date"
              style="width: 100%"
            />
          </template>
          <template v-else-if="hint.ui === 'link'">
            <el-input v-model="form.fields[hint.key]" placeholder="https://..." />
          </template>
          <template v-else-if="hint.ui === 'number'">
            <el-input-number v-model="form.fields[hint.key]" style="width: 200px" />
          </template>
          <template v-else-if="hint.ui === 'star'">
            <el-rate v-model="form.fields[hint.key]" />
          </template>
          <template v-else-if="hint.ui === 'tags'">
            <el-input
              v-model="form.fields[hint.key]"
              placeholder="用逗号分隔多个标签"
            />
            <div v-if="form.fields[hint.key]" class="tag-preview">
              <el-tag
                v-for="(tag, idx) in splitTags(form.fields[hint.key])"
                :key="idx"
                size="small"
                type="info"
                class="preview-tag"
              >
                {{ tag }}
              </el-tag>
            </div>
          </template>
          <template v-else>
            <el-input v-model="form.fields[hint.key]" />
          </template>
        </el-form-item>
      </template>

      <!-- Markdown 正文 -->
      <el-form-item label="正文 (Markdown)">
        <el-input
          v-model="form.content"
          type="textarea"
          :rows="12"
          :placeholder="isEdit ? '支持 Markdown 格式...' : '支持 Markdown 格式的正文内容'"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="saving" @click="handleSave">
        {{ isEdit ? '保存' : '创建' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { useSchema } from '@/composables/useSchema'
import { createNote, updateNote } from '@/api/notesV2'
import type { GenericNote } from '@/types/note'
import type { TypeDef } from '@/types/schema'

interface FieldHint {
  key: string
  label: string
  ui: string
  options?: string[]
}

const props = defineProps<{
  typeId: string
}>()

const emit = defineEmits<{
  saved: [note: GenericNote]
}>()

const { getType } = useSchema()

const visible = ref(false)
const saving = ref(false)
const formRef = ref<FormInstance>()
const editNote = ref<GenericNote | null>(null)

const form = reactive({
  fields: {} as Record<string, any>,
  content: '',
})

const isEdit = computed(() => !!editNote.value)
const type = computed<TypeDef | undefined>(() => getType(props.typeId))

const fieldHints = computed<FieldHint[]>(() => {
  if (!type.value) return []
  const hints = type.value.FieldHints || {}
  return Object.entries(hints).map(([key, hint]) => ({
    key,
    label: hint.Label || key,
    ui: hint.UI || 'text',
    options: hint.Options,
  }))
})

// 将 tags 数组转为逗号分隔字符串用于输入框显示
function tagsToString(v: unknown): string {
  if (Array.isArray(v)) return v.map(String).join(',')
  return String(v ?? '')
}

// 将逗号分隔的字符串 split 为标签数组（供预览使用）
function splitTags(v: unknown): string[] {
  if (Array.isArray(v)) return v.map(String)
  if (typeof v === 'string') return v.split(',').map(s => s.trim()).filter(Boolean)
  return []
}

function open(note?: GenericNote) {
  if (note) {
    editNote.value = note
    const fields: Record<string, any> = { ...note.fields }
    // 将数组类型的 tags 转为逗号分隔字符串，方便输入框编辑
    for (const h of fieldHints.value) {
      if (h.ui === 'tags' && Array.isArray(fields[h.key])) {
        fields[h.key] = tagsToString(fields[h.key])
      }
    }
    form.fields = fields
    form.content = note.content || ''
  } else {
    editNote.value = null
    // 初始化空字段
    const empty: Record<string, any> = {}
    for (const h of fieldHints.value) {
      empty[h.key] = h.ui === 'star' ? 0 : h.ui === 'number' ? 0 : ''
    }
    form.fields = empty
    form.content = ''
  }
  visible.value = true
}

function resetForm() {
  editNote.value = null
  form.fields = {}
  form.content = ''
}

async function handleSave() {
  saving.value = true
  try {
    // 清理空值 & 将 tags 字符串转为数组
    const fields: Record<string, unknown> = {}
    for (const h of fieldHints.value) {
      const v = form.fields[h.key]
      if (v === '' || v === null || v === undefined) continue
      if (h.ui === 'tags' && typeof v === 'string') {
        fields[h.key] = v.split(',').map((s: string) => s.trim()).filter(Boolean)
      } else {
        fields[h.key] = v
      }
    }

    if (isEdit.value && editNote.value) {
      await updateNote(props.typeId, editNote.value.id, fields, form.content || undefined)
      ElMessage.success('笔记已保存')
    } else {
      await createNote(props.typeId, fields, form.content || undefined)
      ElMessage.success('笔记已创建')
    }

    visible.value = false
    emit('saved', editNote.value as GenericNote)
  } catch (e) {
    ElMessage.error('保存失败，请重试')
  } finally {
    saving.value = false
  }
}

defineExpose({ open })
</script>

<style scoped>
.el-form-item {
  margin-bottom: 12px;
}

.tag-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
}

.preview-tag {
  margin-right: 0;
}
</style>
