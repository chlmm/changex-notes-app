<template>
  <span :class="['field-renderer', `field-ui-${ui}`]">
    <!-- 链接 -->
    <a v-if="ui === 'link'" :href="linkUrl" target="_blank" class="field-link">
      {{ displayText }}
      <el-icon><Link /></el-icon>
    </a>

    <!-- 星级评分 -->
    <span v-else-if="ui === 'star'" class="field-star">
      <el-rate
        :model-value="numValue"
        :max="hint?.Max || 5"
        disabled
        show-score
        size="small"
      />
    </span>

    <!-- 标签组 -->
    <span v-else-if="ui === 'tags' && tagValues.length > 0" class="field-tags">
      <el-tag v-for="tag in tagValues" :key="tag" size="small" type="info" class="tag-item">
        {{ tag }}
      </el-tag>
    </span>

    <!-- 徽章（同标签但颜色不同） -->
    <span v-else-if="ui === 'badges' && tagValues.length > 0" class="field-badges">
      <el-tag v-for="b in tagValues" :key="b" size="small" effect="dark" class="badge-item">
        {{ b }}
      </el-tag>
    </span>

    <!-- 引用块 -->
    <blockquote v-else-if="ui === 'quote'" class="field-quote">
      {{ displayText }}
    </blockquote>

    <!-- 日期 -->
    <span v-else-if="ui === 'date'" class="field-date">
      {{ formatDate(displayText) }}
    </span>

    <!-- 图片 -->
    <el-image
      v-else-if="ui === 'image'"
      :src="displayText"
      fit="cover"
      style="width:60px;height:80px;border-radius:4px"
      :preview-src-list="[displayText]"
      preview-teleported
    />

    <!-- 下拉选择（显示值） -->
    <el-tag v-else-if="ui === 'select' && displayText" size="small" effect="plain">
      {{ displayText }}
    </el-tag>

    <!-- 默认：纯文本 -->
    <span v-else class="field-text">{{ displayText || '—' }}</span>
  </span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Link } from '@element-plus/icons-vue'
import type { FieldHint } from '@/types/schema'

const props = defineProps<{
  value: unknown
  hint?: FieldHint | null
}>()

const ui = computed(() => props.hint?.UI || 'text')

const displayText = computed(() => {
  if (props.value === null || props.value === undefined) return ''
  if (typeof props.value === 'string') return props.value
  if (typeof props.value === 'number') return String(props.value)
  if (typeof props.value === 'boolean') return props.value ? '是' : '否'
  return String(props.value)
})

const numValue = computed(() => {
  const v = Number(props.value)
  return isNaN(v) ? 0 : v
})

const tagValues = computed(() => {
  if (Array.isArray(props.value)) return props.value.map(String)
  return []
})

const linkUrl = computed(() => {
  const v = displayText.value
  if (v.startsWith('http')) return v
  return `https://${v}`
})

function formatDate(v: string): string {
  if (!v) return '—'
  if (/^\d{4}-\d{2}-\d{2}/.test(v)) return v.slice(0, 10)
  return v
}
</script>

<style scoped>
.field-renderer {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.field-link {
  color: var(--el-color-primary);
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 2px;
}
.field-link:hover { text-decoration: underline; }
.field-tags, .field-badges {
  display: inline-flex;
  flex-wrap: wrap;
  gap: 4px;
}
.field-star :deep(.el-rate) { height: 20px; }
.field-quote {
  margin: 0;
  padding: 4px 12px;
  border-left: 3px solid var(--el-color-primary);
  color: var(--el-text-color-secondary);
  font-style: italic;
  max-width: 360px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.field-text { word-break: break-all; }
.tag-item, .badge-item { margin: 0; }
.badge-item { --el-tag-bg-color: var(--el-color-primary-light-9); }
</style>
