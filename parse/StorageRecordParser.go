package parse

import (
	"encoding/xml"
	"io"
	"io/ioutil"

	"github.com/goat-project/exporter/record"
)

// StorageRecords parses data from XML format to storage record.
func StorageRecords(file io.Reader) (record.Storages, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return record.Storages{}, err
	}

	var stRecords record.Storages
	err = xml.Unmarshal(data, &stRecords)

	return stRecords, err
}
