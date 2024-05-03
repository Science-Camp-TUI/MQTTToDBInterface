package main

import (
	"encoding/json"
	"os"
)

type config struct {
	MqttHost       string `json:"mqtt_host"`
	MqttPort       int    `json:"mqtt_port"`
	MqttUsername   string `json:"mqtt_username"`
	MqttPassword   string `json:"mqtt_password"`
	MqttTopic      string `json:"mqtt_topic"`
	InfluxDBUrl    string `json:"influxDBUrl"`
	InfluxDBToken  string `json:"influxDBToken"`
	InfluxDBOrga   string `json:"influxDBOrga"`
	InfluxDBBucket string `json:"influxDBBucket"`
}

func readConfig(path string) *config {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var configData config
	err = json.Unmarshal(file, &configData)
	if err != nil {
		panic(err)
	}
	return &configData
}
