package clickhouse

import (
	"context"
	"vms-core/internal/voltronic"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const pigsInsertQuery = `INSERT INTO inverter(
timestamp,
grid_voltage,
grid_frequency,
ac_output_voltage,
ac_output_frequency,
ac_output_apparent_power,
ac_output_active_power,
output_load_percent,
bus_voltage,
battery_voltage,
battery_charging_current,
battery_capacity,
inverter_heat_sink_temperature,
pv_input_current,
pv_input_voltage,
battery_discharge_current,
battery_discharging_power,
eeprom_version,
pv_charging_power,
battery_charging_power,
device_mode
) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`

func NewClient(opts Options) (*Client, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{opts.Addr},
		Auth: clickhouse.Auth{
			Database: opts.Database,
			Username: opts.Username,
			Password: opts.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

type Client struct {
	conn driver.Conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Name() string {
	return "clickHouse"
}

func (c *Client) GeneralStatus(pigs *voltronic.DeviceGeneralStatus, mode string) error {
	return c.conn.Exec(context.Background(), pigsInsertQuery,
		pigs.Timestamp,
		pigs.GridVoltage,
		pigs.GridFrequency,
		pigs.AlternatingCurrentOutputVoltage,
		pigs.AlternatingCurrentOutputFrequency,
		pigs.AlternatingCurrentOutputApparentPower,
		pigs.AlternatingCurrentOutputActivePower,
		pigs.OutputLoadPercent,
		pigs.BusVoltage,
		pigs.BatteryVoltage,
		pigs.BatteryChargingCurrent,
		pigs.BatteryCapacity,
		pigs.InverterHeatSinkTemperature,
		pigs.PhotovoltaicInputCurrent,
		pigs.PhotovoltaicInputVoltage,
		pigs.BatteryDischargeCurrent,
		pigs.BatteryDischargingPower,
		pigs.EepromVersion,
		pigs.PhotovoltaicChargingPower,
		pigs.BatteryChargingPower,
		mode)
}
