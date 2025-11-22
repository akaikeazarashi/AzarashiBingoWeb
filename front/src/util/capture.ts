import html2canvas from "html2canvas";

export default {
  // OGP用に1.91:1の比率になるように画像をキャプチャする
  async captureHtmlOgp(targetElement: HTMLElement) {
    const canvas = await html2canvas(targetElement, { scale: 1 });
    const origWidth = canvas.width;
    const origHeight = canvas.height;
    const origRatio = origWidth / origHeight;

    const TW_WIDTH = 1200;
    const TW_HEIGHT = 630;
    const TW_RATIO = TW_WIDTH / TW_HEIGHT;

    let drawWidth = TW_WIDTH;
    let drawHeight = TW_HEIGHT;

    if (origRatio > TW_RATIO) {
      drawWidth = TW_WIDTH;
      drawHeight = TW_WIDTH / origRatio;
    } else {
      drawHeight = TW_HEIGHT;
      drawWidth = TW_HEIGHT * origRatio;
    }

    const newCanvas = document.createElement("canvas");
    newCanvas.width = TW_WIDTH;
    newCanvas.height = TW_HEIGHT;
    const ctx = newCanvas.getContext("2d")!;

    ctx.fillStyle = "#ffffff";
    ctx.fillRect(0, 0, TW_WIDTH, TW_HEIGHT);

    ctx.drawImage(
      canvas,
      (TW_WIDTH - drawWidth) / 2,
      (TW_HEIGHT - drawHeight) / 2,
      drawWidth,
      drawHeight
    );

    return await new Promise<Blob | null>(r => newCanvas.toBlob(r, "image/png"));
  },
}