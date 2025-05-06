<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter  } from 'vue-router'
import axios from 'axios'
import { useI18n } from 'vue-i18n'
import { useCartStore } from '@/stores/cart'

const router = useRouter()
const cartStore = useCartStore()
const showAddedModal = ref(false)
const addToCart = () => {
  const token = localStorage.getItem('authToken')
  if (!token) {
    // Chuyển hướng đến trang đăng nhập nếu chưa đăng nhập
    router.push('/login')
    return
  }
  if (!book.value) return 
  cartStore.addToCart({
    
    id: book.value.id,
    title: book.value.title,
    author: book.value.author,
    price: book.value.price,
    image: book.value.image,
    category: book.value.category
  })
  showAddedModal.value = true
}

const closeModal = () => {
  showAddedModal.value = false
}
const { t } = useI18n()
const route = useRoute()
const apiUrl = import.meta.env.VITE_API_URL

interface Book {
  id: number
  title: string
  author: string
  price: number
  image: string
  description: string
  published_date: string
  pages: number
  language: string
  isbn: string
  reviews_count?: number
  rating?: number
  category: string
  publisher?: string
  stock:number
}

interface Review {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  rating: number
  comment: string
  user: {
    username: string
    first_name: string
    last_name: string
  }
}
const reviewError = ref<string | null>(null)
// Add these new refs
const reviews = ref<Review[]>([])
const newReview = ref({
  rating: 5,
  comment: ''
})

// Add this function to fetch reviews
const fetchReviews = async (bookId: string | number) => {
  try {
    const response = await axios.get(`${apiUrl}/api/books/${bookId}/reviews`)
    reviews.value = response.data?.reviews?.map((r: any) => ({
      ID: r.ID,
      CreatedAt: r.CreatedAt,
      UpdatedAt: r.UpdatedAt,
      rating: r.rating,
      comment: r.comment,
      user: {
        username: r.user?.username || 'Anonymous',
        first_name: r.user?.first_name || '',
        last_name: r.user?.last_name || ''
      }
    })) || []
  } catch (error) {
    console.error('Error fetching reviews:', error)
  }
}

const isSubmittingReview = ref(false)
const showReviewForm = ref(false)

// Add this function to submit a new review
const submitReview = async () => {
  const token = localStorage.getItem('authToken')
  if (!token) {
    router.push('/login')
    return
  }

  if (!book.value) return

  reviewError.value = null
  isSubmittingReview.value = true
  
  try {
    await axios.post(
      `${apiUrl}/api/books/${book.value.id}/reviews`,
      {
        book_id: book.value.id,
        rating: newReview.value.rating,
        comment: newReview.value.comment
      },
      {
        headers: {
          Authorization: `Bearer ${token}`
        }
      }
    )

    // Refresh reviews after successful submission
    await fetchReviews(book.value.id)
    newReview.value = { rating: 5, comment: '' }
    showReviewForm.value = false
  } catch (error: any) {
    if (axios.isAxiosError(error)) {
      if (error.response?.data?.error === "You have already reviewed this book") {
        reviewError.value = t('book.alreadyReviewed')
      } else {
        reviewError.value = error.response?.data?.message || t('book.reviewError')
      }
    } else {
      reviewError.value = t('book.reviewError')
    }
  } finally {
    isSubmittingReview.value = false
  }
}
// Cập nhật hàm mở form review
const openReviewForm = () => {
  const token = localStorage.getItem('authToken')
  if (!token) {
    router.push('/login')
    return
  }
  showReviewForm.value = true
}

const book = ref<Book | null>(null)
const relatedBooks = ref<Book[]>([])
const isLoading = ref(false)

const shuffleArray = (array: any[]) => {
  for (let i = array.length - 1; i > 0; i--) {
    let j = Math.floor(Math.random() * (i + 1)); // Sử dụng 'let' cho biến j
    [array[i], array[j]] = [array[j], array[i]]; // Hoán đổi các phần tử
  }
}


