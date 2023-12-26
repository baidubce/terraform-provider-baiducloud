	## 1.19.26 (Unreleased)
	
## 1.19.25 (December 26, 2023)
NOTES:
- ADD ET Gateway.

## 1.19.24 (December 25, 2023)
NOTES:
- The nat gateway now supports the DNAT and SNAT eips.

## 1.19.23 (November 27, 2023)
NOTES:
- ADD SCS Security ip.
        
## 1.19.22 (November 24, 2023)
NOTES:
- The routing table now supports import.

## 1.19.21 (November 15, 2023)
NOTES:
- The routing table now supports the configuration of dedicated gateways, including single-line and multi-line routing.

## 1.19.20 (November 13, 2023)
NOTES:
- ADD RDS Security ip.

## 1.19.19 (October 17, 2023)
NOTES:
- update resource/baiducloud_cds storage_type description : add `enhanced_ssd_pl1`.

## 1.19.18 (September 5, 2023)
ENHANCEMENTS:
- resource/baiducloud_instance: Support modifying parameter `instance_spec`.

## 1.19.17 (September 1, 2023)
BUG FIXES:
- resource/baiducloud_bcc_key_pair: Fix an issue where key pair could not be created.

## 1.19.16 (September 1, 2023)
ENHANCEMENTS:
- Hong Kong region (hkg) can now configure BBC, EIP, CCE, BLB, and BOS.

## 1.19.15 (August 30, 2023)
BUG FIXES:
- resource/baiducloud_nat_gateway: Remove default value of parameter `spec` to fix an issue where enhanced gateway could not be created.

## 1.19.14 (August 29, 2023)
ENHANCEMENTS:
- resource/baiducloud_instance: Add parameter `stop_with_no_charge` to support stopping charging after shutdown for postpaid instance without local disks.

## 1.19.13 (August 22, 2023)
FEATURES:
- **New Resource:** `baiducloud_iam_access_key`
- **New Data Source:** `baiducloud_iam_access_keys`

## 1.19.12 (August 17, 2023)
FEATURES:
- **New Resource:** `baiducloud_bcc_key_pair`

ENHANCEMENTS:
- eip support region bd

## 1.19.11 (August 10, 2023)
FEATURES:
- **New Data Source:** `baiducloud_bcc_key_pairs`

ENHANCEMENTS:
- Update `baidubce/bce-sdk-go` to v0.9.155

## 1.19.10 (August 4, 2023)
BUG FIXES:
- resource/baiducloud_peer_conn_acceptor: Fix error when creating cross-account or cross-region resources.

## 1.19.9 (July 12, 2023)
ENHANCEMENTS:
- resource/baiducloud_instance: Add parameter `resource_group_id` to support specifying resource group when creating instance.
- Update `baidubce/bce-sdk-go` to v0.9.153 

## 1.19.8 (July 3, 2023)
BUG FIXES:
- resource/baiducloud_instance: Fix parameter `cds_auto_renew` does not take effect when set to true.

## 1.19.7 (April 25, 2023)
ENHANCEMENTS:
- resource/baiducloud_instance: Add parameter `hostname`, support modification.
- datasource/baiducloud_instance: Add parameters `keypair_id`, `auto_renew`, `instance_ids`, `instance_names`, 
`cds_ids`, `deploy_set_ids`, `security_group_ids`, `payment_timing`, `status`, `tags`, `vpc_id`, `private_ips`.

## 1.19.6 (April 04, 2023)
NOTES:
- Add parameter disk_io_type for rds instance

## 1.19.5 (March 01, 2023)
NOTES:
- Update ccev2 supported k8s version

## 1.19.4 (January 11, 2023)
NOTES:
- Fix some doc mistakes

## 1.19.3 (January 6, 2023)
NOTES:
- Optimization cds datasource

## 1.19.2 (January 3, 2023)
NOTES:
- Fix bcc instance autorenew mistakes

## 1.19.1 (December 30, 2022)
NOTES:
- Fix baiducloud CDS prepaid mistakes

