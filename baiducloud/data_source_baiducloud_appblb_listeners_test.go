package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

const (
	testAccAppBLBListenersDataSourceName          = "data.baiducloud_appblb_listeners.default"
	testAccAppBLBListenersDataSourceAttrKeyPrefix = "listeners.0."
)

//lintignore:AT003
func TestAccBaiduCloudAppBLBListenersDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAppBLBListenersDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					// TCP Listener
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenersDataSourceName+"_TCP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_TCP", testAccAppBLBListenersDataSourceAttrKeyPrefix+"listener_port", "125"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_TCP", testAccAppBLBListenersDataSourceAttrKeyPrefix+"protocol", "TCP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_TCP", testAccAppBLBListenersDataSourceAttrKeyPrefix+"scheduler", "LeastConnection"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_TCP", testAccAppBLBListenersDataSourceAttrKeyPrefix+"tcp_session_timeout", "900"),

					// UDP Listener
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenersDataSourceName+"_UDP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_UDP", testAccAppBLBListenersDataSourceAttrKeyPrefix+"listener_port", "126"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_UDP", testAccAppBLBListenersDataSourceAttrKeyPrefix+"protocol", "UDP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_UDP", testAccAppBLBListenersDataSourceAttrKeyPrefix+"scheduler", "LeastConnection"),

					// HTTP Listener
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenersDataSourceName+"_HTTP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_HTTP", testAccAppBLBListenersDataSourceAttrKeyPrefix+"listener_port", "127"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_HTTP", testAccAppBLBListenersDataSourceAttrKeyPrefix+"protocol", "HTTP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_HTTP", testAccAppBLBListenersDataSourceAttrKeyPrefix+"scheduler", "LeastConnection"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_HTTP", testAccAppBLBListenersDataSourceAttrKeyPrefix+"keep_session", "true"),

					// HTTPS Listener
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenersDataSourceName+"_HTTP"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_HTTPS", testAccAppBLBListenersDataSourceAttrKeyPrefix+"listener_port", "128"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_HTTPS", testAccAppBLBListenersDataSourceAttrKeyPrefix+"protocol", "HTTPS"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_HTTPS", testAccAppBLBListenersDataSourceAttrKeyPrefix+"scheduler", "LeastConnection"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_HTTPS", testAccAppBLBListenersDataSourceAttrKeyPrefix+"keep_session", "true"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_HTTPS", testAccAppBLBListenersDataSourceAttrKeyPrefix+"cert_ids.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_HTTPS", testAccAppBLBListenersDataSourceAttrKeyPrefix+"encryption_type", "userDefind"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_HTTPS", testAccAppBLBListenersDataSourceAttrKeyPrefix+"encryption_protocols.#", "3"),

					// SSL Listener
					testAccCheckBaiduCloudDataSourceId(testAccAppBLBListenersDataSourceName+"_SSL"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_SSL", testAccAppBLBListenersDataSourceAttrKeyPrefix+"listener_port", "129"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_SSL", testAccAppBLBListenersDataSourceAttrKeyPrefix+"protocol", "SSL"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_SSL", testAccAppBLBListenersDataSourceAttrKeyPrefix+"scheduler", "LeastConnection"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_SSL", testAccAppBLBListenersDataSourceAttrKeyPrefix+"cert_ids.#", "1"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_SSL", testAccAppBLBListenersDataSourceAttrKeyPrefix+"encryption_type", "userDefind"),
					resource.TestCheckResourceAttr(testAccAppBLBListenersDataSourceName+"_SSL", testAccAppBLBListenersDataSourceAttrKeyPrefix+"encryption_protocols.#", "3"),
				),
			},
		},
	})
}

