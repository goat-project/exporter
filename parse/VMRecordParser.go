package parse

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/goat-project/exporter/utils"
	"github.com/sirupsen/logrus"

	"github.com/goat-project/exporter/record"
)

const (
	apelVersion = "0.4"
	apelMessage = "APEL-cloud-message: v" + apelVersion
)

// VMRecords parses data from template to vm/server record.
func VMRecords(file io.Reader) (record.VMs, error) {
	// create reader
	reader := bufio.NewReader(file)

	// check APEL format
	if !apelFormat(reader) {
		return record.VMs{}, fmt.Errorf("not APEL format")
	}

	// create VMs structure
	var vms record.VMs

	for {
		// read vm (or %)
		vmb, err := reader.ReadString('%')
		if err != nil {
			break // EOF
		}
		if vmb == "%" {
			continue
		}
		if vmb == "" {
			continue
		}

		// split vm lines
		lines := strings.Split(vmb, "\n")
		if len(lines) == 0 {
			continue
		}

		// create vm map
		vmm := make(map[string]string)

		// split line to save vm parts
		for _, line := range lines {
			part := strings.Split(line, ": ")
			if len(part) == 2 {
				vmm[part[0]] = part[1]
			}
		}

		// save vm to VMs structure
		vms.VMs = append(vms.VMs, record.VM{
			VMUUID:              vmm["VMUUID"],
			SiteName:            vmm["SiteName"],
			CloudComputeService: utils.String(vmm["CloudComputeService"]),
			MachineName:         vmm["MachineName"],
			LocalUserID:         utils.String(vmm["LocalUserId"]),
			LocalGroupID:        utils.String(vmm["LocalGroupId"]),
			GlobalUserName:      utils.String(vmm["GlobalUserName"]),
			Fqan:                utils.String(vmm["FQAN"]),
			Status:              utils.String(vmm["Status"]),
			StartTime:           utils.String(vmm["StartTime"]),
			EndTime:             utils.String(vmm["EndTime"]),
			SuspendDuration:     utils.String(vmm["SuspendDuration"]),
			WallDuration:        utils.String(vmm["WallDuration"]),
			CPUDuration:         utils.String(vmm["CpuDuration"]),
			CPUCount:            utils.StrToUint32(vmm["CpuCount"]),
			NetworkType:         utils.String(vmm["NetworkType"]),
			NetworkInbound:      utils.StrToUint64(vmm["NetworkInbound"]),
			NetworkOutbound:     utils.StrToUint64(vmm["NetworkOutbound"]),
			PublicIPCount:       utils.StrToUint64(vmm["PublicIPCount"]),
			Memory:              utils.StrToUint64(vmm["Memory"]),
			Disk:                utils.StrToUint64(vmm["Disk"]),
			StorageRecordID:     utils.String(vmm["StorageRecordId"]),
			ImageID:             utils.String(vmm["ImageId"]),
			CloudType:           utils.String(vmm["CloudType"]),
			BenchmarkType:       utils.String(vmm["BenchmarkType"]),
			Benchmark:           utils.StrToFloat32(vmm["Benchmark"]),
		})
	}

	return vms, nil
}

// apelFormat reads line from reader and checks message.
func apelFormat(reader *bufio.Reader) bool {
	line, err := reader.ReadString('\n')
	if err != nil {
		logrus.WithField("error", err).Error("unable to read first line")
		return false
	}

	return line == apelMessage+"\n"
}
