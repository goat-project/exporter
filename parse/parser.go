package parse

import (
	"os"

	"github.com/goat-project/exporter/record"

	"github.com/fsnotify/fsnotify"
	"github.com/gabriel-vasile/mimetype"
	"github.com/sirupsen/logrus"
)

// Parser structure with event channel for incoming events and record channel for parsed records.
type Parser struct {
	EventChan  chan fsnotify.Event
	RecordChan chan record.Record
}

const (
	mimePlaintext = "text/plain; charset=utf-8"
	mimeJSON      = "application/json"
	mimeXML       = "text/xml; charset=utf-8"
)

// SetParser sets event and record channels to parser.
func SetParser(eventChan chan fsnotify.Event, recordChan chan record.Record) *Parser {
	return &Parser{
		EventChan:  eventChan,
		RecordChan: recordChan,
	}
}

// Parse takes event from channel, parses content and put to record channel to export to Prometheus.
func (p Parser) Parse() {
	for event := range p.EventChan {
		file, err := os.Open(event.Name)
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err, "file": event.Name}).Error("error open file")
			continue
		}

		mimeType, err := mimetype.DetectFile(file.Name())
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err, "file": event.Name}).Error("error detect mime type")
			continue
		}

		switch mimeType.String() {
		case mimePlaintext:
			vms, err := VMRecords(file)
			if err != nil {
				logrus.WithFields(logrus.Fields{"error": err, "file": file.Name(),
					"type": "vm"}).Error("error parse file")
				closeFile(file)
				continue
			}

			p.RecordChan <- vms
		case mimeJSON:
			ips, err := IPRecords(file)
			if err != nil {
				logrus.WithFields(logrus.Fields{"error": err, "file": file.Name(),
					"type": "ip"}).Error("error parse file")
				closeFile(file)
				continue
			}

			p.RecordChan <- ips
		case mimeXML:
			storages, err := StorageRecords(file)
			if err != nil {
				logrus.WithFields(logrus.Fields{"error": err, "file": file.Name(),
					"type": "st"}).Error("error parse file")
				closeFile(file)
				continue
			}

			p.RecordChan <- storages
		default:
			logrus.WithFields(logrus.Fields{"type": mimeType.String(), "file": file.Name()}).Error("unknown file type")
			closeFile(file)
			continue
		}

		logrus.WithFields(logrus.Fields{"type": mimeType.String(), "file": file.Name()}).Debug("file parsed")
		closeFile(file)
	}
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err, "file": file.Name()}).Error("error close file")
	}
}
