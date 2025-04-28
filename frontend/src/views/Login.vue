<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 p-4">
    <div class="w-full max-w-md bg-white rounded-xl shadow-lg overflow-hidden">
      <!-- Header Section -->
      <div class="bg-gradient-to-r from-secondary to-indigo-800 py-6 px-8 text-center">
        <h1 class="text-2xl font-bold text-white">Welcome Back</h1>
        <p class="text-white/90 mt-1">Sign in to your Book Store account</p>
      </div>

      <!-- Form Section -->
      <div class="p-8">
        <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Email Field -->
        <div>
          <label for="username" class="block text-sm font-medium text-gray-700 mb-1">Username</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
              </svg>
            </div>
            <input
              id="username"
              v-model="form.username"
              name="username"
              type="text"
              required
              class="block w-full pl-10 pr-3 py-3 border border-black rounded-lg focus:ring-2 focus:ring-black focus:border-black placeholder-gray-400 transition text-black"
              placeholder="Your user name"
              @input="clearError"
            />
          </div>
        </div>

        <!-- Password Field -->
        <div>
          <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
              </svg>
            </div>
            <input
              id="password"
              v-model="form.password"
              name="password"
              :type="showPassword ? 'text' : 'password'"
              required
              class="block w-full pl-10 pr-10 py-3 border border-black rounded-lg focus:ring-2 focus:ring-black focus:border-black placeholder-gray-400 transition text-black"
              placeholder="••••••••"
              @input="clearError"
            />
            <button
              type="button"
              class="absolute inset-y-0 right-0 pr-3 flex items-center"
              @click="showPassword = !showPassword"
            >
              <svg v-if="showPassword" class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
              </svg>
              <svg v-else class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.522 5 12 5c4.478 0 8.268 2.943 9.543 7-1.275 4.057-5.065 7-9.543 7-4.478 0-8.268-2.943-9.543-7z"/>
              </svg>
            </button>
          </div>
        </div>

          <!-- Error Message -->
          <div v-if="error" class="text-sm text-red-500">
            {{ error }}
          </div>

          <!-- Submit Button -->
          <div>
            <button
              type="submit"
              class="w-full flex justify-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-secondary hover:bg-indigo-800 focus:outline-none focus:ring-2 focus:ring-secondary transition"
            >
              Sign In
            </button>
          </div>
        </form>

        <!-- Registration Link -->
        <div class="mt-6 text-center text-sm">
          <p class="text-gray-600">
            Don't have an account?
            <a href="/register" class="font-medium text-primary hover:text-secondary">
              Sign up
            </a>
          </p>
        </div>
      </div>

      <!-- Footer -->
      <div class="bg-gray-50 px-8 py-4 text-center">
        <p class="text-xs text-gray-500">
          &copy; 2025 EBookStore. All rights reserved.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">

import { reactive, ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

import axios from 'axios';

interface User {
  id: number;
  email: string;
  username: string;
  role: string;
}

interface LoginSuccessResponse {
  message: string;
  success: true;
  token: string;
  user: User;
}

interface LoginErrorResponse {
  error: string;
  success: false;
}

type LoginResponse = LoginSuccessResponse | LoginErrorResponse;

const apiUrl = import.meta.env.VITE_API_URL

const router = useRouter();

function clearError() {
  error.value = ''
}
const form = reactive({
  username: '',
  password: ''
});

const showPassword = ref(false);
const error = ref('');
const isLoading = ref(false);

async function handleSubmit() {
  if (form.username.toLowerCase() === 'admin') {
    error.value = 'Username "admin" is not allowed for login';
    return;
  }
  try {
    isLoading.value = true;
    error.value = '';

    const response = await axios.post<LoginResponse>(`${apiUrl}/api/auth/login`, {
      username: form.username,
      password: form.password
    });

    if (response.data.success) {
      // Xử lý khi đăng nhập thành công
      const { token, user } = response.data;
      
      // Lưu token vào localStorage hoặc Vuex/Pinia
      localStorage.setItem('authToken', token);
      localStorage.setItem('user', JSON.stringify(user));
      
      
        router.push('/');

    } else {
      // Xử lý khi đăng nhập thất bại
      error.value = response.data.error;
    }
  } catch (err) {
    // Xử lý lỗi network hoặc lỗi server
    if (axios.isAxiosError(err)) {
      // Lỗi từ phản hồi API (4xx/5xx)
      error.value = (err.response?.data as LoginErrorResponse)?.error || 'Lỗi hệ thống';
    } else if (err instanceof Error) {
      error.value = err.message;
    } else {
      error.value = 'Unknown Error';
    }
  } finally {
    isLoading.value = false;
  }
}
</script>

<style scoped>
:root {
  --secondary: #0D0842;
}

/* Áp dụng màu custom */
.bg-secondary {
  background-color: var(--secondary);
}

.text-primary {
  @apply text-blue-600;
}

.text-secondary {
  color: var(--secondary);
}
</style>
