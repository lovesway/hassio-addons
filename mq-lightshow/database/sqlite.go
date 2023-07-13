// Package database provides our database.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	// Sqlite3 for sql package.
	"github.com/lovesway/hassio-addons/mq-lightshow/devicetypes"
	"github.com/lovesway/hassio-addons/mq-lightshow/models"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// Sqlite struct to represent a class.
type Sqlite struct {
	db          *sql.DB
	deviceTypes []models.DeviceType
}

// NewSqlite method to instantiate class/struct.
func NewSqlite(dt devicetypes.DeviceTypes) *Sqlite {
	database := &Sqlite{
		deviceTypes: dt.GetDeviceTypes(),
	}

	return database
}

// SetDeviceTypes to set deviceTypes in db instance.
func (sl *Sqlite) SetDeviceTypes(t []models.DeviceType) {
	sl.deviceTypes = t
}

// global log adapter to support zap sugar.
var log *zap.SugaredLogger

// SetLogger to set log in db instance.
func SetLogger(l *zap.SugaredLogger) {
	log = l
}

// GetShows to return a slice of Show structs.
func (sl *Sqlite) GetShows() ([]models.Show, error) {
	shows := []models.Show{}

	sqlStmt := "SELECT show_id, name, topic, repeat, " +
		"global_delay, global_speed, global_parameter1, global_parameter2 FROM shows"

	rows, err := sl.db.Query(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return shows, err
	}

	defer func() {
		if err = rows.Err(); err != nil {
			log.Errorf("Sqlite GetShows: %v.", err)
		}

		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	for rows.Next() {
		var showID int

		var name string

		var topic string

		var repeat bool

		var globalDelay float32

		var globalSpeed int

		var globalParameter1 string

		var globalParameter2 string

		err = rows.Scan(&showID, &name, &topic, &repeat, &globalDelay, &globalSpeed, &globalParameter1, &globalParameter2)
		if err != nil {
			log.Errorf("%q: %s", err, sqlStmt)

			return shows, err
		}

		ts := models.Show{
			ID:               showID,
			Name:             name,
			Topic:            topic,
			Repeat:           repeat,
			GlobalDelay:      globalDelay,
			GlobalSpeed:      globalSpeed,
			GlobalParameter1: globalParameter1,
			GlobalParameter2: globalParameter2,
		}
		shows = append(shows, ts)
	}

	return shows, err
}

// GetShow to return a single Show struct.
func (sl *Sqlite) GetShow(showID int) (models.Show, error) {
	sqlStmt := fmt.Sprintf("SELECT name, topic, repeat, "+
		"global_delay, global_speed, global_parameter1, global_parameter2 FROM shows where show_id = '%v'", showID)

	var name, topic, globalParameter1, globalParameter2 string

	var repeat bool

	var globalDelay float32

	var globalSpeed int

	err := sl.db.QueryRow(sqlStmt).Scan(
		&name, &topic, &repeat, &globalDelay, &globalSpeed, &globalParameter1, &globalParameter2,
	)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return models.Show{}, err
	}

	show := models.Show{
		ID:               showID,
		Name:             name,
		Topic:            topic,
		Repeat:           repeat,
		GlobalDelay:      globalDelay,
		GlobalSpeed:      globalSpeed,
		GlobalParameter1: globalParameter1,
		GlobalParameter2: globalParameter2,
	}

	return show, err
}

// GetShowByTopic to return a single Show struct.
func (sl *Sqlite) GetShowByTopic(topic string) (models.Show, error) {
	sqlStmt := fmt.Sprintf("SELECT show_id, name, repeat, "+
		"global_delay, global_speed, global_parameter1, global_parameter2 FROM shows where topic='%v'", topic)

	var showID, globalSpeed int

	var name, globalParameter1, globalParameter2 string

	var repeat bool

	var globalDelay float32

	err := sl.db.QueryRow(sqlStmt).Scan(
		&showID, &name, &repeat, &globalDelay, &globalSpeed, &globalParameter1, &globalParameter2,
	)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return models.Show{}, err
	}

	show := models.Show{
		ID:               showID,
		Name:             name,
		Topic:            topic,
		Repeat:           repeat,
		GlobalDelay:      globalDelay,
		GlobalSpeed:      globalSpeed,
		GlobalParameter1: globalParameter1,
		GlobalParameter2: globalParameter2,
	}

	return show, err
}

