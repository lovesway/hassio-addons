package models

type (
	// DeviceType structure.
	DeviceType struct {
		ID       int
		Name     string
		Commands []Command
	}
)
