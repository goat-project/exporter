package export

import (
	"github.com/goat-project/exporter/gauge"
	"github.com/goat-project/exporter/record"
	"github.com/sirupsen/logrus"
)

// Exporter receives records in record channel and exports them using a given gauge.
type Exporter struct {
	RecordChan chan record.Record
	Gauge      *gauge.Gauge
}

// CreateExporter creates exported with record channel and gauges.
func CreateExporter(recordChan chan record.Record, gauge *gauge.Gauge) *Exporter {
	return &Exporter{
		RecordChan: recordChan,
		Gauge:      gauge,
	}
}

// Export exports records based on their type.
func (e Exporter) Export(finished chan bool) {
	for records := range e.RecordChan {
		switch records.(type) {
		case record.IPs:
			e.Gauge.IPGauge.Export(records)
		case record.Storages:
			e.Gauge.StorageGauge.Export(records)
		case record.VMs:
			e.Gauge.VMGauge.Export(records)
		default:
			logrus.Error("unable to export, unknown record type")
		}
	}

	logrus.Info("export finished")
	finished <- true
}
