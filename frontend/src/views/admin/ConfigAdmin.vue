<template>
  <section class="flex md:bg-gray-100 min-h-screen overflow-hidden">
    <!-- Sidebar -->
    <AdminSidebar />

    <!-- Main Content -->
    <div class="flex-grow text-gray-800">
      <!-- Header -->
      <AppHeader />

      <!-- Configuration Content -->
      <main class="p-6 sm:p-10 space-y-6">
        <div class="flex flex-col space-y-6 md:space-y-0 md:flex-row justify-between">
          <div class="mr-6">
            <h1 class="text-4xl font-semibold mb-2">System Configuration</h1>
            <h2 class="text-gray-600 ml-0.5">Manage global settings</h2>
          </div>
        </div>

        <!-- Configuration Form -->
        <section class="grid md:grid-cols-1 xl:grid-cols-1 gap-6">
          <div class="flex flex-col bg-white shadow rounded-lg">
            <div class="px-6 py-5 border-b border-gray-100">
              <h3 class="text-lg font-semibold text-gray-900">Configuration Settings</h3>
            </div>

            <div class="p-6">
              <div v-if="loading" class="flex justify-center items-center py-12">
                <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-500"></div>
              </div>

              <form v-else @submit.prevent="submitForm" class="space-y-6">
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label for="shippingfee" class="block text-sm font-medium text-gray-700 mb-1">Shipping Fee ($)</label>
                    <input
                      v-model="form.ShippingFee"
                      type="number"
                      id="shippingfee"
                      min="0"
                      step="0.01"
                      class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-600"
                      required
                    />
                  </div>

                  <div>
                    <label for="promotion" class="block text-sm font-medium text-gray-700 mb-1">Promotion (%)</label>
                    <input
                      v-model="form.Promotion"
                      type="number"
                      id="promotion"
                      min="0"
                      max="100"
                      class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-600"
                      required
                    />
                  </div>
                </div>

                <div>
                  <label for="promotioninfo" class="block text-sm font-medium text-gray-700 mb-1">Promotion Information</label>
                  <textarea
                    v-model="form.PromotionInfo"
                    id="promotioninfo"
                    rows="3"
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-600"
                  ></textarea>
                </div>

                <div class="flex justify-end space-x-4">
                  <button
                    v-if="configExists"
                    @click="deleteConfig"
                    type="button"
                    class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500"
                  >
                    Delete Configuration
                  </button>
                  
                  <button
                    type="submit"
                    class="px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500"
                  >
                    {{ configExists ? 'Update Configuration' : 'Create Configuration' }}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </section>

        <section class="text-right font-semibold text-gray-500">
          <a href="#" class="text-purple-600 hover:underline">EBook Store</a> System Configuration
        </section>
      </main>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import AdminSidebar from '@/components/AdminSidebar.vue'
import AppHeader from '@/components/HeaderAdmin.vue'

const router = useRouter()

interface SystemConfig {
  ID?: number
  CreatedAt?: string
  UpdatedAt?: string
  DeletedAt?: string | null
  ShippingFee: number
  Promotion: number
  PromotionInfo: string
}

// Refs
const loading = ref(true)
const configExists = ref(false)
const form = ref<SystemConfig>({
  ShippingFee: 0,
  Promotion: 0,
  PromotionInfo: ''
})

// Fetch config from API
const fetchConfig = async () => {
  try {
    const token = localStorage.getItem('adminToken')
    if (!token) {
      router.push('/admin-login')
      return
    }

    const response = await axios.get<SystemConfig>(`${import.meta.env.VITE_API_URL}/api/admin/system-config`, {
      headers: { Authorization: `Bearer ${token}` }
    })

    if (response.data && response.data.ID) {
      form.value = {
        ShippingFee: response.data.ShippingFee,
        Promotion: response.data.Promotion,
        PromotionInfo: response.data.PromotionInfo
      }
      configExists.value = true
    } else {
      configExists.value = false
      // Reset form to default values if no config exists
      form.value = {
        ShippingFee: 0,
        Promotion: 0,
        PromotionInfo: ''
      }
    }
  } catch (error) {
    console.error('Fetch config error:', error)
    // Reset form to default values on error
    form.value = {
      ShippingFee: 0,
      Promotion: 0,
      PromotionInfo: ''
    }
    configExists.value = false
  } finally {
    loading.value = false
  }
}

// Submit form (create or update)
const submitForm = async () => {
  try {
    const token = localStorage.getItem('adminToken')
    if (!token) {
      router.push('/admin-login')
      return
    }

    loading.value = true

    const request = configExists.value
      ? axios.put<SystemConfig>(`${import.meta.env.VITE_API_URL}/api/admin/system-config`, {
          ShippingFee: form.value.ShippingFee,
          Promotion: form.value.Promotion,
          PromotionInfo: form.value.PromotionInfo
        }, {
          headers: { Authorization: `Bearer ${token}` }
        })
      : axios.post<SystemConfig>(`${import.meta.env.VITE_API_URL}/api/admin/system-config`, {
          ShippingFee: form.value.ShippingFee,
          Promotion: form.value.Promotion,
          PromotionInfo: form.value.PromotionInfo
        }, {
          headers: { Authorization: `Bearer ${token}` }
        })

    const response = await request
    configExists.value = true
    form.value = {
      ShippingFee: response.data.ShippingFee,
      Promotion: response.data.Promotion,
      PromotionInfo: response.data.PromotionInfo
    }
    
    alert('Configuration saved successfully!')
  } catch (error) {
    console.error('Error saving configuration:', error)
    alert('Failed to save configuration. Please try again.')
  } finally {
    loading.value = false
  }
}

// Delete configuration
const deleteConfig = async () => {
  if (!confirm('Are you sure you want to delete the system configuration?')) {
    return
  }

  try {
    const token = localStorage.getItem('adminToken')
    if (!token) {
      router.push('/admin-login')
      return
    }

    loading.value = true

    await axios.delete(
      `${import.meta.env.VITE_API_URL}/api/admin/system-config`,
      {
        headers: { Authorization: `Bearer ${token}` }
      }
    )

    // Reset form
    form.value = {
      ShippingFee: 0,
      Promotion: 0,
      PromotionInfo: ''
    }
    configExists.value = false
    alert('Configuration deleted successfully!')
  } catch (error) {
    console.error('Error deleting configuration:', error)
    alert('Failed to delete configuration. Please try again.')
  } finally {
    loading.value = false
  }
}

// Lifecycle hooks
onMounted(() => {
  fetchConfig()
})
</script>