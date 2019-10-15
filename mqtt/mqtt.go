package mqtt

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type connection struct {
	client mqtt.Client
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

	return &connection{
		client: client,
	}, nil
}

func (c *connection) ReportVolume(volume int) {
	token := c.client.Publish("vsx/volume", 0, true, fmt.Sprintf("%d", volume))
	token.Wait()
}

func (c *connection) ReportPower(on bool) {
	token := c.client.Publish("vsx/power", 0, true, fmt.Sprintf("%v", on))
	token.Wait()
}
