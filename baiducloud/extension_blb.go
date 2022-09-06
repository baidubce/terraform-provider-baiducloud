package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/blb"
)

const (
	BLBTCP   = "TCP"
	BLBUDP   = "UDP"
	BLBHTTP  = "HTTP"
	BLBHTTPS = "HTTPS"
	BLBSSL   = "SSL"
)

var BLBProcessingStatus = []string{
	string(blb.BLBStatusCreating),
	string(blb.BLBStatusUpdating),
}

var BLBAvailableStatus = []string{
	string(blb.BLBStatusAvailable),
}

var BLBFailedStatus = []string{
	string(blb.BLBStatusUnavailable),
	string(blb.BLBStatusPaused),
}

var TransportProtocolBLB = []string{TCP, UDP, SSL}
