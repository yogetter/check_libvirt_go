package main

import (
	libvirt "github.com/libvirt/libvirt-go"
	"log"
	"runtime"
	"time"
)

func refreshDomain(conn *libvirt.Connect) []libvirt.Domain {
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	checkError(err)
	return doms
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getVmStats(VM *instance, dom *libvirt.Domain, conn *libvirt.Connect, vmQueue chan *instance) {
	domsPoint := make([]*libvirt.Domain, 1)
	domsPoint[0] = dom
	domStats, err := conn.GetAllDomainStats(domsPoint, 0, 0)
	checkError(err)
	VM.dom = dom
	VM.getName()
	VM.setMemValue()
	VM.setCpuValue(domStats[0].Cpu)
	VM.setBlockStats(domStats[0].Block)
	VM.setInterfaceValue(domStats[0].Net)
	vmQueue <- VM
	domStats[0].Domain.Free()
}

func start() {
	vmQueue := make(chan *instance)
	conn, err := libvirt.NewConnect("qemu:///system")
	checkError(err)
	doms := refreshDomain(conn)
	VMs := make([]instance, len(doms))
	tmp := make([]instance, len(doms))
	for i, dom := range doms {
		go getVmStats(&VMs[i], &dom, conn, vmQueue)
		tmp[i] = *<-vmQueue
		VMs[i].dom.Free()
	}
	time.Sleep(60 * time.Second)
	doms = refreshDomain(conn)
	for i, dom := range doms {
		go getVmStats(&VMs[i], &dom, conn, vmQueue)
		VMs[i] = *<-vmQueue
		VMs[i].setAllValue(tmp[i], CpuCore)
		VMs[i].getValue()
		VMs[i].dom.Free()
		//influx.insertVmInfo(VMs[i])
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
	}
}
