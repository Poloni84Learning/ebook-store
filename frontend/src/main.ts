import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import './assets/main.css'
import router from './router'
import { createI18n } from 'vue-i18n'
import messages from './locales'


const i18n = createI18n({
    legacy: false, // ⚠️ BẮT BUỘC khi dùng Composition API
    locale: 'en',
    fallbackLocale: 'en',
    messages
  })
  
  const app = createApp(App)
  app.use(router)
  app.use(i18n)
  app.mount('#app')