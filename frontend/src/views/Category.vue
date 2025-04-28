<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const apiUrl = import.meta.env.VITE_API_URL

interface Book {
  id: number
  title: string
  author: string
  price: number
  image: string
  rating?: number
  category: string
}

interface CategoryWithBooks {
  name: string
  books: Book[]
}

const categoriesWithBooks = ref<CategoryWithBooks[]>([])
const isLoading = ref(true)

onMounted(async () => {
  try {
    // Lấy danh sách các categories
    const categoriesResponse = await axios.get(`${apiUrl}/api/categories`)
    const categoryNames = categoriesResponse.data.data
    
    // Lấy sách cho mỗi category (giới hạn 4 quyển mỗi loại)
    const categoriesData = await Promise.all(
      categoryNames.map(async (categoryName: string) => {
        try {
          const booksResponse = await axios.get(`${apiUrl}/api/books/by-category`, {
            params: {
              category: categoryName
            }
          })
          
          return {
            name: categoryName,
            books: booksResponse.data.data.slice(0, 4).map((book: any) => ({
              id: book.ID,
              title: book.title,
              author: book.author,
              price: book.price,
              image: `${apiUrl}${book.cover_image}`,
              rating: book.average_rating,
              category: book.category
            }))
          }
        } catch (error) {
          console.error(`Error fetching books for category ${categoryName}:`, error)
          return {
            name: categoryName,
            books: []
          }
        }
      })
    )
    
    categoriesWithBooks.value = categoriesData.filter(category => category.books.length > 0)
  } catch (error) {
    console.error('❌ Error fetching categories:', error)
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div class="w-full bg-gray-50 min-h-screen">
    <!-- Categories Header -->
    <section class="bg-secondary text-white py-12">
      <div class="container mx-auto px-4 text-center">
        <h1 class="text-4xl md:text-5xl font-bold mb-4">{{ t('categories.title') }}</h1>
        <p class="text-lg max-w-2xl mx-auto">{{ t('categories.subtitle') }}</p>
      </div>
    </section>

    <!-- Categories Content -->
    <section class="container mx-auto px-4 py-12">
      <!-- Loading State -->
      <div v-if="isLoading" class="space-y-16">
        <div v-for="n in 3" :key="n" class="animate-pulse">
          <div class="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-6">
            <div v-for="m in 4" :key="m" class="bg-white rounded-lg shadow-md border border-gray-100">
              <div class="aspect-w-3 aspect-h-4 w-full bg-gray-200"></div>
              <div class="p-4 space-y-2">
                <div class="h-5 bg-gray-200 rounded w-3/4"></div>
                <div class="h-4 bg-gray-200 rounded w-1/2"></div>
                <div class="h-5 bg-gray-200 rounded w-1/4"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Categories List -->
      <div v-else class="space-y-16">
        <div v-for="category in categoriesWithBooks" :key="category.name" class="category-section">
            <div class="flex justify-between items-center bg-secondary mb-6 px-4 py-2">
                <h2 class="text-2xl font-bold text-white">{{ category.name }}</h2>
            <router-link 
              :to="`/categories/${encodeURIComponent(category.name)}`" 
              class="text-primary hover:text-primary-dark font-medium flex items-center"
            >
              {{ t('categories.viewAll') }}
              <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
              </svg>
            </router-link>
          </div>

          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-6">
            <div 
              v-for="book in category.books" 
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
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.category-section {
  scroll-margin-top: 1rem;
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