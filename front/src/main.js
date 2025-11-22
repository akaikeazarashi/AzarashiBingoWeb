import './scss/app.scss'

// import 'bootstrap/dist/css/bootstrap.min.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import { useSessionStore } from '@/stores/session'
import { useBingoStore } from '@/stores/bingo'
import { Modal } from '@kouts/vue-modal'
import '@kouts/vue-modal/dist/vue-modal.css'
import Toast from "vue-toastification"
import "vue-toastification/dist/index.css"

const app = createApp(App);

app.use(router);

// モーダル登録
app.component('Modal', Modal);

// トースト登録
app.use(Toast, {
  position: "top-right",
  timeout: 3000,
  closeOnClick: true,
  pauseOnHover: true,
})

const pinia = createPinia();
app.use(pinia);

// ストアの認証状態を初期化
const sessionStore = useSessionStore();
sessionStore.initStore();

// ストアのビンゴデータを初期化
const bingoStore = useBingoStore();
bingoStore.initStore();

app.mount('#app');
