package watch

import (
	"fmt"
	"os"
	"testing"

	"github.com/fsnotify/fsnotify"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// test Watcher
// test no such directory - fatal error
// test empty directory - wait
// test recursive dir adding - num of files
// test events in root - num of files, reactions
// test events in subdirectory - num of files, reactions
// events: add, remove, modify file (rename?)
// test removing of used directory - no change

// TODO finish watcher tests

func TestResources(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Watcher Suite")
}

var _ = Describe("Watcher tests", func() {
	dirPath := "/tmp/goat/service-test"

	var (
		watcher Watcher
		err     error
	)

	JustBeforeEach(func() {
		var w *fsnotify.Watcher
		w, err = fsnotify.NewWatcher()
		if err != nil {
			fmt.Println("Unable to finish test:", err)
			return
		}

		watcher = Watcher{
			Watcher:   w,
			EventChan: make(chan fsnotify.Event),
		}
	})

	AfterEach(func() {
		if err = watcher.Watcher.Close(); err != nil {
			fmt.Println("Unable to close watch:", err)
			return
		}

		if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
			if err = os.RemoveAll(dirPath); err != nil {
				return // report error
			}
		}
	})

	Describe("adding a directory", func() {
		Context("when no such directory", func() {
			BeforeEach(func() {
				if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
					Expect(os.RemoveAll(dirPath)).NotTo(HaveOccurred())
				}

				_, err := os.Stat(dirPath)
				Expect(err).To(HaveOccurred())
			})

			It("should return an error", func() {
				Expect(watcher.AddRootWithSubDirs(dirPath)).To(HaveOccurred())
			})
		})

		Context("when directory is empty", func() {
			BeforeEach(func() {
				Expect(os.MkdirAll(dirPath, 0700)).NotTo(HaveOccurred())
			})

			It("should not return an error", func() {
				Expect(watcher.AddRootWithSubDirs(dirPath)).NotTo(HaveOccurred())
			})
		})

		Context("when directory contains subdirectories", func() {
			BeforeEach(func() {
				Expect(os.MkdirAll(dirPath, 0700)).NotTo(HaveOccurred())
				Expect(os.MkdirAll(dirPath+"/a", 0700)).NotTo(HaveOccurred())
				Expect(os.MkdirAll(dirPath+"/a/a", 0700)).NotTo(HaveOccurred())
				Expect(os.MkdirAll(dirPath+"b", 0700)).NotTo(HaveOccurred())
			})

			It("should not return an error", func() {
				Expect(watcher.AddRootWithSubDirs(dirPath)).NotTo(HaveOccurred())
			})
		})
	})

	Describe("watching directories", func() {
		Context("when no such directory", func() {
			BeforeEach(func() {
				if _, err := os.Stat(dirPath); !os.IsNotExist(err) {
					Expect(os.RemoveAll(dirPath)).NotTo(HaveOccurred())
				}

				_, err := os.Stat(dirPath)
				Expect(err).To(HaveOccurred())
			})

			It("should return an error", func(done Done) {
				go watcher.Watch()

				// no detection for this situation

				close(done)
			}, 0.2)
		})

		Context("when file is modified", func() {
			BeforeEach(func() {
				Expect(os.MkdirAll(dirPath, 0700)).NotTo(HaveOccurred())
			})

			It("should not return an error", func(done Done) {
				Expect(watcher.AddRootWithSubDirs(dirPath)).NotTo(HaveOccurred())

				go watcher.Watch()

				file, err := os.Create(dirPath + "/file.txt")
				Expect(err).NotTo(HaveOccurred())

				_, err = file.Write([]byte("Hello world!"))
				Expect(err).NotTo(HaveOccurred())

				event := <-watcher.EventChan
				Expect(event.Op).To(Equal(fsnotify.Write))
				Expect(event.Name).To(Equal(dirPath + "/file.txt"))

				close(done)
			}, 0.2)
		})
	})
})
