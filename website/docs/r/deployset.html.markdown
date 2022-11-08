---
layout: "baiducloud"
subcategory: "Baidu Cloud Compute (BCC)"
page_title: "BaiduCloud: baiducloud_deployset"
sidebar_current: "docs-baiducloud-resource-deployset"
description: |-
  Use this resource to creat a deployset.
---

# baiducloud_deployset

Use this resource to creat a deployset.

## Example Usage

```hcl
resource "baiducloud_deployset" "default" {
  name     = "terraform-test"
  desc     = "test desc"
  strategy = "HOST_HA"
}
```

## Argument Reference

The following arguments are supported:

* `desc` - (Optional) Description of the deployset.
* `name` - (Optional) Name of the deployset. Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",".", the value must start with a letter, length 1-65.
* `strategy` - (Optional) Strategy of deployset.Available values are HOST_HA, RACK_HA and TOR_HA

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `az_intstance_statis_list` - Availability Zone Instance Statistics List.
  * `bbc_instance_cnt` - Count of BBC instance which is in the deployset.
  * `bbc_instance_ids` - IDs of BBC instance which is in the deployset.
  * `bcc_instance_cnt` - Count of BCC instance which is in the deployset.
  * `bcc_instance_ids` - IDs of BCC instance which is in the deployset.
  * `instance_count` - Count of instance which is in the deployset.
  * `instance_ids` - IDs of instance which is in the deployset.
  * `instance_total` - Total of instance which is in the deployset.
  * `zone_name` - Zone name of deployset.
* `concurrency` - concurrency of deployset.
* `short_id` - deployset short id.
* `uuid` - deployset uuid.


## Import

deployset can be imported, e.g.

```hcl
$ terraform import baiducloud_deployset.default id
```

