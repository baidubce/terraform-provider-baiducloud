package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/cert"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccCertResourceType = "baiducloud_cert"
	testAccCertResourceName = testAccCertResourceType + "." + BaiduCloudTestResourceName
)

func init() {
	resource.AddTestSweepers(testAccCertResourceType, &resource.Sweeper{
		Name: testAccCertResourceType,
		F:    testSweepCerts,
	})
}

func testSweepCerts(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}
	client := rawClient.(*connectivity.BaiduClient)

	raw, err := client.WithCertClient(func(client *cert.Client) (i interface{}, e error) {
		return client.ListCerts()
	})
	if err != nil {
		return fmt.Errorf("get Certs error: %s", err)
	}

	for _, c := range raw.(*cert.ListCertResult).Certs {
		if !strings.HasPrefix(c.CertName, BaiduCloudTestResourceAttrNamePrefix) {
			log.Printf("[INFO] Skipping Cert: %s (%s)", c.CertName, c.CertId)
			continue
		}

		log.Printf("[INFO] Deleting Cert: %s (%s)", c.CertName, c.CertId)

		_, err := client.WithCertClient(func(client *cert.Client) (i interface{}, e error) {
			return nil, client.DeleteCert(c.CertId)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Cert %s (%s)", c.CertName, c.CertId)
		}
	}

	return nil
}

func TestAccBaiduCloudCert(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCertDestory,

		Steps: []resource.TestStep{
			{
				Config: testAccCertConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCertResourceName),
					resource.TestCheckResourceAttr(testAccCertResourceName, "cert_name", BaiduCloudTestResourceAttrNamePrefix+"Cert"),
					resource.TestCheckResourceAttr(testAccCertResourceName, "cert_type", "1"),
				),
			},
			{
				Config: testAccCertConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccCertResourceName),
					resource.TestCheckResourceAttr(testAccCertResourceName, "cert_name", BaiduCloudTestResourceAttrNamePrefix+"CertUpdate"),
					resource.TestCheckResourceAttr(testAccCertResourceName, "cert_type", "1"),
				),
			},
		},
	})
}

func testAccCertDestory(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccCertResourceType {
			continue
		}

		_, err := client.WithCertClient(func(client *cert.Client) (i interface{}, e error) {
			return client.GetCertMeta(rs.Primary.ID)
		})
		if err != nil {
			if NotFoundError(err) || IsExceptedErrors(err, []string{"not exist"}) {
				continue
			}
			return WrapError(err)
		}
		return WrapError(Error("Cert still exist"))
	}

	return nil
}

