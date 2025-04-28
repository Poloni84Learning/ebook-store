<template>
    <aside class="hidden sm:flex sm:flex-col w-60">
      <a href="#" class="inline-flex items-center justify-center h-20 w-60 text-2xl bg-purple-600 hover:bg-purple-500 focus:bg-purple-500">
        EBookStore
      </a>
      <div class="flex-grow flex flex-col justify-between text-gray-500 bg-gray-800">
        <nav class="flex flex-col mx-4 my-6 space-y-4">
          <!-- Dashboard -->
          <router-link 
            to="/admin/dashboard" 
            class="inline-flex items-center justify-start py-3 rounded-lg w-full transition-colors duration-200"
            :class="{'text-purple-600 bg-white': $route.path.startsWith('/admin/dashboard'), 'text-gray-600 hover:bg-gray-100': $route.path !== '/admin'}"
          >
            <span class="sr-only">Dashboard</span>
            <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="ml-4 h-6 w-6">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            <span class="ml-2">Dashboard</span>
          </router-link>
  
          <!-- Books -->
          <router-link 
            to="/admin/books" 
            class="inline-flex items-center justify-start py-3 rounded-lg w-full transition-colors duration-200"
            :class="{'text-purple-600 bg-white': $route.path.startsWith('/admin/books'), 'text-gray-600 hover:bg-gray-100': !$route.path.startsWith('/admin/books')}"
          >
            <span class="sr-only">Books</span>
            <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="ml-4 h-6 w-6">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
            <span class="ml-2">Books</span>
          </router-link>
  
          <!-- Categories -->
          <router-link 
            to="/admin/categories" 
            class="inline-flex items-center justify-start py-3 rounded-lg w-full transition-colors duration-200"
            :class="{'text-purple-600 bg-white': $route.path.startsWith('/admin/categories'), 'text-gray-600 hover:bg-gray-100': !$route.path.startsWith('/admin/categories')}"
          >
            <span class="sr-only">Categories</span>
            <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="ml-4 h-6 w-6">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
            </svg>
            <span class="ml-2">Categories</span>
          </router-link>
  
          <!-- Orders -->
          <router-link 
            to="/admin/orders" 
            class="inline-flex items-center justify-start py-3 rounded-lg w-full transition-colors duration-200"
            :class="{'text-purple-600 bg-white': $route.path.startsWith('/admin/orders'), 'text-gray-600 hover:bg-gray-100': !$route.path.startsWith('/admin/orders')}"
          >
            <span class="sr-only">Orders</span>
            <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="ml-4 h-6 w-6">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
            </svg>
            <span class="ml-2">Orders</span>
          </router-link>
  
          <!-- Reports -->
          <router-link 
            to="/admin/reports" 
            class="inline-flex items-center justify-start py-3 rounded-lg w-full transition-colors duration-200"
            :class="{'text-purple-600 bg-white': $route.path.startsWith('/admin/reports'), 'text-gray-600 hover:bg-gray-100': !$route.path.startsWith('/admin/reports')}"
          >
            <span class="sr-only">Reports</span>
            <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="ml-4 h-6 w-6">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span class="ml-2">Reports</span>
          </router-link>
        </nav>
        <div class="inline-flex items-center justify-center h-20 w-20 border-t border-gray-700">
          <button class="p-3 hover:text-gray-400 hover:bg-gray-700 focus:text-gray-400 focus:bg-gray-700 rounded-lg">
            <span class="sr-only">Settings</span>
            <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="h-6 w-6">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </button>
        </div>
      </div>
    </aside>
  </template>
  
  <script setup lang="ts">
  import { useRoute } from 'vue-router'
  
  const $route = useRoute()
  </script>