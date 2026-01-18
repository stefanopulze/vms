package voltronic

import (
	"fmt"
	"log/slog"
	"reflect"
	"time"
)

func (c *Client) UpdateTime(t time.Time) error {
	cmd := "DAT" + t.Format("060102150405")
	logger.Debug(fmt.Sprintf("Updating time to: %s", cmd))
	return c.SendUpdateCommand(cmd)
}

func (c *Client) UpdateFlags(flags DeviceFlags) error {
	enabled := ""
	disabled := ""

	rf := reflect.TypeOf(flags)
	vf := reflect.ValueOf(flags)
	for i := 0; i < rf.NumField(); i++ {
		field := rf.Field(i)
		flagName, ok := field.Tag.Lookup("flag")
		if !ok {
			continue
		}

		if vf.Field(i).Bool() {
			enabled += flagName
		} else {
			disabled += flagName
		}
	}

	if len(enabled) > 0 {
		logger.Debug(fmt.Sprintf("Enabling flags: %s", enabled))
		if err := c.SendUpdateCommand("PE" + enabled); err != nil {
			slog.Error(fmt.Sprintf("Failed to enable flags: %s", enabled))
			return err
		}
	}

	if len(disabled) > 0 {
		logger.Debug(fmt.Sprintf("Disabling flags: %s", enabled))
		if err := c.SendUpdateCommand("PD" + disabled); err != nil {
			slog.Error(fmt.Sprintf("Failed to disable flags: %s", enabled))
			return err
		}
	}

	return nil
}

func (c *Client) UpdateSourcePriority(v string) error {
	var mode string
	switch v {
	case "usb":
		mode = "00"
	case "sub":
		mode = "01"
	case "sbu":
		mode = "02"
	default:
		return fmt.Errorf("unknown source priority mode: %s", v)
	}

	logger.Debug(fmt.Sprintf("Updating source priority: %s", mode))
	return c.SendUpdateCommand("POP" + mode)
}

func (c *Client) UpdateBatteryRechargeVoltage(v float32) error {
	cmd := fmt.Sprintf("PBCV%0.1f", v)
	logger.Debug(fmt.Sprintf("Updating battery recharge voltage: %s", cmd))
	return c.SendUpdateCommand(cmd)
}

func (c *Client) UpdateBatteryRedischargeVoltage(v float32) error {
	cmd := fmt.Sprintf("PBDV%0.1f", v)
	logger.Debug(fmt.Sprintf("Updating battery redischarge voltage: %s", cmd))
	return c.SendUpdateCommand(cmd)
}

func (c *Client) UpdateChargerPriority(v string) error {
	var mode string
	switch v {
	case "solar_first":
		mode = "01"
	case "solar_utility":
		mode = "02"
	case "only_solar":
		mode = "03"
	default:
		return fmt.Errorf("unknown charger priority mode: %s", v)
	}

	cmd := fmt.Sprintf("PCP%s", mode)
	logger.Debug(fmt.Sprintf("Updating charger priority: %s", cmd))
	return c.SendUpdateCommand(cmd)
}

// UpdateLedUsage Enable/disable LED function
func (c *Client) UpdateLedUsage(enable bool) error {
	v := 0
	if enable {
		v = 1
	}
	cmd := fmt.Sprintf("PLEDE%d", v)
	logger.Debug(fmt.Sprintf("Updating led usage: %s", cmd))
	return c.SendUpdateCommand(cmd)
}

func (c *Client) UpdateWorkingMode(v string) error {
	var mode string
	switch v {
	case "appliance":
		mode = "00"
	case "ups":
		mode = "01"
	default:
		return fmt.Errorf("unknown working mode: %s", v)
	}

	cmd := fmt.Sprintf("PGR%s", mode)
	logger.Debug(fmt.Sprintf("Updating working mode: %s", cmd))
	return c.SendUpdateCommand(cmd)
}
