package connectivity

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Load endpoints from endpoints.xml or environment variables to meet specified application scenario, like private cloud.
type ServiceCode string

const (
	BCCCode    = ServiceCode("BCC")
	VPCCode    = ServiceCode("VPC")
	EIPCode    = ServiceCode("EIP")
	APPBLBCode = ServiceCode("APPBLB")
	BOSCode    = ServiceCode("BOS")
	CERTCode   = ServiceCode("CERT")
	CFCCode    = ServiceCode("CFC")
	CCECode    = ServiceCode("CCE")
	CCEv2Code  = ServiceCode("CCEv2")
	SCSCode    = ServiceCode("SCS")
	RDSCode    = ServiceCode("RDS")
	DTSCode    = ServiceCode("DTS")
	IAMCode    = ServiceCode("IAM")
)

const (
	DefaultBJRegionBccEndPoint = "bcc.bj.baidubce.com"
	DefaultBJRegionEipEndPoint = "eip.bj.baidubce.com"
	DefaultBJRegionBlbEndPoint = "blb.bj.baidubce.com"
	DefaultBJRegionBosEndPoint = "bj.bcebos.com"
	DefaultBJRegionCfcEndPoint = "cfc.bj.baidubce.com"
	DefaultBJRegionCceEndPoint = "cce.bj.baidubce.com"
	DefaultBJRegionScsEndPoint = "redis.bj.baidubce.com"
	DefaultBJRegionRdsEndPoint = "rds.bj.baidubce.com"
	DefaultBJRegionDtsEndPoint = "rds.bj.baidubce.com"

	DefaultBDRegionBccEndPoint = "bcc.bd.baidubce.com"

	DefaultGZRegionBccEndPoint = "bcc.gz.baidubce.com"
	DefaultGZRegionEipEndPoint = "eip.gz.baidubce.com"
	DefaultGZRegionBlbEndPoint = "blb.gz.baidubce.com"
	DefaultGZRegionBosEndPoint = "gz.bcebos.com"
	DefaultGZRegionCfcEndPoint = "cfc.gz.baidubce.com"
	DefaultGZRegionCceEndPoint = "cce.gz.baidubce.com"
	DefaultGZRegionScsEndPoint = "redis.gz.baidubce.com"
	DefaultGZRegionRdsEndPoint = "rds.gz.baidubce.com"

	DefaultSURegionBccEndPoint = "bcc.su.baidubce.com"
	DefaultSURegionEipEndPoint = "eip.su.baidubce.com"
	DefaultSURegionBlbEndPoint = "blb.su.baidubce.com"
	DefaultSURegionBosEndPoint = "su.bcebos.com"
	DefaultSURegionCfcEndPoint = "cfc.su.baidubce.com"
	DefaultSURegionCceEndPoint = "cce.su.baidubce.com"
	DefaultSURegionScsEndPoint = "redis.su.baidubce.com"
	DefaultSURegionRdsEndPoint = "rds.su.baidubce.com"

	DefaultFSHRegionBccEndPoint = "bcc.fsh.baidubce.com"

	DefaultFWHRegionBccEndPoint = "bcc.fwh.baidubce.com"
	DefaultFWHRegionEipEndPoint = "eip.fwh.baidubce.com"
	DefaultFWHRegionBlbEndPoint = "blb.fwh.baidubce.com"
	DefaultFWHRegionBosEndPoint = "fwh.bcebos.com"
	DefaultFWHRegionCfcEndPoint = "cfc.fwh.baidubce.com"
	DefaultFWHRegionCceEndPoint = "cce.fwh.baidubce.com"
	DefaultFWHRegionScsEndPoint = "redis.fwh.baidubce.com"
	DefaultFWHRegionRdsEndPoint = "rds.fwh.baidubce.com"

	DefaultHKGRegionBccEndPoint = "bcc.hkg.baidubce.com"

	DefaultSINRegionBccEndPoint = "bcc.sin.baidubce.com"

	DefaultCERTEndPoint = "certificate.baidubce.com"
	DefaultIAMEndPoint  = "iam.bj.baidubce.com"
)

