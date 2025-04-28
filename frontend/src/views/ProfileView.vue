<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import DefaultAvatar from '@/assets/img/images_placeholder.png'
import axios from 'axios'

const apiUrl = import.meta.env.VITE_API_URL
const router = useRouter()
const user = ref<any>(null)
const isEditing = ref(false)
const isLoading = ref(true)

const editableFields = reactive({
  email: '',
  first_name: '',
  last_name: '',
  avatar_url: '',
  address: '',
  phone: '',
})

const fetchUserInfo = async () => {
  try {
    const token = localStorage.getItem('authToken')
    if (!token) {
      router.push('/login')
      return
    }

    const res = await axios.get(`${apiUrl}/api/user/profile`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (res.data.success) {
      user.value = res.data.user

      // Khi load xong thì gán dữ liệu vào editableFields
      Object.assign(editableFields, {
        email: user.value.email || '',
        first_name: user.value.first_name || '',
        last_name: user.value.last_name || '',
        avatar_url: user.value.avatar_url || '',
        address: user.value.address || '',
        phone: user.value.phone || '',
      })
    } else {
      router.push('/login')
    }
  } catch (error) {
    console.error('Error fetching profile:', error)
    router.push('/login')
  } finally {
    isLoading.value = false
  }
}

const startEditing = () => {
  isEditing.value = true
}

const cancelEditing = () => {
  isEditing.value = false
  // Reset lại editableFields về dữ liệu gốc
  Object.assign(editableFields, {
    email: user.value.email || '',
    first_name: user.value.first_name || '',
    last_name: user.value.last_name || '',
    avatar_url: user.value.avatar_url || '',
    address: user.value.address || '',
    phone: user.value.phone || '',
  })
}

const saveProfile = async () => {
  try {
    const token = localStorage.getItem('authToken')
    if (!token) return

    // Chỉ gửi đúng những trường API yêu cầu
    const payload = {
      first_name: editableFields.first_name,
      last_name: editableFields.last_name,
      email: editableFields.email,
      address: editableFields.address,
      phone: editableFields.phone,
    }

    const res = await axios.put(`${apiUrl}/api/user/profile`, payload, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })

    if (res.data.success) {
      // Cập nhật user sau khi save thành công
      user.value = { ...user.value, ...payload }
      isEditing.value = false
      alert('Profile updated successfully!')
    } else {
      alert('Update failed')
    }
  } catch (error) {
    console.error('Error updating profile:', error)
    alert('Error updating profile')
  }
}


onMounted(() => {
  fetchUserInfo()
})
</script>

<template>
  <div class="min-h-screen bg-gray-100 flex items-center justify-center p-6">
    <div class="bg-white shadow-lg rounded-lg p-8 w-full max-w-2xl">
      <div v-if="isLoading" class="text-center text-gray-500">Loading...</div>

      <div v-else-if="user" class="flex flex-col items-center space-y-6">
        <!-- Avatar -->
        <img
          :src="isEditing ? editableFields.avatar_url || DefaultAvatar : (user.avatar_url || DefaultAvatar)"
          alt="User Avatar"
          class="w-32 h-32 rounded-full border-4 border-primary object-cover"
        />

        <!-- Các trường thông tin -->
        <div class="text-center space-y-2">
            <h2 class="text-2xl text-black font-bold">{{ user.username }}</h2>
           
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 w-full mt-6">
            <!-- Email -->
            <div>
                <label class="block text-sm text-gray-600">Email</label>
            <template v-if="isEditing">
                <input
                v-model="editableFields.email"
                type="email"
                class="border text-black px-3 py-2 rounded w-full"
                :placeholder="user.email || 'Update your email'"
                />
            </template>
            <template v-else>
              <div class="font-semibold text-black">{{ user.email || 'N/A' }}</div>
            </template>
        </div>
        
        <!-- Phone -->
        <div>
          <label class="block text-sm text-gray-600">Phone</label>
          <template v-if="isEditing">
            <input
              v-model="editableFields.phone"
              type="text"
              class="border text-black px-3 py-2 rounded w-full"
              :placeholder="user.phone || 'Update your phone'"
            />
          </template>
          <template v-else>
            <div class="font-semibold text-black">{{ user.phone || 'N/A' }}</div>
          </template>
        </div>

        <!-- First Name -->
        <div>
            <label class="block text-sm text-gray-600">First Name</label>
            <template v-if="isEditing">
                <input
                v-model="editableFields.first_name"
                type="text"
                class="border text-black px-3 py-2 rounded w-full"
                :placeholder="user.first_name || 'Update your first name'"
              />
            </template>
            <template v-else>
              <div class="font-semibold text-black">{{ user.first_name || 'N/A' }}</div>
            </template>
        </div>

          <!-- Last Name -->
          <div>
            <label class="block text-sm text-gray-600">Last Name</label>
            <template v-if="isEditing">
              <input
                v-model="editableFields.last_name"
                type="text"
                class="border text-black px-3 py-2 rounded w-full"
                :placeholder="user.last_name || 'Update your last name'"
              />
            </template>
            <template v-else>
              <div class="font-semibold text-black">{{ user.last_name || 'N/A' }}</div>
            </template>
          </div>

          <!-- Address -->
        <div class="md:col-span-2">
            <label class="block text-sm text-gray-600">Address</label>
                <template v-if="isEditing">
                    <textarea
                    v-model="editableFields.address"
                    class="border text-black px-3 py-2 rounded w-full h-32 resize-none"
                    :placeholder="user.address || 'Update your address'"
                    ></textarea>
                </template>
            <template v-else>
                <div class="font-semibold text-black">{{ user.address || 'N/A' }}</div>
            </template>
        </div>

       
        </div>

        <!-- Buttons -->
        <div class="flex space-x-4 mt-6">
          <template v-if="!isEditing">
            <button
              class="px-6 py-2 bg-primary text-white rounded-md hover:bg-primary-dark transition"
              @click="startEditing"
            >
              Edit Profile
            </button>
          </template>

          <template v-else>
            <button
              class="px-6 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 transition"
              @click="saveProfile"
            >
              Save
            </button>
            <button
              class="px-6 py-2 bg-gray-400 text-white rounded-md hover:bg-gray-500 transition"
              @click="cancelEditing"
            >
              Cancel
            </button>
          </template>
        </div>
      </div>

      <div v-else class="text-center text-red-500">
        Failed to load user profile.
      </div>
    </div>
  </div>
</template>

<style scoped>
.bg-primary {
  background-color: #3B82F6; /* Tailwind màu blue-500 */
}
.bg-primary-dark {
  background-color: #2563EB; /* Tailwind màu blue-600 */
}
.border-primary {
  border-color: #3B82F6;
}
</style>
