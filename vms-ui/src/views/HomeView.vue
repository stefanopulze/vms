<script setup lang="ts">
import {onMounted, ref} from "vue";
import {api, type ChargerSources} from "@/api/api.ts";
import ButtonsToggle from "@/components/ui/ButtonsToggle.vue";
import {chargerSourcePriorityItems, sourcePriorityItems} from "@/components/enums.ts";
import BatteryPct from "@/components/ui/BatteryPct.vue";
import StatusSummary from "@/components/ui/StatusSummary.vue";
import {useInverterStore} from "@/stores/inverter.ts";
import {storeToRefs} from "pinia";

const store = useInverterStore()
const {ratingInfo, generalStatus} = storeToRefs(store)
const pending = ref(true)
const outputSourcePriority = ref(0)
const chargerSourcePriority = ref<ChargerSources>('solar_first')
const formatter = new Intl.DateTimeFormat('it-IT', {
  year: '2-digit',
  month: '2-digit',
  day: '2-digit',
  hour: '2-digit',
  minute: '2-digit',
  second: '2-digit'
})

function updateOutputSourcePriority(v: number, old: number) {
  const mode = (v == 0 ? "usb" : v == 1 ? "sub" : "sbu")
  api.updateOutputSourcePriority(mode).then(
    console.log
  ).catch(err => {
    console.log(err)
    outputSourcePriority.value = old
  })
}

function updateChargerSourcePriority(v: ChargerSources, old: ChargerSources) {
  api.updateChargerSourcePriority(v).then(
    console.log
  ).catch(err => {
    console.log(err)
    chargerSourcePriority.value = old
  })
}

onMounted(async () => {
  store.fetchAll().then(() => {
    outputSourcePriority.value = ratingInfo.value?.outputSourcePriority || 0
    chargerSourcePriority.value = ratingInfo.value?.chargerSourcePriorityEnum || 'solar_first'
  }).finally(() => pending.value = false)
})
</script>

<template>
  <main class="container" v-if="!pending">
    <StatusSummary/>


    <div class="field">
      <label>Output Source Priority</label>
      <ButtonsToggle v-model="outputSourcePriority"
                     :values="sourcePriorityItems"
                     @change="updateOutputSourcePriority"/>
    </div>

    <div class="field">
      <label>Charger source priority</label>
      <ButtonsToggle v-model="chargerSourcePriority"
                     :values="chargerSourcePriorityItems"
                     @change="updateChargerSourcePriority"/>
    </div>

    <div class="field">
      <label>Battery</label>
      <BatteryPct :pct="generalStatus?.batteryCapacity || 0"></BatteryPct>
    </div>

    <div class="field">
      <label>Rating Info</label>
      <pre>{{ ratingInfo?.timestamp ? formatter.format(ratingInfo.timestamp) : '....' }}</pre>
    </div>
  </main>
  <main v-else>Loading...</main>
</template>