func testAccCertConfig() string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  cert_name         = "%s"
  cert_server_data  = "-----BEGIN CERTIFICATE-----\nMIIEGzCCA8CgAwIBAgIQBHVIJNCDJKsC1maaUVgqdjAKBggqhkjOPQQDAjByMQswCQYDVQQGEwJDTjElMCMGA1UEChMcVHJ1c3RBc2lhIFRlY2hub2xvZ2llcywgSW5jLjEdMBsGA1UECxMURG9tYWluIFZhbGlkYXRlZCBTU0wxHTAbBgNVBAMTFFRydXN0QXNpYSBUTFMgRUNDIENBMB4XDTE5MDkwNjAwMDAwMFoXDTIwMDkwNTEyMDAwMFowHzEdMBsGA1UEAxMUdGVzdC55aW5jaGVuZ2ZlbmcuY24wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvLo4ICiTCCAoUwHwYDVR0jBBgwFoAUEoZEZiYIVCaPZTeyKU4mIeCTvtswHQYDVR0OBBYEFAichc0eFh+KdwMYjD7Pbvc8Q80IMB8GA1UdEQQYMBaCFHRlc3QueWluY2hlbmdmZW5nLmNuMA4GA1UdDwEB/wQEAwIHgDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwTAYDVR0gBEUwQzA3BglghkgBhv1sAQIwKjAoBggrBgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzAIBgZngQwBAgEwgZIGCCsGAQUFBwEBBIGFMIGCMDQGCCsGAQUFBzABhihodHRwOi8vc3RhdHVzZi5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tMEoGCCsGAQUFBzAChj5odHRwOi8vY2FjZXJ0cy5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tL1RydXN0QXNpYVRMU0VDQ0NBLmNydDAJBgNVHRMEAjAAMIIBAwYKKwYBBAHWeQIEAgSB9ASB8QDvAHUAu9nfvB+KcbWTlCOXqpJ7RzhXlQqrUugakJZkNo4e0YUAAAFtBK0O6QAABAMARjBEAiAdmHDa5NbRtLx3lc9nQ9G81RZycaqQPMj3+sazAo5vjQIgLNuFD7zperowYJAtetRR4QUi/8dORH087fWBp+Waj5MAdgCHdb/nWXz4jEOZX73zbv9WjUdWNv9KtWDBtOr/XqCDDwAAAW0ErQ9SAAAEAwBHMEUCIQDzdkB41ukE5XQGDTp8N4r+Aw/TZ/FlhPrrZryVGz9RIQIgWiuG2RHKCbh6FtJo62ml9RDYHeW/xA7c5sBBeKkSfG4wCgYIKoZIzj0EAwIDSQAwRgIhALnmf8VUwhxU0dRo2iOlfRb9uFy3hXMceU4IEvsLSwOVAiEAxsfjpOn0JyE943lhWRvjXX8FOm927cI5mbZ5F+p6dAA=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4zCCAsugAwIBAgIQBz/JpHsGAhj24Khq6fw+OzANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4NTdaFw0yNzEyMDgxMjI4NTdaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBFQ0MgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASdQvDzv44jBee0APcvKOWszZsRjc4j+L6DLlYOf9tSgvfOJplfMeDNDZzOQEcJbVPD+yekJQUmObCPOrgMhqMIo4IBTzCCAUswHQYDVR0OBBYEFBKGRGYmCFQmj2U3silOJiHgk77bMB8GA1UdIwQYMBaAFAPeUDVW0Uy7ZvCj4hsbw5eyPdFVMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBCBgNVHR8EOzA5MDegNaAzhjFodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRHbG9iYWxSb290Q0EuY3JsMEwGA1UdIARFMEMwNwYJYIZIAYb9bAECMCowKAYIKwYBBQUHAgEWHGh0dHBzOi8vd3d3LmRpZ2ljZXJ0LmNvbS9DUFMwCAYGZ4EMAQIBMA0GCSqGSIb3DQEBCwUAA4IBAQBZcGGhLE09CbQD5xP93NAuNC85G1BMa1OG2Q01TWvvgp7Qt1wNfRLAnhQT5pb7kRs+E7nM4IS894ufmuL452q8gYaq5HmvOmfhXMmL6K+eICfvyqjb/tSi8iy20ULO/TZhLhPor9tle52Yx811FG4i5vqwPIUEOEJ7pXe6RPVoBiwi4rbLspQGD/vYqrj9OJV4JctoIhhGq+y/sozU6nBXHfhVSD3x+hkOOst6tyRq481IyUWQHcFtwda3gfMnaA3dsag2dtJz33RIJIUfxXmVK7w4YzHOHifn7TYk8iNrDDLtql6vS8FjiUx3kJnI6zge1C9lUHhZ/aD3RiTJrwWI\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgp8yx31T7g0TyZcU4IdJS4px8p0b9FOHqx0uIMwtIjP6gCgYIKoZIzj0DAQehRANCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvL\n-----END PRIVATE KEY-----"
}
`, testAccCertResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"Cert")
}

func testAccCertConfigUpdate() string {
	return fmt.Sprintf(`
