package main

import (
	"fmt"
	"strconv"
	"time"

	"gitlab.local/hassio-addons/mq-lightshow/models"
)

// Executor represents the controller for the UI.
type Executor struct {
	md    Modeler
	mq    *MQController
	gblsZ globals // globals with zero value for usage in comparisons.
}

// Shows tracks all instances of running shows.
var Shows []Running

// NewExecutor provides an instance of Executor.
func NewExecutor(md Modeler, mq *MQController) *Executor {
	return &Executor{
		md: md,
		mq: mq,
	}
}

// ExecuteActionByID to send an action to MQTT.
func (e Executor) ExecuteActionByID(actionID int) {
	action, err := e.md.GetAction(actionID)
	if err != nil {
		log.Error(err.Error())

		return
	}

	for _, device := range action.Devices {
		e.ExecuteAction(device.Topic, action.Command, action.Parameter)
	}
}

// ExecuteAction to send an action to MQTT.
func (e Executor) ExecuteAction(topic string, command string, parameter string) {
	e.mq.SendAction(topic, command, parameter)
}

func (e Executor) waitForSeconds(secondsFloat float32) {
	const oneThousand = 1000
	duration := time.Duration(int64(secondsFloat * oneThousand))

	delay := duration * time.Millisecond

	time.Sleep(delay)
}

// ExecuteActionGroupByID to send a group of actions to MQTT.
func (e Executor) ExecuteActionGroupByID(groupID int) {
	log.Infof("ExecuteActionGroupByID: %v", groupID)

	g, err := e.md.GetGroup(groupID)
	if err != nil {
		log.Error(err.Error())

		return
	}

	actions, err := e.md.GetActions(groupID)
	if err != nil {
		log.Error(err.Error())

		return
	}

	for _, a := range actions {
		e.ExecuteActionByID(a.ID)
	}

	e.waitForSeconds(g.Delay)
}

// ExecuteSceneByID to send a Scene to MQTT.
func (e Executor) ExecuteSceneByID(sceneID int) {
	log.Infof("ExecuteSceneByID: %v", sceneID)

	gs, err := e.md.GetGroups(sceneID)
	if err != nil {
		log.Error(err.Error())

		return
	}

	for _, g := range gs {
		e.ExecuteActionGroupByID(g.ID)
	}
}

type globals struct {
	Delay      float32
	Speed      int
	Parameter1 string
	Parameter2 string
}

func (e Executor) runShow(show models.Show) {
	// check for show globals.
	var gbls globals
	if show.GlobalDelay != e.gblsZ.Delay {
		gbls.Delay = show.GlobalDelay
	}

	if show.GlobalSpeed != e.gblsZ.Speed {
		gbls.Speed = show.GlobalSpeed
	}

	if show.GlobalParameter1 != e.gblsZ.Parameter1 {
		gbls.Parameter1 = show.GlobalParameter1
	}

	if show.GlobalParameter2 != e.gblsZ.Parameter2 {
		gbls.Parameter2 = show.GlobalParameter2
	}

	looping := false

	for {
		for _, cycle := range show.Cycles {
			if !e.IsShowRunning(show.ID) {
				return
			}

			if looping && !cycle.LoopInclude {
				continue
			}

			// check for cycle globals.
			cgbls := gbls
			if cycle.GlobalDelay != e.gblsZ.Delay {
				cgbls.Delay = cycle.GlobalDelay
			}

			if cycle.GlobalSpeed != e.gblsZ.Speed {
				cgbls.Speed = cycle.GlobalSpeed
			}

			if cycle.GlobalParameter1 != e.gblsZ.Parameter1 {
				cgbls.Parameter1 = cycle.GlobalParameter1
			}

			if cycle.GlobalParameter2 != e.gblsZ.Parameter2 {
				cgbls.Parameter2 = cycle.GlobalParameter2
			}

			for i := 1; i <= cycle.SceneCycles; i++ {
				e.runScene(show.ID, cgbls, cycle.Scene)

				if !e.IsShowRunning(show.ID) {
					return
				}
			}

			if cycle.EndDelay > 0 {
				time.Sleep(time.Duration(cycle.EndDelay) * time.Second)
			}
		}

		if !show.Repeat {
			break
		}

		looping = true
	}

	if err := e.StopShow(show.ID); err != nil {
		log.Error(err)
	}
}

func (e Executor) runScene(showID int, gbls globals, scene models.Scene) {
	for _, group := range scene.Groups {
		if !e.IsShowRunning(showID) {
			return
		}

		e.runActionGroup(showID, gbls, group)
	}
}

func (e Executor) runActionGroup(showID int, gbls globals, group models.Group) {
	for _, action := range group.Actions {
		e.runAction(showID, gbls, action)

		if !e.IsShowRunning(showID) {
			return
		}
	}

	if group.GlobalDelay && gbls.Delay != e.gblsZ.Delay {
		e.waitForSeconds(gbls.Delay)
	} else {
		e.waitForSeconds(group.Delay)
	}
}

func (e Executor) runAction(showID int, gbls globals, action models.Action) {
	for _, d := range action.Devices {
		if !e.IsShowRunning(showID) {
			return
		}

		// determine whether or not to use global parameters.
		var parameter string
		if action.GlobalParameter == globalSpeed && gbls.Speed != e.gblsZ.Speed {
			parameter = strconv.Itoa(gbls.Speed)
		} else if action.GlobalParameter == globalParameter1 && gbls.Parameter1 != e.gblsZ.Parameter1 {
			parameter = gbls.Parameter1
		} else if action.GlobalParameter == globalParameter2 && gbls.Parameter2 != e.gblsZ.Parameter2 {
			parameter = gbls.Parameter2
		} else {
			parameter = action.Parameter
		}

		e.ExecuteAction(d.Topic, action.Command, parameter)
	}
}

// Running tracks an instance of a running show.
type Running struct {
	ShowID int
}

// IsShowRunning to determine whether or not a show is running.
func (e Executor) IsShowRunning(showID int) bool {
	for _, running := range Shows {
		if showID == running.ShowID {
			return true
		}
	}

	return false
}

// StartShow to run a tracked show which can be stopped.
func (e Executor) StartShow(showID int) error {
	var err error

	if e.IsShowRunning(showID) {
		return fmt.Errorf("Show already running for showID: %v", showID)
	}

	show, err := e.md.GetShowRecursive(showID)
	if err != nil {
		return err
	}

	log.Infof("Starting Show: %v", show.Name)

	go func() {
		e.runShow(show)
	}()

	Shows = append(Shows, Running{ShowID: showID})

	e.mq.SendShowState(show.Topic, "ON")

	return err
}

// StopShow to stop a running show.
func (e Executor) StopShow(showID int) error {
	var err error

	show, err := e.md.GetShow(showID)
	if err != nil {
		return err
	}

	e.removeShowByIndex(e.getRunningShowIndex(showID))

	log.Infof("Stopping Show: %v", show.Name)

	e.mq.SendShowState(show.Topic, "OFF")

	return err
}

func (e Executor) removeShowByIndex(i int) {
	if !e.showIndexExists(i) {
		log.Error("no show exists for index %v", i)

		return
	}

	Shows[i] = Shows[len(Shows)-1]
	Shows = Shows[:len(Shows)-1]
}

func (e Executor) getRunningShowIndex(showID int) int {
	for i, show := range Shows {
		if showID == show.ShowID {
			return i
		}
	}

	return 0
}

func (e Executor) showIndexExists(index int) bool {
	for i := range Shows {
		if i == index {
			return true
		}
	}

	return false
}
