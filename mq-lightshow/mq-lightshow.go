package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"gitlab.local/hassio-addons/mq-lightshow/database"
	"gitlab.local/hassio-addons/mq-lightshow/devicetypes"
	"gitlab.local/hassio-addons/mq-lightshow/models"
	"go.uber.org/zap"
)

var (
	log     *zap.SugaredLogger // global log adapter to support zap sugar.
	ex      *Executor          // globally accessible instance of Executor.
	version = "development"    // injected by the build process
)

const (
	globalDelay      = "GlobalDelay"
	globalSpeed      = "GlobalSpeed"
	globalParameter1 = "GlobalParameter1"
	globalParameter2 = "GlobalParameter2"
)

func setLogger(l *zap.SugaredLogger) {
	log = l
	database.SetLogger(log)
}

// Config provides basic configuration for http server.
type Config struct {
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// HTTPServer represents the web server.
type HTTPServer struct {
	server *http.Server
	wg     sync.WaitGroup
}

func getHaConfiguration() models.Configuration {
	jsonFile, err := os.Open("data/options.json")
	if err != nil {
		panic(fmt.Sprintf("Error reading Home Assistant config file: %s", err.Error()))
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = jsonFile.Close()
	if err != nil {
		panic(fmt.Sprintf("Error closing Home Assistant config file: %s", err.Error()))
	}

	var c models.Configuration

	err = json.Unmarshal(byteValue, &c)
	if err != nil {
		panic(fmt.Sprintf("Error reading Home Assistant config file: %s", err.Error()))
	}

	return c
}

func main() {
	conf := getHaConfiguration()

	setLogger(GetLogger(conf.LogLevel))

	defer loggerSync()

	log.Infof("Starting mq-lightshow version %s", version)

	dt := devicetypes.NewDeviceTypes()
	db := database.NewSqlite(dt)
	ss := NewStringsToStruct()
	md := NewModler(db)
	mq := NewMQController(md)
	ex = NewExecutor(md, mq)
	ac := NewAPIController(md, ss, db)
	c := NewController(md, db, mq, dt)

	db.InitializeClient()

	mq.MqttConnect(conf)

	const five = 5

	serverCfg := Config{
		Host:         ":8099",
		ReadTimeout:  five * time.Second,
		WriteTimeout: five * time.Second,
	}

	httpServer := Start(serverCfg, ac, c)

	defer func() {
		err := httpServer.Stop()
		if err != nil {
			log.Error(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	log.Info("Shutting down")
	mq.MqttDisconnect()
	db.Disconnect()
}

// Start launches the HTTP Server.
func Start(cfg Config, ac APIController, c Controller) *HTTPServer {
	// setup Context
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := getRouter(ac, c)

	const twenty = 20

	httpServer := HTTPServer{
		server: &http.Server{
			Addr:           cfg.Host,
			Handler:        router,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: 1 << twenty,
		},
	}

	httpServer.wg.Add(1)

	go func() {
		log.Infof("http server started for %v", cfg.Host)

		err := httpServer.server.ListenAndServe()
		if err != nil {
			log.Error(err)
		}

		httpServer.wg.Done()
	}()

	return &httpServer
}

// Stop turns off the HTTP Server.
func (httpServer *HTTPServer) Stop() error {
	// create a context to attempt a graceful 5 second shutdown.
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	log.Infof("http server stopping")

	// attempt the graceful shutdown by closing the listener and completing all inflight requests.
	if err := httpServer.server.Shutdown(ctx); err != nil {
		// looks like we timed out on the graceful shutdown. Force close
		if err := httpServer.server.Close(); err != nil {
			log.Errorf("error stopping http server: %v", err)

			return err
		}
	}

	// wait for the listener to report that it is closed.
	httpServer.wg.Wait()
	log.Info("http server stopped")

	return nil
}
