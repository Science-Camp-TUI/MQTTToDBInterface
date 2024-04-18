package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type observationData struct {
	BirdClassId int       `json:"birdClassId"`
	Confidence  float64   `json:"confidence"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Timestamp   time.Time `json:"timestamp"`
}

func (o *observationData) toJSON() string {
	jsonData, err := json.Marshal(*o)
	if err != nil {
		panic(err)
	}
	return string(jsonData)
}

func decodeMessage(payload []byte) *observationData {
	messageString := string(payload)
	dataPattern := regexp.MustCompile("\"data\":\\[([0-9,]*)\\]")
	dataString := dataPattern.FindStringSubmatch(messageString)
	data := strings.Split(dataString[1], ",")
	var dataBytes []byte

	for _, value := range data {
		intByte, _ := strconv.Atoi(value)
		dataBytes = append(dataBytes, byte(intByte))
	}

	var bytesAsHexString []string
	for _, dataByte := range dataBytes {
		bytesAsHexString = append(bytesAsHexString, hex.EncodeToString([]byte{dataByte}))
	}

	fmt.Println(bytesAsHexString)

	return &observationData{
		BirdClassId: int(binary.LittleEndian.Uint16(dataBytes[0:2])),
		Confidence:  float64(math.Float32frombits(binary.LittleEndian.Uint32(dataBytes[4:8]))),
		Latitude:    float64(math.Float32frombits(binary.LittleEndian.Uint32(dataBytes[8:12]))),
		Longitude:   float64(math.Float32frombits(binary.LittleEndian.Uint32(dataBytes[12:16]))),
		Timestamp:   time.Unix(int64(math.Ceil(math.Float64frombits(binary.LittleEndian.Uint64(dataBytes[16:24])))), 0),
	}
}
