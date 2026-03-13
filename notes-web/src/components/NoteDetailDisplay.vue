<template>
  <!-- 抽屉模式 -->
  <el-drawer
    v-if="openMode === 'drawer'"
    v-model="visible"
    :title="title"
    :size="drawerSize"
    direction="rtl"
    :before-close="handleClose"
  >
    <div class="detail-content">
      <slot :detail="detail" :loading="loading" />
    </div>
  </el-drawer>

  <!-- 对话框模式 -->
  <el-dialog
    v-else-if="openMode === 'dialog'"
    v-model="visible"
    :title="title"
    :width="dialogWidth"
    top="5vh"
    :before-close="handleClose"
  >
    <div class="detail-content">
      <slot :detail="detail" :loading="loading" />
    </div>
  </el-dialog>

  <!-- 页面模式不需要渲染 -->
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useNoteViewSettings, type OpenMode } from '@/composables/useNoteViewSettings'

interface Props {
  // 当前选中项
  item?: any
  // 标题
  title?: string
  // 自定义打开模式（优先级高于全局设置）
  mode?: OpenMode
  // 跳转路由
  routePath?: string
  // 异步加载详情的函数
  loadDetail?: (item: any) => Promise<any>
  // 抽屉尺寸
  drawerSize?: string
  // 对话框宽度
  dialogWidth?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  open: [item: any]
}>()

const router = useRouter()
const { settings } = useNoteViewSettings()

const visible = ref(false)
const detail = ref<any>(null)
const loading = ref(false)

// 实际使用的打开模式
const openMode = computed(() => props.mode || settings.value.openMode)
const drawerSize = computed(() => props.drawerSize || settings.value.drawerSize || '50%')
const dialogWidth = computed(() => props.dialogWidth || settings.value.dialogWidth || '60%')

// 打开详情
async function open(item: any) {
  emit('open', item)

  if (openMode.value === 'page' && props.routePath) {
    // 页面模式：直接跳转
    router.push(props.routePath)
    return
  }

  // 抽屉或对话框模式
  visible.value = true
  loading.value = true
  detail.value = null

  if (props.loadDetail) {
    try {
      detail.value = await props.loadDetail(item)
    } catch (e) {
      console.error('Failed to load detail:', e)
    }
  } else {
    detail.value = item
  }

  loading.value = false
}

// 关闭
function handleClose() {
  visible.value = false
  detail.value = null
  emit('close')
}

// 暴露方法
defineExpose({
  open,
  close: handleClose,
  visible,
  detail,
  loading,
})
</script>

<style scoped>
.detail-content {
  max-height: 80vh;
  overflow-y: auto;
}
</style>
