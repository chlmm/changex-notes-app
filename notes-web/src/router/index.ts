import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' },
  },
  {
    path: '/:typeId',
    name: 'GenericList',
    component: () => import('@/components/views/GenericNoteList.vue'),
    props: true,
    meta: { title: '列表' },
  },
  {
    path: '/:typeId/detail',
    name: 'GenericDetail',
    component: () => import('@/components/views/GenericNoteDetail.vue'),
    props: (route) => ({
      typeId: route.params.typeId as string,
      noteId: (route.query.id as string) || '',
    }),
    meta: { title: '详情' },
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
