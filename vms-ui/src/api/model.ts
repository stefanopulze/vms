import type {ChargerSources} from "@/api/api.ts";

export interface RatingInfo {
  "timestamp": Date,
  "gridRatingVoltage": number,
  "gridRatingCurrent": number,
  "alternatingCurrentRatingVoltage": number,
  "alternatingCurrentRatingFrequency": number,
  "alternatingCurrentRatingCurrent": number,
  "alternatingCurrentRatingApparentPower": number,
  "alternatingCurrentRatingActivePower": number,
  "batteryRatingVoltage": number,
  "batteryRechargeVoltage": number,
  "batteryUnderVoltage": number,
  "batteryBulkVoltage": number,
  "batteryFloatVoltage": number,
  "batteryType": number,
  "maxAlternatingCurrentChargingCurrent": number,
  "maxSolarChargeControllerChargingCurrent": number,
  "inputVoltageRange": number,
  "outputSourcePriority": number,
  "chargerSourcePriority": number,
  "parallelMaxNum": number,
  "machineType": number,
  "topology": number,
  "outputMode": number,
  "batteryRedischargeVoltage": number,
  "photovoltaicOkConditionForParallel": number,
  "photovoltaicPowerBalance": number,
  "maximumChargingTimeAtCVStage": number,
  "operationLogic": number,
  "maxDischargingCurrent": number,
  chargerSourcePriorityEnum: ChargerSources
}

export interface GeneralStatus {
  timestamp: string
  gridVoltage: number
  gridFrequency: number
  alternatingCurrentOutputVoltage: number
  alternatingCurrentOutputFrequency: number
  alternatingCurrentOutputApparentPower: number
  alternatingCurrentOutputActivePower: number
  outputLoadPercent: number
  busVoltage: number
  batteryVoltage: number
  batteryChargingCurrent: number
  batteryCapacity: number
  inverterHeatSinkTemperature: number
  photovoltaicInputCurrent: number
  photovoltaicInputVoltage: number
  batteryVoltageFromSolarChargeController: number
  batteryDischargeCurrent: number
  batteryDischargingPower: number
  deviceStatus8bit: string
  deviceStatus8bitAddSbuPriorityVersion: number
  deviceStatus8bitConfigurationChangedStatus: number
  deviceStatus8bitSolarChargeControllerFirmwareUpdated: number
  deviceStatus8bitLoadStatus: number
  deviceStatus8bitBatteryVoltageToSteadyWhileCharging: number
  deviceStatus8bitCharging: number
  deviceStatus8bitChargingFromSolarChargeController: number
  deviceStatus8bitChargingFromAlternatingCurrent: number
  batteryVoltageOffsetForFansOn: number
  eepromVersion: number
  photovoltaicChargingPower: number
  deviceStatus3bit: string
  deviceStatus3bitFloatingCharging: number
  deviceStatus3bitSwitchOn: number
  deviceStatus3bitReserved: number
  batteryChargingPower: number
}

export interface DeviceInfo {
  firmware: Firmware
  generalModelName: string
  modelName: string
  remotePanelFirmware: Firmware
  secondCpuFirmware: Firmware
  serial: string
  serialNumber: string
}

export interface Firmware {
  Major: number
  Minor: number
}
