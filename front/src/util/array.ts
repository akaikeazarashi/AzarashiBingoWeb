// 配列をidKeyをキーとしたRecordに変換する
function arrayToRecord<T, K extends keyof T>(arr: T[], idKey: K): Record<string, T> {
  return arr.reduce((acc, item) => {
    const rawKey = item[idKey];
    acc[Number(rawKey)] = item;
    return acc;
  }, {} as Record<string, T>);
}

export { arrayToRecord }
