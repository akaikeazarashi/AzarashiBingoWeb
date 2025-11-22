import { defineStore } from 'pinia';
import type { Bingo, StorageBingo } from '@/types/bingo';

interface State {
  bingoList: Record<number, StorageBingo>;
}

export const useBingoStore = defineStore('bingo', {
  state: (): State => {
    return {
      bingoList: {}
    }
  },
  getters: {
    getBingo: (state) => {
      return (id: number): StorageBingo | undefined => {
        return state.bingoList[id] ?? null;
      };
    }
  },
  actions: {
    initStore(): void {
      let bingoDataList = localStorage.getItem('bingoList');
      if (bingoDataList) {
        this.bingoList = JSON.parse(bingoDataList);
      }
    },
    getBingoScore(bingoId: number): number {
      const bingo = this.bingoList[bingoId];
      if (!bingo) {
        return 0;
      }
      return bingo.items.filter(item => item.isChecked).length;
    },
    setBingo(bingoId: number, bingo: Bingo): void {
      const storageBingo = this.toStorageBingo(bingo);
      this.bingoList[bingoId] = storageBingo;
      localStorage.setItem('bingoList', JSON.stringify(this.bingoList));
    },
    removeBingo(bingoId: number): void {
      delete this.bingoList[bingoId]
      localStorage.setItem('bingoList', JSON.stringify(this.bingoList));
    },
    ResetAll(): void {
      localStorage.removeItem('bingoList');
      this.bingoList = {};
    },

    toStorageBingo(bingo: Bingo): StorageBingo {
      const storageItems = [];

      bingo.items.forEach(item => {
        storageItems.push({
          id: item.id,
          orderIndex: item.orderIndex,
          isChecked: item.isChecked ?? false,
        });
      });

      return {
        id: bingo.id,
        items: storageItems,
      };
    }
  }
});
