<template>
    <section class="flex md:bg-gray-100 min-h-screen overflow-hidden">
      <AdminSidebar />
      
      <div class="flex-grow text-gray-800">
        <AppHeader/>
  
        <main class="p-6 sm:p-10 space-y-6 bg-gray-100">
          <div class="bg-white shadow rounded-lg overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-200">
              <h2 class="text-lg font-semibold text-gray-800">Order Management</h2>
            </div>
            
            <div class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Order ID</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Customer</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="order in filteredOrders" :key="order.id">
                    <!-- Order ID -->
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="text-sm font-medium text-gray-900">#{{ order.ID }}</div>
                    </td>
                    
                    <!-- Customer -->
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="text-sm text-gray-900">{{ order.user.username }}</div>
                      <div class="text-sm text-gray-500">{{ order.user.email }}</div>
                    </td>
                    
                    <!-- Date -->
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="text-sm text-gray-500">{{ formatDate(order.CreatedAt) }}</div>
                    </td>
                    
                    <!-- Total -->
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="text-sm text-gray-900">${{ order.total_amount }}</div>
                    </td>
                    
                    <!-- Status -->
                    <td class="px-6 py-4 whitespace-nowrap">
                      <select 
                        v-model="order.status"
                        @change="updateOrderStatus(order)"
                        class="border rounded px-2 py-1 text-sm"
                        :class="getStatusClass(order.status)"
                      >
                        <option value="pending">Pending</option>
                        <option value="processing">Processing</option>
                        <option value="completed">Completed</option>
                        <option value="canceled">Canceled</option>
                      </select>
                    </td>
                    
                    <!-- Actions -->
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <router-link 
                        :to="`/admin/orders/${order.ID}`"
                        class="text-indigo-600 hover:text-indigo-900 mr-3"
                      >
                        View
                      </router-link>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <!-- Pagination -->
            <div class="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
              <div class="text-sm text-gray-500">
                Showing <span class="font-medium">{{ currentPage * itemsPerPage - itemsPerPage + 1 }}</span> to 
                <span class="font-medium">{{ Math.min(currentPage * itemsPerPage, totalOrders) }}</span> of 
                <span class="font-medium">{{ totalOrders }}</span> orders
              </div>
              <div class="flex space-x-2">
                <button @click="prevPage" :disabled="currentPage === 1" 
                        class="px-3 py-1 border rounded-md" 
                        :class="{'opacity-50 cursor-not-allowed': currentPage === 1}">
                  Previous
                </button>
                <button @click="nextPage" :disabled="currentPage * itemsPerPage >= totalOrders" 
                        class="px-3 py-1 border rounded-md" 
                        :class="{'opacity-50 cursor-not-allowed': currentPage * itemsPerPage >= totalOrders}">
                  Next
                </button>
              </div>
            </div>
          </div>
        </main>
      </div>
  
      <!-- Delete Confirmation Modal -->
    </section>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import axios from 'axios'
  import { useRouter } from 'vue-router'
  import AdminSidebar from '@/components/AdminSidebar.vue'
  import AppHeader from '@/components/HeaderAdmin.vue'
  
  const router = useRouter()
  
  // Data
  const orders = ref<any[]>([])
  const searchQuery = ref('')
  const currentPage = ref(1)
  const itemsPerPage = 10

  
  // Computed properties
  const totalOrders = computed(() => orders.value.length)
  const filteredOrders = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage
    const end = start + itemsPerPage
    return orders.value
      .filter(order => 
        order.ID.toString().includes(searchQuery.value) ||
        order.customer_name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        order.customer_email.toLowerCase().includes(searchQuery.value.toLowerCase()))
      .slice(start, end)
  })
  
  // Methods
  const fetchOrders = async () => {
    try {
      const token = localStorage.getItem('adminToken')
      if (!token) {
        router.push('/admin-login')
        return
      }
  
      const response = await axios.get(`${import.meta.env.VITE_API_URL}/api/orders/all`, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      })
  
      orders.value = response.data
    } catch (error) {
      console.error('Error fetching orders:', error)
    }
  }
  
  const searchOrders = () => {
    currentPage.value = 1
  }
  
  const updateOrderStatus = async (order: any) => {
    try {
      const token = localStorage.getItem('adminToken')
      if (!token) {
        router.push('/admin-login')
        return
      }
  
      const response = await axios.put(`${import.meta.env.VITE_API_URL}/api/orders/${order.ID}/status`, 
        { status: order.status },
        {
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${token}`
          }
        }
      )
        console.log(order.status)
      // Show success message or update UI as needed
      const updatedOrder = response.data
    const index = orders.value.findIndex(o => o.ID === updatedOrder.ID)
    if (index !== -1) {
      orders.value[index] = updatedOrder
    }
    } catch (error) {
      console.error('Error updating order status:', error)
      // Revert status if update fails
      const originalOrder = orders.value.find(o => o.id === order.id)
      if (originalOrder) {
        order.status = originalOrder.status
      }
      alert('Failed to update order status. Please try again.')
    }
  }


  
  const nextPage = () => {
    if (currentPage.value * itemsPerPage < totalOrders.value) {
      currentPage.value++
    }
  }
  
  const prevPage = () => {
    if (currentPage.value > 1) {
      currentPage.value--
    }
  }
  
  const formatDate = (dateString: string) => {
    const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'short', day: 'numeric' }
    return new Date(dateString).toLocaleDateString(undefined, options)
  }
  
  const getStatusClass = (status: string) => {
    switch (status) {
      case 'pending':
        return 'bg-yellow-100 text-yellow-800'
      case 'processing':
        return 'bg-blue-100 text-blue-800'
      case 'completed':
        return 'bg-green-100 text-green-800'
      case 'canceled':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }
  
  // Lifecycle
  onMounted(() => {
    fetchOrders()
  })
  </script>