import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 8090,
    host: true,
    proxy: {
      '/client': {
        target: 'http://localhost:9090',
        changeOrigin: true,
      },
      '/user': {
        target: 'http://localhost:9090',
        changeOrigin: true,
      },
    },
  },
});
