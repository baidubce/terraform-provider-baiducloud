/*
The BaiduCloud provider is used to interact with many resources supported by BaiduCloud. The provider needs to be configured with the proper credentials before it can be used.

The BaiduCloud provider is used to interact with the many resources supported by [BaiduCloud](https://cloud.baidu.com).
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

Example Usage

```hcl
# Configure the BaiduCloud Provider
provider "baiducloud" {
  access_key  = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}
```

Resources List

Data Sources
  baiducloud_vpcs
  baiducloud_subnets
  baiducloud_route_rules
  baiducloud_acls
  baiducloud_nat_gateways
  baiducloud_peer_conns
  baiducloud_bos_buckets
  baiducloud_bos_bucket_objects
  baiducloud_appblbs
  baiducloud_appblb_listeners
  baiducloud_appblb_server_groups
  baiducloud_eips
  baiducloud_instances
  baiducloud_cdss
  baiducloud_security_groups
  baiducloud_security_group_rules
  baiducloud_snapshots
  baiducloud_auto_snapshot_policies
  baiducloud_zones
  baiducloud_specs
  baiducloud_images
  baiducloud_certs
  baiducloud_cfc_function
  baiducloud_scs_specs
  baiducloud_scss
  baiducloud_cce_versions
  baiducloud_cce_container_net
  baiducloud_cce_cluster_nodes
  baiducloud_cce_kubeconfig
  baiducloud_ccev2_container_cidr
  baiducloud_ccev2_clusterip_cidr
  baiducloud_ccev2_cluster_instances
  baiducloud_ccev2_instance_group_instances
  baiducloud_dtss

CERT Resources
  baiducloud_cert

EIP Resources
  baiducloud_eip
  baiducloud_eip_association

APPBLB Resources
  baiducloud_appblb
  baiducloud_appblb_server_group
  baiducloud_appblb_listener

BCC Resources
  baiducloud_instance
  baiducloud_security_group
  baiducloud_security_group_rule
  baiducloud_cds
  baiducloud_cds_attachment
  baiducloud_snapshot
  baiducloud_auto_snapshot_policy

VPC Resources
  baiducloud_vpc
  baiducloud_subnet
  baiducloud_route_rule
  baiducloud_acl
  baiducloud_nat_gateway
  baiducloud_nat_snat_rule
  baiducloud_peer_conn
  baiducloud_peer_conn_acceptor

BOS Resources
  baiducloud_bos_bucket
  baiducloud_bos_bucket_object

CFC Resources
  baiducloud_cfc_function
  baiducloud_cfc_alias
  baiducloud_cfc_trigger
  baiducloud_cfc_version

SCS Resources
  baiducloud_scs

DTS Resources
  baiducloud_dts

CCE Resources
  baiducloud_cce_cluster

CCEv2 Resources
  baiducloud_ccev2_cluster
  baiducloud_ccev2_instance_group

IAM Resources
  baiducloud_iam_user
  baiducloud_iam_group
  baiducloud_iam_group_membership
  baiducloud_iam_policy
  baiducloud_iam_user_policy_attachment
  baiducloud_iam_group_policy_attachment
*/
package baiducloud

import (
	"bytes"
	"fmt"
	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/service/cdn"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-baiducloud/baiducloud/connectivity"
)

