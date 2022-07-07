package connectivity

// Region represents BCC region
type Region string

// Constants of region definition
const (
	// Default Region
	DefaultRegion = RegionBeiJing

	// Regions
	RegionBeiJing   = Region("bj")
	RegionBaoDing   = Region("bd")
	RegionGuangZhou = Region("gz")
	RegionSuZhou    = Region("su")
	RegionShangHai  = Region("fsh")
	RegionWuHan     = Region("fwh")
	RegionHongKong  = Region("hkg")
	RegionSingapore = Region("sin")
)
