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
  <!-- Cover Image -->
<div class="mb-4">
  <label class="block text-sm font-semibold text-gray-700 mb-1">Cover Image</label>
  <input
    type="file"
    accept="image/*"
    @change="handleImageChange"
    class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
  />
</div>
            <!-- PDF File Upload -->
            <div class="mb-4">
              <label class="block text-sm font-semibold text-gray-700 mb-1">Upload PDF</label>
              <input
                type="file"
                accept="application/pdf"
                @change="handleFileChange"
                required
                class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
              />
            </div>
            <div class="mb-4">
    <label class="block text-sm font-semibold text-gray-700 mb-1">Table of Contents Pages</label>
    <input
        v-model="form.toc"
        type="text"
        class="p-2 border w-full rounded-md focus:outline-none focus:ring focus:border-blue-300"
        placeholder="e.g. 3,4,5 (comma separated)"
    />
    <p class="text-xs text-gray-500 mt-1">Enter page numbers containing table of contents, separated by commas</p>
</div>
  
              <!-- Submit Button -->
              <button 
                type="submit" 
                :disabled="isSubmitting"
                class="w-full py-2 bg-green-500 hover:bg-green-600 text-white font-bold rounded-md transition-colors"
                :class="{'opacity-50 cursor-not-allowed': isSubmitting}"
              >
              <span v-if="isSubmitting">Submitting...</span>
              <span v-else-if="isProcessing">Processing...</span>
              <span v-else>Add Book</span>
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
    coverImageFile: File | null
    pdfFile: File | null
    toc: string // Thay đổi từ int[] thành string để phù hợp với input
    tocPages: number[] // Thêm trường mới để lưu dạng mảng số
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
    cover_image: '',
    coverImageFile: null,
    pdfFile: null,
    toc: '', // Người dùng nhập chuỗi như "3,4,5"
    tocPages: [] // Sẽ được xử lý khi submit
  })
  
  const categories = ref<Category[]>([])
  const isSubmitting = ref(false)
  
  const processTocInput = (tocString: string): number[] => {
    // Xử lý chuỗi nhập vào thành mảng số
    return tocString
        .split(',') // Tách theo dấu phẩy
        .map(item => item.trim()) // Loại bỏ khoảng trắng
        .filter(item => item !== '') // Loại bỏ item rỗng
        .map(Number) // Chuyển thành số
        .filter(num => !isNaN(num) && num > 0) // Chỉ lấy số hợp lệ
}

  // Check authentication and fetch categories on mount
  onMounted(() => {
    const token = localStorage.getItem('adminToken')
    if (!token) {
      router.push('/admin-login')
      return
    }
    fetchCategories()
  })
  
  // Khi chọn file PDF
  const handleFileChange = (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target.files && target.files.length > 0) {
      form.value.pdfFile = target.files[0]
    }
  }
  const handleImageChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    form.value.coverImageFile = target.files[0] // Giờ sẽ không còn lỗi
    form.value.cover_image = target.files[0].name // Cập nhật cả tên file nếu cần
  }
}

const isProcessing = ref(false)

  // Submit form
  const submitForm = async () => {
    try {
      isSubmitting.value = true
      isProcessing.value = false
      form.value.tocPages = processTocInput(form.value.toc)
      // Chuẩn bị FormData để gửi file + dữ liệu
      const formData = new FormData()
  
      formData.append('title', form.value.title)
      formData.append('author', form.value.author)
      formData.append('description', form.value.description)
      formData.append('price', form.value.price.toString())
      formData.append('stock', form.value.stock.toString())
      formData.append('category', form.value.category)
      formData.append('publisher', form.value.publisher)
      formData.append('isbn', form.value.isbn)
      formData.append('pages', form.value.pages.toString())
      formData.append('language', form.value.language)
      formData.append('toc',form.value.toc)
      // Thêm file ảnh nếu có
      if (form.value.coverImageFile) {
        formData.append('cover_image', form.value.coverImageFile)
      }
      // Append file PDF nếu có
      if (form.value.pdfFile) {
        formData.append('pdf', form.value.pdfFile) // 👈 tên 'pdf' bạn đặt theo API backend yêu cầu
      }
      isProcessing.value = true
      console.log('Submitting book data with file:', formData)
      console.log('🔍 FormData content:')
    for (const [key, value] of formData.entries()) {
      console.log(`${key}:`, value)
    }
      const response = await axios.post(
        `${import.meta.env.VITE_API_URL}/api/books`,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
            Authorization: `Bearer ${localStorage.getItem('adminToken')}`,
          }
        }
      )
  
      if (!response.data.success) {
        throw new Error(response.data.message || 'Failed to add book')
      }
  
      // Redirect thành công
      isProcessing.value = false
      router.push('/admin/books')
    } catch (error: any) {
      console.error('Error adding book:', error)
      alert(error?.message || 'Failed to add book. Please try again.')
    } finally {
      isSubmitting.value = false
      isProcessing.value = false
    }
  }
  
  // Fetch categories
  const fetchCategories = async () => {
    try {
      const response = await axios.get(
        `${import.meta.env.VITE_API_URL}/api/categories`,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('adminToken')}`,
          }
        }
      )
  
      if (response.data.success) {
        categories.value = response.data.data.map((name: string, index: number) => ({
          id: index,
          name: name
        }))
  
        // Set default
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
  