import { createRouter, createWebHistory } from 'vue-router'

// 动态导入组件
const Login = () => import('../views/Login.vue')
const Layout = () => import('../views/Layout.vue')
const Dashboard = () => import('../views/Dashboard.vue')
const Factories = () => import('../views/Factories.vue')
const Products = () => import('../views/Products.vue')
const Statistics = () => import('../views/Statistics.vue')

const routes = [
  { path: '/login', name: 'login', component: Login, meta: { public: true, title: '登录' } },
  {
    path: '/',
    component: Layout,
    children: [
      { path: '', name: 'dashboard', component: Dashboard, meta: { title: '首页' } },
      { path: 'factories', name: 'factories', component: Factories, meta: { title: '工厂管理' } },
      { path: 'products', name: 'products', component: Products, meta: { title: '商品管理' } },
      { path: 'statistics', name: 'statistics', component: Statistics, meta: { title: '数据统计' } }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  if (to.meta.public) return next()
  const token = localStorage.getItem('token')
  if (!token) return next({ path: '/login', query: { redirect: to.fullPath } })
  next()
})

export default router


