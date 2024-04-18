package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func writeData(data *observationData) {
	token := "IRg-qH2cxLCDeFKiNHD1CeUpdCiHhAzl0z5Bmj445rE4qG6kFWgwQwgZN16f4i9mSmBtQlNFuk2umRzTX6JbBg=="
	url := "http://localhost:8086"
	client := influxdb2.NewClient(url, token)

	org := "SciCampTest"
	bucket := "BirdData"
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
	}

	fmt.Println("Data written successfully")
	/*for value := 0; value < 5; value++ {
		tags := map[string]string{
			"tagname1": "tagvalue1",
		}
		fields := map[string]interface{}{
			"field1": value,
		}
		point := write.NewPoint("measurement1", tags, fields, time.Now())
		time.Sleep(1 * time.Millisecond) // separate points by 1 second

		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			fmt.Println(err)
		}
		fmt.Println("done")
	}*/

	/*queryAPI := client.QueryAPI(org)
		query := `from(bucket: "TestData")
	            |> range(start: -10m)
	            |> filter(fn: (r) => r._measurement == "measurement1")`
		results, err := queryAPI.Query(context.Background(), query)
		if err != nil {
			log.Fatal(err)
		}
		for results.Next() {
			fmt.Println(results.Record())
		}
		if err := results.Err(); err != nil {
			log.Fatal(err)
		}*/
}
