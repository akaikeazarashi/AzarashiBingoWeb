<template>
  <div class="progress-container">
    <div class="progress-label">
      {{ currentValue }} / {{ maxValue }}
    </div>
    <div class="progress-bar-background">
      <div class="progress-bar-current" :style="{ width: `${valueRatio * 100}%`}"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  currentValue: number;
  maxValue: number;
}>();

const valueRatio = computed((): number => {
  if (props.currentValue <= 0) {
    return 0;
  }
  return props.currentValue / props.maxValue;
});

</script>

<style lang="scss" scoped>
.progress-container {
  width: 100%;
  height: 20px;
  margin: 10px 0;
  position: relative;
}
.progress-label {
  position: absolute;
  inset: 0;
  z-index: 3;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 12px;
  font-weight: bold;
  color: #ffffff;
  text-shadow: 0 0 3px rgba(0, 0, 0, 0.8);
  pointer-events: none;
}
.progress-bar-background {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  border-radius: 10px;
  background-color: #ccc;
}
.progress-bar-current {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: #53ff5c;
  transition: width 0.3s ease;
}
</style>
