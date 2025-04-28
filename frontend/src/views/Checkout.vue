<script setup lang="ts">
import { onMounted } from 'vue'
import { useCartStore } from '@/stores/cart'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const router = useRouter()
const { t } = useI18n()
const cartStore = useCartStore()
const { items, subtotal, total } = storeToRefs(cartStore)

const shippingCost = 3.99
const currentStep = ref(1)
const isLoading = ref(true)
const form = ref({
  name: '',
  email: '',
  address: '',
  phone: '',
  city: '',
  zip: ''
})

const errors = ref({
  name: '',
  email: '',
  address: '',
  phone: '',
  city: '',
  zip: ''
})

const validateForm = () => {
  let isValid = true
  
  if (!form.value.name) {
    errors.value.name = t('checkout.errors.required')
    isValid = false
  }
  
  if (!form.value.email) {
    errors.value.email = t('checkout.errors.required')
    isValid = false
  } else if (!/^\S+@\S+\.\S+$/.test(form.value.email)) {
    errors.value.email = t('checkout.errors.email')
    isValid = false
  }
  
  return isValid
}

const placeOrder = async () => {
  if (!validateForm()) return;
  
  try {
    isLoading.value = true;
    
    // Prepare order items
    const orderItems = items.value.map(item => ({
      book_id: item.id,
      quantity: item.quantity
    }));
    
    // Determine payment method
    let paymentMethodValue = '';
    switch (paymentMethod.value) {
      case 'credit':
        paymentMethodValue = 'Card';
        break;
      case 'bankTransfer':
        paymentMethodValue = 'BankTransfer';
        break;
      case 'cod':
        paymentMethodValue = 'COD';
        break;
      default:
        paymentMethodValue = 'COD';
    }
    
    // Get auth token
    const token = localStorage.getItem('authToken');
    if (!token) {
      throw new Error('No authentication token found');
    }
    const order = {
        order_items: orderItems,
        payment_method: paymentMethodValue,
    }
    // Send order to API
    try {
    const response = await axios.post(
      `${import.meta.env.VITE_API_URL}/api/orders`, order, { 
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        }
      }
    );
    
    console.log('Order response:', response);
    if (response.data.success) {
      // Clear cart on successful order
      cartStore.clearCart();
      
      // Redirect to order history
      router.push('/order-history');
    } else {
      throw new Error(response.data.message || 'Failed to place order');
    }
} catch (error) {
  console.error('Error in API call:', error);
}
  } catch (error) {
    console.error('Order failed:', error);
    alert('Failed to place order. Please try again.');
  } finally {
    isLoading.value = false;
  }
};
const fetchUserProfile = async () => {
  try {
    const token = localStorage.getItem('authToken')
    if (!token) return
    
    const res = await axios.get(`${import.meta.env.VITE_API_URL}/api/user/profile`, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    
    if (res.data.success) {
      const user = res.data.user
      form.value = {
        name: `${user.first_name} ${user.last_name}`.trim(),
        email: user.email,
        address: user.address,
        phone: user.phone,
        city: '',
        zip: ''
      }
    }
  } catch (error) {
    console.error('Error fetching user profile:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
    const token = localStorage.getItem('authToken')
  if (!token) {
    router.push('/login')
  }
  fetchUserProfile()
})

const paymentMethod = ref('credit') // credit, bankTransfer, cod
const cardDetails = ref({
  number: '',
  name: '',
  expiry: '',
  cvc: ''
})
const saveCard = ref(false)
const paymentErrors = ref({
  number: '',
  expiry: '',
  cvc: ''
})


const handleCardNumberInput = (e: Event) => {
  // Format số thẻ: 4242 4242 4242 4242
  const input = e.target as HTMLInputElement
  let value = input.value.replace(/\s/g, '')
  if (value.length > 16) value = value.slice(0, 16)
  let formatted = ''
  for (let i = 0; i < value.length; i++) {
    if (i > 0 && i % 4 === 0) formatted += ' '
    formatted += value[i]
  }
  cardDetails.value.number = formatted
}

const handleExpiryInput = (e: Event) => {
  // Format ngày hết hạn: MM/YY
  const input = e.target as HTMLInputElement
  let value = input.value.replace(/\D/g, '')
  if (value.length > 4) value = value.slice(0, 4)
  if (value.length > 2) {
    value = `${value.slice(0, 2)}/${value.slice(2)}`
  }
  cardDetails.value.expiry = value
}
</script>

<template>
  <div class="container mx-auto px-4 py-8">
    <div v-if="isLoading" class="text-center py-12">
      <svg class="animate-spin h-8 w-8 text-primary mx-auto" viewBox="0 0 24 24">
        <!-- SVG loading spinner -->
      </svg>
      <p class="mt-2">Loading your information...</p>
    </div>
    <h1 class="text-3xl font-bold text-secondary mb-8">{{ t('checkout.title') }}</h1>
    
    <!-- Progress Steps -->
    <div class="flex justify-between mb-12">
      <div 
        v-for="step in 3" 
        :key="step" 
        class="flex flex-col items-center"
        :class="{ 'text-primary': currentStep >= step, 'text-gray-400': currentStep < step }"
      >
        <div class="w-10 h-10 rounded-full border-2 flex items-center justify-center mb-2"
          :class="{
            'border-primary bg-primary text-white': currentStep >= step,
            'border-gray-300': currentStep < step
          }">
          {{ step }}
        </div>
        <span class="text-sm">{{ t(`checkout.steps.${['shipping', 'payment', 'review'][step-1]}`) }}</span>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Form Column -->
      <div class="lg:col-span-2">
        <!-- Shipping Information -->
        <div v-show="currentStep === 1" class="bg-white text-black rounded-lg shadow p-6 mb-6">
          <h2 class="text-xl font-semibold mb-4">{{ t('checkout.steps.shipping') }}</h2>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                {{ t('checkout.form.name') }}
              </label>
              <input 
                v-model="form.name" 
                type="text" 
                class="w-full px-4 py-2 text-black border rounded-lg"
                :class="{ 'border-red-500': errors.name }"
              >
              <p v-if="errors.name" class="text-red-500 text-sm mt-1">{{ errors.name }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700  mb-1">
                {{ t('checkout.form.email') }}
              </label>
              <input 
                v-model="form.email" 
                type="email" 
                class="w-full px-4 py-2 text-black border rounded-lg"
                :class="{ 'border-red-500': errors.email }"
              >
              <p v-if="errors.email" class="text-red-500 text-sm mt-1">{{ errors.email }}</p>
            </div>
            
        <!-- Address -->
        <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  {{ t('checkout.form.address') }}
                </label>
                <input 
                  v-model="form.address" 
                  type="text" 
                  class="w-full px-4 py-2 border rounded-lg"
                  :class="{ 'border-red-500': errors.address }"
                >
                <p v-if="errors.address" class="text-red-500 text-sm mt-1">{{ errors.address }}</p>
              </div>
              
              <!-- Phone -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Phone Number
                </label>
                <input 
                  v-model="form.phone" 
                  type="tel" 
                  class="w-full px-4 py-2 border rounded-lg"
                  :class="{ 'border-red-500': errors.phone }"
                >
                <p v-if="errors.phone" class="text-red-500 text-sm mt-1">{{ errors.phone }}</p>
              </div>
              
              <!-- City & ZIP (không có trong API) -->
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    City
                  </label>
                  <input 
                    v-model="form.city" 
                    type="text" 
                    class="w-full px-4 py-2 border rounded-lg"
                  >
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    ZIP Code
                  </label>
                  <input 
                    v-model="form.zip" 
                    type="text" 
                    class="w-full px-4 py-2 border rounded-lg"
                  >
                </div>
              </div>
          </div>
        </div>

          <!-- Payment Method Selection -->
  <div v-show="currentStep === 2" class="bg-white rounded-lg shadow p-6">
    <h2 class="text-xl font-semibold mb-4">{{ t('checkout.payment.title') }}</h2>
    
    <div class="space-y-4">
      <!-- Payment Method Options -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <button
          v-for="method in ['credit', 'bankTransfer', 'cod']"
          :key="method"
          @click="paymentMethod = method"
          class="border rounded-lg p-4 text-center hover:border-primary transition-colors"
          :class="{
            'border-primary text-black ring-1 ring-primary': paymentMethod === method,
            'border-gray-300 text-black': paymentMethod !== method
          }"
        >
          {{ t(`checkout.payment.methods.${method}`) }}
        </button>
      </div>

      <!-- Credit Card Form -->
      <div v-if="paymentMethod === 'credit'" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            {{ t('checkout.payment.card.number') }}
          </label>
          <input
            v-model="cardDetails.number"
            @input="handleCardNumberInput"
            type="text"
            placeholder="4242 4242 4242 4242"
            class="w-full px-4 py-2 text-black border rounded-lg"
            :class="{ 'border-red-500': paymentErrors.number }"
          >
          <p v-if="paymentErrors.number" class="text-red-500 text-sm mt-1">
            {{ paymentErrors.number }}
          </p>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            {{ t('checkout.payment.card.name') }}
          </label>
          <input
            v-model="cardDetails.name"
            type="text"
            placeholder="John Doe"
            class="w-full px-4 py-2 text-black border rounded-lg"
          >
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              {{ t('checkout.payment.card.expiry') }}
            </label>
            <input
              v-model="cardDetails.expiry"
              @input="handleExpiryInput"
              type="text"
              placeholder="MM/YY"
              class="w-full px-4 py-2 text-black border rounded-lg"
              :class="{ 'border-red-500': paymentErrors.expiry }"
            >
            <p v-if="paymentErrors.expiry" class="text-red-500 text-sm mt-1">
              {{ paymentErrors.expiry }}
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              {{ t('checkout.payment.card.cvc') }}
            </label>
            <input
              v-model="cardDetails.cvc"
              type="text"
              placeholder="123"
              maxlength="3"
              class="w-full px-4 text-black  py-2 border rounded-lg"
              :class="{ 'border-red-500': paymentErrors.cvc }"
            >
            <p v-if="paymentErrors.cvc" class="text-red-500 text-sm mt-1">
              {{ paymentErrors.cvc }}
            </p>
          </div>
        </div>

        <div class="flex items-center mt-4">
          <input
            v-model="saveCard"
            type="checkbox"
            id="saveCard"
            class="h-4 w-4 text-primary rounded border-gray-300 focus:ring-primary"
          >
          <label for="saveCard" class="ml-2 text-sm text-gray-700">
            {{ t('checkout.payment.saveCard') }}
          </label>
        </div>
      </div>

      <!-- BankTransfer Notice -->
      <div v-if="paymentMethod === 'bankTransfer'" class="bg-blue-50 p-4 rounded-lg">
        <p class="text-blue-800">You will be redirected to Bank Transfer to complete your payment</p>
      </div>

      <!-- COD Notice -->
      <div v-if="paymentMethod === 'cod'" class="bg-gray-50 p-4 rounded-lg">
        <p class="text-gray-800">Pay with cash upon delivery</p>
      </div>
    </div>
  </div>
<!-- Review Order Step -->
<div v-show="currentStep === 3" class="bg-white rounded-lg shadow p-6 mb-6">
  <h2 class="text-xl font-semibold mb-6">{{ t('checkout.steps.review') }}</h2>
  
  <!-- Shipping Information Review -->
  <div class="mb-8">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-medium text-gray-800">{{ t('checkout.review.shippingTitle') }}</h3>
      <button 
        @click="currentStep = 1"
        class="text-sm text-primary hover:text-primary-dark"
      >
        {{ t('common.edit') }}
      </button>
    </div>
    
    <div class="bg-gray-50 text-black rounded-lg p-4">
      <p class="font-medium">{{ form.name }}</p>
      <p>{{ form.address }}</p>
      <p>{{ form.city }}, {{ form.zip }}</p>
      <p class="mt-2">{{ form.phone }}</p>
      <p>{{ form.email }}</p>
    </div>
  </div>

  <!-- Payment Method Review -->
  <div class="mb-8">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-medium text-gray-800">{{ t('checkout.review.paymentTitle') }}</h3>
      <button 
        @click="currentStep = 2"
        class="text-sm text-primary hover:text-primary-dark"
      >
        {{ t('common.edit') }}
      </button>
    </div>
    
    <div class="bg-gray-50 text-black rounded-lg p-4">
      <template v-if="paymentMethod === 'credit'">
        <p class="font-medium">{{ t('checkout.payment.methods.credit') }}</p>
        <p class="text-black mt-2">•••• •••• •••• {{ cardDetails.number.slice(-4) }}</p>
        <p>{{ cardDetails.name }}</p>
        <p>Expires {{ cardDetails.expiry }}</p>
      </template>
      
      <template v-else-if="paymentMethod === 'bankTransfer'">
        <p class="font-medium">{{ t('checkout.payment.methods.bankTransfer') }}</p>
        <p class="mt-2 text-sm">{{ t('checkout.review.bankTransferNotice') }}</p>
      </template>
      
      <template v-else>
        <p class="font-medium">{{ t('checkout.payment.methods.cod') }}</p>
        <p class="mt-2 text-sm">{{ t('checkout.review.codNotice') }}</p>
      </template>
    </div>
  </div>


  <!-- Order Summary -->
  <div class="mt-6 text-black border-t border-gray-200 pt-6">
    <div class="flex justify-between mb-2">
      <span class="text-gray-600">{{ t('checkout.orderSummary.subtotal') }}</span>
      <span>${{ subtotal.toFixed(2) }}</span>
    </div>
    <div class="flex text-black justify-between mb-2">
      <span class="text-gray-600">{{ t('checkout.orderSummary.shipping') }}</span>
      <span>${{ shippingCost.toFixed(2) }}</span>
    </div>
    <div class="flex justify-between font-bold text-black text-lg mt-4 pt-4 border-t border-gray-200">
      <span>{{ t('checkout.orderSummary.total') }}</span>
      <span>${{ total.toFixed(2) }}</span>
    </div>
  </div>
</div>
        <!-- Navigation Buttons -->
        <div class="flex justify-between mt-6">
          <button 
            v-if="currentStep > 1" 
            @click="currentStep--"
            class="px-6 py-2 border text-black border-gray-300 rounded-lg hover: bg-gray-100"
          >
            Back
          </button>
          <button 
            v-if="currentStep < 3" 
            @click="currentStep++"
            class="px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary-dark ml-auto"
          >
            Continue
          </button>
        </div>
      </div>

      <!-- Order Summary Column -->
      <div class="lg:col-span-1">
        <div class="bg-white rounded-lg shadow p-6 sticky top-4">
          <h2 class="text-xl text-black font-semibold mb-4">{{ t('checkout.orderSummary.title') }}</h2>
          
          <div class="space-y-4 mb-6">
            <div v-for="item in items" :key="item.id" class="flex items-center">
              <img :src="item.image" :alt="item.title" class="w-16 h-16 object-cover rounded-lg mr-4">
              <div class="flex-1">
                <h3 class="text-black font-medium">{{ item.title }}</h3>
                <p class="text-sm text-gray-600">{{ item.quantity }} × ${{ item.price.toFixed(2) }}</p>
              </div>
              <span class="text-black font-medium">${{ (item.price * item.quantity).toFixed(2) }}</span>
            </div>
          </div>

          <div class="border-t border-gray-200 pt-4 space-y-2">
            <div class="flex text-black justify-between">
              <span>{{ t('checkout.orderSummary.subtotal') }}</span>
              <span>${{ subtotal.toFixed(2) }}</span>
            </div>
            <div class="flex text-black justify-between">
              <span>{{ t('checkout.orderSummary.shipping') }}</span>
              <span>${{ shippingCost.toFixed(2) }}</span>
            </div>
            <div class="flex text-black justify-between font-bold text-lg mt-4">
              <span>{{ t('checkout.orderSummary.total') }}</span>
              <span>${{ total.toFixed(2) }}</span>
            </div>
          </div>

          <button 
            v-if="currentStep === 3"
            @click="placeOrder"
            class="w-full bg-primary text-white py-3 rounded-lg font-bold hover:bg-primary-dark mt-6"
          >
            {{ t('checkout.orderSummary.checkoutButton') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>