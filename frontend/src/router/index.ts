import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import BookInfo from '@/views/BookInfo.vue' 
import Checkout from '@/views/Checkout.vue'
import OrderHistory from '@/views/OrderHistory.vue'
import AdminLogin from '@/views/admin/AdminLogin.vue'

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/login',
    name: 'Login',
    component: Login,
    meta: { hideLayout: true }
  },
  { path: '/register',
    name: 'Register',
    component: Register,
    meta: { hideLayout: true }
  },
  { path: '/checkout',
    name: 'Checkout',
    component: Checkout,
    meta: { hideLayout: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/ProfileView.vue')
    
  },
  {
    path: '/books',
    name: 'Books',
    component: () => import('@/views/Book.vue')
  },
  {
    path: '/categories',
    name: 'Categories',
    component: () => import('@/views/Category.vue')
  },
  {
    path: '/books/:id', 
    name: 'Book-info',
    component: BookInfo, 
    props: true 
  },
  {
    path: '/books/search',
    name: 'BookSearch',
    component: () => import('@/views/BookSearch.vue')
  },
  {
    path: '/categories/:categoryName',
    name: 'CategoryInfo',
    component: () => import('@/views/CategoryInfo.vue'),
    props: true 
  },
  {
    path: '/cart',
    name: 'Cart',
    component: () => import('@/views/Cart.vue').then(m => m.default || m)
  },
  {
    path:"/order-history",
    name:"Order History",
    component: OrderHistory
  },
  {
    path:"/admin-login",
    name:"Admin Login",
    component: AdminLogin,
    meta: { hideLayout: true }
  },
  
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
