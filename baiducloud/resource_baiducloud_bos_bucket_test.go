package baiducloud

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	testAccBosBucketResourceType               = "baiducloud_bos_bucket"
	BaiduCloudTestBucketResourceAttrNamePrefix = "test-baiduacc-"
	testAccBosBucketResourceName               = testAccBosBucketResourceType + "." + BaiduCloudTestResourceName
	testAccBosBucketResourceAttrName           = BaiduCloudTestBucketResourceAttrNamePrefix + "bucket"
)

func init() {
	resource.AddTestSweepers(testAccBosBucketResourceType, &resource.Sweeper{
		Name: testAccBosBucketResourceType,
		F:    testSweepBosBuckets,
	})
}

func testSweepBosBuckets(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("get BaiduCloud client error: %s", err)
	}

	client := rawClient.(*connectivity.BaiduClient)

	raw, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.ListBuckets()
	})
	if err != nil {
		return fmt.Errorf("list buckets error: %v", err)
	}

	result, _ := raw.(*api.ListBucketsResult)
	if err != nil {
		return fmt.Errorf("get buckets error: %v", err)
	}

	for _, buc := range result.Buckets {
		if !strings.HasPrefix(buc.Name, BaiduCloudTestResourceAttrNamePrefix) {
			log.Printf("[INFO] Skipping bucket: %s", buc.Name)
			continue
		}

		log.Printf("[INFO] Deleting bucket: %s", buc.Name)
		_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.DeleteBucket(buc.Name)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete bucket %s", buc.Name)
		}
	}

	return nil
}

func TestAccBaiduCloudBosBucket(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccBosBucketDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccBosBucketConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBosBucketResourceName),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "bucket", testAccBosBucketResourceAttrName),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "acl", "public-read-write"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.0.id", "test-rc"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.0.status", "enabled"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.0.resource.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.0.destination.0.bucket", BaiduCloudTestBucketResourceAttrNamePrefix+"bucket-peer"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.0.replicate_deletes", "disabled"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "force_destroy", "false"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "logging.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "logging.0.target_bucket", testAccBosBucketResourceAttrName),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "logging.0.target_prefix", "logs/"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.0.id", "test-lr"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.0.status", "enabled"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.0.resource.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.0.condition.0.time.0.date_greater_than", "2029-12-31T00:00:00Z"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.0.action.0.name", "DeleteObject"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "storage_class", "COLD"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "server_side_encryption_rule", "AES256"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "website.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "website.0.index_document", "index.html"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "website.0.error_document", "err.html"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "cors_rule.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "cors_rule.0.allowed_origins.0", "https://www.baidu.com"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "cors_rule.0.allowed_methods.0", "GET"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "cors_rule.0.max_age_seconds", "1800"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "copyright_protection.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "copyright_protection.0.resource.#", "1"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "location"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "owner_id"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "owner_name"),
				),
			},
			{
				ResourceName:            testAccBosBucketResourceName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccBosBucketConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBosBucketResourceName),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "bucket", testAccBosBucketResourceAttrName),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "acl", "public-read"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.0.id", "test-rc"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.0.status", "enabled"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.0.resource.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.0.destination.0.bucket", BaiduCloudTestBucketResourceAttrNamePrefix+"bucket-peer"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.0.replicate_deletes", "enabled"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "force_destroy", "true"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "logging.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "logging.0.target_bucket", testAccBosBucketResourceAttrName),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "logging.0.target_prefix", "logs-update/"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.0.id", "test-lr"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.0.status", "disabled"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.0.resource.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.0.condition.0.time.0.date_greater_than", "$(lastModified)+P7D"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.0.action.0.name", "DeleteObject"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "server_side_encryption_rule", "AES256"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "website.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "website.0.index_document", "index02.html"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "website.0.error_document", "err02.html"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "cors_rule.#", "1"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "cors_rule.0.allowed_origins.0", "https://www.baidu.com"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "cors_rule.0.allowed_methods.0", "POST"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "cors_rule.0.max_age_seconds", "1800"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "copyright_protection.0.resource.#", "1"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "location"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "owner_id"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "owner_name"),
				),
			},
			{
				Config: testAccBosBucketConfigUpdate02(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBaiduCloudDataSourceId(testAccBosBucketResourceName),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "bucket", testAccBosBucketResourceAttrName),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "acl", "public-read"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "replication_configuration.#", "0"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "logging.#", "0"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "lifecycle_rule.#", "0"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "website.#", "0"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "cors_rule.#", "0"),
					resource.TestCheckResourceAttr(testAccBosBucketResourceName, "copyright_protection.0.resource.#", "0"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "location"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "creation_date"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "owner_id"),
					resource.TestCheckResourceAttrSet(testAccBosBucketResourceName, "owner_name"),
				),
			},
		},
	})
}

func testAccBosBucketDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.BaiduClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != testAccBosBucketResourceType {
			continue
		}

		raw, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return bosClient.DoesBucketExist(rs.Primary.ID)
		})
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		exist, _ := raw.(bool)
		if !exist {
			continue
		}

		return WrapError(Error("Bos Bucket still exist"))
	}

	return nil
}

func testAccBosBucketConfig() string {
	return fmt.Sprintf(`
resource "baiducloud_bos_bucket" "my-bucket" {
  bucket = "%s"
  acl = "public-read-write"
}

resource "baiducloud_bos_bucket" "default" {
  bucket = "%s"
  acl = "public-read-write"

  replication_configuration {
    id = "test-rc"
    status = "enabled"
    resource = ["%s"]
    destination {
      bucket = "${baiducloud_bos_bucket.my-bucket.bucket}"
    }
	replicate_history {
	  bucket = "${baiducloud_bos_bucket.my-bucket.bucket}"
	  storage_class = "COLD"
	}
    replicate_deletes = "disabled"
  }

  force_destroy = false

  logging {
    target_bucket = "%s"
    target_prefix = "logs/"
  }

  lifecycle_rule {
    id = "test-lr"
    status =  "enabled"
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
`, BaiduCloudTestBucketResourceAttrNamePrefix+"bucket-peer", testAccBosBucketResourceAttrName,
		testAccBosBucketResourceAttrName+"/*", testAccBosBucketResourceAttrName,
		testAccBosBucketResourceAttrName+"/*", testAccBosBucketResourceAttrName+"/*",
	)
}

func testAccBosBucketConfigUpdate() string {
	return fmt.Sprintf(`
resource "baiducloud_bos_bucket" "my-bucket" {
  bucket = "%s"
  acl = "public-read-write"
}

resource "baiducloud_bos_bucket" "default" {
  bucket = "%s"
  acl = "public-read"

  replication_configuration {
    id = "test-rc"
    status = "enabled"
    resource = ["%s"]
    destination {
      bucket = "${baiducloud_bos_bucket.my-bucket.bucket}"
    }
	replicate_history {
	  bucket = "${baiducloud_bos_bucket.my-bucket.bucket}"
	  storage_class = "COLD"
	}
    replicate_deletes = "enabled"
  }

  force_destroy = true

  logging {
    target_bucket = "%s"
    target_prefix = "logs-update/"
  }

  lifecycle_rule {
    id = "test-lr"
    status =  "disabled"
    resource = ["%s"]
    condition {
      time {
        date_greater_than = "$(lastModified)+P7D"
      }
    }
    action {
      name = "DeleteObject"
    }
  }

  storage_class = "STANDARD"

  server_side_encryption_rule = "AES256"

  website{
    index_document = "index02.html"
    error_document = "err02.html"
  }

  cors_rule {
    allowed_origins = ["https://www.baidu.com"]
    allowed_methods = ["POST"]
    max_age_seconds = 1800
  }

  copyright_protection {
    resource = ["%s"]
  }
}
`, BaiduCloudTestBucketResourceAttrNamePrefix+"bucket-peer", testAccBosBucketResourceAttrName,
		testAccBosBucketResourceAttrName+"/*", testAccBosBucketResourceAttrName,
		testAccBosBucketResourceAttrName+"/*", testAccBosBucketResourceAttrName+"/update*",
	)
}

func testAccBosBucketConfigUpdate02() string {
	return fmt.Sprintf(`
resource "baiducloud_bos_bucket" "my-bucket" {
  bucket = "%s"
  acl = "public-read-write"
}

resource "baiducloud_bos_bucket" "default" {
  bucket = "%s"
  acl = "public-read"
}
`, BaiduCloudTestBucketResourceAttrNamePrefix+"bucket-peer", testAccBosBucketResourceAttrName,
	)
}
