<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'

interface Book {
  ID: number
  title: string
  author: string
  cover_image: string
  // Các trường khác nếu cần
}

interface ChatMessage {
  text: string
  sender: 'user' | 'bot'
  books?: Book[]
  
}

const router = useRouter()
const isOpen = ref(false)
const messages = ref<ChatMessage[]>([])
const newMessage = ref('')
const isLoading = ref(false)
const apiUrl = import.meta.env.VITE_API_URL
const messagesEndRef = ref<HTMLDivElement | null>(null)

const scrollToBottom = () => {
  nextTick(() => {
    messagesEndRef.value?.scrollIntoView({ behavior: 'smooth' })
  })
}

// Load messages and send welcome message
onMounted(() => {
  const savedMessages = localStorage.getItem('chatHistory')
  if (savedMessages) {
    messages.value = JSON.parse(savedMessages)
    
  } else {
    messages.value.push({
      text: 'Enter the keyword you want to search in depth',
      sender: 'bot'
    })
  }
  
  
  scrollToBottom()
})

watch(messages, (newMessages) => {
  localStorage.setItem('chatHistory', JSON.stringify(newMessages))
  scrollToBottom()

}, { deep: true })
const getCorrectImageUrl = (path: string) => {
  if (!path) return '/placeholder.jpg';
  if (path.startsWith('http')) return path; // Nếu đã là URL đầy đủ
  
  const base = import.meta.env.VITE_API_URL || 'http://localhost:8081';
  return `${base}${path.startsWith('/') ? path : `/${path}`}`;
};
const searchBooks = async (query: string): Promise<Book[]> => {
  isLoading.value = true
  try {
    const response = await fetch(`${apiUrl}/api/books/search-helper?q=${encodeURIComponent(query)}`)
    if (!response.ok) throw new Error('Search failed')
    const { data } = await response.json()
    return data.slice(0, 4) // Lấy tối đa 4 quyển sách
  } catch (error) {
    console.error('Search error:', error)
    return []
  } finally {
    isLoading.value = false
  }
}

const sendMessage = async () => {
  const token = localStorage.getItem('authToken')
  if (!token) {
    router.push('/login')
    return
  }
  
  if (newMessage.value.trim() === '') return
  
  // Thêm tin nhắn người dùng
  const userMessage: ChatMessage = {
    text: newMessage.value,
    sender: 'user'
  }
  messages.value.push(userMessage)
  
  // Tìm kiếm sách
  const books = await searchBooks(newMessage.value)
  newMessage.value = ''
  
  // Thêm phản hồi từ bot
  if (books.length > 0) {
    messages.value.push({
      text: `Found ${books.length} books related to "${userMessage.text}":`,
      sender: 'bot',
      books: books
    })
  } else {
    messages.value.push({
      text: `No books found for "${userMessage.text}". Try another keyword.`,
      sender: 'bot'
    })
  }
}

const clearChat = () => {
  messages.value = [{
    text: 'Enter the keyword you want to search in depth',
    sender: 'bot'
  }]
}

const openBookInNewTab = (bookId: number) => {
  window.open(`/books/${bookId}`, '_blank')
}
</script>

<template>
  <div class="fixed bottom-8 right-8 z-50">
    <!-- Chat button -->
    <button
      v-if="!isOpen"
      @click="isOpen = true"
      class="bg-[#FFCE1A] hover:bg-[#0D0842] text-black hover:text-white rounded-full w-14 h-14 flex items-center justify-center shadow-lg transition-all duration-200"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
      </svg>
    </button>

    <!-- Chat modal -->
    <div
      v-if="isOpen"
      class="bg-white rounded-lg shadow-xl w-80 h-[500px] flex flex-col"
    >
      <!-- Header -->
      <div class="bg-[#FFCE1A] p-4 rounded-t-lg flex justify-between items-center">
        <h3 class="font-bold text-lg">Book Search Helper</h3>
        <div class="flex space-x-2">
          <button @click="clearChat" class="text-black hover:text-gray-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
          <button @click="isOpen = false" class="text-black hover:text-gray-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Messages -->
      <div class="flex-1 p-4 overflow-y-auto">
        <div v-for="(message, index) in messages" :key="index" class="mb-4">
          <div v-if="message.sender === 'bot'" class="flex items-start">
            <div class="bg-gray-100 p-3 rounded-lg max-w-[80%]">
              <p class="text-black">{{ message.text }}</p>
              <div v-if="message.books" class="mt-2 space-y-2">
                <div 
                  v-for="book in message.books" 
                  :key="book.ID"
                  @click="openBookInNewTab(book.ID)"
                  class="flex items-start p-2 hover:bg-gray-200 rounded cursor-pointer"
                >
                  <img 
                    :src="getCorrectImageUrl(book.cover_image) || '/default-book-cover.jpg'" 
                    class="w-10 h-14 object-cover rounded mr-2"
                  >
                  <div>
                    <p class="text-sm text-black font-medium">{{ book.title }}</p>
                    <p class="text-xs text-gray-600">{{ book.author }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="flex justify-end">
            <div class="bg-[#FFCE1A] p-3 rounded-lg max-w-[80%]">
              <p class="text-black">{{ message.text }}</p>
            </div>
          </div>
        </div>
        <div ref="messagesEndRef"></div>
        <div v-if="isLoading" class="flex justify-center py-2">
          <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-[#FFCE1A]"></div>
        </div>
      </div>

      <!-- Input area -->
      <div class="p-4 border-t">
        <form @submit.prevent="sendMessage" class="flex">
          <input
            v-model="newMessage"
            type="text"
            placeholder="Type your keyword..."
            class="flex-1 border text-black rounded-l-lg px-4 py-2 focus:outline-none focus:ring-1 focus:ring-[#FFCE1A]"
          />
          <button
            type="submit"
            :disabled="isLoading"
            class="bg-[#FFCE1A] hover:bg-[#0D0842] text-black hover:text-white px-4 py-2 rounded-r-lg disabled:opacity-50"
          >
            <svg v-if="isLoading" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 animate-spin" viewBox="0 0 24 24">
              <path fill="currentColor" d="M12,4V2A10,10 0 0,0 2,12H4A8,8 0 0,1 12,4Z" />
            </svg>
            <span v-else>Search</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>