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
          
          <div class="relative w-full max-w-md sm:-ml-2">
            <svg aria-hidden="true" viewBox="0 0 20 20" fill="currentColor" class="absolute h-6 w-6 mt-2.5 ml-2 text-gray-400">
              <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
            </svg>
            <input 
              v-model="searchQuery"
              type="text" 
              placeholder="Search books..." 
              class="py-2 pl-10 pr-4 w-full border-4 border-transparent placeholder-gray-400 focus:bg-gray-50 rounded-lg"
              @input="searchBooks"
            />
          </div>
          
          <div class="flex flex-shrink-0 items-center ml-auto">
            <router-link to="/admin/books" class="inline-flex px-5 py-3 text-purple-600 hover:text-purple-700 focus:text-purple-700 hover:bg-purple-100 focus:bg-purple-100 border border-purple-600 rounded-md mr-5">
              <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="flex-shrink-0 h-5 w-5 -ml-1 mt-0.5 mr-2">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
              </svg>
              View Ranking
            </router-link>
            <router-link 
              to="/admin/books/add" 
              class="inline-flex px-5 py-3 text-white bg-purple-600 hover:bg-purple-700 focus:bg-purple-700 rounded-md ml-6"
            >
              <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="flex-shrink-0 h-6 w-6 text-white -ml-1 mr-2">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Add New Book
            </router-link>
          </div>
        </header>
  
        <main class="p-6 sm:p-10 space-y-6 bg-gray-100">
          <div class="bg-white shadow rounded-lg overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-200">
              <h2 class="text-lg font-semibold text-gray-800">Book Management</h2>
            </div>
            
            <div class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Cover</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Title</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Author</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stock</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="book in filteredBooks" :key="book.id">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <img :src="getCorrectImageUrl(book.cover_image)" 
                           alt="Book cover" 
                           class="h-12 w-9 object-cover rounded"
                           >
                    </td>
                    
                    <!-- Title -->
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div v-if="!book.editing" class="text-sm font-medium text-gray-900">{{ book.title }}</div>
                      <input v-else v-model="book.editData.title" type="text" class="w-full p-1 border rounded">
                    </td>
                    
                    <!-- Author -->
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div v-if="!book.editing" class="text-sm text-gray-500">{{ book.author }}</div>
                      <input v-else v-model="book.editData.author" type="text" class="w-full p-1 border rounded">
                    </td>
                    
                    <!-- Price -->
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div v-if="!book.editing" class="text-sm text-gray-500">${{ book.price }}</div>
                      <input v-else v-model="book.editData.price" type="number" step="0.01" class="w-full p-1 border rounded">
                    </td>
                    
                    <!-- Stock -->
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div v-if="!book.editing" class="text-sm text-gray-500">{{ book.stock }}</div>
                      <input v-else v-model="book.editData.stock" type="number" class="w-full p-1 border rounded">
                    </td>
                    
                    <!-- Actions -->
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <div class="flex space-x-2">
                        <button v-if="!book.editing" 
                                @click="startEditing(book)"
                                class="text-indigo-600 hover:text-indigo-900">
                          Edit
                        </button>
                        <button v-if="!book.editing" 
                                @click="confirmDelete(book)"
                                class="text-red-600 hover:text-red-900">
                          Hide
                        </button>
                        
                        <button v-if="book.editing" 
                                @click="saveChanges(book)"
                                class="text-green-600 hover:text-green-900">
                          Save
                        </button>
                        <button v-if="book.editing" 
                                @click="cancelEditing(book)"
                                class="text-gray-600 hover:text-gray-900">
                          Cancel
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <!-- Pagination -->
            <div class="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
              <div class="text-sm text-gray-500">
                Showing <span class="font-medium">{{ currentPage * itemsPerPage - itemsPerPage + 1 }}</span> to 
                <span class="font-medium">{{ Math.min(currentPage * itemsPerPage, totalBooks) }}</span> of 
                <span class="font-medium">{{ totalBooks }}</span> books
              </div>
              <div class="flex space-x-2">
                <button @click="prevPage" :disabled="currentPage === 1" 
                        class="px-3 py-1 border rounded-md" 
                        :class="{'opacity-50 cursor-not-allowed': currentPage === 1}">
                  Previous
                </button>
                <button @click="nextPage" :disabled="currentPage * itemsPerPage >= totalBooks" 
                        class="px-3 py-1 border rounded-md" 
                        :class="{'opacity-50 cursor-not-allowed': currentPage * itemsPerPage >= totalBooks}">
                  Next
                </button>
              </div>
            </div>
          </div>
        </main>
      </div>
  
      <!-- Delete Confirmation Modal -->
      <div v-if="showDeleteModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
        <div class="bg-white rounded-lg p-6 max-w-md w-full">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Confirm Delete</h3>
          <p class="text-gray-600 mb-6">Are you sure you want to delete "{{ bookToDelete?.title }}"? This action cannot be undone.</p>
          <div class="flex justify-end space-x-3">
            <button @click="showDeleteModal = false" class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50">
              Cancel
            </button>
            <button @click="deleteBook" class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700">
              Delete
            </button>
          </div>
        </div>
      </div>
    </section>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import axios from 'axios'
  import { useRouter } from 'vue-router'
  import AdminSidebar from '@/components/AdminSidebar.vue'
  
  const router = useRouter()
  
  // Data
  const books = ref<any[]>([])
  const searchQuery = ref('')
  const currentPage = ref(1)
  const itemsPerPage = 10
  const showDeleteModal = ref(false)
  const bookToDelete = ref<any>(null)
  
  // Computed properties
  const totalBooks = computed(() => books.value.length)
  const filteredBooks = computed(() => {
    const start = (currentPage.value - 1) * itemsPerPage
    const end = start + itemsPerPage
    return books.value
      .filter(book => 
        book.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
        book.author.toLowerCase().includes(searchQuery.value.toLowerCase())
      .slice(start, end))
  })
  
  // Methods
  const fetchBooks = async () => {
    try {
      const token = localStorage.getItem('adminToken')
      if (!token) {
        router.push('/login')
        return
      }
  
      const response = await axios.get(`${import.meta.env.VITE_API_URL}/api/books?limit=50`, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      })
  
      books.value = response.data.data.map((book: any) => ({
        ...book,
        editing: false,
        editData: { ...book }
      }))
    } catch (error) {
      console.error('Error fetching books:', error)
    }
  }
  
  const searchBooks = () => {
    currentPage.value = 1
  }
  
  const startEditing = (book: any) => {
    book.editing = true
    book.editData = { ...book }
  }
  
  const cancelEditing = (book: any) => {
    book.editing = false
  }
  
  const saveChanges = async (book: any) => {
    try {
      const token = localStorage.getItem('adminToken')
      if (!token) {
        router.push('/admin-login')
        return
      }
      // Thêm log để kiểm tra dữ liệu trước khi gửi
    console.log('Data being sent to update book:', {
      bookId: book.ID,
      updateData: book.editData
    })

    // Log chi tiết từng field
    console.log('Update details:', {
      title: book.editData.title,
      author: book.editData.author,
      price: book.editData.price,
      stock: book.editData.stock,
      // Thêm các trường khác nếu cần
    })
      await axios.put(`${import.meta.env.VITE_API_URL}/api/books/${book.ID}`, book.editData, {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`
        }
      })
  
      // Update the book data
      Object.assign(book, book.editData)
      book.editing = false
    } catch (error) {
      console.error('Error updating book:', error)
      alert('Failed to update book. Please try again.')
    }
  }
  
  const confirmDelete = (book: any) => {
    bookToDelete.value = book
    showDeleteModal.value = true
  }
  
  const deleteBook = async () => {
    try {
      const token = localStorage.getItem('adminToken')
      if (!token) {
        router.push('/admin-login')
        return
      }
  
      await axios.delete(`${import.meta.env.VITE_API_URL}/api/books/${bookToDelete.value.ID}`, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      })
  
      // Remove the book from the list
      books.value = books.value.filter(b => b.id !== bookToDelete.value.id)
      showDeleteModal.value = false
      window.location.reload()

    } catch (error) {
      console.error('Error deleting book:', error)
      alert('Failed to delete book. Please try again.')
    }
  }
  
  const nextPage = () => {
    if (currentPage.value * itemsPerPage < totalBooks.value) {
      currentPage.value++
    }
  }
  
  const prevPage = () => {
    if (currentPage.value > 1) {
      currentPage.value--
    }
  }
  
  const getCorrectImageUrl = (path: string) => {
  if (!path) return '/placeholder.jpg';
  if (path.startsWith('http')) return path; // Nếu đã là URL đầy đủ
  
  const base = import.meta.env.VITE_API_URL || 'http://localhost:8081';
  return `${base}${path.startsWith('/') ? path : `/${path}`}`;
};
  // Lifecycle
  onMounted(() => {
    fetchBooks()
  })
  </script>