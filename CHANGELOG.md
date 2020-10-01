## 1.11.0 (Unreleased)
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