const fetchBook = async (id: string | number) => {
  isLoading.value = true
  try {
    // Fetch book details
    const response = await axios.get(`${apiUrl}/api/books/${id}`)
    const rawBook = response.data?.data

    book.value = {
      id: rawBook.ID,
      title: rawBook.title,
      author: rawBook.author,
      price: rawBook.price,
      image: `${apiUrl}${rawBook.cover_image}`,
      description: rawBook.description,
      published_date: rawBook.published_date,
      pages: rawBook.pages,
      language: rawBook.language,
      category: rawBook.category,
      isbn: rawBook.isbn,
      rating: rawBook.average_rating,
      reviews_count: response.data.reviews_count,
      stock:rawBook.stock
    }

    // Fetch related books
    const relatedResponse = await axios.get(`${apiUrl}/api/books/search`, {
      params: {
        author: rawBook.author,
        category: rawBook.category
      }
    })

    relatedBooks.value = relatedResponse.data?.map((b: any) => ({
      id: b.ID,
      title: b.title,
      author: b.author,
      price: b.price,
      image: `${apiUrl}${b.cover_image}`,
      description: '',
      published_date: '',
      pages: 0,
      language: '',
      isbn: '',
      category: '',
    })) || []

    // Shuffle the relatedBooks array and pick the first 4
    shuffleArray(relatedBooks.value)
    relatedBooks.value = relatedBooks.value.slice(0, 4)
    
    // Scroll to top smoothly
    window.scrollTo({ top: 0, behavior: 'smooth' })

    await fetchReviews(id)

  } catch (error) {
    console.error('❌ Error fetching book details:', error)
  } finally {
    isLoading.value = false
  }
}

// Khi component mount lần đầu
onMounted(() => {
  fetchBook(route.params.id as string)
})

// Khi route param id thay đổi (click sách mới)
watch(() => route.params.id, (newId) => {
  if (newId) {
    fetchBook(newId as string)
  }
})
</script>

