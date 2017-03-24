package main

import (
	libvirt "github.com/libvirt/libvirt-go"
	"log"
	"runtime"
	"time"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func start() {
	conn, err := libvirt.NewConnect("qemu:///system")
	checkError(err)
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	checkError(err)
	VMs := make([]instance, len(doms))

	for i, dom := range doms {
		VMs[i].dom = &dom
		VMs[i].setMemValue()
		VMs[i].setCpuValue(CpuCore, conn)
		VMs[i].getDevice()
		//VMs[i].getValue()
		VMs[i].setInterfaceValue(conn)
		influx.insertVmInfo(VMs[i])
		VMs[i].dom.Free()
	}
	conn.Close()
}

var influx db
var CpuCore int

func main() {
	influx = db{}
	influx.init()
	CpuCore = runtime.NumCPU()

	for {
		log.Println("Start collect VM's information")
		start()
		log.Println("End of collect")
		time.Sleep(60 * time.Second)
	}
}
