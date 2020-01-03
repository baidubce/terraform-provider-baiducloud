package connectivity

// Region represents BCC region
type Region string

// Constants of region definition
const (
	// Default Region
	DefaultRegion = RegionBeiJing

	// Regions
	RegionBeiJing   = Region("bj")
	RegionSuZhou    = Region("su")
	RegionGuangZhou = Region("gz")
	RegionWuHan     = Region("fwh")
)
