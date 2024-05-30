package mongodb

const (
	InstanceStatusRunning           = "RUNNING"
	InstanceStatusCreating          = "CREATING"
	InstanceStatusRestarting        = "RESTARTING"
	InstanceStatusClassChanging     = "CLASS_CHANGING"
	InstanceStatusNodeCreating      = "NODE_CREATING"
	InstanceStatusNodeRestarting    = "NODE_RESTARTING"
	InstanceStatusNodeClassChanging = "NODE_CLASS_CHANGING"
	InstanceStatusBackuping         = "BACKUPING"

	StorageTypeSSD         = "CDS_PREMIUM_SSD"
	StorageTypeEnhancedSSD = "CDS_ENHANCED_SSD"
	StorageTypeLocal       = "LOCAL_DISK"
)
