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
	BCCCode       = ServiceCode("BCC")
	VPCCode       = ServiceCode("VPC")
	EIPCode       = ServiceCode("EIP")
	APPBLBCode    = ServiceCode("APPBLB")
	BLBCode       = ServiceCode("BLB")
	BOSCode       = ServiceCode("BOS")
	CERTCode      = ServiceCode("CERT")
	CFCCode       = ServiceCode("CFC")
	CCECode       = ServiceCode("CCE")
	CCEv2Code     = ServiceCode("CCEv2")
	SCSCode       = ServiceCode("SCS")
	RDSCode       = ServiceCode("RDS")
	DTSCode       = ServiceCode("DTS")
	IAMCode       = ServiceCode("IAM")
	CDNCode       = ServiceCode("CDN")
	LOCALDNSCode  = ServiceCode("LOCALDNS")
	BBCCode       = ServiceCode("BBC")
	VPNCode       = ServiceCode("VPN")
	ENICode       = ServiceCode("ENI")
	CFSCode       = ServiceCode("CFS")
	SMSCode       = ServiceCode("SMS")
	BLSCode       = ServiceCode("BLS")
	BECCode       = ServiceCode("BEC")
	ETGATEWAYCode = ServiceCode("ETGATEWAY")
	DNSCode       = ServiceCode("DNS")
)

