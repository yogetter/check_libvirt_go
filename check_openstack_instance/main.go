package main

import (
	libvirt "github.com/libvirt/libvirt-go"
	"runtime"
)

func checError(err error) {
	if err != nil {
	}
}
func start() {
	influx := db{}
	influx.init()
	CpuCore := runtime.NumCPU()
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
	}
	defer conn.Close()

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
	}
	VMs := make([]instance, len(doms))

	for i, dom := range doms {
		VMs[i].dom = &dom
		VMs[i].setMemValue()
		VMs[i].setCpuValue(CpuCore, conn)
		//VMs[i].getValue()
		influx.insertVmInfo(VMs[i])
		dom.Free()
	}

}
func main() {
	start()
}
