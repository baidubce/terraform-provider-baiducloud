package baiducloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testAccBosBucketsDataSourceName          = "data.baiducloud_bos_buckets.default"
	testAccBosBucketsDataSourceAttrKeyPrefix = "buckets.0."
)

//lintignore:AT003
func TestAccBaiduCloudBosBucketsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBosBucketsDataSourceConfig(BaiduCloudTestResourceTypeNameBosBucket),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBosBucketsDataSourceName),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"bucket", BaiduCloudTestResourceTypeNameBosBucket+"-bucket-new"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"acl", "public-read-write"),
					resource.TestCheckResourceAttrSet(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"location"),
					resource.TestCheckResourceAttrSet(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"creation_date"),
					resource.TestCheckResourceAttrSet(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"owner_id"),
					resource.TestCheckResourceAttrSet(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"owner_name"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"replication_configuration.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"replication_configuration.0.id", "test-rc"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"replication_configuration.0.status", "enabled"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"replication_configuration.0.resource.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"replication_configuration.0.destination.0.bucket", BaiduCloudTestResourceTypeNameBosBucket+"-bucket-peer"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"replication_configuration.0.replicate_deletes", "disabled"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"logging.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"logging.0.target_bucket", BaiduCloudTestResourceTypeNameBosBucket+"-bucket-new"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"logging.0.target_prefix", "logs/"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"lifecycle_rule.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"lifecycle_rule.0.id", "test-lr"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"lifecycle_rule.0.status", "enabled"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"lifecycle_rule.0.resource.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"lifecycle_rule.0.condition.0.time.0.date_greater_than", "2029-12-31T00:00:00Z"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"lifecycle_rule.0.action.0.name", "DeleteObject"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"storage_class", "COLD"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"server_side_encryption_rule", "AES256"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"website.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"website.0.index_document", "index.html"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"website.0.error_document", "err.html"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"cors_rule.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"cors_rule.0.allowed_origins.0", "https://www.baidu.com"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"cors_rule.0.allowed_methods.0", "GET"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"cors_rule.0.max_age_seconds", "1800"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"copyright_protection.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketsDataSourceName, testAccBosBucketsDataSourceAttrKeyPrefix+"copyright_protection.0.resource.#", "1"),
				),
			},
		},
	})
}

func testAccBosBucketsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "baiducloud_bos_bucket" "peer" {
  bucket = "%s"
  acl    = "public-read-write"
}

resource "baiducloud_bos_bucket" "default" {
  bucket = "%s"
  acl    = "public-read-write"

  logging {
    target_bucket = "%s"
    target_prefix = "logs/"
  }

  replication_configuration {
    id       = "test-rc"
    status   = "enabled"
    resource = ["%s"]
    destination {
      bucket = baiducloud_bos_bucket.peer.bucket
    }
    replicate_deletes = "disabled"
  }

  force_destroy = true

  lifecycle_rule {
    id       = "test-lr"
    status   =  "enabled"
    resource = ["%s"]
    condition {
      time {
        date_greater_than = "2029-12-31T00:00:00Z"
      }
    }
    action {
      name = "DeleteObject"
    }
  }

  storage_class = "COLD"

  server_side_encryption_rule = "AES256"

  website{
    index_document = "index.html"
    error_document = "err.html"
  }

  cors_rule {
    allowed_origins = ["https://www.baidu.com"]
    allowed_methods = ["GET"]
    max_age_seconds = 1800
  }

  copyright_protection {
    resource = ["%s"]
  }
}

data "baiducloud_bos_buckets" "default" {
  bucket = baiducloud_bos_bucket.default.bucket

  filter {
    name = "acl"
    values = ["public-read-write"]
  }
}
`, name+"-bucket-peer", name+"-bucket-new", name+"-bucket-new",
		name+"-bucket-new"+"/*", name+"-bucket-new"+"/*", name+"-bucket-new"+"/*")
}
