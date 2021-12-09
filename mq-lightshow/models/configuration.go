package models

type (
	// Configuration structure. This all goes into the db with primitive map conversion so strings only.
	Configuration struct {
		MQTTHost string
		MQTTUser string
		MQTTPass string
		LogLevel string
	}
)
