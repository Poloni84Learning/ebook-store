<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import DefaultAvatar from '@/assets/img/images_placeholder.png'

const router = useRouter()
const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:3000'

interface UserData {
  name: string
  role: string
  // Add other user properties as needed
}

const props = defineProps({
  userName: {
    type: String,
    default: 'Pham Anh Tu'
  },
  userRole: {
    type: String,
    default: 'Admin'
  }
})

const searchQuery = ref('')
const showUserMenu = ref(false)

const onSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ 
      path: '/books/search', 
      query: { input: searchQuery.value } 
    })
    searchQuery.value = ''
  }
}

const logout = async () => {
  try {
    const token = localStorage.getItem('adminToken') || localStorage.getItem('staffToken')
    
    if (token) {
      await axios.post(`${apiUrl}/api/auth/logout`, null, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      })
    }
    
    // Clear all auth-related data
    localStorage.removeItem('adminToken')
    localStorage.removeItem('staffToken')
    localStorage.removeItem('userData')
    
    // Redirect to login page
    router.push('/login').then(() => {
      window.location.reload() // Ensure all state is cleared
    })
    
  } catch (error) {
    console.error('Logout failed:', error)
    // Still clear local storage even if API call fails
    localStorage.removeItem('adminToken')
    localStorage.removeItem('staffToken')
    localStorage.removeItem('userData')
    router.push('/admin-login')
  }
}

// Fetch user data if not passed as props
const userData = ref<UserData | null>(null)
const fetchUserData = async () => {
  try {
    const token = localStorage.getItem('adminToken') || localStorage.getItem('staffToken')
    if (!token) return

    const response = await axios.get(`${apiUrl}/api/user/profile`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    if (response.data.success) {
      userData.value = response.data.user
      // Store user data in localStorage if needed
      localStorage.setItem('userData', JSON.stringify(response.data.user))
    }
  } catch (error) {
    console.error('Failed to fetch user data:', error)
  }
}

// Call fetchUserData on component mount if needed
//  fetchUserData()
</script>

<template>
  <header class="flex items-center h-20 px-6 sm:px-10 bg-white">
    <!-- Menu button (mobile) -->
    <button 
      class="block sm:hidden relative flex-shrink-0 p-2 mr-2 text-gray-600 hover:bg-gray-100 hover:text-gray-800 focus:bg-gray-100 focus:text-gray-800 rounded-full"
      aria-label="Menu"
    >
      <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="h-6 w-6">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7" />
      </svg>
    </button>
    
    <!-- Search bar -->
    <div class="relative w-full max-w-md sm:-ml-2">
      <svg aria-hidden="true" viewBox="0 0 20 20" fill="currentColor" class="absolute h-6 w-6 mt-2.5 ml-2 text-gray-400">
        <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
      </svg>
      <input 
        v-model="searchQuery"
        @keyup.enter="onSearch"
        type="text" 
        role="search" 
        placeholder="Search..." 
        class="py-2 pl-10 pr-4 w-full border-2 border-black placeholder-gray-400 focus:bg-gray-50 focus:border-purple-500 rounded-lg" 
      />
    </div>
    
    <!-- User menu -->
    <div class="flex flex-shrink-0 items-center ml-auto">
      <div class="relative">
        <button 
          @click="showUserMenu = !showUserMenu"
          class="inline-flex items-center p-2 hover:bg-gray-100 focus:bg-gray-100 rounded-lg"
        >
          <span class="sr-only">User Menu</span>
          <div class="hidden md:flex md:flex-col md:items-end md:leading-tight">
            <span class="text-black font-semibold">{{ userData?.name || userName }}</span>
            <span class="text-sm text-gray-600">{{ userData?.role || userRole }}</span>
          </div>
          <span class="h-8 w-8 ml-2 sm:ml-3 mr-2 bg-gray-100 rounded-full overflow-hidden">
            <img :src="DefaultAvatar" alt="user profile photo" class="h-full w-full object-cover"/>
          </span>
          <svg aria-hidden="true" viewBox="0 0 20 20" fill="currentColor" class="hidden sm:block h-6 w-6 text-gray-300">
            <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg> 
        </button>
        
        <!-- Dropdown menu -->
        <div 
          v-if="showUserMenu"
          class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-50"
        >
          <button 
            @click="logout"
            class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 w-full text-left flex items-center"
          >
            <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="h-5 w-5 mr-2">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
            Logout
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<style scoped>
/* Add custom styles if needed */
</style>