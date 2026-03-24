package humanize

import "strings"

func Mode(mode string) string {
	switch mode {
	case "line_mode":
		return "line"
	case "battery_mode":
		return "battery"
	default:
		return mode
	}
}

func OutputSourceFull(source string) string {
	switch source {
	case "usb":
		return "Utility SolarProduction Battery"
	case "sub":
		return "SolarProduction Utility Battery"
	case "sbu":
		return "SolarProduction Battery Utility"
	default:
		return strings.ToUpper(source)
	}
}
