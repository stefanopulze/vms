# VMS Core

### Remote Debugging
```bash
# From remote host
socat -d -d TCP-LISTEN:5000,reuseaddr,fork FILE:/dev/ttyUSB1,raw,echo=0,b2400

# From localhost
socat -d -d PTY,link=/tmp/ttyV0,rawer,echo=0 TCP:192.168.1.42:5000
```

## Ideas
### Scenarios
```yaml
scenes:
  - name: Winter
    sourcePriority: sbu
    chargerSourcePriority: solar_utility
    batteryRechargeVoltage: 41
    batteryRedischargeVoltage: 48

  - name: Fog
    sourcePriority: usb
    chargerSourcePriority: solar_first

  - name: Summer
    sourcePriority: sbu
    chargerSourcePriority: only_solar
```

- [ ] Add support for multiple scenarios

## ClickHouse
```sql
CREATE TABLE inverter(
	timestamp DateTime,
	grid_voltage Float32,
	grid_frequency Float32,
	ac_output_voltage Float32,
	ac_output_frequency Float32,
	ac_output_apparent_power Int32,
	ac_output_active_power Int32,
	output_load_percent Int32,
	bus_voltage Int32,
	battery_voltage Float32,
	battery_charging_current Int32,
	battery_capacity Int32,
	inverter_heat_sink_temperature Int32,
	pv_input_current Float32,
	pv_input_voltage Float32,
	battery_discharge_current Int32,
	battery_discharging_power Float32,
	eeprom_version Int32,
	pv_charging_power Int32,
	battery_charging_power Float32,
	device_mode String
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp);
```

## Database

```sql
CREATE TABLE daily_usage (
    id BIGSERIAL PRIMARY KEY ,
    date date,
    line int,
    home int,
    garage int,
    inverter int,
    solar_production int,
    offgrid_percent int
);

create index daily_usage_date_idx on daily_usage (date);
```