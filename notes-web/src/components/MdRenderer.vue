<template>
  <div class="markdown-body" ref="containerRef" v-html="renderedContent"></div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'
import katex from 'katex'
import 'katex/dist/katex.min.css'

const props = defineProps<{
  content: string
}>()

const containerRef = ref<HTMLElement | null>(null)

// 创建 markdown-it 实例
const md: MarkdownIt = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str: string, lang: string): string {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return `<pre class="hljs"><code>${hljs.highlight(str, { language: lang, ignoreIllegals: true }).value}</code></pre>`
      } catch (e) {
        console.error(e)
      }
    }
    return `<pre class="hljs"><code>${md.utils.escapeHtml(str)}</code></pre>`
  },
})

// 渲染内容
const renderedContent = computed(() => {
  if (!props.content) return ''
  
  let content = props.content
  
  // 处理 LaTeX 公式 - 使用占位符避免 markdown-it 干扰
  // 块级公式 $$...$$
  content = content.replace(/\$\$([\s\S]+?)\$\$/g, (_match, formula) => {
    return `<div class="katex-block-placeholder" data-formula="${encodeURIComponent(formula)}"></div>`
  })
  // 行内公式 $...$
  content = content.replace(/\$([^$\n]+)\$/g, (_match, formula) => {
    return `<span class="katex-inline-placeholder" data-formula="${encodeURIComponent(formula)}"></span>`
  })
  
  return md.render(content)
})

// 渲染 KaTeX 公式
function renderKatex() {
  if (!containerRef.value) return
  
  // 渲染行内公式
  containerRef.value.querySelectorAll('.katex-inline-placeholder').forEach((node) => {
    const el = node as HTMLElement
    const formula = decodeURIComponent(el.dataset.formula || '')
    if (formula) {
      try {
        katex.render(formula, el, { throwOnError: false })
        el.classList.remove('katex-inline-placeholder')
      } catch (e) {
        el.textContent = formula
        console.error('KaTeX error:', e)
      }
    }
  })
  
  // 渲染块级公式
  containerRef.value.querySelectorAll('.katex-block-placeholder').forEach((node) => {
    const el = node as HTMLElement
    const formula = decodeURIComponent(el.dataset.formula || '')
    if (formula) {
      try {
        katex.render(formula, el, { displayMode: true, throwOnError: false })
        el.classList.remove('katex-block-placeholder')
      } catch (e) {
        el.textContent = formula
        console.error('KaTeX error:', e)
      }
    }
  })
}

// 监听内容变化，重新渲染 KaTeX
watch(() => props.content, () => {
  nextTick(renderKatex)
}, { immediate: true })
</script>

<style scoped>
.markdown-body {
  line-height: 1.8;
  color: #303133;
}

.markdown-body :deep(h1) {
  font-size: 24px;
  margin: 20px 0 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
}

.markdown-body :deep(h2) {
  font-size: 20px;
  margin: 16px 0 12px;
}

.markdown-body :deep(h3) {
  font-size: 16px;
  margin: 12px 0 8px;
}

.markdown-body :deep(p) {
  margin: 8px 0;
}

.markdown-body :deep(ul),
.markdown-body :deep(ol) {
  padding-left: 24px;
}

.markdown-body :deep(li) {
  margin: 4px 0;
}

.markdown-body :deep(code) {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Fira Code', monospace;
  font-size: 0.9em;
}

.markdown-body :deep(pre) {
  margin: 12px 0;
  border-radius: 8px;
  overflow-x: auto;
}

.markdown-body :deep(pre code) {
  background: transparent;
  padding: 0;
}

.markdown-body :deep(blockquote) {
  border-left: 4px solid #409eff;
  padding-left: 16px;
  margin: 12px 0;
  color: #606266;
}

.markdown-body :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 12px 0;
}

.markdown-body :deep(th),
.markdown-body :deep(td) {
  border: 1px solid #dcdfe6;
  padding: 8px 12px;
  text-align: left;
}

.markdown-body :deep(th) {
  background: #f5f7fa;
  font-weight: 600;
}

.markdown-body :deep(.katex) {
  font-size: 1.1em;
}

.markdown-body :deep(.katex-block-placeholder) {
  text-align: center;
  margin: 16px 0;
  overflow-x: auto;
  padding: 8px 0;
}
</style>
