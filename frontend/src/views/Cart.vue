<script setup lang="ts">
import { useCartStore } from '@/stores/cart'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { computed,ref, onMounted  } from 'vue'
const router = useRouter()
const cartStore = useCartStore()

// Sử dụng storeToRefs để giữ reactivity cho các getters
const { items, totalItems } = storeToRefs(cartStore)

// Tính toán các giá trị trực tiếp từ store
const subtotal = computed(() => cartStore.subtotal)
const total = computed(() => cartStore.total)
const shippingCost = 3.99

const showDeleteModal = ref(false)
const showClearModal = ref(false)
const showCheckoutModal = ref(false)

const itemToDelete = ref<number | null>(null)

onMounted(() => {
  const token = localStorage.getItem('authToken')
  if (!token) {
    router.push('/login')
  }
})

const confirmClearCart = () => {
  showClearModal.value = true
}

const clearCart = () => {
  cartStore.clearCart()
  showClearModal.value = false
}

const cancelClearCart = () => {
  showClearModal.value = false
}

const confirmDelete = (id: number) => {
  // Đóng tất cả modal khác (nếu có)
  showClearModal.value = false;
  showCheckoutModal.value = false;
  
  // Mở modal xóa
  itemToDelete.value = id;
  showDeleteModal.value = true;
}

const deleteItem = () => {
  if (itemToDelete.value !== null) {
    cartStore.removeFromCart(itemToDelete.value)
  }
  showDeleteModal.value = false
}

const cancelDelete = () => {
  showDeleteModal.value = false
  itemToDelete.value = null
}

const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = 'https://via.placeholder.com/200x300?text=Book+Cover'
}

const increaseQuantity = (id: number) => {
  const item = cartStore.items.find(item => item.id === id)
  if (item) {
    cartStore.updateQuantity(id, item.quantity + 1)
  }
}

const decreaseQuantity = (id: number) => {
  const item = cartStore.items.find(item => item.id === id)
  if (item && item.quantity > 1) {
    cartStore.updateQuantity(id, item.quantity - 1)
  }
}

const confirmCheckout = () => {
  showCheckoutModal.value = true
}

const proceedToCheckout = () => {
  showCheckoutModal.value = false
  router.push('/checkout')
}

const cancelCheckout = () => {
  showCheckoutModal.value = false
}
</script>

