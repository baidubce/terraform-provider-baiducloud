package connectivity

import (
	"fmt"
	"github.com/baidubce/bce-sdk-go/auth"
	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/cce"
	ccev2 "github.com/baidubce/bce-sdk-go/services/cce/v2"
	"github.com/baidubce/bce-sdk-go/services/cert"
	"github.com/baidubce/bce-sdk-go/services/cfc"
	"github.com/baidubce/bce-sdk-go/services/dts"
	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/baidubce/bce-sdk-go/services/iam"
	"github.com/baidubce/bce-sdk-go/services/localDns"
	"github.com/baidubce/bce-sdk-go/services/rds"
	"github.com/baidubce/bce-sdk-go/services/scs"
	"github.com/baidubce/bce-sdk-go/services/sts"
	"github.com/baidubce/bce-sdk-go/services/sts/api"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/baidubce/bce-sdk-go/util/log"
	"sync"
)

// BaiduClient of BaiduCloud
type BaiduClient struct {
	config   *Config
	Region   Region
	Endpoint string

	Credentials *auth.BceCredentials

	bccConn      *bcc.Client
	vpcConn      *vpc.Client
	eipConn      *eip.Client
	appBlbConn   *appblb.Client
	bosConn      *bos.Client
	certConn     *cert.Client
	cfcConn      *cfc.Client
	scsConn      *scs.Client
	cceConn      *cce.Client
	ccev2Conn    *ccev2.Client
	rdsConn      *rds.Client
	dtsConn      *dts.Client
	iamConn      *iam.Client
	localDnsConn *localDns.Client
}

type ApiVersion string

var goSdkMutex = sync.RWMutex{} // The Go SDK is not thread-safe

var providerVersion = "1.14.0"

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
	// Initialize the BCC client if necessary
	if client.bccConn == nil {
		client.WithCommonClient(BCCCode)
		bccClient, err := bcc.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			goSdkMutex.Unlock()
			return nil, err
		}
		bccClient.Config.Credentials = client.Credentials
		bccClient.Config.UserAgent = buildUserAgent()
		client.bccConn = bccClient
	}
	goSdkMutex.Unlock()
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
		vpcClient.Config.UserAgent = buildUserAgent()
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
		eipClient.Config.UserAgent = buildUserAgent()
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
		appBlbClient.Config.UserAgent = buildUserAgent()
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
		bosClient.Config.UserAgent = buildUserAgent()
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
		certClient.Config.UserAgent = buildUserAgent()
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
		cfcClient.Config.UserAgent = buildUserAgent()
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
		scsClient.Config.UserAgent = buildUserAgent()
		client.scsConn = scsClient
	}

	return do(client.scsConn)
}

func (client *BaiduClient) WithCCEClient(do func(*cce.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CCE client if necessary
	if client.cceConn == nil {
		client.WithCommonClient(CCECode)
		cceClient, err := cce.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		cceClient.Config.Credentials = client.Credentials
		cceClient.Config.UserAgent = buildUserAgent()
		client.cceConn = cceClient
	}

	return do(client.cceConn)
}

func (client *BaiduClient) WithCCEv2Client(do func(*ccev2.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the CCEv2 client if necessary
	if client.ccev2Conn == nil {
		client.WithCommonClient(CCEv2Code)
		ccev2Client, err := ccev2.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		ccev2Client.Config.Credentials = client.Credentials
		ccev2Client.Config.UserAgent = buildUserAgent()
		client.ccev2Conn = ccev2Client
	}

	return do(client.ccev2Conn)
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
		rdsClient.Config.UserAgent = buildUserAgent()
		client.rdsConn = rdsClient
	}

	return do(client.rdsConn)

}

func (client *BaiduClient) WithDtsClient(do func(*dts.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the DTS client if necessary
	if client.dtsConn == nil {
		client.WithCommonClient(DTSCode)
		dtsClient, err := dts.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		dtsClient.Config.UserAgent = buildUserAgent()
		client.dtsConn = dtsClient
	}

	return do(client.dtsConn)
}

func (client *BaiduClient) WithIamClient(do func(*iam.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the IAM client if necessary
	if client.iamConn == nil {
		client.WithCommonClient(IAMCode)
		iamClient, err := iam.NewClientWithEndpoint(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey,
			client.Endpoint)
		if err != nil {
			return nil, err
		}
		iamClient.Config.UserAgent = buildUserAgent()
		client.iamConn = iamClient
	}

	return do(client.iamConn)
}

func (client *BaiduClient) WithLocalDnsClient(do func(*localDns.Client) (interface{}, error)) (interface{}, error) {
	goSdkMutex.Lock()
	defer goSdkMutex.Unlock()

	// Initialize the LOCALDNS client if necessary
	if client.localDnsConn == nil {
		client.WithCommonClient(LOCALDNSCode)
		localDnsClient, err := localDns.NewClient(client.Credentials.AccessKeyId, client.Credentials.SecretAccessKey, client.Endpoint)
		if err != nil {
			return nil, err
		}
		localDnsClient.Config.Credentials = client.Credentials
		localDnsClient.Config.UserAgent = buildUserAgent()
		client.localDnsConn = localDnsClient
	}

	return do(client.localDnsConn)
}

func buildUserAgent() string {
	return fmt.Sprintf("terraform-provider-baiducloud/%s", providerVersion)
}
