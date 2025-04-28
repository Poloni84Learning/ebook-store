<template>
  <section class="flex md:bg-gray-100 min-h-screen overflow-hidden ">
    <!-- Sidebar -->
    <AdminSidebar />

    <!-- Main Content -->
    <div class="flex-grow text-gray-800">
      <!-- Header -->
      <AppHeader/>

      <!-- Dashboard Content -->
      <main class="p-6 sm:p-10 space-y-6">
        <div class="flex flex-col space-y-6 md:space-y-0 md:flex-row justify-between">
          <div class="mr-6">
            <h1 class="text-4xl font-semibold mb-2">Dashboard</h1>
            <h2 class="text-gray-600 ml-0.5">Book Store Inventory</h2>
            
          </div>
          <!-- Select Box cho time_range -->
            
          
          <div class="flex flex-col md:flex-row items-start justify-end -mb-3">
            <router-link to="/admin/books" class="inline-flex px-5 py-3 text-purple-600 hover:text-purple-700 focus:text-purple-700 hover:bg-purple-100 focus:bg-purple-100 border border-purple-600 rounded-md mb-3">
              <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="flex-shrink-0 h-5 w-5 -ml-1 mt-0.5 mr-2">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
              </svg>
              Manage Books
            </router-link>
            <router-link to="/admin/books/add" class="inline-flex px-5 py-3 text-white bg-purple-600 hover:bg-purple-700 focus:bg-purple-700 rounded-md ml-6 mb-3">
              <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="flex-shrink-0 h-6 w-6 text-white -ml-1 mr-2">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Add New Book
            </router-link>
          </div>
        </div>

        <!-- Stats Cards -->
        <section class="grid md:grid-cols-4 xl:grid-cols-4 gap-6">
          <div class="flex items-center p-8 bg-white shadow rounded-lg">
            <div class="inline-flex flex-shrink-0 items-center justify-center h-8 w-8 text-purple-600 bg-purple-100 rounded-full mr-6">
              <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="h-6 w-6">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
              </svg>
            </div>
            <div>
              <span class="block text-xl font-bold">{{ totalBooksSold  }}</span>
              <span class="block text-gray-500">Products</span>
            </div>
          </div>
          
          <div class="flex items-center p-8 bg-white shadow rounded-lg">
            <div class="inline-flex flex-shrink-0 items-center justify-center h-8 w-8 text-green-600 bg-green-100 rounded-full mr-6">
              <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="h-6 w-6">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </div>
            <div>
              <span class="block text-xl font-bold">${{ totalAmount  }}</span>
              <span class="block text-gray-500">Total Sales</span>
            </div>
          </div>
          <div class="flex items-center p-8 bg-white shadow rounded-lg">
            <div class="inline-flex flex-shrink-0 items-center justify-center h-8 w-8 text-blue-600 bg-blue-100 rounded-full mr-6">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <div>
                
              <span class="block text-xl font-bold">{{ totalOrders }}</span>
              <span class="block text-gray-500">Total Orders</span>
            </div>
          </div>
          <div class="flex items-center p-8 bg-white shadow rounded-lg">
            <div class="inline-flex flex-shrink-0 items-center justify-center h-8 w-8 text-red-600 bg-red-100 rounded-full mr-6">
                <svg aria-hidden="true" fill="currentColor" viewBox="0 0 20 20" class="h-6 w-6 text-yellow-400">
  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.286 3.948a1 1 0 00.95.69h4.146c.969 0 1.371 1.24.588 1.81l-3.36 2.44a1 1 0 00-.364 1.118l1.286 3.948c.3.921-.755 1.688-1.54 1.118l-3.36-2.44a1 1 0 00-1.175 0l-3.36 2.44c-.785.57-1.84-.197-1.54-1.118l1.286-3.948a1 1 0 00-.364-1.118l-3.36-2.44c-.783-.57-.38-1.81.588-1.81h4.146a1 1 0 00.95-.69l1.286-3.948z" />
</svg>

            </div>
            <div>
              <span class="inline-block text-xl font-bold"></span>
              <span class="inline-block text-xl text-black font-semibold">{{ trendingBooks }} Books</span>
              <span class="block text-gray-500">Top rated</span>
            </div>
          </div>
          
          
        </section>

        <section class="grid md:grid-cols-1 xl:grid-cols-2 xl:grid-rows-1 xl:grid-flow-col gap-6">
  <!-- Orders Chart -->
<div class="flex flex-col bg-white shadow rounded-lg">
  <!-- Header row - tiêu đề và select box cùng hàng -->
  <div class="px-6 py-5 border-b border-gray-100 flex items-center justify-between">
    <h3 class="text-lg font-semibold text-gray-900">{{ t('Overview') }}</h3>
    <div class="flex items-center space-x-2">
      <select 
        v-model="timeRange" 
        @change="onTimeRangeChange"
        class="rounded-lg border border-gray-300 px-3 py-2 text-sm text-gray-700 focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary">
        <option value="week">{{ t('Week') }}</option>
        <option value="month">{{ t('Month') }}</option>
        <option value="year">{{ t('Year') }}</option>
      </select>
    </div>
  </div>

  <!-- Chart container -->
  <div class="p-4 flex-1 min-h-[300px]">
    <div class="chart-container" style="height: 400px; width: 100%">
    <canvas ref="chartCanvas"></canvas>
  </div>
  </div>
