package models

type (
	// Group structure.
	Group struct {
		ID          int
		SceneID     int
		Delay       float32
		GlobalDelay bool
		Order       int
		Actions     []Action
	}
)
