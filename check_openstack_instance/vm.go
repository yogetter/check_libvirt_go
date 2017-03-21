package main

import (
	"fmt"
	libvirt "github.com/libvirt/libvirt-go"
	"time"
)

type instance struct {
	Id       string
	Total    uint64
	UnUsed   uint64
	Used     uint64
	CpuUsage float32
	dom      *libvirt.Domain
}

func (s *instance) setCpuValue(CpuCore int, conn *libvirt.Connect) {
	info, err := s.dom.GetInfo()
	if err != nil {
	}
	startTime := info.CpuTime
	s.refreshDomain(conn)
	time.Sleep(10 * time.Second)
	info, err = s.dom.GetInfo()
	if err != nil {
	}
	endTime := info.CpuTime
	usedTime := (endTime - startTime) / 1000
	s.CpuUsage = float32(usedTime) / float32((10 * 1000000 * CpuCore))
	s.CpuUsage *= 100
}

func (s *instance) setMemValue() {
	id, err := s.dom.GetUUIDString()
	if err != nil {
	}
	mem, err := s.dom.MemoryStats(10, 0)
	s.Id = id
	if err != nil {
	} else {
		for _, stat := range mem {
			if stat.Tag == 4 {
				s.UnUsed = stat.Val
			} else if stat.Tag == 6 {
				s.Total = stat.Val
			}
			s.Used = s.Total - s.UnUsed
		}
	}
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
	if err != nil {
	}
	s.dom = dom
}
