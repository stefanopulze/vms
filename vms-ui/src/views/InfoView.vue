<script setup lang="ts">
import {onMounted, ref} from "vue";
import {api} from "@/api/api.ts";
import type {DeviceInfo} from "@/api/model.ts";

const deviceInfo = ref<DeviceInfo>()

onMounted(async () => {
  try {
    deviceInfo.value = await api.fetchInfo()
    console.log(deviceInfo.value)
  } catch (e) {
    console.error(e)
  }
})
</script>

<template>
  <section class="container">
    <dl v-if="deviceInfo" class="device-info">
      <dt class="label">Model Name</dt>
      <dd>{{ deviceInfo.modelName }}</dd>
      <dt class="label">General model name</dt>
      <dd>{{ deviceInfo.generalModelName }}</dd>
      <dt class="label">Serial</dt>
      <dd>{{ deviceInfo.serial }}</dd>
      <dt class="label">Serial Number</dt>
      <dd>{{ deviceInfo.serialNumber }}</dd>
      <dt class="label">Firmware</dt>
      <dd>{{ deviceInfo.firmware.Major }}.{{ deviceInfo.firmware.Minor }}</dd>
      <template v-if="deviceInfo.secondCpuFirmware && deviceInfo.secondCpuFirmware.Major > 0">
        <dt class="label">Second CPU Firmware</dt>
        <dd>{{ deviceInfo.secondCpuFirmware.Major }}.{{ deviceInfo.secondCpuFirmware.Minor }}</dd>
      </template>
      <dt class="label">Remote panel firmware</dt>
      <dd>{{ deviceInfo.remotePanelFirmware.Major }}.{{ deviceInfo.remotePanelFirmware.Minor }}</dd>
    </dl>
  </section>
</template>

<style>
.device-info {
  margin-block-start: 0;

  dt {
    margin-block-end: .3rem;
  }

  dd {
    margin-inline-start: 0;
    margin-block-end: 1rem;
  }
}
</style>
