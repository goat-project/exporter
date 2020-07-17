package parse

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//  Tests:
//  correct JSON file					x
//  wrong file type						x
//  wrong content:
//         - missing required item		x
//         - missing optional item		x
//         - missing string value		x
//         - missing *string value		x
//         - missing int64 value		x
//         - missing int value			x
//         - missing byte value			x
//         - unknown content			x
//         - empty						x

var _ = Describe("IP record parser tests", func() {
	dirPath := "test-data/ip/"

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
				fileName = "0000_correctJSON_20"
			})

			It("should not return an error", func() {
				data, err := IPRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Ips)).To(Equal(20))
			})
		})

		Context("when file has another mime type", func() {
			BeforeEach(func() {
				fileName = "0001_XML"
			})

			It("should return an error", func() {
				data, err := IPRecords(file)

				Expect(err).To(HaveOccurred())
				Expect(len(data.Ips)).To(Equal(0))
			})
		})
	})

	Describe("file content", func() {
		Context("when required item (int64) is missing", func() {
			BeforeEach(func() {
				fileName = "0002_missing_MeasurementTime_int64"
			})

			It("should not return an error", func() {
				data, err := IPRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Ips)).To(Equal(1))
			})
		})

		Context("when optional item (*string) is missing", func() {
			BeforeEach(func() {
				fileName = "0003_missing_CloudComputeService_pstring"
			})

			It("should not return an error", func() {
				data, err := IPRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Ips)).To(Equal(1))
			})
		})

		Context("when item (string) is missing", func() {
			BeforeEach(func() {
				fileName = "0004_missing_LocalUser_string"
			})

			It("should not return an error", func() {
				data, err := IPRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Ips)).To(Equal(1))
			})
		})

		Context("when item (int) is missing", func() {
			BeforeEach(func() {
				fileName = "0005_missing_IPCount_int"
			})

			It("should not return an error", func() {
				data, err := IPRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Ips)).To(Equal(1))
			})
		})

		Context("when item (byte) is missing", func() {
			BeforeEach(func() {
				fileName = "0006_missing_IPVersion_byte"
			})

			It("should not return an error", func() {
				data, err := IPRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.Ips)).To(Equal(1))
			})
		})

		Context("when content is not matching IP structure", func() {
			BeforeEach(func() {
				fileName = "0007_many_json_items"
			})

			It("should return an error", func() {
				data, err := IPRecords(file)

				Expect(err).To(HaveOccurred())
				Expect(data.Ips).To(BeNil())
			})
		})

		Context("when content is empty", func() {
			BeforeEach(func() {
				fileName = "0008_empty"
			})

			It("should return an error", func() {
				data, err := IPRecords(file)

				Expect(err).To(HaveOccurred()) // EOF
				Expect(data.Ips).To(BeNil())
			})
		})

		Context("when JSON format is wrong", func() {
			BeforeEach(func() {
				fileName = "0009_wrong_format"
			})

			It("should return an error", func() {
				data, err := IPRecords(file)

				Expect(err).To(HaveOccurred())
				Expect(data.Ips).To(BeNil())
			})
		})
	})
})
