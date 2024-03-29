---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "baiducloud_bcc_key_pair Resource - terraform-provider-baiducloud"
subcategory: "Baidu Cloud Compute (BCC)"
description: |-
  Use this resource to manage BCC key pair.
  More information can be found in the Developer Guide https://cloud.baidu.com/doc/BCC/s/ykckicewc.
---

# baiducloud_bcc_key_pair (Resource)

Use this resource to manage BCC key pair. 

More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/BCC/s/ykckicewc).

## Example Usage

```terraform
# Create a new key pair and save the private key to a file
resource "baiducloud_bcc_key_pair" "example" {
    name = "example-key-pair"
    description = "created by terraform"
    private_key_file = "./private-key.txt"
}

# Import an existing public key
resource "baiducloud_bcc_key_pair" "example" {
    name = "example-key-pair"
    description = "created by terraform"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCI6n..."
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of key pair.

### Optional

- `description` (String) The description of key pair.
- `private_key_file` (String) The path of the file in which to save the private key.
- `public_key` (String) The public key of keypair. This field can be set to import an existing public key.

### Read-Only

- `created_time` (String) The creation time of key pair.
- `fingerprint` (String) The fingerprint of key pair.
- `id` (String) The ID of this resource.
- `instance_count` (Number) The number of instances bound to key pair.
- `region_id` (String) The id of the region to which key pair belongs.

## Import

Import is supported using the following syntax:

```shell
terraform import baiducloud_bcc_key_pair.example k-7U8ZiXP8
```
