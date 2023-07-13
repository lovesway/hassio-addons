package main

import (
	"fmt"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/lovesway/hassio-addons/mq-lightshow/models"
)

// MQController struct to represent a class.
type MQController struct {
	mc                          MQTT.Client
	messages                    []models.Message
	subscribeInitIgnoreMessages bool
	f                           MQTT.MessageHandler
	md                          Modeler
}

// NewMQController method to instantiate class/struct.
func NewMQController(md Modeler) *MQController {
	mqc := &MQController{
		md:                          md,
		subscribeInitIgnoreMessages: true,
	}

	mqc.f = func(client MQTT.Client, msg MQTT.Message) {
		if mqc.subscribeInitIgnoreMessages {
			log.Debug("mqtt init: ignoring message while initializing")

			return
		}

		const maxMessagesTracked = 50 // only keep the last 50 messages
		if len(mqc.messages) >= maxMessagesTracked {
			mqc.messages = mqc.messages[1:]
		}

		mqc.messages = append(mqc.messages, models.Message{Topic: msg.Topic(), Message: string(msg.Payload())})

		cmd := string(msg.Payload())
		if cmd == "ON" || cmd == "OFF" {
			topicShowSplit := strings.Split(msg.Topic(), "/")

			const three = 3
			if len(topicShowSplit) < three {
				return
			}

			topicShow := topicShowSplit[2]

			show, err := md.GetShowByTopic(topicShow)
			if err != nil {
				log.Errorf("mqtt error: cannot get show by topic: %v", topicShow)

				return
			}

			if cmd == "ON" {
				err = ex.StartShow(show.ID)
				if err != nil {
					log.Error(err.Error())
				}
			} else if cmd == "OFF" {
				err = ex.StopShow(show.ID)
				if err != nil {
					log.Error(err.Error())
				}
			}

			return
		}

		log.Debug("TOPIC: %s", msg.Topic())
		log.Debug("MSG: %s", msg.Payload())
	}

	return mqc
}

// GetMessages returns last 100 messages.
func (mqc *MQController) GetMessages() []models.Message {
	return mqc.messages
}

// MqttConnect to the mqtt server.
func (mqc *MQController) MqttConnect(config models.Configuration) {
	if config.MQTTHost == "" || config.MQTTUser == "" || config.MQTTPass == "" {
		log.Error("Connect requested, but needed values are missing.")

		return
	}

	opts := MQTT.NewClientOptions().AddBroker(config.MQTTHost)
	opts.SetClientID("mq-lightshow")
	opts.SetUsername(config.MQTTUser)
	opts.SetPassword(config.MQTTPass)
	opts.SetDefaultPublishHandler(mqc.f)

	mqc.mc = MQTT.NewClient(opts)
	if token := mqc.mc.Connect(); token.Wait() && token.Error() != nil {
		log.Errorf("%v", token.Error())
	} else {
		log.Info("client connected")
	}

	mqc.subscribeInitIgnoreMessages = true
	mqc.SubscribeShows()

	go mqc.resetSubscribeInit()
}

func (mqc *MQController) resetSubscribeInit() {
	const five = 5

	time.Sleep(five * time.Second)

	mqc.subscribeInitIgnoreMessages = false
}

// MqttDisconnect from the mqtt server.
func (mqc *MQController) MqttDisconnect() {
	if !mqc.IsConnected() {
		log.Info("client was already disconnected")

		return
	}

	log.Info("client disconnecting")

	const eightHundred = 800

	mqc.mc.Disconnect(eightHundred)
}

// IsConnected determines the state of the mqtt server connection.
func (mqc *MQController) IsConnected() bool {
	if mqc.mc == nil {
		return false
	}

	return mqc.mc.IsConnected()
}

// SubscribeShows method.
func (mqc *MQController) SubscribeShows() {
	shows, err := mqc.md.GetShows()
	if err != nil {
		log.Errorf("error subscribing: ", err.Error())

		return
	}

	for _, show := range shows {
		if show.Topic != "" {
			mqc.SendShowState(show.Topic, "OFF")
			t := fmt.Sprintf("mqlightshow/show/%v/cmnd/#", show.Topic)
			mqc.Subscribe(t)
		}
	}
}

// Subscribe to a topic.
func (mqc *MQController) Subscribe(topic string) {
	if token := mqc.mc.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Error(token.Error().Error())
	}
}

// SendAction to send an action message to the mqtt server.
func (mqc *MQController) SendAction(topic string, command string, parameter string) {
	_topic := fmt.Sprintf("%s/cmnd/%s", topic, command)
	_token := mqc.mc.Publish(_topic, 0, false, parameter)
	_token.Wait()
}

// SendShowState to send an action message to the mqtt server.
func (mqc *MQController) SendShowState(topicShow string, state string) {
	_topic := fmt.Sprintf("mqlightshow/show/%s/stat", topicShow)
	_token := mqc.mc.Publish(_topic, 0, true, state)
	_token.Wait()
}
