<template>
  <div class="bingo-panel-wrapper">
    <section id="bingo-panel" :style="{'--size': size}">
      <template v-for="item in items">
        <div class="bingo-panel-item" :class="{selected: item.isSelected, checked: item.isChecked, achieved: item.isAchieved}" @click="selectedItem(item)">
          <p :class="{'font-small': 20 <= item.name.length}">{{ item.name }}</p>
        </div>
      </template>
    </section>
  </div>
</template>

<script setup lang="ts">
import type { Bingo, BingoItem } from '../types/bingo';

defineOptions({
  name: "BingoPanel"
});

const props = withDefaults(
  defineProps<{
    size: number;
    items: BingoItem[];
    isSelectable?: boolean;
  }>(),
  {
    isSelectable: true
  }
);

interface Emits {
  (event: "selectedItem", item: BingoItem): void;
}

const emit = defineEmits<Emits>();

// ビンゴマス選択時の処理
const selectedItem = (item: BingoItem): void => {
  if (!props.isSelectable) {
    return;
  }
  if (item.isSelected) {
    return;
  }

  emit("selectedItem", item);

  item.isSelected = true;
  setTimeout(() => {
    item.isSelected = false;
  }, 300);
}
</script>

<style lang="scss" scoped>
.bingo-panel-wrapper {
  overflow-x: auto;
  width: 100%;
}

#bingo-panel {
  width: calc(100px * var(--size));
  --cell-size: min(100px, calc(100% / var(--size)));
  margin: 0 auto;
  border: solid 5px #ddd;
  display: grid;
  /* gridの列幅を指定する(repeatでビンゴサイズに合わせた分だけサイズを出力) */
  grid-template-columns: repeat(var(--size), var(--cell-size));
  grid-auto-rows: var(--cell-size);
  justify-content: center;
  align-content: start;
}

.bingo-panel-item {
  box-sizing: border-box;
  aspect-ratio: 1 / 1;
  width: 100%;
  height: 100%;
  display: flex;
  text-align: center;
  align-items: center;
  justify-content: center;
  border: 1px solid #ccc;
  background: #fff;
  overflow: hidden;
  user-select: none;
  padding: 5px;

  p {
    margin: 0;

    &.font-small {
      font-size: 0.8rem !important;
    }
  }

  &.selected {
    animation: selected 300ms;
  }
  &.checked {
    background-image: url("/images/bingocheck.png");
    background-position: 50%;
    background-size: cover;
  }
  &.achieved {
    animation: achieved 1s infinite;
  }
}

@keyframes selected {
  0%, 100% { outline: solid 3px rgba($color: #ccc, $alpha: 0); outline-offset: -3px; }
  50% { outline: solid 3px rgb(143, 255, 168); outline-offset: -3px; }
}

@keyframes achieved {
  0%, 100% { background-color: #fff; }
  50% { background-color: #ffeb3b; }
}
</style>