<template>
  <div class="flex flex-col h-full bg-white rounded-xl shadow-lg overflow-hidden">
    <!-- Cart Header -->
    <div class="flex items-center justify-between p-6 border-b border-gray-200">
      <h2 class="text-2xl font-bold text-gray-800">{{ $t('cart.title', {count: totalItems}) }}</h2>
      <button
        @click="confirmClearCart"
        type="button"
        class="px-4 py-2 bg-red-100 text-red-600 rounded-lg font-medium hover:bg-red-200 transition-colors"
        :disabled="items.length === 0"
      >
        {{ $t('cart.clearCart') }}
      </button>

      <!-- Modal xác nhận clear cart -->
      <div v-if="showClearModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
          <div class="flex items-start mb-4">
            <div class="flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100">
              <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-medium text-gray-900">{{ $t('cart.clearModal.title') }}</h3>
              <div class="mt-2">
                <p class="text-sm text-gray-500">{{ $t('cart.clearModal.description') }}</p>
              </div>
            </div>
          </div>
          
          <div class="flex justify-end space-x-4 mt-6">
            <button
              @click="cancelClearCart"
              type="button"
              class="px-4 py-2 border text-black border-gray-300 rounded-lg hover: bg-gray-100"
            >
              {{ $t('common.cancel') }}
            </button>
            <button
              @click="clearCart"
              type="button"
              class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
            >
              {{ $t('cart.clearModal.confirm') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  
    <!-- Cart Items -->
    <div class="flex-1 overflow-y-auto p-6">
      <ul v-if="items.length > 0" class="divide-y divide-gray-200">
        <li v-for="item in items" :key="item.id" class="py-6 flex">
          <div class="h-32 w-24 flex-shrink-0 overflow-hidden rounded-lg shadow-sm border border-gray-200">
            <img
              :src="item.image"
              :alt="item.title"
              class="h-full w-full object-cover"
              @error="handleImageError"
            />
          </div>
  
          <div class="ml-4 flex flex-1 flex-col">
            <div class="flex justify-between text-base font-medium text-gray-900">
              <h3>
                <router-link :to="`/books/${item.id}`" class="hover:text-primary">{{ item.title }}</router-link>
              </h3>
              <p class="ml-4">${{ (item.price * item.quantity).toFixed(2) }}</p>
            </div>
            
            <div class="mt-1 text-sm text-gray-600">
              <p><span class="font-semibold">{{ $t('book.category') }}:</span> {{ item.category }}</p>
              <p><span class="font-semibold">{{ $t('book.author') }}:</span> {{ item.author }}</p>
              <p><span class="font-semibold">{{ $t('book.price') }}:</span> ${{ item.price.toFixed(2) }} {{ $t('book.each') }}</p>
            </div>
  
            <div class="flex items-end justify-between mt-4">
              <div class="flex items-center border border-gray-300 rounded-lg">
                <button 
                  @click="decreaseQuantity(item.id)"
                  class="px-3 py-1 text-gray-600 hover: bg-gray-100"
                  :disabled="item.quantity <= 1"
                >−</button>
                <span class="px-3 text-black">{{ item.quantity }}</span>
                <button 
                  @click="increaseQuantity(item.id)"
                  class="px-3 py-1 text-gray-600 hover: bg-gray-100"
                >+</button>
              </div>
              
              <button 
                @click="confirmDelete(item.id)"
                type="button"
                class="font-medium text-red-600 hover:text-red-500 flex items-center px-4 py-2 rounded-md transition-all duration-300 hover:animate-pulse"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-1" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clip-rule="evenodd" />
                </svg>
                {{ $t('common.remove') }}
              </button>

              
            </div>
          </div>
        </li>
      </ul>
      <div v-else class="text-center py-12">
        <svg xmlns="http://www.w3.org/2000/svg" class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        <h3 class="mt-2 text-lg font-medium text-gray-900">{{ $t('cart.empty.title') }}</h3>
        <p class="mt-1 text-gray-500">{{ $t('cart.empty.description') }}</p>
        <div class="mt-6">
          <router-link to="/" class="btn-primary inline-flex items-center">
            {{ $t('cart.continueShopping') }}
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 ml-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
            </svg>
          </router-link>
        </div>
      </div>
      <!-- Modal xác nhận xóa -->
      <div v-if="showDeleteModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
                  <div class="flex items-start mb-4">
                    <div class="flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100">
                      <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                      </svg>
                    </div>
                    <div class="ml-4">
                      <h3 class="text-lg font-medium text-gray-900">{{ $t('cart.deleteModal.title') }}</h3>
                      <div class="mt-2">
                        <p class="text-sm text-gray-500">{{ $t('cart.deleteModal.description') }}</p>
                      </div>
                    </div>
                  </div>
                  
                  <div class="flex justify-end space-x-4 mt-6">
                    <button
                      @click="cancelDelete"
                      type="button"
                      class="px-4 py-2 text-black border border-gray-300 rounded-lg hover: bg-gray-100"
                    >
                      {{ $t('common.cancel') }}
                    </button>
                    <button
                      @click="deleteItem"
                      type="button"
                      class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
                    >
                      {{ $t('common.remove') }}
                    </button>
                  </div>
                </div>
              </div>
      <!-- Empty Cart State -->
      
    </div>
    
    <!-- Cart Summary -->
    <div v-if="items.length > 0" class="border-t border-gray-200 p-6">
      <div class="flex justify-between text-lg font-semibold text-gray-900 mb-2">
        <p>{{ $t('cart.subtotal') }}</p>
        <p>${{ subtotal.toFixed(2) }}</p>
      </div>
      <div class="flex justify-between text-sm text-gray-600 mb-4">
        <p>{{ $t('cart.shipping') }}</p>
        <p>${{ shippingCost.toFixed(2) }}</p>
      </div>
      <div class="flex justify-between text-lg font-bold text-gray-900 mb-6">
        <p>{{ $t('cart.total') }}</p>
        <p>${{ total.toFixed(2) }}</p>
      </div>
  
      <p class="text-sm text-gray-500 mb-6">{{ $t('cart.shippingNote') }}</p>
  
      <button
        @click="confirmCheckout"
        class="w-full bg-primary text-secondary py-3 rounded-lg font-bold hover:bg-secondary hover:text-white transition-colors"
      >
        {{ $t('cart.checkout') }}
      </button>
  
      <div class="mt-4 text-center">
        <router-link to="/" class="text-sm font-medium text-primary hover:text-secondary inline-flex items-center">
          {{ $t('cart.continueShopping') }}
          <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
          </svg>
        </router-link>
      </div>
    </div>
    <!-- Checkout Confirmation Modal -->
    <div v-if="showCheckoutModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4">
        <div class="flex items-start mb-4">
          <div class="flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-blue-100">
            <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <div class="ml-4">
            <h3 class="text-lg font-medium text-gray-900">{{ $t('cart.checkoutModal.title') }}</h3>
            <div class="mt-2">
              <p class="text-sm text-gray-500">{{ $t('cart.checkoutModal.description') }}</p>
              <div class="mt-4 bg-gray-50 p-4 rounded-lg">
                <div class="flex justify-between mb-2">
                  <span class="text-gray-600">{{ $t('cart.subtotal') }}</span>
                  <span class="text-black font-medium">${{ subtotal.toFixed(2) }}</span>
                </div>
                <div class="flex justify-between mb-2">
                  <span class="text-gray-600">{{ $t('cart.shipping') }}</span>
                  <span class="text-black font-medium">${{ shippingCost.toFixed(2) }}</span>
                </div>
                <div class="flex justify-between pt-2 border-t border-gray-200">
                  <span class="text-gray-900 font-semibold">{{ $t('cart.total') }}</span>
                  <span class="text-black text-gray-900 font-bold">${{ total.toFixed(2) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="flex justify-end space-x-4 mt-6">
          <button
            @click="cancelCheckout"
            type="button"
            class="px-4 py-2 text-black border border-gray-300 rounded-lg hover: bg-gray-100"
          >
            {{ $t('common.cancel') }}
          </button>
          <button
            @click="proceedToCheckout"
            type="button"
            class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-secondary"
          >
            {{ $t('cart.checkoutModal.confirm') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.btn-primary {
  @apply bg-primary px-6 py-2 rounded-md text-base font-bold hover:bg-secondary hover:text-white transition-all duration-200 cursor-pointer;
}
</style>