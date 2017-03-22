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

func (d *db) init() {
	//read config
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(d)
	checkError(err)
	log.Println("DB URL:", d.Url)
	log.Println("DB Name:", d.Db)
	log.Println("DB Username:", d.Username)
	log.Println("DB Password:", d.Password)

}

func (d *db) insertVmInfo(VM instance) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
<<<<<<< HEAD
		Addr:     d.Url,
		Username: d.Username,
		Password: d.Password,
=======
		Addr:     "http://localhost:8086",
		Username: d.username,
		Password: d.password,
>>>>>>> 4e50b27e980d1dc6fc72498388b29eeadb959ac0
	})
	checkError(err)
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  d.Db,
		Precision: "s",
	})
	checkError(err)
	// Create a point and add to batch
	tags := map[string]string{"uuid": VM.Id}
	fields := map[string]interface{}{
		"Total":    VM.Total,
		"Used":     VM.Used,
		"UnUsed":   VM.UnUsed,
		"CpuUsage": VM.CpuUsage,
	}
	log.Println("Send VM information:", tags, fields)
	pt, err := client.NewPoint("vm_usage", tags, fields, time.Now())
	checkError(err)
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
}
