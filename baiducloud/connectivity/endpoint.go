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
)

const (
	DefaultBJRegionBccEndPoint = "bcc.bj.baidubce.com"
	DefaultBJRegionEipEndPoint = "eip.bj.baidubce.com"
	DefaultBJRegionBlbEndPoint = "blb.bj.baidubce.com"
	DefaultBJRegionBosEndPoint = "bj.bcebos.com"

	DefaultGZRegionBccEndPoint = "bcc.gz.baidubce.com"
	DefaultGZRegionEipEndPoint = "eip.gz.baidubce.com"
	DefaultGZRegionBlbEndPoint = "blb.gz.baidubce.com"
	DefaultGZRegionBosEndPoint = "gz.bcebos.com"

	DefaultSURegionBccEndPoint = "bcc.su.baidubce.com"
	DefaultSURegionEipEndPoint = "eip.su.baidubce.com"
	DefaultSURegionBlbEndPoint = "blb.su.baidubce.com"
	DefaultSURegionBosEndPoint = "su.bcebos.com"

	DefaultCERTEndPoint = "certificate.baidubce.com"
)

var (
	// Default Region Endpoints
	DefaultRegionEndpoints = map[Region]map[ServiceCode]string{
		RegionHuaBei:  RegionBJEndpoints,
		RegionHuaNan:  RegionGZEndpoints,
		RegionHuaDong: RegionSUEndpoints,
	}

	// Region BJ Service Endpoints
	RegionBJEndpoints = map[ServiceCode]string{
		BCCCode:    DefaultBJRegionBccEndPoint,
		VPCCode:    DefaultBJRegionBccEndPoint,
		EIPCode:    DefaultBJRegionEipEndPoint,
		APPBLBCode: DefaultBJRegionBlbEndPoint,
		BOSCode:    DefaultBJRegionBosEndPoint,
		CERTCode:   DefaultCERTEndPoint,
	}

	// Region GZ Service Endpoints
	RegionGZEndpoints = map[ServiceCode]string{
		BCCCode:    DefaultGZRegionBccEndPoint,
		VPCCode:    DefaultGZRegionBccEndPoint,
		EIPCode:    DefaultGZRegionEipEndPoint,
		APPBLBCode: DefaultGZRegionBlbEndPoint,
		BOSCode:    DefaultGZRegionBosEndPoint,
		CERTCode:   DefaultCERTEndPoint,
	}

	// Region SU Service Endpoints
	RegionSUEndpoints = map[ServiceCode]string{
		BCCCode:    DefaultSURegionBccEndPoint,
		VPCCode:    DefaultSURegionBccEndPoint,
		EIPCode:    DefaultSURegionEipEndPoint,
		APPBLBCode: DefaultSURegionBlbEndPoint,
		BOSCode:    DefaultSURegionBosEndPoint,
		CERTCode:   DefaultCERTEndPoint,
	}
)

//xml
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
