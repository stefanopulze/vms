<script setup lang="ts">
import {onMounted, ref} from "vue";
import {api, type ChargerSources} from "@/api/api.ts";
import ButtonsToggle from "@/components/ui/ButtonsToggle.vue";
import {chargerSourcePriorityItems, sourcePriorityItems} from "@/components/enums.ts";
import BatteryPct from "@/components/ui/BatteryPct.vue";
import StatusSummary from "@/components/ui/StatusSummary.vue";
import {useInverterStore} from "@/stores/inverter.ts";
import {storeToRefs} from "pinia";
import BatteryAcChargingCurrent from "@/components/ui/BatteryAcChargingCurrent.vue";

const store = useInverterStore()
const {ratingInfo, generalStatus} = storeToRefs(store)
const pending = ref(true)
const outputSourcePriority = ref(0)
const chargerSourcePriority = ref<ChargerSources>('solar_first')
const pendingStatus = ref<'pending' | 'success' | 'error'>('pending')
const formatter = new Intl.DateTimeFormat('it-IT', {
  year: '2-digit',
  month: '2-digit',
  day: '2-digit',
  hour: '2-digit',
  minute: '2-digit',
  second: '2-digit'
})
let timerId = -1

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

function fetchData() {
  clearTimeout(timerId)
  store.fetchAll().then(() => {
    outputSourcePriority.value = ratingInfo.value?.outputSourcePriority || 0
    chargerSourcePriority.value = ratingInfo.value?.chargerSourcePriorityEnum || 'solar_first'
    pendingStatus.value = 'success'
  }).catch(err => {
    console.error(err.message)
    pendingStatus.value = 'error'
  })
    .finally(() => {
      pending.value = false
      // timerId = setTimeout(fetchData, 5_000)
    })
}

onMounted(fetchData)
</script>

<template>
  <div class="update-status" :class="{'success': pendingStatus === 'success'}"></div>
  <main class="container" v-if="!pending">
    <StatusSummary/>

    <div class="field">
      <label>Battery</label>
      <BatteryPct :pct="generalStatus?.batteryCapacity || 0"></BatteryPct>
    </div>

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
      <label>Max AC Charging current</label>
      <BatteryAcChargingCurrent :value="ratingInfo?.maxAlternatingCurrentChargingCurrent || 0"/>
    </div>
  </main>
  <main v-else>Loading...</main>
</template>


<style>
.update-status {
  background-color: #f0f0f0;
  font-size: .9rem;
  height: 4px;

  &.success {
    background-color: #c5fdc5;
  }
}
</style>
