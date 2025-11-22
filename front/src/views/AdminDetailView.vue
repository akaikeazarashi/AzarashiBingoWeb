<template>
  <div>
    <section id="admin-detail-container">
      <h2 class="text-center">ビンゴ登録・更新</h2>
      <div id="admin-form-container" class="row justify-content-center">
        <div class="col-md-6 form-area">
          <div class="row">
            <div class="col-md-3 align-content-center">
              <label for="bingo-name" class="form-label">ビンゴ名</label>
            </div>
            <div class="col-md">
              <input type="text" id="bingo-name" name="bingo_name" class="form-control" v-model="bingoName">
            </div>
          </div>
          <div class="row">
            <div class="col-md-3 align-content-center">
              <label for="description" class="form-label">説明</label>
            </div>
            <div class="col-md">
              <textarea type="text" id="description" name="description" class="form-control" v-model="description"></textarea>
            </div>
          </div>
          <div class="row">
            <div class="col-md-3 align-content-center">
              <label for="size" class="form-label">サイズ</label>
            </div>
            <div class="col-md">
              <select id="item-category" class="form-select" v-model="bingoSize" :disabled="isEdit()" @change="changeSize()">
                <option v-for="val of Array.from(sizeList)" :key="val" :value="val">{{ val }}</option>
              </select>
              <p v-if="isEdit()">編集の場合はサイズを変更できません</p>
            </div>
          </div>
          <div class="row">
            <div class="col">
              <BingoPanel :size="bingoSize" :items="items" v-on:selectedItem="selectedBingoItem"></BingoPanel>
            </div>
          </div>
          <div class="row">
            <div class="col" style="text-align: center;">
              <button type="button" class="btn btn-primary" @click="shuffleItems">シャッフル</button>
            </div>
          </div>
        </div>
        <div class="col-md-2 offset-md-1 form-area" style="height: max-content;">
          <div class="row">
            <button type="button" class="btn btn-primary" :disabled="isSubmitting" @click="submitForm()">登録</button>
          </div>
          <div class="row">
            <router-link to="/admin/list" type="button" class="btn btn-danger" :disabled="isSubmitting">もどる</router-link>
          </div>
        </div>
      </div>
    </section>

    <Modal v-model="isShowModal" title="ビンゴアイテム名設定" :enableClose="false">
      <label for="bingo-name" class="form-label">アイテム名</label>
      <input type="text" name="item_name" class="form-control" v-model="modalItemName">
      <button type="button" class="btn btn-primary mt-3" @click="submitItemName(modalItemName)">確定</button>
      <button type="button" class="btn btn-danger mt-3 ms-3" @click="cancelItemName()">キャンセル</button>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import common from '../common';
import api from '../util/api';
import router from "../router";
import BingoPanel from '../components/BingoPanel.vue';
import { ref, onMounted } from 'vue';
import { useToast } from "vue-toastification";
import type { BingoItem } from '../types/bingo';

const sizeList = common.getSizeList();

const props = defineProps<{
  bingoId?: number | string,
}>();

const isSubmitting = ref(false)
const bingoId = ref(0)
const bingoName = ref('')
const description = ref('')
const bingoSize = ref(3);
const items = ref<BingoItem[]>([]);
const isShowModal = ref(false);
const modalItem = ref<BingoItem|null>(null);
const modalItemName = ref('');
const toast = useToast();

// ビンゴアイテム初期化
const initBingoItems = () => {
  items.value = [];
  let index = 0;
  for (let row = 0; row < bingoSize.value; row++) {
    for (let col = 0; col < bingoSize.value; col++) {
      items.value.push({ id: 0, name: '', orderIndex: index });
      console.log(index);
      index++;
    }
  }
};

// ビンゴサイズ変更時の処理
const changeSize = () => {
  // サイズに合わせてアイテム配列を初期化
  initBingoItems();
};

// ビンゴ表でアイテム選択時に呼び出されるメソッド
const selectedBingoItem = (item: BingoItem) => {
  // 名前変更用のモーダル表示
  modalItemName.value = item.name;
  modalItem.value = item;
  isShowModal.value = true;
}

// モーダルでアイテム名確定時
const submitItemName = (name: string) => {
  if (!modalItem.value) {
    return;
  }
  modalItem.value.name = name;
  modalItem.value = null;
  isShowModal.value = false;
}

// モーダルキャンセル時
const cancelItemName = () => {
  if (!modalItem.value) {
    return;
  }
  modalItem.value = null;
  isShowModal.value = false;
}

// フォーム送信時
const submitForm = async () => {
  if (bingoName.value === '') {
    alert('ビンゴ名が入力されていません');
    return;
  }

  // ビンゴアイテム数の妥当性チェック
  if (items.value.length != (bingoSize.value * bingoSize.value)) {
    alert('ビンゴアイテム数とサイズが不一致');
    return;
  }

  // 各マスの入力チェック
  let inputError = false;
  items.value.forEach((item) => {
    if (!item.name || item.name === '') {
      toast.error(`"${item.orderIndex} 番目のアイテムの名前が指定されていません"`);
      inputError = true;
    }
  });
  if (inputError) {
    return;
  }

  const payload = {
    bingoId: bingoId.value,
    name: bingoName.value,
    description: description.value,
    size: bingoSize.value,
    items: items.value,
  }

  isSubmitting.value = true

  try {
   const response = await api.put('/api/admin/put', payload);
   const responseData = response.data;

    if (!responseData.result) {
      throw new Error('サーバー通信エラー');
    }

    // 少し待機
    await common.sleep(1000);

    // ビンゴ一覧画面にリダイレクト
    router.push('/admin/list');
  } catch(e) {
    alert('エラーが発生しました。' + e);
    throw e;
  } finally {
    isSubmitting.value = false;
  }
};

// サーバーからビンゴデータを取得する
const fetchItem = async () => {
  const paramBingoId = Number(props.bingoId ?? 0);

  // idが渡されている場合はサーバーからデータ取得
  if (paramBingoId && paramBingoId !== 0) {
    const response = await api.get(`/api/admin/detail/${paramBingoId}`);
    const data = response.data;

    if (!response.status) {
      throw new Error('サーバー通信エラー');
    }

    bingoId.value = data.item.id;
    bingoName.value = data.item.name;
    description.value = data.item.description;
    bingoSize.value = data.item.size;
    items.value = data.item.items;
  } else {
    // 新規登録の場合
    bingoId.value = 0;
    bingoName.value = '';
    description.value = '';
    bingoSize.value = 3;

    initBingoItems();
  }
}

// 編集かどうか
const isEdit = (): boolean => {
  const paramBingoId = Number(props.bingoId ?? 0);
  return paramBingoId !== 0
}

// ビンゴアイテムをシャッフル
const shuffleItems = (): void => {
  for (let i = items.value.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [items.value[i], items.value[j]] = [items.value[j], items.value[i]];
  }

  items.value.forEach((item: BingoItem, index: number) => {
    item.orderIndex = index;
  });
}

onMounted(fetchItem)
</script>

<style lang="scss" scoped>
#admin-detail-container {
  margin-top: 30px;
}

.form-area {
  padding: 20px;
  border: solid 4px #dfcd93;
  border-radius: 10px;
  background-color: #fffdf5;
}

#admin-form-container {
  margin-top: 30px;

  .row {
    margin-bottom: 20px;
  }
}
@media screen and (max-width: 992px) {
}
</style>