resource "%s" "%s" {
  cert_name         = "%s"
  cert_server_data  = "-----BEGIN CERTIFICATE-----\nMIIEHjCCA8SgAwIBAgIQD7e2kCM5IFr1AhZhtHco3DAKBggqhkjOPQQDAjByMQswCQYDVQQGEwJDTjElMCMGA1UEChMcVHJ1c3RBc2lhIFRlY2hub2xvZ2llcywgSW5jLjEdMBsGA1UECxMURG9tYWluIFZhbGlkYXRlZCBTU0wxHTAbBgNVBAMTFFRydXN0QXNpYSBUTFMgRUNDIENBMB4XDTE5MDkwNjAwMDAwMFoXDTIwMDkwNTEyMDAwMFowIDEeMBwGA1UEAxMVdGVzdDIueWluY2hlbmdmZW5nLmNuMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE5aIFysizmk3WriZXuYXzgcqcF7ORRPFIxQXvYTDGuuR9ybqBkT3zCt7n7YUW3z9AN4ux1Yxj2VnGM79YpPszGqOCAowwggKIMB8GA1UdIwQYMBaAFBKGRGYmCFQmj2U3silOJiHgk77bMB0GA1UdDgQWBBSoycYcJp+vvxdIWaM9QS4IchsYKDAgBgNVHREEGTAXghV0ZXN0Mi55aW5jaGVuZ2ZlbmcuY24wDgYDVR0PAQH/BAQDAgeAMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjBMBgNVHSAERTBDMDcGCWCGSAGG/WwBAjAqMCgGCCsGAQUFBwIBFhxodHRwczovL3d3dy5kaWdpY2VydC5jb20vQ1BTMAgGBmeBDAECATCBkgYIKwYBBQUHAQEEgYUwgYIwNAYIKwYBBQUHMAGGKGh0dHA6Ly9zdGF0dXNmLmRpZ2l0YWxjZXJ0dmFsaWRhdGlvbi5jb20wSgYIKwYBBQUHMAKGPmh0dHA6Ly9jYWNlcnRzLmRpZ2l0YWxjZXJ0dmFsaWRhdGlvbi5jb20vVHJ1c3RBc2lhVExTRUNDQ0EuY3J0MAkGA1UdEwQCMAAwggEFBgorBgEEAdZ5AgQCBIH2BIHzAPEAdgDuS723dc5guuFCaR+r4Z5mow9+X7By2IMAxHuJeqj9ywAAAW0FFJZQAAAEAwBHMEUCIDq3C14Mq4CaueNUWVIBKI3HGphyj4JqRKVvfGP4qBR4AiEAsgc3/WUucxBeK/+2vQJmFgE+kUwAa3ZGgoq4fmKsxlcAdwCHdb/nWXz4jEOZX73zbv9WjUdWNv9KtWDBtOr/XqCDDwAAAW0FFJa9AAAEAwBIMEYCIQDoRpKHe+ljJ6JmJoMzK3IE+f3AfLrN5f07D9eRIwqBNQIhAMYw+Sn8HZ53sxE5ttkJGetSu4mUf1bqrXG7CoSo5rjFMAoGCCqGSM49BAMCA0gAMEUCIQDjzWnH6V/OHVQvPZuaNXD6P/U4rdoUvhLnqoFkrRZxYAIgU7qPXUAOdwAWy0LuINOz0OmoXc5angeJAqK67hULNI4=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4zCCAsugAwIBAgIQBz/JpHsGAhj24Khq6fw+OzANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4NTdaFw0yNzEyMDgxMjI4NTdaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBFQ0MgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASdQvDzv44jBee0APcvKOWszZsRjc4j+L6DLlYOf9tSgvfOJplfMeDNDZzOQEcJbVPD+yekJQUmObCPOrgMhqMIo4IBTzCCAUswHQYDVR0OBBYEFBKGRGYmCFQmj2U3silOJiHgk77bMB8GA1UdIwQYMBaAFAPeUDVW0Uy7ZvCj4hsbw5eyPdFVMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBCBgNVHR8EOzA5MDegNaAzhjFodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRHbG9iYWxSb290Q0EuY3JsMEwGA1UdIARFMEMwNwYJYIZIAYb9bAECMCowKAYIKwYBBQUHAgEWHGh0dHBzOi8vd3d3LmRpZ2ljZXJ0LmNvbS9DUFMwCAYGZ4EMAQIBMA0GCSqGSIb3DQEBCwUAA4IBAQBZcGGhLE09CbQD5xP93NAuNC85G1BMa1OG2Q01TWvvgp7Qt1wNfRLAnhQT5pb7kRs+E7nM4IS894ufmuL452q8gYaq5HmvOmfhXMmL6K+eICfvyqjb/tSi8iy20ULO/TZhLhPor9tle52Yx811FG4i5vqwPIUEOEJ7pXe6RPVoBiwi4rbLspQGD/vYqrj9OJV4JctoIhhGq+y/sozU6nBXHfhVSD3x+hkOOst6tyRq481IyUWQHcFtwda3gfMnaA3dsag2dtJz33RIJIUfxXmVK7w4YzHOHifn7TYk8iNrDDLtql6vS8FjiUx3kJnI6zge1C9lUHhZ/aD3RiTJrwWI\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQg4vsAo5xhUZD92opgs+dSIDFHgFjikrZylNHvSSIyJjegCgYIKoZIzj0DAQehRANCAATlogXKyLOaTdauJle5hfOBypwXs5FE8UjFBe9hMMa65H3JuoGRPfMK3ufthRbfP0A3i7HVjGPZWcYzv1ik+zMa\n-----END PRIVATE KEY-----"
}
`, testAccCertResourceType, BaiduCloudTestResourceName, BaiduCloudTestResourceAttrNamePrefix+"CertUpdate")
}
