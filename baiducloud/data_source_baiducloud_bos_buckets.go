/*
Use this data source to query BOS bucket list.

Example Usage

```hcl
data "baiducloud_bos_buckets" "default" {}

output "buckets" {
 value = "${data.baiducloud_bos_buckets.default.buckets}"
}
```
*/
package baiducloud

import (
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func dataSourceBaiduCloudBosBuckets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBaiduCloudBosBucketsRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Description: "Name of the bucket to retrieve.",
				Optional:    true,
			},
			"output_file": {
				Type:        schema.TypeString,
				Description: "Output file for saving result.",
				Optional:    true,
				ForceNew:    true,
			},

			// Attributes used for result
			"buckets": {
				Type:        schema.TypeList,
				Description: "List of buckets.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:        schema.TypeString,
							Description: "Name of the bucket.",
							Computed:    true,
						},
						"location": {
							Type:        schema.TypeString,
							Description: "Bucket location of the bucket.",
							Computed:    true,
						},
						"creation_date": {
							Type:        schema.TypeString,
							Description: "Creation date of the bucket.",
							Computed:    true,
						},

						"owner_id": {
							Type:        schema.TypeString,
							Description: "Owner id of the bucket.",
							Computed:    true,
						},
						"owner_name": {
							Type:        schema.TypeString,
							Description: "Owner name of the bucket.",
							Computed:    true,
						},

						"acl": {
							Type:        schema.TypeString,
							Description: "Acl of the bucket.",
							Computed:    true,
						},

						"replication_configuration": {
							Type:        schema.TypeList,
							Description: "Replication configuration of the bucket.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "ID of the replication configuration.",
										Computed:    true,
									},
									"status": {
										Type:        schema.TypeString,
										Description: "Status of the replication configuration.",
										Computed:    true,
									},
									"resource": {
										Type:        schema.TypeList,
										Description: "Resource of the replication configuration.",
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"destination": {
										Type:        schema.TypeList,
										Description: "Destination of the replication configuration.",
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:        schema.TypeString,
													Description: "Destination bucket name of the replication configuration.",
													Computed:    true,
												},
												"storage_class": {
													Type:        schema.TypeString,
													Description: "Destination storage class of the replication configuration.",
													Computed:    true,
												},
											},
										},
									},
									"replicate_history": {
										Type:        schema.TypeList,
										Description: "Configuration of the replicate history.",
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:        schema.TypeString,
													Description: "Destination bucket name of the replication configuration.",
													Computed:    true,
												},
												"storage_class": {
													Type:        schema.TypeString,
													Description: "Destination storage class of the replication configuration.",
													Computed:    true,
												},
											},
										},
									},
									"replicate_deletes": {
										Type:        schema.TypeString,
										Description: "Whether to enable the delete synchronization.",
										Computed:    true,
									},
								},
							},
						},

						"logging": {
							Type:        schema.TypeList,
							Description: "Logging of the bucket.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_bucket": {
										Type:        schema.TypeString,
										Description: "Target bucket name of the logging.",
										Computed:    true,
									},
									"target_prefix": {
										Type:        schema.TypeString,
										Description: "Target prefix of the logging.",
										Computed:    true,
									},
								},
							},
						},

						"lifecycle_rule": {
							Type:        schema.TypeList,
							Description: "Configuration of object lifecycle management.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "ID of the lifecycle rule.",
										Computed:    true,
									},
									"status": {
										Type:        schema.TypeString,
										Description: "Status of the lifecycle rule.",
										Computed:    true,
									},
									"resource": {
										Type:        schema.TypeList,
										Description: "Resource of the lifecycle rule.",
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"condition": {
										Type:        schema.TypeList,
										Description: "Condition of the lifecycle rule.",
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time": {
													Type:        schema.TypeList,
													Description: "Condition time, implemented by the date_greater_than.",
													Computed:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"date_greater_than": {
																Type:        schema.TypeString,
																Description: "Support absolute time date and relative time days.",
																Computed:    true,
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
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Description: "Action name of the lifecycle rule.",
													Computed:    true,
												},
												"storage_class": {
													Type:        schema.TypeString,
													Description: "Storage class of the action.",
													Computed:    true,
												},
											},
										},
									},
								},
							},
						},

						"storage_class": {
							Type:        schema.TypeString,
							Description: "Storage class of the bucket.",
							Computed:    true,
						},
						"server_side_encryption_rule": {
							Type:        schema.TypeString,
							Description: "Encryption of the bucket.",
							Computed:    true,
						},
						"website": {
							Type:        schema.TypeList,
							Description: "Website of the BOS bucket.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"index_document": {
										Type:        schema.TypeString,
										Description: "Baiducloud BOS returns this index document when requests are made to the root domain or any of the subfolders.",
										Computed:    true,
									},
									"error_document": {
										Type:        schema.TypeString,
										Description: "An absolute path to the document to return in case of a 404 error.",
										Computed:    true,
									},
								},
							},
						},

						"cors_rule": {
							Type:        schema.TypeList,
							Description: "Configuration of the Cross-Origin Resource Sharing.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_headers": {
										Type:        schema.TypeList,
										Description: "Indicate which headers are allowed.",
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"allowed_methods": {
										Type:        schema.TypeList,
										Description: "Indicate which methods are allowed.",
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"allowed_origins": {
										Type:        schema.TypeList,
										Description: "Indicate which origins are allowed.",
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"allowed_expose_headers": {
										Type:        schema.TypeList,
										Description: "Indicate which expose headers are allowed.",
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"max_age_seconds": {
										Type:        schema.TypeInt,
										Description: "Indicate time in seconds that browser can cache the response for a preflight request.",
										Computed:    true,
									},
								},
							},
						},

						"copyright_protection": {
							Type:        schema.TypeList,
							Description: "Configuration of the copyright protection.",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource": {
										Type:        schema.TypeList,
										Description: "Resources to be protected for copyright.",
										Computed:    true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceBaiduCloudBosBucketsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	var (
		bucket     string
		outputFile string
	)
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	}
	if v, ok := d.GetOk("output_file"); ok {
		outputFile = v.(string)
	}

	action := "Query BOS buckets " + bucket

	raw, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.ListBuckets()
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_buckets", action, BCESDKGoERROR)
	}

	result, _ := raw.(*api.ListBucketsResult)
	bucketsResult := make([]interface{}, 0)
	for _, buc := range result.Buckets {
		if bucket != "" && bucket != buc.Name {
			continue
		}

		bucMap, err := dataSourceBaiduCloudBosBucketsReadBucketData(meta, buc.Name)
		if err != nil {
			return err
		}

		bucMap["bucket"] = buc.Name
		bucMap["location"] = buc.Location
		bucMap["creation_date"] = buc.CreationDate
		bucMap["owner_id"] = result.Owner.Id
		bucMap["owner_name"] = result.Owner.DisplayName

		bucketsResult = append(bucketsResult, bucMap)
	}

	d.Set("buckets", bucketsResult)
	d.SetId(resource.UniqueId())

	if outputFile != "" {
		if err := writeToFile(outputFile, bucketsResult); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_buckets", action, BCESDKGoERROR)
		}
	}

	return nil
}

