package main

import (
	"fmt"
	libvirt "github.com/libvirt/libvirt-go"
	"strings"
)

type instance struct {
	Id        string
	Name      string
	MemTotal  int64
	MemUnUsed int64
	MemUsed   int64
	CpuUsage  float32
	CpuTime   uint64
	InBytes   int64
	OutBytes  int64
	NicDevice string
	BkTotal   []int64
	BkDevice  []string
	BkWBytes  []int64
	BkRBytes  []int64
	dom       *libvirt.Domain
}

func (s *instance) getName() {
	xml, err := s.dom.GetXMLDesc(1)
	checkError(err)
	tmp := strings.SplitAfter(xml, "<nova:name>")[1]
	s.Name = strings.Split(tmp, "</nova:name>")[0]
}

func (s *instance) getBlockDev() {
	xml, err := s.dom.GetXMLDesc(1)
	checkError(err)
	//Get HDD
	blk_devs := strings.Count(xml, "'vd")
	s.BkDevice = make([]string, blk_devs)
	for i := 0; i < blk_devs; i++ {
		s.BkDevice[i] = "vd" + string(i+97)
	}
	s.BkWBytes = make([]int64, len(s.BkDevice), len(s.BkDevice))
	s.BkRBytes = make([]int64, len(s.BkDevice), len(s.BkDevice))
	s.BkTotal = make([]int64, len(s.BkDevice), len(s.BkDevice))
}

func (s *instance) getNicDev() {
	xml, err := s.dom.GetXMLDesc(1)
	checkError(err)
	//Get Nic
	tmp := strings.SplitAfter(xml, "<interface type='bridge'>")[1]
	tmp = strings.SplitAfter(tmp, "<target dev='")[1]
	tmp = strings.Split(tmp, "'")[1]
 	s.NicDevice = strings.Split(tmp, "'")[1]
}

func (s *instance) setBlockStats() {
	i := 0
	for _, dev := range s.BkDevice {
		stats, err := s.dom.BlockStats(dev)
		checkError(err)
		info, err := s.dom.GetBlockInfo(dev, 0)
		checkError(err)
		s.BkWBytes[i] = stats.WrBytes
		s.BkRBytes[i] = stats.RdBytes
		s.BkTotal[i] = int64(info.Capacity)
		i++

	}
}

func (s *instance) setCpuValue(CpuCore int) {
	info, err := s.dom.GetInfo()
	checkError(err)
	s.CpuTime = info.CpuTime
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

func (s *instance) setInterfaceValue() {
	ifstat, err := s.dom.InterfaceStats(s.NicDevice)
	checkError(err)
	s.InBytes = ifstat.RxBytes
	s.OutBytes = ifstat.TxBytes
}

func (s instance) getValue() {
	fmt.Println("VM：")
	fmt.Println("Uuid: ", s.Id)
	fmt.Println("Name: ", s.Name)
	fmt.Println("Total: ", s.MemTotal)
	fmt.Println("Used: ", s.MemUsed)
	fmt.Println("UnUsed: ", s.MemUnUsed)
	fmt.Println("CPU: ", s.CpuUsage)
	fmt.Println("WrBytes: ", s.BkWBytes)
	fmt.Println("BkDevice: ", s.BkDevice)
	fmt.Println("BkTotal: ", s.BkTotal)
}

func (s *instance) setAllValue(tmp instance, CpuCore int) {
	usedTime := (s.CpuTime - tmp.CpuTime) / 1000
	s.CpuUsage = float32(usedTime) / float32((60 * 1000000 * CpuCore))
	s.CpuUsage *= 100
	s.InBytes = (s.InBytes - tmp.InBytes) / 60
	s.OutBytes = (s.OutBytes - tmp.OutBytes) / 60
	for i := 0; i < len(s.BkDevice); i++ {
		s.BkWBytes[i] = s.BkWBytes[i] - tmp.BkWBytes[i]
		s.BkRBytes[i] = s.BkRBytes[i] - tmp.BkRBytes[i]

	}
}
