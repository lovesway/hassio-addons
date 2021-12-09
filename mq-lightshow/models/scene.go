package models

type (
	// Scene structure.
	Scene struct {
		ID             int
		Name           string
		AllowedDevices []Device
		Groups         []Group
	}
)