<template>
  <div class="w-full bg-gray-50 min-h-screen">
    <!-- Book Details Section -->
    <section v-if="book" class="container mx-auto px-4 py-12">
      <div class="flex flex-col lg:flex-row gap-12">
        <!-- Book Cover -->
        <div class="lg:w-1/3 flex justify-center">
          <div class="w-full max-w-md bg-white p-6 rounded-xl shadow-lg border border-gray-100">
            <img 
              :src="book.image" 
              :alt="book.title" 
              class="w-full h-auto object-cover rounded-lg"
            >
          </div>
        </div>

        <!-- Book Info -->
        <div class="lg:w-2/3">
          <h1 class="text-3xl md:text-4xl font-bold text-gray-800 mb-2">{{ book.title }}</h1>
          <p class="text-xl text-gray-600 mb-6">by {{ book.author }}</p>
          
          <div class="flex items-center mb-6">
            <div class="flex items-center mr-6">
              <svg class="w-6 h-6 text-yellow-400 mr-1" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"></path>
              </svg>
              <span class="text-gray-700 font-medium">{{ book.rating?.toFixed(1) || 'N/A' }} ({{ book.reviews_count || 0 }} reviews)</span>
            </div>
            <span class="text-gray-500">{{ book.pages }} pages</span>
          </div>

          <div class="bg-white p-6 rounded-xl shadow-md border border-gray-100 mb-8">
            <div class="flex justify-between items-start mb-6">
              <div>
                <span class="text-3xl font-bold text-primary">${{ book.price.toFixed(2) }}</span>
                <p class="text-green-600 text-sm mt-1">In Stock</p>
              </div>
              <!-- Nút Add to Cart -->
  <button 
    @click="addToCart"
    class="btn-primary px-8 py-3 text-lg"
    :disabled="book.stock <= 0"
    :class="{'opacity-50 cursor-not-allowed': book.stock <= 0}"
  >
    {{ t('book.addToCart') }}
  </button>

  <!-- Modal Added to Cart -->
  <div v-if="showAddedModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
      <div class="flex items-center mb-4">
        <svg class="w-8 h-8 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
        </svg>
        <h3 class="text-xl text-black font-bold">{{ t('cart.addedTitle') }}</h3>
      </div>
      
      <p class="mb-6">{{ t('cart.addedMessage') }}</p>
      
      <div class="flex justify-end space-x-4">
        <button
          @click="closeModal"
          class="px-4 py-2 text-black border border-gray-300 rounded-lg hover:bg-gray-100"
        >
          {{ t('cart.continueShopping') }}
        </button>
        <router-link
          to="/cart"
          @click="closeModal"
          class="btn-primary px-4 py-2 rounded-lg"
        >
          {{ t('cart.goToCart') }}
        </router-link>
      </div>
    </div>
  </div>
            </div>

            <div class="grid grid-cols-2 md:grid-cols-3 gap-4 text-sm">
            
              <div>
                <p class="text-gray-500">{{ t('book.language') }}</p>
                <p class="font-medium text-black">{{ book.language }}</p>
              </div>
              <div>
                <p class="text-gray-500">{{ t('book.isbn') }}</p>
                <p class="font-medium text-black">{{ book.isbn }}</p>
              </div>
              <div>
                <p class="text-gray-500">{{ t('book.category') }}</p>
                <p class="font-medium text-black">{{ book.category}}</p>
              </div>
            </div>
          </div>

          <!-- Description -->
          <div class="bg-white p-6 rounded-xl shadow-md border border-gray-100 mb-8">
            <h2 class="text-xl font-bold text-gray-800 mb-4">{{ t('book.description') }}</h2>
            <p class="text-gray-700 whitespace-pre-line">{{ book.description }}</p>
          </div>
        </div>
      </div>

      <!-- Related Books -->
      <div v-if="relatedBooks.length > 0" class="mt-16">
        <h2 class="text-2xl font-bold text-secondary mb-8">{{ t('book.relatedBooks') }}</h2>
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-6">
          <div 
            v-for="relatedBook in relatedBooks" 
            :key="relatedBook.id" 
            class="bg-white rounded-lg shadow-md hover:shadow-lg transition-all duration-300 border border-gray-100"
          >
            <router-link :to="`/books/${relatedBook.id}`" class="block">
              <img :src="relatedBook.image" :alt="relatedBook.title" class="w-full rounded-tl-lg rounded-tr-lg h-48 object-cover">
              <div class="p-4">
                <h3 class="font-medium text-gray-800 line-clamp-2 mb-1" style="min-height: 3em;">{{ relatedBook.title }}</h3>
                <p class="text-gray-600 text-sm line-clamp-1 mb-2">{{ relatedBook.author }}</p>
                <span class="font-bold text-primary">${{ relatedBook.price.toFixed(2) }}</span>
              </div>
            </router-link>
          </div>
        </div>
      </div>
    </section>

    <!-- Loading State -->
    <section v-if="isLoading" class="container mx-auto px-4 py-12">
      <div class="flex flex-col lg:flex-row gap-12">
        <div class="lg:w-1/3">
          <div class="bg-gray-200 rounded-xl w-full h-96 animate-pulse"></div>
        </div>
        <div class="lg:w-2/3 space-y-6">
          <div class="h-12 bg-gray-200 rounded w-3/4 animate-pulse"></div>
          <div class="h-6 bg-gray-200 rounded w-1/2 animate-pulse"></div>
          <div class="h-6 bg-gray-200 rounded w-1/4 animate-pulse"></div>
          <div class="h-32 bg-gray-200 rounded animate-pulse"></div>
          <div class="h-48 bg-gray-200 rounded animate-pulse"></div>
        </div>
      </div>
    </section>
        <!-- Reviews Section -->
        <section v-if="book" class="container mx-auto px-4 py-8">
      <div class="bg-white p-6 rounded-xl shadow-md border border-gray-100">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-2xl font-bold text-gray-800">
            {{ t('book.reviews') }} ({{ reviews.length }})
          </h2>
          <button 
            v-if="!showReviewForm" 
            @click="showReviewForm = true"
            class="btn-primary px-4 py-2 text-sm"
          >
            {{ t('book.addReview') }}
          </button>
        </div>

        <!-- Review Form -->
        <div v-if="showReviewForm" class="mb-8 p-4 bg-gray-50 rounded-lg">
                  <!-- Hiển thị thông báo lỗi -->
        <div v-if="reviewError" class="mb-4 p-3 bg-red-50 text-red-600 rounded-lg">
          {{ reviewError }}
        </div>
          <h3 class="text-lg font-medium mb-4">{{ t('book.writeReview') }}</h3>
          <div class="mb-4">
            <label class="block text-gray-700 mb-2">{{ t('book.rating') }}</label>
            <div class="flex items-center">
              <button 
                v-for="star in 5" 
                :key="star" 
                @click="newReview.rating = star"
                class="focus:outline-none"
              >
                <svg 
                  class="w-6 h-6" 
                  :class="star <= newReview.rating ? 'text-yellow-400' : 'text-gray-300'" 
                  fill="currentColor" 
                  viewBox="0 0 20 20"
                >
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"></path>
                </svg>
              </button>
            </div>
          </div>
          <div class="mb-4">
            <label class="block text-gray-700 mb-2">{{ t('book.comment') }}</label>
            <textarea
              v-model="newReview.comment"
              class="text-black w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              rows="4"
              :placeholder="t('book.reviewPlaceholder')"
            ></textarea>
          </div>
          <div class="flex justify-end space-x-3">
            <button
              @click="showReviewForm = false"
              class="px-4 py-2 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-100"
            >
              {{ t('common.cancel') }}
            </button>
            <button
              @click="submitReview"
              :disabled="isSubmittingReview || !newReview.comment"
              class="btn-primary px-4 py-2 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ isSubmittingReview ? t('common.submitting') : t('common.submit') }}
            </button>
          </div>
        </div>

