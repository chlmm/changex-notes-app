<template>
  <el-card class="generic-note-card" shadow="hover" @click="$emit('click')">
    <!-- Title -->
    <div class="card-title">{{ title }}</div>

    <!-- Subtitle / description line -->
    <div v-if="subtitle" class="card-subtitle">
      {{ subtitle }}
    </div>

    <!-- Key fields -->
    <div v-if="cardFields.length" class="card-fields">
      <div v-for="f in cardFields" :key="f.key" class="card-field">
        <span class="field-label">{{ f.label }}:</span>
        <FieldRenderer :value="note.fields[f.key]" :hint="f.hint" />
      </div>
    </div>

    <!-- Tags -->
    <div v-if="tags.length" class="card-tags">
      <el-tag v-for="t in tags" :key="t" size="small" type="info">{{ t }}</el-tag>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { GenericNote } from '@/types/note'
import type { FieldHint } from '@/types/schema'
import FieldRenderer from '@/components/fields/FieldRenderer.vue'

const props = defineProps<{
  note: GenericNote
  cardFields: { key: string; label: string; hint: FieldHint | null }[]
}>()

defineEmits(['click'])

const title = computed(() =>
  String(
    props.note.fields['title'] ||
    props.note.fields['source.title'] ||
    props.note.fields['source.name'] ||
    props.note.fields['name'] ||
    props.note.id
  )
)

const subtitle = computed(() => {
  const desc = props.note.fields['description'] || props.note.fields['summary']
  if (!desc) return ''
  const s = String(desc)
  return s.length > 80 ? s.slice(0, 80) + '...' : s
})

const tags = computed(() => {
  const t = props.note.fields['tags']
  if (Array.isArray(t)) return t.slice(0, 4).map(String)
  return []
})
</script>

<style scoped>
.generic-note-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  height: 100%;
}
.generic-note-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}
.card-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.card-subtitle {
  font-size: 13px;
  color: #909399;
  margin-bottom: 10px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.card-fields {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 8px;
}
.card-field {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}
.field-label {
  color: #909399;
  font-size: 12px;
  min-width: 40px;
  flex-shrink: 0;
}
.card-tags {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  margin-top: 8px;
}
</style>
