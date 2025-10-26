<template>
  <div
    v-if="visible"
    class="floating-tooltip"
    :style="tooltipStyle"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
import { reactive, onMounted, onUnmounted, CSSProperties } from "vue";

const props = defineProps<{
  visible?: boolean;
}>();

const tooltipStyle = reactive<CSSProperties>({
  position: "fixed",
  left: "0px",
  top: "0px",
  pointerEvents: "none",
  zIndex: "9999",
});

function handleMouseMove(event: MouseEvent) {
  const offsetX = 15;
  const offsetY = 15;
  tooltipStyle.left = `${event.clientX + offsetX}px`;
  tooltipStyle.top = `${event.clientY + offsetY}px`;
}

onMounted(() => window.addEventListener("mousemove", handleMouseMove));
onUnmounted(() => window.removeEventListener("mousemove", handleMouseMove));
</script>

<style scoped>

</style>
