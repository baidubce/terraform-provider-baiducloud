package baiducloud

import "fmt"

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

var eipStillInUsed = fmt.Errorf("eip is still in used")

var EIPProcessingStatus = []string{EIPStatusCreating, EIPStatusBinding, EIPStatusUnBinding, EIPStatusUpdating}
var EIPFailedStatus = []string{EIPStatusPaused, EIPStatusUnavailable}
