<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { RouterView } from 'vue-router'

const { t, locale } = useI18n()
</script>

<template>
  <div class="min-h-screen bg-blackBG font-sans">
    <nav class="bg-secondary text-white p-4 shadow-md">
      <div class="container mx-auto flex justify-between items-center">
        <router-link to="/" class="text-2xl font-bold">BookStore</router-link>
        <div class="hidden md:flex space-x-6">
          <router-link to="/" class="hover:text-primary">{{ t('home') }}</router-link>
          <router-link to="/books" class="hover:text-primary">{{ t('books') }}</router-link>
          <router-link to="/categories" class="hover:text-primary">{{ t('categories') }}</router-link>
          <router-link to="/admin" class="hover:text-primary">{{ t('admin') }}</router-link>
        </div>
        <div class="flex items-center space-x-4">
          <button class="p-2 rounded-full hover:bg-gray-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </button>
          <button class="p-2 rounded-full hover:bg-gray-700">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
            </svg>
          </button>
          <router-link to="/login" class="bg-primary text-secondary px-4 py-2 rounded-md font-bold">{{ t('login') }}</router-link>
          <select v-model="locale" class="border px-2 py-1 rounded bg-secondary text-white">
            <option value="en">English</option>
            <option value="ru">Русский</option>
          </select>
        </div>
      </div>
    </nav>

    <main class="container mx-auto">
      <RouterView />
    </main>

    <footer class="bg-secondary text-white py-8 mt-12">
      <div class="container mx-auto px-4">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div>
            <h3 class="text-xl font-bold mb-4">BookStore</h3>
            <p>{{ t('footer.description') }}</p>
          </div>
          <div>
            <h4 class="font-bold mb-4">{{ t('footer.quickLinks') }}</h4>
            <ul class="space-y-2">
              <li><router-link to="/" class="hover:text-primary">{{ t('home') }}</router-link></li>
              <li><router-link to="/books" class="hover:text-primary">{{ t('books') }}</router-link></li>
              <li><router-link to="/categories" class="hover:text-primary">{{ t('categories') }}</router-link></li>
              <li><router-link to="/about" class="hover:text-primary">{{ t('about') }}</router-link></li>
            </ul>
          </div>
          <div>
            <h4 class="font-bold mb-4">{{ t('footer.customerService') }}</h4>
            <ul class="space-y-2">
              <li><router-link to="/contact" class="hover:text-primary">{{ t('contact') }}</router-link></li>
              <li><router-link to="/faq" class="hover:text-primary">{{ t('faq') }}</router-link></li>
              <li><router-link to="/shipping" class="hover:text-primary">{{ t('shipping') }}</router-link></li>
              <li><router-link to="/returns" class="hover:text-primary">{{ t('returns') }}</router-link></li>
            </ul>
          </div>
          <div>
            <h4 class="font-bold mb-4">{{ t('footer.newsletter') }}</h4>
            <p class="mb-2">{{ t('footer.subscribeText') }}</p>
            <div class="flex">
              <input type="email" :placeholder="t('footer.emailPlaceholder')" class="px-4 py-2 rounded-l-md w-full text-gray-800">
              <button class="bg-primary text-secondary px-4 py-2 rounded-r-md font-bold">{{ t('footer.subscribe') }}</button>
            </div>
          </div>
        </div>
        <div class="border-t border-gray-700 mt-8 pt-6 text-center">
          <p>&copy; {{ new Date().getFullYear() }} BookStore. {{ t('footer.rights') }}</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<style>
/* Import fonts */
@import url('https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100..900;1,100..900&family=Nunito+Sans:ital,opsz,wght@0,6..12,200..1000;1,6..12,200..1000&display=swap');

.btn-primary {
  @apply bg-[#FFCE1A] px-12 py-2 rounded-md text-base font-bold hover:bg-[#0D0842] hover:text-white transition-all duration-200 cursor-pointer;
}
</style>