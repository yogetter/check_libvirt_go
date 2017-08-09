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

func getVmStats(VM *instance, dom *libvirt.Domain) {
	VM.dom = dom
	VM.getName()
	VM.getNicDev()
	VM.getBlockDev()
	VM.setMemValue()
	VM.setCpuValue(CpuCore)
	VM.setBlockStats()
	VM.setInterfaceValue()
}

func start() {
	conn, err := libvirt.NewConnect("qemu:///system")
	checkError(err)
	doms := refreshDomain(conn)
	VMs := make([]instance, len(doms))
	tmp := make([]instance, len(doms))
	for i, dom := range doms {
		getVmStats(&VMs[i], &dom)
		tmp[i] = VMs[i]
		VMs[i].dom.Free()
	}
	time.Sleep(60 * time.Second)
	doms = refreshDomain(conn)
	for i, dom := range doms {
		getVmStats(&VMs[i], &dom)
		VMs[i].setAllValue(tmp[i], CpuCore)
		//VMs[i].getValue()
		VMs[i].dom.Free()
		influx.insertVmInfo(VMs[i])
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
