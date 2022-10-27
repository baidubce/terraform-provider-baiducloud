---
layout: "baiducloud"
subcategory: "CCE"
page_title: "BaiduCloud: baiducloud_cce_versions"
sidebar_current: "docs-baiducloud-datasource-cce_versions"
description: |-
  Use this data source to list cce support kubernetes versions.
---

# baiducloud_cce_versions

Use this data source to list cce support kubernetes versions.

## Example Usage

```hcl
data "baiducloud_cce_versions" "default" {}

output "versions" {
  value = "${data.baiducloud_cce_versions.default.versions}"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional, ForceNew) Output file for saving result.
* `version_regex` - (Optional, ForceNew) Regex pattern of the search version

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `versions` - Useful kubernetes version list


