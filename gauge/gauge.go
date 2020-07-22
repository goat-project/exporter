package gauge

// Gauge represents all gauges.
type Gauge struct {
	VMGauge      *VMGauge
	IPGauge      *IPGauge
	StorageGauge *StorageGauge
}

// CreateAll creates all gauges.
func CreateAll() *Gauge {
	return &Gauge{
		VMGauge:      NewVMGauge(),
		IPGauge:      NewIPGauge(),
		StorageGauge: NewStorageGauge(),
	}
}

// RegistryAll registers all gauges.
func (g Gauge) RegistryAll() {
	g.VMGauge.Register()
	g.IPGauge.Register()
	g.StorageGauge.Register()
}