const (
	DefaultBJRegionBccEndPoint = "bcc.bj.baidubce.com"
	DefaultBJRegionBbcEndPoint = "bbc.bj.baidubce.com"
	DefaultBJRegionVpnEndPoint = "bcc.bj.baidubce.com"
	DefaultBJRegionEniEndPoint = "bcc.bj.baidubce.com"
	DefaultBJRegionEipEndPoint = "eip.bj.baidubce.com"
	DefaultBJRegionBlbEndPoint = "blb.bj.baidubce.com"
	DefaultBJRegionBosEndPoint = "bj.bcebos.com"
	DefaultBJRegionCfcEndPoint = "cfc.bj.baidubce.com"
	DefaultBJRegionCfsEndPoint = "cfs.bj.baidubce.com"
	DefaultBJRegionCceEndPoint = "cce.bj.baidubce.com"
	DefaultBJRegionScsEndPoint = "redis.bj.baidubce.com"
	DefaultBJRegionRdsEndPoint = "rds.bj.baidubce.com"
	DefaultBJRegionDtsEndPoint = "rds.bj.baidubce.com"
	DefaultBJRegionSmsEndPoint = "smsv3.bj.baidubce.com"
	DefaultBJRegionBlsEndPoint = "bls-log.bj.baidubce.com"
	DefaultBJRegionBecEndPoint = "bec.baidubce.com"

	DefaultBDRegionBccEndPoint = "bcc.bd.baidubce.com"
	DefaultBDRegionBbcEndPoint = "bbc.bd.baidubce.com"
	DefaultBDRegionVpnEndPoint = "bcc.bd.baidubce.com"
	DefaultBDRegionEniEndPoint = "bcc.bd.baidubce.com"
	DefaultBDRegionCfsEndPoint = "cfs.bd.baidubce.com"
	DefaultBDRegionEipEndPoint = "eip.bd.baidubce.com"

	DefaultGZRegionBccEndPoint = "bcc.gz.baidubce.com"
	DefaultGZRegionBbcEndPoint = "bbc.gz.baidubce.com"
	DefaultGZRegionVpnEndPoint = "bcc.gz.baidubce.com"
	DefaultGZRegionEipEndPoint = "eip.gz.baidubce.com"
	DefaultGZRegionEniEndPoint = "bcc.gz.baidubce.com"
	DefaultGZRegionBlbEndPoint = "blb.gz.baidubce.com"
	DefaultGZRegionBosEndPoint = "gz.bcebos.com"
	DefaultGZRegionCfcEndPoint = "cfc.gz.baidubce.com"
	DefaultGZRegionCceEndPoint = "cce.gz.baidubce.com"
	DefaultGZRegionScsEndPoint = "redis.gz.baidubce.com"
	DefaultGZRegionRdsEndPoint = "rds.gz.baidubce.com"
	DefaultGZRegionBlsEndPoint = "bls-log.gz.baidubce.com"

	DefaultSURegionBccEndPoint = "bcc.su.baidubce.com"
	DefaultSURegionBbcEndPoint = "bbc.su.baidubce.com"
	DefaultSURegionVPNEndPoint = "bcc.su.baidubce.com"
	DefaultSURegionEipEndPoint = "eip.su.baidubce.com"
	DefaultSURegionEniEndPoint = "bcc.su.baidubce.com"
	DefaultSURegionBlbEndPoint = "blb.su.baidubce.com"
	DefaultSURegionBosEndPoint = "su.bcebos.com"
	DefaultSURegionCfcEndPoint = "cfc.su.baidubce.com"
	DefaultSURegionCfsEndPoint = "cfs.su.baidubce.com"
	DefaultSURegionCceEndPoint = "cce.su.baidubce.com"
	DefaultSURegionScsEndPoint = "redis.su.baidubce.com"
	DefaultSURegionRdsEndPoint = "rds.su.baidubce.com"
	DefaultSURegionSmsEndPoint = "smsv3.su.baidubce.com"
	DefaultSURegionBlsEndPoint = "bls-log.su.baidubce.com"

	DefaultFSHRegionBccEndPoint = "bcc.fsh.baidubce.com"
	DefaultFSHRegionVPNEndPoint = "bcc.fsh.baidubce.com"
	DefaultFSHRegionEniEndPoint = "bcc.fsh.baidubce.com"

	DefaultFWHRegionBccEndPoint = "bcc.fwh.baidubce.com"
	DefaultFWHRegionBbcEndPoint = "bbc.fwh.baidubce.com"
	DefaultFWHRegionVPNEndPoint = "bbc.fwh.baidubce.com"
	DefaultFWHRegionEipEndPoint = "eip.fwh.baidubce.com"
	DefaultFWHRegionEniEndPoint = "bcc.fwh.baidubce.com"
	DefaultFWHRegionBlbEndPoint = "blb.fwh.baidubce.com"
	DefaultFWHRegionBosEndPoint = "fwh.bcebos.com"
	DefaultFWHRegionCfcEndPoint = "cfc.fwh.baidubce.com"
	DefaultFWHRegionCceEndPoint = "cce.fwh.baidubce.com"
	DefaultFWHRegionCfsEndPoint = "cfs.fwh.baidubce.com"
	DefaultFWHRegionScsEndPoint = "redis.fwh.baidubce.com"
	DefaultFWHRegionRdsEndPoint = "rds.fwh.baidubce.com"

	DefaultHKGRegionBccEndPoint = "bcc.hkg.baidubce.com"
	DefaultHKGRegionBbcEndPoint = "bbc.hkg.baidubce.com"
	DefaultHKGRegionVPNEndPoint = "bcc.hkg.baidubce.com"
	DefaultHKGRegionEniEndPoint = "bcc.hkg.baidubce.com"
	DefaultHKGRegionCfsEndPoint = "cfs.hkg.baidubce.com"
	DefaultHKGRegionEipEndPoint = "eip.hkg.baidubce.com"
	DefaultHKGRegionCceEndPoint = "cce.hkg.baidubce.com"
	DefaultHKGRegionBlbEndPoint = "blb.hkg.baidubce.com"
	DefaultHKGRegionBosEndPoint = "hkg.bcebos.com"

	DefaultSINRegionBccEndPoint = "bcc.sin.baidubce.com"
	DefaultSINRegionVPNEndPoint = "bcc.sin.baidubce.com"
	DefaultSINRegionEniEndPoint = "bcc.sin.baidubce.com"

	DefaultCERTEndPoint = "certificate.baidubce.com"
	DefaultIAMEndPoint  = "iam.bj.baidubce.com"
	DefaultCDNEndPoint  = "cdn.baidubce.com"
	DefaultDNSEndPoint  = "dns.baidubce.com"
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
		CERTCode: DefaultCERTEndPoint,
		IAMCode:  DefaultIAMEndPoint,
		CDNCode:  DefaultCDNEndPoint,
		DNSCode:  DefaultDNSEndPoint,

		BCCCode:       DefaultBJRegionBccEndPoint,
		BBCCode:       DefaultBJRegionBbcEndPoint,
		VPNCode:       DefaultBJRegionVpnEndPoint,
		VPCCode:       DefaultBJRegionBccEndPoint,
		EIPCode:       DefaultBJRegionEipEndPoint,
		ENICode:       DefaultBJRegionEniEndPoint,
		APPBLBCode:    DefaultBJRegionBlbEndPoint,
		BLBCode:       DefaultBJRegionBlbEndPoint,
		BOSCode:       DefaultBJRegionBosEndPoint,
		CFCCode:       DefaultBJRegionCfcEndPoint,
		CCECode:       DefaultBJRegionCceEndPoint,
		CCEv2Code:     DefaultBJRegionCceEndPoint,
		SCSCode:       DefaultBJRegionScsEndPoint,
		RDSCode:       DefaultBJRegionRdsEndPoint,
		DTSCode:       DefaultBJRegionDtsEndPoint,
		LOCALDNSCode:  DefaultBJRegionBccEndPoint,
		CFSCode:       DefaultBJRegionCfsEndPoint,
		SMSCode:       DefaultBJRegionSmsEndPoint,
		BLSCode:       DefaultBJRegionBlsEndPoint,
		BECCode:       DefaultBJRegionBecEndPoint,
		ETGATEWAYCode: DefaultBJRegionBccEndPoint,
	}

	// Region BD Service Endpoints
	RegionBDEndpoints = map[ServiceCode]string{
		CERTCode: DefaultCERTEndPoint,
		IAMCode:  DefaultIAMEndPoint,
		CDNCode:  DefaultCDNEndPoint,
		DNSCode:  DefaultDNSEndPoint,

		BCCCode:       DefaultBDRegionBccEndPoint,
		VPCCode:       DefaultBDRegionBccEndPoint,
		VPNCode:       DefaultBDRegionVpnEndPoint,
		ENICode:       DefaultBDRegionEniEndPoint,
		BBCCode:       DefaultBDRegionBbcEndPoint,
		CFSCode:       DefaultBDRegionCfsEndPoint,
		EIPCode:       DefaultBDRegionEipEndPoint,
		ETGATEWAYCode: DefaultBDRegionBccEndPoint,
	}

	// Region GZ Service Endpoints
	RegionGZEndpoints = map[ServiceCode]string{
		CERTCode: DefaultCERTEndPoint,
		IAMCode:  DefaultIAMEndPoint,
		CDNCode:  DefaultCDNEndPoint,
		DNSCode:  DefaultDNSEndPoint,

		BCCCode:       DefaultGZRegionBccEndPoint,
		VPNCode:       DefaultGZRegionVpnEndPoint,
		BBCCode:       DefaultGZRegionBbcEndPoint,
		VPCCode:       DefaultGZRegionBccEndPoint,
		EIPCode:       DefaultGZRegionEipEndPoint,
		ENICode:       DefaultGZRegionEniEndPoint,
		APPBLBCode:    DefaultGZRegionBlbEndPoint,
		BLBCode:       DefaultGZRegionBlbEndPoint,
		BOSCode:       DefaultGZRegionBosEndPoint,
		CFCCode:       DefaultGZRegionCfcEndPoint,
		CCECode:       DefaultGZRegionCceEndPoint,
		CCEv2Code:     DefaultGZRegionCceEndPoint,
		SCSCode:       DefaultGZRegionScsEndPoint,
		RDSCode:       DefaultGZRegionRdsEndPoint,
		BLSCode:       DefaultGZRegionBlsEndPoint,
		ETGATEWAYCode: DefaultGZRegionBccEndPoint,
	}

	// Region SU Service Endpoints
	RegionSUEndpoints = map[ServiceCode]string{
		CERTCode: DefaultCERTEndPoint,
		IAMCode:  DefaultIAMEndPoint,
		CDNCode:  DefaultCDNEndPoint,
		DNSCode:  DefaultDNSEndPoint,

		BCCCode:       DefaultSURegionBccEndPoint,
		BBCCode:       DefaultSURegionBbcEndPoint,
		VPNCode:       DefaultSURegionVPNEndPoint,
		VPCCode:       DefaultSURegionBccEndPoint,
		EIPCode:       DefaultSURegionEipEndPoint,
		ENICode:       DefaultSURegionEniEndPoint,
		APPBLBCode:    DefaultSURegionBlbEndPoint,
		BLBCode:       DefaultSURegionBlbEndPoint,
		BOSCode:       DefaultSURegionBosEndPoint,
		CFCCode:       DefaultSURegionCfcEndPoint,
		CCECode:       DefaultSURegionCceEndPoint,
		CCEv2Code:     DefaultSURegionCceEndPoint,
		SCSCode:       DefaultSURegionScsEndPoint,
		RDSCode:       DefaultSURegionRdsEndPoint,
		CFSCode:       DefaultSURegionCfsEndPoint,
		SMSCode:       DefaultSURegionSmsEndPoint,
		BLSCode:       DefaultSURegionBlsEndPoint,
		ETGATEWAYCode: DefaultSURegionBccEndPoint,
	}

	// Region FSH Service Endpoints
	RegionFSHEndpoints = map[ServiceCode]string{
		CERTCode: DefaultCERTEndPoint,
		IAMCode:  DefaultIAMEndPoint,
		CDNCode:  DefaultCDNEndPoint,
		DNSCode:  DefaultDNSEndPoint,

		BCCCode: DefaultFSHRegionBccEndPoint,
		VPNCode: DefaultFSHRegionVPNEndPoint,
		VPCCode: DefaultFSHRegionBccEndPoint,
		ENICode: DefaultFSHRegionEniEndPoint,
	}

	// Region FWH Service Endpoints
	RegionFWHEndpoints = map[ServiceCode]string{
		CERTCode: DefaultCERTEndPoint,
		IAMCode:  DefaultIAMEndPoint,
		CDNCode:  DefaultCDNEndPoint,
		DNSCode:  DefaultDNSEndPoint,

		BCCCode:       DefaultFWHRegionBccEndPoint,
		BBCCode:       DefaultFWHRegionBbcEndPoint,
		VPNCode:       DefaultFWHRegionVPNEndPoint,
		VPCCode:       DefaultFWHRegionBccEndPoint,
		EIPCode:       DefaultFWHRegionEipEndPoint,
		ENICode:       DefaultFWHRegionEniEndPoint,
		APPBLBCode:    DefaultFWHRegionBlbEndPoint,
		BLBCode:       DefaultFWHRegionBlbEndPoint,
		BOSCode:       DefaultFWHRegionBosEndPoint,
		CFCCode:       DefaultFWHRegionCfcEndPoint,
		CCECode:       DefaultFWHRegionCceEndPoint,
		CCEv2Code:     DefaultFWHRegionCceEndPoint,
		SCSCode:       DefaultFWHRegionScsEndPoint,
		RDSCode:       DefaultFWHRegionRdsEndPoint,
		CFSCode:       DefaultFWHRegionCfsEndPoint,
		ETGATEWAYCode: DefaultFWHRegionBccEndPoint,
	}

	// Region HKG Service Endpoints
	RegionHKGEndpoints = map[ServiceCode]string{
		CERTCode: DefaultCERTEndPoint,
		IAMCode:  DefaultIAMEndPoint,
		CDNCode:  DefaultCDNEndPoint,
		DNSCode:  DefaultDNSEndPoint,

		BCCCode:       DefaultHKGRegionBccEndPoint,
		BBCCode:       DefaultHKGRegionBbcEndPoint,
		VPCCode:       DefaultHKGRegionBccEndPoint,
		VPNCode:       DefaultHKGRegionVPNEndPoint,
		ENICode:       DefaultHKGRegionEniEndPoint,
		CFSCode:       DefaultHKGRegionCfsEndPoint,
		EIPCode:       DefaultHKGRegionEipEndPoint,
		CCECode:       DefaultHKGRegionCceEndPoint,
		CCEv2Code:     DefaultHKGRegionCceEndPoint,
		BLBCode:       DefaultHKGRegionBlbEndPoint,
		APPBLBCode:    DefaultHKGRegionBlbEndPoint,
		BOSCode:       DefaultHKGRegionBosEndPoint,
		ETGATEWAYCode: DefaultHKGRegionBccEndPoint,
	}

	// Region SIN Service Endpoints
	RegionSINEndpoints = map[ServiceCode]string{
		CERTCode: DefaultCERTEndPoint,
		IAMCode:  DefaultIAMEndPoint,
		CDNCode:  DefaultCDNEndPoint,
		DNSCode:  DefaultDNSEndPoint,

		BCCCode: DefaultSINRegionBccEndPoint,
		VPCCode: DefaultSINRegionBccEndPoint,
		VPNCode: DefaultSINRegionVPNEndPoint,
		ENICode: DefaultSINRegionEniEndPoint,
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
