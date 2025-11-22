<template>
  <div>
    <div class="row justify-content-center">
      <div id="top-contents" class="row text-center justify-content-center">
        <div class="col-xl-6">
          <h1>アザラシビンゴ</h1>
          <p>挑戦するビンゴを選んでください！</p>
          <Transition name="fade" mode="out-in">
            <section class="row bingo-panel-card-container" key="bingo-list" v-if="!isLoading">
              <div class="bingo-panel-card card" v-for="bingo in bingoList">
                <div class="card-header">{{ bingo.name }}</div>
                <div class="card-body">
                  <div class="card-text">
                    <p>{{ bingo.description }}</p>
                    <div class="row bingo-panel-card-footer">
                      <div class="btn-start">
                        <RouterLink class="btn btn-primary" :to="`/bingo/${bingo.id}`" @click="startBingo(bingo)">挑戦する</RouterLink>
                      </div>
                      <div class="bingo-progress">
                        <BingoProgress :max-value="bingo.items.length" :current-value="getBingoScore(bingo.id)"></BingoProgress>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </section>
            <section class="loading-container" key="loading" v-else>
              <img src="/images/loading.gif">
            </section>
          </Transition>
          <section class="notice-area">
            <h2>アザラシビンゴについて</h2>
            <ul>
              <li>ページ表示が上手くいかない場合はリロードをお試しください。</li>
              <li>データは閲覧中のブラウザ内に保存されます。<br />ブラウザや端末を変更すると別のデータになりますのでご注意ください。</li>
              <li>データを初期化する場合は<a href="" @click="onClickResetData">こちら</a>をクリックしてください。</li>
              <li>ビンゴの内容や大きさは予告なく変更・削除される可能性があります。</li>
              <li>本サイトの内容の無断転載・複製を禁じます。</li>
            </ul>
            <h2 class="mt-5">参考サイト様(外部リンク)</h2>
            <ul>
              <li><a href="https://sites.google.com/view/azarashisetsu/%E3%82%A2%E3%82%B6%E3%83%A9%E6%96%BD%E8%A8%AD%E3%81%BE%E3%81%A8%E3%82%81" target="_blank">アザラシデータベース</a></li>
            </ul>
          </section>
        </div>
      </div>
    </div>
    <Modal v-model="isShowModal" title="データ初期化確認">
      <p>初期化すると元に戻せません。よろしいですか？</p>
      <div style="text-align: center;">
        <button type="button" class="btn btn-primary mt-3" @click="submitDataReset()">初期化する</button>
        <button type="button" class="btn btn-danger mt-3 ms-3" @click="cancelDataReset()">キャンセル</button>
      </div>
    </Modal>    
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { RouterLink } from 'vue-router';
import api from '@/util/api';
import router from "../router";
import BingoProgress from '@/components/BingoProgress.vue';
import { useBingoStore } from '@/stores/bingo';
import type { Bingo } from '@/types/bingo';

const bingoList = ref<Bingo[]>([]);
const bingoStore = useBingoStore();
const isLoading = ref(true);
const isShowModal = ref(false);

// ビンゴ一覧の初期化
const initList = async () => {
  isLoading.value = true;

  const response = await api.get('/api/items');
  const data = response.data;

  if (!data.result) {
    throw new Error('リスト取得に失敗しました。');
  }

  bingoList.value = data.bingoList.map((item: Bingo) => ({
    id: item.id,
    name: item.name,
    description: item.description,
    size: item.size,
    items: item.items
  }));

  isLoading.value = false;
};

// ビンゴ開始処理
const startBingo = (bingo: Bingo): void => {
  router.push(`/bingo/${bingo.id}`);
}

// ビンゴスコア取得
const getBingoScore = (id: number): number => {
  return bingoStore.getBingoScore(id);
}

// ビンゴデータリセット押下時
const onClickResetData = (event: Event) => {
  event.preventDefault();
  isShowModal.value = true;
}

// ビンゴデータリセット確定時
const submitDataReset = () => {
  bingoStore.ResetAll();
  isShowModal.value = false;
}

// ビンゴデータリセットキャンセル時
const cancelDataReset = () => {
  isShowModal.value = false;
}

onMounted(initList);
</script>

<style lang="scss">
#top-contents {
  margin: 20px 0;
  padding: 0;

  .bingo-panel-card-container {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    justify-content: flex-start;
  }

  .bingo-panel-card {
    /* 両サイドのgap分を考慮してcontainerのgap/2を引く */
    flex: 0 0 calc(50% - 10px);
    max-width: calc(50% - 10px);
    margin: 0;
    text-align: left;
    padding: 0;
  }
  @media (max-width: 991px) {
    .bingo-panel-card {
      /* ブレイクポイント以下の場合は100%を指定して1カラムにする */
      flex: 0 0 100%;
      max-width: 100%;
    }
  }

  .card-header {
    background-color: rgb(179, 225, 255);
  }

  .btn-start {
    width: 115px;
  }

  .bingo-progress {
    padding: 0;
    width: 125px;
  }

  .loading-container {
    min-height: 50vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .notice-area {
    text-align: left;
    margin-top: 50px;
  }
}

@media screen and (max-width: 1200px) {
  #top-contents {
    margin-top: 20px;
  }

  .bingo-panel-card-container {
    margin: 0 10px;
  }

  .bingo-panel-card-footer {
    display: flex;
  }
}
</style>

