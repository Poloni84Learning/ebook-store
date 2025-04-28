<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="w-full max-w-sm bg-white rounded-lg shadow-md overflow-hidden">
      <!-- Header Section -->
      <div class="bg-gray-800 py-5 px-6 text-center">
        <h1 class="text-xl font-semibold text-white">Management Portal</h1>
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
              class="block text-black w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-gray-500 focus:border-gray-500"
              placeholder="Enter your username"
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
              class="block text-black w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-gray-500 focus:border-gray-500"
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
import { reactive, ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';

interface User {
  id: number;
  username: string;
  email: string;
  role: 'admin' | 'staff';
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

const router = useRouter();
const apiUrl = import.meta.env.VITE_API_URL;

const form = reactive({
  username: '',
  password: ''
});

const error = ref('');
const isLoading = ref(false);

// Kiểm tra token khi component được mount
onMounted(() => {
  checkExistingSession();
});

function checkExistingSession() {
  const token = localStorage.getItem('adminToken') || localStorage.getItem('staffToken');
  const userData = localStorage.getItem('userData');
  
  if (token && userData) {
    const user = JSON.parse(userData);
    redirectBasedOnRole(user.role);
  }
}

function redirectBasedOnRole(role: string) {
  switch (role) {
    case 'admin':
      router.push('/admin/dashboard');
      break;
    case 'staff':
      router.push('/staff/orders');
      break;
    default:
      clearAuthData();
      router.push('/admin-login');
  }
}

function clearAuthData() {
  localStorage.removeItem('adminToken');
  localStorage.removeItem('staffToken');
  localStorage.removeItem('userData');
}

async function handleSubmit() {
  try {
    isLoading.value = true;
    error.value = '';

    const response = await axios.post<LoginResponse>(`${apiUrl}/api/auth/staff-login`, {
      username: form.username,
      password: form.password
    });

    if (response.data.success) {
      const { token, user } = response.data;
      
      // Lưu token và user data tùy theo role
      if (user.role === 'admin') {
        localStorage.setItem('adminToken', token);
      } else {
        localStorage.setItem('staffToken', token);
      }
      
      localStorage.setItem('userData', JSON.stringify(user));
      
      redirectBasedOnRole(user.role);
    } else {
      error.value = response.data.error;
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