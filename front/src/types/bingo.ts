// ビンゴデータ
type Bingo = {
  id: number;
  name: string;
  description: string;
  size: number;
  items: BingoItem[];
};

// ビンゴデータの1アイテム
type BingoItem = {
  id: number;
  name: string;
  orderIndex: number;
  isSelected?: boolean;
  isChecked?: boolean;
  isAchieved?: boolean;
};

// 保存用ビンゴデータ
type StorageBingo = {
  id: number;
  items: StorageBingoItem[];
};

// 保存用ビンゴデータの1アイテム
type StorageBingoItem = {
  id: number;
  orderIndex: number;
  isChecked: boolean;
};

export type { Bingo, BingoItem, StorageBingo, StorageBingoItem };
