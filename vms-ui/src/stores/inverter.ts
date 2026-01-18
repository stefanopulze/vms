import {defineStore} from "pinia";
import {ref} from "vue";
import type {GeneralStatus, RatingInfo} from "@/api/model.ts";
import {api} from "@/api/api.ts";


export const useInverterStore = defineStore('inverter', () => {
  const ratingInfo = ref<RatingInfo>()
  const generalStatus = ref<GeneralStatus>()

  async function fetchRatingInfo() {
    ratingInfo.value = await api.fetchRatingInfo()
  }

  async function fetchGeneralStatus() {
    generalStatus.value = await api.fetchGeneralStatus()
  }

  async function fetchAll() {
    return Promise.all([fetchRatingInfo(), fetchGeneralStatus()])
  }

  return {
    ratingInfo,
    generalStatus,
    fetchAll,
    fetchRatingInfo,
    fetchGeneralStatus
  }
})
