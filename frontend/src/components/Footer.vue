<script setup lang="ts">
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
</script>

<template>
  <footer class="bg-secondary text-white py-8 mt-12">
    <div class="container mx-auto px-4">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <div>
          <h3 class="text-xl font-bold mb-4">EBookStore</h3>
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
</template>
