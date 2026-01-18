package event

type Type string

const (
	UpdateFlags       Type = "update_flags"
	ReadGeneralStatus Type = "read_general_status"
	ReadRatingInfo    Type = "read_rating_info"
	ReadMode          Type = "read_mode"
	ReadWarnings      Type = "read_warnings"
)
