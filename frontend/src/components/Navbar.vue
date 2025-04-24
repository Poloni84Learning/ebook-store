<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t, locale } = useI18n()
</script>

<template>
  <nav class="bg-secondary text-white p-4 shadow-md">
    <div class="container mx-auto flex justify-between items-center">
      <router-link to="/" class="text-2xl font-bold">EBookStore</router-link>
      <div class="hidden md:flex space-x-6">
        <router-link to="/" class="hover:text-primary">{{ t('home.title') }}</router-link>
        <router-link to="/books" class="hover:text-primary">{{ t('books') }}</router-link>
        <router-link to="/categories" class="hover:text-primary">{{ t('categories') }}</router-link>
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
        <select v-model="locale" class="border px-2 py-1 rounded bg-white text-black">
          <option value="en">English</option>
          <option value="ru">Русский</option>
        </select>
      </div>
    </div>
  </nav>
</template>
