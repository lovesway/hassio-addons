package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lovesway/hassio-addons/mq-lightshow/database"
	"github.com/lovesway/hassio-addons/mq-lightshow/models"
)

var (
	errNoShowID   = errors.New("no showid given")
	errNoCycleID  = errors.New("no cycleid given")
	errNoSceneID  = errors.New("no sceneid given")
	errNoGroupID  = errors.New("no groupid given")
	errNoActionID = errors.New("no actionid given")
	errNoSortID   = errors.New("no sortid given")
)

// APIController represents the controller for the API.
type APIController struct {
	md Modeler
	ss StringsToStruct
	db *database.Sqlite
}

// NewAPIController provides an instance of APIController.
func NewAPIController(md Modeler, ss StringsToStruct, db *database.Sqlite) APIController {
	return APIController{
		md: md,
		ss: ss,
		db: db,
	}
}

// Response object.
type Response struct {
	Status  int16
	Error   bool
	Message string
}

// ResponseData object.
type ResponseData struct {
	Status  int16
	Error   bool
	Message string
	Data    interface{}
}

func getResponseError(message string) Response {
	if message == "" {
		message = "error"
	}

	return Response{
		Status:  http.StatusInternalServerError,
		Error:   true,
		Message: message,
	}
}

func getResponse(message string) Response {
	if message == "" {
		message = "success"
	}

	return Response{
		Status:  http.StatusOK,
		Error:   false,
		Message: message,
	}
}

func getResponseData() ResponseData {
	return ResponseData{
		Status:  http.StatusOK,
		Error:   false,
		Message: "success",
		Data:    nil,
	}
}

