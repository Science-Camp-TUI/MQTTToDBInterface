package main

import (
	"crypto/tls"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
)

func core(configPath string) {
	configData := readConfig(configPath)

	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("%s://%s:%d", configData.MqttProtocol, configData.MqttHost, configData.MqttPort))
	//opts.SetClientID("BirdDBReader")
	opts.SetUsername(configData.MqttUsername)
	opts.SetPassword(configData.MqttPassword)
	if configData.MqttProtocol == "ssl" {
		sslConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		opts.SetTLSConfig(sslConfig)
	}

	opts.SetDefaultPublishHandler(func(client MQTT.Client, message MQTT.Message) {
		fmt.Printf("Unhandled message: %s\n", string(message.Payload()))
	})

	// Create a new MQTT client
	client := MQTT.NewClient(opts)
	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	qos := 2
	for _, topics := range configData.MqttTopic {

		mqttMessageHandler := func(client MQTT.Client, msg MQTT.Message) {
			receivedData := decodeMessage(msg.Payload(), topics.Topic)
			if receivedData == nil {
				return
			}
			//fmt.Println(receivedData.toJSON())
			writeData(receivedData, configData, topics.Bucket)

			fmt.Println(receivedData.toJSON())

			//os.Exit(-1)
		}

		if token := client.Subscribe(topics.Topic, byte(qos), mqttMessageHandler); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// Infinite loop to receive MQTT messages
	for {
		// Add your message handling logic here
	}
}
