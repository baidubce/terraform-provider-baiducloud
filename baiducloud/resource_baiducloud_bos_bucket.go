/*
Provide a resource to create a BOS Bucket.

Example Usage

```hcl
Private bucket
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"
  acl = "private"
}
```

Using replication configuration
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  replication_configuration {
    id = "test-rc"
    status = "enabled"
    resource = ["test-terraform/abc*"]
    destination {
      bucket = "test-terraform"
    }
    replicate_deletes = "disabled"
  }
}
```

Using logging
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  logging {
    target_bucket = "test-terraform"
    target_prefix = "logs/"
  }
}
```

Using lifecycle rule
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  lifecycle_rule {
	id = "test-id01"
	status =  "enabled"
	resource = ["test-terraform/abc*"]
	condition {
	  time {
	   date_greater_than = "2019-09-07T00:00:00Z"
	  }
	}
	action {
	  name = "DeleteObject"
	}
  }
  lifecycle_rule {
	id = "test-id02"
	status =  "enabled"
	resource = ["test-terraform/def*"]
	condition {
	  time {
	   date_greater_than = "$(lastModified)+P7D"
	  }
	}
	action {
	  name = "DeleteObject"
	}
  }
}
```

Using website
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  website{
    index_document = "index.html"
    error_document = "err.html"
  }
}
```

Using cors rule
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  cors_rule {
    allowed_origins = ["https://www.baidu.com"]
    allowed_methods = ["GET"]
    max_age_seconds = 1800
  }
}
```

Using copyright protection
```hcl
resource "baiducloud_bos_bucket" "default" {
  bucket = "${var.bucket}"

  copyright_protection {
    resource = ["test-terraform/abc*"]
  }
}
```

Import

BOS bucket can be imported, e.g.

```hcl
$ terraform import baiducloud_bos_bucket.default bucket_id
```
*/
package baiducloud

