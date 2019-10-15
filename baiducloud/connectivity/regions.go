package connectivity

// Region represents BCC region
type Region string

// Constants of region definition
const (
	// Default Region
	DefaultRegion = RegionHuaBei

	// Regions
	RegionHuaBei  = Region("bj")
	RegionHuaDong = Region("su")
	RegionHuaNan  = Region("gz")
)
