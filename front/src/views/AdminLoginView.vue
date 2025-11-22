<template>
  <section id="admin-login-container">
    <div class="row justify-content-center">
      <div class="col-md-3">
        <div class="row">
          <div class="col-md">
            <label for="login-userid" class="form-label">ユーザーid</label>
            <input type="text" id="login-userid" class="form-control" v-model="userId" @keydown.enter="onClickLogin">
            <div class="form-check mt-2">
              <input type="checkbox" id="checkbox-save-user-id" class="form-check-input" v-model="isSaveUserId">
              <label class="form-check-label" for="checkbox-save-user-id">ユーザーidを保存する</label>
            </div>
          </div>
        </div>
        <div class="row">
          <div class="col-md">
            <label for="login-password" class="form-label">パスワード</label>
            <input type="password" id="login-password" class="form-control" v-model="password" @keydown.enter="onClickLogin">
          </div>
        </div>
        <div class="row">
          <div class="col-md">
            <p v-text="errorMsg" class="text-danger"></p>
            <button type="button" class="btn btn-primary" @click="onClickLogin" :disabled="isSubmitting">ログイン</button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import api from '../util/api';
import common from '../common';
import router from "../router";
import { useSessionStore } from '../stores/session';

const userId = ref('');
const password = ref('');
const errorMsg = ref('');
const isSaveUserId = ref(false);
const isSubmitting = ref(false);

const sessionStorage = useSessionStore();

// ログインボタン押下時
const onClickLogin = async () => {
  errorMsg.value = '';

  // 処理中はログインボタンを無効化
  isSubmitting.value = true;

  try {
    const response = await api.post('/api/admin/login', {
      'user_id': userId.value,
      'password': password.value
    });

    const responseData = response.data;

    if (!responseData.result) {
      throw new Error('ログインに失敗しました。しばらく経ってから再度試してください。');
    }

    // セッション管理用ストアを認証状態にする
    sessionStorage.onLogined();

    if (isSaveUserId.value) {
      sessionStorage.saveAdminUserId(userId.value);
    } else {
      sessionStorage.removeAdminUserId();
    }

    // 作品一覧画面にリダイレクト
    await common.sleep(100);
    router.push('/admin/list');
  } catch(e) {
    errorMsg.value = 'エラーが発生しました。再度ログインしてください。'
    throw e
  } finally {
    isSubmitting.value = false;
  }
}

onMounted(() => {
  // 既にログインしている場合は一覧ページへリダイレクト
  if (sessionStorage.isAuthed) {
    router.push('/admin/list');
  }

  const storageUserId = sessionStorage.getSavedAdminUserId();
  if (storageUserId) {
    userId.value = storageUserId;
    isSaveUserId.value = true;
  } else {
    isSaveUserId.value = false;
  }
});
</script>

<style scoped lang="scss">
#admin-login-container {
  margin-top: 20px;

  .row {
    margin-bottom: 20px;
  }
}

@media screen and (max-width: 992px) {
}
</style>

