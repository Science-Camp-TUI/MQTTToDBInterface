package main

import (
	"crypto/tls"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
)

func main() {
	mqttConfig := readConfig("config.py")

	sslConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("ssl://%s:%d", mqttConfig.mqttHost, mqttConfig.mqttPort))
	opts.SetClientID("BirdDBReader")
	opts.SetUsername(mqttConfig.mqttUsername)
	opts.SetPassword(mqttConfig.mqttPassword)
	opts.SetTLSConfig(sslConfig)

	mqttMessageHandler := func(client MQTT.Client, msg MQTT.Message) {
		receivedData := decodeMessage(msg.Payload())
		fmt.Println(receivedData)
		//os.Exit(-1)
	}

	opts.SetDefaultPublishHandler(mqttMessageHandler)

	// Create a new MQTT client
	client := MQTT.NewClient(opts)
	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := mqttConfig.mqttTopics[0]
	qos := 2
	if token := client.Subscribe(topic, byte(qos), nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// Infinite loop to receive MQTT messages
	for {
		// Add your message handling logic here
	}
}
