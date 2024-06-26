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

	qos := 1
	for _, topic := range configData.MqttTopic {
		mqttMessageHandler := genHandleFunc(topic.Topic, topic.Bucket, configData)

		if token := client.Subscribe(topic.Topic, byte(qos), mqttMessageHandler); token.Wait() && token.Error() != nil {
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

func genHandleFunc(topic, bucket string, configData *config) func(client MQTT.Client, msg MQTT.Message) {
	return func(client MQTT.Client, msg MQTT.Message) {
		receivedData := decodeMessage(msg.Payload(), topic)
		if receivedData == nil {
			return
		}
		//fmt.Println(receivedData.toJSON())
		writeData(receivedData, configData, bucket, topic)

		//fmt.Println(receivedData.toJSON())

		//os.Exit(-1)
	}
}
