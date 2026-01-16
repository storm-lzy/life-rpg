import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    port: 5173,
    host: true,
    allowedHosts: ['frp-sea.com', '.frp-sea.com','frp-use.com','.frp-use.com'],
    proxy: {
      '/api': {
        target: 'http://frp-use.com:42714',
        changeOrigin: true,
      }
    }
  }
})
