package models

type (
	// Message from mqtt.
	Message struct {
		Topic   string
		Message string
	}
)
