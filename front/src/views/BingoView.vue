<template>
  <section class="bingo-detail" v-if="bingo">
    <div class="row bingo-view-area">
      <div class="col-md-10">
        <div class="bingo-detail-container">
          <h2>{{ bingo.name }}</h2>
          <p>{{ bingo.description }}</p>
        </div>
        <div id="bingo-panel-container" class="bingo-panel-container">
          <div id="score-container">
            <BingoProgress :max-value="bingo.items.length" :current-value="score"></BingoProgress>
          </div>
          <BingoPanel :size="bingo.size" :items="bingo.items" :isSelectable="isPanelSelectable" v-on:selectedItem="selectedBingoItem"></BingoPanel>
        </div>
        <div class="result-area">
          <Transition name="fade" mode="out-in">
            <div class="bingo-submit-button-container" v-if="!isSubmitted || !bingoResult">
              <p>マスを選択するとチェックを切り替えられます。<br/>達成しているマスを全部チェックしたらボタンをクリック！</p>
              <button type="button" class="btn btn-primary btn-submit" @click="submitBingo">これでOK！</button>
            </div>
            <div class="bingo-result-container" v-else>
              <div v-if="0 < bingoResult.achievedCount">
                <p>
                  <span class="bingo-result-count-label">{{ bingoResult.achievedCount }}</span>
                  <span class="bingo-result-label">ビンゴ</span><br />
                  達成しました！
                </p>
                <p v-if="bingoResult.isAllAchieved">
                  すべて達成したあなたは<b>完全にアザラシ</b>です！<br />おめでとうございます！
                </p>
                <p v-else>
                  残りのビンゴも達成できるようにこれからも<b>アザラシ</b>になりましょう！
                </p>
              </div>
              <div v-else>
                <p>ビンゴできませんでした。<br />達成できるようにこれからも<b>アザラシ</b>になりましょう！</p>
              </div>
              <div class="row">
                <div class="col-md-6">
                  <button type="button" class="btn btn-share-x" :disabled="isSnsSubmitted" @click="postSns">
                    <img src="/images/x-logo-white.png" alt="x_twitter" />
                    <span> 結果をX(Twitter)に投稿する</span>
                  </button>
                  <button type="button" class="btn btn-primary" @click="reEditBingo">チェックし直す</button>
                </div>
              </div>
            </div>
          </Transition>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import api from '@/util/api';
import { arrayToRecord } from '@/util/array';
import router from "../router";
import BingoPanel from '@/components/BingoPanel.vue';
import BingoProgress from '@/components/BingoProgress.vue';
import { useBingoStore } from '@/stores/bingo';
import capture from "@/util/capture";
import type { Bingo, BingoItem } from '@/types/bingo';

const bingo = ref<Bingo>(null);
const isPanelSelectable = ref(true);
const isSubmitted = ref(false);
const isSnsSubmitted = ref(false);
const bingoResult = ref(null);
const shareBlob = ref<Blob | null>(null);
const bingoStore = useBingoStore();

const score = computed((): number => {
  return bingo.value.items.filter(x => x.isChecked).length;
});

const props = defineProps<{
  bingoId?: number | string,
}>();

// ビンゴデータ初期化
const initBingoPanel = async () => {
  const paramBingoId = Number(props.bingoId ?? 0);

  // ビンゴidが指定されない場合はトップページへ戻す
  if (paramBingoId === 0) {
    router.push('/');
    return;
  }

  // サーバーからビンゴデータ取得
  // (ビンゴ詳細画面に直アクセスされた場合を考慮して都度取得する)
  const response = await api.get(`/api/item/${paramBingoId}`);
  const data = response.data;

  if (!data.result) {
    throw new Error(`ビンゴデータの取得に失敗しました (id:${paramBingoId})`);
  }

  const responseBingo = data.item as Bingo;
  const storageBingo = bingoStore.getBingo(paramBingoId);

  // サーバーのビンゴデータを正とし、ローカルの達成状態をマージする
  if (storageBingo) {
    const storageItems = arrayToRecord(storageBingo.items, "orderIndex");
    responseBingo.items.forEach(item => {
      item.isChecked = storageItems[item.orderIndex]?.isChecked ?? false;
    });
  }

  bingo.value = responseBingo;
}

