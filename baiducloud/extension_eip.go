package baiducloud

const (
	InstanceTypeBCC = "BCC"
	InstanceTypeBLB = "BLB"
	InstanceTypeVPN = "VPN"
	InstanceTypeNAT = "NAT"
)

const (
	EIPStatusCreating    = "creating"
	EIPStatusAvailable   = "available"
	EIPStatusBinded      = "binded"
	EIPStatusBinding     = "binding"
	EIPStatusUnBinding   = "unbinding"
	EIPStatusUpdating    = "updating"
	EIPStatusPaused      = "paused"
	EIPStatusUnavailable = "unavailable"
)

var EIPProcessingStatus = []string{EIPStatusCreating, EIPStatusBinding, EIPStatusUnBinding, EIPStatusUpdating}
var EIPFailedStatus = []string{EIPStatusPaused, EIPStatusUnavailable}
