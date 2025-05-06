<template>
    <section class="flex md:bg-gray-100 min-h-screen overflow-hidden">
      <!-- Sidebar -->
      <AdminSidebar />
  
      <!-- Main Content -->
      <div class="flex-grow text-gray-800">
        <!-- Header -->
        <AppHeader />
  
        <!-- User Management Content -->
        <main class="p-6 sm:p-10 space-y-6">
          <div class="flex flex-col space-y-6 md:space-y-0 md:flex-row justify-between">
            <div class="mr-6">
              <h1 class="text-4xl font-semibold mb-2">User Management</h1>
              <h2 class="text-gray-600 ml-0.5">Manage staff and customer accounts</h2>
            </div>
            
            <button 
              @click="showCreateStaffModal = true"
              class="inline-flex px-5 py-3 text-white bg-purple-600 hover:bg-purple-700 focus:bg-purple-700 rounded-md mb-3">
              <svg aria-hidden="true" fill="none" viewBox="0 0 24 24" stroke="currentColor" class="flex-shrink-0 h-6 w-6 text-white -ml-1 mr-2">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Create New Staff
            </button>
          </div>
  
          <!-- User Table -->
          <section class="grid md:grid-cols-1 xl:grid-cols-1 gap-6">
            <div class="bg-white shadow rounded-lg overflow-hidden">
              <div class="px-6 py-4 border-b border-gray-100">
                <h3 class="text-lg font-semibold text-gray-900">All Users</h3>
              </div>
              
              <div class="overflow-x-auto">
                <table class="min-w-full divide-y divide-gray-200">
                  <thead class="bg-gray-50">
                    <tr>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Username</th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
                      <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                    </tr>
                  </thead>
                  <tbody class="bg-white divide-y divide-gray-200">
                    <tr v-if="loadingUsers">
                      <td colspan="5" class="px-6 py-4 text-center">
                        <div class="flex justify-center items-center">
                          <div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-purple-500"></div>
                        </div>
                      </td>
                    </tr>
                    <tr v-else-if="filteredUsers.length === 0">
                      <td colspan="5" class="px-6 py-4 text-center text-gray-500">No users found</td>
                    </tr>
                    <tr v-for="user in filteredUsers" :key="user.id">
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ user.id }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ user.username }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ user.email }}</td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        <span 
                          :class="{
                            'bg-green-100 text-green-800': user.role === 'staff',
                            'bg-blue-100 text-blue-800': user.role === 'customer'
                          }"
                          class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                          {{ user.role }}
                        </span>
                      </td>
                      <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        <div class="flex space-x-2">
                          <select 
                            v-model="user.role" 
                            @change="updateUserRole(user)"
                            :disabled="user.role === 'admin'"
                            class="rounded border border-gray-300 py-1 px-2 text-sm focus:outline-none focus:ring-2 focus:ring-purple-600">
                            <option value="staff">Staff</option>
                            <option value="customer">Customer</option>
                          </select>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </section>
  
          <section class="text-right font-semibold text-gray-500">
            <a href="#" class="text-purple-600 hover:underline">EBook Store</a> User Management
          </section>
        </main>
      </div>
  
      <!-- Create Staff Modal -->
      <div v-if="showCreateStaffModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center p-4 z-50">
        <div class="bg-white rounded-lg shadow-xl max-w-md w-full">
          <div class="p-6">
            <div class="flex justify-between items-center mb-4">
              <h3 class="text-lg font-semibold text-gray-900">Create New Staff Account</h3>
              <button @click="showCreateStaffModal = false" class="text-gray-400 hover:text-gray-500">
                <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <form @submit.prevent="createStaffAccount" class="space-y-4">
              <div>
                <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
                <input 
                  v-model="newStaff.username"
                  type="text" 
                  id="username" 
                  required
                  class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-purple-600 focus:border-purple-600">
              </div>
              
              <div>
                <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
                <input 
                  v-model="newStaff.email"
                  type="email" 
                  id="email" 
                  required
                  class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-purple-600 focus:border-purple-600">
              </div>
              
              <div>
                <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
                <input 
                  v-model="newStaff.password"
                  type="password" 
                  id="password" 
                  required
                  minlength="8"
                  class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-purple-600 focus:border-purple-600">
                <p class="mt-1 text-xs text-gray-500">Password must be at least 8 characters</p>
              </div>
              
              <div class="flex justify-end space-x-3 pt-4">
                <button 
                  @click="showCreateStaffModal = false"
                  type="button"
                  class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500">
                  Cancel
                </button>
                <button 
                  type="submit"
                  :disabled="creatingStaff"
                  class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-purple-600 hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 disabled:opacity-50">
                  <span v-if="creatingStaff">Creating...</span>
                  <span v-else>Create Staff</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </section>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted, computed } from 'vue'
  import { useRouter } from 'vue-router'
  import axios from 'axios'
  import AdminSidebar from '@/components/AdminSidebar.vue'
  import AppHeader from '@/components/HeaderAdmin.vue'
  
  const router = useRouter()
  
  // State
  const users = ref<any[]>([])
  const loadingUsers = ref(false)
  const showCreateStaffModal = ref(false)
  const creatingStaff = ref(false)
  const newStaff = ref({
    username: '',
    email: '',
    password: ''
  })
  
  // Computed
  const filteredUsers = computed(() => {
    return users.value.filter(user => user.role !== 'admin')
  })
  
  // Methods
  const fetchUsers = async () => {
    try {
      loadingUsers.value = true
      const token = localStorage.getItem('adminToken')
      if (!token) {
        router.push('/admin-login')
        return
      }
  
      const response = await axios.get(
        `${import.meta.env.VITE_API_URL}/api/admin/users`,
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      )
  
      users.value = response.data.users
    } catch (error) {
      console.error('Error fetching users:', error)
      alert('Failed to fetch users. Please try again.')
    } finally {
      loadingUsers.value = false
    }
  }
  
  const createStaffAccount = async () => {
    try {
      creatingStaff.value = true
      const token = localStorage.getItem('adminToken')
      if (!token) {
        router.push('/admin-login')
        return
      }
  
      await axios.post(
        `${import.meta.env.VITE_API_URL}/api/admin/users`,
        newStaff.value,
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      )
  
      // Reset form and close modal
      newStaff.value = { username: '', email: '', password: '' }
      showCreateStaffModal.value = false
      
      // Refresh user list
      await fetchUsers()
      
      alert('Staff account created successfully!')
    } catch (error) {
      console.error('Error creating staff account:', error)
      alert('Failed to create staff account. Please try again.')
    } finally {
      creatingStaff.value = false
    }
  }
  
  const updateUserRole = async (user: any) => {
    if (user.role === 'admin') return // Don't allow changing admin roles
    
    if (!confirm(`Are you sure you want to change ${user.username}'s role to ${user.role}?`)) {
      // Reset to original role if user cancels
      await fetchUsers()
      return
    }
  
    try {
      const token = localStorage.getItem('adminToken')
      if (!token) {
        router.push('/admin-login')
        return
      }
  
      await axios.put(
        `${import.meta.env.VITE_API_URL}/api/admin/users/${user.id}/role`,
        { role: user.role },
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      )
  
      alert('User role updated successfully!')
    } catch (error) {
      console.error('Error updating user role:', error)
      alert('Failed to update user role. Please try again.')
      // Reset to original role on error
      await fetchUsers()
    }
  }
  
  // Lifecycle hooks
  onMounted(() => {
    fetchUsers()
  })
  </script>
  
  <style scoped>
  /* Custom styles if needed */
  </style>