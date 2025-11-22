import { fileURLToPath, URL } from 'node:url'

import { path } from 'path'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

export default defineConfig(({ command, mode }) => {
  const isProd = mode === 'production'
  return {
    base: isProd ? "/static/" : "/",
    plugins: [
      vue({
        reactivityTransform: true,
        template: {
            transformAssetUrls: {
                includeAbsolute: false,
            },
        },
      }),
      vueDevTools(),
    ],
    server: {
      host: true,
      port: 5173,
      hmr: {
        protocol: 'ws',
        host: 'localhost',
        port: 5173,
        clientPort: 5173
      },
      proxy: {
        // APIをGoサーバーへ転送
        '/api': 'http://localhost:8080',
      }
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
        vue: 'vue/dist/vue.esm-bundler.js',
      },
    },
  }
})
