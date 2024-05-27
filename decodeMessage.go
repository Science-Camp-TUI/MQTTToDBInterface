package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type observationData struct {
	BirdClassId int       `json:"birdClassId"`
	Confidence  float64   `json:"confidence"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Timestamp   time.Time `json:"timestamp"`
}
type mqttMessage struct {
	BaseStations []struct {
		BsEui      int64   `json:"bsEui"`
		EqSnr      float64 `json:"eqSnr"`
		Mode       string  `json:"mode"`
		Profile    string  `json:"profile"`
		Rssi       float64 `json:"rssi"`
		RxTime     int64   `json:"rxTime"`
		Snr        float64 `json:"snr"`
		Subpackets struct {
			Frequency []float64 `json:"frequency"`
			Rssi      []float64 `json:"rssi"`
			Snr       []float64 `json:"snr"`
		} `json:"subpackets"`
	} `json:"baseStations"`
	Cnt         int         `json:"cnt"`
	Components  interface{} `json:"components"`
	Data        []int       `json:"data"`
	DlAck       bool        `json:"dlAck"`
	DlOpen      bool        `json:"dlOpen"`
	Format      int         `json:"format"`
	Meta        interface{} `json:"meta"`
	ResponseExp bool        `json:"responseExp"`
	TypeEui     int         `json:"typeEui"`
}

func (o *observationData) toJSON() string {
	jsonData, err := json.Marshal(*o)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(jsonData)
}

func decodeMessage(payload []byte, topic string) *observationData {
	var receivedData mqttMessage
	err := json.Unmarshal(payload, &receivedData)
	if err != nil {
		fmt.Println(Red+"Unmarshal error"+Reset, err)
		return nil
	}

	data := receivedData.Data
	var dataBytes []byte

	for _, value := range data {
		dataBytes = append(dataBytes, byte(value))
	}

	if len(dataBytes) != 24 {
		fmt.Println(fmt.Sprintf("%sByte number mismatch%s on %s, expected 24 got %d", Red, Reset, topic, len(dataBytes)))
		return nil
	}

	//fmt.Println(bytesAsHexString)

	obs := &observationData{
		BirdClassId: int(binary.LittleEndian.Uint16(dataBytes[0:2])),
		Confidence:  float64(math.Float32frombits(binary.LittleEndian.Uint32(dataBytes[4:8]))),
		Latitude:    float64(math.Float32frombits(binary.LittleEndian.Uint32(dataBytes[8:12]))),
		Longitude:   float64(math.Float32frombits(binary.LittleEndian.Uint32(dataBytes[12:16]))),
		Timestamp:   time.Unix(int64(math.Ceil(math.Float64frombits(binary.LittleEndian.Uint64(dataBytes[16:24])))), 0),
	}

	if obs.Confidence > 1 || obs.Confidence < 0 {
		fmt.Println(fmt.Sprintf("%sConfidence mismatch%s on %s, expected 0<=conf<=1 got %f", Red, Reset, topic, obs.Confidence))
		return nil
	}

	if obs.BirdClassId < 0 || obs.Confidence > 6521 {
		fmt.Println(fmt.Sprintf("%sBirdID mismatch%s on %s, expected 0<=BirdID<=6521 got %d", Red, Reset, topic, obs.BirdClassId))
		return nil
	}

	return obs
}
