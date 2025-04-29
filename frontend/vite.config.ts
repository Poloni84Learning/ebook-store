import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': '/src' // Chỉ cần thế này trong Vite
    }
  },
  server: {
    proxy: {
      '/api':{
        target: 'http://localhost:8081', // Địa chỉ backend
        changeOrigin: true,
      },
    },
    host: '0.0.0.0',  // Cho phép truy cập từ bên ngoài container
    port: 5173,
    strictPort: true,  // Bắt buộc dùng đúng cổng
    watch: {
      usePolling: true  // Cần thiết khi dùng Docker trên Windows/macOS
    }
  },
  })
