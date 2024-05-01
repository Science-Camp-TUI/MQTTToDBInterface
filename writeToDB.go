package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func writeData(data *observationData, configData *config) {
	token := configData.InfluxDBToken
	url := configData.InfluxDBUrl
	client := influxdb2.NewClient(url, token)

	org := configData.InfluxDBOrga
	bucket := configData.InfluxDBBucket
	writeAPI := client.WriteAPIBlocking(org, bucket)

	tags := map[string]string{
		"birdId": fmt.Sprintf("%d", data.BirdClassId),
	}
	fields := map[string]interface{}{
		"confidence": data.Confidence,
		"latitude":   data.Latitude,
		"longitude":  data.Longitude,
	}

	dataPoint := write.NewPoint("BirdObservation", tags, fields, data.Timestamp)
	if err := writeAPI.WritePoint(context.Background(), dataPoint); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Data written successfully")
	}
}
