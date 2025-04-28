<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import DefaultAvatar from '@/assets/img/images_placeholder.png' // avatar mặc định
import axios from 'axios' // nhớ đã cài `axios`

const { t, locale } = useI18n()
const router = useRouter()
const apiUrl = import.meta.env.VITE_API_URL
const isAuthenticated = ref(false)
const user = ref<any>(null)
const showDropdown = ref(false)
const showSearch = ref(false)
const searchQuery = ref('')

const searchInputRef = ref<HTMLElement | null>(null)
const searchButtonRef = ref<HTMLElement | null>(null)

const fetchUserInfo = async () => {
  try {
    const token = localStorage.getItem('authToken')
    if (!token) return

    const res = await axios.get(`${apiUrl}/api/user/profile`, { // <-- chỉnh URL nếu khác
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (res.data.success) {
      isAuthenticated.value = true
      user.value = res.data.user
    } else {
      isAuthenticated.value = false
      user.value = null
    }
  } catch (error) {
    console.error('Error fetching user info:', error)
    isAuthenticated.value = false
    user.value = null
  }
}
const onSearch = () => {
  if (searchQuery.value.trim()) {
    // Điều hướng đến trang tìm kiếm với tham số query 'input'
    router.push({ path: '/books/search', query: { input: searchQuery.value } })
    searchQuery.value = ''
    showSearch.value = false
  }
}
const logout = () => {
  localStorage.removeItem('authToken')
  localStorage.removeItem('user')
  router.push('/').then(() => {
    window.location.reload()
  })
}
const dropdownRef = ref<HTMLElement | null>(null)

const handleClickOutside = (event: MouseEvent) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    showDropdown.value = false
  }
  if (searchInputRef.value && !searchInputRef.value.contains(event.target as Node) && 
      !searchButtonRef.value?.contains(event.target as Node)) {
    showSearch.value = false
  }

}

onMounted(() => {
  fetchUserInfo()
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})


const categories = ref<any[]>([])
const showCategoryDropdown = ref(false)

const fetchCategories = async () => {
  try {
    const res = await axios.get(`${apiUrl}/api/categories`) 
    if (res.data.success) { 
      categories.value = res.data.data.map((name: string, index: number) => ({
        id: index,
        name: name
      }))
    }
  } catch (error) {
    console.error('Error fetching categories:', error)
  }
}


onMounted(() => {
  fetchUserInfo()
  fetchCategories()
  document.addEventListener('click', handleClickOutside)
})

</script>

<template>
  <nav class="bg-secondary text-white p-4 shadow-md">
    <div class="container mx-auto flex justify-between items-center">
      <router-link to="/" class="text-2xl font-bold">EBookStore</router-link>

      <div class="hidden md:flex space-x-6">
      <!-- Home -->
      <div class="relative group">
        <router-link
          to="/"
          class="block px-4 py-2 rounded hover:bg-primary hover:text-white transition"
        >
          {{ t('home.title') }}
        </router-link>
      </div>

      <!-- Books -->
      <div class="relative group">
        <router-link
          to="/books"
          class="block px-4 py-2 rounded hover:bg-primary hover:text-white transition"
        >
          {{ t('books') }}
        </router-link>
      </div>

      <!-- Categories with Dropdown -->
      <div
        class="relative group"
        @mouseenter="showCategoryDropdown = true"
        @mouseleave="showCategoryDropdown = false"
      >
        <router-link
          to="/categories"
          class="block px-4 py-2 rounded hover:bg-primary hover:text-white transition"
        >
          {{ t('categories') }}
        </router-link>

        <!-- Dropdown -->
        <div
        v-if="showCategoryDropdown"
        class="absolute left-0 top-full bg-white text-black rounded-md shadow-lg z-50 p-4"
        style="min-width: 300px;" 
        >
        <div class="grid grid-cols-2 gap-4">
          <router-link
            v-for="cat in categories"
            :key="cat.id"
            :to="`/categories/${cat.name}`"
            class="block hover:bg-gray-100 px-2 py-1 rounded"
            @click="showCategoryDropdown = false"
          >
            {{ cat.name }}
          </router-link>
        </div>
      </div>

</div>

    </div>

      
      <div class="flex items-center space-x-4">

        <div class="relative">
          <!-- Nút kính lúp -->
          <button 
            ref="searchButtonRef"
            @click.stop="showSearch = !showSearch" 
            class="p-2 rounded-full hover:bg-gray-700"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </button>

          <!-- Thanh tìm kiếm - đặt bên trái nút kính lúp -->
          <transition name="search-slide">
            <div 
              v-if="showSearch" 
              ref="searchInputRef" 
              class="absolute right-full top-0 mr-2"
              @click.stop
            >
              <input
                type="text"
                placeholder="Search books, authors, categories..."
                class="min-w-[250px] px-4 py-2 rounded-md border text-black focus:outline-none focus:ring focus:border-primary"
                v-model="searchQuery"
                @keyup.enter="onSearch"
                />
            </div>
          </transition>
        </div>
        <!-- Nút giỏ hàng -->
        <button @click="router.push('/cart')" class="p-2 rounded-full hover:bg-gray-700">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
          </svg>
        </button>

        <div class="relative">
        <!-- Nếu đã đăng nhập -->
        <template v-if="isAuthenticated">
          <img
            :src="user?.avatar_url || DefaultAvatar"
            alt="User Avatar"
            class="w-10 h-10 rounded-full border-2 border-white cursor-pointer hover:opacity-90 transition"
            @click.stop="showDropdown = !showDropdown"
          />

          
          <!-- Dropdown Menu -->
          <transition name="fade">
            <div
              v-if="showDropdown"
              ref="dropdownRef"
              class="absolute w-48 bg-white rounded-md overflow-hidden shadow-xl z-50"
              style="left: 0; right: 0; top: calc(100% + 8px); margin-left: auto; margin-right: auto;"
            >
              <router-link
                to="/profile"
                class="block px-4 py-2 text-gray-800 hover:bg-gray-100 "
                @click="showDropdown = false"
              >
              {{ t('profile') }}
              </router-link>
              <router-link
                to="/order-history"
                class="block px-4 py-2 text-gray-800 hover:bg-gray-100 "
                @click="showDropdown = false"
              >
              {{ t('ordersHistory') }}

              </router-link>
              <router-link
                to="/"
                @click="logout"
                class="block text-left px-4 py-2 text-gray-800 hover:bg-gray-100"
              >
              {{ t('logout') }}
              </router-link>
            </div>
          </transition>

        </template>

        <!-- Nếu chưa đăng nhập -->
        <template v-else>
          <router-link to="/login" class="bg-primary text-secondary px-4 py-2 rounded-md font-bold">
            {{ t('login') }}
          </router-link>
        </template>
      </div>

        <!-- Language Switcher -->
        <select v-model="locale" class="border px-2 py-1 rounded bg-white text-black">
          <option value="en">English</option>
          <option value="ru">Русский</option>
        </select>
      </div>
    </div>
  </nav>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

.fade-enter-active, .fade-leave-active {
  transition: all 0.3s ease;
}
.fade-enter-from {
  opacity: 0;
  transform: translateX(20px);
}
.fade-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

/* Animation cho thanh search */
.search-slide-enter-active, .search-slide-leave-active {
  transition: all 0.3s ease;
}
.search-slide-enter-from {
  opacity: 0;
  transform: translateX(10px);
}
.search-slide-leave-to {
  opacity: 0;
  transform: translateX(10px);
}

</style>
