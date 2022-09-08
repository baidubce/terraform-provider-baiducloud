---
layout: "baiducloud"
page_title: "BaiduCloud: baiducloud_bbc_instances"
sidebar_current: "docs-baiducloud-datasource-bbc_instances"
description: |-
  Use this data source to query BCC Instance list.
---

# baiducloud_bbc_instances

Use this data source to query BCC Instance list.

## Example Usage

```hcl
data "baiducloud_bbc_instances" "data_bbc_instance" {
  internal_ip = "172.16.16.4"
}

output "instances" {
 value = "${data.baiducloud_bbc_instances.data_bbc_instance.instances}"
}

```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional, ForceNew) only support filter string/int/bool value
* `internal_ip` - (Optional) internal ip.
* `output_file` - (Optional, ForceNew) Output file of the bbc instances search result
* `vpc_id` - (Optional) bbc vpc id.

The `filter` object supports the following:

* `name` - (Required) filter variable name
* `values` - (Required) filter variable value list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instances` - The result of the bbc instances list.
  * `create_time` - BBC create time.
  * `deployset_id` - The deployset ID of the BBC instance.
  * `description` - The description of the BBC instance.
  * `expire_time` - BBC expire time.
  * `flavor_id` - The flavor ID of the BBC instance.
  * `has_alive` - BBC instance has alive.
  * `host_id` - The host ID of the BBC instance.
  * `hostname` - The hostname of the BBC instance.
  * `image_id` - The image ID of the BBC instance.
  * `instance_id` - The ID of the BBC instance.
  * `instance_name` - The name of the BBC instance.
  * `internal_ip` - Internal IP assigned to the instance.
  * `network_capacity_in_mbps` - Public network bandwidth(Mbps) of the instance.
  * `payment_timing` - Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.
  * `public_ip` - The public IP of the BBC instance.
  * `rdma_ip` - The rdma IP of the BBC instance.
  * `region` - The region of the BBC instance.
  * `status` - The status of the instance.Include starting, running, stopped, deleted
  * `switch_id` - The switch ID of the BBC instance.
  * `uuid` - The UUID of the BBC instance.
  * `zone` - zone name which BBC instance belong to.


