package voltronic

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
	"vms-core/internal/utils"
)

func (c *Client) QueryFirmware() (*Firmware, error) {
	data, err := c.SendCommand("QVFW")
	if err != nil {
		return nil, err
	}

	// VERFW:00046.82
	return parseFirmwareVersion(data, "VERFW")
}

func (c *Client) QuerySecondCpuFirmware() (*Firmware, error) {
	data, err := c.SendCommand("QVFW2")
	if err != nil {
		return nil, err
	}

	// VERFW:00046.82
	return parseFirmwareVersion(data, "VERFW2")
}

func (c *Client) QueryRemotePanelFirmware() (*Firmware, error) {
	data, err := c.SendCommand("QVFW3")
	if err != nil {
		return nil, err
	}

	// VERFW:00046.82
	return parseFirmwareVersion(data, "VERFW3")
}

func (c *Client) QueryPIRI() (*DeviceRatingInfo, error) {
	data, err := c.SendCommand("QPIRI")
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(string(data[1:]))

	c.ratingInfoMux.Lock()
	c.ratingInfo.Timestamp = time.Now()
	c.ratingInfo.GridRatingVoltage = utils.ParseFloatAsInt16(fields[0])                 // 230
	c.ratingInfo.GridRatingCurrent = utils.ParseFloat(fields[1])                        // 34.7
	c.ratingInfo.AlternatingCurrentRatingVoltage = utils.ParseFloatAsInt16(fields[2])   // 230
	c.ratingInfo.AlternatingCurrentRatingFrequency = utils.ParseFloatAsInt16(fields[3]) // 50
	c.ratingInfo.AlternatingCurrentRatingCurrent = utils.ParseFloat(fields[4])          // 34.7
	c.ratingInfo.AlternatingCurrentRatingApparentPower = utils.ParseInt(fields[5])      // 8000
	c.ratingInfo.AlternatingCurrentRatingActivePower = utils.ParseInt(fields[6])        // 8000
	c.ratingInfo.BatteryRatingVoltage = utils.ParseFloat(fields[7])                     // 48
	c.ratingInfo.BatteryRechargeVoltage = utils.ParseFloat(fields[8])                   // 48
	c.ratingInfo.BatteryUnderVoltage = utils.ParseFloat(fields[9])                      // 45
	c.ratingInfo.BatteryBulkVoltage = utils.ParseFloat(fields[10])                      // 53.2
	c.ratingInfo.BatteryFloatVoltage = utils.ParseFloat(fields[11])                     // 53.2
	c.ratingInfo.BatteryType = utils.ParseInt(fields[12])                               // 3 0=AGM, 1=Flooded, 2=User, 3=PYL, 4=SH
	c.ratingInfo.MaxAlternatingCurrentChargingCurrent = utils.ParseInt(fields[13])      // 120
	c.ratingInfo.MaxSolarChargeControllerChargingCurrent = utils.ParseInt(fields[14])   // 1
	c.ratingInfo.InputVoltageRange = utils.ParseInt(fields[15])                         // 2 0=Appliance, 1=UPS
	c.ratingInfo.OutputSourcePriority = utils.ParseInt(fields[16])                      // 3 0=Utility first, 1=solar, 2=SBU
	c.ratingInfo.ChargerSourcePriority = utils.ParseInt(fields[17])                     // 9
	c.ratingInfo.ParallelMaxNum = utils.ParseInt(fields[18])                            // 1
	c.ratingInfo.MachineType = uint8(utils.ParseInt(fields[19]))                        // "0"
	c.ratingInfo.Topology = utils.ParseInt(fields[20])                                  // 0
	c.ratingInfo.OutputMode = utils.ParseInt(fields[21])                                // 48
	c.ratingInfo.BatteryRedischargeVoltage = utils.ParseFloat(fields[22])               // 48.5 (considerando "48.5" o separato)
	c.ratingInfo.PhotovoltaicOkConditionForParallel = utils.ParseInt(fields[23])        // 0
	c.ratingInfo.PhotovoltaicPowerBalance = utils.ParseInt(fields[24])                  // 1
	c.ratingInfo.MaximumChargingTimeAtCVStage = utils.ParseInt(fields[25])              // 480
	c.ratingInfo.OperationLogic = utils.ParseInt(fields[26])                            // 0
	c.ratingInfo.MaxDischargingCurrent = utils.ParseInt(fields[27])                     // 0 (se presente)
	c.ratingInfoMux.Unlock()

	return &c.ratingInfo, err
}

