package connectivity

// Config Constants
const (
	LogDir = "./logs/"
)

// Config Service Endpoints
type ConfigEndpoints map[ServiceCode]string

// Config of BaiduCloud
type Config struct {
	AccessKey string
	SecretKey string
	Region    Region

	// Config Service Endpoints Map
	ConfigEndpoints ConfigEndpoints
}
