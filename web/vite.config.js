import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/': 'http://0.0.0.0:10001',
    }
  },
  build: {
    outDir: '../cmd/cli/ui'
  }
})
