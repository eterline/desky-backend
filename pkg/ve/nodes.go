package ve

import (
	"context"

	"github.com/luthermonson/go-proxmox"
)

type Virt interface {
	Shutdown()
}

type VENode struct {
	*proxmox.Node
	context.Context
}

type NodeData struct {
	Name   string
	CPU    NodeCPU
	Memory NodeMem
	Uptime uint64
}

type NodeCPU struct {
	CPUs      int
	Cores     int
	Frequency int
}

type NodeMem struct {
	Total       int64
	UsedPrecent int
}

func Node(session *proxmox.Client, node string, ctx context.Context) (VENode, error) {
	n, err := session.Node(context.Background(), node)
	return VENode{n, ctx}, err
}

func (node *VENode) Data() NodeData {
	return NodeData{
		Name:   node.Name,
		CPU:    cpuInf(node.CPUInfo),
		Memory: memInf(node.Memory),
		Uptime: node.Uptime,
	}
}

func cpuInf(cpu proxmox.CPUInfo) NodeCPU {
	return NodeCPU{
		CPUs:      cpu.CPUs,
		Cores:     cpu.Cores,
		Frequency: int(cpu.MHZ),
	}
}

func memInf(mem proxmox.Memory) NodeMem {
	precent := float32(mem.Used) / float32(mem.Total) * 100
	return NodeMem{
		Total:       int64(mem.Total),
		UsedPrecent: int(precent),
	}
}
