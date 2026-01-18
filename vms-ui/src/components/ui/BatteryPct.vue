<script setup lang="ts">
defineProps<{ pct: number }>()
</script>

<template>
  <div class="battery-pct">
    <progress :value="pct" max="100"></progress>
    <span>{{ pct }}%</span>
  </div>
</template>

<style>
.battery-pct {
  display: flex;
  gap: 1rem;
  align-items: center;

  progress {
    width: 100%;
    display: block;
    height: 30px;
  }
}

</style>
