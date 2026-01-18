package voltronic

import (
	"slices"
	"testing"
	"vms-core/testutils"
)

func TestNewClient(t *testing.T) {
	port := testutils.NewDummySerial()
	port.MockCommand("QPI", testutils.FromHex("515049BEAC0D"), testutils.FromHex("28504933309A0B0D"))
	client := NewClient(port)
	data, err := client.SendCommand("QPI")
	if err != nil {
		t.Fatal(err)
	}

	cmd := port.LastCommand()
	if !slices.Equal(cmd.Request, testutils.FromHex("515049BEAC0D")) {
		t.Error("invalid command")
	}

	if !slices.Equal(cmd.Response, testutils.FromHex("28504933309A0B0D")) {
		t.Error("invalid response")
	}

	if !slices.Equal(data, testutils.FromHex("2850493330")) {
		t.Error("invalid payload")
	}

	if "(PI30" != string(data) {
		t.Error("invalid payload")
	}
}

func TestClient_GetFirmware(t *testing.T) {
	port := testutils.NewDummySerial()
	port.MockCommand("QVFW", testutils.FromHex("5156465762990D"), testutils.FromHex("2856455246573A30303034362E383247F80D"))
	client := NewClient(port)
	data, err := client.QueryFirmware()
	if err != nil {
		t.Fatal(err)
	}

	if data.Major != 46 {
		t.Error("invalid major firmware version")
	}

	if data.Minor != 82 {
		t.Error("invalid minor firmware version")
	}
}

func TestClient_QueryPIRI(t *testing.T) {
	port := testutils.NewDummySerial()
	port.MockCommand("QPIRI",
		testutils.FromHex("5150495249F8540D"),
		testutils.FromHex("283233302E302033342E37203233302E302035302E302033342E37203830303020383030302034382E302034382E302034352E302035332E322035332E32203320303032203132302031203220332039203031203020302034382E3520302031203438302030203030309C920D"),
	)
	client := NewClient(port)
	data, err := client.QueryPIRI()
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "GridRatingVoltage", uint16(230), data.GridRatingVoltage)
	assertEqual(t, "GridRatingCurrent", 34.7, data.GridRatingCurrent)
	assertEqual(t, "AlternatingCurrentRatingVoltage", uint16(230), data.AlternatingCurrentRatingVoltage)
	assertEqual(t, "AlternatingCurrentRatingFrequency", uint16(50), data.AlternatingCurrentRatingFrequency)
	assertEqual(t, "AlternatingCurrentRatingCurrent", 34.7, data.AlternatingCurrentRatingCurrent)
	assertEqual(t, "AlternatingCurrentRatingApparentPower", 8000, data.AlternatingCurrentRatingApparentPower)
	assertEqual(t, "AlternatingCurrentRatingActivePower", 8000, data.AlternatingCurrentRatingActivePower)
	assertEqual(t, "BatteryRatingVoltage", 48.0, data.BatteryRatingVoltage)
	assertEqual(t, "BatteryRechargeVoltage", 48.0, data.BatteryRechargeVoltage)
	assertEqual(t, "BatteryUnderVoltage", 45.0, data.BatteryUnderVoltage)
	assertEqual(t, "BatteryBulkVoltage", 53.2, data.BatteryBulkVoltage)
	assertEqual(t, "BatteryFloatVoltage", 53.2, data.BatteryFloatVoltage)
	assertEqual(t, "BatteryType", 3, data.BatteryType)
	assertEqual(t, "MaxAlternatingCurrentChargingCurrent", 2, data.MaxAlternatingCurrentChargingCurrent)
	assertEqual(t, "MaxSolarChargeControllerChargingCurrent", 120, data.MaxSolarChargeControllerChargingCurrent)
	assertEqual(t, "InputVoltageRange", 1, data.InputVoltageRange)
	assertEqual(t, "OutputSourcePriority", 2, data.OutputSourcePriority)
	assertEqual(t, "ChargerSourcePriority", 3, data.ChargerSourcePriority)
	assertEqual(t, "ParallelMaxNum", 9, data.ParallelMaxNum)
	assertEqual(t, "MachineType", uint8(1), data.MachineType)
	assertEqual(t, "Topology", 0, data.Topology)
	assertEqual(t, "OutputMode", 0, data.OutputMode)
	assertEqual(t, "BatteryRedischargeVoltage", 48.5, data.BatteryRedischargeVoltage)
	assertEqual(t, "PhotovoltaicOkConditionForParallel", 0, data.PhotovoltaicOkConditionForParallel)
	assertEqual(t, "PhotovoltaicPowerBalance", 1, data.PhotovoltaicPowerBalance)
	assertEqual(t, "MaximumChargingTimeAtCVStage", 480, data.MaximumChargingTimeAtCVStage)
	assertEqual(t, "OperationLogic", 0, data.OperationLogic)
	assertEqual(t, "MaxDischargingCurrent", 0, data.MaxDischargingCurrent)

	t.Logf("%+v", data)
}

func TestClient_QueryPIGS(t *testing.T) {
	port := testutils.NewDummySerial()
	port.MockCommand("QPIGS",
		testutils.FromHex("5150494753B7A90D"),
		testutils.FromHex("283234342E322035302E30203233302E302035302E302030333638203033343020303034203339352035302E3130203031352030353620303032362030332E38203332302E332030302E30302030303030302030303031303131302030302030302030313233302030313061AC0D"),
	)
	client := NewClient(port)
	data, err := client.QueryPIGS()
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, "GridVoltage", 244.2, data.GridVoltage)
	assertEqual(t, "GridFrequency", 50.0, data.GridFrequency)
	assertEqual(t, "AlternatingCurrentOutputVoltage", 230.0, data.AlternatingCurrentOutputVoltage)
	assertEqual(t, "AlternatingCurrentOutputFrequency", 50.0, data.AlternatingCurrentOutputFrequency)
}

func assertEqual(t *testing.T, field string, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("%s expected %v, got %v", field, expected, actual)
	}
}
