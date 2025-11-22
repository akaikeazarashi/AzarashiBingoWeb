<template>
  <section id="admin-list-container">
    <div class="row justify-content-center">
      <div class="col-md-8">
        <div class="header-button-container d-flex gap-3">
          <router-link to="/admin/new" class="btn btn-primary" :disabled="isSubmitting">新規登録</router-link>
          <button type="button" class="btn btn-primary" :disabled="isSubmitting" @click="onClickImport">インポート</button>
          <input type="file" ref="fileInput" accept="application/json" style="display:none;" @change="handleImportFileChange" />
        </div>
        <table class="table table-striped">
          <thead>
            <tr>
              <th scope="col" class="col-md-1">id</th>
              <th scope="col" class="col-md-2">ビンゴ名</th>
              <th scope="col" class="col-md-4">説明</th>
              <th scope="col" class="col-md-1">サイズ</th>
              <th scope="col" class="col-md-3">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="bingo in bingoList" :key="bingo.id">
              <td class="list-col-center">{{ bingo.id }}</td>
              <td>{{ bingo.name }}</td>
              <td>{{ bingo.description }}</td>
              <td class="list-col-center">{{ bingo.size }}</td>
              <td class="list-col-center">
                <div class="list-button-container d-flex gap-3">
                  <router-link :to="{name: 'admin_edit', params: {bingoId: bingo.id}}" class="btn btn-primary" :disabled="isSubmitting">編集</router-link>
                  <button type="button" class="btn btn-danger" :disabled="isSubmitting" @click="onClickDelete(bingo.id)">削除</button>
                  <a :href="`/api/admin/export/${bingo.id}`" class="btn btn-success" :disabled="isSubmitting">エクスポート</a>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import common from "../common";
import api from '../util/api';
import { useToast } from "vue-toastification";
import type { Bingo } from '../types/bingo';

const bingoList = ref<Bingo[]>([]);
const isSubmitting = ref(false);
const fileInput = ref<HTMLInputElement | null>(null);
const toast = useToast();

// ビンゴ一覧表示
const showList = async () => {
  const response = await api.get('/api/admin/list');
  const data = response.data;

  if (!data.result) {
    throw new Error('リスト取得に失敗しました。');
  }

  bingoList.value = data.bingoList.map((item: Bingo) => ({
    id: item.id,
    name: item.name,
    description: item.description,
    size: item.size,
  }));
}

// ビンゴ削除
const onClickDelete = async (bingoId: number) => {
  if (!confirm("本当に削除しますか？")) {
    return;
  }

  // 処理中はボタンを無効化
  isSubmitting.value = true;

  try {
    const response = await api.delete(`/api/admin/delete/${bingoId}`);
    const data = response.data;

    if (!data.result) {
      toast.error('削除に失敗しました');
      throw new Error('削除に失敗しました');
    }
  } finally {
    isSubmitting.value = false;
  }

  toast.success('削除に成功しました');

  // 少し待機してから作品一覧画面を再描画
  await common.sleep(500)
  await showList()
}

// インポートボタンクリック時
const onClickImport = async () => {
  // ファイル選択ダイアログ表示
  fileInput.value?.click()

}

// インポートファイル選択時
const handleImportFileChange = async (e: Event) => {
  const target = e.target as HTMLInputElement
  if (!target.files || target.files.length === 0) {
    return;
  }

  const file = target.files[0];

  try {
    // ファイルからJSONデータを用意
    const text = await file.text();
    const jsonData = JSON.parse(text)

    // サーバーに送信
    const response = await api.put(`/api/admin/import`, jsonData);
    const data = response.data;

    if (!data.result) {
      toast.error('インポートに失敗しました');
      throw new Error('インポートに失敗しました');
    }
  } finally {
    isSubmitting.value = false;
  }

  toast.success('インポートに成功しました');

  // 少し待機してから作品一覧画面を再描画
  await common.sleep(500)
  await showList()
}

onMounted(showList);
</script>

<style lang="scss">
#admin-list-container {
  th {
    text-align: center;
  }
  .header-button-container {
    margin: 20px 0 10px;
  }
  .list-col-center {
    text-align: center;
  }
  .list-button-container {
    justify-content: center;
  }
  .img-wrapper {
    width: 50%;
    aspect-ratio: 1 / 1;
    overflow: hidden;

    img {
      box-shadow: 0 2px 5px #999;
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }
}

@media screen and (max-width: 992px) {
}
</style>