func (c *Client) QueryPIGS() (*DeviceGeneralStatus, error) {
	data, err := c.SendCommand("QPIGS")
	if err != nil {
		return nil, err
	}

	//slog.Debug(string(data))
	fields := strings.Fields(string(data[1:]))

	c.generalStatusMux.Lock()
	c.generalStatus.Timestamp = time.Now()
	c.generalStatus.GridVoltage = utils.ParseFloat(fields[0])                                                                   // 226.9
	c.generalStatus.GridFrequency = utils.ParseFloat(fields[1])                                                                 // 49.9
	c.generalStatus.AlternatingCurrentOutputVoltage = utils.ParseFloat(fields[2])                                               // 229.8
	c.generalStatus.AlternatingCurrentOutputFrequency = utils.ParseFloat(fields[3])                                             // 49.9
	c.generalStatus.AlternatingCurrentOutputApparentPower = utils.ParseInt(fields[4])                                           // 0575
	c.generalStatus.AlternatingCurrentOutputActivePower = utils.ParseInt(fields[5])                                             // 0499
	c.generalStatus.OutputLoadPercent = utils.ParseInt(fields[6])                                                               // 007
	c.generalStatus.BusVoltage = utils.ParseInt(fields[7])                                                                      // 361
	c.generalStatus.BatteryVoltage = utils.ParseFloat(fields[8])                                                                //  49.30
	c.generalStatus.BatteryChargingCurrent = utils.ParseInt(fields[9])                                                          // 000
	c.generalStatus.BatteryChargingPower = float64(c.generalStatus.BatteryChargingCurrent) * c.generalStatus.BatteryVoltage     // 000
	c.generalStatus.BatteryCapacity = utils.ParseInt(fields[10])                                                                // 082
	c.generalStatus.InverterHeatSinkTemperature = utils.ParseInt(fields[11])                                                    // 0025
	c.generalStatus.PhotovoltaicInputCurrent = utils.ParseFloat(fields[12])                                                     //  00.4
	c.generalStatus.PhotovoltaicInputVoltage = utils.ParseFloat(fields[13])                                                     // 225.9
	c.generalStatus.BatteryVoltageFromSolarChargeController = utils.ParseFloat(fields[14])                                      // 00.00
	c.generalStatus.BatteryDischargeCurrent = utils.ParseInt(fields[15])                                                        // 00008
	c.generalStatus.BatteryDischargingPower = float64(c.generalStatus.BatteryDischargeCurrent) * c.generalStatus.BatteryVoltage // 00008
	c.generalStatus.DeviceStatus8Bit = fields[16]                                                                               // 00010110
	c.generalStatus.BatteryVoltageOffsetForFansOn = utils.ParseInt(fields[17])                                                  // 00
	c.generalStatus.EepromVersion = utils.ParseInt(fields[18])                                                                  // 00
	c.generalStatus.PhotovoltaicChargingPower = utils.ParseInt(fields[19])                                                      // 00099
	c.generalStatus.DeviceStatus3Bit = fields[20]                                                                               // 010
	//DeviceStatus3BitFloatingCharging:        utils.ParseInt(fields[20][:1]),
	//DeviceStatus3BitSwitchOn:                utils.ParseInt(fields[20][1:2]),
	//DeviceStatus3BitReserved:                utils.ParseInt(fields[20][2:]),

	c.generalStatusMux.Unlock()

	return &c.generalStatus, err
}

func (c *Client) QueryMode() (*DeviceMode, error) {
	data, err := c.SendCommand("QMOD")
	if err != nil {
		return nil, err
	}

	var mode string
	switch string(data[1:]) {
	case "P":
		mode = "power_on"
	case "S":
		mode = "standby"
	case "L":
		mode = "line_mode"
	case "B":
		mode = "battery_mode"
	case "F":
		mode = "fault_mode"
	case "D":
		mode = "shutdown_mode"
	case "C":
		mode = "charging_mode"
	case "Y":
		mode = "bypass_mode"
	case "E":
		mode = "eco_mode"
	default:
		return nil, fmt.Errorf("unknown mode: %s", string(data[1:]))
	}

	return &DeviceMode{
		Timestamp: time.Now(),
		Mode:      mode,
	}, nil
}

