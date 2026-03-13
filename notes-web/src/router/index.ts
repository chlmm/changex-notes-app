import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' },
  },
  {
    path: '/books',
    name: 'Books',
    component: () => import('@/views/Books.vue'),
    meta: { title: '书籍学习' },
  },
  {
    path: '/books/:type/:name',
    name: 'BookDetail',
    component: () => import('@/views/BookDetail.vue'),
    meta: { title: '书籍详情' },
  },
  {
    path: '/videos',
    name: 'Videos',
    component: () => import('@/views/Videos.vue'),
    meta: { title: '视频学习' },
  },
  {
    path: '/videos/:id',
    name: 'VideoDetail',
    component: () => import('@/views/VideoDetail.vue'),
    meta: { title: '视频详情' },
  },
  {
    path: '/knowledge',
    name: 'Knowledge',
    component: () => import('@/views/Knowledge.vue'),
    meta: { title: '知识库' },
  },
  {
    path: '/skills',
    name: 'Skills',
    component: () => import('@/views/Skills.vue'),
    meta: { title: '技能库' },
  },
  {
    path: '/problems',
    name: 'Problems',
    component: () => import('@/views/Problems.vue'),
    meta: { title: '问题与解决' },
  },
  {
    path: '/index',
    name: 'Index',
    component: () => import('@/views/Index.vue'),
    meta: { title: '索引与收藏' },
  },
  {
    path: '/quotes',
    name: 'Quotes',
    component: () => import('@/views/Quotes.vue'),
    meta: { title: '金句收藏' },
  },
  {
    path: '/github',
    name: 'GitHub',
    component: () => import('@/views/GitHub.vue'),
    meta: { title: 'GitHub 项目收藏' },
  },
  {
    path: '/anime',
    name: 'Anime',
    component: () => import('@/views/Anime.vue'),
    meta: { title: '动漫收藏' },
  },
  {
    path: '/movies',
    name: 'Movies',
    component: () => import('@/views/Movies.vue'),
    meta: { title: '电影收藏' },
  },
  {
    path: '/games',
    name: 'Games',
    component: () => import('@/views/Games.vue'),
    meta: { title: '游戏收藏' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  document.title = `${to.meta.title || '笔记库'} - Notes Web`
  next()
})

export default router
