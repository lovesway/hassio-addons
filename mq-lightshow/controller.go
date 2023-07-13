package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/lovesway/hassio-addons/mq-lightshow/database"
	"github.com/lovesway/hassio-addons/mq-lightshow/devicetypes"
	"github.com/lovesway/hassio-addons/mq-lightshow/models"
)

// Controller represents the controller for the UI.
type Controller struct {
	md Modeler
	db *database.Sqlite
	mq *MQController
	dt []models.DeviceType
}

// NewController provides an instance of Controller.
func NewController(md Modeler, db *database.Sqlite, mq *MQController, dts devicetypes.DeviceTypes) Controller {
	return Controller{
		md: md,
		db: db,
		mq: mq,
		dt: dts.GetDeviceTypes(),
	}
}

// PageInfo represents common page items, mostly used in the header and footer.
type PageInfo struct {
	Title              string
	Posted             bool
	ScenesLinkEnabled  bool
	DevicesLinkEnabled bool
	MQTTLinkEnabled    bool
}

// httpRedirect is a handler for hassio ingress.
func httpRedirect(w http.ResponseWriter, r *http.Request, url string) {
	ingressPath := r.Header.Get("X-Ingress-Path")
	urlWithIngressPath := fmt.Sprintf("%v/%v", ingressPath, url)
	http.Redirect(w, r, urlWithIngressPath, http.StatusSeeOther)
}

func httpErrorHandler(w http.ResponseWriter, message string) {
	log.Errorf("httpError: %v", message)
	http.Error(w, message, http.StatusInternalServerError)
}

