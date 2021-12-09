package models

// Show structure.
type Show struct {
	Running          bool
	Repeat           bool
	GlobalDelay      float32
	ID               int
	GlobalSpeed      int
	Name             string
	Topic            string
	GlobalParameter1 string
	GlobalParameter2 string
	Cycles           []Cycle
}
