package record

// VM represents parsed vm/server record.
type VM struct {
	VMUUID              string
	SiteName            string
	CloudComputeService *string
	MachineName         string
	LocalUserID         *string
	LocalGroupID        *string
	GlobalUserName      *string
	Fqan                *string
	Status              *string
	StartTime           *string
	EndTime             *string
	SuspendDuration     *string
	WallDuration        *string
	CPUDuration         *string
	CPUCount            uint32
	NetworkType         *string
	NetworkInbound      *uint64
	NetworkOutbound     *uint64
	PublicIPCount       *uint64
	Memory              *uint64
	Disk                *uint64
	BenchmarkType       *string
	Benchmark           *float32
	StorageRecordID     *string
	ImageID             *string
	CloudType           *string
}

// VMs represents vms structure parsed from APEL template where virtual machine/server records are wrapped.
type VMs struct {
	VMs []VM
}
