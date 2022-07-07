package baiducloud

const (
	SCSStatusPrecreate      = "Precreat"
	SCSStatusCreating       = "Creating"
	SCSStatusRunning        = "Running"
	SCSStatusRebooting      = "Rebooting"
	SCSStatusPausing        = "Pausing"
	SCSStatusPaused         = "Paused"
	SCSStatusDeleted        = "Deleted"
	SCSStatusDeleting       = "Deleting"
	SCSStatusFailed         = "Failed"
	SCSStatusModifying      = "Modifying"
	SCSStatusModifyFailed   = "Modifyfailed"
	SCSStatusBackuping      = "Backuping"
	SCSStatusAzTransforming = "Aztransforming"
	SCSStatusExpire         = "Expire"
	SCSStatusFlushing       = "Flushing"
	SCSStatusFlushFailed    = "Flush failed"
	SCSStatusIsolated       = "isolated"
)

func SCSEngineIntegers() map[string]int {
	return map[string]int{
		"memcache": 1,
		"redis":    2,
		"PegaDB":   3,
	}
}
