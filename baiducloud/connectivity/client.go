package connectivity

import (
	"sync"

	"github.com/baidubce/bce-sdk-go/services/appblb"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/cert"
	"github.com/baidubce/bce-sdk-go/services/eip"
	"github.com/baidubce/bce-sdk-go/services/vpc"
	"github.com/baidubce/bce-sdk-go/util/log"
)

// BaiduClient of BaiduCloud
type BaiduClient struct {
	config    *Config
	AccessKey string
	SecretKey string
	Region    Region
	Endpoint  string

	bccConn    *bcc.Client
	vpcConn    *vpc.Client
	eipConn    *eip.Client
	appBlbConn *appblb.Client
	bosConn    *bos.Client
	certConn   *cert.Client
}

type ApiVersion string

var goSdkMutex = sync.RWMutex{} // The Go SDK is not thread-safe

// Client for BaiduCloudClient
func (c *Config) Client() (*BaiduClient, error) {
	return &BaiduClient{
		config:    c,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
		Region:    c.Region,
	}, nil
}

func (client *BaiduClient) WithCommonClient(serviceCode ServiceCode) *BaiduClient {
	// TODO: log set
	log.SetLogLevel(log.DEBUG)
	log.SetLogHandler(log.FILE)
	log.SetLogDir(LogDir)

	accessKey := client.config.AccessKey
	secretKey := client.config.SecretKey
	region := client.config.Region
	if region == "" {
		region = DefaultRegion
	}
	endpoint, _ := client.config.ConfigEndpoints[serviceCode]
	if endpoint == "" {
		endpoint = loadEndpoint(region, serviceCode)
	}
	client.Endpoint = endpoint

	if client.AccessKey == "" {
		client.AccessKey = accessKey
	}
	if client.SecretKey == "" {
		client.SecretKey = secretKey
	}
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
		bccClient, _ := bcc.NewClient(client.AccessKey, client.SecretKey, client.Endpoint)
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
		vpcClient, _ := vpc.NewClient(client.AccessKey, client.SecretKey, client.Endpoint)
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
		eipClient, _ := eip.NewClient(client.AccessKey, client.SecretKey, client.Endpoint)
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
		appBlbClient, _ := appblb.NewClient(client.AccessKey, client.SecretKey, client.Endpoint)
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
		bosClient, _ := bos.NewClient(client.AccessKey, client.SecretKey, client.Endpoint)
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
		certClient, _ := cert.NewClient(client.AccessKey, client.SecretKey, client.Endpoint)
		client.certConn = certClient
	}

	return do(client.certConn)
}