// AddShow to db.
func (sl *Sqlite) AddShow(s models.Show) (int, error) {
	sqlStmt := fmt.Sprintf(
		"INSERT INTO shows(name, topic, repeat, global_delay, global_speed, global_parameter1, global_parameter2)"+
			" values('%s', '%s', '%t', '%v', '%v', '%s', '%s')",
		s.Name, s.Topic, s.Repeat, s.GlobalDelay, s.GlobalSpeed, s.GlobalParameter1, s.GlobalParameter2,
	)

	res, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return 0, err
	}

	id, _ := res.LastInsertId()

	return int(id), err
}

// SetShow to update a Show.
func (sl *Sqlite) SetShow(s models.Show) error {
	sqlStmt := fmt.Sprintf(
		"UPDATE shows set name='%s', topic='%s', repeat='%t', "+
			"global_delay='%v', global_speed='%v', global_parameter1='%s', global_parameter2='%s' where show_id='%v'",
		s.Name, s.Topic, s.Repeat, s.GlobalDelay, s.GlobalSpeed, s.GlobalParameter1, s.GlobalParameter2, s.ID,
	)

	_, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// DeleteShow function.
func (sl *Sqlite) DeleteShow(showID int) error {
	sqlStmt := fmt.Sprintf("DELETE from shows_cycles where show_id=%v", showID)

	_, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return err
	}

	sqlStmt = fmt.Sprintf("DELETE from shows where show_id=%v", showID)

	_, err = sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// GetShowCycles to return a slice of Cycle structs.
func (sl *Sqlite) GetShowCycles(showID int) ([]models.Cycle, error) {
	cycles := []models.Cycle{}

	rows, err := sl.db.Query("SELECT cycle_id, show_id, scene_id, cycles, end_delay, loop_include, "+
		"global_delay, global_speed, global_parameter1, global_parameter2 FROM shows_cycles where show_id = ?", showID)
	if err != nil {
		log.Error(err)

		return cycles, err
	}

	defer func() {
		if err = rows.Err(); err != nil {
			log.Errorf("Sqlite GetShowCycles: %v.", err)
		}

		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	for rows.Next() {
		var cycleID, showID, sceneID, cyclesVal, globalSpeed int

		var endDelay, globalDelay float32

		var loopInclude bool

		var globalParameter1, globalParameter2 string

		err = rows.Scan(
			&cycleID, &showID, &sceneID, &cyclesVal, &endDelay, &loopInclude,
			&globalDelay, &globalSpeed, &globalParameter1, &globalParameter2,
		)
		if err != nil {
			log.Error(err)

			return cycles, err
		}

		cycle := models.Cycle{
			ID:               cycleID,
			ShowID:           showID,
			SceneID:          sceneID,
			SceneCycles:      cyclesVal,
			EndDelay:         endDelay,
			LoopInclude:      loopInclude,
			GlobalDelay:      globalDelay,
			GlobalSpeed:      globalSpeed,
			GlobalParameter1: globalParameter1,
			GlobalParameter2: globalParameter2,
		}
		cycles = append(cycles, cycle)
	}

	return cycles, err
}

// GetShowCycle to return a single Show struct.
func (sl *Sqlite) GetShowCycle(cycleID int) (models.Cycle, error) {
	sqlStmt := fmt.Sprintf("SELECT show_id, scene_id, cycles, end_delay, loop_include, "+
		"global_delay, global_speed, global_parameter1, global_parameter2 FROM shows_cycles where cycle_id = '%v'", cycleID)

	var showID, sceneID, cycles, globalSpeed int

	var endDelay, globalDelay float32

	var loopInclude bool

	var globalParameter1, globalParameter2 string

	err := sl.db.QueryRow(sqlStmt).Scan(
		&showID, &sceneID, &cycles, &endDelay, &loopInclude, &globalDelay, &globalSpeed, &globalParameter1, &globalParameter2,
	)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return models.Cycle{}, err
	}

	cycle := models.Cycle{
		ID:               cycleID,
		ShowID:           showID,
		SceneID:          sceneID,
		SceneCycles:      cycles,
		EndDelay:         endDelay,
		LoopInclude:      loopInclude,
		GlobalDelay:      globalDelay,
		GlobalSpeed:      globalSpeed,
		GlobalParameter1: globalParameter1,
		GlobalParameter2: globalParameter2,
	}

	return cycle, err
}

// AddShowCycle to db.
func (sl *Sqlite) AddShowCycle(c models.Cycle) (int, error) {
	sqlStmt := fmt.Sprintf(
		"INSERT INTO shows_cycles(show_id, scene_id, cycles, end_delay, loop_include, global_delay, global_speed, "+
			"global_parameter1, global_parameter2) values('%v', '%v', '%v', '%v', '%t', '%v', '%v', '%s', '%s')",
		c.ShowID, c.SceneID, c.SceneCycles, c.EndDelay, c.LoopInclude,
		c.GlobalDelay, c.GlobalSpeed, c.GlobalParameter1, c.GlobalParameter2,
	)

	res, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

// SetShowCycle to update a Cycle.
func (sl *Sqlite) SetShowCycle(c models.Cycle) error {
	sqlStmt := fmt.Sprintf("UPDATE shows_cycles set show_id='%v', scene_id='%v', "+
		"cycles='%v', end_delay='%v', loop_include='%t', global_delay='%v', global_speed='%v', "+
		"global_parameter1='%s', global_parameter2='%s' where cycle_id='%v'",
		c.ShowID, c.SceneID, c.SceneCycles, c.EndDelay, c.LoopInclude, c.GlobalDelay,
		c.GlobalSpeed, c.GlobalParameter1, c.GlobalParameter2, c.ID,
	)

	_, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// DeleteShowCycle function.
func (sl *Sqlite) DeleteShowCycle(cycleID int) error {
	sqlStmt := fmt.Sprintf("DELETE from shows_cycles where cycle_id=%v", cycleID)

	_, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// GetScenes to return a slice of Scene structs.
func (sl *Sqlite) GetScenes() ([]models.Scene, error) {
	s := []models.Scene{}

	sqlStmt := "SELECT scene_id, name, allowed_devices FROM scenes"

	rows, err := sl.db.Query(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return s, err
	}

	defer func() {
		if err = rows.Err(); err != nil {
			log.Errorf("Sqlite GetScenes: %v.", err)
		}

		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	for rows.Next() {
		var sceneID int

		var name, allowedDevicesString string

		err = rows.Scan(&sceneID, &name, &allowedDevicesString)
		if err != nil {
			log.Errorf("%q: %s", err, sqlStmt)

			return s, err
		}

		allowedDevices := []models.Device{}

		ds := strings.Split(allowedDevicesString, " ")
		for _, v := range ds {
			if v != "" {
				tID, err := strconv.Atoi(v)
				if err != nil {
					log.Errorf("%q: %s", err, sqlStmt)

					continue
				}

				allowedDevices = append(allowedDevices, sl.GetDevice(tID))
			}
		}

		s = append(s, models.Scene{ID: sceneID, Name: name, AllowedDevices: allowedDevices})
	}

	return s, err
}

// GetScene to return a single Scene struct.
func (sl *Sqlite) GetScene(sceneID int) (models.Scene, error) {
	sqlStmt := fmt.Sprintf("SELECT name, allowed_devices FROM scenes where scene_id = '%v'", sceneID)

	var name, allowedDevicesString string

	err := sl.db.QueryRow(sqlStmt).Scan(&name, &allowedDevicesString)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return models.Scene{}, err
	}

	allowedDevices := []models.Device{}

	ds := strings.Split(allowedDevicesString, " ")
	for _, v := range ds {
		if v == "" {
			continue
		}

		tID, err := strconv.Atoi(v)
		if err != nil {
			log.Errorf("%q: %s", err, sqlStmt)

			return models.Scene{}, err
		}

		allowedDevices = append(allowedDevices, sl.GetDevice(tID))
	}

	scene := models.Scene{
		ID:             sceneID,
		Name:           name,
		AllowedDevices: allowedDevices,
	}

	return scene, err
}

// AddScene to db.
func (sl *Sqlite) AddScene(s models.Scene) (int, error) {
	var allowedDevices string

	for _, d := range s.AllowedDevices {
		sid := strconv.Itoa(d.ID)
		if allowedDevices == "" {
			allowedDevices = sid
		} else {
			allowedDevices = allowedDevices + " " + sid
		}
	}

	sqlStmt := fmt.Sprintf("INSERT INTO scenes(name, allowed_devices) values('%s', '%s')", s.Name, allowedDevices)

	res, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

// SetScene to update a scene.
func (sl *Sqlite) SetScene(s models.Scene) error {
	var allowedDevices string

	for _, d := range s.AllowedDevices {
		sid := strconv.Itoa(d.ID)
		if allowedDevices == "" {
			allowedDevices = sid
		} else {
			allowedDevices = allowedDevices + " " + sid
		}
	}

	sqlStmt := fmt.Sprintf(
		"UPDATE scenes set name='%s', allowed_devices='%s' where scene_id='%v'",
		s.Name, allowedDevices, s.ID,
	)

	_, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// DeleteScene to delete a Scene.
func (sl *Sqlite) DeleteScene(sceneID int) error {
	scene, err := sl.GetScene(sceneID)
	if err != nil {
		return err
	}

	for _, group := range scene.Groups {
		for _, action := range group.Actions {
			err := sl.DeleteAction(action.ID)
			if err != nil {
				log.Error(err)
			}
		}

		err := sl.DeleteGroup(group.ID)
		if err != nil {
			log.Error(err)
		}
	}

	sqlStmt := fmt.Sprintf("DELETE from scenes where scene_id=%v", sceneID)

	_, err = sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// GetGroups to return a slice of Group structs.
func (sl *Sqlite) GetGroups(sceneID int) ([]models.Group, error) {
	s := []models.Group{}

	rows, err := sl.db.Query(
		"SELECT group_id, delay, global_delay, `order` FROM scenes_group where scene_id=? ORDER BY `order`",
		sceneID,
	)
	if err != nil {
		log.Error(err)

		return s, err
	}

	defer func() {
		if err = rows.Err(); err != nil {
			log.Errorf("Sqlite GetGroups: %v.", err)
		}

		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	for rows.Next() {
		var groupID, order int

		var delay float32

		var globalDelay bool

		err = rows.Scan(&groupID, &delay, &globalDelay, &order)
		if err != nil {
			log.Error(err)
		}

		sr := models.Group{
			ID:          groupID,
			SceneID:     sceneID,
			Delay:       delay,
			GlobalDelay: globalDelay,
			Order:       order,
		}
		s = append(s, sr)
	}

	return s, err
}

// GetGroup to return a single Group struct.
func (sl *Sqlite) GetGroup(groupID int) (models.Group, error) {
	sqlStmt := fmt.Sprintf(
		"SELECT scene_id, delay, global_delay, `order` FROM scenes_group where group_id = '%v'",
		groupID,
	)

	var sceneID, order int

	var delay float32

	var globalDelay bool

	err := sl.db.QueryRow(sqlStmt).Scan(&sceneID, &delay, &globalDelay, &order)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return models.Group{}, err
	}

	group := models.Group{
		ID:          groupID,
		SceneID:     sceneID,
		Delay:       delay,
		GlobalDelay: globalDelay,
		Order:       order,
	}

	return group, err
}

// GetGroupOrderNext to get the next order value.
func (sl *Sqlite) GetGroupOrderNext(sceneID int) (int, error) {
	var count int

	row := sl.db.QueryRow("SELECT count(*) FROM scenes_group where scene_id = ?", sceneID)

	err := row.Scan(&count)
	if err != nil {
		log.Error(err)
	} else {
		count++
	}

	return count, err
}

// AddGroup to add a group.
func (sl *Sqlite) AddGroup(g models.Group) (int, error) {
	sqlStmt := fmt.Sprintf(
		"INSERT INTO scenes_group(scene_id, delay, global_delay, 'order') "+
			"values('%v', '%v', '%v', '%v')",
		g.SceneID, g.Delay, g.GlobalDelay, g.Order,
	)

	res, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Error(err.Error())

		return 0, err
	}

	return int(id), err
}

// SetGroup to update a Group.
func (sl *Sqlite) SetGroup(g models.Group) error {
	sqlStmt := fmt.Sprintf(
		"UPDATE scenes_group set delay='%v', global_delay='%v', `order`='%v' where group_id='%v'",
		g.Delay, g.GlobalDelay, g.Order, g.ID,
	)

	_, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// SortGroup to update the sort order of groups.
func (sl *Sqlite) SortGroup(sceneID int, groupID int, sortID int) error {
	rows, err := sl.db.Query(
		"SELECT group_id FROM scenes_group where scene_id=? and group_id!=? ORDER BY `order`",
		sceneID, groupID,
	)
	if err != nil {
		log.Error(err)

		return err
	}

	type groupNeedingUpdate struct {
		GroupID int
		Order   int
	}

	var groupsNeedingUpdate []groupNeedingUpdate

	defer func() {
		if err = rows.Err(); err != nil {
			log.Errorf("Sqlite SortGroup: %v.", err)
		}

		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	i := 1

	for rows.Next() {
		if i == sortID {
			i++
		}

		var thisGroupID int

		err = rows.Scan(&thisGroupID)
		if err != nil {
			log.Error(err)

			return err
		}

		g := groupNeedingUpdate{
			GroupID: thisGroupID,
			Order:   i,
		}

		groupsNeedingUpdate = append(groupsNeedingUpdate, g)
		i++
	}

	for _, gnu := range groupsNeedingUpdate {
		sqlStmt := fmt.Sprintf("UPDATE scenes_group set `order`='%v' where group_id='%v'", gnu.Order, gnu.GroupID)

		_, err := sl.db.Exec(sqlStmt)
		if err != nil {
			log.Errorf("%q: %s", err, sqlStmt)

			return err
		}
	}

	sqlStmt := fmt.Sprintf("UPDATE scenes_group set `order`='%v' where group_id='%v'", sortID, groupID)

	_, err = sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// DeleteGroup function.
func (sl *Sqlite) DeleteGroup(groupID int) error {
	sqlStmt := fmt.Sprintf("DELETE from scenes_group where group_id=%v", groupID)

	_, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// GetActions to return a slice of Action structs.
func (sl *Sqlite) GetActions(groupID int) ([]models.Action, error) {
	s := []models.Action{}

	rows, err := sl.db.Query(
		"SELECT action_id, devices, command, parameter, global_parameter, `order` "+
			"FROM scenes_action where group_id=? ORDER BY `order`",
		groupID,
	)
	if err != nil {
		log.Error(err)

		return s, err
	}

	defer func() {
		if err = rows.Err(); err != nil {
			log.Errorf("Sqlite GetActions: %v.", err)
		}

		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	for rows.Next() {
		var actionID, order int

		var devicesString, command, parameter, globalParameter string

		err = rows.Scan(&actionID, &devicesString, &command, &parameter, &globalParameter, &order)
		if err != nil {
			log.Error(err)

			return s, err
		}

		devices := []models.Device{}

		ds := strings.Split(devicesString, " ")
		for _, v := range ds {
			tID, err := strconv.Atoi(v)
			if err != nil {
				log.Error(err)

				continue
			}

			devices = append(devices, sl.GetDevice(tID))
		}

		sa := models.Action{
			ID:              actionID,
			GroupID:         groupID,
			Devices:         devices,
			Command:         command,
			Parameter:       parameter,
			GlobalParameter: globalParameter,
			Order:           order,
		}
		s = append(s, sa)
	}

	return s, err
}

// GetAction to return a single Action struct.
func (sl *Sqlite) GetAction(actionID int) (models.Action, error) {
	sqlStmt := fmt.Sprintf(
		"SELECT group_id, devices, command, parameter, global_parameter, `order` "+
			"FROM scenes_action where action_id = '%v'",
		actionID,
	)

	var groupID, order int

	var devicesString, command, parameter, globalParameter string

	err := sl.db.QueryRow(sqlStmt).Scan(&groupID, &devicesString, &command, &parameter, &globalParameter, &order)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return models.Action{}, err
	}

	devices := []models.Device{}

	ds := strings.Split(devicesString, " ")
	for _, v := range ds {
		tID, err := strconv.Atoi(v)
		if err != nil {
			log.Errorf("%q: %s", err, sqlStmt)

			return models.Action{}, err
		}

		devices = append(devices, sl.GetDevice(tID))
	}

	action := models.Action{
		ID:              actionID,
		GroupID:         groupID,
		Devices:         devices,
		Command:         command,
		Parameter:       parameter,
		GlobalParameter: globalParameter,
		Order:           order,
	}

	return action, err
}

// GetActionOrderNext to get the next order value.
func (sl *Sqlite) GetActionOrderNext(groupID int) (int, error) {
	var count int

	row := sl.db.QueryRow("SELECT count(*) FROM scenes_action where group_id = ?", groupID)

	err := row.Scan(&count)
	if err != nil {
		log.Error(err)
	} else {
		count++
	}

	return count, err
}

// AddAction to db.
func (sl *Sqlite) AddAction(a models.Action) (int, error) {
	var devices string

	for _, d := range a.Devices {
		sid := strconv.Itoa(d.ID)
		if devices == "" {
			devices = sid
		} else {
			devices = devices + " " + sid
		}
	}

	sqlStmt := fmt.Sprintf(
		"INSERT INTO scenes_action(group_id, devices, command, parameter, global_parameter, 'order') "+
			"values('%v', '%s', '%s', '%s', '%s', '%v')",
		a.GroupID, devices, a.Command, a.Parameter, a.GlobalParameter, a.Order,
	)

	res, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Error(err)

		return 0, err
	}

	return int(id), err
}

// SetAction to update a Action.
func (sl *Sqlite) SetAction(a models.Action) error {
	var devices string

	for _, d := range a.Devices {
		sid := strconv.Itoa(d.ID)
		if devices == "" {
			devices = sid
		} else {
			devices = devices + " " + sid
		}
	}

	sqlStmt := fmt.Sprintf(
		"UPDATE scenes_action set devices='%s', command='%s', parameter='%s', "+
			"global_parameter='%s', `order`='%v' where action_id='%v'",
		devices, a.Command, a.Parameter, a.GlobalParameter, a.Order, a.ID,
	)

	_, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// SortAction function.
func (sl *Sqlite) SortAction(groupID int, actionID int, sortID int) error {
	rows, err := sl.db.Query(
		"SELECT action_id FROM scenes_action where group_id=? and action_id!=? ORDER BY `order`",
		groupID, actionID,
	)
	if err != nil {
		log.Error(err)

		return err
	}

	type actionNeedingUpdate struct {
		ActionID int
		Order    int
	}

	var actionsNeedingUpdate []actionNeedingUpdate

	defer func() {
		if err = rows.Err(); err != nil {
			log.Errorf("Sqlite SortAction: %v.", err)
		}

		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	i := 1
	for rows.Next() {
		if i == sortID {
			i++
		}

		var thisActionID int

		err = rows.Scan(&thisActionID)
		if err != nil {
			log.Error(err)

			return err
		}

		a := actionNeedingUpdate{
			ActionID: thisActionID,
			Order:    i,
		}

		actionsNeedingUpdate = append(actionsNeedingUpdate, a)
		i++
	}

	for _, anu := range actionsNeedingUpdate {
		sqlStmt := fmt.Sprintf("UPDATE scenes_action set `order`='%v' where action_id='%v'", anu.Order, anu.ActionID)

		_, err := sl.db.Exec(sqlStmt)
		if err != nil {
			log.Errorf("%q: %s", err, sqlStmt)

			return err
		}
	}

	sqlStmt := fmt.Sprintf("UPDATE scenes_action set `order`='%v' where action_id='%v'", sortID, actionID)

	_, err = sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// DeleteAction function.
func (sl *Sqlite) DeleteAction(actionID int) error {
	sqlStmt := fmt.Sprintf("DELETE from scenes_action where action_id=%v", actionID)

	_, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	return err
}

// AddDevice to add a new device.
func (sl *Sqlite) AddDevice(d models.Device) (insertID int) {
	sqlStmt := fmt.Sprintf("INSERT INTO devices(name, topic, type) values('%s', '%s', %v)", d.Name, d.Topic, d.Type.ID)

	res, err := sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}

	id, _ := res.LastInsertId()

	return int(id)
}

// SetDevice to update a device.
func (sl *Sqlite) SetDevice(d models.Device) {
	sqlStmt := fmt.Sprintf(
		"UPDATE devices set name='%s', topic='%s', type='%v' where device_id='%v'",
		d.Name, d.Topic, d.Type.ID, d.ID,
	)

	if _, err := sl.db.Exec(sqlStmt); err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}
}

// DeleteDevice function.
func (sl *Sqlite) DeleteDevice(deviceID int) {
	sqlStmt := fmt.Sprintf("DELETE from devices where device_id=%v", deviceID)

	if _, err := sl.db.Exec(sqlStmt); err != nil {
		log.Errorf("%q: %s", err, sqlStmt)
	}
}

// GetDevice to return a single Device struct.
func (sl *Sqlite) GetDevice(deviceID int) models.Device {
	d := models.Device{}
	sqlStmt := fmt.Sprintf("SELECT name, topic, type FROM devices where device_id = '%v'", deviceID)

	var name, topic string

	var typeID int

	err := sl.db.QueryRow(sqlStmt).Scan(&name, &topic, &typeID)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return d
	}

	d.ID = deviceID
	d.Name = name
	d.Topic = topic
	d.Type = sl.GetDeviceType(typeID)

	return d
}

const deviceSelect = "SELECT device_id, name, topic, type FROM devices"

// GetDevices to return a slice of Device structs.
func (sl *Sqlite) GetDevices() []models.Device {
	d := []models.Device{}

	sqlStmt := deviceSelect

	rows, err := sl.db.Query(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return d
	}

	defer func() {
		if err = rows.Err(); err != nil {
			log.Errorf("Sqlite GetDevices: %v.", err)
		}

		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	for rows.Next() {
		var deviceID, typeID int

		var name, topic string

		err = rows.Scan(&deviceID, &name, &topic, &typeID)
		if err != nil {
			log.Errorf("%q: %s", err, sqlStmt)

			return d
		}

		d = append(d, models.Device{ID: deviceID, Name: name, Topic: topic, Type: sl.GetDeviceType(typeID)})
	}

	return d
}

// GetDevicesWithActionSelected to return a slice of Device structs with select flags set.
func (sl *Sqlite) GetDevicesWithActionSelected(actionID int) []models.Device {
	d := []models.Device{}

	a, err := sl.GetAction(actionID)
	if err != nil {
		log.Error(err)

		return d
	}

	sqlStmt := deviceSelect

	rows, err := sl.db.Query(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return d
	}

	defer func() {
		if err = rows.Err(); err != nil {
			log.Errorf("Sqlite GetDevicesWithActionSelected: %v.", err)
		}

		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	for rows.Next() {
		var deviceID, typeID int

		var name, topic string

		err = rows.Scan(&deviceID, &name, &topic, &typeID)
		if err != nil {
			log.Errorf("%q: %s", err, sqlStmt)

			return d
		}

		var selected bool
		selected = false

		for _, dev := range a.Devices {
			if dev.ID == deviceID {
				selected = true
			}
		}

		d = append(
			d,
			models.Device{ID: deviceID, Name: name, Topic: topic, Type: sl.GetDeviceType(typeID), Selected: selected},
		)
	}

	return d
}

// ReturnDevicesWithActionSelected to return a slice of Device structs with select flags set.
func (sl *Sqlite) ReturnDevicesWithActionSelected(action models.Action, devices []models.Device) []models.Device {
	dr := []models.Device{}

	var selected bool
	for _, device := range devices {
		selected = false
		td := device

		for _, adevice := range action.Devices {
			if adevice.ID == device.ID {
				selected = true

				break
			}
		}

		td.Selected = selected
		dr = append(dr, td)
	}

	return dr
}

// GetDevicesWithSceneSelected to return a slice of Device structs with select flags set.
func (sl *Sqlite) GetDevicesWithSceneSelected(sceneID int) ([]models.Device, error) {
	d := []models.Device{}

	s, err := sl.GetScene(sceneID)
	if err != nil {
		return d, err
	}

	sqlStmt := deviceSelect

	rows, err := sl.db.Query(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return d, err
	}

	defer func() {
		if err = rows.Err(); err != nil {
			log.Errorf("Sqlite GetDevicesWithSceneSelected: %v.", err)
		}

		err := rows.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	for rows.Next() {
		var deviceID, typeID int

		var name, topic string

		err = rows.Scan(&deviceID, &name, &topic, &typeID)
		if err != nil {
			log.Errorf("%q: %s", err, sqlStmt)

			return d, err
		}

		var selected bool
		selected = false

		for _, dev := range s.AllowedDevices {
			if dev.ID == deviceID {
				selected = true
			}
		}

		d = append(
			d,
			models.Device{ID: deviceID, Name: name, Topic: topic, Type: sl.GetDeviceType(typeID), Selected: selected},
		)
	}

	return d, err
}

// GetDeviceTypes to return a slice of DeviceType structs.
func (sl *Sqlite) GetDeviceTypes() []models.DeviceType {
	return sl.deviceTypes
}

// GetDeviceType to return a single DeviceType struct.
func (sl *Sqlite) GetDeviceType(deviceID int) models.DeviceType {
	for _, d := range sl.deviceTypes {
		if d.ID == deviceID {
			return models.DeviceType{
				ID:       d.ID,
				Name:     d.Name,
				Commands: d.Commands,
			}
		}
	}

	return models.DeviceType{}
}

// Connect creates database connection.
func (sl *Sqlite) connect() {
	dbf := "data/sqlite.db"

	var err error

	sl.db, err = sql.Open("sqlite3", dbf)
	if err != nil {
		log.Error(err)
	}

	log.Infof("sqlite connected to %s", dbf)
}

// Disconnect database connection.
func (sl *Sqlite) Disconnect() {
	log.Info("disconnecting")

	if err := sl.db.Close(); err != nil {
		log.Error(err)
	}
}

func (sl *Sqlite) bootstrap() {
	// initialize database if needed.
	sqlStmt := "SELECT name FROM sqlite_master WHERE type='table' AND name='devices';"
	row := sl.db.QueryRow(sqlStmt)

	var name string

	err := row.Scan(&name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Infof("%v table does not exist, initializing database", name)
		}
	} else {
		log.Infof("%v table already exists, skipping database init", name)

		return
	}

	// create required tables.
	sqlStmt = "CREATE TABLE devices (device_id INTEGER PRIMARY KEY, name TEXT, topic TEXT, type INTEGER);"

	_, err = sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return
	}

	sqlStmt = "CREATE TABLE shows (show_id INTEGER PRIMARY KEY, name TEXT, topic TEXT, repeat TEXT, " +
		"global_delay INTEGER, global_speed INTEGER, global_parameter1 TEXT, global_parameter2 TEXT);"

	_, err = sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return
	}

	sqlStmt = "CREATE TABLE shows_cycles (cycle_id INTEGER PRIMARY KEY, show_id INTEGER, scene_id INTEGER, " +
		"cycles INTEGER, end_delay INTEGER, loop_include TEXT, global_delay INTEGER, global_speed INTEGER, " +
		"global_parameter1 TEXT, global_parameter2 TEXT, 'order' INTEGER);"

	_, err = sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return
	}

	sqlStmt = "CREATE TABLE scenes (scene_id INTEGER PRIMARY KEY, name TEXT, allowed_devices TEXT);"

	_, err = sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return
	}

	sqlStmt = "CREATE TABLE scenes_group (group_id INTEGER PRIMARY KEY, scene_id INTEGER, " +
		"delay INTEGER, global_delay INTEGER, 'order' INTEGER);"

	_, err = sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return
	}

	sqlStmt = "CREATE TABLE scenes_action (action_id INTEGER PRIMARY KEY, group_id INTEGER, devices TEXT, " +
		"command TEXT, parameter TEXT, global_parameter TEXT, 'order' INTEGER);"

	_, err = sl.db.Exec(sqlStmt)
	if err != nil {
		log.Errorf("%q: %s", err, sqlStmt)

		return
	}
}

// InitializeClient is a replacement for init() because we need to set log variable first.
func (sl *Sqlite) InitializeClient() {
	sl.connect()
	sl.bootstrap()
}
