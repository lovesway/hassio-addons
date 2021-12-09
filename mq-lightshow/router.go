package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getRouter(ac APIController, c Controller) *mux.Router {
	// setup route handlers
	router := mux.NewRouter()
	router = getRoutesAPI(router, ac)
	router = getRoutesUI(router, c)

	return router
}

func getRoutesAPI(router *mux.Router, ac APIController) *mux.Router {
	router.HandleFunc("/api/v1/shows", ac.Shows).Methods("GET")
	router.HandleFunc("/api/v1/show", ac.ShowCreate).Methods("POST")
	router.HandleFunc("/api/v1/show/{showID}", ac.Show).Methods("GET")
	router.HandleFunc("/api/v1/show/{showID}/start", ac.ShowStart).Methods("POST")
	router.HandleFunc("/api/v1/show/{showID}/stop", ac.ShowStop).Methods("POST")
	router.HandleFunc("/api/v1/show/{showID}/configure", ac.ShowConfigure).Methods("POST")
	router.HandleFunc("/api/v1/show/{showID}/delete", ac.ShowDelete).Methods("POST")
	router.HandleFunc("/api/v1/show/{showID}/cycles", ac.ShowCycles).Methods("GET")
	router.HandleFunc("/api/v1/show/{showID}/cycle", ac.ShowCycleCreate).Methods("POST")
	router.HandleFunc("/api/v1/show/{showID}/cycle/{cycleID}", ac.ShowCycle).Methods("GET")
	router.HandleFunc("/api/v1/show/{showID}/cycle/{cycleID}/edit", ac.ShowCycleEdit).Methods("POST")
	router.HandleFunc("/api/v1/show/{showID}/cycle/{cycleID}/delete", ac.ShowCycleDelete).Methods("POST")
	router.HandleFunc("/api/v1/scenes", ac.Scenes).Methods("GET")
	router.HandleFunc("/api/v1/scene", ac.SceneCreate).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}", ac.Scene).Methods("GET")
	router.HandleFunc("/api/v1/scene/{sceneID}/configure", ac.SceneConfigure).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}/run", ac.SceneRun).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}/duplicate", ac.SceneDuplicate).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}/delete", ac.SceneDelete).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}/groups", ac.SceneGroups).Methods("GET")
	router.HandleFunc("/api/v1/scene/{sceneID}/group", ac.SceneGroupCreate).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}/group/{groupID}", ac.SceneGroup).Methods("GET")
	router.HandleFunc("/api/v1/scene/{sceneID}/group/{groupID}/configure", ac.SceneGroupConfigure).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}/group/{groupID}/run", ac.SceneGroupRun).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}/group/{groupID}/duplicate", ac.SceneGroupDuplicate).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}/group/{groupID}/sort/{sortID}", ac.SceneGroupSort).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}/group/{groupID}/delete", ac.SceneGroupDelete).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}/group/{groupID}/actions", ac.SceneGroupActions).Methods("GET")
	router.HandleFunc("/api/v1/scene/{sceneID}/group/{groupID}/action", ac.SceneGroupActionCreate).Methods("POST")
	router.HandleFunc("/api/v1/scene/{sceneID}/group/{groupID}/action/{actionID}", ac.SceneGroupAction).Methods("GET")
	router.HandleFunc(
		"/api/v1/scene/{sceneID}/group/{groupID}/action/{actionID}/edit", ac.SceneGroupActionEdit,
	).Methods("POST")
	router.HandleFunc(
		"/api/v1/scene/{sceneID}/group/{groupID}/action/{actionID}/run", ac.SceneGroupActionRun,
	).Methods("POST")
	router.HandleFunc(
		"/api/v1/scene/{sceneID}/group/{groupID}/action/{actionID}/duplicate", ac.SceneGroupActionDuplicate,
	).Methods("POST")
	router.HandleFunc(
		"/api/v1/scene/{sceneID}/group/{groupID}/action/{actionID}/sort/{sortID}", ac.SceneGroupActionSort,
	).Methods("POST")
	router.HandleFunc(
		"/api/v1/scene/{sceneID}/group/{groupID}/action/{actionID}/delete", ac.SceneGroupActionDelete,
	).Methods("POST")

	return router
}

func getRoutesUI(router *mux.Router, c Controller) *mux.Router {
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("www/js/"))))
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("www/css/"))))
	router.PathPrefix("/icons/").Handler(http.StripPrefix("/icons/", http.FileServer(http.Dir("www/icons/"))))
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("www/images/"))))

	router.HandleFunc("/", c.ShowsHandler)
	router.HandleFunc("/shows", c.ShowsHandler)
	router.HandleFunc("/shows-add", c.ShowsAddHandler)
	router.HandleFunc("/shows-configure", c.ShowsConfigureHandler)
	router.HandleFunc("/shows-cycles", c.ShowsCyclesHandler)
	router.HandleFunc("/shows-cycles-add", c.ShowsCyclesAddHandler)
	router.HandleFunc("/shows-cycles-edit", c.ShowsCyclesEditHandler)
	router.HandleFunc("/scenes", c.ScenesHandler)
	router.HandleFunc("/scenes-add", c.ScenesAddHandler)
	router.HandleFunc("/scenes-configure", c.ScenesConfigureHandler)
	router.HandleFunc("/scenes-groups", c.ScenesGroupsHandler)
	router.HandleFunc("/scenes-groups-add", c.ScenesGroupsAddHandler)
	router.HandleFunc("/scenes-groups-configure", c.ScenesGroupsConfigureHandler)
	router.HandleFunc("/scenes-groups-actions", c.ScenesActionsHandler)
	router.HandleFunc("/scenes-groups-actions-add", c.ScenesActionsAddHandler)
	router.HandleFunc("/scenes-groups-actions-edit", c.ScenesActionsEditHandler)
	router.HandleFunc("/devices", c.DevicesHandler)
	router.HandleFunc("/devices-add", c.DevicesAddHandler)
	router.HandleFunc("/devices-edit", c.DevicesEditHandler)
	router.HandleFunc("/devices-delete", c.DevicesDeleteHandler)
	router.HandleFunc("/mqtt", c.MqttHandler)
	router.HandleFunc("/mqtt-log", c.MqttLogHandler)

	return router
}
