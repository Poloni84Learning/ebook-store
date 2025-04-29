import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import BookInfo from '@/views/BookInfo.vue' 
import Checkout from '@/views/Checkout.vue'
import OrderHistory from '@/views/OrderHistory.vue'
import AdminLogin from '@/views/admin/AdminLogin.vue'
import Dashboard from '@/views/admin/Dashboard.vue'
import AddBook from '@/views/admin/AddBook.vue'

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/login',
    name: 'Login',
    component: Login,
    meta: { hideLayout: true, hideHelper: true}
  },
  { path: '/register',
    name: 'Register',
    component: Register,
    meta: { hideLayout: true, hideHelper: true }
  },
  { path: '/checkout',
    name: 'Checkout',
    component: Checkout,
    meta: { hideLayout: true, hideHelper: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/ProfileView.vue'),
    meta: { hideHelper: true }
    
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
    component: () => import('@/views/Cart.vue').then(m => m.default || m),
    meta: { hideHelper: true }
  },
  {
    path:"/order-history",
    name:"Order History",
    component: OrderHistory,
    meta: { hideHelper: true }
  },
  {
    path:"/admin-login",
    name:"Admin Login",
    component: AdminLogin,
    meta: { hideLayout: true, hideHelper: true }
  },
  {
    path:"/admin/dashboard",
    name:"Admin Dashboard",
    component: Dashboard,
    meta: { hideLayout: true, hideHelper: true }
  },
  {
    path: '/admin/books',
    name: 'Admin Books',
    component: () => import('@/views/admin/BookAdmin.vue'),
    meta: { hideLayout: true, hideHelper: true }
  },
  {
    path:"/admin/books/add",
    name:"Admin Add Book",
    component: AddBook,
    meta: { hideLayout: true, hideHelper: true }
  },
  {
    path: '/admin/books/view',
    name: 'Book Management',
    component: () => import('@/views/admin/ViewBook.vue'),
    meta: { hideLayout: true, hideHelper: true }
  },
  {
    path: '/admin/orders',
    name: 'OrderManagement',
    component: () => import('@/views/admin/OrderManagement.vue'),
    meta: { hideLayout: true, hideHelper: true }
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
