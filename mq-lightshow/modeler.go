package main

import (
	"github.com/lovesway/hassio-addons/mq-lightshow/database"
	"github.com/lovesway/hassio-addons/mq-lightshow/models"
)

// Modeler struct to represent a class.
type Modeler struct {
	db *database.Sqlite
}

// NewModler method to instantiate class/struct.
func NewModler(db *database.Sqlite) Modeler {
	return Modeler{
		db: db,
	}
}

// GetShows to return a slice of Show objects.
func (md *Modeler) GetShows() ([]models.Show, error) {
	shows, err := md.db.GetShows()
	if err != nil {
		return []models.Show{}, err
	}

	for i, show := range shows {
		s := &shows[i]
		s.Running = ex.IsShowRunning(show.ID)
	}

	return shows, err
}

// GetShowsRecursive to return a slice of Show objects with child objects populated.
func (md *Modeler) GetShowsRecursive() ([]models.Show, error) {
	shows, err := md.GetShows()
	if err != nil {
		return []models.Show{}, err
	}

	for i, show := range shows {
		cycles, err := md.GetShowCyclesRecursive(show.ID)
		if err != nil {
			return []models.Show{}, err
		}

		s := &shows[i]
		s.Cycles = cycles
	}

	return shows, err
}

// GetShow to return a Show.
func (md *Modeler) GetShow(showID int) (models.Show, error) {
	show, err := md.db.GetShow(showID)
	if err != nil {
		return models.Show{}, err
	}

	show.Running = ex.IsShowRunning(show.ID)

	return show, err
}

// GetShowRecursive to return a Show object with child objects populated.
func (md *Modeler) GetShowRecursive(showID int) (models.Show, error) {
	show, err := md.GetShow(showID)
	if err != nil {
		return models.Show{}, err
	}

	cycles, err := md.GetShowCyclesRecursive(show.ID)
	if err != nil {
		return models.Show{}, err
	}

	show.Cycles = cycles

	return show, err
}

// GetShowByTopic to get a show by its mqtt topic.
func (md *Modeler) GetShowByTopic(topic string) (models.Show, error) {
	show, err := md.db.GetShowByTopic(topic)
	if err != nil {
		return models.Show{}, err
	}

	show.Running = ex.IsShowRunning(show.ID)

	return show, err
}

// AddShow to add a Show.
func (md *Modeler) AddShow(show models.Show) (int, error) {
	return md.db.AddShow(show)
}

// DeleteShow to delete a Show.
func (md *Modeler) DeleteShow(showID int) error {
	return md.db.DeleteShow(showID)
}

// SetShow to update a Show.
func (md *Modeler) SetShow(show models.Show) error {
	return md.db.SetShow(show)
}

// GetShowCycles to return a slice of Cycle objects.
func (md *Modeler) GetShowCycles(showID int) ([]models.Cycle, error) {
	return md.db.GetShowCycles(showID)
}

// GetShowCyclesRecursive to return a slice of Cycle objects with child objects populated.
func (md *Modeler) GetShowCyclesRecursive(showID int) ([]models.Cycle, error) {
	cycles, err := md.GetShowCycles(showID)
	if err != nil {
		return []models.Cycle{}, err
	}

	for i, cycle := range cycles {
		scene, err := md.GetSceneRecursive(cycle.SceneID)
		if err != nil {
			return []models.Cycle{}, err
		}

		c := &cycles[i]
		c.Scene = scene
	}

	return cycles, err
}

// GetShowCycle to return a Cycle.
func (md *Modeler) GetShowCycle(cycleID int) (models.Cycle, error) {
	return md.db.GetShowCycle(cycleID)
}

// AddShowCycle to add a Cycle.
func (md *Modeler) AddShowCycle(cycle models.Cycle) (int, error) {
	return md.db.AddShowCycle(cycle)
}

// SetShowCycle to update a Cycle.
func (md *Modeler) SetShowCycle(cycle models.Cycle) error {
	return md.db.SetShowCycle(cycle)
}

// DeleteShowCycle to delete a Cycle.
func (md *Modeler) DeleteShowCycle(cycleID int) error {
	return md.db.DeleteShowCycle(cycleID)
}

// GetScenes to return a slice of Scene objects.
func (md *Modeler) GetScenes() ([]models.Scene, error) {
	return md.db.GetScenes()
}

