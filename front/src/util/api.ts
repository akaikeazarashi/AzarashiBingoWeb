import axios from 'axios'
import router from "../router";
import { useSessionStore } from '../stores/session';

const client = axios.create();

client.interceptors.response.use(
  response => response,
  error => {
    const status = error?.response?.status

    // 401エラーの場合は認証状態を解除
    if (status === 401) {
      useSessionStore().onUnauthenticated();
      router.push('/admin/login');
    }
    return Promise.reject(error)
  }
)

export default {
  async get(url: string, credentials: boolean = true) {
    return client.get(url, {
      withCredentials: credentials
    });
  },

  async post(url: string, data: Object, credentials: boolean = true) {
    return client.post(url, data, {
      withCredentials: credentials
    });
  },

  async postWithHeader(url: string, data: Object, header: Object) {
    return client.post(url, data, header);
  },

  async put(url: string, data: Object, credentials: boolean = true) {
    return client.put(url, data, {
      withCredentials: credentials
    });
  },

  async delete(url: string, credentials: boolean = true) {
    return client.delete(url, {
      withCredentials: credentials
    });
  }
}
