<template>
    <div class="min-h-screen flex items-center justify-center bg-gray-100">
      <div class="w-full max-w-sm bg-white rounded-lg shadow-md overflow-hidden">
        <!-- Header Section - Đơn giản hơn -->
        <div class="bg-gray-800 py-5 px-6 text-center">
          <h1 class="text-xl font-semibold text-white">Admin Portal</h1>
        </div>
  
        <!-- Form Section -->
        <div class="p-6">
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <!-- Username Field -->
            <div>
              <label for="username" class="block text-sm font-medium text-gray-700 mb-1">Username</label>
              <input
                id="username"
                v-model="form.username"
                name="username"
                type="text"
                required
                class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-gray-500 focus:border-gray-500"
                placeholder="Admin username"
                autofocus
              />
            </div>
  
            <!-- Password Field -->
            <div>
              <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
              <input
                id="password"
                v-model="form.password"
                name="password"
                type="password"
                required
                class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-gray-500 focus:border-gray-500"
                placeholder="••••••••"
              />
            </div>
  
            <!-- Error Message -->
            <div v-if="error" class="text-sm text-red-600">
              {{ error }}
            </div>
  
            <!-- Submit Button -->
            <div>
              <button
                type="submit"
                class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-800 hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
                :disabled="isLoading"
              >
                <span v-if="!isLoading">Sign In</span>
                <span v-else>Processing...</span>
              </button>
            </div>
          </form>
        </div>
  
        
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { reactive, ref } from 'vue';
  import { useRouter } from 'vue-router';
  import axios from 'axios';
  
  const router = useRouter();
  const apiUrl = import.meta.env.VITE_API_URL;
  
  const form = reactive({
    username: '',
    password: ''
  });
  
  const error = ref('');
  const isLoading = ref(false);
  
  async function handleSubmit() {
    // Kiểm tra username phải bắt đầu bằng admin_
    if (!form.username.startsWith('admin_')) {
      error.value = 'Only admin accounts can access this portal';
      return;
    }
  
    try {
      isLoading.value = true;
      error.value = '';
  
      const response = await axios.post(`${apiUrl}/api/auth/staff-login`, {
        username: form.username,
        password: form.password
      });
  
      if (response.data.success) {
        localStorage.setItem('adminToken', response.data.token);
        localStorage.setItem('adminData', JSON.stringify(response.data.user));
        router.push('/admin/dashboard');
      } else {
        error.value = response.data.error || 'Authentication failed';
      }
    } catch (err) {
      error.value = axios.isAxiosError(err) 
        ? err.response?.data?.error || 'Server error'
        : 'Network error';
    } finally {
      isLoading.value = false;
    }
  }
  </script>
  
  <style scoped>
  /* Tối giản style, chỉ giữ những gì cần thiết */
  </style>