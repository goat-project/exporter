package gauge

import (
	"fmt"
	"time"

	"github.com/goat-project/exporter/utils"

	"github.com/goat-project/exporter/record"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// VMGauge represents virtual machine/server gauges exported to Prometheus.
type VMGauge struct {
	Timestamp       *prometheus.GaugeVec
	StartTime       *prometheus.GaugeVec
	EndTime         *prometheus.GaugeVec
	SuspendDuration *prometheus.GaugeVec
	WallDuration    *prometheus.GaugeVec
	CPUDuration     *prometheus.GaugeVec
	CPUCount        *prometheus.GaugeVec
	NetworkInbound  *prometheus.GaugeVec
	NetworkOutbound *prometheus.GaugeVec
	PublicIPCount   *prometheus.GaugeVec
	Memory          *prometheus.GaugeVec
	Disk            *prometheus.GaugeVec
}

// NewVMGauge creates new vm/server gauge.
func NewVMGauge() *VMGauge {
	vmg := VMGauge{}

	vmg.Timestamp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "Timestamp",
		Help:      "represents time when the measurements were exported to the Prometheus.",
	},
		[]string{
			"VMUUID",
			"SiteName",
			"CloudComputeService",
			"MachineName",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
			"FQAN",
			"Status",
			"Benchmark",
			"BenchmarkType",
			"StorageRecordId",
			"ImageId",
			"CloudType",
		},
	)

	vmg.StartTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "StartTime",
		Help:      "represents the time when the given virtual machine/server was started.",
	},
		[]string{
			"VMUUID",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
		},
	)

	vmg.EndTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "EndTime",
		Help:      "represents the time when the given virtual machine/server was finished (or recorded).",
	},
		[]string{
			"VMUUID",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
		},
	)

	vmg.SuspendDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "SuspendDuration",
		Help: "represents the time when the given virtual machine/server was suspended. " +
			"The value is counted as END_TIME - START_TIME - WALL_DURATION",
	},
		[]string{
			"VMUUID",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
		},
	)

	vmg.WallDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "WallDuration",
		Help:      "represents the time when the given virtual machine/server was running.",
	},
		[]string{
			"VMUUID",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
		},
	)

	vmg.CPUDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "CPUDuration",
		Help:      "represents the time when the given CPU was running. Same as WallDuration.",
	},
		[]string{
			"VMUUID",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
		},
	)

	vmg.CPUCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "CPUCount",
		Help:      "represents the number of CPUs.",
	},
		[]string{
			"VMUUID",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
		},
	)

	vmg.NetworkInbound = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "NetworkInbound",
		Help:      "represents network inbound.",
	},
		[]string{
			"VMUUID",
			"NetworkType",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
		},
	)

	vmg.NetworkOutbound = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "NetworkOutbound",
		Help:      "represents network outbound.",
	},
		[]string{
			"VMUUID",
			"NetworkType",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
		},
	)

	vmg.PublicIPCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "PublicIPCount",
		Help:      "represents the number of used public IPs.",
	},
		[]string{
			"VMUUID",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
		},
	)

	vmg.Memory = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "Memory",
		Help:      "represents the size of memory.",
	},
		[]string{
			"VMUUID",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
		},
	)

	vmg.Disk = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "vm",
		Name:      "Disk",
		Help:      "represents the size of disks.",
	},
		[]string{
			"VMUUID",
			"LocalUserID",
			"LocalGroupID",
			"GlobalUserName",
		},
	)

	return &vmg
}

// Register registers vm/server gauge.
func (vmg *VMGauge) Register() {
	gauges := []prometheus.Collector{
		vmg.Timestamp,
		vmg.StartTime,
		vmg.EndTime,
		vmg.SuspendDuration,
		vmg.WallDuration,
		vmg.CPUDuration,
		vmg.CPUCount,
		vmg.NetworkInbound,
		vmg.NetworkOutbound,
		vmg.PublicIPCount,
		vmg.Memory,
		vmg.Disk,
	}

	prometheus.MustRegister(gauges...)

	logrus.WithField("resource", "vm").Debug("gauges registered")
}

