<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()

const isOpen = ref(false)

const messages = ref<Array<{text: string, sender: 'user' | 'bot'}>>([])
const newMessage = ref('')

// Load messages from localStorage
onMounted(() => {
  const savedMessages = localStorage.getItem('chatHistory')
  if (savedMessages) {
    messages.value = JSON.parse(savedMessages)
  }
  const token = localStorage.getItem('authToken')
    if (!token) {
        messages.value = []
  }
  
})

// Save messages to localStorage when they change
watch(messages, (newMessages) => {
  localStorage.setItem('chatHistory', JSON.stringify(newMessages))
}, { deep: true })

const sendMessage = () => {
    const token = localStorage.getItem('authToken')
    if (!token) {
    router.push('/login')
  }
  if (newMessage.value.trim() === '') return
  
  // Add user message
  messages.value.push({
    text: newMessage.value,
    sender: 'user'
  })
  
  // Simulate bot response
  setTimeout(() => {
    messages.value.push({
      text: `I received your message: "${newMessage.value}"`,
      sender: 'bot'
    })
    newMessage.value = ''
  }, 500)
}

const clearChat = () => {
  messages.value = []
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
        <h3 class="font-bold text-lg">Helper Chat</h3>
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
            </div>
          </div>
          <div v-else class="flex justify-end">
            <div class="bg-[#FFCE1A] p-3 rounded-lg max-w-[80%]">
              <p class="text-black">{{ message.text }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Input area -->
      <div class="p-4 border-t">
        <form @submit.prevent="sendMessage" class="flex">
          <input
            v-model="newMessage"
            type="text"
            placeholder="Type your message..."
            class="flex-1 border text-black rounded-l-lg px-4 py-2 focus:outline-none focus:ring-1 focus:ring-[#FFCE1A]"
          />
          <button
            type="submit"
            class="bg-[#FFCE1A] hover:bg-[#0D0842] text-black hover:text-white px-4 py-2 rounded-r-lg"
          >
            Send
          </button>
        </form>
      </div>
    </div>
  </div>
</template>