const (
	PROVIDER_ACCESS_KEY = "BAIDUCLOUD_ACCESS_KEY"
	PROVIDER_SECRET_KEY = "BAIDUCLOUD_SECRET_KEY"
	PROVIDER_REGION     = "BAIDUCLOUD_REGION"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_ACCESS_KEY, nil),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(PROVIDER_SECRET_KEY, nil),
				Description: descriptions["secret_key"],
				Sensitive:   true,
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc(PROVIDER_REGION, nil),
				Description:  descriptions["region"],
				InputDefault: "bj",
			},
			"endpoints": endpointsSchema(),

			"assume_role": assumeRoleSchema(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"baiducloud_vpcs":                           dataSourceBaiduCloudVpcs(),
			"baiducloud_subnets":                        dataSourceBaiduCloudSubnets(),
			"baiducloud_route_rules":                    dataSourceBaiduCloudRouteRules(),
			"baiducloud_acls":                           dataSourceBaiduCloudAcls(),
			"baiducloud_nat_gateways":                   dataSourceBaiduCloudNatGateways(),
			"baiducloud_nat_snat_rules":                 dataSourceBaiduCloudNatSnatRules(),
			"baiducloud_peer_conns":                     dataSourceBaiduCloudPeerConns(),
			"baiducloud_bos_buckets":                    dataSourceBaiduCloudBosBuckets(),
			"baiducloud_bos_bucket_objects":             dataSourceBaiduCloudBosBucketObjects(),
			"baiducloud_appblbs":                        dataSourceBaiduCloudAppBLBs(),
			"baiducloud_appblb_listeners":               dataSourceBaiduCloudAppBLBListeners(),
			"baiducloud_appblb_server_groups":           dataSourceBaiduCloudAppBLBServerGroups(),
			"baiducloud_blbs":                           dataSourceBaiduCloudBLBs(),
			"baiducloud_blb_listeners":                  dataSourceBaiduCloudBLBListeners(),
			"baiducloud_blb_backend_servers":            dataSourceBaiduCloudBLBBackendServer(),
			"baiducloud_blb_securitygroups":             dataSourceBaiduCloudBLBSecurityGroups(),
			"baiducloud_certs":                          dataSourceBaiduCloudCerts(),
			"baiducloud_eips":                           dataSourceBaiduCloudEips(),
			"baiducloud_instances":                      dataSourceBaiduCloudInstances(),
			"baiducloud_cdss":                           dataSourceBaiduCloudCDSs(),
			"baiducloud_security_groups":                dataSourceBaiduCloudSecurityGroups(),
			"baiducloud_security_group_rules":           dataSourceBaiduCloudSecurityGroupRules(),
			"baiducloud_snapshots":                      dataSourceBaiduCloudSnapshots(),
			"baiducloud_auto_snapshot_policies":         dataSourceBaiduCloudAutoSnapshotPolicies(),
			"baiducloud_zones":                          dataSourceBaiduCloudZones(),
			"baiducloud_specs":                          dataSourceBaiduCloudSpecs(),
			"baiducloud_images":                         dataSourceBaiduCloudImages(),
			"baiducloud_cfc_function":                   dataSourceBaiduCloudCFCFunction(),
			"baiducloud_scs_specs":                      dataSourceBaiduCloudScsSpecs(),
			"baiducloud_scss":                           dataSourceBaiduCloudScss(),
			"baiducloud_cce_versions":                   dataSourceBaiduCloudCceKubernetesVersion(),
			"baiducloud_cce_container_net":              dataSourceBaiduCloudCceContainerNet(),
			"baiducloud_cce_cluster_nodes":              dataSourceBaiduCloudCCEClusterNodes(),
			"baiducloud_ccev2_container_cidr":           dataSourceBaiduCloudCCEv2ContainerCIDRs(),
			"baiducloud_ccev2_clusterip_cidr":           dataSourceBaiduCloudCCEv2ClusterIPCidrs(),
			"baiducloud_ccev2_cluster_instances":        dataSourceBaiduCloudCCEv2ClusterInstances(),
			"baiducloud_ccev2_instance_group_instances": dataSourceBaiduCloudCCEv2InstanceGroupInstances(),
			"baiducloud_cce_kubeconfig":                 dataSourceBaiduCloudCceKubeConfig(),
			"baiducloud_rdss":                           dataSourceBaiduCloudRdss(),
			"baiducloud_dtss":                           dataSourceBaiduCloudDtss(),
			"baiducloud_cdn_domains":                    cdn.DataSourceDomains(),
			"baiducloud_localdns_privatezones":          dataSourceBaiduCloudLocalDnsPrivateZones(),
			"baiducloud_localdns_vpcs":                  dataSourceBaiduCloudLocalDnsVpcs(),
			"baiducloud_localdns_records":               dataSourceBaiduCloudPrivateZoneDNSRecords(),
			"baiducloud_bbc_images":                     dataSourceBaiduCloudBbcImages(),
			"baiducloud_bbc_flavors":                    dataSourceBaiduCloudBbcFlavors(),
			"baiducloud_bbc_instances":                  dataSourceBaiduCloudBbcInstances(),
			"baiducloud_deploysets":                     dataSourceBaiduCloudDeploySets(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"baiducloud_instance":                    resourceBaiduCloudInstance(),
			"baiducloud_cds":                         resourceBaiduCloudCDS(),
			"baiducloud_cds_attachment":              resourceBaiduCloudCDSAttachment(),
			"baiducloud_snapshot":                    resourceBaiduCloudSnapshot(),
			"baiducloud_auto_snapshot_policy":        resourceBaiduCloudAutoSnapshotPolicy(),
			"baiducloud_vpc":                         resourceBaiduCloudVpc(),
			"baiducloud_subnet":                      resourceBaiduCloudSubnet(),
			"baiducloud_route_rule":                  resourceBaiduCloudRouteRule(),
			"baiducloud_security_group":              resourceBaiduCloudSecurityGroup(),
			"baiducloud_security_group_rule":         resourceBaiduCloudSecurityGroupRule(),
			"baiducloud_eip":                         resourceBaiduCloudEip(),
			"baiducloud_eip_association":             resourceBaiduCloudEipAssociation(),
			"baiducloud_acl":                         resourceBaiduCloudAcl(),
			"baiducloud_nat_gateway":                 resourceBaiduCloudNatGateway(),
			"baiducloud_nat_snat_rule":               resourceBaiduCloudNatSnatRule(),
			"baiducloud_blb":                         resourceBaiduCloudBLB(),
			"baiducloud_blb_listener":                resourceBaiduCloudBlbListener(),
			"baiducloud_blb_backend_server":          resourceBaiduCloudBlbBackendServer(),
			"baiducloud_blb_securitygroup":           resourceBaiduCloudBlbSecurityGroup(),
			"baiducloud_appblb":                      resourceBaiduCloudAppBLB(),
			"baiducloud_peer_conn":                   resourceBaiduCloudPeerConn(),
			"baiducloud_peer_conn_acceptor":          resourceBaiduCloudPeerConnAcceptor(),
			"baiducloud_appblb_server_group":         resourceBaiduCloudAppBlbServerGroup(),
			"baiducloud_appblb_listener":             resourceBaiduCloudAppBlbListener(),
			"baiducloud_bos_bucket":                  resourceBaiduCloudBosBucket(),
			"baiducloud_bos_bucket_object":           resourceBaiduCloudBucketObject(),
			"baiducloud_cert":                        resourceBaiduCloudCert(),
			"baiducloud_cfc_function":                resourceBaiduCloudCFCFunction(),
			"baiducloud_cfc_alias":                   resourceBaiduCloudCFCAlias(),
			"baiducloud_cfc_version":                 resourceBaiduCloudCFCVersion(),
			"baiducloud_cfc_trigger":                 resourceBaiduCloudCFCTrigger(),
			"baiducloud_scs":                         resourceBaiduCloudScs(),
			"baiducloud_cce_cluster":                 resourceBaiduCloudCCECluster(),
			"baiducloud_ccev2_cluster":               resourceBaiduCloudCCEv2Cluster(),
			"baiducloud_ccev2_instance":              resourceBaiduCloudCCEv2Instance(),
			"baiducloud_ccev2_instance_group":        resourceBaiduCloudCCEv2InstanceGroup(),
			"baiducloud_rds_instance":                resourceBaiduCloudRdsInstance(),
			"baiducloud_rds_readonly_instance":       resourceBaiduCloudRdsReadOnlyInstance(),
			"baiducloud_rds_account":                 resourceBaiduCloudRdsAccount(),
			"baiducloud_dts":                         resourceBaiduCloudDts(),
			"baiducloud_iam_user":                    resourceBaiduCloudIamUser(),
			"baiducloud_iam_group":                   resourceBaiduCloudIamGroup(),
			"baiducloud_iam_group_membership":        resourceBaiduCloudIamGroupMembership(),
			"baiducloud_iam_policy":                  resourceBaiduCloudIamPolicy(),
			"baiducloud_iam_user_policy_attachment":  resourceBaiduCloudIamUserPolicyAttachment(),
			"baiducloud_iam_group_policy_attachment": resourceBaiduCloudIamGroupPolicyAttachment(),
			"baiducloud_cdn_domain":                  cdn.ResourceDomain(),
			"baiducloud_cdn_domain_config_cache":     cdn.ResourceDomainConfigCache(),
			"baiducloud_localdns_privatezone":        resourceBaiduCloudLocalDnsPrivateZone(),
			"baiducloud_localdns_vpc":                resourceBaiduCloudLocalDnsVpc(),
			"baiducloud_localdns_record":             resourceBaiduCloudPrivateZoneRecord(),
			"baiducloud_bbc_instance":                resourceBaiduCloudBccInstance(),
			"baiducloud_bbc_image":                   resourceBaiduCloudBbcImage(),
			"baiducloud_deployset":                   resourceBaiduCloudDeploySet(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "The Access Key of BaiduCloud for API operations. You can retrieve this from the 'Security Management' section of the BaiduCloud console.",

		"secret_key": "The Secret key of BaiduCloud for API operations. You can retrieve this from the 'Security Management' section of the BaiduCloud console.",

		"region": "The region where BaiduCloud operations will take place. Examples are bj, su, gz, etc.",

		"assume_role_name": "The role name for assume role.",

		"assume_role_account_id": "The main account id for assume role account.",

		"assume_role_user_id": "The user id for assume role.",

		"assume_role_acl": "The acl for this assume role.",

		"bcc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BCC endpoints.",

		"vpc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom VPC endpoints.",

		"eip_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom EIP endpoints.",

		"appblb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BLB endpoints.",

		"blb_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BLB endpoints.",

		"bos_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BOS endpoints.",

		"cfc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CFC endpoints.",

		"scs_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom SCS endpoints.",

		"cce_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CCE endpoints.",

		"ccev2_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CCEv2 endpoints.",

		"rds_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RDS endpoints.",

		"dts_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DTS endpoints.",

		"cdn_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CDN endpoints.",

		"bbc_endpoint": "Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BBC endpoints.",
	}
}

func endpointsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"bcc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["bcc_endpoint"],
				},
				"vpc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["vpc_endpoint"],
				},
				"eip": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["eip_endpoint"],
				},
				"appblb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["appblb_endpoint"],
				},
				"blb": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["blb_endpoint"],
				},
				"bos": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["bos_endpoint"],
				},
				"cfc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cfc_endpoint"],
				},
				"scs": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["scs_endpoint"],
				},
				"cce": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cce_endpoint"],
				},
				"ccev2": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["ccev2_endpoint"],
				},
				"rds": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["rds_endpoint"],
				},
				"dts": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["dts_endpoint"],
				},
				"cdn": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["cdn_endpoint"],
				},
				"bbc": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "",
					Description: descriptions["bbc_endpoint"],
				},
			},
		},
		Set: endpointsToHash,
	}
}

func endpointsToHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["bcc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["vpc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["eip"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["appblb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["blb"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["bos"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cfc"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["scs"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cce"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ccev2"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["rds"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["dts"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["cdn"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["bbc"].(string)))
	return hashcode.String(buf.String())
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	accessKey, ok := d.GetOk("access_key")
	if !ok {
		accessKey = os.Getenv(PROVIDER_ACCESS_KEY)
	}
	secretKey, ok := d.GetOk("secret_key")
	if !ok {
		secretKey = os.Getenv(PROVIDER_SECRET_KEY)
	}
	region, ok := d.GetOk("region")
	if !ok {
		region = os.Getenv(PROVIDER_REGION)
	}

	config := connectivity.Config{
		AccessKey: accessKey.(string),
		SecretKey: secretKey.(string),
		Region:    connectivity.Region(region.(string)),
	}

	assumeRoleList, ok := d.GetOk("assume_role")
	if ok {
		if assumeRoles, ok := assumeRoleList.([]interface{}); ok && len(assumeRoles) > 0 {
			assumeRole := assumeRoles[0].(map[string]interface{})

			if accountId, ok := assumeRole["account_id"]; ok {
				config.AssumeRoleAccountId = accountId.(string)
			}

			if roleName, ok := assumeRole["role_name"]; ok {
				config.AssumeRoleRoleName = roleName.(string)
			}

			if userId, ok := assumeRole["user_id"]; ok {
				config.AssumeRoleUserId = userId.(string)
			}

			if acl, ok := assumeRole["acl"]; ok {
				config.AssumeRoleAcl = acl.(string)
			}
		}
	}

	endpointsSet := d.Get("endpoints").(*schema.Set)

	for _, endpointsSetI := range endpointsSet.List() {
		endpoints := endpointsSetI.(map[string]interface{})
		config.ConfigEndpoints = make(map[connectivity.ServiceCode]string)
		config.ConfigEndpoints[connectivity.BCCCode] = strings.TrimSpace(endpoints["bcc"].(string))
		config.ConfigEndpoints[connectivity.VPCCode] = strings.TrimSpace(endpoints["vpc"].(string))
		config.ConfigEndpoints[connectivity.EIPCode] = strings.TrimSpace(endpoints["eip"].(string))
		config.ConfigEndpoints[connectivity.APPBLBCode] = strings.TrimSpace(endpoints["appblb"].(string))
		config.ConfigEndpoints[connectivity.BLBCode] = strings.TrimSpace(endpoints["blb"].(string))
		config.ConfigEndpoints[connectivity.BOSCode] = strings.TrimSpace(endpoints["bos"].(string))
		config.ConfigEndpoints[connectivity.BOSCode] = strings.TrimSpace(endpoints["cfc"].(string))
		config.ConfigEndpoints[connectivity.SCSCode] = strings.TrimSpace(endpoints["scs"].(string))
		config.ConfigEndpoints[connectivity.CCECode] = strings.TrimSpace(endpoints["cce"].(string))
		config.ConfigEndpoints[connectivity.CCEv2Code] = strings.TrimSpace(endpoints["ccev2"].(string))
		config.ConfigEndpoints[connectivity.RDSCode] = strings.TrimSpace(endpoints["rds"].(string))
		config.ConfigEndpoints[connectivity.DTSCode] = strings.TrimSpace(endpoints["dts"].(string))
		config.ConfigEndpoints[connectivity.CDNCode] = strings.TrimSpace(endpoints["cdn"].(string))
		config.ConfigEndpoints[connectivity.BBCCode] = strings.TrimSpace(endpoints["bbc"].(string))
	}

	client, err := config.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Assume role configurations, for more information, please refer to https://cloud.baidu.com/doc/IAM/s/Qjwvyc8ov",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: descriptions["assume_role_name"],
				},

				"account_id": {
					Type:        schema.TypeString,
					Required:    true,
					Description: descriptions["assume_role_account_id"],
				},

				"user_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_user_id"],
				},

				"acl": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: descriptions["assume_role_acl"],
				},
			},
		},
	}
}