// ビンゴ表でアイテム選択時に呼び出されるメソッド
const selectedBingoItem = (item: BingoItem) => {
  item.isChecked = !item.isChecked;
}

// ビンゴ確定
const submitBingo = async () => {
  isSubmitted.value = true;

  // ビンゴパネルを操作できなくする
  isPanelSelectable.value = false;

  // サーバーから結果取得
  const response = await api.post(`/api/submit`, bingo.value);
  const data = response.data;

  if (!data.result) {
    isSubmitted.value = false;
    throw new Error('ビンゴデータの送信に失敗しました');
  }

  // ストアに保存
  const bingoStore = useBingoStore();
  bingoStore.setBingo(bingo.value.id, bingo.value);

  // SNS投稿用のキャプチャを取得
  // (結果反映後はマスが点滅するためここで事前に取っておく)
  const element = document.getElementById("bingo-panel");
  const blob = await capture.captureHtmlOgp(element);
  if (!blob) {
    throw new Error('ビンゴ画像のキャプチャに失敗');
  }
  shareBlob.value = blob

  // 結果をビンゴ表に反映
  bingo.value.items.forEach(item => {
    item.isAchieved = data.bingoResult.achievedItemIdList.includes(item.id);
  });

  bingoResult.value = data.bingoResult;
}

// ビンゴ表を再編集する
const reEditBingo = () => {
  // ビンゴ表の結果表示を戻す
  bingo.value.items.forEach(item => {
    item.isAchieved = false;
  });

  bingoResult.value = null;
  isSubmitted.value = false;
  isPanelSelectable.value = true;
}

// SNS投稿
const postSns = async () => {
  isSnsSubmitted.value = true;

  try {
    // S3アップロードAPIへ送信
    const formData = new FormData();
    formData.append("bingoId", String(bingo.value.id));
    formData.append("file", shareBlob.value, "bingo.png");

    // サーバーに画像を送信
    const response = await api.postWithHeader("/api/bingo/upload", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    const data = response.data;

    if (!data.result) {
      isSnsSubmitted.value = false;
      throw new Error('ビンゴ画像のアップロードに失敗しました');
    }

    // 投稿用URLを開く
    const text = encodeURIComponent(`【アザラシビンゴ】「${bingo.value.name}」で${bingoResult.value.achievedCount}ビンゴ達成しました！`);
    const postUrl = `https://x.com/intent/post?text=${text}&url=${data.url}`;
    window.open(postUrl, '_blank');
  } finally {
    isSnsSubmitted.value = false;
  }
}

onMounted(initBingoPanel);
</script>

<style lang="scss" scoped>
.bingo-view-area {
  margin: 20px 0;
  justify-content: center;
  align-items: center;
}

.bingo-detail-container {
  margin: 20px 0;
  text-align: center;
}

.bingo-panel-container {
  margin: 20px 0;
  text-align: center;
}

#score-container {
  margin: 0 auto;
  max-width: 200px;
}

.result-area {
  min-height: 250px;
}

.bingo-submit-button-container {
  text-align: center;

  .btn-submit {
    padding: 10px 30px;
  }
}

.bingo-result-container {
  text-align: center;

  .row {
    justify-content: center;
  }

  .bingo-result-label {
    font-size: 1.5rem;
  }

  .bingo-result-count-label {
    font-size: 2.5rem;
    font-weight: bold;
    padding: 0 10px;
    color: red;
  }

  button {
    margin: 10px 10px;
  }
}

.btn-share-x {
  background-color: #000000;
  vertical-align: middle;
  font-size: 1rem;
  color: #FFF;

  & > span {
    vertical-align: middle;
  }

  & > img {
    vertical-align: middle;
    max-width: 30px;
  }
}
</style>
