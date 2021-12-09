package models

type (
	// Cycle structure.
	Cycle struct {
		ID               int
		ShowID           int
		SceneID          int
		SceneCycles      int
		EndDelay         float32
		LoopInclude      bool
		GlobalDelay      float32
		GlobalSpeed      int
		GlobalParameter1 string
		GlobalParameter2 string
		Scene            Scene
	}
)
