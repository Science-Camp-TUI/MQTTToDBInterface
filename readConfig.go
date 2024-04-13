package main

import (
	"os"
	"regexp"
	"strconv"
)

type config struct {
	mqttHost     string
	mqttPort     int
	mqttUsername string
	mqttPassword string
	mqttTopics   []string
}

func readConfig(path string) *config {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	fileString := string(file)
	hostPattern, portPattern, usernamePattern, passwordPattern := regexp.MustCompile("mqtt_host\\s*=\\s*\\\"(.+?)\\\""), regexp.MustCompile(`mqtt_port\s*=\s*(\d+)`), regexp.MustCompile("mqtt_username\\s*=\\s*\\\"(.+?)\\\""), regexp.MustCompile("mqtt_password\\s*=\\s*\\\"(.+?)\\\"")
	hostResults, portResults, usernameResults, passwordResults := hostPattern.FindAllStringSubmatch(fileString, -1), portPattern.FindAllStringSubmatch(fileString, -1), usernamePattern.FindAllStringSubmatch(fileString, -1), passwordPattern.FindAllStringSubmatch(fileString, -1)

	host, port, username, password := "127.0.0.1", 8883, "", ""

	if hostResults != nil {
		host = hostResults[0][1]
	}
	if portResults != nil {
		port, err = strconv.Atoi(portResults[0][1])
		if err != nil {
			panic(err)
		}
	}
	if usernameResults != nil {
		username = usernameResults[0][1]
	}
	if passwordResults != nil {
		password = passwordResults[0][1]
	}

	return &config{
		mqttHost:     host,
		mqttPort:     port,
		mqttUsername: username,
		mqttPassword: password,
		mqttTopics:   nil,
	}
}
