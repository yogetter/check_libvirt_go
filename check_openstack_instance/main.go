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
	influx := db{}
	influx.init()
	CpuCore := runtime.NumCPU()
	conn, err := libvirt.NewConnect("qemu:///system")
	checkError(err)
	defer conn.Close()

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	checkError(err)
	VMs := make([]instance, len(doms))

	for i, dom := range doms {
		VMs[i].dom = &dom
		VMs[i].setMemValue()
		VMs[i].setCpuValue(CpuCore, conn)
		VMs[i].getDevice()
		VMs[i].setInterfaceValue(conn)
		influx.insertVmInfo(VMs[i])
		dom.Free()
	}
	conn.Close()
}
func main() {
	for {
		log.Println("Start collect VM's information")
		start()
		log.Println("End of collect")
		time.Sleep(60 * time.Second)
	}
}