// Export exports vm/server gauges to Prometheus.
func (vmg *VMGauge) Export(rec record.Record) {
	vms := rec.(record.VMs)

	for _, vm := range vms.VMs {
		vmg.Timestamp.With(labelForVMTimestamp(vm)).Set(float64(time.Now().Unix()))

		labelNework := prometheus.Labels{
			"VMUUID":      vm.VMUUID,
			"NetworkType": "",
		}

		if vm.NetworkType != nil {
			labelNework["NetworkType"] = *vm.NetworkType
		}

		if vm.StartTime != nil {
			vmg.StartTime.With(labelForVM(vm)).Set(utils.StrToF64(*vm.StartTime))
		}

		if vmg.EndTime != nil {
			vmg.EndTime.With(labelForVM(vm)).Set(utils.StrToF64(*vm.EndTime))
		}

		if vm.SuspendDuration != nil {
			vmg.SuspendDuration.With(labelForVM(vm)).Set(utils.StrToF64(*vm.SuspendDuration))
		}

		if vm.WallDuration != nil {
			vmg.WallDuration.With(labelForVM(vm)).Set(utils.StrToF64(*vm.WallDuration))
		}

		if vm.CPUDuration != nil {
			vmg.CPUDuration.With(labelForVM(vm)).Set(utils.StrToF64(*vm.CPUDuration))
		}

		vmg.CPUCount.With(labelForVM(vm)).Set(float64(vm.CPUCount))

		if vm.NetworkInbound != nil {
			vmg.NetworkInbound.With(labelNework).Set(float64(*vm.NetworkInbound))
		}

		if vm.NetworkOutbound != nil {
			vmg.NetworkOutbound.With(labelNework).Set(float64(*vm.NetworkOutbound))
		}

		if vm.PublicIPCount != nil {
			vmg.PublicIPCount.With(labelForVM(vm)).Set(float64(*vm.PublicIPCount))
		}

		if vm.Memory != nil {
			vmg.Memory.With(labelForVM(vm)).Set(float64(*vm.Memory))
		}

		if vm.Disk != nil {
			vmg.Disk.With(labelForVM(vm)).Set(float64(*vm.Disk))
		}
	}
}

func labelForVMTimestamp(vm record.VM) prometheus.Labels {
	labels := prometheus.Labels{
		"VMUUID":              vm.VMUUID,
		"SiteName":            vm.SiteName,
		"MachineName":         vm.MachineName,
		"CloudComputeService": "",
		"LocalUserID":         "",
		"LocalGroupID":        "",
		"GlobalUserName":      "",
		"FQAN":                "",
		"Status":              "",
		"Benchmark":           "",
		"BenchmarkType":       "",
		"StorageRecordId":     "",
		"ImageId":             "",
		"CloudType":           "",
	}

	if vm.CloudComputeService != nil {
		labels["CloudComputeService"] = *vm.CloudComputeService
	}

	if vm.LocalUserID != nil {
		labels["LocalUserID"] = *vm.LocalUserID
	}

	if vm.LocalGroupID != nil {
		labels["LocalGroupID"] = *vm.LocalGroupID
	}

	if vm.GlobalUserName != nil {
		labels["GlobalUserName"] = *vm.GlobalUserName
	}

	if vm.Fqan != nil {
		labels["FQAN"] = *vm.Fqan
	}

	if vm.Status != nil {
		labels["Status"] = *vm.Status
	}

	if vm.Benchmark != nil {
		labels["Benchmark"] = fmt.Sprintf("%f", *vm.Benchmark)
	}

	if vm.BenchmarkType != nil {
		labels["BenchmarkType"] = *vm.BenchmarkType
	}

	if vm.StorageRecordID != nil {
		labels["StorageRecordId"] = *vm.StorageRecordID
	}

	if vm.ImageID != nil {
		labels["ImageId"] = *vm.ImageID
	}

	if vm.CloudType != nil {
		labels["CloudType"] = *vm.CloudType
	}

	return labels
}

func labelForVM(vm record.VM) prometheus.Labels {
	labels := prometheus.Labels{
		"VMUUID":         vm.VMUUID,
		"LocalUserID":    "",
		"LocalGroupID":   "",
		"GlobalUserName": "",
	}

	if vm.LocalUserID != nil {
		labels["LocalUserID"] = *vm.LocalUserID
	}

	if vm.LocalGroupID != nil {
		labels["LocalGroupID"] = *vm.LocalGroupID
	}

	if vm.GlobalUserName != nil {
		labels["GlobalUserName"] = *vm.GlobalUserName
	}

	return labels
}
