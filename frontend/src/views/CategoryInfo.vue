<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()
const apiUrl = import.meta.env.VITE_API_URL

interface Book {
  id: number
  title: string
  author: string
  price: number
  image: string
  rating?: number
  description: string
  pages: number
}

const categoryName = ref('')
const books = ref<Book[]>([])
const isLoading = ref(true)

onMounted(async () => {
  try {
    categoryName.value = decodeURIComponent(route.params.categoryName as string)
    
    const response = await axios.get(`${apiUrl}/api/books/by-category`, {
      params: {
        category: categoryName.value
      }
    })
    
    books.value = response.data.data.map((book: any) => ({
      id: book.ID,
      title: book.title,
      author: book.author,
      price: book.price,
      image: `${apiUrl}${book.cover_image}`,
      rating: book.average_rating,
      description: book.description,
      pages: book.pages
    }))
    
  } catch (error) {
    console.error('❌ Error fetching category books:', error)
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div class="w-full bg-gray-50 min-h-screen">
    <!-- Category Header -->
    <section class="bg-secondary text-white py-12">
      <div class="container mx-auto px-4 text-center">
        <h1 class="text-4xl md:text-5xl font-bold mb-2">{{ categoryName }}</h1>
        <p class="text-lg">{{ books.length }} {{ t('category.books') }}</p>
      </div>
    </section>

    <!-- Books List -->
    <section class="container mx-auto px-4 py-12">
      <!-- Loading State -->
      <div v-if="isLoading" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-6">
        <div v-for="n in 8" :key="n" class="bg-white rounded-lg shadow-md border border-gray-100 animate-pulse">
          <div class="aspect-w-3 aspect-h-4 w-full bg-gray-200"></div>
          <div class="p-4 space-y-2">
            <div class="h-5 bg-gray-200 rounded w-3/4"></div>
            <div class="h-4 bg-gray-200 rounded w-1/2"></div>
            <div class="h-5 bg-gray-200 rounded w-1/4"></div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="books.length === 0" class="text-center py-12">
        <h3 class="text-xl font-medium text-gray-700">{{ t('category.noBooks') }}</h3>
      </div>

      <!-- Books Grid -->
      <div v-else class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-6">
        <div 
          v-for="book in books" 
          :key="book.id" 
          class="bg-white rounded-lg shadow-md hover:shadow-lg transition-all duration-300 border border-gray-100"
        >
          <router-link :to="`/books/${book.id}`" class="block">
            <img :src="book.image" :alt="book.title" class="w-full rounded-tl-lg rounded-tr-lg h-48 object-cover">
            <div class="p-4">
              <h3 class="font-medium text-gray-800 line-clamp-2 mb-1" style="min-height: 3em;">{{ book.title }}</h3>
              <p class="text-gray-600 text-sm line-clamp-1 mb-2">{{ book.author }}</p>
              <div class="flex justify-between items-center">
                <span class="font-bold text-primary">${{ book.price.toFixed(2) }}</span>
                <div v-if="book.rating" class="flex items-center">
                  <svg class="w-4 h-4 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"></path>
                  </svg>
                  <span class="ml-1 text-gray-600 text-xs">{{ book.rating.toFixed(1) }}</span>
                </div>
              </div>
            </div>
          </router-link>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
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

.aspect-w-3 {
  position: relative;
  padding-bottom: 133.33%;
}

.aspect-w-3 > * {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}
</style>