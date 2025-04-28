import { defineStore } from 'pinia'

interface Book {
  id: number
  title: string
  author: string
  price: number
  image: string
  category: string
}

interface CartItem extends Book {
  quantity: number
}

export const useCartStore = defineStore('cart', {
  state: () => ({
    items: [] as CartItem[],
  }),
  
  actions: {
    // Initialize cart from localStorage
    initializeCart() {
      const savedCart = localStorage.getItem('cart')
      if (savedCart) {
        this.items = JSON.parse(savedCart)
      }
    },
    
    // Save cart to localStorage
    saveToLocalStorage() {
      localStorage.setItem('cart', JSON.stringify(this.items))
    },
    
    addToCart(book: Book) {
      const existingItem = this.items.find(item => item.id === book.id)
      
      if (existingItem) {
        existingItem.quantity++
      } else {
        this.items.push({ ...book, quantity: 1 })
      }
      
      this.saveToLocalStorage()
    },
    
    removeFromCart(id: number) {
      this.items = this.items.filter(item => item.id !== id)
      this.saveToLocalStorage()
    },
    
    updateQuantity(id: number, newQuantity: number) {
      const item = this.items.find(item => item.id === id)
      if (item) {
        item.quantity = newQuantity > 0 ? newQuantity : 1
        this.saveToLocalStorage()
      }
    },
    
    clearCart() {
      this.items = []
      localStorage.removeItem('cart')
    }
  },
  
  getters: {
    totalItems(): number {
      return this.items.reduce((total, item) => total + item.quantity, 0)
    },
    subtotal(): number {
      return this.items.reduce((sum, item) => sum + (item.price * item.quantity), 0)
    },
    total(): number {
      const shipping = 3.99
      return this.subtotal + shipping
    }
  }
})