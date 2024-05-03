package main

import (
	"crypto/tls"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net"
	"os"
	"os/signal"
	"sync"
)

func core(m *sync.Mutex, cons *[]*net.Conn, configPath string) {
	configData := readConfig(configPath)

	sslConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("ssl://%s:%d", configData.MqttHost, configData.MqttPort))
	opts.SetClientID("BirdDBReader")
	opts.SetUsername(configData.MqttUsername)
	opts.SetPassword(configData.MqttPassword)
	opts.SetTLSConfig(sslConfig)

	mqttMessageHandler := func(client MQTT.Client, msg MQTT.Message) {
		receivedData := decodeMessage(msg.Payload())
		//fmt.Println(receivedData.toJSON())
		writeData(receivedData, configData)

		for i, conn := range *cons {
			_, err := (*conn).Write([]byte(receivedData.toJSON()))
			if err != nil {
				m.Lock()
				*cons = append((*cons)[0:i], (*cons)[i+1:len(*cons)]...)
				fmt.Println("Remove listener")
				m.Unlock()
			}
		}

		//os.Exit(-1)
	}

	opts.SetDefaultPublishHandler(mqttMessageHandler)

	// Create a new MQTT client
	client := MQTT.NewClient(opts)
	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	qos := 2
	if token := client.Subscribe(configData.MqttTopic, byte(qos), mqttMessageHandler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// Infinite loop to receive MQTT messages
	for {
		// Add your message handling logic here
	}
}
