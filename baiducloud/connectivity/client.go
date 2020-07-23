package connectivity

import (
	"sync"

	"github.com/baidubce/bce-sdk-go/auth"
	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/cce"
	"github.com/baidubce/bce-sdk-go/services/cert"
	"github.com/baidubce/bce-sdk-go/services/cfc"
	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/baidubce/bce-sdk-go/services/rds"
	"github.com/baidubce/bce-sdk-go/services/scs"
	"github.com/baidubce/bce-sdk-go/services/sts"
	"github.com/baidubce/bce-sdk-go/services/sts/api"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/baidubce/bce-sdk-go/util/log"
)

// BaiduClient of BaiduCloud
type BaiduClient struct {
	config   *Config
	Region   Region
	Endpoint string

	Credentials *auth.BceCredentials

	bccConn    *bcc.Client
	vpcConn    *vpc.Client
	eipConn    *eip.Client
	appBlbConn *appblb.Client
	bosConn    *bos.Client
	certConn   *cert.Client
	cfcConn    *cfc.Client
	scsConn    *scs.Client
	cceConn    *cce.Client
	rdsConn    *rds.Client
}

type ApiVersion string

var goSdkMutex = sync.RWMutex{} // The Go SDK is not thread-safe

// Client for BaiduCloudClient
func (c *Config) Client() (*BaiduClient, error) {
	client := &BaiduClient{
		config: c,
		Region: c.Region,
	}

	if c.AssumeRoleAccountId != "" && c.AssumeRoleRoleName != "" {
		stsClient, err := sts.NewClient(c.AccessKey, c.SecretKey)
		if err != nil {
			return nil, err
		}

		args := &api.AssumeRoleArgs{
			AccountId: c.AssumeRoleAccountId,
			RoleName:  c.AssumeRoleRoleName,
			UserId:    c.AssumeRoleUserId,
			Acl:       c.AssumeRoleAcl,
		}
		assumeRole, err := stsClient.AssumeRole(args)
		if err != nil {
			return nil, err
		}

		stsCredential, err := auth.NewSessionBceCredentials(
			assumeRole.AccessKeyId,
			assumeRole.SecretAccessKey,
			assumeRole.SessionToken)
		if err != nil {
			return nil, err
		}

		client.Credentials = stsCredential
	} else {
		credentials, err := auth.NewBceCredentials(c.AccessKey, c.SecretKey)
		if err != nil {
			return nil, err
		}

		client.Credentials = credentials
	}

	return client, nil
}

func (client *BaiduClient) WithCommonClient(serviceCode ServiceCode) *BaiduClient {
	log.SetLogLevel(log.DEBUG)
	log.SetLogHandler(log.NONE)
	//log.SetLogDir(LogDir)

	region := client.config.Region
	if region == "" {
		region = DefaultRegion
	}
	endpoint, _ := client.config.ConfigEndpoints[serviceCode]
	if endpoint == "" {
		endpoint = loadEndpoint(region, serviceCode)
	}
	client.Endpoint = endpoint

	if client.Region == "" {
		client.Region = region
	}

	return client
}

func (client *BaiduClient) WithBccClient(do func(*bcc.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the BCC client if necessary
	if client.bccConn == nil {
		client.WithCommonClient(BCCCode)
		bccClient, err := bcc.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		bccClient.Config.Credentials = client.Credentials

		client.bccConn = bccClient
	}

	return do(client.bccConn)
}

func (client *BaiduClient) WithVpcClient(do func(*vpc.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the VPC client if necessary
	if client.vpcConn == nil {
		client.WithCommonClient(VPCCode)
		vpcClient, err := vpc.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		vpcClient.Config.Credentials = client.Credentials

		client.vpcConn = vpcClient
	}

	return do(client.vpcConn)
}

func (client *BaiduClient) WithEipClient(do func(*eip.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the EIP client if necessary
	if client.eipConn == nil {
		client.WithCommonClient(EIPCode)
		eipClient, err := eip.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		eipClient.Config.Credentials = client.Credentials

		client.eipConn = eipClient
	}

	return do(client.eipConn)
}

func (client *BaiduClient) WithAppBLBClient(do func(*appblb.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the APPBLB client if necessary
	if client.appBlbConn == nil {
		client.WithCommonClient(APPBLBCode)
		appBlbClient, err := appblb.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		appBlbClient.Config.Credentials = client.Credentials

		client.appBlbConn = appBlbClient
	}

	return do(client.appBlbConn)
}

func (client *BaiduClient) WithBosClient(do func(*bos.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the BOS client if necessary
	if client.bosConn == nil {
		client.WithCommonClient(BOSCode)
		bosClient, err := bos.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		bosClient.Config.Credentials = client.Credentials

		client.bosConn = bosClient
	}

	return do(client.bosConn)
}

func (client *BaiduClient) WithCertClient(do func(*cert.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CERT client if necessary
	if client.certConn == nil {
		client.WithCommonClient(CERTCode)
		certClient, err := cert.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		certClient.Config.Credentials = client.Credentials

		client.certConn = certClient
	}

	return do(client.certConn)
}

func (client *BaiduClient) WithCFCClient(do func(*cfc.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CFC client if necessary
	if client.cfcConn == nil {
		client.WithCommonClient(CFCCode)
		cfcClient, err := cfc.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		cfcClient.Config.Credentials = client.Credentials

		client.cfcConn = cfcClient
	}

	return do(client.cfcConn)
}

func (client *BaiduClient) WithScsClient(do func(*scs.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the SCS client if necessary
	if client.scsConn == nil {
		client.WithCommonClient(SCSCode)
		scsClient, err := scs.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		scsClient.Config.Credentials = client.Credentials

		client.scsConn = scsClient
	}

	return do(client.scsConn)
}

func (client *BaiduClient) WithCCEClient(do func(*cce.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CFC client if necessary
	if client.cceConn == nil {
		client.WithCommonClient(CCECode)
		cceClient, err := cce.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		cceClient.Config.Credentials = client.Credentials

		client.cceConn = cceClient
	}

	return do(client.cceConn)
}

func (client *BaiduClient) WithRdsClient(do func(*rds.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the RDS client if necessary
	if client.rdsConn == nil {
		client.WithCommonClient(RDSCode)
		rdsClient, err := rds.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		rdsClient.Config.Credentials = client.Credentials

		client.rdsConn = rdsClient
	}

	return do(client.rdsConn)

}
