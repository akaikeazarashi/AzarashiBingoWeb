<template>
  <nav class="navbar navbar-expand-lg navbar-light sticky-top">
    <div class="container-fluid">
      <router-link :to="{ path: '/', hash: '#app' }" class="navbar-brand" tabindex="-1">{{ appName }}</router-link>

      <div id="navbar-content" class="collapse navbar-collapse">
        <ul class="navbar-nav ms-auto me-3">
          <transition name="fade">
            <button type="button" class="btn btn-primary" v-if="isAuthed" @click="onClickLogout">ログアウト</button>
          </transition>
        </ul>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import router from "@/router";
import common from "@/common";
import { computed } from "vue";
import { useSessionStore } from "@/stores/session";

const appName = common.getAppName();

const isAuthed = computed(() => {
  return useSessionStore().isAuthed;
})

// ログアウトクリック時
const onClickLogout = () => {
  const sessionStore = useSessionStore();
  sessionStore.logout();

  router.push('/admin');
}
</script>

<style lang="scss">
.navbar {
  background-color: #c3e2ff;
  border-bottom: 3px dashed #9bcfff;
}
</style>
