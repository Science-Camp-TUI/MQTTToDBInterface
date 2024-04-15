package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type observationData struct {
	birdClassId int
	confidence  float64
	latitude    float64
	longitude   float64
	timestamp   time.Time
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
		birdClassId: int(binary.LittleEndian.Uint16(dataBytes[0:2])),
		confidence:  float64(math.Float32frombits(binary.LittleEndian.Uint32(dataBytes[4:8]))),
		latitude:    float64(math.Float32frombits(binary.LittleEndian.Uint32(dataBytes[8:12]))),
		longitude:   float64(math.Float32frombits(binary.LittleEndian.Uint32(dataBytes[12:16]))),
		timestamp:   time.Unix(int64(math.Ceil(math.Float64frombits(binary.LittleEndian.Uint64(dataBytes[16:24])))), 0),
	}
}