func testAccAppBLBListenersDataSourceConfig() string {
	return fmt.Sprintf(`
data "baiducloud_specs" "default" {}

data "baiducloud_zones" "default" {}

data "baiducloud_images" "default" {
  image_type = "System"
}

resource "baiducloud_instance" "default" {
  name                  = "%s"
  image_id              = data.baiducloud_images.default.images.0.id
  availability_zone     = data.baiducloud_zones.default.zones.0.zone_name
  cpu_count             = data.baiducloud_specs.default.specs.0.cpu_count
  memory_capacity_in_gb = data.baiducloud_specs.default.specs.0.memory_size_in_gb
  billing = {
    payment_timing = "Postpaid"
  }
}

resource "baiducloud_vpc" "default" {
  name        = "%s"
  description = "test"
  cidr        = "192.168.0.0/24"
}

resource "baiducloud_subnet" "default" {
  name        = "%s"
  zone_name   = data.baiducloud_zones.default.zones.0.zone_name
  cidr        = "192.168.0.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "test description"
}

resource "baiducloud_appblb" "default" {
  depends_on  = [baiducloud_instance.default]
  name        = "%s"
  description = ""
  vpc_id      = baiducloud_vpc.default.id
  subnet_id   = baiducloud_subnet.default.id
}

resource "baiducloud_cert" "default" {
  cert_name         = "%s"
  cert_server_data  = "-----BEGIN CERTIFICATE-----\nMIIEGzCCA8CgAwIBAgIQBHVIJNCDJKsC1maaUVgqdjAKBggqhkjOPQQDAjByMQswCQYDVQQGEwJDTjElMCMGA1UEChMcVHJ1c3RBc2lhIFRlY2hub2xvZ2llcywgSW5jLjEdMBsGA1UECxMURG9tYWluIFZhbGlkYXRlZCBTU0wxHTAbBgNVBAMTFFRydXN0QXNpYSBUTFMgRUNDIENBMB4XDTE5MDkwNjAwMDAwMFoXDTIwMDkwNTEyMDAwMFowHzEdMBsGA1UEAxMUdGVzdC55aW5jaGVuZ2ZlbmcuY24wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvLo4ICiTCCAoUwHwYDVR0jBBgwFoAUEoZEZiYIVCaPZTeyKU4mIeCTvtswHQYDVR0OBBYEFAichc0eFh+KdwMYjD7Pbvc8Q80IMB8GA1UdEQQYMBaCFHRlc3QueWluY2hlbmdmZW5nLmNuMA4GA1UdDwEB/wQEAwIHgDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwTAYDVR0gBEUwQzA3BglghkgBhv1sAQIwKjAoBggrBgEFBQcCARYcaHR0cHM6Ly93d3cuZGlnaWNlcnQuY29tL0NQUzAIBgZngQwBAgEwgZIGCCsGAQUFBwEBBIGFMIGCMDQGCCsGAQUFBzABhihodHRwOi8vc3RhdHVzZi5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tMEoGCCsGAQUFBzAChj5odHRwOi8vY2FjZXJ0cy5kaWdpdGFsY2VydHZhbGlkYXRpb24uY29tL1RydXN0QXNpYVRMU0VDQ0NBLmNydDAJBgNVHRMEAjAAMIIBAwYKKwYBBAHWeQIEAgSB9ASB8QDvAHUAu9nfvB+KcbWTlCOXqpJ7RzhXlQqrUugakJZkNo4e0YUAAAFtBK0O6QAABAMARjBEAiAdmHDa5NbRtLx3lc9nQ9G81RZycaqQPMj3+sazAo5vjQIgLNuFD7zperowYJAtetRR4QUi/8dORH087fWBp+Waj5MAdgCHdb/nWXz4jEOZX73zbv9WjUdWNv9KtWDBtOr/XqCDDwAAAW0ErQ9SAAAEAwBHMEUCIQDzdkB41ukE5XQGDTp8N4r+Aw/TZ/FlhPrrZryVGz9RIQIgWiuG2RHKCbh6FtJo62ml9RDYHeW/xA7c5sBBeKkSfG4wCgYIKoZIzj0EAwIDSQAwRgIhALnmf8VUwhxU0dRo2iOlfRb9uFy3hXMceU4IEvsLSwOVAiEAxsfjpOn0JyE943lhWRvjXX8FOm927cI5mbZ5F+p6dAA=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIID4zCCAsugAwIBAgIQBz/JpHsGAhj24Khq6fw+OzANBgkqhkiG9w0BAQsFADBhMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBDQTAeFw0xNzEyMDgxMjI4NTdaFw0yNzEyMDgxMjI4NTdaMHIxCzAJBgNVBAYTAkNOMSUwIwYDVQQKExxUcnVzdEFzaWEgVGVjaG5vbG9naWVzLCBJbmMuMR0wGwYDVQQLExREb21haW4gVmFsaWRhdGVkIFNTTDEdMBsGA1UEAxMUVHJ1c3RBc2lhIFRMUyBFQ0MgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASdQvDzv44jBee0APcvKOWszZsRjc4j+L6DLlYOf9tSgvfOJplfMeDNDZzOQEcJbVPD+yekJQUmObCPOrgMhqMIo4IBTzCCAUswHQYDVR0OBBYEFBKGRGYmCFQmj2U3silOJiHgk77bMB8GA1UdIwQYMBaAFAPeUDVW0Uy7ZvCj4hsbw5eyPdFVMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYIKwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBCBgNVHR8EOzA5MDegNaAzhjFodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vRGlnaUNlcnRHbG9iYWxSb290Q0EuY3JsMEwGA1UdIARFMEMwNwYJYIZIAYb9bAECMCowKAYIKwYBBQUHAgEWHGh0dHBzOi8vd3d3LmRpZ2ljZXJ0LmNvbS9DUFMwCAYGZ4EMAQIBMA0GCSqGSIb3DQEBCwUAA4IBAQBZcGGhLE09CbQD5xP93NAuNC85G1BMa1OG2Q01TWvvgp7Qt1wNfRLAnhQT5pb7kRs+E7nM4IS894ufmuL452q8gYaq5HmvOmfhXMmL6K+eICfvyqjb/tSi8iy20ULO/TZhLhPor9tle52Yx811FG4i5vqwPIUEOEJ7pXe6RPVoBiwi4rbLspQGD/vYqrj9OJV4JctoIhhGq+y/sozU6nBXHfhVSD3x+hkOOst6tyRq481IyUWQHcFtwda3gfMnaA3dsag2dtJz33RIJIUfxXmVK7w4YzHOHifn7TYk8iNrDDLtql6vS8FjiUx3kJnI6zge1C9lUHhZ/aD3RiTJrwWI\n-----END CERTIFICATE-----"
  cert_private_data = "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgp8yx31T7g0TyZcU4IdJS4px8p0b9FOHqx0uIMwtIjP6gCgYIKoZIzj0DAQehRANCAAR+aGvOdizh+oAWwT6829WdcZw7oBJVU1UvKQdm7dW/7SIdrMEWq6NIWaERMKkLD6gQ6Y5KFV9oDQdSocGBtBvL\n-----END PRIVATE KEY-----"
}

resource "baiducloud_appblb_listener" "default_TCP" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 125
  protocol      = "TCP"
  scheduler     = "LeastConnection"
}

resource "baiducloud_appblb_listener" "default_UDP" {
  blb_id         = baiducloud_appblb.default.id
  listener_port  = 126
  protocol       = "UDP"
  scheduler      = "LeastConnection"
}

resource "baiducloud_appblb_listener" "default_HTTP" {
  blb_id        = baiducloud_appblb.default.id
  listener_port = 127
  protocol      = "HTTP"
  scheduler     = "LeastConnection"
  keep_session  = true
}

resource "baiducloud_appblb_listener" "default_HTTPS" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 128
  protocol             = "HTTPS"
  scheduler            = "LeastConnection"
  keep_session         = true
  cert_ids             = [baiducloud_cert.default.id]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}

resource "baiducloud_appblb_listener" "default_SSL" {
  blb_id               = baiducloud_appblb.default.id
  listener_port        = 129
  protocol             = "SSL"
  scheduler            = "LeastConnection"
  cert_ids             = [baiducloud_cert.default.id]
  encryption_protocols = ["sslv3", "tlsv10", "tlsv11"]
  encryption_type      = "userDefind"
}

data "baiducloud_appblb_listeners" "default_TCP" {
  blb_id        = baiducloud_appblb.default.id
  protocol      = baiducloud_appblb_listener.default_TCP.protocol
  listener_port = baiducloud_appblb_listener.default_TCP.listener_port
}

data "baiducloud_appblb_listeners" "default_UDP" {
  blb_id        = baiducloud_appblb.default.id
  protocol      = baiducloud_appblb_listener.default_UDP.protocol
  listener_port = baiducloud_appblb_listener.default_UDP.listener_port
}

data "baiducloud_appblb_listeners" "default_HTTP" {
  blb_id        = baiducloud_appblb.default.id
  protocol      = baiducloud_appblb_listener.default_HTTP.protocol
  listener_port = baiducloud_appblb_listener.default_HTTP.listener_port
}

data "baiducloud_appblb_listeners" "default_HTTPS" {
  blb_id        = baiducloud_appblb.default.id
  protocol      = baiducloud_appblb_listener.default_HTTPS.protocol
  listener_port = baiducloud_appblb_listener.default_HTTPS.listener_port
}

data "baiducloud_appblb_listeners" "default_SSL" {
  blb_id        = baiducloud_appblb.default.id
  protocol      = baiducloud_appblb_listener.default_SSL.protocol
  listener_port = baiducloud_appblb_listener.default_SSL.listener_port
}
`, BaiduCloudTestResourceAttrNamePrefix+"BCC",
		BaiduCloudTestResourceAttrNamePrefix+"VPC",
		BaiduCloudTestResourceAttrNamePrefix+"Subnet",
		BaiduCloudTestResourceAttrNamePrefix+"APPBLB",
		BaiduCloudTestResourceAttrNamePrefix+"Cert")
}
