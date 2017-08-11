package main

import (
	"fmt"
	libvirt "github.com/libvirt/libvirt-go"
	"strings"
)

type instance struct {
	Id         string
	Name       string
	MemTotal   int64
	MemUnUsed  int64
	MemUsed    int64
	CpuUsage   float32
	CpuTime    uint64
	NetStats   []libvirt.DomainStatsNet
	BlockStats []libvirt.DomainStatsBlock
	dom        *libvirt.Domain
}

func (s *instance) getName() {
	xml, err := s.dom.GetXMLDesc(1)
	checkError(err)
	tmp := strings.SplitAfter(xml, "<nova:name>")[1]
	s.Name = strings.Split(tmp, "</nova:name>")[0]
}

func (s *instance) setBlockStats(Block []libvirt.DomainStatsBlock) {
	s.BlockStats = make([]libvirt.DomainStatsBlock, len(Block))
	s.BlockStats = Block
}

func (s *instance) setCpuValue(Cpu *libvirt.DomainStatsCPU) {
	s.CpuTime = Cpu.Time
}

func (s *instance) setMemValue() {
	id, err := s.dom.GetUUIDString()
	checkError(err)
	mem, err := s.dom.MemoryStats(10, 0)
	s.Id = id
	checkError(err)
	for _, stat := range mem {
		if stat.Tag == 4 {
			s.MemUnUsed = int64(stat.Val * 1024)
		} else if stat.Tag == 6 {
			s.MemTotal = int64(stat.Val * 1024)
		}
		s.MemUsed = s.MemTotal - s.MemUnUsed
	}
}

func (s *instance) setInterfaceValue(Net []libvirt.DomainStatsNet) {
	s.NetStats = make([]libvirt.DomainStatsNet, len(Net))
	s.NetStats = Net
}

func (s instance) getValue() {
	fmt.Println("VM：")
	fmt.Println("Uuid: ", s.Id)
	fmt.Println("Name: ", s.Name)
	fmt.Println("Total: ", s.MemTotal)
	fmt.Println("Used: ", s.MemUsed)
	fmt.Println("UnUsed: ", s.MemUnUsed)
	fmt.Println("CPU: ", s.CpuUsage)
	fmt.Println("BlockStats: ", s.BlockStats)
	fmt.Println("NetStats: ", s.NetStats)
}

func (s *instance) setAllValue(tmp instance, CpuCore int) {
	usedTime := (s.CpuTime - tmp.CpuTime) / 1000
	s.CpuUsage = float32(usedTime) / float32((60 * 1000000 * CpuCore))
	s.CpuUsage *= 100
	s.NetStats[0].RxBytes = (s.NetStats[0].RxBytes - tmp.NetStats[0].RxBytes) / 60
	s.NetStats[0].TxBytes = (s.NetStats[0].TxBytes - tmp.NetStats[0].TxBytes) / 60
	for i := 0; i < len(s.BlockStats); i++ {
		s.BlockStats[i].WrBytes = s.BlockStats[i].WrBytes - tmp.BlockStats[i].WrBytes
		s.BlockStats[i].RdBytes = s.BlockStats[i].RdBytes - tmp.BlockStats[i].RdBytes
	}
}
