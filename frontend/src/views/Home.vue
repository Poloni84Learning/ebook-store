<script setup lang="ts">
import { ref , onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'
import { useRouter } from 'vue-router';

const router = useRouter();
const { t } = useI18n();

const apiUrl = import.meta.env.VITE_API_URL
interface Book {
  id: number
  title: string
  author: string
  price: number
  image: string
}
const featuredBooks = ref<Book[]>([])
const mostReviewedBooks = ref<Book[]>([])
const mostOrderedBooks = ref<Book[]>([])

function handleSignUp() {
  const token = localStorage.getItem('authToken');
  if (token) {
    alert(t('home.alreadyLoggedIn')); // Hoặc dùng Toast, Modal đẹp hơn
  } else {
    router.push('/register');
  }
}

onMounted(async () => {
  try {
    const responseFeaturedBooks = await axios.get(`${apiUrl}/api/books?limit=4`)

    const rawFeaturedBooks = responseFeaturedBooks.data?.data ?? []

    featuredBooks.value = rawFeaturedBooks.map((featuredBooks: any): Book => ({
      id: featuredBooks.ID,
      title: featuredBooks.title,
      author: featuredBooks.author,
      price: featuredBooks.price,
      image: `${apiUrl}${featuredBooks.cover_image}`
    }))

    const responseMostReviewedBooks = await axios.get(`${apiUrl}/api/books/most-reviewed?limit=4`)

    const rawMostReviewedBooks = responseMostReviewedBooks.data?.data ?? []

    mostReviewedBooks.value = rawMostReviewedBooks.map((mostReviewedBooks: any): Book => ({
      id: mostReviewedBooks.id,
      title: mostReviewedBooks.title,
      author: mostReviewedBooks.author,
      price: mostReviewedBooks.price,
      image: `${apiUrl}${mostReviewedBooks.cover_image}`
    }))

    const responseMostOrderedBooks = await axios.get(`${apiUrl}/api/books/top-selling?limit=4`)

    const rawMostOrderedBooks = responseMostOrderedBooks.data?.data ?? []

    mostOrderedBooks.value = rawMostOrderedBooks.map((mostOrderedBooks: any): Book => ({
      id: mostOrderedBooks.id,
      title: mostOrderedBooks.title,
      author: mostOrderedBooks.author,
      price: mostOrderedBooks.price,
      image: `${apiUrl}${mostOrderedBooks.cover_image}`
    }))




  } catch (error) {
    console.error('❌ Error fetching books:', error)
  }
})
</script>

<template>
  <div class ="w-full">
    <!-- Hero Section -->
    <section class="bg-secondary text-white py-16">
      <div class="container mx-auto px-4 flex flex-col md:flex-row items-center">
        <div class="md:w-1/2 mb-8 md:mb-0">
          <h1 class="text-4xl md:text-5xl font-bold mb-4">{{ t('home.heroTitle') }}</h1>
          <p class="text-lg mb-6">{{ t('home.heroSubtitle') }}</p>
          <router-link to="/books" class="btn-primary">{{ t('home.browseBooks') }}</router-link>
        </div>
        <div class="md:w-1/2">
          <img src="/assets/images/bookshelf.jpg" :alt="t('home.heroImageAlt')" class="rounded-lg shadow-xl">
        </div>
      </div>
    </section>

<!-- Featured Books -->
<section class="bg-white py-12 w-full">
  <div class="container mx-auto px-4">
    <h2 class="text-3xl font-bold mb-8 text-secondary">{{ t('home.featuredBooks') }}</h2>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
      <router-link 
        v-for="book in featuredBooks" 
        :key="book.id" 
        :to="`/books/${book.id}`"
        class="bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-2xl transition-all duration-300 border border-gray-100 block"
      >
        <img :src="book.image" :alt="book.title" class="w-full h-64 object-cover">
        <div class="p-5">
          <h3 class="font-bold text-lg mb-2 text-gray-800">{{ book.title }}</h3>
          <p class="text-gray-600 mb-3">{{ book.author }}</p>
          <div class="flex justify-between items-center">
            <span class="font-bold text-primary text-lg">${{ book.price.toFixed(2) }}</span>
          </div>
        </div>
      </router-link>
    </div>
  </div>
</section>

<!-- Most Reviewed Books -->
<section class="bg-white py-12 w-full">
  <div class="container mx-auto px-4">
    <h2 class="text-3xl font-bold mb-8 text-secondary">{{ t('home.mostReviewedBooks') }}</h2>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
      <router-link 
        v-for="book in mostReviewedBooks" 
        :key="book.id" 
        :to="`/books/${book.id}`"
        class="bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-2xl transition-all duration-300 border border-gray-100 block"
      >
        <img :src="book.image" :alt="book.title" class="w-full h-64 object-cover">
        <div class="p-5">
          <h3 class="font-bold text-lg mb-2 text-gray-800">{{ book.title }}</h3>
          <p class="text-gray-600 mb-3">{{ book.author }}</p>
          <div class="flex justify-between items-center">
            <span class="font-bold text-primary text-lg">${{ book.price.toFixed(2) }}</span>
          </div>
        </div>
      </router-link>
    </div>
  </div>
</section>

<!-- Most Ordered Books -->
<section class="bg-white py-12 w-full">
  <div class="container mx-auto px-4">
    <h2 class="text-3xl font-bold mb-8 text-secondary">{{ t('home.mostOrderedBooks') }}</h2>
    <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
      <router-link 
        v-for="book in mostOrderedBooks" 
        :key="book.id" 
        :to="`/books/${book.id}`"
        class="bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-2xl transition-all duration-300 border border-gray-100 block"
      >
        <img :src="book.image" :alt="book.title" class="w-full h-64 object-cover">
        <div class="p-5">
          <h3 class="font-bold text-lg mb-2 text-gray-800">{{ book.title }}</h3>
          <p class="text-gray-600 mb-3">{{ book.author }}</p>
          <div class="flex justify-between items-center">
            <span class="font-bold text-primary text-lg">${{ book.price.toFixed(2) }}</span>
          </div>
        </div>
      </router-link>
    </div>
  </div>
</section>

    <!-- Call to Action -->
    <section class="py-12 bg-gray-100">
  <div class="container mx-auto px-4 text-center">
    <h2 class="text-3xl font-bold mb-4 text-secondary">{{ t('home.ctaTitle') }}</h2>
    <p class="text-lg mb-6 max-w-2xl mx-auto text-secondary">{{ t('home.ctaSubtitle') }}</p>
    
    <button @click="handleSignUp" class="btn-primary inline-block">
      {{ t('home.signUpNow') }}
    </button>
  </div>
</section>

  </div>
</template>

<style scoped>
/* Add any component-specific styles here */
</style>