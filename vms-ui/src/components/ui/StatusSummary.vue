<script setup lang="ts">
import {formatSIUnit} from "@/utils/number.ts";
import {useInverterStore} from "@/stores/inverter.ts";
import {storeToRefs} from "pinia";

const store = useInverterStore()
const {generalStatus} = storeToRefs(store)
</script>

<template>
  <div class="status-summary">
    <div class="stat">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
        <path
          d="M12 18C8.68629 18 6 15.3137 6 12C6 8.68629 8.68629 6 12 6C15.3137 6 18 8.68629 18 12C18 15.3137 15.3137 18 12 18ZM12 16C14.2091 16 16 14.2091 16 12C16 9.79086 14.2091 8 12 8C9.79086 8 8 9.79086 8 12C8 14.2091 9.79086 16 12 16ZM11 1H13V4H11V1ZM11 20H13V23H11V20ZM3.51472 4.92893L4.92893 3.51472L7.05025 5.63604L5.63604 7.05025L3.51472 4.92893ZM16.9497 18.364L18.364 16.9497L20.4853 19.0711L19.0711 20.4853L16.9497 18.364ZM19.0711 3.51472L20.4853 4.92893L18.364 7.05025L16.9497 5.63604L19.0711 3.51472ZM5.63604 16.9497L7.05025 18.364L4.92893 20.4853L3.51472 19.0711L5.63604 16.9497ZM23 11V13H20V11H23ZM4 11V13H1V11H4Z"></path>
      </svg>
      <span class="watt">{{ formatSIUnit(generalStatus?.photovoltaicChargingPower, "w") }}</span>
    </div>

    <div class="stat">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
        <path
          d="M6 18H18V6H6V18ZM14 20H10V22H8V20H5C4.44772 20 4 19.5523 4 19V16H2V14H4V10H2V8H4V5C4 4.44772 4.44772 4 5 4H8V2H10V4H14V2H16V4H19C19.5523 4 20 4.44772 20 5V8H22V10H20V14H22V16H20V19C20 19.5523 19.5523 20 19 20H16V22H14V20ZM8 8H16V16H8V8Z"></path>
      </svg>
      <span class="watt">{{ formatSIUnit(generalStatus?.alternatingCurrentOutputActivePower, "w") }}</span>
    </div>

    <div class="stat">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
        <path
          d="M8 19H3C2.44772 19 2 18.5523 2 18V6C2 5.44772 2.44772 5 3 5H9.625L8.45833 7H4V17H8V19ZM12.375 19L13.5417 17H18V7H14V5H19C19.5523 5 20 5.44772 20 6V18C20 18.5523 19.5523 19 19 19H12.375ZM21 9H23V15H21V9ZM12 11H15L10 19V13H7L12 5V11Z"></path>
      </svg>
      <span class="watt">{{ formatSIUnit(generalStatus?.batteryChargingPower, "w") }}</span>
    </div>
  </div>
</template>

<style>
.status-summary {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  justify-items: center;
  margin-block: .5rem 2rem;

  svg {
    width: 32px;
    display: block;
    margin-block-end: .3rem;
  }

  .stat {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  span {
    font-weight: 500;
    font-size: 1.2rem;
  }
}
</style>
