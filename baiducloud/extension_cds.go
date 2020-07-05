package baiducloud

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/bcc/api"
)

const CDSNotAttachedErrorCode = "Volume.DiskNotAttachedInstance"

var CDSFailedStatus = []string{
	string(api.VolumeStatusDELETING),
	string(api.VolumeStatusDELETED),
	string(api.VolumeStatusERROR),
	string(api.VolumeStatusEXPIRED),
	string(api.VolumeStatusNOTAVAILABLE),
}

var CDSProcessingStatus = []string{
	string(api.VolumeStatusCREATING),
	string(api.VolumeStatusATTACHING),
	string(api.VolumeStatusDETACHING),
	string(api.VolumeStatusRECHARGING),
	string(api.VolumeStatusSCALING),
	string(api.VolumeStatusSNAPSHOTPROCESSING),
	string(api.VolumeStatusIMAGEPROCESSING),
}

var cdsStillInUsed = fmt.Errorf("cds is still in used")
