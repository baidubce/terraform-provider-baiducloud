---
layout: "baiducloud"
subcategory: "BBC"
page_title: "BaiduCloud: baiducloud_bbc_instance"
sidebar_current: "docs-baiducloud-resource-bbc_instance"
description: |-
  Use this resource to create a BBC instance.
---

# baiducloud_bbc_instance

Use this resource to create a BBC instance.

~> **NOTE:** The terminate operation of BBC does NOT take effect immediately，maybe takes for several minites.

## Example Usage

```hcl
data "baiducloud_bbc_images" "bbc_images" {
  image_type = "BbcSystem"
  os_name    = "CentOS"
}
data "baiducloud_security_groups" "sg" {
  filter {
    name   = "name"
    values = ["default"]
  }
}
data "baiducloud_subnets" "subnets" {
  filter {
    name   = "zone_name"
    values = ["cn-bj-d"]
  }
  filter {
    name   = "name"
    values = ["系统预定义子网D"]
  }
}
data "baiducloud_bbc_flavors" "bbc_flavors" {
  filter {
    name   = "flavor_id"
    values = ["BBC-I4-01S"]
  }
}
resource "baiducloud_bbc_instance" "bbc_instance2" {
  action         = "start"
  payment_timing = "Postpaid"
  flavor_id            = "${data.baiducloud_bbc_flavors.bbc_flavors.flavors.0.flavor_id}"
  image_id             = "${data.baiducloud_bbc_images.bbc_images.images.0.id}"
  name                 = "terraform_test1"
  purchase_count       = 1
  raid                 = "Raid5"
  zone_name            = "cn-bj-d"
  root_disk_size_in_gb = 40
  security_groups      = [
    "${data.baiducloud_security_groups.sg.security_groups.0.id}",
    "${data.baiducloud_security_groups.sg.security_groups.1.id}",
  ]
  tags = {
    "testKey" = "terraform_test"
  }
  description = "terraform_test"
  admin_pass  = "terraform123456"
}
```
If you want to create a prepaid BBC, use following properties
```hcl
payment_timing = "Prepaid"
reservation             = {
  reservation_length    = 1
  reservation_time_unit = "Month"
}
```
## Argument Reference

The following arguments are supported:

* `flavor_id` - (Required) Id of the BBC Flavor.
* `image_id` - (Required) Id of the BBC Image.
* `name` - (Required) BBC name.Support for uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",".", the value must start with a letter, length 1-65.
* `purchase_count` - (Optional) The number of BBC instances created (purchased) in batch. It must be an integer greater than 0. It is an optional parameter. The default value is 1.
* `raid` - (Required) Type of the raid to start. Available values are Raid5, NoRaid.
* `root_disk_size_in_gb` - (Required) The system disk size of the BBC instance to be created.
* `security_groups` - (Required) Security groups of the bbc instance.
* `zone_name` - (Required, ForceNew) The naming convention of zonename is "country-region-availability area", in lowercase, for example, Beijing availability area A is "cn-bj-a"“
* `action` - (Optional) action.Available values are "start" and "stop" 
* `admin_pass` - (Optional) name.
* `auto_renew_time_unit` - (Optional) Monthly payment or annual payment, month is "month" and year is "year".
* `auto_renew_time` - (Optional) The automatic renewal time is 1-9 per month and 1-3 per year.
* `deploy_set_id` - (Optional) deploy set of bbc.
* `description` - (Optional) description.
* `hostname` - (Optional) Hostname is not specified by default. Hostname only supports lowercase letters, numbers and -. Special characters. It must start with a letter. Special symbols cannot be used consecutively. It does not support starting or ending with special symbols. The length is 2-64.
* `subnet_id` - (Optional) Id of bbc subnet.
* `tags` - (Optional, ForceNew) Tags, do not support modify
* `payment_timing` - (Optional) Payment timing of billing, which can be Prepaid or Postpaid. The default is Postpaid.
* `reservation` - (Optional) Reservation of the bbc instance.

The `reservation` object supports the following:

* `reservation_length` - (Optional) The reservation length that you will pay for your resource. It is valid when payment_timing is Prepaid. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36].Default value is 1.
* `reservation_time_unit` - (Optional) The reservation time unit that you will pay for your resource. It is valid when payment_timing is Prepaid. The value can only be month currently, which is also the default value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Create time of the BBC instance.
* `expire_time` - Expire time of the BBC instance.
* `has_alive` - Whether the instance has alive.
* `host_name` - Host name of instance, only supports lowercase letters, numbers and - . Special characters, must start with a letter, special symbols cannot be used consecutively, do not support the beginning or end of special symbols, length 2-64
* `internal_ip` - Internal IP assigned to the instance.
* `network_capacity_in_mbps` - Public network bandwidth(Mbps) of the instance.
* `public_ip` - Public ip of BBC instance.
* `raid_id` - Id of the raid.
* `rdma_ip` - Rdma IP of instance.
* `region` - Region of instance.
* `status` - The status of the instance.Include starting, running, stopped, deleted
* `switch_id` - Switch id the BBC instance associated.
* `uuid` - Uuid of the instance.


## Import

BBC instance can be imported, e.g.

```hcl
$ terraform import baiducloud_bbc_instance.my-server id
```