## 1.19.0 (December 29, 2022)
NOTES:
- Support Baidu Edge Computing (BEC)

FEATURES:
- **New Resource:** `baiducloud_bec_vm_instance`
- **New Data Source:** `baiducloud_bec_nodes`
- **New Data Source:** `baiducloud_bec_vm_instances`

## 1.18.4 (December 16, 2022)
NOTES:
- fix baiducloud EIP auto renew mistakes

## 1.18.3 (December 2, 2022)
NOTES:
- fix baiducloud BCC prepaid mistakes

## 1.18.2 (November 14, 2022)
NOTES:
- Support baiducloud BLS Log Store

## 1.18.1 (November 14, 2022)
NOTES:
- Support baiducloud SCS security IPs

ENHANCEMENTS:
- resource/baiducloud_scs: Add parameter security IPs

## 1.18.0 (November 08, 2022)
NOTES:
- Support baiducloud Cloud File Storage (CFS)

FEATURES:
- **New Data Source:** `baiducloud_cfss`
- **New Data Source:** `baiducloud_cfs_mount_targets`
- **New Resource:** `baiducloud_cfs`
- **New Resource:** `baiducloud_cfs_mount_target`

## 1.17.1 (November 07, 2022)
NOTES:
- Support baiducloud SMS Template

## 1.17.0 (November 02, 2022)
NOTES:
- Support BaiduCloud SNIC (Service Network Interface Card)

FEATURES:
- **New Data Source:** `baiducloud_snics`
- **New Data Source:** `baiducloud_snic_public_services`
- **New Resource:** `baiducloud_snic`

## 1.16.3 (October 31, 2022)
NOTES:
- Support baiducloud SMS signature

## 1.16.2 (October 27, 2022)
NOTES:
- Support baiducloud ENI

FEATURES:
- **New Data Source:** `baiducloud_enis`
- **New Resource:** `baiducloud_eni`
- **New Resource:** `baiducloud_eni_attachment`

ENHANCEMENTS:
- resource/baiducloud_specs: Parameter additions and changes

## 1.16.1 (October 25, 2022)
BUG FIXES:
- Fix subcategory error for CDN/SCS docs.

## 1.16.0 (October 25, 2022)
FEATURES:
- **New Data Source:** `baiducloud_cdn_domain_certificate`
- **New Resource:** `baiducloud_cdn_domain_config_origin`
- **New Resource:** `baiducloud_cdn_domain_config_advanced`
- **New Resource:** `baiducloud_cdn_domain_config_https`

ENHANCEMENTS:
- resource/baiducloud_cdn_domain_config_cache: Merge configuration read API calls.
- resource/baiducloud_cdn_domain_config_acl: Merge configuration read API calls.
- resource/baiducloud_scs: Delete validation for `disk_flavor`.

## 1.15.12 (October 24, 2022)
NOTES:
- Support baiducloud VPN

FEATURES:
- **New Data Source:** `baiducloud_vpn_gateways`
- **New Data Source:** `baiducloud_vpn_conns`
- **New Resource:** `baiducloud_vpn_gateway`
- **New Resource:** `baiducloud_vpn_conn`

ENHANCEMENTS:
- BCE Go SDK update to v0.9.137

## 1.15.11 (October 17, 2022)
NOTES:
- Support baidu blb security group

## 1.15.10 (September 26, 2022)
NOTES:
- Support baidu blb backend server

## 1.15.9 (September 23, 2022)
NOTES:
- Support baiducloud deploy set

FEATURES:
- **New Data Source:** `baiducloud_deploysets`
- **New Resource:** `baiducloud_deployset`

ENHANCEMENTS:
- Support baiducloud BCC param deploy_set_ids

## 1.15.8 (September 19, 2022)
NOTES:
- Support baidu blb listener

## 1.15.7 (September 16, 2022)
FEATURES:
- **New Resource:** `baiducloud_cdn_domain_config_acl`

ENHANCEMENTS:
- resource/baiducloud_cdn_domain: Add `weight`, `isp` arguments for `origin`

