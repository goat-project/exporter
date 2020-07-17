package parse

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/goat-project/exporter/record"
)

// IPRecords parses data from JSON format to IPs.
func IPRecords(file io.Reader) (record.IPs, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return record.IPs{}, err
	}

	var ipRecords record.IPs
	err = json.Unmarshal(data, &ipRecords)

	return ipRecords, err
}
