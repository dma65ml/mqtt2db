package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var script string

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	vm.Set("topic", string(msg.Topic()))
	vm.Set("msg", string(msg.Payload()))
	vm.Run(script)
	records := [][]string{}
	if value, err := vm.Get("value"); err == nil {
		value, _ := value.ToString()
		slideValue := strings.Split(value, ",")
		for len(slideValue) >= 4 {
			records = append(records, slideValue[:4])
			slideValue = slideValue[4:]
		}
	}
	log.Printf("Records : %s\n", records)

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func initMQTT(config *Config) mqtt.Client {
	// init MQTT client
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Mqtt.Broker, config.Mqtt.Port))
	opts.SetClientID(config.Mqtt.ClientId)
	opts.SetUsername(config.Mqtt.Username)
	opts.SetPassword(config.Mqtt.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe Client
	sub(client, config.Mqtt.Topic)

	script = config.Mqtt.Script

	return client
}
