package main

import (
	"fmt"
	libvirt "github.com/libvirt/libvirt-go"
	"strings"
	"time"
)

type instance struct {
	Id       string
	Total    int64
	UnUsed   int64
	Used     int64
	CpuUsage float32
	InBytes  int64
	OutBytes int64
	Device   string
	Hostname string
	dom      *libvirt.Domain
}

func (s *instance) setCpuValue(CpuCore int, conn *libvirt.Connect) {
	info, err := s.dom.GetInfo()
	checkError(err)
	startTime := info.CpuTime
	s.refreshDomain(conn)
	time.Sleep(5 * time.Second)
	info, err = s.dom.GetInfo()
	checkError(err)
	endTime := info.CpuTime
	usedTime := (endTime - startTime) / 1000
	s.CpuUsage = float32(usedTime) / float32((5 * 1000000 * CpuCore))
	s.CpuUsage *= 100
}

func (s *instance) setMemValue() {
	id, err := s.dom.GetUUIDString()
	checkError(err)
	mem, err := s.dom.MemoryStats(10, 0)
	s.Id = id
	checkError(err)
	for _, stat := range mem {
		if stat.Tag == 4 {
			s.UnUsed = int64(stat.Val * 1024)
		} else if stat.Tag == 6 {
			s.Total = int64(stat.Val * 1024)
		}
		s.Used = s.Total - s.UnUsed
	}
}
func (s *instance) getDevice() {
	xml, err := s.dom.GetXMLDesc(1)
	checkError(err)
	tmp := strings.SplitAfter(xml, "<interface type='bridge'>")[1]
	tmp = strings.SplitAfter(tmp, "<target dev=")[1]
	tmp = strings.SplitAfter(tmp, "'")[1]
	s.Device = tmp[0 : len(tmp)-1]
}
func (s *instance) setInterfaceValue(conn *libvirt.Connect) {
	ifstat, err := s.dom.InterfaceStats(s.Device)
	checkError(err)
	inBefore := ifstat.RxBytes
	outBefore := ifstat.TxBytes
	time.Sleep(1 * time.Second)
	s.refreshDomain(conn)
	ifstat, err = s.dom.InterfaceStats(s.Device)
	checkError(err)
	s.InBytes = ifstat.RxBytes - inBefore
	s.OutBytes = ifstat.TxBytes - outBefore
}
func (s instance) getValue() {
	fmt.Println("VM：")
	fmt.Println("Uuid: ", s.Id)
	fmt.Println("Total: ", s.Total)
	fmt.Println("Used: ", s.Used)
	fmt.Println("UnUsed: ", s.UnUsed)
	fmt.Println("CPU: ", s.CpuUsage)
}
func (s *instance) refreshDomain(conn *libvirt.Connect) {
	dom, err := conn.LookupDomainByUUIDString(s.Id)
	checkError(err)
	s.dom.Free()
	s.dom = dom

}
