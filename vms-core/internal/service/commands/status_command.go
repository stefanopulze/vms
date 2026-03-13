package commands

import (
	"context"
	"fmt"
	"vms-core/internal/cache"
	"vms-core/internal/humanize"
	"vms-core/internal/infrastructure/telegram"
)

var _ Command = (*StatusCommand)(nil)

func NewStatusCommand(tc *telegram.Client, qs *cache.QuerySnapshot) *StatusCommand {
	return &StatusCommand{
		telegram:      tc,
		querySnapshot: qs,
	}
}

type StatusCommand struct {
	telegram      *telegram.Client
	querySnapshot *cache.QuerySnapshot
}

func (s StatusCommand) GetPattern() string {
	return "/status"
}

func (s StatusCommand) StartNewSession() (int64, error) {
	mode := s.querySnapshot.GetMode()
	piri := s.querySnapshot.GetRatingInfo()
	pigs := s.querySnapshot.GetGeneralStatus()

	var message string
	if pigs == nil || piri == nil {
		message = "No data available. Please wait a few seconds."
	} else {
		message = fmt.Sprintf("Mode: %s\n", humanize.Mode(mode))
		message += fmt.Sprintf("Output Source Priority: %s\n", humanize.OutputSourceFull(piri.OutputSourcePriorityEnum()))
		message += fmt.Sprintf("Solar production: %dw\n", pigs.PhotovoltaicChargingPower)
		message += fmt.Sprintf("Battery")
		message += fmt.Sprintf("\tCapacity: %d%%\n", pigs.BatteryCapacity)
		message += fmt.Sprintf("\tCharging: %.2fw\n", pigs.BatteryChargingPower)
		message += fmt.Sprintf("\tDischarging: %.2fw", pigs.BatteryDischargingPower)
	}

	return s.telegram.Send(context.Background(), message)
}

func (s StatusCommand) HandleCallback(_ *telegram.CallbackQuery) error {
	return nil
}

func (s StatusCommand) NeedCallback() bool {
	return false
}
