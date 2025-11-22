import { defineStore } from 'pinia';

interface State {
  authed: boolean;
}

export const useSessionStore = defineStore('session', {
  state: (): State => {
    return {
      authed: false
    }
  },
  getters: {
    isAuthed(state): boolean {
      return state.authed;
    },
  },
  actions: {
    initStore() {
      this.authed = sessionStorage.getItem('authenticated') === 'true';
    },
    onLogined() {
      sessionStorage.setItem('authenticated', 'true');
      this.authed = true;
    },
    onUnauthenticated() {
      this.logout();
    },
    logout() {
      sessionStorage.removeItem('authenticated');
      this.authed = false;
    },
    getSavedAdminUserId(): string {
      const storageUserId = localStorage.getItem('admin_user_id');
      if (!storageUserId) {
        return null;
      }
      return storageUserId;
    },
    saveAdminUserId(userId: string): void {
      localStorage.setItem('admin_user_id', userId);
    },
    removeAdminUserId(): void {
      localStorage.removeItem('admin_user_id');
    }
  }
});
