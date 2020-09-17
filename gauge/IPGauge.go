package gauge

import (
	"time"

	"github.com/goat-project/exporter/record"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// IPGauge represents IP gauges exported to Prometheus.
type IPGauge struct {
	Timestamp       *prometheus.GaugeVec
	MeasurementTime *prometheus.GaugeVec
	IPCount         *prometheus.GaugeVec
}

// NewIPGauge create new IP gauge.
func NewIPGauge() *IPGauge {
	ipg := IPGauge{}

	ipg.Timestamp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "ip",
		Name:      "Timestamp",
		Help:      "represents time when the measurements were exported to the Prometheus.",
	},
		[]string{
			"SiteName",
			"CloudComputeService",
			"CloudType",
			"LocalUser",
			"LocalGroup",
			"GlobalUserName",
			"FQAN",
			"IPVersion",
		},
	)

	ipg.MeasurementTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "ip",
		Name:      "MeasurementTime",
		Help:      "represents time when the measurements were recorded.",
	},
		[]string{
			"LocalUser",
			"LocalGroup",
			"GlobalUserName",
		},
	)

	ipg.IPCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "ip",
		Name:      "IPCount",
		Help:      "represents the number of IPs owned by a given user.",
	},
		[]string{
			"LocalUser",
			"LocalGroup",
			"GlobalUserName",
		},
	)

	return &ipg
}

// Register registers IP gauge.
func (ipg *IPGauge) Register() {
	gauges := []prometheus.Collector{
		ipg.Timestamp,
		ipg.MeasurementTime,
		ipg.IPCount,
	}

	prometheus.MustRegister(gauges...)

	logrus.WithField("resource", "ip").Debug("gauges registered")
}

// Export exports IP gauges to Prometheus.
func (ipg *IPGauge) Export(rec record.Record) {
	ips := rec.(record.IPs)

	for _, ip := range ips.Ips {
		label := prometheus.Labels{
			"LocalUser":      ip.LocalUser,
			"LocalGroup":     ip.LocalGroup,
			"GlobalUserName": ip.GlobalUserName,
		}

		labelTimestamp := prometheus.Labels{
			"SiteName":            ip.SiteName,
			"CloudType":           ip.CloudType,
			"LocalUser":           ip.LocalUser,
			"LocalGroup":          ip.LocalGroup,
			"GlobalUserName":      ip.GlobalUserName,
			"FQAN":                ip.FQAN,
			"IPVersion":           string(ip.IPVersion),
			"CloudComputeService": "",
		}

		if ip.CloudComputeService != nil {
			labelTimestamp["CloudComputeService"] = *ip.CloudComputeService
		}

		ipg.Timestamp.With(labelTimestamp).Set(float64(time.Now().Unix()))

		ipg.MeasurementTime.With(label).Set(float64(ip.MeasurementTime))

		ipg.IPCount.With(label).Set(float64(ip.IPCount))
	}
}
