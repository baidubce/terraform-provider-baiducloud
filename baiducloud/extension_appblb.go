package baiducloud

import "github.com/baidubce/bce-sdk-go/services/appblb"

const (
	TCP   = "TCP"
	UDP   = "UDP"
	HTTP  = "HTTP"
	HTTPS = "HTTPS"
	SSL   = "SSL"
)

var APPBLBProcessingStatus = []string{
	string(appblb.BLBStatusCreating),
	string(appblb.BLBStatusUpdating),
}

var APPBLBAvailableStatus = []string{
	string(appblb.BLBStatusAvailable),
}

var APPBLBFailedStatus = []string{
	string(appblb.BLBStatusUnavailable),
	string(appblb.BLBStatusPaused),
}

var TransportProtocol = []string{TCP, UDP, SSL}