// ShowsHandler function.
func (c Controller) ShowsHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("www/shows.tpl", "www/base.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	type data struct{ PageInfo PageInfo }

	tplErr := tpl.ExecuteTemplate(w, "base", data{PageInfo: PageInfo{Title: "MQ Light Show"}})
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// ShowsAddHandler controller.
func (c Controller) ShowsAddHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("www/shows-add.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	type data struct{ PageInfo PageInfo }

	tplErr := tpl.ExecuteTemplate(w, "content", data{PageInfo: PageInfo{Title: "Adding Show"}})
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// ShowsConfigureHandler controller.
func (c Controller) ShowsConfigureHandler(w http.ResponseWriter, r *http.Request) {
	showID := getRequestShowID(r)
	if showID == 0 {
		httpErrorHandler(w, "Url Param 'showID' is missing")

		return
	}

	tpl, err := template.ParseFiles("www/shows-configure.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Configuring Show", ScenesLinkEnabled: true}

	type data struct {
		PageInfo    PageInfo
		Show        models.Show
		GlobalDelay string // string conversion to maintain empty value comparison ability.
		GlobalSpeed string // string conversion to maintain empty value comparison ability.
	}

	show, err := c.md.GetShow(showID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	dat := data{
		PageInfo: pi,
		Show:     show,
	}

	var sz models.Show // for null struct comparison.
	if dat.Show.GlobalDelay != sz.GlobalDelay {
		dat.GlobalDelay = fmt.Sprintf("%v", dat.Show.GlobalDelay)
	}

	if dat.Show.GlobalSpeed != sz.GlobalSpeed {
		dat.GlobalSpeed = fmt.Sprintf("%v", dat.Show.GlobalSpeed)
	}

	tplErr := tpl.ExecuteTemplate(w, "content", dat)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

func getRequestShowID(r *http.Request) int {
	keys, ok := r.URL.Query()["showID"]
	if !ok || len(keys[0]) < 1 {
		return 0
	}

	showIDString := keys[0]

	showID, err := strconv.Atoi(showIDString)
	if err != nil {
		log.Error(err)

		return 0
	}

	return showID
}

// ShowsCyclesHandler function.
func (c Controller) ShowsCyclesHandler(w http.ResponseWriter, r *http.Request) {
	showID := getRequestShowID(r)
	if showID == 0 {
		httpErrorHandler(w, "Url Param 'showID' is missing")

		return
	}

	tpl, err := template.ParseFiles("www/shows-cycles.tpl", "www/base.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Show Cycles"}

	type data struct {
		PageInfo PageInfo
		Show     models.Show
	}

	show, err := c.md.GetShow(showID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	dat := data{
		PageInfo: pi,
		Show:     show,
	}

	tplErr := tpl.ExecuteTemplate(w, "base", dat)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// ShowsCyclesAddHandler controller.
func (c Controller) ShowsCyclesAddHandler(w http.ResponseWriter, r *http.Request) {
	showID := getRequestShowID(r)
	if showID == 0 {
		httpErrorHandler(w, "Url Param 'showID' is missing")

		return
	}

	tpl, err := template.ParseFiles("www/shows-cycles-add.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Adding Cycle"}

	type data struct {
		PageInfo PageInfo
		ShowID   int
		Scenes   []models.Scene
	}

	scenes, err := c.md.GetScenes()
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	dat := data{
		PageInfo: pi,
		ShowID:   showID,
		Scenes:   scenes,
	}

	tplErr := tpl.ExecuteTemplate(w, "content", dat)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// ShowsCyclesEditHandler controller.
func (c Controller) ShowsCyclesEditHandler(w http.ResponseWriter, r *http.Request) {
	cycleID := getRequestCycleID(r)
	if cycleID == 0 {
		httpErrorHandler(w, "Url Param 'cycleID' is missing")

		return
	}

	tpl, err := template.ParseFiles("www/shows-cycles-edit.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Editing Show Scene Cycle"}

	type data struct {
		PageInfo    PageInfo
		Scenes      []models.Scene
		Cycle       models.Cycle
		GlobalDelay string // string conversion to maintain empty value comparison ability.
		GlobalSpeed string // string conversion to maintain empty value comparison ability.
	}

	cycle, err := c.md.GetShowCycle(cycleID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	scenes, err := c.md.GetScenes()
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	dat := data{
		PageInfo: pi,
		Scenes:   scenes,
		Cycle:    cycle,
	}

	var cz models.Cycle // for null struct comparison.
	if dat.Cycle.GlobalDelay != cz.GlobalDelay {
		dat.GlobalDelay = fmt.Sprintf("%v", dat.Cycle.GlobalDelay)
	}

	if dat.Cycle.GlobalSpeed != cz.GlobalSpeed {
		dat.GlobalSpeed = fmt.Sprintf("%v", dat.Cycle.GlobalSpeed)
	}

	tplErr := tpl.ExecuteTemplate(w, "content", dat)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

func getRequestCycleID(r *http.Request) int {
	keys, ok := r.URL.Query()["cycleID"]
	if !ok || len(keys[0]) < 1 {
		log.Error("Url Param 'cycleID' is missing")

		return 0
	}

	cycleIDString := keys[0]

	cycleID, err := strconv.Atoi(cycleIDString)
	if err != nil {
		log.Error(err)

		return 0
	}

	return cycleID
}

// ScenesHandler controller.
func (c Controller) ScenesHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("www/scenes.tpl", "www/base.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Scenes", ScenesLinkEnabled: true}

	type data struct {
		PageInfo PageInfo
		Scenes   []models.Scene
	}

	scenes, err := c.md.GetScenes()
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	dat := data{
		PageInfo: pi,
		Scenes:   scenes,
	}

	tplErr := tpl.ExecuteTemplate(w, "base", dat)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// ScenesAddHandler controller.
func (c Controller) ScenesAddHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("www/scenes-add.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Adding Scene", ScenesLinkEnabled: true}

	type data struct {
		PageInfo PageInfo
		Devices  []models.Device
	}

	tplErr := tpl.ExecuteTemplate(w, "content", data{PageInfo: pi, Devices: c.db.GetDevices()})
	if tplErr != nil {
		log.Error(tplErr)
	}
}

func getRequestSceneID(r *http.Request) int {
	keys, ok := r.URL.Query()["sceneID"]
	if !ok || len(keys[0]) < 1 {
		log.Error("Url Param 'sceneID' is missing")

		return 0
	}

	sceneIDString := keys[0]

	sceneID, err := strconv.Atoi(sceneIDString)
	if err != nil {
		log.Error(err)

		return 0
	}

	return sceneID
}

// ScenesConfigureHandler controller.
func (c Controller) ScenesConfigureHandler(w http.ResponseWriter, r *http.Request) {
	sceneID := getRequestSceneID(r)
	if sceneID == 0 {
		httpErrorHandler(w, "Url Param 'sceneID' is missing")

		return
	}

	tpl, err := template.ParseFiles("www/scenes-configure.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Configuring Scene", ScenesLinkEnabled: true}

	type data struct {
		PageInfo PageInfo
		Scene    models.Scene
		Devices  []models.Device
	}

	scene, err := c.md.GetScene(sceneID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	devices, err := c.db.GetDevicesWithSceneSelected(sceneID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	dat := data{
		PageInfo: pi,
		Scene:    scene,
		Devices:  devices,
	}

	tplErr := tpl.ExecuteTemplate(w, "content", dat)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// ScenesGroupsHandler controller.
func (c Controller) ScenesGroupsHandler(w http.ResponseWriter, r *http.Request) {
	sceneID := getRequestSceneID(r)
	if sceneID == 0 {
		httpErrorHandler(w, "Url Param 'sceneID' is missing")

		return
	}

	tpl, err := template.ParseFiles("www/scenes-groups.tpl", "www/base.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Actions", ScenesLinkEnabled: true}

	type data struct {
		PageInfo     PageInfo
		Scene        models.Scene
		TotalSeconds string
		TotalMinutes string
	}

	scene, err := c.md.GetScene(sceneID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	groups, err := c.md.GetGroups(sceneID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	var totals float32
	for _, sg := range groups {
		totals += sg.Delay
	}

	const secondsInMinute = 60

	d := data{
		PageInfo:     pi,
		Scene:        scene,
		TotalSeconds: fmt.Sprintf("%.1f", totals),
		TotalMinutes: fmt.Sprintf("%.2f", totals/secondsInMinute),
	}

	tplErr := tpl.ExecuteTemplate(w, "base", d)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// ScenesGroupsAddHandler controller.
func (c Controller) ScenesGroupsAddHandler(w http.ResponseWriter, r *http.Request) {
	sceneID := getRequestSceneID(r)
	if sceneID == 0 {
		httpErrorHandler(w, "Url Param 'sceneID' is missing")

		return
	}

	tpl, err := template.ParseFiles("www/scenes-groups-add.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Adding Scene Group", ScenesLinkEnabled: true}

	type data struct {
		PageInfo  PageInfo
		SceneID   int
		OrderNext int
	}

	orderNext, err := c.md.GetGroupOrderNext(sceneID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	dat := data{
		PageInfo:  pi,
		SceneID:   sceneID,
		OrderNext: orderNext,
	}

	tplErr := tpl.ExecuteTemplate(w, "content", dat)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

func getRequestGroupID(r *http.Request) int {
	keys, ok := r.URL.Query()["groupID"]
	if !ok || len(keys[0]) < 1 {
		log.Error("Url Param 'groupID' is missing")

		return 0
	}

	groupIDString := keys[0]

	groupID, err := strconv.Atoi(groupIDString)
	if err != nil {
		log.Error(err)

		return 0
	}

	return groupID
}

// ScenesGroupsConfigureHandler controller.
func (c Controller) ScenesGroupsConfigureHandler(w http.ResponseWriter, r *http.Request) {
	sceneID := getRequestSceneID(r)
	if sceneID == 0 {
		httpErrorHandler(w, "Url Param 'sceneID' is missing")

		return
	}

	groupID := getRequestGroupID(r)
	if groupID == 0 {
		httpErrorHandler(w, "Url Param 'groupID' is missing")

		return
	}

	tpl, err := template.ParseFiles("www/scenes-groups-configure.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Editing Scene Group", ScenesLinkEnabled: true}

	type data struct {
		PageInfo PageInfo
		SceneID  int
		Group    models.Group
	}

	group, err := c.md.GetGroup(groupID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	dat := data{
		PageInfo: pi,
		SceneID:  sceneID,
		Group:    group,
	}

	tplErr := tpl.ExecuteTemplate(w, "content", dat)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// ScenesActionsHandler controller.
func (c Controller) ScenesActionsHandler(w http.ResponseWriter, r *http.Request) {
	sceneID := getRequestSceneID(r)
	if sceneID == 0 {
		httpErrorHandler(w, "Url Param 'sceneID' is missing")

		return
	}

	groupID := getRequestGroupID(r)
	if groupID == 0 {
		httpErrorHandler(w, "Url Param 'groupID' is missing")

		return
	}

	tpl, err := template.ParseFiles("www/scenes-groups-actions.tpl", "www/base.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Actions", ScenesLinkEnabled: true}

	type data struct {
		PageInfo PageInfo
		Scene    models.Scene
		Group    models.Group
		Actions  []models.Action
	}

	scene, err := c.md.GetScene(sceneID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	group, err := c.md.GetGroup(groupID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	actions, err := c.md.GetActions(groupID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	dat := data{
		PageInfo: pi,
		Scene:    scene,
		Group:    group,
		Actions:  actions,
	}

	tplErr := tpl.ExecuteTemplate(w, "base", dat)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// ScenesActionsAddHandler controller.
func (c Controller) ScenesActionsAddHandler(w http.ResponseWriter, r *http.Request) {
	sceneID := getRequestSceneID(r)
	if sceneID == 0 {
		httpErrorHandler(w, "Url Param 'sceneID' is missing")

		return
	}

	groupID := getRequestGroupID(r)
	if groupID == 0 {
		httpErrorHandler(w, "Url Param 'groupID' is missing")

		return
	}

	tpl, err := template.ParseFiles("www/scenes-groups-actions-add.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Adding Action", ScenesLinkEnabled: true}

	type data struct {
		PageInfo    PageInfo
		SceneID     int
		GroupID     int
		Action      models.Action
		Commands    []models.Command
		Devices     []models.Device
		DeviceTypes []models.DeviceType
		OrderNext   int
	}

	scene, err := c.md.GetScene(sceneID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	orderNext, err := c.md.GetActionOrderNext(groupID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	dat := data{
		PageInfo:    pi,
		SceneID:     sceneID,
		GroupID:     groupID,
		DeviceTypes: c.dt,
		OrderNext:   orderNext,
	}

	if len(scene.AllowedDevices) != 0 {
		dat.Devices = scene.AllowedDevices
	} else {
		dat.Devices = c.db.GetDevices()
	}

	tplErr := tpl.ExecuteTemplate(w, "content", dat)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// ScenesActionsEditHandler controller.
func (c Controller) ScenesActionsEditHandler(w http.ResponseWriter, r *http.Request) {
	sceneID := getRequestSceneID(r)
	if sceneID == 0 {
		httpErrorHandler(w, "Url Param 'sceneID' is missing")

		return
	}

	groupID := getRequestGroupID(r)
	if groupID == 0 {
		httpErrorHandler(w, "Url Param 'groupID' is missing")

		return
	}

	actionID := getRequestActionID(r)
	if actionID == 0 {
		httpErrorHandler(w, "Url Param 'actionID' is missing")

		return
	}

	tpl, err := template.ParseFiles("www/scenes-groups-actions-edit.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Editing Action", ScenesLinkEnabled: true}

	type data struct {
		PageInfo    PageInfo
		SceneID     int
		GroupID     int
		Action      models.Action
		Commands    []models.Command
		Devices     []models.Device
		DeviceTypes []models.DeviceType
	}

	action, err := c.md.GetAction(actionID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	dat := data{
		PageInfo:    pi,
		SceneID:     sceneID,
		GroupID:     groupID,
		Action:      action,
		DeviceTypes: c.dt,
	}

	scene, err := c.md.GetScene(sceneID)
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	if len(scene.AllowedDevices) != 0 {
		dat.Devices = c.db.ReturnDevicesWithActionSelected(dat.Action, scene.AllowedDevices)
	} else {
		dat.Devices = c.db.GetDevicesWithActionSelected(actionID)
	}

	tplErr := tpl.ExecuteTemplate(w, "content", dat)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

func getRequestActionID(r *http.Request) int {
	keys, ok := r.URL.Query()["actionID"]
	if !ok || len(keys[0]) < 1 {
		log.Error("Url Param 'actionID' is missing")

		return 0
	}

	IDstring := keys[0]

	actionID, err := strconv.Atoi(IDstring)
	if err != nil {
		log.Error(err)

		return 0
	}

	return actionID
}

// DevicesHandler controller.
func (c Controller) DevicesHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("www/devices.tpl", "www/base.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Devices", DevicesLinkEnabled: true}

	type data struct {
		PageInfo PageInfo
		Devices  []models.Device
	}

	tplErr := tpl.ExecuteTemplate(w, "base", data{PageInfo: pi, Devices: c.db.GetDevices()})
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// DevicesAddHandler controller.
func (c Controller) DevicesAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost { // handle posted input
		tid, terr := strconv.Atoi(r.PostFormValue("type"))
		if terr != nil {
			httpErrorHandler(w, terr.Error())

			return
		}

		d := models.Device{
			Name:  r.PostFormValue("name"),
			Topic: r.PostFormValue("topic"),
			Type:  models.DeviceType{ID: tid},
		}

		c.db.AddDevice(d)
		httpRedirect(w, r, "devices")

		return
	}

	tpl, err := template.ParseFiles("www/devices-add.tpl", "www/base.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Adding Device", DevicesLinkEnabled: true}

	type data struct {
		PageInfo    PageInfo
		DeviceTypes []models.DeviceType
	}

	tplErr := tpl.ExecuteTemplate(w, "base", data{PageInfo: pi, DeviceTypes: c.db.GetDeviceTypes()})
	if tplErr != nil {
		log.Error(tplErr)
	}
}

func getRequestDeviceID(r *http.Request) int {
	keys, ok := r.URL.Query()["deviceID"]
	if !ok || len(keys[0]) < 1 {
		log.Error("Url Param 'deviceID' is missing")

		return 0
	}

	IDstring := keys[0]

	deviceID, err := strconv.Atoi(IDstring)
	if err != nil {
		log.Error(err)

		return 0
	}

	return deviceID
}

// DevicesDeleteHandler controller.
func (c Controller) DevicesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := getRequestDeviceID(r)
	if deviceID == 0 {
		return
	}

	c.db.DeleteDevice(deviceID)

	httpRedirect(w, r, "devices")
}

// DevicesEditHandler controller.
func (c Controller) DevicesEditHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := getRequestDeviceID(r)
	if deviceID == 0 {
		return
	}

	if r.Method == http.MethodPost {
		tid, terr := strconv.Atoi(r.PostFormValue("type"))
		if terr != nil {
			httpErrorHandler(w, terr.Error())

			return
		}

		d := models.Device{
			ID:    deviceID,
			Name:  r.PostFormValue("name"),
			Topic: r.PostFormValue("topic"),
			Type:  c.db.GetDeviceType(tid),
		}

		c.db.SetDevice(d)
		httpRedirect(w, r, "devices")

		return
	}

	tpl, err := template.ParseFiles("www/devices-edit.tpl", "www/base.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "Editing Device", DevicesLinkEnabled: true}

	type data struct {
		PageInfo    PageInfo
		Device      models.Device
		DeviceTypes []models.DeviceType
	}

	tplErr := tpl.ExecuteTemplate(
		w, "base", data{PageInfo: pi, Device: c.db.GetDevice(deviceID), DeviceTypes: c.db.GetDeviceTypes()},
	)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// MqttHandler function.
func (c Controller) MqttHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_, ok := r.URL.Query()["disconnect"]
		if ok {
			c.mq.MqttDisconnect()
			time.Sleep(1 * time.Second)
			httpRedirect(w, r, "mqtt")

			return
		}

		_, ok = r.URL.Query()["connect"]
		if ok {
			c.mq.MqttConnect(getHaConfiguration())
			time.Sleep(1 * time.Second)
			httpRedirect(w, r, "mqtt")

			return
		}
	}

	tpl, err := template.ParseFiles("www/mqtt.tpl", "www/base.tpl")
	if err != nil {
		httpErrorHandler(w, err.Error())

		return
	}

	pi := PageInfo{Title: "MQTT"}
	pi.MQTTLinkEnabled = true

	type data struct {
		PageInfo      PageInfo
		MQTTConnected bool
		MQTTHost      string
	}

	conf := getHaConfiguration()

	d := data{PageInfo: pi}
	d.MQTTConnected = c.mq.IsConnected()
	d.MQTTHost = conf.MQTTHost

	tplErr := tpl.ExecuteTemplate(w, "base", d)
	if tplErr != nil {
		log.Error(tplErr)
	}
}

// MqttLogHandler function.
func (c Controller) MqttLogHandler(w http.ResponseWriter, r *http.Request) {
	d := c.mq.GetMessages()

	j, err := json.Marshal(d)
	if err != nil {
		log.Errorf("Cannot encode to JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(j)
	if err != nil {
		log.Error(err)
	}
}
