package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/radhus/vsxknob/handler"
)

const (
	topicState = "vsx/state"
	topicSet   = "vsx/set"

	maxVolume = 185
)

type state struct {
	Power  bool    `json:"power"`
	Volume float64 `json:"volume"`
	Muted  bool    `json:"muted"`
	Source string  `json:"source"`

	volumeSet bool
	powerSet  bool
	mutedSet  bool
	sourceSet bool
}

type request struct {
	Power  *bool    `json:"power,omitempty"`
	Volume *float64 `json:"volume,omitempty"`
	Muted  *bool    `json:"muted,omitempty"`
	Source *string  `json:"source,omitempty"`
}

type connection struct {
	client mqtt.Client
	setter handler.Setter

	lastState state
}

func New(url string) (*connection, error) {
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	opts := mqtt.NewClientOptions().AddBroker(url).SetClientID("vsxknob")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("Failed to connect to MQTT: %w", token.Error())
	}

	connection := &connection{
		client: client,
	}

	if token := client.Subscribe(topicSet, 0, connection.subscribe); token.Wait() && token.Error() != nil {
		client.Disconnect(0)
		return nil, fmt.Errorf("Failed to subscribe to MQTT: %w", token.Error())
	}

	return connection, nil
}

func (c *connection) Setter(setter handler.Setter) {
	c.setter = setter
}

func (c *connection) subscribe(_ mqtt.Client, message mqtt.Message) {
	if c.setter == nil {
		log.Println("Cannot handle MQTT message without a Setter")
		return
	}

	req := request{}
	if err := json.Unmarshal(message.Payload(), &req); err != nil {
		log.Println("Couldn't unmarshal request:", err)
	}

	if req.Power != nil {
		c.setter.SetPower(*req.Power)
	}

	if req.Volume != nil {
		volume := int(*req.Volume * maxVolume)
		c.setter.SetVolume(volume)
	}

	if req.Muted != nil {
		c.setter.SetMute(*req.Muted)
	}

	if req.Source != nil {
		c.setter.SetSource(*req.Source)
	}
}

func (c *connection) publishState(newState state) {
	if newState == c.lastState {
		return
	}
	c.lastState = newState

	ready := newState.powerSet && newState.volumeSet && newState.mutedSet && newState.sourceSet
	if !ready {
		return
	}

	payload, err := json.Marshal(newState)
	if err != nil {
		log.Println("Couldn't marshal JSON:", err)
	}

	token := c.client.Publish(topicState, 0, true, payload)
	token.Wait()
}

func (c *connection) ReportVolume(rawVolume int) {
	newState := c.lastState
	newState.Volume = float64(rawVolume) / maxVolume
	newState.volumeSet = true
	c.publishState(newState)
}

func (c *connection) ReportPower(on bool) {
	newState := c.lastState
	newState.Power = on
	newState.powerSet = true
	c.publishState(newState)
}

func (c *connection) ReportMuted(muted bool) {
	newState := c.lastState
	newState.Muted = muted
	newState.mutedSet = true
	c.publishState(newState)
}

func (c *connection) ReportSource(source string) {
	newState := c.lastState
	newState.Source = source
	newState.sourceSet = true
	c.publishState(newState)
}
