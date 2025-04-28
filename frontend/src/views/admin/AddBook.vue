<template>
    <section class="flex md:bg-gray-100 min-h-screen overflow-hidden">
      <AdminSidebar />
      
      <div class="flex-grow text-gray-800">
        <header class="flex items-center h-20 px-6 sm:px-10 bg-white">
          <button class="block sm:hidden relative flex-shrink-0 p-2 mr-2 text-gray-600 hover:bg-gray-100 hover:text-gray-800 focus:bg-gray-100 focus:text-gray-800 rounded-full">
            <span class="sr-only">Menu</span>
            <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="h-6 w-6">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7" />
            </svg>
          </button>
          
          <div class="flex-grow">
            <h1 class="text-2xl font-semibold">Add New Book</h1>
          </div>
        </header>
  
        <main class="p-6 sm:p-10 bg-gray-100">
          <div class="max-w-lg mx-auto md:p-6 p-3 bg-white rounded-lg shadow-md">
            <h2 class="text-2xl font-bold text-gray-800 mb-4">Add New Book</h2>
  
            <form @submit.prevent="submitForm" class="space-y-4">
              <!-- Title -->
              <div class="mb-4">
                <label class="block text-sm font-semibold text-gray-700 mb-1">Title</label>
                <input
                  v-model="form.title"
                  type="text"
                  required
                  class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
                  placeholder="Book Title"
                />
              </div>
  
              <!-- Author -->
              <div class="mb-4">
                <label class="block text-sm font-semibold text-gray-700 mb-1">Author</label>
                <input
                  v-model="form.author"
                  type="text"
                  required
                  class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
                  placeholder="Author Name"
                />
              </div>
  
              <!-- Description -->
              <div class="mb-4">
                <label class="block text-sm font-semibold text-gray-700 mb-1">Description</label>
                <textarea
                  v-model="form.description"
                  rows="3"
                  class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
                  placeholder="Book description"
                ></textarea>
              </div>
  
              <!-- Price -->
              <div class="mb-4">
                <label class="block text-sm font-semibold text-gray-700 mb-1">Price ($)</label>
                <input
                  v-model="form.price"
                  type="number"
                  min="0"
                  step="0.01"
                  required
                  class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
                  placeholder="19.99"
                />
              </div>
  
              <!-- Stock -->
              <div class="mb-4">
                <label class="block text-sm font-semibold text-gray-700 mb-1">Stock</label>
                <input
                  v-model="form.stock"
                  type="number"
                  min="0"
                  required
                  class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
                  placeholder="10"
                />
              </div>
  
              <!-- Category -->
              <div class="mb-4">
                <label class="block text-sm font-semibold text-gray-700 mb-1">Category</label>
                <select
                  v-model="form.category"
                  required
                  class="w-full p-2 border rounded-md focus:outline-none focus:ring focus:border-blue-300"
                >
                  <option 
                    v-for="category in categories" 
                    :key="category.id" 
                    :value="category.name"
                  >
                    {{ category.name }}
                  </option>
                </select>
              </div>
  
              <!-- Publisher -->
              <div class="mb-4">
                <label class="block text-sm font-semibold text-gray-700 mb-1">Publisher</label>
                <input
                  v-model="form.publisher"
                  type="text"
                  required
                  class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
                  placeholder="Publisher Name"
                />
              </div>
  
              <!-- ISBN -->
              <div class="mb-4">
                <label class="block text-sm font-semibold text-gray-700 mb-1">ISBN</label>
                <input
                  v-model="form.isbn"
                  type="text"
                  required
                  class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
                  placeholder="9780132350884"
                />
              </div>
  
              <!-- Pages -->
              <div class="mb-4">
                <label class="block text-sm font-semibold text-gray-700 mb-1">Pages</label>
                <input
                  v-model="form.pages"
                  type="number"
                  min="1"
                  required
                  class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
                  placeholder="464"
                />
              </div>
  
              <!-- Language -->
              <div class="mb-4">
                <label class="block text-sm font-semibold text-gray-700 mb-1">Language</label>
                <input
                  v-model="form.language"
                  type="text"
                  required
                  class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
                  placeholder="English"
                />
              </div>
  
            
  
              <!-- Submit Button -->
              <button 
                type="submit" 
                :disabled="isSubmitting"
                class="w-full py-2 bg-green-500 hover:bg-green-600 text-white font-bold rounded-md transition-colors"
                :class="{'opacity-50 cursor-not-allowed': isSubmitting}"
              >
                <span v-if="!isSubmitting">Add Book</span>
                <span v-else>Adding...</span>
              </button>
            </form>
          </div>
        </main>
      </div>
    </section>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import axios from 'axios'
  import { useRouter } from 'vue-router'
  import AdminSidebar from '@/components/AdminSidebar.vue'
  
  interface Category {
    id: number
    name: string
  }
  
  interface BookForm {
    title: string
    author: string
    description: string
    price: number
    stock: number
    category: string
    publisher: string
    isbn: string
    pages: number
    language: string
    cover_image: string
  }
  
  const router = useRouter()
  
  // Form data
  const form = ref<BookForm>({
    title: '',
    author: '',
    description: '',
    price: 0,
    stock: 0,
    category: '',
    publisher: '',
    isbn: '',
    pages: 0,
    language: 'English',
    cover_image: ''
  })
  
  const categories = ref<Category[]>([])
  const selectedFile = ref<File | null>(null)
  const isSubmitting = ref(false)
  
  // Check authentication and fetch categories on mount
  onMounted(() => {
    const token = localStorage.getItem('adminToken')
    if (!token) {
      router.push('/admin-login')
      return
    }
    fetchCategories()
  })
  
  const handleFileChange = (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target.files && target.files[0]) {
      selectedFile.value = target.files[0]
    }
  }
  
  const submitForm = async () => {
    try {
      isSubmitting.value = true
      
      // First upload the image if selected
      let coverImageUrl = ''

  
      // Prepare book data
      const bookData = {
        title: form.value.title,
        author: form.value.author,
        description: form.value.description,
        price: parseFloat(form.value.price.toString()),
        stock: parseInt(form.value.stock.toString()),
        category: form.value.category,
        publisher: form.value.publisher,
        isbn: form.value.isbn,
        pages: parseInt(form.value.pages.toString()),
        language: form.value.language,
        cover_image: coverImageUrl
      }
  
      console.log('Submitting book data:', bookData)
  
      // Submit book data
      const response = await axios.post(
        `${import.meta.env.VITE_API_URL}/api/books`, 
        bookData, 
        {
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${localStorage.getItem('adminToken')}`
          }
        }
      )
  
      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to add book')
      }
  
      // Redirect to books list after successful submission
      router.push('/admin/books')
    } catch (error) {
      console.error('Error adding book:', error)
      alert(error || 'Failed to add book. Please try again.')
    } finally {
      isSubmitting.value = false
    }
  }
  
  const fetchCategories = async () => {
    try {
      const response = await axios.get(
        `${import.meta.env.VITE_API_URL}/api/categories`,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('adminToken')}`
          }
        }
      )
  
      if (response.data.success) {
        categories.value = response.data.data.map((name: string, index: number) => ({
          id: index,
          name: name
        }))
        
        // Set default category if categories exist
        if (categories.value.length > 0) {
          form.value.category = categories.value[0].name
        }
      }
    } catch (error) {
      console.error('Error fetching categories:', error)
      alert('Failed to load categories. Please try again later.')
    }
  }
  </script>