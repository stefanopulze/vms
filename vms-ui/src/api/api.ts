import ky from "ky";
import type {BatteryAcChargingCurrentValues, DeviceInfo, GeneralStatus, RatingInfo} from "@/api/model.ts";
import {customJsonReceiver} from "@/api/utils.ts";

const kyApi = ky.create({
  parseJson: customJsonReceiver,
  retry: {
    statusCodes: [408, 502, 503, 504]
  }
})

class Api {
  async fetchInfo() {
    return await kyApi.get<DeviceInfo>('/api/inverter/info').json()
  }

  async fetchMode() {
    return await kyApi.get('/api/inverter/mode').json()
  }

  async fetchGeneralStatus() {
    return await kyApi.get<GeneralStatus>('/api/inverter/general-status').json()
  }

  async fetchRatingInfo() {
    return await kyApi.get<RatingInfo>('/api/inverter/rating-info').json()
  }

  async updateOutputSourcePriority(source: 'usb' | 'sub' | 'sbu') {
    return await kyApi.put<RatingInfo>('/api/inverter/source-priority', {
      json: {source}
    }).json()
  }

  async updateChargerSourcePriority(source: ChargerSources) {
    return await kyApi.put<RatingInfo>('/api/inverter/charger-source-priority', {
      json: {source}
    }).json()
  }

  async fetchBatteryAcChargingCurrentValues() {
    return await kyApi.get<BatteryAcChargingCurrentValues>('/api/inverter/max-ac-charging-current-values').json()
  }

  async updateBatteryAcChargingCurrent(v: number) {
    return await kyApi.put('/api/inverter/max-ac-charging-current', {
      json: {current: v}
    }).text()
  }
}

export type ChargerSources = 'solar_first' | 'solar_utility' | 'only_solar'

export const api = new Api()