import (
	"encoding/json"
	"time"

	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudBosBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudBosBucketCreate,
		Read:   resourceBaiduCloudBosBucketRead,
		Update: resourceBaiduCloudBosBucketUpdate,
		Delete: resourceBaiduCloudBosBucketDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Description: "Name of the bucket.",
				Required:    true,
				ForceNew:    true,
			},
			"acl": {
				Type:         schema.TypeString,
				Description:  "Canned ACL to apply, available values are private, public-read and public-read-write. Default to private.",
				Optional:     true,
				Default:      api.CANNED_ACL_PRIVATE,
				ValidateFunc: validateBOSBucketACL(),
			},
			"replication_configuration": {
				Type:        schema.TypeList,
				Description: "Replication configuration of the BOS bucket.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the replication configuration.",
							Required:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the replication configuration. Valid values are enabled and disabled.",
							Required:    true,
							ValidateFunc: validation.StringInSlice([]string{
								api.STATUS_ENABLED,
								api.STATUS_DISABLED,
							}, false),
						},
						"resource": {
							Type:        schema.TypeSet,
							Description: "Resource of the replication configuration. The configuration format of the resource is {$bucket_name/<effective object prefix>}, which must start with \"$bucket_name\"+\"/\"",
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: resourceHash,
						},
						"destination": {
							Type:        schema.TypeList,
							Description: "Destination of the replication configuration.",
							Required:    true,
							MaxItems:    1,
							MinItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Description: "Destination bucket name of the replication configuration.",
										Required:    true,
									},
									"storage_class": {
										Type:         schema.TypeString,
										Description:  "Destination storage class of the replication configuration, the parameter does not need to be configured if it is consistent with the storage class of the source bucket, if you need to specify the storage class separately, it can be COLD, STANDARD, STANDARD_IA.",
										Optional:     true,
										ValidateFunc: validateBOSBucketRCStorageClass(),
									},
								},
							},
						},
						"replicate_history": {
							Type:        schema.TypeList,
							Description: "Configuration of the replicate history. The bucket name in replicate history needs to be the same as the bucket name in the destination above. After the history file is copied, all the objects of the inventory are copied to the destination bucket synchronously. The history file copy range is not referenced to the resource.",
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Description: "Destination bucket name of the replication configuration.",
										Required:    true,
									},
									"storage_class": {
										Type:         schema.TypeString,
										Description:  "Destination storage class of the replication configuration, the parameter does not need to be configured if it is consistent with the storage class of the source bucket, if you need to specify the storage class separately, it can be COLD, STANDARD, STANDARD_IA.",
										Optional:     true,
										ValidateFunc: validateBOSBucketRCStorageClass(),
									},
								},
							},
						},
						"replicate_deletes": {
							Type:        schema.TypeString,
							Description: "Whether to enable the delete synchronization, which can be enabled, disabled.",
							Required:    true,
							ValidateFunc: validation.StringInSlice([]string{
								api.STATUS_DISABLED,
								api.STATUS_ENABLED,
							}, false),
						},
					},
				},
			},

			"logging": {
				Type:        schema.TypeList,
				Description: "Settings of the bucket logging.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_bucket": {
							Type:        schema.TypeString,
							Description: "Target bucket name that will receive the log data.",
							Required:    true,
						},
						"target_prefix": {
							Type:        schema.TypeString,
							Description: "Target prefix for the log data.",
							Optional:    true,
						},
					},
				},
			},

			"lifecycle_rule": {
				Type:        schema.TypeList,
				Description: "Configuration of object lifecycle management.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "ID of the lifecycle rule. The id must be unique and cannot be repeated in the same bucket. The system will automatically generate one when it is not specified.",
							Optional:    true,
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of the lifecycle rule, which can be enabled, disabled. The rule cannot take effect when the status is disabled.",
							Required:    true,
							ValidateFunc: validation.StringInSlice([]string{
								api.STATUS_ENABLED,
								api.STATUS_DISABLED,
							}, false),
						},
						"resource": {
							Type:        schema.TypeSet,
							Description: "Resource of the lifecycle rule. For example, samplebucket/prefix/* will be valid for the object prefixed with prefix/ in samplebucket; samplebucket/* will be valid for all objects in samplebucket.",
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: resourceHash,
						},
						"condition": {
							Type:        schema.TypeList,
							Description: "Condition of the lifecycle rule, only the time form is supported currently.",
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time": {
										Type:        schema.TypeList,
										Description: "The condition time, implemented by the date_greater_than.",
										Required:    true,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												// It can be date or days in different situations.
												"date_greater_than": {
													Type:        schema.TypeString,
													Description: "Support absolute time date and relative time days. The absolute time date format is yyyy-mm-ddThh:mm:ssZ,eg. 2019-09-07T00:00:00Z. The absolute time is UTC time, which must end at 00:00:00(UTC 0 point); the description of relative time days follows ISO8601, and the minimum time granularity supported is days, such as: $(lastModified)+P7D indicates the time of object 7 days after last-modified.",
													Required:    true,
												},
											},
										},
									},
								},
							},
						},
						"action": {
							Type:        schema.TypeList,
							Description: "Action of the lifecycle rule.",
							Required:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Description:  "Action name of the lifecycle rule, which can be Transition, DeleteObject and AbortMultipartUpload.",
										Required:     true,
										ValidateFunc: validateBOSBucketLifecycleRuleActionName(),
									},
									"storage_class": {
										Type:         schema.TypeString,
										Description:  "When the action is Transition, it can be set to STANDARD_IA or COLD or ARCHIVE, indicating that it is changed from the original storage type to low frequency storage or cold storage or archive storage.",
										Optional:     true,
										ValidateFunc: validateBOSBucketLifecycleRuleActionStorage(),
									},
								},
							},
						},
					},
				},
			},

			"storage_class": {
				Type:         schema.TypeString,
				Description:  "Storage class of the BOS bucket, available values are STANDARD, STANDARD_IA, COLD or ARCHIVE.",
				Optional:     true,
				Default:      api.STORAGE_CLASS_STANDARD,
				ValidateFunc: validateBOSBucketStorageClass(),
			},

			"server_side_encryption_rule": {
				Type:         schema.TypeString,
				Description:  "Encryption rule for the server side, which can only be AES256 currently.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"AES256"}, false),
			},

			"website": {
				Type:        schema.TypeList,
				Description: "Website of the BOS bucket.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index_document": {
							Type:        schema.TypeString,
							Description: "Baiducloud BOS returns this index document when requests are made to the root domain or any of the subfolders.",
							Optional:    true,
						},
						"error_document": {
							Type:        schema.TypeString,
							Description: "An absolute path to the document to return in case of a 404 error.",
							Optional:    true,
						},
					},
				},
			},

			"cors_rule": {
				Type:        schema.TypeList,
				Description: "Configuration of the Cross-Origin Resource Sharing. Up to 100 rules are allowed per bucket, if there are multiple configurations, the execution order is from top to bottom.",
				Optional:    true,
				MaxItems:    100,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_headers": {
							Type:        schema.TypeList,
							Description: "Specifies which headers are allowed.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"allowed_methods": {
							Type:        schema.TypeList,
							Description: "Specifies which methods are allowed. Can be GET,PUT,DELETE,POST or HEAD.",
							Required:    true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateHttpMethod(),
							},
						},
						"allowed_origins": {
							Type:        schema.TypeList,
							Description: "Specifies which origins are allowed, containing up to one * wildcard.",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"allowed_expose_headers": {
							Type:        schema.TypeList,
							Description: "Specifies which expose headers are allowed.",
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"max_age_seconds": {
							Type:        schema.TypeInt,
							Description: "Specifies time in seconds that browser can cache the response for a preflight request.",
							Optional:    true,
						},
					},
				},
			},

			"copyright_protection": {
				Type:        schema.TypeList,
				Description: "Configuration of the copyright protection.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource": {
							Type:        schema.TypeSet,
							Description: "The resources to be protected for copyright.",
							Required:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"force_destroy": {
				Type:        schema.TypeBool,
				Description: "Whether to force delete the bucket and related objects when the bucket is not empty. Default to false.",
				Optional:    true,
				Default:     false,
			},

			// compute attribute
			"location": {
				Type:        schema.TypeString,
				Description: "Location of the BOS bucket.",
				Computed:    true,
			},
			"creation_date": {
				Type:        schema.TypeString,
				Description: "Creation date of the BOS bucket.",
				Computed:    true,
			},
			"owner_id": {
				Type:        schema.TypeString,
				Description: "Owner ID of the BOS bucket.",
				Computed:    true,
			},
			"owner_name": {
				Type:        schema.TypeString,
				Description: "Owner name of the BOS bucket.",
				Computed:    true,
			},
		},
	}
}

func resourceBaiduCloudBosBucketCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	bucket := d.Get("bucket").(string)
	action := "Create Bucket " + bucket

	location, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.PutBucket(bucket)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}
	addDebug(action, location)
	d.SetId(bucket)

	return resourceBaiduCloudBosBucketUpdate(d, meta)
}

func resourceBaiduCloudBosBucketRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bosService := &BosService{client}

	bucket := d.Id()
	action := "Query Bucket " + bucket

	// set bucket field first
	d.Set("bucket", bucket)

	// read bucket detail with ListBuckets api
	raw, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.ListBuckets()
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}
	addDebug(action, raw)

	result, _ := raw.(*api.ListBucketsResult)
	d.Set("owner_id", result.Owner.Id)
	d.Set("owner_name", result.Owner.DisplayName)
	for _, b := range result.Buckets {
		if b.Name == bucket {
			d.Set("location", b.Location)
			d.Set("creation_date", b.CreationDate)
			break
		}
	}

	// read bucket acl
	acl, err := bosService.resourceBaiduCloudBosBucketReadAcl(bucket)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}
	d.Set("acl", acl)

	// read replication configuration
	rc, err := bosService.resourceBaiduCloudBosBucketReadReplicationConfigure(bucket)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}
	d.Set("replication_configuration", rc)

	// read logging
	logging, err := bosService.resourceBaiduCloudBosBucketReadLogging(bucket)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}
	d.Set("logging", logging)

	// read lifecycle rules
	lcRules, err := bosService.resourceBaiduCloudBosBucketReadLifecycle(bucket)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}
	d.Set("lifecycle_rule", lcRules)

	// read storage class
	raw, err = client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetBucketStorageclass(bucket)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}
	d.Set("storage_class", raw.(string))

	// read server_side_encryption_rule
	raw, err = client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetBucketEncryption(bucket)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}
	d.Set("server_side_encryption_rule", raw.(string))

	// read website
	website, err := bosService.resourceBaiduCloudBosBucketReadWebsite(bucket)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}
	d.Set("website", website)

	// read cors
	cors, err := bosService.resourceBaiduCloudBosBucketReadCors(bucket)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}
	d.Set("cors_rule", cors)

	// read copyright protection
	copyright, err := bosService.resourceBaiduCloudBosBucketReadCopyright(bucket)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}
	d.Set("copyright_protection", copyright)

	return nil
}

func resourceBaiduCloudBosBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	bucket := d.Id()
	action := "Update Bucket " + bucket
	addDebug(action, bucket)

	d.Partial(true)

	// update acl
	if d.HasChange("acl") {
		if err := resourceBaiduCloudBosBucketAclUpdate(d, client); err != nil {
			return err
		}
		d.SetPartial("acl")
	}

	// update replication
	if d.HasChange("replication_configuration") {
		if err := resourceBaiduCloudBosBucketReplicationConfigurationUpdate(d, client); err != nil {
			return err
		}
		d.SetPartial("replication_configuration")
	}

	if d.HasChange("force_destroy") {
		d.SetPartial("force_destroy")
	}

	// update logging
	if d.HasChange("logging") {
		if err := resourceBaiduCloudBosBucketLoggingUpdate(d, client); err != nil {
			return err
		}
		d.SetPartial("logging")
	}

	// update lifecycle rules
	if d.HasChange("lifecycle_rule") {
		if err := resourceBaiduCloudBosBucketLifecycleUpdate(d, client); err != nil {
			return err
		}
		d.SetPartial("lifecycle_rule")
	}

	// update storage class
	if d.HasChange("storage_class") {
		sc := d.Get("storage_class").(string)
		_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.PutBucketStorageclass(bucket, sc)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
		}
		d.SetPartial("storage_class")
	}

	// update server_side_encryption_rule
	if d.HasChange("server_side_encryption_rule") {
		if err := resourceBaiduCloudBosBucketEncryptionUpdate(d, client); err != nil {
			return err
		}
		d.SetPartial("server_side_encryption_rule")
	}

	// update website
	if d.HasChange("website") {
		if err := resourceBaiduCloudBosBucketWebsiteUpdate(d, client); err != nil {
			return err
		}
		d.SetPartial("website")
	}

	// update cors
	if d.HasChange("cors_rule") {
		if err := resourceBaiduCloudBosBucketCorsUpdate(d, client); err != nil {
			return err
		}
		d.SetPartial("cors_rule")
	}

	// update copyright protection
	if d.HasChange("copyright_protection") {
		if err := resourceBaiduCloudBosBucketCopyrightProtectionUpdate(d, client); err != nil {
			return err
		}
		d.SetPartial("copyright_protection")
	}

	d.Partial(false)

	return resourceBaiduCloudBosBucketRead(d, meta)
}

func resourceBaiduCloudBosBucketDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	bosService := BosService{client}

	bucket := d.Id()
	action := "Delete Bucket " + bucket

	errRetry := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, errDelete := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.DeleteBucket(bucket)
		})
		if errDelete != nil {
			forceDestroy := d.Get("force_destroy").(bool)
			if !forceDestroy {
				return resource.NonRetryableError(errDelete)
			}

			if IsExceptedErrors(errDelete, []string{"BucketNotEmpty"}) {
				action += " objects"

				objects, err := bosService.ListAllObjects(bucket, "")
				if err != nil {
					return resource.NonRetryableError(err)
				}
				addDebug(action, objects)
				for _, obj := range objects {
					_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
						return nil, bosClient.DeleteObject(bucket, obj.Key)
					})
					if err != nil {
						return resource.NonRetryableError(err)
					}
				}

				return resource.RetryableError(errDelete)
			}

			if IsExceptedErrors(errDelete, []string{"ReplicationStatusNotEmpty"}) {
				_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
					return nil, bosClient.DeleteBucketReplication(bucket, "")
				})
				if err != nil {
					return resource.NonRetryableError(err)
				}

				return resource.RetryableError(errDelete)
			}

			return resource.NonRetryableError(errDelete)
		}

		return nil
	})
	if errRetry != nil {
		return WrapErrorf(errRetry, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}

	return nil
}

func flattenBaiduCloudBucketReplicationConfiguration(getBucketReplicationResult *api.GetBucketReplicationResult) []map[string]interface{} {
	replicationConfiguration := make([]map[string]interface{}, 0, 1)

	if getBucketReplicationResult == nil {
		return replicationConfiguration
	}

	m := make(map[string]interface{})

	m["id"] = getBucketReplicationResult.Id
	m["status"] = getBucketReplicationResult.Status
	m["resource"] = getBucketReplicationResult.Resource

	if getBucketReplicationResult.Destination != nil {
		des := make(map[string]interface{})
		des["bucket"] = getBucketReplicationResult.Destination.Bucket
		des["storage_class"] = getBucketReplicationResult.Destination.StorageClass

		m["destination"] = []interface{}{des}
	}

	if getBucketReplicationResult.ReplicateHistory != nil {
		rh := make(map[string]interface{})
		rh["bucket"] = getBucketReplicationResult.ReplicateHistory.Bucket
		rh["storage_class"] = getBucketReplicationResult.ReplicateHistory.StorageClass

		m["replicate_history"] = []interface{}{rh}
	}

	m["replicate_deletes"] = getBucketReplicationResult.ReplicateDeletes

	replicationConfiguration = append(replicationConfiguration, m)

	return replicationConfiguration
}

func resourceHash(v interface{}) int {
	return hashcode.String(v.(string))
}

func resourceBaiduCloudBosBucketAclUpdate(d *schema.ResourceData, client *connectivity.BaiduClient) error {
	action := "Update BOS Bucket acl"
	bucket := d.Get("bucket").(string)

	acl := d.Get("acl").(string)
	if _, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return nil, bosClient.PutBucketAclFromCanned(bucket, acl)
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}

	return nil
}

func resourceBaiduCloudBosBucketReplicationConfigurationUpdate(d *schema.ResourceData, client *connectivity.BaiduClient) error {
	bucket := d.Get("bucket").(string)
	o, n := d.GetChange("replication_configuration")

	action := "Update BOS Bucket replication configuration"

	oldRC := o.([]interface{})
	newRC := n.([]interface{})

	// delete old replication configuration
	if len(oldRC) != 0 {
		status, ok := oldRC[0].(map[string]interface{})["status"]
		if ok && (status == api.STATUS_ENABLED || len(newRC) == 0) {
			_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
				return nil, bosClient.DeleteBucketReplication(bucket, "")
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
			}
		}
	}

	// create new replication configuration
	if len(newRC) != 0 {
		rc := newRC[0].(map[string]interface{})
		args := &api.PutBucketReplicationArgs{
			Destination: &api.BucketReplicationDescriptor{},
		}

		args.Id = rc["id"].(string)
		args.Status = rc["status"].(string)

		resource := rc["resource"].(*schema.Set).List()
		args.Resource = make([]string, len(resource))
		for i, res := range resource {
			args.Resource[i] = res.(string)
		}

		des, _ := rc["destination"].([]interface{})[0].(map[string]interface{})
		if desBucket, ok := des["bucket"]; ok {
			args.Destination.Bucket = desBucket.(string)
		}
		if sc, ok := des["storage_class"]; ok {
			args.Destination.StorageClass = sc.(string)
		}

		if rh, ok := rc["replicate_history"]; ok {
			hs := rh.([]interface{})

			if len(hs) > 0 {
				args.ReplicateHistory = &api.BucketReplicationDescriptor{}
				his := hs[0].(map[string]interface{})
				if hisBucket, ok := his["bucket"]; ok {
					args.ReplicateHistory.Bucket = hisBucket.(string)
				}
				if sc, ok := his["storage_class"]; ok {
					args.ReplicateHistory.StorageClass = sc.(string)
				}
			}
		}

		args.ReplicateDeletes = rc["replicate_deletes"].(string)

		_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.PutBucketReplicationFromStruct(bucket, args, "")
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
		}
	}

	return nil
}

