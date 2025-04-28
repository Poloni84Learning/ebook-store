<script setup lang="ts">
import { ref, onMounted, computed, watchEffect } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n()
const apiUrl = import.meta.env.VITE_API_URL

interface Book {
  id: number
  title: string
  author: string
  price: number
  image: string
  rating?: number
  reviews_count?: number
}

const books = ref<Book[]>([])
const isLoading = ref(true)
const searchQuery = ref('')
const currentPage = ref(1)
const selectedSort = ref('newest')

// Responsive items per page
const itemsPerPage = ref(6) // Mặc định 2x3 cho mobile
const updateItemsPerPage = () => {
  const width = window.innerWidth
  if (width >= 1280) itemsPerPage.value = 18 // 6x3
  else if (width >= 1024) itemsPerPage.value = 15 // 5x3
  else if (width >= 768) itemsPerPage.value = 12 // 4x3
  else if (width >= 640) itemsPerPage.value = 9 // 3x3
  else itemsPerPage.value = 6 // 2x3
}

onMounted(async () => {
  try {
    updateItemsPerPage()
    window.addEventListener('resize', updateItemsPerPage)
    
    const response = await axios.get(`${apiUrl}/api/books?limit=50`)
    const rawBooks = response.data?.data ?? []
    
    books.value = rawBooks.map((book: any): Book => ({
      id: book.ID,
      title: book.title,
      author: book.author,
      price: book.price,
      image: `${apiUrl}${book.cover_image}`,
      rating: book.average_rating,
      reviews_count: book.reviews_count
    }))
  } catch (error) {
    console.error('❌ Error fetching books:', error)
  } finally {
    isLoading.value = false
  }
})

const filteredBooks = computed(() => {
  let result = [...books.value]
  
  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(book => 
      book.title.toLowerCase().includes(query) || 
      book.author.toLowerCase().includes(query)
    )
  }
  
  // Sort books
  switch (selectedSort.value) {
    case 'newest':
      result.sort((a, b) => b.id - a.id)
      break
    case 'price-low':
      result.sort((a, b) => a.price - b.price)
      break
    case 'price-high':
      result.sort((a, b) => b.price - a.price)
      break
    case 'rating':
      result.sort((a, b) => (b.rating || 0) - (a.rating || 0))
      break
  }
  
  return result
})

const paginatedBooks = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  return filteredBooks.value.slice(start, start + itemsPerPage.value)
})

const totalPages = computed(() => {
  return Math.ceil(filteredBooks.value.length / itemsPerPage.value)
})
</script>

<template>
  <div class="w-full bg-gray-50 min-h-screen">
    <!-- Page Header -->
    <section class="bg-secondary text-white py-12">
      <div class="container mx-auto px-4 text-center">
        <h1 class="text-4xl md:text-5xl font-bold mb-4">{{ t('books.pageTitle') }}</h1>
        <p class="text-lg max-w-2xl mx-auto">{{ t('books.pageSubtitle') }}</p>
      </div>
    </section>

    <!-- Main Content -->
    <section class="container mx-auto px-4 py-8">
      <!-- Search and Filter Bar -->
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-8 gap-4">
        <div class="w-full md:w-1/3">
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="t('books.searchPlaceholder')"
              class="w-full pl-10 pr-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
            >
            <svg class="absolute left-3 top-2.5 h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>
        </div>
        
        <div class="flex items-center gap-4 w-full md:w-auto">
          <label for="sort" class="text-gray-700 whitespace-nowrap">{{ t('books.sortBy') }}:</label>
          <select
            id="sort"
            v-model="selectedSort"
            class="text-black rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent py-2 px-3"
          >
            <option value="newest">{{ t('books.newest') }}</option>
            <option value="price-low">{{ t('books.priceLow') }}</option>
            <option value="price-high">{{ t('books.priceHigh') }}</option>
            <option value="rating">{{ t('books.topRated') }}</option>
          </select>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
        <div v-for="n in itemsPerPage" :key="n" class="bg-white rounded-lg shadow-md border border-gray-100 animate-pulse">
          <div class="aspect-w-3 aspect-h-4 w-full overflow-hidden bg-gray-200"></div>
          <div class="p-3 space-y-2">
            <div class="h-5 bg-gray-200 rounded w-3/4"></div>
            <div class="h-4 bg-gray-200 rounded w-1/2"></div>
            <div class="h-5 bg-gray-200 rounded w-1/4"></div>
          </div>
        </div>
      </div>

      <!-- Books Grid -->
      <div v-else>
        <div v-if="paginatedBooks.length === 0" class="text-center py-12">
          <h3 class="text-xl font-medium text-gray-700">{{ t('books.noBooksFound') }}</h3>
          <p class="text-gray-500 mt-2">{{ t('books.tryDifferentSearch') }}</p>
        </div>

        <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
          <div 
            v-for="book in paginatedBooks" 
            :key="book.id" 
            class="bg-white rounded-lg shadow-md hover:shadow-lg transition-all duration-300 border border-gray-100 flex flex-col"
          >
            <router-link :to="`/books/${book.id}`" class="flex flex-col h-full">
              <div class="aspect-w-3 aspect-h-4 w-full overflow-hidden">
                <img 
                  :src="book.image" 
                  :alt="book.title" 
                  class="object-cover rounded-tl-lg rounded-tr-lg w-full h-full"
                >
              </div>
              
              <div class="p-3 flex-grow flex flex-col">
                <h3 
                  class="font-medium text-gray-800 hover:text-primary transition-colors line-clamp-2 mb-1"
                  style="min-height: 3em;"
                >
                  {{ book.title }}
                </h3>
                
                <p class="text-gray-600 text-sm line-clamp-1 mb-2">
                  {{ book.author }}
                </p>
                
                <div class="mt-auto flex justify-between items-center">
                  <span class="font-bold text-primary text-base">${{ book.price.toFixed(2) }}</span>
                  <div v-if="book.rating" class="flex items-center">
                    <svg class="w-4 h-4 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292z"></path>
                    </svg>
                    <span class="ml-1 text-gray-600 text-xs">{{ book.rating.toFixed(1) }} ({{ book.reviews_count }})</span>
                  </div>
                </div>
              </div>
            </router-link>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex justify-center mt-12">
          <nav class="flex items-center gap-1">
            <button
              @click="currentPage = Math.max(1, currentPage - 1)"
              :disabled="currentPage === 1"
              class="px-3 py-1 rounded-lg border border-gray-300 text-black hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              &laquo;
            </button>
            
            <button
              v-for="page in totalPages"
              :key="page"
              @click="currentPage = page"
              :class="{'bg-primary text-white': currentPage === page, 'hover:bg-gray-100': currentPage !== page}"
              class="px-4 py-1 text-black rounded-lg border border-gray-300"
            >
              {{ page }}
            </button>
            
            <button
              @click="currentPage = Math.min(totalPages, currentPage + 1)"
              :disabled="currentPage === totalPages"
              class="px-3 py-1 rounded-lg border border-gray-300 text-black hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              &raquo;
            </button>
          </nav>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.aspect-w-3 {
  position: relative;
  padding-bottom: 133.33%; /* 4:3 aspect ratio */
}

.aspect-w-3 > * {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>