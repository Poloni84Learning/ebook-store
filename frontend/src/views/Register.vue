<template>
    <div class="min-h-screen flex items-center justify-center bg-gray-50 p-4">
      <div class="w-full max-w-md bg-white rounded-xl shadow-lg overflow-hidden">
        <!-- Header Section -->
        <div class="bg-gradient-to-r from-secondary to-indigo-800 py-6 px-8 text-center">
          <h1 class="text-2xl font-bold text-white">Create Account</h1>
          <p class="text-white/90 mt-1">Sign up for your Book Store account</p>
        </div>
  
        <!-- Form Section -->
        <div class="p-8">
          <form @submit.prevent="handleSubmit" class="space-y-6">
            <!-- Username Field -->
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
                  type="text"
                  required
                  class="block w-full pl-10 pr-3 py-3 border border-black rounded-lg focus:ring-2 focus:ring-black focus:border-black placeholder-gray-400 transition text-black"
                  placeholder="Choose a username"
                  @input="clearError"
                />
              </div>
            </div>
  
            <!-- Email Field -->
            <div>
              <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12h2a2 2 0 002-2V7a2 2 0 00-2-2h-2M8 12H6a2 2 0 01-2-2V7a2 2 0 012-2h2m0 0h8m-8 0v4m8-4v4"/>
                  </svg>
                </div>
                <input
                  id="email"
                  v-model="form.email"
                  type="email"
                  required
                  class="block w-full pl-10 pr-3 py-3 border border-black rounded-lg focus:ring-2 focus:ring-black focus:border-black placeholder-gray-400 transition text-black"
                  placeholder="you@example.com"
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
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242"/>
                  </svg>
                  <svg v-else class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.522 5 12 5c4.478 0 8.268 2.943 9.543 7-1.275 4.057-5.065 7-9.543 7-4.478 0-8.268-2.943-9.543-7z"/>
                  </svg>
                </button>
              </div>
            </div>
  
            <!-- Confirm Password Field -->
            <div>
              <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">Confirm Password</label>
              <input
                id="confirmPassword"
                v-model="form.confirmPassword"
                :type="showPassword ? 'text' : 'password'"
                required
                class="block w-full py-3 px-3 border border-black rounded-lg focus:ring-2 focus:ring-black focus:border-black placeholder-gray-400 transition text-black"
                placeholder="••••••••"
                @input="clearError"
              />
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
                Sign Up
              </button>
            </div>
          </form>
  
          <!-- Login Link -->
          <div class="mt-6 text-center text-sm">
            <p class="text-gray-600">
              Already have an account?
              <a href="/login" class="font-medium text-primary hover:text-secondary">
                Sign in
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
  
  const apiUrl = import.meta.env.VITE_API_URL;
  
  const router = useRouter();
  
  const form = reactive({
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  
  const showPassword = ref(false);
  const error = ref('');
  const isLoading = ref(false);
  
  function clearError() {
    error.value = '';
  }
  onMounted(() => {
  const token = localStorage.getItem('authToken');
  if (token) {
    router.push('/');
  }
});
  async function handleSubmit() {
  if (/^(admin|admin_|staff_)/i.test(form.username)) {

    error.value = 'This username is not allowed for login';
    return;
}
    try {
      isLoading.value = true;
      error.value = '';
  
      if (form.password !== form.confirmPassword) {
        error.value = 'Passwords do not match';
        return;
      }
  
      const response = await axios.post(`${apiUrl}/api/auth/register`, {
        username: form.username,
        email: form.email,
        password: form.password
      });
  
      // Giả sử đăng ký thành công thì chuyển về trang login
      router.push('/login');
    } catch (err) {
      if (axios.isAxiosError(err)) {
        error.value = err.response?.data?.error || 'Server error';
      } else if (err instanceof Error) {
        error.value = err.message;
      } else {
        error.value = 'Unknown error';
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
  