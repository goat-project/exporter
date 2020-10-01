package gauge

import (
	"time"

	"github.com/goat-project/exporter/utils"

	"github.com/goat-project/exporter/record"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// StorageGauge represents storage gauges exported to Prometheus.
type StorageGauge struct {
	Timestamp                 *prometheus.GaugeVec
	CreateTime                *prometheus.GaugeVec
	FileCount                 *prometheus.GaugeVec
	StartTime                 *prometheus.GaugeVec
	EndTime                   *prometheus.GaugeVec
	ResourceCapacityUsed      *prometheus.GaugeVec
	LogicalCapacityUsed       *prometheus.GaugeVec
	ResourceCapacityAllocated *prometheus.GaugeVec
}

// NewStorageGauge creates storage gauge.
func NewStorageGauge() *StorageGauge {
	stg := StorageGauge{}

	stg.Timestamp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "st",
		Name:      "Timestamp",
		Help:      "represents time when the measurements were exported to the Prometheus.",
	},
		[]string{
			"RecordId",
			"StorageSystem",
			"Site",
			"StorageShare",
			"StorageMedia",
			"StorageClass",
			"DirectoryPath",
			"LocalUser",
			"LocalGroup",
			"UserIdentity",
			"Group",
			"GroupAttribute",
			"GroupAttributeType",
		},
	)

	stg.FileCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "st",
		Name:      "FileCount",
		Help:      "represents the number of files.",
	},
		[]string{
			"RecordId",
			"Site",
			"LocalUser",
			"LocalGroup",
			"UserIdentity",
		},
	)

	stg.ResourceCapacityUsed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "st",
		Name:      "ResourceCapacityUsed",
		Help:      "represents the amount of resource capacity used.",
	},
		[]string{
			"RecordId",
			"Site",
			"LocalUser",
			"LocalGroup",
			"UserIdentity",
		},
	)

	stg.LogicalCapacityUsed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "st",
		Name:      "LogicalCapacityUsed",
		Help:      "represents the amount of logical capacity used.",
	},
		[]string{
			"RecordId",
			"Site",
			"LocalUser",
			"LocalGroup",
			"UserIdentity",
		},
	)

	stg.ResourceCapacityAllocated = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "st",
		Name:      "ResourceCapacityAllocated",
		Help:      "represents the amount of resource capacity allocated.",
	},
		[]string{
			"RecordId",
			"Site",
			"LocalUser",
			"LocalGroup",
			"UserIdentity",
		},
	)

	stg.CreateTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "st",
		Name:      "CreateTime",
		Help:      "represents the time when the measurements were recorded.",
	},
		[]string{
			"RecordId",
			"Site",
			"LocalUser",
			"LocalGroup",
			"UserIdentity",
		},
	)

	stg.StartTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "st",
		Name:      "StartTime",
		Help:      "represents the time when the given storage was created/registered.",
	},
		[]string{
			"RecordId",
			"Site",
			"LocalUser",
			"LocalGroup",
			"UserIdentity",
		},
	)

	stg.EndTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "st",
		Name:      "EndTime",
		Help:      "represents the time when the given storage was finished (or recorded).",
	},
		[]string{
			"RecordId",
			"Site",
			"LocalUser",
			"LocalGroup",
			"UserIdentity",
		},
	)

	return &stg
}

// Register registers storage gauge.
func (stg *StorageGauge) Register() {
	gauges := []prometheus.Collector{
		stg.Timestamp,
		stg.CreateTime,
		stg.FileCount,
		stg.StartTime,
		stg.EndTime,
		stg.ResourceCapacityUsed,
		stg.LogicalCapacityUsed,
		stg.ResourceCapacityAllocated,
	}

	prometheus.MustRegister(gauges...)

	logrus.WithField("resource", "st").Debug("gauges registered")
}

// Export exports storage gauges to Prometheus.
func (stg *StorageGauge) Export(rec record.Record) {
	storages := rec.(record.Storages)

	for _, storage := range storages.Storages {
		label := prometheus.Labels{
			"RecordId":     storage.RecordID,
			"Site":         "",
			"LocalUser":    "",
			"LocalGroup":   "",
			"UserIdentity": "",
		}

		if storage.Site != nil {
			label["Site"] = *storage.Site
		}

		if storage.LocalUser != nil {
			label["LocalUser"] = *storage.LocalUser
		}

		if storage.LocalGroup != nil {
			label["LocalGroup"] = *storage.LocalGroup
		}

		if storage.UserIdentity != nil {
			label["UserIdentity"] = *storage.UserIdentity
		}

		stg.Timestamp.With(labelForStorageTimestamp(storage)).Set(float64(time.Now().Unix()))

		stg.CreateTime.With(label).Set(float64(storage.CreateTime.Unix()))

		if storage.FileCount != nil {
			stg.FileCount.With(label).Set(utils.StrToF64(*storage.FileCount))
		}

		stg.StartTime.With(label).Set(float64(storage.StartTime.Unix()))

		stg.EndTime.With(label).Set(float64(storage.EndTime.Unix()))

		stg.ResourceCapacityUsed.With(label).Set(float64(storage.ResourceCapacityUsed))

		if storage.LogicalCapacityUsed != nil {
			stg.LogicalCapacityUsed.With(label).Set(float64(*storage.LogicalCapacityUsed))
		}

		if storage.ResourceCapacityAllocated != nil {
			stg.ResourceCapacityAllocated.With(label).Set(float64(*storage.ResourceCapacityAllocated))
		}

	}
}

func labelForStorageTimestamp(storage record.Storage) prometheus.Labels {
	labels := prometheus.Labels{
		"RecordId":           storage.RecordID,
		"StorageSystem":      storage.StorageSystem,
		"Site":               "",
		"StorageShare":       "",
		"StorageMedia":       "",
		"StorageClass":       "",
		"DirectoryPath":      "",
		"LocalUser":          "",
		"LocalGroup":         "",
		"UserIdentity":       "",
		"Group":              "",
		"GroupAttribute":     "",
		"GroupAttributeType": "",
	}

	if storage.Site != nil {
		labels["Site"] = *storage.Site
	}

	if storage.StorageShare != nil {
		labels["StorageShare"] = *storage.StorageShare
	}

	if storage.StorageMedia != nil {
		labels["StorageMedia"] = *storage.StorageMedia
	}

	if storage.StorageClass != nil {
		labels["StorageClass"] = *storage.StorageClass
	}

	if storage.DirectoryPath != nil {
		labels["DirectoryPath"] = *storage.DirectoryPath
	}

	if storage.LocalUser != nil {
		labels["LocalUser"] = *storage.LocalUser
	}

	if storage.LocalGroup != nil {
		labels["LocalGroup"] = *storage.LocalGroup
	}

	if storage.UserIdentity != nil {
		labels["UserIdentity"] = *storage.UserIdentity
	}

	if storage.Group != nil {
		labels["Group"] = *storage.Group
	}

	if storage.GroupAttribute != nil {
		labels["GroupAttribute"] = *storage.GroupAttribute
	}

	if storage.GroupAttributeType != nil {
		labels["GroupAttributeType"] = *storage.GroupAttributeType
	}

	return labels
}
