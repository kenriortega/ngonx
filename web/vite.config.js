import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    // port: 3000,
    // proxy: {
    //   '/api': {
    //     target: 'http://0.0.0.0:10001',
    //     changeOrigin: true,
    //     ws: true,
    //     rewrite: (path) => path.replace(/^\/api/, '')
    //   },
    // }
  },
  build: {
    outDir: '../cmd/cli/ui'
  }
})
