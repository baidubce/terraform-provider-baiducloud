---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "baiducloud_bec_vm_instance Resource - terraform-provider-baiducloud"
subcategory: "Baidu Edge Computing (BEC)"
description: |-
  Use this resource to manage BEC VM Instance.
  More information can be found in the Developer Guide https://cloud.baidu.com/doc/BEC/s/jknpo0evo.
---

# baiducloud_bec_vm_instance (Resource)

Use this resource to manage BEC VM Instance. 

More information can be found in the [Developer Guide](https://cloud.baidu.com/doc/BEC/s/jknpo0evo).

## Example Usage

```terraform
resource "baiducloud_bec_vm_instance" "example" {

    service_id = "s-jw3y1234"
    vm_name = "vm_name_example"
    host_name = "host-name_example"
    region_id = "cn-maanshan-ct"

    cpu = 4
    memory = 8
    image_type = "bcc"
    image_id = "m-sqj56gCj"

    system_volume {
        name = "system_volume"
        size_in_gb = 40
        volume_type = "NVME"
    }
    data_volume {
        name = "data_volume1"
        size_in_gb = 20
        volume_type = "SATA"
    }

    need_public_ip = true
    need_ipv6_public_ip = true
    bandwidth = 10

    dns_config {
        dns_type = "DEFAULT"
    }

    key_config {
        type = "bccKeyPair"
        bcc_key_pair_id_list = ["k-cTaMVJcD", "k-9QQp6luE"]
    }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cpu` (Number) CPU core count of the vm instance. At least 1 core.
- `dns_config` (Block List, Min: 1, Max: 1) DNS config of the vm instance. (see [below for nested schema](#nestedblock--dns_config))
- `image_id` (String) ID of the image.
- `key_config` (Block List, Min: 1, Max: 1) Password or keypair config of the vm instance. (see [below for nested schema](#nestedblock--key_config))
- `memory` (Number) Memory size (GB) of the vm instance. At least 1 GB.
- `region_id` (String) Node ID, composed of lowercase letters of [`country`-`city`-`isp`]. Can be obtained through data source `baiducloud_bec_nodes`.
- `service_id` (String) ID of the vm instance group.
- `system_volume` (Block List, Min: 1, Max: 1) System volume config of the vm instance. (see [below for nested schema](#nestedblock--system_volume))

### Optional

- `bandwidth` (Number) Public network bandwidth size (Mbps).
- `data_volume` (Block List) Data volume config of the vm instance. (see [below for nested schema](#nestedblock--data_volume))
- `host_name` (String) Host name of the vm instance. If empty, system will assign one.
- `image_type` (String) Valid values: `bec`(public image or bec custom image), `bcc`(bcc custom image)
- `need_ipv6_public_ip` (Boolean) Whether to open IPv6 public network. Defaults to `false`.
- `need_public_ip` (Boolean) Whether to open public network. Defaults to `false`.
- `spec` (String) Specification family.
- `vm_name` (String) Name of the vm instance. If empty, system will assign one.

### Read-Only

- `create_time` (String) Creation time of the vm instance.
- `id` (String) The ID of this resource.
- `internal_ip` (String) Local network IPv4 address of the vm instance.
- `ipv6_public_ip` (String) Public network IPv6 address of the vm instance.
- `public_ip` (String) Public network IPv4 address of the vm instance.
- `status` (String) Status of the vm instance. Possible values: `CREATING`, `RUNNING`, `STOPPING`, `STOPPED`, `RESTARTING`, `REINSTALLING`, `STARTING`, `IMAGING`, `FAILED`, `UNKNOWN`

<a id="nestedblock--dns_config"></a>
### Nested Schema for `dns_config`

Required:

- `dns_type` (String) DNS type. Valid values: `NONE`(no DNS config), `DEFAULT`(114.114.114.114 for domestic nodes, 8.8.8.8 for overseas nodes), `LOCAL`(local dns of node), `CUSTOMIZE`

Optional:

- `dns_address` (List of String) Custom DNS address.


<a id="nestedblock--key_config"></a>
### Nested Schema for `key_config`

Required:

- `type` (String) Valid values: `bccKeyPair`, `password`

Optional:

- `admin_pass` (String, Sensitive) Length of the password is limited to 8 to 32 characters. Letters, numbers and symbols must exist at the same time, and the symbols are limited to `!@#$%^+*()`
- `bcc_key_pair_id_list` (List of String) Key pair ID list.


<a id="nestedblock--system_volume"></a>
### Nested Schema for `system_volume`

Required:

- `name` (String) Name of the disk.
- `size_in_gb` (Number) Size (GB) of the disk.
- `volume_type` (String) Type of the disk. Valid values: `NVME`(SSD), `SATA`(HDD).

Read-Only:

- `pvc_name` (String) PVC name of the disk.


<a id="nestedblock--data_volume"></a>
### Nested Schema for `data_volume`

Required:

- `name` (String) Name of the disk.
- `size_in_gb` (Number) Size (GB) of the disk.
- `volume_type` (String) Type of the disk. Valid values: `NVME`(SSD), `SATA`(HDD).

Read-Only:

- `pvc_name` (String) PVC name of the disk.

## Import

Import is supported using the following syntax:

```shell
terraform import baiducloud_bec_vm_instance.example vm-abs1sd13
```
