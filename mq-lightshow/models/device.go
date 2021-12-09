package models

type (
	// Device structure.
	Device struct {
		ID       int
		Name     string
		Topic    string
		Type     DeviceType
		Selected bool
	}
)
