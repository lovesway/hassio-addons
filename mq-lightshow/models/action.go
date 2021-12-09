// Package models provides our models
package models

type (
	// Action structure.
	Action struct {
		ID              int
		GroupID         int
		Devices         []Device
		Command         string
		Parameter       string
		GlobalParameter string
		Order           int
	}
)
