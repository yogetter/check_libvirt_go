package main

import (
	"encoding/json"
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"os"
	"time"
)

type db struct {
	Url      string
	Db       string
	Username string
	Password string
}

var Hostname string

func (d *db) init() {
	//read config
	file, _ := os.Open("db_conf.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(d)
	checkError(err)
	Hostname, err = os.Hostname()
	checkError(err)
	log.Println("DB URL:", d.Url)
	log.Println("DB Name:", d.Db)
	log.Println("DB Username:", d.Username)
	log.Println("DB Password:", d.Password)
	file.Close()
}

func (d *db) insertVmInfo(VM instance) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     d.Url,
		Username: d.Username,
		Password: d.Password,
	})
	checkError(err)
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  d.Db,
		Precision: "s",
	})
	checkError(err)
	// Create a point and add to batch
	tags := map[string]string{"uuid": VM.Id, "Hostname": Hostname}
	fields := map[string]interface{}{
		"Total":    VM.Total,
		"Used":     VM.Used,
		"UnUsed":   VM.UnUsed,
		"CpuUsage": VM.CpuUsage,
		"Rx":       VM.InBytes,
		"Tx":       VM.OutBytes,
	}
	log.Println("Send VM information:", tags, fields)
	pt, err := client.NewPoint("vm_usage", tags, fields, time.Now())
	checkError(err)
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
	c.Close()
}
