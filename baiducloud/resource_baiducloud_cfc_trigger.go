/*
Provide a resource to create a CFC Function Trigger.

Example Usage

```hcl
resource "baiducloud_cfc_trigger" "http-trigger" {
  source_type   = "cfc-http-trigger/v1/CFCAPI"
  target        = "function_brn"
  resource_path = "/aaabbs"
  method        = ["GET","PUT"]
  auth_type     = "iam"
}
```

```
*/
package baiducloud

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/cfc"
	"github.com/baidubce/bce-sdk-go/services/cfc/api"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

func resourceBaiduCloudCFCTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceBaiduCloudCFCTriggerCreate,
		Read:   resourceBaiduCloudCFCTriggerRead,
		Update: resourceBaiduCloudCFCTriggerUpdate,
		Delete: resourceBaiduCloudCFCTriggerDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"source_type": {
				Type:         schema.TypeString,
				Description:  "CFC Funtion Trigger source type, support bos/http/crontab/dueros/duedge/cdn",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"bos", "http", "crontab", "dueros", "duedge", "cdn"}, false),
			},
			"target": {
				Type:        schema.TypeString,
				Description: "CFC Function Trigger target, it should be function brn",
				Required:    true,
				ForceNew:    true,
			},
			"relation_id": {
				Type:        schema.TypeString,
				Description: "CFC Function Trigger relation id",
				Computed:    true,
			},
			"sid": {
				Type:        schema.TypeString,
				Description: "CFC Funtion Trigger sid",
				Computed:    true,
			},

			// bos & cdn trigger
			"status": {
				Type:             schema.TypeString,
				Description:      "CFC Funtion Trigger status if source_type is bos or cdn",
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"enabled", "disabled"}, false),
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"bos", "cdn"}),
			},
			// bos trigger
			"resource": {
				Type:             schema.TypeString,
				Description:      "CFC Function Trigger resource if source_type is bos",
				Optional:         true,
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"bos"}),
			},
			"bucket": {
				Type:             schema.TypeString,
				Description:      "CFC Function Trigger source bucket if source_type is bos",
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"bos"}),
			},
			"bos_event_type": {
				Type:        schema.TypeList,
				Description: "CFC Function Trigger bos event type",
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(BosTriggerEventTypeList, false),
				},
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"bos"}),
			},

			// bos && crontab trigger
			"name": {
				Type:             schema.TypeString,
				Description:      "CFC Function Trigger name if source_type is bos or crontab",
				Optional:         true,
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"bos", "crontab"}),
			},

			// cdn trigger
			"cdn_event_type": {
				Type:             schema.TypeString,
				Description:      "CFC Function Trigger cdn event type",
				Optional:         true,
				ValidateFunc:     validation.StringInSlice(CDNTriggerEventTypeList, false),
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"cdn"}),
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "CFC Function Trigger domain if source_type is cdn",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"cdn"}),
			},
			"remark": {
				Type:             schema.TypeString,
				Description:      "CFC Function Trigger remark if source_type is cdn",
				Optional:         true,
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"cdn"}),
			},

			// http trigger
			"resource_path": {
				Type:        schema.TypeString,
				Description: "CFC Function Trigger resource path if source_type is http",
				Optional:    true,
				ValidateFunc: func(i interface{}, s string) (strings []string, errors []error) {
					path := i.(string)
					if path[0] != '/' {
						errors = append(errors, fmt.Errorf("resource_path should start with '/' "))
					}

					return
				},
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"http"}),
			},
			"method": {
				Type:             schema.TypeSet,
				Description:      "CFC Function Trigger method if source_type is http",
				Optional:         true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
				},
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"http"}),
			},
			"auth_type": {
				Type:             schema.TypeString,
				Description:      "CFC Function Trigger auth type if source_type is http, support anonymous or iam",
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"anonymous", "iam"}, false),
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"http"}),
			},

			// crontab trigger
			"input": {
				Type:             schema.TypeString,
				Description:      "CFC Function Trigger input if source_type is crontab",
				Optional:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"crontab"}),
			},
			"schedule_expression": {
				Type:             schema.TypeString,
				Description:      "CFC Function Trigger schedule expression if source_type is crontab",
				Optional:         true,
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"crontab"}),
			},
			"enabled": {
				Type:             schema.TypeString,
				Description:      "CFC Function Trigger enabled if source_type is crontab",
				Optional:         true,
				ValidateFunc:     validation.StringInSlice([]string{"Enabled", "Disabled"}, false),
				DiffSuppressFunc: cfcTriggerSourceTypeSuppressFunc([]string{"crontab"}),
			},
		},
	}
}

func resourceBaiduCloudCFCTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	createArgs, err := buildBaiduCloudCreateCFCTriggerArgs(d)
	if err != nil {
		return WrapError(err)
	}

	action := "Create CFC Function " + createArgs.Target + " trigger " + string(createArgs.Source)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return client.CreateTrigger(createArgs)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		response, _ := raw.(*api.CreateTriggerResult)

		d.SetId(base64.StdEncoding.EncodeToString([]byte(response.Relation.RelationId)))
		d.Set("relation_id", response.Relation.RelationId)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_trigger", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudCFCTriggerRead(d, meta)
}

func resourceBaiduCloudCFCTriggerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)
	cfcService := CFCService{client}

	functionBrn := d.Get("target").(string)
	relationId := d.Get("relation_id").(string)
	action := "Query Function " + functionBrn + " trigger " + relationId

	functionRelation, err := cfcService.CFCGetTriggerByFunction(functionBrn, relationId)
	if err != nil {
		d.SetId("")
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_trigger", action, BCESDKGoERROR)
	}

	d.Set("relation_id", functionRelation.RelationId)
	d.Set("sid", functionRelation.Sid)
	d.Set("target", functionRelation.Target)

	if strings.HasPrefix(string(functionRelation.Source), "bos/") {
		d.Set("source_type", "bos")
	} else {
		d.Set("source_type", functionRelation.Source)
	}

	data := functionRelation.Data.(map[string]interface{})
	switch functionRelation.Source {
	case api.SourceType("bos"):
		if v, ok := data["Resource"].(string); ok {
			d.Set("resource", v)
		}
		d.Set("status", data["Status"].(string))
		d.Set("name", data["Name"].(string))
		d.Set("bos_event_type", data["EventType"].([]interface{}))
		d.Set("bucket", string(functionRelation.Source)[4:])
		d.Set("source_type", "bos")
	case api.SourceTypeHTTP:
		methods := strings.Split(data["Method"].(string), ",")
		if err := d.Set("method", methods); err != nil {
			return err
		}
		d.Set("resource_path", data["ResourcePath"].(string))
		d.Set("auth_type", data["AuthType"].(string))
		d.Set("source_type", "http")
	case api.SourceTypeCDN:
		d.Set("cdn_event_type", data["EventType"].(string))
		d.Set("status", data["Status"].(string))
		d.Set("source_type", "cdn")

		if v, ok := data["Remark"].(string); ok {
			d.Set("remark", v)
		}

		if v, ok := data["Domains"].([]interface{}); ok {
			d.Set("domains", v)
		}
	case api.SourceTypeCrontab:
		d.Set("name", data["Name"].(string))
		d.Set("input", data["Input"].(string))
		d.Set("schedule_expression", data["ScheduleExpression"].(string))
		d.Set("enabled", data["Enabled"].(string))
		d.Set("source_type", "crontab")
	case api.SourceTypeDuEdge:
		d.Set("source_type", "duedge")
	case api.SourceTypeDuerOS:
		d.Set("source_type", "dueros")
	}

	return nil
}

func resourceBaiduCloudCFCTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	source := d.Get("source_type").(string)
	if !checkCFCTriggerDataHasChange(d, source) {
		return nil
	}

	updateArgs, err := buildBaiduCloudUpdateCFCTriggerArgs(d)
	if err != nil {
		return WrapError(err)
	}
	action := "Update function " + updateArgs.Target + " trigger " + updateArgs.RelationId

	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		raw, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return client.UpdateTrigger(updateArgs)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, raw)

		response, _ := raw.(*api.UpdateTriggerResult)
		d.Set("relation_id", response.Relation.RelationId)

		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_trigger", action, BCESDKGoERROR)
	}

	return resourceBaiduCloudCFCTriggerRead(d, meta)
}

func resourceBaiduCloudCFCTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.BaiduClient)

	sourceType := d.Get("source_type").(string)
	deleteArgs := &api.DeleteTriggerArgs{
		Target:     d.Get("target").(string),
		Source:     SourceTypeMap[sourceType],
		RelationId: d.Get("relation_id").(string),
	}

	if sourceType == "bos" {
		if value, ok := d.GetOk("bucket"); ok {
			deleteArgs.Source = api.SourceType("bos/" + value.(string))
		}
	}

	action := "Delete CFC Function " + deleteArgs.Target + " trigger " + deleteArgs.RelationId
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithCFCClient(func(client *cfc.Client) (i interface{}, e error) {
			return nil, client.DeleteTrigger(deleteArgs)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{bce.EINTERNAL_ERROR}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		addDebug(action, raw)
		return nil
	})

	if err != nil {
		if IsExceptedErrors(err, ObjectNotFound) {
			return nil
		}

		return WrapErrorf(err, DefaultErrorMsg, "baiducloud_cfc_trigger", action, BCESDKGoERROR)
	}

	return nil
}

func buildBaiduCloudCreateCFCTriggerArgs(d *schema.ResourceData) (*api.CreateTriggerArgs, error) {
	sourceType := d.Get("source_type").(string)
	result := &api.CreateTriggerArgs{
		Source: SourceTypeMap[sourceType],
		Target: d.Get("target").(string),
	}

	data, err := buildBaiduCloudCreateCFCTriggerData(d, sourceType)
	if err != nil {
		return nil, err
	}
	result.Data = data
	if sourceType == "bos" {
		if value, ok := d.GetOk("bucket"); ok {
			result.Source = api.SourceType("bos/" + value.(string))
		} else {
			return nil, triggerDataRequiredError("bucket", sourceType)
		}
	}

	return result, nil
}