<!-- Reviews List -->
<div v-if="reviews.length > 0" class="space-y-6">
  <div v-for="review in reviews" :key="review.ID" class="border-b border-gray-200 pb-6 last:border-0">
    <div class="flex items-start mb-3">
      <div class="w-10 h-10 rounded-full bg-gray-300 flex items-center justify-center mr-4">
        <span class="text-gray-600 font-medium">
          {{ review.user.first_name.charAt(0) || review.user.username.charAt(0) }}
        </span>
      </div>
      <div>
        <h4 class="font-medium text-gray-800">
          {{ review.user.first_name }} {{ review.user.last_name }} 
          <span class="text-gray-500 text-sm">(@{{ review.user.username }})</span>
        </h4>
        <div class="flex items-center mt-1">
          <div class="flex mr-2">
            <svg 
              v-for="star in 5" 
              :key="star" 
              class="w-4 h-4" 
              :class="star <= review.rating ? 'text-yellow-400' : 'text-gray-300'" 
              fill="currentColor" 
              viewBox="0 0 20 20"
            >
              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"></path>
            </svg>
          </div>
          <span class="text-gray-500 text-sm">
            {{ new Date(review.CreatedAt).toLocaleDateString() }}
          </span>
        </div>
      </div>
    </div>
    <p class="text-gray-700 whitespace-pre-line pl-14">{{ review.comment }}</p>
  </div>
</div>
        <div v-else class="text-center py-8 text-gray-500">
          {{ t('book.noReviews') }}
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.btn-primary {
  @apply bg-primary text-white font-medium rounded-lg hover:bg-primary-dark transition-colors;
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