// Shows will return a list of Show objects.
func (ac APIController) Shows(w http.ResponseWriter, r *http.Request) {
	shows, err := ac.md.GetShows()
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponseData()
	re.Data = shows

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// Show will return a Show object.
func (ac APIController) Show(w http.ResponseWriter, r *http.Request) {
	showID, err := getShowIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	show, err := ac.md.GetShow(showID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponseData()
	re.Data = show

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

type showStrings struct {
	Name             string
	Topic            string
	Repeat           string
	GlobalDelay      string
	GlobalSpeed      string
	GlobalParameter1 string
	GlobalParameter2 string
}

func getShowIDFromRequest(r *http.Request) (int, error) {
	var err error

	var showID int

	v := mux.Vars(r)
	showIDString := v["showID"]

	if showIDString == "" {
		return showID, errNoShowID
	}

	showIDInt, err := strconv.Atoi(showIDString)
	if err != nil {
		return showID, err
	}

	return showIDInt, err
}

// ShowCreate will create a show.
func (ac APIController) ShowCreate(w http.ResponseWriter, r *http.Request) {
	var dd showStrings

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dd)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	show := models.Show{}

	show, err = ac.ss.Show(dd, show)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	_, err = ac.md.AddShow(show)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Show created successfully")
	re.Status = http.StatusCreated

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// ShowConfigure will update a Show.
func (ac APIController) ShowConfigure(w http.ResponseWriter, r *http.Request) {
	showID, err := getShowIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	var dd showStrings

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&dd)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	show := models.Show{}

	show, err = ac.ss.Show(dd, show)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	show.ID = showID

	err = ac.md.SetShow(show)
	if err != nil {
		log.Error(err)
	}

	re := getResponse("Show updated successfully")
	re.Status = http.StatusCreated

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// ShowDelete will delete a Show.
func (ac APIController) ShowDelete(w http.ResponseWriter, r *http.Request) {
	showID, err := getShowIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	err = ac.md.DeleteShow(showID)
	if err != nil {
		log.Error(err)
	}

	re := getResponse("Show deleted successfully")
	re.Status = http.StatusNoContent

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// ShowStart will execute a show.
func (ac APIController) ShowStart(w http.ResponseWriter, r *http.Request) {
	showID, err := getShowIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	if showID == 0 {
		re := getResponseError("No showID given")

		jsonErr := json.NewEncoder(w).Encode(re)
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	err = ex.StartShow(showID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	jsonErr := json.NewEncoder(w).Encode(getResponse("Show started"))
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// ShowStop will stop a show execution.
func (ac APIController) ShowStop(w http.ResponseWriter, r *http.Request) {
	showID, err := getShowIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	if showID == 0 {
		re := getResponseError("No showID given")

		jsonErr := json.NewEncoder(w).Encode(re)
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	err = ex.StopShow(showID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	jsonErr := json.NewEncoder(w).Encode(getResponse("Show stopped"))
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// ShowCycles will return a list of Cycle objects for given showID.
func (ac APIController) ShowCycles(w http.ResponseWriter, r *http.Request) {
	showID, err := getShowIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	data, err := ac.md.GetShowCyclesRecursive(showID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponseData()
	re.Data = data

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// ShowCycle will return a Cycle object for given cycleID.
func (ac APIController) ShowCycle(w http.ResponseWriter, r *http.Request) {
	cycleID, err := getCycleIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	data, err := ac.md.GetShowCycle(cycleID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponseData()
	re.Data = data

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

type cycleStrings struct {
	ShowID           string
	SceneID          string
	SceneCycles      string
	EndDelay         string
	LoopInclude      string
	GlobalDelay      string
	GlobalSpeed      string
	GlobalParameter1 string
	GlobalParameter2 string
}

func getCycleIDFromRequest(r *http.Request) (int, error) {
	var err error

	var cycleID int

	v := mux.Vars(r)
	cycleIDString := v["cycleID"]

	if cycleIDString == "" {
		return cycleID, errNoCycleID
	}

	cycleIDInt, err := strconv.Atoi(cycleIDString)
	if err != nil {
		return cycleID, err
	}

	return cycleIDInt, err
}

// ShowCycleCreate will create a show.
func (ac APIController) ShowCycleCreate(w http.ResponseWriter, r *http.Request) {
	showID, err := getShowIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	var dd cycleStrings

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&dd)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	cycle := models.Cycle{}

	cycle, err = ac.ss.Cycle(dd, cycle)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	cycle.ShowID = showID

	_, err = ac.md.AddShowCycle(cycle)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Cycle created successfully")
	re.Status = http.StatusCreated

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// ShowCycleEdit will update a Cycle.
func (ac APIController) ShowCycleEdit(w http.ResponseWriter, r *http.Request) {
	showID, err := getShowIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	cycleID, err := getCycleIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	var dd cycleStrings

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&dd)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	cycle := models.Cycle{}

	cycle, err = ac.ss.Cycle(dd, cycle)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	cycle.ID = cycleID
	cycle.ShowID = showID

	err = ac.md.SetShowCycle(cycle)
	if err != nil {
		log.Error(err)
	}

	re := getResponse("Cycle updated successfully")
	re.Status = http.StatusCreated

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// ShowCycleDelete will delete a Cycle.
func (ac APIController) ShowCycleDelete(w http.ResponseWriter, r *http.Request) {
	cycleID, err := getCycleIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	err = ac.md.DeleteShowCycle(cycleID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Cycle deleted successfully")
	re.Status = 204

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

type sceneStrings struct {
	Name             string
	AllowedDeviceIDs []string
}

func getSceneIDFromRequest(r *http.Request) (int, error) {
	var err error

	var sceneID int

	v := mux.Vars(r)
	sceneIDString := v["sceneID"]

	if sceneIDString == "" {
		return sceneID, errNoSceneID
	}

	sceneIDInt, err := strconv.Atoi(sceneIDString)
	if err != nil {
		return sceneID, err
	}

	return sceneIDInt, err
}

// Scenes will return a list of Scene objects.
func (ac APIController) Scenes(w http.ResponseWriter, r *http.Request) {
	scenes, err := ac.md.GetScenes()
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponseData()
	re.Data = scenes

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneCreate will create a new Scene.
func (ac APIController) SceneCreate(w http.ResponseWriter, r *http.Request) {
	var dd sceneStrings

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err := dec.Decode(&dd)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	scene := models.Scene{}

	scene, err = ac.ss.Scene(dd, scene)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	for _, deviceIDString := range dd.AllowedDeviceIDs {
		deviceID, err := strconv.Atoi(deviceIDString)
		if err != nil {
			jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
			if jsonErr != nil {
				log.Error(jsonErr)
			}

			return
		}

		scene.AllowedDevices = append(scene.AllowedDevices, ac.db.GetDevice(deviceID))
	}

	_, err = ac.md.AddScene(scene)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Scene created successfully")
	re.Status = http.StatusCreated

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// Scene will return a Scene.
func (ac APIController) Scene(w http.ResponseWriter, r *http.Request) {
	sceneID, err := getSceneIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	scene, err := ac.md.GetScene(sceneID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponseData()
	re.Data = scene

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneConfigure will update a Scene.
func (ac APIController) SceneConfigure(w http.ResponseWriter, r *http.Request) {
	sceneID, err := getSceneIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	var dd sceneStrings

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err = dec.Decode(&dd)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	scene := models.Scene{}

	scene, err = ac.ss.Scene(dd, scene)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	for _, deviceIDString := range dd.AllowedDeviceIDs {
		deviceID, err := strconv.Atoi(deviceIDString)
		if err != nil {
			jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
			if jsonErr != nil {
				log.Error(jsonErr)
			}

			return
		}

		scene.AllowedDevices = append(scene.AllowedDevices, ac.db.GetDevice(deviceID))
	}

	scene.ID = sceneID

	err = ac.md.SetScene(scene)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Scene configured successfully")
	re.Status = http.StatusCreated

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneRun will execute a Scene.
func (ac APIController) SceneRun(w http.ResponseWriter, r *http.Request) {
	sceneID, err := getSceneIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	go ex.ExecuteSceneByID(sceneID)

	jsonErr := json.NewEncoder(w).Encode(getResponse("Scene run successful"))
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneDuplicate will duplicate a Scene.
func (ac APIController) SceneDuplicate(w http.ResponseWriter, r *http.Request) {
	sceneID, err := getSceneIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	scene, err := ac.md.GetScene(sceneID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	scene.Name += " (dupe)"

	dupeSceneID, err := ac.md.AddScene(scene)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	groups, err := ac.md.GetGroups(sceneID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))

		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	for _, group := range groups {
		group.SceneID = dupeSceneID

		dupeGroupID, err := ac.md.AddGroup(group)
		if err != nil {
			jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
			if jsonErr != nil {
				log.Error(jsonErr)
			}

			return
		}

		actions, err := ac.md.GetActions(group.ID)
		if err != nil {
			jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
			if jsonErr != nil {
				log.Error(jsonErr)
			}

			return
		}

		// duplicate actions.
		for _, action := range actions {
			action.GroupID = dupeGroupID

			_, err := ac.md.AddAction(action)
			if err != nil {
				jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
				if jsonErr != nil {
					log.Error(jsonErr)
				}

				return
			}
		}
	}

	re := getResponse("Scene duplicated successfully")
	re.Status = 204

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneDelete will delete a Scene.
func (ac APIController) SceneDelete(w http.ResponseWriter, r *http.Request) {
	sceneID, err := getSceneIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	err = ac.md.DeleteScene(sceneID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Scene deleted successfully")
	re.Status = 204

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

type groupStrings struct {
	SceneID     string
	Delay       string
	GlobalDelay string
	Order       string
}

func getGroupIDFromRequest(r *http.Request) (int, error) {
	var err error

	var groupID int

	v := mux.Vars(r)
	groupIDString := v["groupID"]

	if groupIDString == "" {
		return groupID, errNoGroupID
	}

	groupIDInt, err := strconv.Atoi(groupIDString)
	if err != nil {
		return groupID, err
	}

	return groupIDInt, err
}

// SceneGroups will return a list of Group objects.
func (ac APIController) SceneGroups(w http.ResponseWriter, r *http.Request) {
	sceneID, err := getSceneIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	groups, err := ac.md.GetGroups(sceneID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponseData()
	re.Data = groups

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroup will return a Group object.
func (ac APIController) SceneGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := getGroupIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	group, err := ac.md.GetGroup(groupID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponseData()
	re.Data = group

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupCreate will create a new Group.
func (ac APIController) SceneGroupCreate(w http.ResponseWriter, r *http.Request) {
	sceneID, err := getSceneIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	var dd groupStrings

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err = dec.Decode(&dd)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	group := models.Group{}

	group, err = ac.ss.Group(dd, group)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	orderNext, err := ac.md.GetGroupOrderNext(sceneID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	group.SceneID = sceneID
	group.Order = orderNext

	_, err = ac.md.AddGroup(group)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Group created successfully")
	re.Status = http.StatusCreated

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupConfigure will update a Group.
func (ac APIController) SceneGroupConfigure(w http.ResponseWriter, r *http.Request) {
	sceneID, err := getSceneIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	groupID, err := getGroupIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	var dd groupStrings

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	err = dec.Decode(&dd)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	group := models.Group{}

	group, err = ac.ss.Group(dd, group)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	group.ID = groupID
	group.SceneID = sceneID

	err = ac.md.SetGroup(group)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Group configured successfully")
	re.Status = http.StatusCreated

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupRun will execute a Group.
func (ac APIController) SceneGroupRun(w http.ResponseWriter, r *http.Request) {
	groupID, err := getGroupIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	go ex.ExecuteActionGroupByID(groupID)

	jsonErr := json.NewEncoder(w).Encode(getResponse("Group run successful"))
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupDuplicate will duplicate a Group.
func (ac APIController) SceneGroupDuplicate(w http.ResponseWriter, r *http.Request) {
	sceneID, err := getSceneIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	groupID, err := getGroupIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	group, err := ac.md.GetGroup(groupID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	actions, err := ac.md.GetActions(groupID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	orderNext, err := ac.md.GetGroupOrderNext(sceneID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	group.ID = 0
	group.Order = orderNext

	dupeGroupID, err := ac.md.AddGroup(group)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	// duplicate actions.
	for _, action := range actions {
		action.GroupID = dupeGroupID

		_, err := ac.md.AddAction(action)
		if err != nil {
			jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
			if jsonErr != nil {
				log.Error(jsonErr)
			}

			return
		}
	}

	re := getResponse("Group duplicated successfully")
	re.Status = 204

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupSort will change the sort order of a Group.
func (ac APIController) SceneGroupSort(w http.ResponseWriter, r *http.Request) {
	sceneID, err := getSceneIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	groupID, err := getGroupIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	sortID, err := getSortIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	err = ac.md.SortGroup(sceneID, groupID, sortID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Group sort order changed successfully")
	re.Status = 204

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupDelete will delete a Group.
func (ac APIController) SceneGroupDelete(w http.ResponseWriter, r *http.Request) {
	groupID, err := getGroupIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	group, err := ac.md.GetGroup(groupID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	for _, a := range group.Actions {
		err = ac.md.DeleteAction(a.ID)
		if err != nil {
			jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
			if jsonErr != nil {
				log.Error(jsonErr)
			}

			return
		}
	}

	err = ac.md.DeleteGroup(groupID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Group deleted successfully")
	re.Status = http.StatusNoContent

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

type actionStrings struct {
	GroupID         string
	DeviceIDs       []string
	Command         string
	Parameter       string
	GlobalParameter string
}

func getActionIDFromRequest(r *http.Request) (int, error) {
	var err error

	var actionID int

	v := mux.Vars(r)
	actionIDString := v["actionID"]

	if actionIDString == "" {
		return actionID, errNoActionID
	}

	actionIDInt, err := strconv.Atoi(actionIDString)
	if err != nil {
		return actionID, err
	}

	return actionIDInt, err
}

// SceneGroupActions will return a list of Action objects.
func (ac APIController) SceneGroupActions(w http.ResponseWriter, r *http.Request) {
	groupID, err := getGroupIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	actions, err := ac.md.GetActions(groupID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponseData()
	re.Data = actions

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupAction will return an Action object.
func (ac APIController) SceneGroupAction(w http.ResponseWriter, r *http.Request) {
	actionID, err := getActionIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	action, err := ac.md.GetAction(actionID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponseData()
	re.Data = action

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupActionCreate will create a new Action.
func (ac APIController) SceneGroupActionCreate(w http.ResponseWriter, r *http.Request) {
	groupID, err := getGroupIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	var dd actionStrings

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&dd)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	action := models.Action{}

	action, err = ac.ss.Action(dd, action)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	orderNext, err := ac.md.GetActionOrderNext(groupID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	action.GroupID = groupID
	action.Order = orderNext

	for _, deviceIDString := range dd.DeviceIDs {
		deviceID, err := strconv.Atoi(deviceIDString)
		if err != nil {
			jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
			if jsonErr != nil {
				log.Error(jsonErr)
			}

			return
		}

		action.Devices = append(action.Devices, ac.db.GetDevice(deviceID))
	}

	_, err = ac.md.AddAction(action)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Action created successfully")
	re.Status = http.StatusCreated

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupActionEdit will update an Action.
func (ac APIController) SceneGroupActionEdit(w http.ResponseWriter, r *http.Request) {
	groupID, err := getGroupIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	actionID, err := getActionIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	var dd actionStrings

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&dd)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Error(err.Error())
		}
	}()

	action := models.Action{}

	action, err = ac.ss.Action(dd, action)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	action.ID = actionID
	action.GroupID = groupID

	existingAction, err := ac.md.GetAction(actionID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	action.Order = existingAction.Order

	for _, deviceIDString := range dd.DeviceIDs {
		deviceID, err := strconv.Atoi(deviceIDString)
		if err != nil {
			jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
			if jsonErr != nil {
				log.Error(jsonErr)
			}

			return
		}

		action.Devices = append(action.Devices, ac.db.GetDevice(deviceID))
	}

	err = ac.md.SetAction(action)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Action updated successfully")
	re.Status = http.StatusCreated

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupActionRun will execute an Action.
func (ac APIController) SceneGroupActionRun(w http.ResponseWriter, r *http.Request) {
	actionID, err := getActionIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	go ex.ExecuteActionByID(actionID)

	jsonErr := json.NewEncoder(w).Encode(getResponse("Action run successful"))
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupActionDuplicate will duplicate an Action.
func (ac APIController) SceneGroupActionDuplicate(w http.ResponseWriter, r *http.Request) {
	groupID, err := getGroupIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	actionID, err := getActionIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	action, err := ac.md.GetAction(actionID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	orderNext, err := ac.md.GetActionOrderNext(groupID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	action.ID = 0
	action.Order = orderNext

	_, err = ac.md.AddAction(action)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Action duplicated successfully")
	re.Status = http.StatusNoContent

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

func getSortIDFromRequest(r *http.Request) (int, error) {
	var err error

	var sortID int

	v := mux.Vars(r)
	sortIDString := v["sortID"]

	if sortIDString == "" {
		return sortID, errNoSortID
	}

	sortIDInt, err := strconv.Atoi(sortIDString)
	if err != nil {
		return sortID, err
	}

	return sortIDInt, err
}

// SceneGroupActionSort will change the sort order of an Action.
func (ac APIController) SceneGroupActionSort(w http.ResponseWriter, r *http.Request) {
	groupID, err := getGroupIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	actionID, err := getActionIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	sortID, err := getSortIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	err = ac.md.SortAction(groupID, actionID, sortID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Action sort order changed successfully")
	re.Status = 204

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}

// SceneGroupActionDelete will delete an Action.
func (ac APIController) SceneGroupActionDelete(w http.ResponseWriter, r *http.Request) {
	actionID, err := getActionIDFromRequest(r)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	err = ac.md.DeleteAction(actionID)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(getResponseError(err.Error()))
		if jsonErr != nil {
			log.Error(jsonErr)
		}

		return
	}

	re := getResponse("Action deleted successfully")
	re.Status = http.StatusNoContent

	jsonErr := json.NewEncoder(w).Encode(re)
	if jsonErr != nil {
		log.Error(jsonErr)
	}
}