func (c *Client) QueryWarning() (*DeviceWarning, error) {
	data, err := c.SendCommand("QPIWS")
	if err != nil {
		return nil, err
	}

	warning := &DeviceWarning{
		Timestamp:              time.Now(),
		PVLoss:                 parseBoolean(data[1:2]),
		InverterFault:          parseBoolean(data[2:3]),
		BusOver:                parseBoolean(data[3:4]),
		BusUnder:               parseBoolean(data[4:5]),
		BusSoftFail:            parseBoolean(data[5:6]),
		LineFail:               parseBoolean(data[6:7]),
		OpvShort:               parseBoolean(data[7:8]),
		InverterVoltageTooLow:  parseBoolean(data[8:9]),
		InverterVoltageTooHigh: parseBoolean(data[9:10]),
		OverTemperature:        parseBoolean(data[10:11]),
		FanLocked:              parseBoolean(data[11:12]),
		BatteryVoltageHigh:     parseBoolean(data[12:13]),
		BatteryVoltageLow:      parseBoolean(data[13:14]),
		BatteryUnderShutdown:   parseBoolean(data[14:15]),
		BatteryDerating:        parseBoolean(data[15:16]),
		OverLoad:               parseBoolean(data[16:17]),
		EepromFault:            parseBoolean(data[17:18]),
		InvertOverCurrent:      parseBoolean(data[18:19]),
		InvertSoftFail:         parseBoolean(data[19:20]),
		SelfTestFail:           parseBoolean(data[20:21]),
		OpDCVoltageOver:        parseBoolean(data[21:22]),
		BatteryOpen:            parseBoolean(data[22:23]),
		CurrentSensorFail:      parseBoolean(data[23:24]),
		PVVoltageHigh:          parseBoolean(data[26:27]),
		PVOverCurrent:          parseBoolean(data[27:28]),
		DCDCOverCurrent:        parseBoolean(data[30:31]),
	}

	return warning, err
}

func (c *Client) QueryTime() (*time.Time, error) {
	data, err := c.SendCommand("QT")
	if err != nil {
		return nil, err
	}

	parse, err := time.Parse("20060102150405", string(data[1:]))
	if err != nil {
		return nil, err
	}

	return &parse, err
}

func (c *Client) QuerySerialNumber() (string, error) {
	s, err := c.SendCommand("QSID")
	if err != nil {
		return "", err
	}

	return string(s[1:]), nil
}

func (c *Client) QuerySerial() (string, error) {
	s, err := c.SendCommand("QID")
	if err != nil {
		return "", err
	}

	return string(s[1:]), nil
}

func (c *Client) QueryModelName() (string, error) {
	s, err := c.SendCommand("QMN")
	if err != nil {
		return "", err
	}

	return string(s[1:]), nil
}

func (c *Client) QueryGeneralModelName() (string, error) {
	s, err := c.SendCommand("QGMN")
	if err != nil {
		return "", err
	}

	return string(s[1:]), nil
}

func (c *Client) QueryFlags() (*DeviceFlags, error) {
	data, err := c.SendCommand("QFLAG")
	if err != nil {
		return nil, err
	}

	flags := &DeviceFlags{}
	enabled := true
	for _, b := range data[1:] {
		if b == 'E' {
			enabled = true
			continue
		} else if b == 'D' {
			enabled = false
			continue
		}

		switch b {
		case 'a':
			flags.SilenceBuzzer = enabled
		case 'b':
			flags.OverloadBypass = enabled
		case 'j':
			flags.PowerSaving = enabled
		case 'k':
			flags.LcdEscape = enabled
		case 'u':
			flags.OverloadRestart = enabled
		case 'v':
			flags.OverTemperatureRestart = enabled
		case 'x':
			flags.Backlight = enabled
		case 'y':
			flags.AlarmPrimarySourceInterrupt = enabled
		case 'z':
			flags.FaultCodeRecord = enabled
		default:
			slog.Warn(fmt.Sprintf("unknown flag: %c", b))
		}
	}

	return flags, nil
}
