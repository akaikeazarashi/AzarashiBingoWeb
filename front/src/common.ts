const sizeList = [3, 4, 5, 6, 7, 8];

export default {
  /**
   * アプリ名取得
   */
  getAppName() {
    return 'アザラシビンゴwebアプリ';
  },

  /**
   * ビンゴサイズのリストを取得する
   */
  getSizeList() {
    return sizeList;
  },

  /**
   * 指定したミリ秒待機する
   */
  sleep(msec) {
    return new Promise(resolve => setTimeout(resolve, msec))
  }
}
