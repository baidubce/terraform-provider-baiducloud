---
layout: "baiducloud"
subcategory: "BCC"
page_title: "BaiduCloud: baiducloud_snapshot"
sidebar_current: "docs-baiducloud-resource-snapshot"
description: |-
  Provide a resource to create a Snapshot.
---

# baiducloud_snapshot

Provide a resource to create a Snapshot.

## Example Usage

```hcl
resource "baiducloud_snapshot" "my-snapshot" {
  name        = "${var.name}"
  description = "${var.description}"
  volume_id   = "v-Trb3rQXa"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of the snapshot, which supports uppercase and lowercase letters, numbers, Chinese and special characters, such as "-","_","/",",".", and the value must start with a letter, length 1-65.
* `volume_id` - (Required, ForceNew) Volume id of the snapshot, this value will be nil if volume has been released.
* `description` - (Optional, ForceNew) Description of the snapshot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_method` - Creation method of the snapshot.
* `create_time` - Creation time of the snapshot.
* `size_in_gb` - Size of the snapshot in GB.
* `status` - Status of the snapshot.


## Import

Snapshot can be imported, e.g.

```hcl
$ terraform import baiducloud_snapshot.default id
```

