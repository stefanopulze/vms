package model

import "vms-core/internal/voltronic"

type UpdateMode struct {
	Mode string `json:"mode"`
}

type UpdateSourcePriority struct {
	Source string `json:"source"`
}

type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

type ApiRatingInfo struct {
	*voltronic.DeviceRatingInfo
	SourcePriorityEnum        string `json:"sourcePriorityEnum"`
	ChargerSourcePriorityEnum string `json:"chargerSourcePriorityEnum"`
}
