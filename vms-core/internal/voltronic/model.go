package voltronic

import "time"

type Firmware struct {
	Major int
	Minor int
}

type DeviceRatingInfo struct {
	Timestamp                               time.Time `json:"timestamp"`
	GridRatingVoltage                       uint16    `json:"gridRatingVoltage"`
	GridRatingCurrent                       float64   `json:"gridRatingCurrent"`
	AlternatingCurrentRatingVoltage         uint16    `json:"alternatingCurrentRatingVoltage"`
	AlternatingCurrentRatingFrequency       uint16    `json:"alternatingCurrentRatingFrequency"`
	AlternatingCurrentRatingCurrent         float64   `json:"alternatingCurrentRatingCurrent"`
	AlternatingCurrentRatingApparentPower   int       `json:"alternatingCurrentRatingApparentPower"`
	AlternatingCurrentRatingActivePower     int       `json:"alternatingCurrentRatingActivePower"`
	BatteryRatingVoltage                    float64   `json:"batteryRatingVoltage"`
	BatteryRechargeVoltage                  float64   `json:"batteryRechargeVoltage"`
	BatteryUnderVoltage                     float64   `json:"batteryUnderVoltage"`
	BatteryBulkVoltage                      float64   `json:"batteryBulkVoltage"`
	BatteryFloatVoltage                     float64   `json:"batteryFloatVoltage"`
	BatteryType                             int       `json:"batteryType"`
	MaxAlternatingCurrentChargingCurrent    int       `json:"maxAlternatingCurrentChargingCurrent"`
	MaxSolarChargeControllerChargingCurrent int       `json:"maxSolarChargeControllerChargingCurrent"`
	InputVoltageRange                       int       `json:"inputVoltageRange"`
	OutputSourcePriority                    int       `json:"outputSourcePriority"`
	ChargerSourcePriority                   int       `json:"chargerSourcePriority"`
	ParallelMaxNum                          int       `json:"parallelMaxNum"`
	MachineType                             uint8     `json:"machineType"`
	Topology                                int       `json:"topology"`
	OutputMode                              int       `json:"outputMode"`
	BatteryRedischargeVoltage               float64   `json:"batteryRedischargeVoltage"`
	PhotovoltaicOkConditionForParallel      int       `json:"photovoltaicOkConditionForParallel"`
	PhotovoltaicPowerBalance                int       `json:"photovoltaicPowerBalance"`
	MaximumChargingTimeAtCVStage            int       `json:"maximumChargingTimeAtCVStage"`
	OperationLogic                          int       `json:"operationLogic"`
	MaxDischargingCurrent                   int       `json:"maxDischargingCurrent"`
}

type DeviceGeneralStatus struct {
	Timestamp                                            time.Time `json:"timestamp"`
	GridVoltage                                          float64   `json:"gridVoltage"`
	GridFrequency                                        float64   `json:"gridFrequency"`
	AlternatingCurrentOutputVoltage                      float64   `json:"alternatingCurrentOutputVoltage"`
	AlternatingCurrentOutputFrequency                    float64   `json:"alternatingCurrentOutputFrequency"`
	AlternatingCurrentOutputApparentPower                int       `json:"alternatingCurrentOutputApparentPower"`
	AlternatingCurrentOutputActivePower                  int       `json:"alternatingCurrentOutputActivePower"`
	OutputLoadPercent                                    int       `json:"outputLoadPercent"`
	BusVoltage                                           int       `json:"busVoltage"`
	BatteryVoltage                                       float64   `json:"batteryVoltage"`
	BatteryChargingCurrent                               int       `json:"batteryChargingCurrent"`
	BatteryCapacity                                      int       `json:"batteryCapacity"`
	InverterHeatSinkTemperature                          int       `json:"inverterHeatSinkTemperature"`
	PhotovoltaicInputCurrent                             float64   `json:"photovoltaicInputCurrent"`
	PhotovoltaicInputVoltage                             float64   `json:"photovoltaicInputVoltage"`
	BatteryVoltageFromSolarChargeController              float64   `json:"batteryVoltageFromSolarChargeController"`
	BatteryDischargeCurrent                              int       `json:"batteryDischargeCurrent"`
	BatteryDischargingPower                              float64   `json:"batteryDischargingPower"`
	DeviceStatus8Bit                                     string    `json:"deviceStatus8bit"`
	DeviceStatus8BitAddSbuPriorityVersion                int       `json:"deviceStatus8bitAddSbuPriorityVersion"`
	DeviceStatus8BitConfigurationChangedStatus           int       `json:"deviceStatus8bitConfigurationChangedStatus"`
	DeviceStatus8BitSolarChargeControllerFirmwareUpdated int       `json:"deviceStatus8bitSolarChargeControllerFirmwareUpdated"`
	DeviceStatus8BitLoadStatus                           int       `json:"deviceStatus8bitLoadStatus"`
	DeviceStatus8BitBatteryVoltageToSteadyWhileCharging  int       `json:"deviceStatus8bitBatteryVoltageToSteadyWhileCharging"`
	DeviceStatus8BitCharging                             int       `json:"deviceStatus8bitCharging"`
	DeviceStatus8BitChargingFromSolarChargeController    int       `json:"deviceStatus8bitChargingFromSolarChargeController"`
	DeviceStatus8BitChargingFromAlternatingCurrent       int       `json:"deviceStatus8bitChargingFromAlternatingCurrent"`
	BatteryVoltageOffsetForFansOn                        int       `json:"batteryVoltageOffsetForFansOn"`
	EepromVersion                                        int       `json:"eepromVersion"`
	PhotovoltaicChargingPower                            int       `json:"photovoltaicChargingPower"`
	DeviceStatus3Bit                                     string    `json:"deviceStatus3bit"`
	DeviceStatus3BitFloatingCharging                     int       `json:"deviceStatus3bitFloatingCharging"`
	DeviceStatus3BitSwitchOn                             int       `json:"deviceStatus3bitSwitchOn"`
	DeviceStatus3BitReserved                             int       `json:"deviceStatus3bitReserved"`
	BatteryChargingPower                                 float64   `json:"batteryChargingPower"`
}