func buildBaiduCloudUpdateCFCTriggerArgs(d *schema.ResourceData) (*api.UpdateTriggerArgs, error) {
	sourceType := d.Get("source_type").(string)
	result := &api.UpdateTriggerArgs{
		Source:     SourceTypeMap[sourceType],
		Target:     d.Get("target").(string),
		RelationId: d.Get("relation_id").(string),
	}

	data, err := buildBaiduCloudCreateCFCTriggerData(d, sourceType)
	if err != nil {
		return nil, err
	}
	result.Data = data
	if sourceType == "bos" {
		if value, ok := d.GetOk("bucket"); ok {
			result.Source = api.SourceType("bos/" + value.(string))
		} else {
			return nil, triggerDataRequiredError("bucket", sourceType)
		}
	}

	return result, nil
}

func buildBaiduCloudCreateCFCTriggerData(d *schema.ResourceData, sourceType string) (interface{}, error) {
	switch sourceType {
	case "bos":
		data := &api.BosTriggerData{}
		if value, ok := d.GetOk("resource"); ok {
			data.Resource = value.(string)
		} else {
			return nil, triggerDataRequiredError("resource", sourceType)
		}

		if value, ok := d.GetOk("status"); ok {
			data.Status = value.(string)
		} else {
			return nil, triggerDataRequiredError("status", sourceType)
		}

		if value, ok := d.GetOk("bos_event_type"); ok {
			eventType := value.([]interface{})
			for _, e := range eventType {
				data.EventType = append(data.EventType, api.BosEventType(e.(string)))
			}
		} else {
			return nil, triggerDataRequiredError("event_type", sourceType)
		}

		if value, ok := d.GetOk("name"); ok {
			data.Name = value.(string)
		} else {
			return nil, triggerDataRequiredError("name", sourceType)
		}

		return data, nil

	case "http":
		data := &api.HttpTriggerData{}
		if value, ok := d.GetOk("resource_path"); ok {
			data.ResourcePath = value.(string)
		} else {
			return nil, triggerDataRequiredError("resource_path", sourceType)
		}

		if value, ok := d.GetOk("method"); ok {
			methods := make([]string, 0)
			for _, m := range value.(*schema.Set).List() {
				methods = append(methods,  m.(string))
			}
			data.Method = strings.Join(methods, ",")
		} else {
			return nil, triggerDataRequiredError("method", sourceType)
		}

		if value, ok := d.GetOk("auth_type"); ok {
			data.AuthType = value.(string)
		} else {
			return nil, triggerDataRequiredError("auth_type", sourceType)
		}
		return data, nil

	case "cdn":
		data := &api.CDNTriggerData{}
		if value, ok := d.GetOk("remark"); ok {
			data.Remark = value.(string)
		}

		if value, ok := d.GetOk("status"); ok {
			data.Status = value.(string)
		} else {
			return nil, triggerDataRequiredError("status", sourceType)
		}

		if value, ok := d.GetOk("dimains"); ok {
			eventType := value.([]interface{})
			for _, e := range eventType {
				data.Domains = append(data.Domains, e.(string))
			}
		}

		if value, ok := d.GetOk("cdn_event_type"); ok {
			data.EventType = api.CDNEventType(value.(string))
		} else {
			return nil, triggerDataRequiredError("name", sourceType)
		}
		return data, nil
	case "crontab":
		data := &api.CrontabTriggerData{
			Brn: d.Get("target").(string),
		}

		if value, ok := d.GetOk("name"); ok {
			data.Name = value.(string)
		} else {
			return nil, triggerDataRequiredError("name", sourceType)
		}

		if value, ok := d.GetOk("input"); ok {
			data.Input = value.(string)
		}

		if value, ok := d.GetOk("schedule_expression"); ok {
			data.ScheduleExpression = value.(string)
		} else {
			return nil, triggerDataRequiredError("schedule_expression", sourceType)
		}

		if value, ok := d.GetOk("enabled"); ok {
			data.Enabled = value.(string)
		} else {
			return nil, triggerDataRequiredError("enabled", sourceType)
		}

		return data, nil
	}

	return nil, nil
}

func checkCFCTriggerDataHasChange(d *schema.ResourceData, sourceType string) bool {
	switch sourceType {
	case "bos":
		return d.HasChange("resource") || d.HasChange("status") || d.HasChange("bos_event_type") || d.HasChange("name")
	case "http":
		return d.HasChange("resource_path") || d.HasChange("method") || d.HasChange("auth_type")
	case "cdn":
		return d.HasChange("remark") || d.HasChange("status") || d.HasChange("domains") || d.HasChange("cdn_event_type")
	case "crontab":
		return d.HasChange("name") || d.HasChange("enabled") || d.HasChange("input") || d.HasChange("brn") || d.HasChange("schedule_expression")
	}

	return false
}

func triggerDataRequiredError(param string, sourceType string) error {
	return fmt.Errorf("%s is required if CFC Trigger source is %v", param, sourceType)
}
