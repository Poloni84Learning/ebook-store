<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const featuredBooks = ref([
  {
    id: 1,
    title: 'Book Title 1',
    author: 'Author Name',
    price: 19.99,
    image: 'https://via.placeholder.com/200x300'
  },
  {
    id: 2,
    title: 'Book Title 2',
    author: 'Author Name',
    price: 24.99,
    image: 'https://via.placeholder.com/200x300'
  },
  {
    id: 3,
    title: 'Book Title 3',
    author: 'Author Name',
    price: 15.99,
    image: 'https://via.placeholder.com/200x300'
  },
  {
    id: 4,
    title: 'Book Title 4',
    author: 'Author Name',
    price: 29.99,
    image: 'https://via.placeholder.com/200x300'
  }
])

const toggleFavorite = (bookId: number) => {
  // Implement favorite functionality
  console.log('Toggled favorite for book:', bookId)
}
</script>

<template>
  <div>
    <!-- Hero Section -->
    <section class="bg-secondary text-white py-16">
      <div class="container mx-auto px-4 flex flex-col md:flex-row items-center">
        <div class="md:w-1/2 mb-8 md:mb-0">
          <h1 class="text-4xl md:text-5xl font-bold mb-4">{{ t('home.heroTitle') }}</h1>
          <p class="text-lg mb-6">{{ t('home.heroSubtitle') }}</p>
          <router-link to="/books" class="btn-primary">{{ t('home.browseBooks') }}</router-link>
        </div>
        <div class="md:w-1/2">
          <img src="https://via.placeholder.com/500x300" :alt="t('home.heroImageAlt')" class="rounded-lg shadow-xl">
        </div>
      </div>
    </section>

    <!-- Featured Books -->
    <section class="py-12 container mx-auto px-4">
      <h2 class="text-3xl font-bold mb-8 text-secondary">{{ t('home.featuredBooks') }}</h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <div 
          v-for="book in featuredBooks" 
          :key="book.id" 
          class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition-shadow duration-300"
        >
          <img :src="book.image" :alt="book.title" class="w-full h-64 object-cover">
          <div class="p-4">
            <h3 class="font-bold text-lg mb-2">{{ book.title }}</h3>
            <p class="text-gray-600 mb-2">{{ book.author }}</p>
            <div class="flex justify-between items-center">
              <span class="font-bold text-primary">${{ book.price.toFixed(2) }}</span>
              <button 
                @click="toggleFavorite(book.id)" 
                class="text-gray-400 hover:text-favorite transition-colors duration-200"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Call to Action -->
    <section class="py-12 bg-gray-100">
      <div class="container mx-auto px-4 text-center">
        <h2 class="text-3xl font-bold mb-4 text-secondary">{{ t('home.ctaTitle') }}</h2>
        <p class="text-lg mb-6 max-w-2xl mx-auto">{{ t('home.ctaSubtitle') }}</p>
        <router-link to="/register" class="btn-primary inline-block">{{ t('home.signUpNow') }}</router-link>
      </div>
    </section>
  </div>
</template>

<style scoped>
/* Add any component-specific styles here */
</style>