## 1.15.6 (September 14, 2022)
NOTES:
- go sdk update to v0.9.135

ENHANCEMENTS:
- Support baiducloud NAT param cu_num

BUG FIXES:
- Fix when instance_spec are provided, cpu_count, instance_type and memory_size_in_gb still diff

## 1.15.5 (September 08, 2022)
FEATURES:
- **New Data Source:** `baiducloud_bbc_instances`
- **New Data Source:** `baiducloud_bbc_images`
- **New Data Source:** `baiducloud_bbc_flavors`
- **New Resource:** `baiducloud_bbc_instance`
- **New Resource:** `baiducloud_bbc_image`

ENHANCEMENTS:
- Support baiducloud BCC param instance_spec

BUG FIXES:
- Fix filter's bug.

## 1.15.3 (September 06, 2022)
NOTES:
- Support baidu blb

## 1.15.0 (August 24, 2022)
FEATURES:
- **New Data Source:** `baiducloud_cdn_domains`
- **New Resource:** `baiducloud_cdn_domain`
- **New Resource:** `baiducloud_cdn_domain_config_cache`

BUG FIXES:
- Fix provider crash when setting custom endpoint.

## 1.14.1 (August 10, 2022)
NOTES:
- Support baidu localdns

## 1.14.0 (July 11, 2022)
ENHANCEMENTS:
- resource/baiducloud_scs: Attribute `engine` now supports new value `PegaDB`
- resource/baiducloud_scs: Attribute `tag` now supports specifying when creating.
- resource/baiducloud_scs: Add new attributes `client_auth`, `store_type`, `enable_read_only`, `disk_type`, `replication_resize_type`, `reservation_length`, `reservation_time_unit`
- resource/baiducloud_scs: Add new attributes `disk_flavor`, `replication_info`, and both support modification
- resource/baiducloud_scs: Remove `ForceNew` behavior from attribute `engine_version`, `shard_num`, `vpc_id`
- resource/baiducloud_instance: Remove `ForceNew` behavior from attribute `card_count` 

BUG FIXES:
- resource/baiducloud_scs: Attribute `billing` is deprecated, use `payment_timing`, `reservation_length`, `reservation_time_unit` instead.

## 1.13.0 (July 07, 2022)
ENHANCEMENTS:
- BCC/VPC related service now support region `bd`(BaoDing), `fsh`(ShangHai), `hkg`(HongKong), `sin`(Singapore) 

## 1.12.9 (June 30, 2022)
NOTES:
- ADD bcc max parallelism note

## 1.12.8 (June 30, 2022)
NOTES:
- Repair bcc instance markdown

## 1.12.7 (June 30, 2022)
NOTES:
- Improve BCC purchasing efficiency

## 1.12.6 (June 01, 2022)
NOTES:
- BCC add user_data
- Repair markdown

## 1.12.5 (May 26, 2022)
NOTES:
- Add nat_snat_rule
- Update go sdk
- Update eip

## 1.12.4 (May 24, 2022)
NOTES:
- add terraform-registry-manifest.json

## 1.12.3 (May 24, 2022)
NOTES:
- update tag

## 1.12.2 (May 24, 2022)
NOTES:
- update Release

## 1.12.1 (May 23, 2022)
NOTES:
- Repair SCS Document

## 1.12.0 (August 12, 2021)

NOTES:
- Repair and delete the security group and check whether deletion is allowed
- When repairing and deleting VPC, check whether deletion is allowed
- After adding and modifying the BCC instance subnet, the status is changesubnet
- SCS delete and add isolated status
- The added status during BCC creation is the deleted instantaneous status

## 1.11.3 (April 23, 2021)

NOTES:
- provider: Remove ValidateFunc for argument bandwidth_in_mbps of eip

## 1.11.2 (April 22, 2021)

NOTES:
- provider: Fix go mod vendor for new bce-sdk-go v0.9.62

## 1.11.1 (April 22, 2021)

BUG FIXES:
- provider: Fix wrong argument verify when resize eip