// GetScenesRecursive to return a slice of Scene objects with recursive objects populated.
func (md *Modeler) GetScenesRecursive() ([]models.Scene, error) {
	scenes, err := md.GetScenes()
	if err != nil {
		return []models.Scene{}, err
	}

	for i, scene := range scenes {
		groups, err := md.GetGroups(scene.ID)
		if err != nil {
			return []models.Scene{}, err
		}

		s := &scenes[i]
		s.Groups = groups
	}

	return scenes, err
}

// GetScene to return a Scene.
func (md *Modeler) GetScene(sceneID int) (models.Scene, error) {
	return md.db.GetScene(sceneID)
}

// GetSceneRecursive to return a Scene object with recursive objects populated.
func (md *Modeler) GetSceneRecursive(sceneID int) (models.Scene, error) {
	scene, err := md.GetScene(sceneID)
	if err != nil {
		return scene, err
	}

	groups, err := md.GetGroupsRecursive(scene.ID)
	if err != nil {
		return scene, err
	}

	scene.Groups = groups

	return scene, err
}

// AddScene to add a Scene.
func (md *Modeler) AddScene(scene models.Scene) (int, error) {
	return md.db.AddScene(scene)
}

// SetScene to update a scene.
func (md *Modeler) SetScene(scene models.Scene) error {
	return md.db.SetScene(scene)
}

// DeleteScene to delete a Scene.
func (md *Modeler) DeleteScene(sceneID int) error {
	return md.db.DeleteScene(sceneID)
}

// GetGroups to return a slice of Group objects.
func (md *Modeler) GetGroups(sceneID int) ([]models.Group, error) {
	return md.db.GetGroups(sceneID)
}

// GetGroupsRecursive to return a slice of Group objects.
func (md *Modeler) GetGroupsRecursive(sceneID int) ([]models.Group, error) {
	groups, err := md.GetGroups(sceneID)
	if err != nil {
		return []models.Group{}, err
	}

	for i, group := range groups {
		g := &groups[i]

		actions, err := md.GetActions(group.ID)
		if err != nil {
			return []models.Group{}, err
		}

		g.Actions = actions
	}

	return groups, err
}

// GetGroup to return a Group object.
func (md *Modeler) GetGroup(groupID int) (models.Group, error) {
	return md.db.GetGroup(groupID)
}

// GetGroupRecursive to return a Group object with recursive objects populated.
func (md *Modeler) GetGroupRecursive(groupID int) (models.Group, error) {
	group, err := md.GetGroup(groupID)
	if err != nil {
		return models.Group{}, err
	}

	actions, err := md.GetActions(group.ID)
	if err != nil {
		return models.Group{}, err
	}

	group.Actions = actions

	return group, err
}

// GetGroupOrderNext to return the next available order.
func (md *Modeler) GetGroupOrderNext(sceneID int) (int, error) {
	return md.db.GetGroupOrderNext(sceneID)
}

// AddGroup to add a group.
func (md *Modeler) AddGroup(group models.Group) (int, error) {
	return md.db.AddGroup(group)
}

// SetGroup to update a Group.
func (md *Modeler) SetGroup(group models.Group) error {
	return md.db.SetGroup(group)
}

// SortGroup to update the sort order of groups.
func (md *Modeler) SortGroup(sceneID int, groupID int, sortID int) error {
	return md.db.SortGroup(sceneID, groupID, sortID)
}

// DeleteGroup function.
func (md *Modeler) DeleteGroup(groupID int) error {
	return md.db.DeleteGroup(groupID)
}

// GetActions to return a slice of Action objects.
func (md *Modeler) GetActions(groupID int) ([]models.Action, error) {
	return md.db.GetActions(groupID)
}

// GetAction to return an Action object.
func (md *Modeler) GetAction(actionID int) (models.Action, error) {
	return md.db.GetAction(actionID)
}

// AddAction to add an Action object.
func (md *Modeler) AddAction(action models.Action) (int, error) {
	return md.db.AddAction(action)
}

// SetAction to update an Action object.
func (md *Modeler) SetAction(action models.Action) error {
	return md.db.SetAction(action)
}

// SortAction to update the sort order of Action objects.
func (md *Modeler) SortAction(groupID int, actionID int, sortID int) error {
	return md.db.SortAction(groupID, actionID, sortID)
}

// GetActionOrderNext to get the next order value.
func (md *Modeler) GetActionOrderNext(groupID int) (int, error) {
	return md.db.GetActionOrderNext(groupID)
}

// DeleteAction to delete an Action object.
func (md *Modeler) DeleteAction(actionID int) error {
	return md.db.DeleteAction(actionID)
}
