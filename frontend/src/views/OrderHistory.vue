<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'

const { t } = useI18n()


const orders = ref<any[]>([])
const isLoading = ref(true)

const fetchOrders = async () => {
  try {
    const token = localStorage.getItem('authToken')
    if (!token) return
    
    const response = await axios.get(`${import.meta.env.VITE_API_URL}/api/orders`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    
    orders.value = response.data || []
  } catch (error) {
    console.error('Error fetching orders:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchOrders()
})
const getCorrectImageUrl = (path: string) => {
  if (!path) return '/placeholder.jpg';
  if (path.startsWith('http')) return path; // Nếu đã là URL đầy đủ
  
  const base = import.meta.env.VITE_API_URL || 'http://localhost:8081';
  return `${base}${path.startsWith('/') ? path : `/${path}`}`;
};

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const getTotalItems = (orderItems: any[]) => {
  return orderItems.reduce((total, item) => total + item.Quantity, 0)
}
</script>

<template>
  <div class="container mx-auto p-6">
    <h1 class="text-2xl text-black font-semibold mb-4">{{ t('orderHistory.title') }}</h1>
    
    <div v-if="isLoading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-primary"></div>
    </div>
    
    <div v-else>
      <div v-if="orders.length === 0" class="bg-white rounded-lg shadow p-6 text-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
        <p class="mt-4 text-lg text-gray-600">{{ t('orderHistory.emptyMessage') }}</p>
      </div>
      
      <div v-else class="space-y-6">
        <div v-for="order in [...orders].reverse()" :key="order.ID" class="bg-white rounded-lg shadow-md overflow-hidden">
          <div class="p-4 border-b border-gray-200">
            <div class="flex justify-between items-start">
              <div>
                <h2 class="text-black font-bold text-lg">{{ t('orderHistory.orderId') }}: #{{ order.ID }}</h2>
                <p class="text-gray-600">{{ t('orderHistory.orderDate') }}: {{ formatDate(order.CreatedAt) }}</p>
                <p class="text-gray-600">{{ t('orderHistory.itemsCount') }}: {{ getTotalItems(order.order_items) }}</p>
              </div>
              <div class="text-right">
                <span class="px-3 py-1 rounded-full text-sm font-medium"
                  :class="{
                    'bg-green-100 text-green-800': order.status === 'pending',
                    'bg-blue-100 text-blue-800': order.status === 'processing',
                    'bg-yellow-100 text-yellow-800': order.status === 'completed',
                    'bg-red-100 text-red-800': order.status === 'cancelled'
                  }">
                  {{ t(`orderHistory.status.${order.status}`) }}
                </span>
                <p class="mt-2 font-bold text-lg">${{ order.total_amount.toFixed(2) }}</p>
              </div>
            </div>
          </div>
          
          <div class="p-4">
            <div class="mb-4">
              <h3 class="text-black font-semibold mb-2">{{ t('orderHistory.customerInfo') }}</h3>
              <p class="text-gray-600">{{ t('common.name') }}: {{ order.user.username }}</p>
              <p class="text-gray-600">{{ t('common.email') }}: {{ order.user.email }}</p>
              <p class="text-gray-600">{{ t('orderHistory.paymentMethod') }}: {{ order.payment_method }}</p>
            </div>
            
            <div class="mt-4">
              <h3 class="text-black font-semibold mb-2">{{ t('orderHistory.orderItems') }}</h3>
              <div class="space-y-3">
                <div v-for="item in order.order_items" :key="item.ID" class="flex items-center border-b border-gray-100 pb-3">
                  <div class="flex-shrink-0 w-16 h-16 bg-gray-200 rounded-md overflow-hidden">
                    <img 
                      :src="getCorrectImageUrl(item.Book.cover_image)"
                      :alt="item.Book.title"
                      class="w-full h-full object-cover"
                    >
                  </div>
                  <div class="ml-4 flex-1">
                    <h4 class="text-black font-medium">{{ item.Book.title }}</h4>
                    <p class="text-sm text-gray-600">{{ t('common.by') }} {{ item.Book.author }}</p>
                    <p class="text-sm text-gray-600">{{ t('common.quantity') }}: {{ item.Quantity }}</p>
                    <p class="text-sm text-gray-600">${{ item.Price.toFixed(2) }} {{ t('common.each') }}</p>
                  </div>
                  <div class="text-right">
                    <p class="font-medium">${{ (item.Price * item.Quantity).toFixed(2) }}</p>
                  </div>
                </div>
              </div>
            </div>
            
            <div class="mt-6 pt-4 border-t border-gray-200">
              <div class="text-black flex justify-between font-bold text-lg">
                <span>{{ t('orderHistory.total') }}</span>
                <span>${{ order.total_amount.toFixed(2) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>