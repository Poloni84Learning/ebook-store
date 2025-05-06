<template>
    <section class="flex md:bg-gray-100 min-h-screen overflow-hidden">
      <!-- Admin Sidebar -->
      <AdminSidebar />
      
      <!-- Main Content -->
      <div class="flex-grow text-gray-800">
        <AppHeader/>
        <!-- Header -->
        <header class="flex items-center h-20 px-6 sm:px-10 bg-white">
          <button class="block sm:hidden relative flex-shrink-0 p-2 mr-2 text-gray-600 hover:bg-gray-100 hover:text-gray-800 focus:bg-gray-100 focus:text-gray-800 rounded-full">
            <span class="sr-only">Menu</span>
            <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="h-6 w-6">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7" />
            </svg>
          </button>
          
          <div class="flex-grow">
            <h1 class="text-2xl font-semibold">Book Management</h1>
          </div>
          
          <div class="flex flex-shrink-0 items-center ml-auto">
            <router-link to="/admin/books/view" class="inline-flex px-5 py-3 text-purple-600 hover:text-purple-700 focus:text-purple-700 hover:bg-purple-100 focus:bg-purple-100 border border-purple-600 rounded-md mr-5">
              <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="flex-shrink-0 h-5 w-5 -ml-1 mt-0.5 mr-2">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
              </svg>
              View Books
            </router-link>
            
            <router-link 
              to="/admin/books/add" 
              class="inline-flex px-5 py-3 text-white bg-purple-600 hover:bg-purple-700 focus:bg-purple-700 rounded-md"
            >
              <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="flex-shrink-0 h-6 w-6 text-white -ml-1 mr-2">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Add New Book
            </router-link>
            
          </div>
        </header>
  
        <!-- Main Content Area -->
        <main class="flex-grow p-6 sm:p-10 space-y-6 bg-gray-100">
          <!-- Time Period Selector -->
          <div class="flex justify-end">
            <select 
              v-model="selectedPeriod"
              @change="fetchAllData"
              class="rounded-lg border-2 border-gray-300 p-2 text-gray-700 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-200"
            >
              <option value="yesterday">Yesterday</option>
              <option value="day-before">Day Before</option>
              <option value="last-week">Last Week</option>
              <option value="last-two-weeks">Last Two Weeks</option>
            </select>
          </div>
  
          <!-- Three Column Layout -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- Top Authors Column -->
            <div class="bg-white shadow rounded-lg overflow-hidden">
              <div class="px-6 py-4 border-b border-gray-200">
                <h2 class="text-lg font-semibold text-gray-800">Top Authors</h2>
              </div>
              <div class="divide-y divide-gray-200">
                <div 
                  v-for="(author, index) in topAuthors" 
                  :key="'author-' + index"
                  class="px-6 py-4 flex justify-between items-center hover:bg-gray-50"
                >
                  <div>
                    <h3 class="text-md font-medium text-gray-800">{{ author.author }}</h3>
                    <p class="text-sm text-gray-500">{{ author.completed_orders_count }} orders</p>
                  </div>
                  <span class="text-lg font-semibold text-purple-600">#{{ index + 1 }}</span>
                </div>
                <div v-if="topAuthors.length === 0" class="px-6 py-4 text-center text-gray-500">
                  No data available
                </div>
              </div>
            </div>
  
            <!-- Top Categories Column -->
            <div class="bg-white shadow rounded-lg overflow-hidden">
              <div class="px-6 py-4 border-b border-gray-200">
                <h2 class="text-lg font-semibold text-gray-800">Top Categories</h2>
              </div>
              <div class="divide-y divide-gray-200">
                <div 
                  v-for="(category, index) in topCategories" 
                  :key="'category-' + index"
                  class="px-6 py-4 flex justify-between items-center hover:bg-gray-50"
                >
                  <div>
                    <h3 class="text-md font-medium text-gray-800">{{ category.category }}</h3>
                    <p class="text-sm text-gray-500">{{ category.completed_orders_count }} orders</p>
                  </div>
                  <span class="text-lg font-semibold text-purple-600">#{{ index + 1 }}</span>
                </div>
                <div v-if="topCategories.length === 0" class="px-6 py-4 text-center text-gray-500">
                  No data available
                </div>
              </div>
            </div>
  
            <!-- Top Books Column -->
            <div class="bg-white shadow rounded-lg overflow-hidden">
              <div class="px-6 py-4 border-b border-gray-200">
                <h2 class="text-lg font-semibold text-gray-800">Top Books</h2>
              </div>
              <div class="divide-y divide-gray-200">
                <div 
                  v-for="(book, index) in topBooks" 
                  :key="'book-' + book.id"
                  class="px-6 py-4 flex items-center hover:bg-gray-50"
                >
                  <img 
                    :src="getCorrectImageUrl(book.cover_image)" 
                    :alt="book.title"
                    class="w-12 h-16 object-cover rounded"
                    @error="handleImageError"
                  >
                  <div class="ml-4 flex-grow">
                    <h3 class="text-md font-medium text-gray-800">{{ book.title }}</h3>
                    <p class="text-sm text-gray-500">{{ book.author }}</p>
                    <p class="text-sm font-semibold text-purple-600">${{ book.price }}</p>
                  </div>
                  <div class="ml-4 text-right">
                    <span class="block text-sm text-gray-500">{{ book.completed_orders_count }} orders</span>
                    <span class="block text-lg font-semibold text-purple-600">#{{ index + 1 }}</span>
                  </div>
                </div>
                <div v-if="topBooks.length === 0" class="px-6 py-4 text-center text-gray-500">
                  No data available
                </div>
              </div>
            </div>
          </div>
        </main>
      </div>
    </section>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted, watch } from 'vue'
  import axios from 'axios'
  import { useRouter, useRoute  } from 'vue-router'
  import AdminSidebar from '@/components/AdminSidebar.vue' // Import the sidebar component
