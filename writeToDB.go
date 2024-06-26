package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func writeData(data *observationData, configData *config, bucket, topic string) {
	token := configData.InfluxDBToken
	url := configData.InfluxDBUrl
	client := influxdb2.NewClient(url, token)

	org := configData.InfluxDBOrga
	writeAPI := client.WriteAPIBlocking(org, bucket)

	tags := map[string]string{
		"birdId":    fmt.Sprintf("%d", data.BirdClassId),
		"latitude":  fmt.Sprintf("%f", data.Latitude),
		"longitude": fmt.Sprintf("%f", data.Longitude),
	}
	fields := map[string]interface{}{
		"confidence": data.Confidence,
	}

	dataPoint := write.NewPoint("BirdObservation", tags, fields, data.Timestamp)
	if err := writeAPI.WritePoint(context.Background(), dataPoint); err != nil {
		fmt.Println(Red+"WriteDB error"+Reset, err)
	} else {
		fmt.Println(fmt.Sprintf("%sData written successfully:%s %s -> %s %s", Green, Reset, topic, bucket, data.toJSON()))
	}
}
