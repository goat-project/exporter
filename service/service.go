package service

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/goat-project/exporter/export"
	"github.com/goat-project/exporter/record"

	"github.com/goat-project/exporter/parse"

	"github.com/goat-project/exporter/watch"

	"github.com/fsnotify/fsnotify"

	"github.com/goat-project/exporter/gauge"

	"github.com/goat-project/exporter/constants"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Serve accountable to Prometheus.
func Serve() {
	eventChan := make(chan fsnotify.Event)
	recordChan := make(chan record.Record)

	exportFinished := make(chan bool, 1)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		fmt.Println()
		fmt.Println(sig)

		close(eventChan)
		close(recordChan)

		<-exportFinished
		os.Exit(0)
	}()

	w, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.WithField("error", err).Error("error create watch")
		return
	}

	defer func() {
		err = w.Close()
		if err != nil {
			logrus.WithField("error", err).Error("error close watch")
			return
		}
	}()

	parser := parse.SetParser(eventChan, recordChan)

	watcher := watch.Watcher{
		Watcher:   w,
		EventChan: eventChan,
	}

	err = watcher.AddRootWithSubDirs(viper.GetString(constants.CfgDirectoryPath))
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "dir": viper.GetString(constants.CfgDirectoryPath)}).Error(
			"error add directory")
	}

	gauges := gauge.CreateAll()
	gauges.RegistryAll()

	exporter := export.CreateExporter(recordChan, gauges)

	go parser.Parse()
	go watcher.Watch()
	go exporter.Export(exportFinished)

	http.Handle("/metrics", promhttp.Handler())
	if err = http.ListenAndServe(viper.GetString(constants.CfgPrometheusEndpoint), nil); err != nil {
		logrus.WithFields(logrus.Fields{"error": err,
			"endpoint": viper.GetString(constants.CfgPrometheusEndpoint)}).Fatal("error listen and serve")
	}
}
