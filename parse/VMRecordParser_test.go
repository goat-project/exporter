package parse

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//  Tests:
//  correct APEL file					x
//  wrong file type						x
//  wrong content:
//         - missing required item		x
//         - missing optional item		x
//         - missing string value		x
//         - missing *string value		x
//         - missing uint32 value		x
//         - missing *uint64 value		x
//         - missing *float32 value		x
//         - unknown content			x
//         - empty						x
// 		   - not tuple					x
// 		   - missing value				x
// 		   - missing colon :			x
// 		   - missing name				x
// 		   - missing %%					x

var _ = Describe("VM record parser tests", func() {
	dirPath := "test-data/vm/"

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
				fileName = "0000_correctAPEL_10"
			})

			It("should not return an error", func() {
				data, err := VMRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(10))
			})
		})

		Context("when file has another mime type", func() {
			BeforeEach(func() {
				fileName = "0001_XML"
			})

			It("should return an error", func() {
				data, err := VMRecords(file)

				Expect(err).To(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(0))
			})
		})
	})

	Describe("file content", func() {
		Context("when required item (string) is missing", func() {
			BeforeEach(func() {
				fileName = "0002_missing_VMUUID_string"
			})

			It("should not return an error", func() {
				data, err := VMRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(1))
			})
		})

		Context("when optional item (*string) is missing", func() {
			BeforeEach(func() {
				fileName = "0003_missing_FQAN_pstring"
			})

			It("should not return an error", func() {
				data, err := VMRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(1))
			})
		})

		Context("when item (uint32) is missing", func() {
			BeforeEach(func() {
				fileName = "0004_missing_CPUCount_uint32"
			})

			It("should not return an error", func() {
				data, err := VMRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(1))
			})
		})

		Context("when item (*uint64) is missing", func() {
			BeforeEach(func() {
				fileName = "0005_missing_Disk_puint64"
			})

			It("should not return an error", func() {
				data, err := VMRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(1))
			})
		})

		Context("when item (*float32) is missing", func() {
			BeforeEach(func() {
				fileName = "0006_missing_Benchmark_pfloat32"
			})

			It("should not return an error", func() {
				data, err := VMRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(1))
			})
		})

		Context("when content is not matching VM structure", func() {
			BeforeEach(func() {
				fileName = "0007_many_items"
			})

			It("should not return an error", func() {
				data, err := VMRecords(file)

				Expect(err).NotTo(HaveOccurred()) // unknown items are ignored
				Expect(len(data.VMs)).To(Equal(1))
			})
		})

		Context("when content is empty", func() {
			BeforeEach(func() {
				fileName = "0008_empty_apel"
			})

			It("should return an error", func() {
				data, err := VMRecords(file)

				Expect(err).To(HaveOccurred()) // EOF
				Expect(data.VMs).To(BeNil())
			})
		})
	})

	Describe("APEL format", func() {
		Context("when value is missing", func() {
			BeforeEach(func() {
				fileName = "0009_missing_value"
			})

			It("should not return an error", func() {
				data, err := VMRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(1)) // ignore
			})
		})

		Context("when item name is missing", func() {
			BeforeEach(func() {
				fileName = "0010_missing_item_name"
			})

			It("should not return an error", func() {
				data, err := VMRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(1)) // ignore
			})
		})

		Context("when colon (:) is missing", func() {
			BeforeEach(func() {
				fileName = "0011_missing_colon"
			})

			It("should not return an error", func() {
				data, err := VMRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(1)) // ignore
			})
		})

		Context("when percent symbol is missing", func() {
			BeforeEach(func() {
				fileName = "0012_missing_percent_symbol"
			})

			It("should not return an error", func() {
				data, err := VMRecords(file)

				Expect(err).NotTo(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(0)) // unable to parse last VM
			})
		})

		Context("when APEL header is missing", func() {
			BeforeEach(func() {
				fileName = "0013_missing_APEL_header"
			})

			It("should return an error", func() {
				data, err := VMRecords(file)

				Expect(err).To(HaveOccurred())
				Expect(len(data.VMs)).To(Equal(0)) // unable to parse last VM
			})
		})
	})
})
