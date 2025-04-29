<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import Navbar from './components/Navbar.vue'
import Footer from './components/Footer.vue'
import HelperChat from './components/HelperChat.vue'
import { useCartStore } from '@/stores/cart'
import { computed } from 'vue'

const cartStore = useCartStore()
cartStore.initializeCart()
const route = useRoute()


const showHelperChat = computed(() => {
  return !route.meta.hideHelper
})
</script>

<template>
  <div class="min-h-screen bg-blackBG font-sans">
    <Navbar v-if="!route.meta.hideLayout" />
    <main class="w-full px-0">
      <RouterView />
    </main>
    <Footer v-if="!route.meta.hideLayout" />
    <HelperChat v-if="showHelperChat" />
  </div>
</template>

<style>
/* Import fonts */
@import url('https://fonts.googleapis.com/css2?family=Montserrat:ital,wght@0,100..900;1,100..900&family=Nunito+Sans:ital,opsz,wght@0,6..12,200..1000;1,6..12,200..1000&display=swap');

.btn-primary {
  @apply bg-[#FFCE1A] px-12 py-2 rounded-md text-base font-bold hover:bg-[#0D0842] hover:text-white transition-all duration-200 cursor-pointer;
}
</style>