</div>
</section>

        
        <section class="text-right font-semibold text-gray-500">
          <a href="#" class="text-purple-600 hover:underline">EBook Store</a> Dashboard
        </section>
      </main>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick,onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { Chart } from 'chart.js/auto'
import type { ChartConfiguration } from 'chart.js/auto'
import AdminSidebar from '@/components/AdminSidebar.vue' // Import the sidebar component
import DefaultAvatar from '@/assets/img/images_placeholder.png'
import AppHeader from '@/components/HeaderAdmin.vue'

const router = useRouter()
const { t } = useI18n()

// Refs
const totalOrders = ref(0)
const trendingBooks = ref(0)
const totalAmount = ref(0)
const totalBooksSold = ref(0)
const isLoading = ref(true)
const timeRange = ref<'week' | 'month' | 'year'>('week')
const chartCanvas = ref<HTMLCanvasElement | null>(null)
const chartInstance = ref<Chart | null>(null)

// Interfaces
interface OrderTrendData {
  date?: string
  week_start?: string
  month?: string
  count: number
}

// Methods
const fetchRevenueData = async (): Promise<OrderTrendData[]> => {
  try {
    const token = localStorage.getItem('adminToken')
    if (!token) {
      router.push('/admin-login')
      return []
    }

    const response = await axios.get(
      `${import.meta.env.VITE_API_URL}/api/admin/dashboard/order-trend`,
      {
        headers: { Authorization: `Bearer ${token}` },
        params: { time_range: timeRange.value }
      }
    )
    return response.data
  } catch (error) {
    console.error('Error fetching revenue data:', error)
    return []
  }
}

const renderChart = async () => {
  await nextTick() // Đảm bảo DOM đã render
  
  if (!chartCanvas.value) {
    console.error('Canvas element not found')
    return
  }

  const ctx = chartCanvas.value.getContext('2d')
  if (!ctx) {
    console.error('Failed to get canvas context')
    return
  }

  // Clean up previous chart
  if (chartInstance.value) {
    chartInstance.value.destroy()
    chartInstance.value = null
  }

  const data = await fetchRevenueData()
  const { labels, xAxisTitle } = getChartLabels(data)

  const config: ChartConfiguration<'bar'> = {
    type: 'bar',
    data: {
      labels,
      datasets: [{
        label: t('Orders'),
        data: data.map(item => item.count),
        backgroundColor: '#4f46e5',
        borderColor: '#4f46e5',
        borderWidth: 1
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: {
          callbacks: {
            label: (context) => `${context.parsed.y} ${t('orders')}`
          }
        }
      },
      scales: {
        y: {
          beginAtZero: true,
          title: { display: true, text: t('Number of Orders') }
        },
        x: {
          title: { display: true, text: xAxisTitle }
        }
      }
    }
  }

  chartInstance.value = new Chart(ctx, config)
}

const getChartLabels = (data: OrderTrendData[]): { labels: string[]; xAxisTitle: string } => {
  switch (timeRange.value) {
    case 'week':
      return {
        labels: data.map(item => item.date?.split('-').reverse().join('/') || ''),
        xAxisTitle: t('Day')
      }
    case 'month':
      return {
        labels: data.map(item => t('Week') + ` ${item.week_start?.split('-')[2]}`),
        xAxisTitle: t('Week')
      }
    case 'year':
      return {
        labels: data.map(item => item.month?.replace('-', '/') || ''),
        xAxisTitle: t('Month')
      }
    default:
      return { labels: [], xAxisTitle: '' }
  }
}

const fetchStats = async () => {
  try {
    const token = localStorage.getItem('adminToken')
    if (!token) return

    const response = await axios.get(`${import.meta.env.VITE_API_URL}/api/admin/dashboard/total-stats`, {
      headers: { Authorization: `Bearer ${token}` },
      params: { time_range: timeRange.value }
    })

    totalOrders.value = response.data.total_orders || 0
    totalAmount.value = response.data.total_amount || 0
    totalBooksSold.value = response.data.total_books_sold || 0
  } catch (error) {
    console.error('Error fetching stats:', error)
  }
}

const fetchTotalOrders = async () => {
  try {
    const token = localStorage.getItem('adminToken')
    if (!token) return

    const response = await axios.get(`${import.meta.env.VITE_API_URL}/api/admin/dashboard/total-order`, {
      headers: { Authorization: `Bearer ${token}` },
      params: { time_range: timeRange.value }
    })

    trendingBooks.value = response.data.total_orders || 0
  } catch (error) {
    console.error('Error fetching total orders:', error)
  }
}

const onTimeRangeChange = () => {
  isLoading.value = true
  Promise.all([fetchStats(), fetchTotalOrders(), renderChart()])
    .finally(() => isLoading.value = false)
}

// Lifecycle hooks
onMounted(() => {
  const token = localStorage.getItem('adminToken')
  if (!token) router.push('/admin-login')
  Promise.all([fetchStats(), fetchTotalOrders(), renderChart()])
    .finally(() => isLoading.value = false)
})

watch(timeRange, () => {
  isLoading.value = true
  onTimeRangeChange()
})

// Clean up
onBeforeUnmount(() => {
  if (chartInstance.value) {
    chartInstance.value.destroy()
  }
})
</script>

<style>
.chart-container {
  position: relative;
  width: 100%;
  height: 100%;
}

.chart-container canvas {
  display: block;
  width: 100% !important;
  height: 100% !important;
}
</style>