var (
	// Default Region Endpoints
	DefaultRegionEndpoints = map[Region]map[ServiceCode]string{
		RegionBeiJing:   RegionBJEndpoints,
		RegionBaoDing:   RegionBDEndpoints,
		RegionGuangZhou: RegionGZEndpoints,
		RegionSuZhou:    RegionSUEndpoints,
		RegionShangHai:  RegionFSHEndpoints,
		RegionWuHan:     RegionFWHEndpoints,
		RegionHongKong:  RegionHKGEndpoints,
		RegionSingapore: RegionSINEndpoints,
	}

	// Region BJ Service Endpoints
	RegionBJEndpoints = map[ServiceCode]string{
		BCCCode:    DefaultBJRegionBccEndPoint,
		VPCCode:    DefaultBJRegionBccEndPoint,
		EIPCode:    DefaultBJRegionEipEndPoint,
		APPBLBCode: DefaultBJRegionBlbEndPoint,
		BOSCode:    DefaultBJRegionBosEndPoint,
		CERTCode:   DefaultCERTEndPoint,
		CFCCode:    DefaultBJRegionCfcEndPoint,
		CCECode:    DefaultBJRegionCceEndPoint,
		CCEv2Code:  DefaultBJRegionCceEndPoint,
		SCSCode:    DefaultBJRegionScsEndPoint,
		RDSCode:    DefaultBJRegionRdsEndPoint,
		DTSCode:    DefaultBJRegionDtsEndPoint,
		IAMCode:    DefaultIAMEndPoint,
	}

	// Region BD Service Endpoints
	RegionBDEndpoints = map[ServiceCode]string{
		BCCCode: DefaultBDRegionBccEndPoint,
		VPCCode: DefaultBDRegionBccEndPoint,
	}

	// Region GZ Service Endpoints
	RegionGZEndpoints = map[ServiceCode]string{
		BCCCode:    DefaultGZRegionBccEndPoint,
		VPCCode:    DefaultGZRegionBccEndPoint,
		EIPCode:    DefaultGZRegionEipEndPoint,
		APPBLBCode: DefaultGZRegionBlbEndPoint,
		BOSCode:    DefaultGZRegionBosEndPoint,
		CERTCode:   DefaultCERTEndPoint,
		CFCCode:    DefaultGZRegionCfcEndPoint,
		CCECode:    DefaultGZRegionCceEndPoint,
		CCEv2Code:  DefaultGZRegionCceEndPoint,
		SCSCode:    DefaultGZRegionScsEndPoint,
		RDSCode:    DefaultGZRegionRdsEndPoint,
		IAMCode:    DefaultIAMEndPoint,
	}

	// Region SU Service Endpoints
	RegionSUEndpoints = map[ServiceCode]string{
		BCCCode:    DefaultSURegionBccEndPoint,
		VPCCode:    DefaultSURegionBccEndPoint,
		EIPCode:    DefaultSURegionEipEndPoint,
		APPBLBCode: DefaultSURegionBlbEndPoint,
		BOSCode:    DefaultSURegionBosEndPoint,
		CERTCode:   DefaultCERTEndPoint,
		CFCCode:    DefaultSURegionCfcEndPoint,
		CCECode:    DefaultSURegionCceEndPoint,
		CCEv2Code:  DefaultSURegionCceEndPoint,
		SCSCode:    DefaultSURegionScsEndPoint,
		RDSCode:    DefaultSURegionRdsEndPoint,
		IAMCode:    DefaultIAMEndPoint,
	}

	// Region FSH Service Endpoints
	RegionFSHEndpoints = map[ServiceCode]string{
		BCCCode: DefaultFSHRegionBccEndPoint,
		VPCCode: DefaultFSHRegionBccEndPoint,
	}

	// Region FWH Service Endpoints
	RegionFWHEndpoints = map[ServiceCode]string{
		BCCCode:    DefaultFWHRegionBccEndPoint,
		VPCCode:    DefaultFWHRegionBccEndPoint,
		EIPCode:    DefaultFWHRegionEipEndPoint,
		APPBLBCode: DefaultFWHRegionBlbEndPoint,
		BOSCode:    DefaultFWHRegionBosEndPoint,
		CERTCode:   DefaultCERTEndPoint,
		CFCCode:    DefaultFWHRegionCfcEndPoint,
		CCECode:    DefaultFWHRegionCceEndPoint,
		CCEv2Code:  DefaultFWHRegionCceEndPoint,
		SCSCode:    DefaultFWHRegionScsEndPoint,
		RDSCode:    DefaultFWHRegionRdsEndPoint,
		IAMCode:    DefaultIAMEndPoint,
	}

	// Region HKG Service Endpoints
	RegionHKGEndpoints = map[ServiceCode]string{
		BCCCode: DefaultHKGRegionBccEndPoint,
		VPCCode: DefaultHKGRegionBccEndPoint,
	}

	// Region SIN Service Endpoints
	RegionSINEndpoints = map[ServiceCode]string{
		BCCCode: DefaultSINRegionBccEndPoint,
		VPCCode: DefaultSINRegionBccEndPoint,
	}
)

// Endpoints xml struct
type Endpoints struct {
	Endpoint []Endpoint `xml:"Endpoint"`
}

type Endpoint struct {
	Name     string   `xml:"name,attr"`
	Regions  Regions  `xml:"Region"`
	Products Products `xml:"Products"`
}

type Regions struct {
	Region string `xml:"Region"`
}

type Products struct {
	Product []Product `xml:"Product"`
}

type Product struct {
	ProductName string `xml:"ProductName"`
	DomainName  string `xml:"DomainName"`
}

func loadEndpointFromEnvOrXML(region Region, serviceCode ServiceCode) string {
	endpoint := strings.TrimSpace(os.Getenv(fmt.Sprintf("%s_ENDPOINT", string(serviceCode))))
	if endpoint != "" {
		return endpoint
	}

	// Load current path endpoint file endpoints.xml, if failed, it will load from environment variables TF_ENDPOINT_PATH
	data, err := ioutil.ReadFile("./endpoints.xml")
	if err != nil || len(data) <= 0 {
		d, e := ioutil.ReadFile(os.Getenv("TF_ENDPOINT_PATH"))
		if e != nil {
			return ""
		}
		data = d
	}
	var endpoints Endpoints
	err = xml.Unmarshal(data, &endpoints)
	if err != nil {
		return ""
	}
	for _, endpoint := range endpoints.Endpoint {
		if endpoint.Regions.Region == string(region) {
			for _, product := range endpoint.Products.Product {
				if strings.ToLower(product.ProductName) == strings.ToLower(string(serviceCode)) {
					return strings.TrimSpace(product.DomainName)
				}
			}
		}
	}

	return ""
}

func loadEndpoint(region Region, serviceCode ServiceCode) string {
	endpoint := loadEndpointFromEnvOrXML(region, serviceCode)
	if endpoint == "" {
		endpoint = DefaultRegionEndpoints[region][serviceCode]
	}

	return endpoint
}