## 1.11.0 (February 27, 2021)

FEATURES:
* **New Resource:** `resource_baiducloud_iam_user`
* **New Resource:** `resource_baiducloud_iam_group`
* **New Resource:** `resource_baiducloud_iam_group_membership`
* **New Resource:** `resource_baiducloud_iam_policy`
* **New Resource:** `resource_baiducloud_iam_user_policy_attachment`
* **New Resource:** `resource_baiducloud_iam_group_policy_attachment`

## 1.10.3 (November 30, 2020)

BUG FIXES:
- provider: Fix a bug that failed to create RDS read-only instance with self built VPC

## 1.10.2 (November 06, 2020)

BUG FIXES:
- provider: Fix wrong usage of nonRetryable error

NOTES:
- provider: Fix wrong information in the provider index document 

## 1.10.1 (October 10, 2020)

NOTES:
- provider: Add a note about the debugging usage of the provider.
- provider: Add a publish script for provider.

## 1.10.0 (September 30, 2020)

FEATURES:

* **New Resource:** `resource_baiducloud_ccev2_instance`

## 1.9.0 (September 04, 2020)

FEATURES:

* **New Data Source:** `data_source_baiducloud_ccev2_cluster_instances`
* **New Data Source:** `data_source_baiducloud_ccev2_instance_group_instances`

## 1.8.0 (August 25, 2020)

FEATURES:

* **New Data Source:** `data_source_baiducloud_ccev2_clusterip_cidr`
* **New Data Source:** `data_source_baiducloud_ccev2_container_cidr`

* **New Resource:** `resource_baiducloud_ccev2_cluster`
* **New Resource:** `resource_baiducloud_ccev2_instance_group`

## 1.7.0 (July 29, 2020)

FEATURES:

* **New Data Source:** `data_source_baiducloud_dtss`
* **New Resource:** `resource_baiducloud_dts`

## 1.6.1 (July 24, 2020)

NOTES:
- provider: Fix wrong information in the provider index document 

## 1.6.0 (July 23, 2020)

ENHANCEMENTS:
- provider: Support assume role

## 1.5.0 (July 14, 2020)

FEATURES:

* **New Data Source:** `data_source_baiducloud_rdss`
* **New Resource:** `resource_baiducloud_rds_account`
* **New Resource:** `resource_baiducloud_rds_instance`
* **New Resource:** `resource_baiducloud_rds_readonly_instance`

BUG FIXES:
- resource/resource_baiducloud_scs: Fix forceNew behavior of parameters "proxy_num"、"replication_num"、"port"、"engine_version"、"subnets"

## 1.4.1 (July 10, 2020)

ENHANCEMENTS:
- resource/baiducloud_instance: support keypair_id when creating instances
- resource/baiducloud_cce_cluster: support keypair_id when creating cluster nodes

## 1.4.0 (July 06, 2020)

FEATURES:

* **New Data Source:** `baiducloud_scs_specs`
* **New Data Source:** `baiducloud_scss`
* **New Resource:** `baiducloud_scs`

## 1.3.1 (July 06, 2020)

BUG FIXES:
- resource/baiducloud_cce_cluster: Fix test tf config error
- datasource/baiducloud_cce_cluster_nodes: Fix test tf config error
- datasource/baiducloud_cce_kubeconfig: Fix test tf config error

## 1.3.0 (July 05, 2020)

FEATURES:

* **New Data Source:** `baiducloud_cce_versions`
* **New Data Source:** `baiducloud_cce_container_net`
* **New Data Source:** `baiducloud_cce_cluster_nodes`
* **New Data Source:** `baiducloud_cce_kubeconfig`
* **New Resource:** `baiducloud_cce_cluster`

BUG FIXES:
- resource/baiducloud_cds: Fix the problem of before delete not check status
- resource/baiducloud_cds_attachment: Fix the problem of detach error if cds has been detached or deleted before
- resource/baiducloud_eip: Fix the problem of befor delete not check status
- resource/baiducloud_eip_association: Fix the problem of unbind error if eip has been unbind or deleted before
- resource/baiducloud_security_group_rule: Fix the problem of create rule error if protocol is all and port_range is not 1-65535
- datasource/baiducloud_instances: Fix the problem of wrong parameter "card_count" type

