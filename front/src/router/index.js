import { createRouter, createWebHistory } from 'vue-router';
import TopView from '../views/TopView.vue';
import BingoView from '../views/BingoView.vue';
import AdminLoginView from '../views/AdminLoginView.vue';
import AdminListView from '../views/AdminListView.vue';
import AdminDetailView from '../views/AdminDetailView.vue';
import { useSessionStore } from '../stores/session';

const routes = [
  { path: '/', name: 'top', component: TopView },
  { path: '/list', name: 'list', component: TopView },
  { path: '/bingo', name: 'bingo', component: BingoView },
  { path: '/bingo/:bingoId', name: 'bingo_detail', component: BingoView, props: true },
  { path: '/admin', name: 'admin', component: AdminLoginView },
  { path: '/admin/login', name: 'admin_login', component: AdminLoginView },
  { path: '/admin/list', name: 'admin_list', component: AdminListView },
  { path: '/admin/new', name: 'admin_new', component: AdminDetailView },
  { path: '/admin/edit/:bingoId', name: 'admin_edit', component: AdminDetailView, props: true },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior (to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth',
      }
    } else {
      return { x: 0, y: 0 }
    }
  },
})

// 管理画面ページのトークンチェック
router.beforeEach((to, from, next) => {
  const sessionStore = useSessionStore();
  const isAdminRoute = to.path.startsWith('/admin') && to.path !== '/admin/login';

  if (isAdminRoute && !sessionStore.isAuthed) {
    console.warn(`invalid Authenticated (path:${to.path})`)
    next('/admin/login');
  } else {
    next();
  }
})

export default router
