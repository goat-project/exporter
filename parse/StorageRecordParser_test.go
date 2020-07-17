package parse

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//  Tests:
//  correct XML file					x
//  wrong file type						x
//  wrong content:
//         - missing required item		x
//         - missing optional item		x
//         - missing string value		x
//         - missing *string value		x
//         - missing time value			x
//         - missing uint64 value		x
//         - missing *uint64 value		x
//         - unknown content			x
//         - empty						x
//         - wrong XML format			x

var _ = Describe("Storage record parser tests", func() {
	dirPath := "test-data/st/"

	var (
		err      error
		file     *os.File
		fileName string
	)

	JustBeforeEach(func() {
		_, err = os.Stat(filepath.Join(dirPath, filepath.Clean(fileName)))
		Expect(err).NotTo(HaveOccurred())

		file, err = os.Open(filepath.Join(dirPath, filepath.Clean(fileName)))
		if err != nil {
			fmt.Println("error open file:", err)
		}
	})

	AfterEach(func() {
		err = file.Close()
		if err != nil {
			fmt.Println("error close file:", err, file.Name())
		}
	})

	Describe("parsing file", func() {
		Context("when file is correct", func() {
			BeforeEach(func() {
				fileName = "0000_correctXML_10"
			})

			It("should not return an error", func() {
				data, err := StorageRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Storages)).To(Equal(10))
			})
		})

		Context("when file has another mime type", func() {
			BeforeEach(func() {
				fileName = "0001_JSON"
			})

			It("should return an error", func() {
				data, err := StorageRecords(file)

				Expect(err).To(HaveOccurred())
				Expect(len(data.Storages)).To(Equal(0))
			})
		})
	})

	Describe("file content", func() {
		Context("when required item (string) is missing", func() {
			BeforeEach(func() {
				fileName = "0002_missing_RECORD_ID_string"
			})

			It("should not return an error", func() {
				data, err := StorageRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Storages)).To(Equal(1))
			})
		})

		Context("when optional item (*string) is missing", func() {
			BeforeEach(func() {
				fileName = "0003_missing_STORAGE_SHARE_pstring"
			})

			It("should not return an error", func() {
				data, err := StorageRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Storages)).To(Equal(1))
			})
		})

		Context("when item (time) is missing", func() {
			BeforeEach(func() {
				fileName = "0004_missing_START_TIME_time"
			})

			It("should not return an error", func() {
				data, err := StorageRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Storages)).To(Equal(1))
			})
		})

		Context("when item (uint64) is missing", func() {
			BeforeEach(func() {
				fileName = "0005_missing_RES_CAP_USED_uint64"
			})

			It("should not return an error", func() {
				data, err := StorageRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Storages)).To(Equal(1))
			})
		})

		Context("when item (*uint64) is missing", func() {
			BeforeEach(func() {
				fileName = "0006_missing_LOG_CAP_USED_puint64"
			})

			It("should not return an error", func() {
				data, err := StorageRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Storages)).To(Equal(1))
			})
		})

		Context("when content is not matching storage structure", func() {
			BeforeEach(func() {
				fileName = "0007_many_xml_items"
			})

			It("should not return an error", func() {
				data, err := StorageRecords(file)

				Expect(err).NotTo(HaveOccurred()) // unknown items are ignored
				Expect(len(data.Storages)).To(Equal(1))
			})
		})

		Context("when content is empty", func() {
			BeforeEach(func() {
				fileName = "0008_empty"
			})

			It("should return an error", func() {
				data, err := StorageRecords(file)

				Expect(err).To(HaveOccurred()) // EOF
				Expect(data.Storages).To(BeNil())
			})
		})

		Context("when XML format is wrong", func() {
			BeforeEach(func() {
				fileName = "0009_wrong_format"
			})

			It("should return an error", func() {
				data, err := StorageRecords(file)

				Expect(err).To(HaveOccurred())
				Expect(data.Storages).To(BeNil())
			})
		})
	})
})
