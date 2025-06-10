import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'node:path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@app': path.resolve(__dirname, './src/app'),
      '@features': path.resolve(__dirname, './src/features'),
      '@shared': path.resolve(__dirname, './src/shared'),
      '@design-system': path.resolve(__dirname, './src/design-system'),
      '@assets': path.resolve(__dirname, './src/assets'),
      '@hooks': path.resolve(__dirname, './src/shared/hooks'),
      '@types': path.resolve(__dirname, './src/shared/types'),
      '@utils': path.resolve(__dirname, './src/shared/utils'),
      '@services': path.resolve(__dirname, './src/shared/services'),
      '@components': path.resolve(__dirname, './src/shared/components'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL || 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'blueprint': [
            '@blueprintjs/core',
            '@blueprintjs/table',
            '@blueprintjs/select',
            '@blueprintjs/icons',
            '@blueprintjs/popover2'
          ],
          'vendor': ['react', 'react-dom', 'react-router-dom'],
          'graphql': ['@apollo/client', 'graphql'],
        },
      },
    },
  },
})
