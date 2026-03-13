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
		return "Utility Solar Battery"
	case "sub":
		return "Solar Utility Battery"
	case "sbu":
		return "Solar Battery Utility"
	default:
		return strings.ToUpper(source)
	}
}
