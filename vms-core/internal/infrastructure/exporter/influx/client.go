package influx

import (
	"context"
	"time"
	"vms-core/internal/infrastructure/exporter"
	"vms-core/internal/voltronic"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
)

var _ exporter.Client = (*Client)(nil)

func NewClient(opts Options) (*Client, error) {
	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:         opts.Host,
		Token:        opts.Token,
		Database:     opts.Database,
		Organization: opts.Organization,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

type Client struct {
	client *influxdb3.Client
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Name() string {
	return "influx"
}

func (c *Client) GeneralStatus(pigs *voltronic.DeviceGeneralStatus, mode string) error {
	tags := map[string]string{}
	fields := map[string]interface{}{}

	fields["grid_voltage"] = pigs.GridVoltage
	fields["grid_frequency"] = pigs.GridFrequency
	fields["ac_output_voltage"] = pigs.AlternatingCurrentOutputVoltage
	fields["ac_output_frequency"] = pigs.AlternatingCurrentOutputFrequency
	fields["ac_output_apparent_power"] = pigs.AlternatingCurrentOutputApparentPower
	fields["ac_output_active_power"] = pigs.AlternatingCurrentOutputActivePower
	fields["output_load_percent"] = pigs.OutputLoadPercent
	fields["bus_voltage"] = pigs.BusVoltage
	fields["battery_voltage"] = pigs.BatteryVoltage
	fields["battery_charging_current"] = pigs.BatteryChargingCurrent
	fields["battery_capacity"] = pigs.BatteryCapacity
	fields["inverter_heat_sink_temperature"] = pigs.InverterHeatSinkTemperature
	fields["pv_input_current"] = pigs.PhotovoltaicInputCurrent
	fields["pv_input_voltage"] = pigs.PhotovoltaicInputVoltage
	//fields["battery_voltage_from_solar_chargeController"] = pigs.BatteryVoltageFromSolarChargeController
	fields["battery_discharge_current"] = pigs.BatteryDischargeCurrent
	fields["battery_discharging_power"] = pigs.BatteryDischargingPower
	//fields["battery_voltage_offsetForFansOn"] = pigs.BatteryVoltageOffsetForFansOn
	fields["eeprom_version"] = pigs.EepromVersion
	fields["pv_charging_power"] = pigs.PhotovoltaicChargingPower
	fields["battery_charging_power"] = pigs.BatteryChargingPower
	fields["device_mode"] = mode

	point := influxdb3.NewPoint("inverter", tags, fields, pigs.Timestamp)

	return c.client.WritePoints(context.Background(), []*influxdb3.Point{point})
}

func (c *Client) Query(ctx context.Context, query string, params map[string]any) ([]map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	iterator, err := c.client.QueryWithParameters(ctx, query, params)
	if err != nil {
		return nil, err
	}

	values := make([]map[string]any, 0)
	for iterator.Next() {
		value := iterator.Value()
		values = append(values, value)
	}

	return values, nil
}

func (c *Client) WritePoints(ctx context.Context, points ...*influxdb3.Point) error {
	return c.client.WritePoints(ctx, points)
}
