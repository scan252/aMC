<script setup lang="ts">
import { onMounted, ref } from 'vue'
import PageHeader from '../components/PageHeader.vue'
import * as bindings from '../../bindings/github.com/scan252/aMC/windows'

const { NewsService } = bindings as any

interface ForumPost {
  id: string
  title: string
  summary: string
  time: string
  url: string
}

const posts = ref<ForumPost[]>([])
const loading = ref(true)
const lastError = ref('')

onMounted(async () => {
  try {
    posts.value = ((await NewsService.News(20)) as any[]) ?? []
  } catch (e: any) {
    lastError.value = String(e?.message ?? e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <PageHeader title="资讯日历" subtitle="官方公告聚合 · 演示数据（真实模式待接口联调）" />

    <div v-if="loading" class="loading glass">加载中…</div>
    <div v-else-if="lastError" class="msg glass">{{ lastError }}</div>
    <div v-else class="list">
      <a
        v-for="(p, i) in posts"
        :key="p.id"
        class="post glass"
        :class="`a d${Math.min(i + 1, 6)}`"
        :href="p.url"
        target="_blank"
      >
        <div class="post-main">
          <b>{{ p.title }}</b>
          <p>{{ p.summary }}</p>
        </div>
        <span class="time">{{ p.time }}</span>
      </a>
    </div>
  </div>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  overflow-y: auto;
}
.loading,
.msg {
  padding: 30px;
  text-align: center;
  color: var(--txt-3);
  font-size: 13px;
}
.list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.post {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  padding: 16px 20px;
  text-decoration: none;
  color: var(--txt-1);
  transition: all 0.2s;
}
.post:hover {
  transform: translateY(-1px);
  border-color: rgba(237, 190, 90, 0.25);
}
.post b {
  display: block;
  font-size: 14px;
  margin-bottom: 5px;
}
.post p {
  font-size: 12px;
  color: var(--txt-3);
  line-height: 1.6;
}
.time {
  font-family: var(--font-mono);
  font-size: 11.5px;
  color: var(--txt-3);
  flex-shrink: 0;
}
</style>