## 1.2.0 (April 13, 2020)

NOTES:
- resource/baiducloud_cfc: Fix wrong information in the baiducloud_cfc document

ENHANCEMENTS:
- resource/baiducloud_instance: support start/stop operation of the instance
- resource/baiducloud_eip: support start/stop auto renew if eip is prepaid

## 1.1.0 (March 12, 2020)

ENHANCEMENTS:
- datasource/*: datasource add support "filter"

BUG FIXES:
- datasource/baiducloud_zones: Fix the problems of unavailable parameter "name_regex" ([#1](https://github.com/terraform-providers/terraform-provider-baiducloud/issues/1))
- datasource/baiducloud_specs: Fix the problems of unavailable parameter "name_regex" ([#1](https://github.com/terraform-providers/terraform-provider-baiducloud/issues/1))
- datasource/baiducloud_vpcs: Fix the problems of parameter "vpc_id" in tf cannot be empty. ([#4](https://github.com/terraform-providers/terraform-provider-baiducloud/issues/4))

## 1.0.0 (February 19, 2020)

FEATURES:

* **New Data Source:** `baiducloud_certs`
* **New Data Source:** `baiducloud_eips`
* **New Data Source:** `baiducloud_appblb`
* **New Data Source:** `baiducloud_appblb_server_groups`
* **New Data Source:** `baiducloud_appblb_listeners`
* **New Data Source:** `baiducloud_instances`
* **New Data Source:** `baiducloud_cdss`
* **New Data Source:** `baiducloud_security_groups`
* **New Data Source:** `baiducloud_security_group_rules`
* **New Data Source:** `baiducloud_snapshots`
* **New Data Source:** `baiducloud_auto_snapshot_policies`
* **New Data Source:** `baiducloud_zones`
* **New Data Source:** `baiducloud_specs`
* **New Data Source:** `baiducloud_images`
* **New Data Source:** `baiducloud_vpcs`
* **New Data Source:** `baiducloud_subnets`
* **New Data Source:** `baiducloud_route_rules`
* **New Data Source:** `baiducloud_acls`
* **New Data Source:** `baiducloud_nat_gateways`
* **New Data Source:** `baiducloud_peer_conns`
* **New Data Source:** `baiducloud_bos_buckets`
* **New Data Source:** `baiducloud_bos_bucket_objects`
* **New Data Source:** `baiducloud_cfc_function`

* **New Resource:** `baiducloud_cert`
* **New Resource:** `baiducloud_eip`
* **New Resource:** `baiducloud_eip_association`
* **New Resource:** `baiducloud_appblb`
* **New Resource:** `baiducloud_appblb_server_group`
* **New Resource:** `baiducloud_appblb_listener`
* **New Resource:** `baiducloud_instance`
* **New Resource:** `baiducloud_security_group`
* **New Resource:** `baiducloud_security_group_rule`
* **New Resource:** `baiducloud_cds`
* **New Resource:** `baiducloud_cds_attachment`
* **New Resource:** `baiducloud_snapshot`
* **New Resource:** `baiducloud_auto_snapshot_policy`
* **New Resource:** `baiducloud_vpc`
* **New Resource:** `baiducloud_subnet`
* **New Resource:** `baiducloud_route_rule`
* **New Resource:** `baiducloud_acl`
* **New Resource:** `baiducloud_nat_gateway`
* **New Resource:** `baiducloud_peer_conn`
* **New Resource:** `baiducloud_peer_conn_acceptor`
* **New Resource:** `baiducloud_bos_bucket`
* **New Resource:** `baiducloud_bos_bucket_object`
* **New Resource:** `baiducloud_cfc_alias`
* **New Resource:** `baiducloud_cfc_function`
* **New Resource:** `baiducloud_cfc_trigger`
* **New Resource:** `baiducloud_cfc_version`
