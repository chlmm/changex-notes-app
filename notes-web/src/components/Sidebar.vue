<template>
  <div class="sidebar">
    <div class="logo">
      <el-icon size="24"><Notebook /></el-icon>
      <span>笔记库</span>
    </div>

    <el-menu
      :default-active="activeMenu"
      router
      class="sidebar-menu"
    >
      <el-menu-item index="/">
        <el-icon><HomeFilled /></el-icon>
        <span>首页</span>
      </el-menu-item>

      <!-- 按分类动态生成菜单 -->
      <template v-for="cat in categoriesWithTypes" :key="cat.ID">
        <el-sub-menu v-if="cat.Types.length > 1" :index="cat.ID">
          <template #title>
            <el-icon><component :is="cat.Icon" /></el-icon>
            <span>{{ cat.Label }}</span>
          </template>
          <el-menu-item
            v-for="t in cat.Types"
            :key="t.ID"
            :index="'/' + t.ID"
          >
            {{ t.Label }}
          </el-menu-item>
        </el-sub-menu>

        <el-menu-item
          v-else
          :index="'/' + cat.Types[0].ID"
        >
          <el-icon><component :is="cat.Types[0].Icon" /></el-icon>
          <span>{{ cat.Types[0].Label }}</span>
        </el-menu-item>
      </template>
    </el-menu>

    <div class="sidebar-footer">
      <el-tag size="small" type="info">共 {{ totalNotes }} 条笔记</el-tag>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useSchema } from '@/composables/useSchema'
import { getStats } from '@/api/notesV2'
import { Notebook, HomeFilled } from '@element-plus/icons-vue'
import type { CategoryDef, TypeDef } from '@/types/schema'

const route = useRoute()
const { categories, loadSchema, getTypesByCategory } = useSchema()

const totalNotes = ref(0)

interface CategoryWithTypes extends CategoryDef {
  Types: TypeDef[]
}

const categoriesWithTypes = computed<CategoryWithTypes[]>(() => {
  return categories.value.map(cat => ({
    ...cat,
    Types: getTypesByCategory(cat.ID),
  })).filter(c => c.Types.length > 0)
})

const activeMenu = computed(() => {
  const path = route.path
  if (path === '/') return '/'
  const parts = path.split('/')
  return '/' + parts[1]
})

onMounted(async () => {
  await loadSchema()
  getStats().then(s => {
    totalNotes.value = Object.values(s).reduce((a, b) => a + b, 0)
  }).catch(() => {})
})
</script>

<style scoped>
.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 20px;
  font-size: 18px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  border-bottom: 1px solid var(--el-border-color-light);
}

.sidebar-menu {
  flex: 1;
  border-right: none;
}

.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--el-border-color-light);
}
</style>
