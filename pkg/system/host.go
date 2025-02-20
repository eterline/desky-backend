package system

import (
	"context"
	"time"

	cpuPs "github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/sensors"

	hostPs "github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

type SystemService struct {
	ctx       context.Context
	CancelCtx context.CancelFunc
}

func New() *SystemService {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &SystemService{
		ctx:       ctx,
		CancelCtx: cancel,
	}
}

func (hs *SystemService) RAMInfo() (ram *RAMInfo) {

	ram = &RAMInfo{}

	stat, err := mem.VirtualMemoryWithContext(hs.ctx)
	if err == nil {
		ram = &RAMInfo{
			Total:      stat.Total,
			Used:       stat.Used,
			Avail:      stat.Available,
			UsePercent: stat.UsedPercent,
		}
	}

	return ram
}

func (hs *SystemService) HostInfo() (host *HostInfo) {

	host = &HostInfo{}

	info, err := hostPs.InfoWithContext(hs.ctx)
	if err == nil {
		host = &HostInfo{
			Name:         info.Hostname,
			Uptime:       Uptime(),
			OS:           info.OS,
			ProcessCount: info.Procs,
			VirtSystem:   info.VirtualizationSystem,
			Addrs:        HostAddrs(),
		}
	}

	return host
}

func (hs *SystemService) CPUInfo() (cpu *CPUInfo) {

	cpu = &CPUInfo{}

	stats, err := cpuPs.InfoWithContext(hs.ctx)
	if err == nil {

		load, _ := cpuPs.PercentWithContext(hs.ctx, 1*time.Second, false)

		cpu.Name = stats[0].ModelName
		cpu.Model = stats[0].Model
		cpu.Cache = stats[0].CacheSize
		cpu.Load = load[0]

		var found bool

		for _, core := range stats {

			found = false

			for _, c := range cpu.Cores {
				if c.ID == core.CoreID {
					found = true
					break
				}
			}

			if !found {
				cpu.Cores = append(cpu.Cores, CpuCore{
					ID:      core.CoreID,
					FreqMhz: core.Mhz,
				})
			}

		}

		cpu.ThreadCount = uint64(len(stats))
		cpu.CoreCount = uint64(len(cpu.Cores))
	}

	return cpu
}

func (hs *SystemService) Temperatures() (data []SensorInfo) {

	stat, err := sensors.TemperaturesWithContext(hs.ctx)
	if err == nil {

		data = make([]SensorInfo, len(stat))

		for i, sens := range stat {
			data[i] = SensorInfo{
				Key:     sens.SensorKey,
				Current: sens.Temperature,
				Max:     sens.High,
			}
		}

	}

	return data
}

func (hs *SystemService) Load() (avg *AverageLoad) {

	avg = &AverageLoad{}

	res, err := load.AvgWithContext(hs.ctx)
	if err == nil {
		avg = &AverageLoad{
			Load1:  res.Load1,
			Load5:  res.Load5,
			Load15: res.Load15,
		}
	}

	return avg
}