import AppHeader from '@/components/HeaderAdmin.vue'

  const router = useRouter()
  const activeTab = ref('view')
  const route = useRoute()
watch(
  () => route.path,
  (newPath) => {
    if (newPath === '/admin/books/view') {
      activeTab.value = 'view'
    }
  },
  { immediate: true }
)
  // Data refs
  const selectedPeriod = ref('last-week')
  const topAuthors = ref<Array<{author: string, completed_orders_count: number}>>([])
  const topCategories = ref<Array<{category: string, completed_orders_count: number}>>([])
  const topBooks = ref<Array<{
    id: number,
    title: string,
    author: string,
    cover_image: string,
    price: number,
    completed_orders_count: number
  }>>([])
  
  // Handle image loading errors
  const handleImageError = (event: Event) => {
    const img = event.target as HTMLImageElement
    img.src = '/images/default-book-cover.jpg'
  }
  const getCorrectImageUrl = (path: string) => {
  if (!path) return '/placeholder.jpg';
  if (path.startsWith('http')) return path; // Nếu đã là URL đầy đủ
  
  const base = import.meta.env.VITE_API_URL || 'http://localhost:8081';
  return `${base}${path.startsWith('/') ? path : `/${path}`}`;
};
  // Fetch all data
  const fetchAllData = async () => {
    try {
      const token = localStorage.getItem('adminToken')
      if (!token) {
        router.push('/admin-login')
        return
      }
  
      const [authorsRes, categoriesRes, booksRes] = await Promise.all([
        axios.get(`${import.meta.env.VITE_API_URL}/api/admin/dashboard/top-authors`, {
          headers: { Authorization: `Bearer ${token}` },
          params: { period: selectedPeriod.value }
        }),
        axios.get(`${import.meta.env.VITE_API_URL}/api/admin/dashboard/top-categories`, {
          headers: { Authorization: `Bearer ${token}` },
          params: { period: selectedPeriod.value }
        }),
        axios.get(`${import.meta.env.VITE_API_URL}/api/admin/dashboard/top-books`, {
          headers: { Authorization: `Bearer ${token}` },
          params: { period: selectedPeriod.value }
        })
      ])
  
      topAuthors.value = authorsRes.data.data || []
      topCategories.value = categoriesRes.data.data || []
      topBooks.value = booksRes.data.data || []
  
    } catch (error) {
      console.error('Error fetching data:', error)
      // Handle error (show toast, etc.)
    }
  }
  
  // Fetch data on component mount
  onMounted(() => {
    
    fetchAllData()
  })
  </script>