type DeviceWarning struct {
	Timestamp              time.Time `json:"timestamp"`
	PVLoss                 bool      `json:"pvLoss"`
	InverterFault          bool      `json:"inverterFault"`
	BusOver                bool      `json:"busOver"`
	BusUnder               bool      `json:"busUnder"`
	BusSoftFail            bool      `json:"busSoftFail"`
	LineFail               bool      `json:"lineFail"`
	OpvShort               bool      `json:"opvShort"`
	InverterVoltageTooLow  bool      `json:"inverterVoltageTooLow"`
	InverterVoltageTooHigh bool      `json:"inverterVoltageTooHigh"`
	OverTemperature        bool      `json:"overTemperature"`
	FanLocked              bool      `json:"fanLocked"`
	BatteryVoltageHigh     bool      `json:"batteryVoltageHigh"`
	BatteryVoltageLow      bool      `json:"batteryVoltageLow"`
	BatteryUnderShutdown   bool      `json:"batteryUnderShutdown"`
	BatteryDerating        bool      `json:"batteryDerating"`
	OverLoad               bool      `json:"overLoad"`
	EepromFault            bool      `json:"eepromFault"`
	InvertOverCurrent      bool      `json:"invertOverCurrent"`
	InvertSoftFail         bool      `json:"invertSoftFail"`
	SelfTestFail           bool      `json:"selfTestFail"`
	OpDCVoltageOver        bool      `json:"opDCVoltageOver"`
	BatteryOpen            bool      `json:"batteryOpen"`
	CurrentSensorFail      bool      `json:"currentSensorFail"`
	PVVoltageHigh          bool      `json:"pvVoltageHigh"`
	PVOverCurrent          bool      `json:"pvOverCurrent"`
	DCDCOverCurrent        bool      `json:"dcdcOverCurrent"`
}

func (dri DeviceRatingInfo) ChargerSourcePriorityEnum() string {
	switch dri.ChargerSourcePriority {
	case 0:
		return "utility"
	case 1:
		return "solar_first"
	case 2:
		return "solar_utility"
	case 3:
		return "only_solar"
	default:
		return "n.d."
	}
}

func (dri DeviceRatingInfo) OutputSourcePriorityEnum() string {
	switch dri.OutputSourcePriority {
	case 0:
		return "usb"
	case 1:
		return "sub"
	case 2:
		return "sbu"
	default:
		return "n.d."
	}
}

type DeviceMode struct {
	Timestamp time.Time `json:"timestamp"`
	Mode      string    `json:"mode"`
}

type DeviceFlags struct {
	AlarmPrimarySourceInterrupt bool `json:"alarmPrimarySourceInterrupt" flag:"y"`
	Backlight                   bool `json:"backlight" flag:"x"`
	FaultCodeRecord             bool `json:"faultCodeRecord" flag:"z"`
	LcdEscape                   bool `json:"lcdEscape" flag:"k"`
	OverloadBypass              bool `json:"overloadBypass" flag:"b"`
	OverloadRestart             bool `json:"overloadRestart" flag:"u"`
	OverTemperatureRestart      bool `json:"overTemperatureRestart" flag:"v"`
	PowerSaving                 bool `json:"powerSaving" flag:"j"`
	SilenceBuzzer               bool `json:"silenceBuzzer" flag:"a"`
}
