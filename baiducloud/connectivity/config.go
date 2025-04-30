package connectivity

// Config Constants
const (
	LogDir = "./logs/"
)

// Config Service Endpoints
type ConfigEndpoints map[ServiceCode]string

// Config of BaiduCloud
type Config struct {
	AccessKey    string
	SecretKey    string
	SessionToken string
	Region       Region

	// assume role
	AssumeRoleRoleName  string
	AssumeRoleAccountId string
	AssumeRoleUserId    string
	AssumeRoleAcl       string

	// Config Service Endpoints Map
	ConfigEndpoints ConfigEndpoints
}