func resourceBaiduCloudBosBucketLoggingUpdate(d *schema.ResourceData, client *connectivity.BaiduClient) error {
	action := "Update BOS Bucket logging"
	bucket := d.Get("bucket").(string)

	logging, ok := d.GetOk("logging")
	if !ok {
		// delete directly
		_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.DeleteBucketLogging(bucket)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
		}
	} else {
		// put logging again
		l := logging.([]interface{})[0].(map[string]interface{})

		args := &api.PutBucketLoggingArgs{}
		args.TargetBucket = l["target_bucket"].(string)
		if pre, ok := l["target_prefix"]; ok {
			args.TargetPrefix = pre.(string)
		}

		_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.PutBucketLoggingFromStruct(bucket, args)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
		}
	}

	return nil
}

func resourceBaiduCloudBosBucketLifecycleUpdate(d *schema.ResourceData, client *connectivity.BaiduClient) error {
	action := "Update BOS Bucket lifecycle rules"
	bucket := d.Get("bucket").(string)

	lifecycleRules := d.Get("lifecycle_rule").([]interface{})
	if len(lifecycleRules) == 0 {
		// delete directly
		_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.DeleteBucketLifecycle(bucket)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
		}
	}

	rules := make([]*api.LifecycleRuleType, 0, len(lifecycleRules))
	for _, lifecycleRule := range lifecycleRules {
		lcr := lifecycleRule.(map[string]interface{})

		rule := &api.LifecycleRuleType{}

		if id, ok := lcr["id"]; ok {
			rule.Id = id.(string)
		}
		if status, ok := lcr["status"]; ok {
			rule.Status = status.(string)
		}
		if resource, ok := lcr["resource"]; ok {
			res := resource.(*schema.Set).List()
			rule.Resource = make([]string, 0, len(res))
			for _, r := range res {
				rule.Resource = append(rule.Resource, r.(string))
			}
		}
		if condition, ok := lcr["condition"]; ok {
			rule.Condition = api.LifecycleConditionType{
				Time: api.LifecycleConditionTimeType{},
			}
			cond := condition.([]interface{})[0].(map[string]interface{})
			if tim, ok := cond["time"]; ok {
				dateGreaterThan := tim.([]interface{})[0].(map[string]interface{})
				rule.Condition.Time.DateGreaterThan = dateGreaterThan["date_greater_than"].(string)
			}
		}
		if action, ok := lcr["action"]; ok {
			rule.Action = api.LifecycleActionType{}
			act := action.([]interface{})[0].(map[string]interface{})
			if name, ok := act["name"]; ok {
				rule.Action.Name = name.(string)
			}
			if sc, ok := act["storage_class"]; ok {
				rule.Action.StorageClass = sc.(string)
			}
		}

		rules = append(rules, rule)
	}

	args := struct {
		Rules []*api.LifecycleRuleType `json:"rule"`
	}{
		Rules: rules,
	}
	jsonBytes, jsonErr := json.Marshal(args)
	if jsonErr != nil {
		return WrapErrorf(jsonErr, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}

	_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return nil, bosClient.PutBucketLifecycleFromString(bucket, string(jsonBytes))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}

	return nil
}

