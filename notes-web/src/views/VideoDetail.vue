<template>
  <div class="video-detail">
    <div class="back-btn">
      <el-button @click="$router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回列表
      </el-button>
    </div>

    <el-skeleton v-if="loading" :rows="10" animated />

    <template v-else-if="video">
      <!-- 头部信息 -->
      <el-card class="info-card">
        <div class="video-header">
          <div class="video-info">
            <h1 class="video-title">{{ video.title }}</h1>
            <div class="video-meta">
              <span class="meta-item">
                <el-icon><VideoPlay /></el-icon>
                {{ video.platform }}
              </span>
              <span class="meta-item">
                <el-icon><Calendar /></el-icon>
                {{ video.created }}
              </span>
              <a :href="video.url" target="_blank" class="meta-link">
                <el-icon><Link /></el-icon>
                查看原视频
              </a>
            </div>
            <div class="video-tags">
              <el-tag v-for="tag in video.tags" :key="tag" size="small" type="info">
                {{ tag }}
              </el-tag>
            </div>
          </div>
          <div class="video-status">
            <StatusTag :status="video.status || ''" />
          </div>
        </div>
      </el-card>

      <!-- B站播放器 -->
      <el-card v-if="video.platform === 'bilibili'" class="player-card">
        <iframe
          :src="`//player.bilibili.com/player.html?bvid=${video.videoId}&page=1&autoplay=0`"
          class="bilibili-player"
          scrolling="no"
          border="0"
          frameborder="no"
          framespacing="0"
          allowfullscreen="true"
        />
      </el-card>

      <!-- 内容区域 -->
      <el-card class="content-card">
        <MdRenderer :content="video.content || ''" />
      </el-card>
    </template>

    <el-empty v-else description="视频不存在" />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useNotesStore } from '@/stores/notes'
import { useAsync } from '@/composables/useAsync'
import MdRenderer from '@/components/MdRenderer.vue'
import StatusTag from '@/components/StatusTag.vue'
import type { VideoNote } from '@/types/note'

const route = useRoute()
const notesStore = useNotesStore()

const { data: video, loading, execute } = useAsync<VideoNote | null>()

onMounted(async () => {
  const id = route.params.id as string
  await execute(() => notesStore.getVideoDetail(id) as Promise<VideoNote | null>)
})
</script>

<style scoped>
.video-detail {
  max-width: 1000px;
  margin: 0 auto;
}

.back-btn {
  margin-bottom: 16px;
}

.info-card {
  margin-bottom: 20px;
}

.video-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.video-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
}

.video-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #909399;
}

.meta-link {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #409eff;
  text-decoration: none;
}

.meta-link:hover {
  text-decoration: underline;
}

.video-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.video-status {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 8px;
}

.player-card {
  margin-bottom: 20px;
}

.bilibili-player {
  width: 100%;
  height: 500px;
}

.content-card {
  margin-bottom: 20px;
}
</style>
