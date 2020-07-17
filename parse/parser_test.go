package parse

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	"github.com/fsnotify/fsnotify"
	"github.com/goat-project/exporter/record"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Parser Suite")
}

// Tests
//  - file exist                x
//  - file does not exist       x
//  - plaintext - ok/nok        x
//  - xml - ok/nok              x
//  - json - ok/nok             x
//  - another mime type         x
//  - undetectable mime type    x
//  - close closed file         x

var _ = Describe("Record parser tests", func() {
	dirPath := "test-data/"

	var (
		eventChan  chan fsnotify.Event
		recordChan chan record.Record

		parser *Parser

		hook *test.Hook
	)

	JustBeforeEach(func() {
		eventChan = make(chan fsnotify.Event)
		recordChan = make(chan record.Record)

		parser = SetParser(eventChan, recordChan)

		hook = test.NewGlobal()
	})

	AfterEach(func() {
		close(eventChan)
		close(recordChan)

		hook.Reset()
	})

	Describe("parsing file", func() {
		Context("when file is correct XML", func() {
			It("should not return an error", func(done Done) {
				go parser.Parse()

				parser.EventChan <- fsnotify.Event{
					Name: filepath.Join(dirPath, filepath.Clean("st/0000_correctXML_10")),
				}

				rec := <-parser.RecordChan

				storages := rec.(record.Storages)

				Expect(len(storages.Storages)).To(Equal(10))
				close(done)
			}, 0.2)
		})

		Context("when file is correct JSON", func() {
			It("should not return an error", func(done Done) {
				go parser.Parse()

				parser.EventChan <- fsnotify.Event{
					Name: filepath.Join(dirPath, filepath.Clean("ip/0000_correctJSON_20")),
				}

				rec := <-parser.RecordChan

				ips := rec.(record.IPs)

				Expect(len(ips.Ips)).To(Equal(20))
				close(done)
			}, 0.2)
		})

		Context("when file is correct APEL", func() {
			It("should not return an error", func(done Done) {
				go parser.Parse()

				parser.EventChan <- fsnotify.Event{
					Name: filepath.Join(dirPath, filepath.Clean("vm/0000_correctAPEL_10")),
				}

				rec := <-parser.RecordChan

				vms := rec.(record.VMs)

				Expect(len(vms.VMs)).To(Equal(10))
				close(done)
			}, 0.2)
		})

		Context("when file is not correct XML", func() {
			It("should return an error", func(done Done) {
				go parser.Parse()

				parser.EventChan <- fsnotify.Event{
					Name: filepath.Join(dirPath, filepath.Clean("st/0009_wrong_format")),
				}

				for len(hook.Entries) == 0 { // wait while hook is written
				}

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("error parse file"))

				close(done)
			}, 0.2)
		})

		Context("when file is not correct JSON", func() {
			It("should return an error", func(done Done) {
				go parser.Parse()

				parser.EventChan <- fsnotify.Event{
					Name: filepath.Join(dirPath, filepath.Clean("ip/0009_wrong_format")),
				}

				for len(hook.Entries) == 0 { // wait while hook is written
				}

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("error parse file"))

				close(done)
			}, 0.2)
		})

		Context("when file is not correct APEL", func() {
			It("should return an error", func(done Done) {
				go parser.Parse()

				parser.EventChan <- fsnotify.Event{
					Name: filepath.Join(dirPath, filepath.Clean("vm/0013_missing_APEL_header")),
				}

				for len(hook.Entries) == 0 { // wait while hook is written
				}

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("error parse file"))

				close(done)
			}, 0.2)
		})

		Context("when file does not exist", func() {
			It("should return an error", func(done Done) {
				go parser.Parse()

				parser.EventChan <- fsnotify.Event{
					Name: "asdf",
				}

				for len(hook.Entries) == 0 { // wait while hook is written
				}

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("error open file"))

				close(done)
			}, 0.2)
		})

		Context("when file has another mime type", func() {
			It("should return an error", func(done Done) {
				go parser.Parse()

				parser.EventChan <- fsnotify.Event{
					Name: filepath.Join(dirPath, filepath.Clean("text.csv")),
				}

				for len(hook.Entries) == 0 { // wait while hook is written
				}

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("unknown file type"))

				close(done)
			}, 0.2)
		})

		Context("when file has undetectable mime type", func() {
			It("should return an error", func(done Done) {
				go parser.Parse()

				parser.EventChan <- fsnotify.Event{
					Name: "",
				}

				// TODO find out unsupported mime type

				close(done)
			}, 0.2)
		})
	})

	Describe("closing file", func() {
		Context("when file is closed", func() {
			It("should return an error", func() {
				file, err := os.Open(filepath.Join(dirPath, filepath.Clean("text.csv")))
				Expect(err).NotTo(HaveOccurred())

				err = file.Close()
				Expect(err).NotTo(HaveOccurred())

				closeFile(file)

				Expect(hook.LastEntry().Level).To(Equal(logrus.ErrorLevel))
				Expect(hook.LastEntry().Message).To(Equal("error close file"))
			})
		})
	})
})
