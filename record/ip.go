package record

// IP represents parsed IP record.
type IP struct {
	MeasurementTime     int64
	SiteName            string
	CloudComputeService *string
	CloudType           string
	LocalUser           string
	LocalGroup          string
	GlobalUserName      string
	FQAN                string
	IPVersion           byte
	IPCount             int
}

// IPs represents Ips structure parsed from JSON where IP records are wrapped.
type IPs struct {
	Ips []IP
}
