package main

import (
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"time"
)

type db struct {
	db       string
	username string
	password string
}

func (d *db) init() {
	d.db = "openstack_vm"
	d.username = "admin"
	d.password = "admin"

}

func (d *db) insertVmInfo(VM instance) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://172.22.131.70:8086",
		Username: d.username,
		Password: d.password,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  d.db,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{"uuid": VM.Id}
	fields := map[string]interface{}{
		"Total":  VM.Total,
		"Used":   VM.Used,
		"UnUsed": VM.UnUsed,
	}

	pt, err := client.NewPoint("mem_usage", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
}