func dataSourceBaiduCloudBosBucketsReadBucketData(meta interface{}, bucket string) (map[string]interface{}, error) {
	action := "read bos bucket data " + bucket

	client := meta.(*connectivity.BaiduClient)
	bosService := &BosService{client}

	// read basic
	bucMap := make(map[string]interface{})

	// read bucket acl
	acl, err := bosService.resourceBaiduCloudBosBucketReadAcl(bucket)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_buckets", action, BCESDKGoERROR)
	}
	bucMap["acl"] = acl

	// read replication configuration
	rc, err := bosService.resourceBaiduCloudBosBucketReadReplicationConfigure(bucket)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_buckets", action, BCESDKGoERROR)
	}
	bucMap["replication_configuration"] = rc

	// read logging
	logging, err := bosService.resourceBaiduCloudBosBucketReadLogging(bucket)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_buckets", action, BCESDKGoERROR)
	}
	bucMap["logging"] = logging

	// read lifecycle rules
	lcRules, err := bosService.resourceBaiduCloudBosBucketReadLifecycle(bucket)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_buckets", action, BCESDKGoERROR)
	}
	bucMap["lifecycle_rule"] = lcRules

	// read storage class
	raw, err := client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetBucketStorageclass(bucket)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_buckets", action, BCESDKGoERROR)
	}
	bucMap["storage_class"] = raw.(string)

	// read server_side_encryption_rule
	raw, err = client.WithBosClient(func(bosClient *bos.Client) (i interface{}, e error) {
		return bosClient.GetBucketEncryption(bucket)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_buckets", action, BCESDKGoERROR)
	}
	bucMap["server_side_encryption_rule"] = raw.(string)

	// read website
	website, err := bosService.resourceBaiduCloudBosBucketReadWebsite(bucket)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_buckets", action, BCESDKGoERROR)
	}
	bucMap["website"] = website

	// read cors
	cors, err := bosService.resourceBaiduCloudBosBucketReadCors(bucket)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_buckets", action, BCESDKGoERROR)
	}
	bucMap["cors_rule"] = cors

	// read copyright protection
	copyright, err := bosService.resourceBaiduCloudBosBucketReadCopyright(bucket)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "baiducloud_bos_buckets", action, BCESDKGoERROR)
	}
	bucMap["copyright_protection"] = copyright

	return bucMap, nil
}
