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
	hostResults, portResults, usernameResults, passwordResults := hostPattern.FindStringSubmatch(fileString), portPattern.FindStringSubmatch(fileString), usernamePattern.FindStringSubmatch(fileString), passwordPattern.FindStringSubmatch(fileString)

	topicPattern := regexp.MustCompile(`mqtt_topics\s*=\s*\[\s*([^]]*)\s*]`)
	topicResults := topicPattern.FindStringSubmatch(fileString)
	singleTopicPattern := regexp.MustCompile("\"(.+)\"")
	singleTopicsResult := singleTopicPattern.FindAllStringSubmatch(topicResults[1], -1)

	var topics []string
	host, port, username, password := "127.0.0.1", 8883, "", ""

	if hostResults != nil {
		host = hostResults[1]
	}
	if portResults != nil {
		port, err = strconv.Atoi(portResults[1])
		if err != nil {
			panic(err)
		}
	}
	if usernameResults != nil {
		username = usernameResults[1]
	}
	if passwordResults != nil {
		password = passwordResults[1]
	}
	for _, singleTopic := range singleTopicsResult {
		topics = append(topics, singleTopic[1])
	}
	return &config{
		mqttHost:     host,
		mqttPort:     port,
		mqttUsername: username,
		mqttPassword: password,
		mqttTopics:   topics,
	}
}
