package utils

func ModeToHuman(mode string) string {
	switch mode {
	case "line_mode":
		return "line"
	case "battery_mode":
		return "battery"
	default:
		return mode
	}
}
