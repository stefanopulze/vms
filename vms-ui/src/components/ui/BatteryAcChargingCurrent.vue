<script setup lang="ts">
import {onMounted, ref} from "vue";
import {api} from "@/api/api.ts";
import ButtonsToggle from "@/components/ui/ButtonsToggle.vue";

const props = defineProps<{ value: number | undefined }>()

const currentValue = ref<number>(props.value || 0)
const values = ref<{ label: string, value: number}[]>([])

async function fetchValues() {
  const response = await api.fetchBatteryAcChargingCurrentValues()
  values.value = response.values.slice(0, 4).map(v => ({label:  "" + v, value: v}))
}

async function updateChargingCurrent(v: number) {
  await api.updateBatteryAcChargingCurrent(v)
}

onMounted(fetchValues)
</script>

<template>
  <ButtonsToggle v-model="currentValue" :values="values" @change="updateChargingCurrent"/>
</template>