func resourceBaiduCloudBosBucketEncryptionUpdate(d *schema.ResourceData, client *connectivity.BaiduClient) error {
	action := "Update BOS Bucket encryption"
	bucket := d.Get("bucket").(string)

	var err error
	rule, ok := d.GetOk("server_side_encryption_rule")
	if ok {
		// add rule
		_, err = client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.PutBucketEncryption(bucket, rule.(string))
		})

	} else {
		// delete rule
		_, err = client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.DeleteBucketEncryption(bucket)
		})
	}
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}

	return nil
}

func resourceBaiduCloudBosBucketWebsiteUpdate(d *schema.ResourceData, client *connectivity.BaiduClient) error {
	bucket := d.Get("bucket").(string)
	action := "Update BOS Bucket website"

	var err error
	website, ok := d.GetOk("website")
	if ok {
		// put directly
		web := website.([]interface{})[0].(map[string]interface{})

		args := &api.PutBucketStaticWebsiteArgs{
			Index:    web["index_document"].(string),
			NotFound: web["error_document"].(string),
		}
		_, err = client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.PutBucketStaticWebsiteFromStruct(bucket, args)
		})
	} else {
		// delete website
		_, err = client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.DeleteBucketStaticWebsite(bucket)
		})
	}
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}

	return nil
}

func resourceBaiduCloudBosBucketCorsUpdate(d *schema.ResourceData, client *connectivity.BaiduClient) error {
	bucket := d.Get("bucket").(string)
	action := "Update BOS Bucket CORS"

	cors := d.Get("cors_rule").([]interface{})
	if len(cors) == 0 {
		// delete cors
		if _, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.DeleteBucketCors(bucket)
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
		}

		return nil
	}

	// put cors
	args := &api.PutBucketCorsArgs{
		CorsConfiguration: make([]api.BucketCORSType, 0, len(cors)),
	}
	for _, cor := range cors {
		bucketCors := api.BucketCORSType{}
		corsMap := cor.(map[string]interface{})

		for k, v := range corsMap {
			if k == "max_age_seconds" {
				bucketCors.MaxAgeSeconds = int64(v.(int))
				continue
			}

			items := make([]string, 0, len(v.([]interface{})))
			for _, item := range v.([]interface{}) {
				if str, ok := item.(string); ok {
					items = append(items, str)
				}
			}
			switch k {
			case "allowed_headers":
				bucketCors.AllowedHeaders = items
			case "allowed_methods":
				bucketCors.AllowedMethods = items
			case "allowed_origins":
				bucketCors.AllowedOrigins = items
			case "allowed_expose_headers":
				bucketCors.AllowedExposeHeaders = items
			}
		}

		args.CorsConfiguration = append(args.CorsConfiguration, bucketCors)
	}

	_, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return nil, bosClient.PutBucketCorsFromStruct(bucket, args)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}

	return nil
}

func resourceBaiduCloudBosBucketCopyrightProtectionUpdate(d *schema.ResourceData, client *connectivity.BaiduClient) error {
	bucket := d.Get("bucket").(string)
	action := "Update BOS Bucket copyright protection"

	rawCP := d.Get("copyright_protection").([]interface{})
	if len(rawCP) == 0 {
		// delete
		if _, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
			return nil, bosClient.DeleteBucketCopyrightProtection(bucket)
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
		}

		return nil
	}

	rawMap := rawCP[0].(map[string]interface{})
	rawResources := rawMap["resource"].(*schema.Set).List()
	resources := make([]string, 0, len(rawResources))
	for _, rawRes := range rawResources {
		resources = append(resources, rawRes.(string))
	}

	if _, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return nil, bosClient.PutBucketCopyrightProtection(bucket, resources...)
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_bucket", action, BCESDKGoERROR)
	}

	return nil
}
