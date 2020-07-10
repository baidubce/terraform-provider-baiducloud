provider "baiducloud" {}

data "baiducloud_zones" "defaultA" {
  name_regex = ".*a$"
}

data "baiducloud_zones" "defaultB" {
  name_regex = ".*b$"
}

data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}

resource "baiducloud_vpc" "default" {
  name        = var.vpc-name
  description = var.description
  cidr        = "192.168.0.0/16"
}

resource "baiducloud_subnet" "defaultA" {
  name        = var.subnet_name_a
  zone_name   = data.baiducloud_zones.defaultA.zones.0.zone_name
  cidr        = var.subnet_cidr_a
  vpc_id      = baiducloud_vpc.default.id
  description = var.description

  tags = {
    "testKey"  = "testValue"
    "testKey2" = "testValue2"
  }
}

resource "baiducloud_subnet" "defaultB" {
  name        = var.subnet_name_b
  zone_name   = data.baiducloud_zones.defaultB.zones.0.zone_name
  cidr        = var.subnet_cidr_b
  vpc_id      = baiducloud_vpc.default.id
  description = var.description

  tags = {
    "testKey"  = "testValue"
    "testKey2" = "testValue2"
  }
}

resource "baiducloud_security_group" "default" {
  name   = var.scurity-group-name
  vpc_id = baiducloud_vpc.default.id
}

# for more detail CCE Cluster security group rule config
# please refer to https://cloud.baidu.com/doc/CCE/s/Fjwvy1cid
resource "baiducloud_security_group_rule" "default" {
  security_group_id = baiducloud_security_group.default.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "ingress"
}

resource "baiducloud_security_group_rule" "default2" {
  security_group_id = baiducloud_security_group.default.id
  remark            = "remark"
  protocol          = "all"
  port_range        = "1-65535"
  direction         = "egress"
}

data "baiducloud_cce_versions" "default" {
  version_regex = ".*13.*"
}

data "baiducloud_cce_container_net" "default" {
  vpc_id   = baiducloud_vpc.default.id
  vpc_cidr = baiducloud_vpc.default.cidr
}


# example for managed cluster
resource "baiducloud_cce_cluster" "default_managed" {
  cluster_name = var.cce-name

  # available zone value support zoneA/zoneB/...
  # For example, if use baiducloud_zones to get all avalilable zone in beijing region
  # and you can get [cn-bj-a, cn-bj-b, cn-bj-c, cn-bj-d, cn-bj-e],
  # then you can set available zone in value zoneA/zoneB/zoneC/zoneD/zoneE
  main_available_zone = "zoneA"
  version             = data.baiducloud_cce_versions.default.versions.0
  container_net       = "172.18.0.0/16"

  # optional parameters
  advanced_options = {
    # if you set advanced_options, all parameters below are required
    kube_proxy_mode = "ipvs"
    dns_mode        = "CoreDNS"
    cni_mode        = "cni"
    cni_type        = "VPC_SECONDARY_IP_VETH"
    max_pod_num     = "256"
  }
  delete_eip_cds   = "true"
  delete_snapshots = "true"

  worker_config {
    # NOTICE: only count and subnet_uuid map are variable parameter
    # count is map, key is available zone, value is worker nodes count in this zone
    count = {
      "zoneA" : 2
      #"zoneB" : 3
    }
    subnet_uuid = {
      "zoneA" : baiducloud_subnet.defaultA.id
      "zoneB" : baiducloud_subnet.defaultB.id
    }

    # for more detail config, please refer to https://cloud.baidu.com/doc/CCE/s/Ujwvy1fxs#%E5%88%9B%E5%BB%BA%E9%9B%86%E7%BE%A4
    instance_type     = "10"
    cpu               = 1
    memory            = 2
    security_group_id = baiducloud_security_group.default.id
    product_type      = "postpay"
    image_id          = data.baiducloud_images.default.images.0.id
    image_type        = "common"

    # optional param
    root_disk_size_in_gb   = 100
    root_disk_storage_type = "ssd"
    admin_pass             = "baiduPasswd@123"

    # please set your keypair id
    #keypair_id             = "k-xxxxxxx"

    cds_disks {
      volume_type     = "sata"
      disk_size_in_gb = 10
    }

    eip = {
      bandwidth_in_mbps = 100

      # support netraffic/bandwidth
      sub_product_type = "netraffic"
    }
  }
}

# example for independent cluster
#resource "baiducloud_cce_cluster" "default_independent" {
# cluster_name        = var.cce-name
#  main_available_zone = "zoneA"
#  version             = data.baiducloud_cce_versions.default.versions.0
#  container_net       = "172.16.0.0/16"
#
#  advanced_options = {
#    kube_proxy_mode = "ipvs"
#    dns_mode        = "CoreDNS"
#    cni_mode        = "cni"
#    cni_type        = "VPC_SECONDARY_IP_VETH"
#    max_pod_num     = "256"
#  }
#
#  delete_eip_cds   = "true"
#  delete_snapshots = "true"
#
#
#  worker_config {
#    count = {
#      "zoneA" : 2
#      #"zoneB" : 3
#    }
#
#    instance_type = "10"
#    cpu           = 1
#    memory        = 2
#    subnet_uuid = {
#      "zoneA" : baiducloud_subnet.defaultA.id
#      "zoneB" : baiducloud_subnet.defaultB.id
#    }
#    security_group_id = baiducloud_security_group.defualt.id
#    product_type      = "postpay"
#    image_id          = data.baiducloud_images.default.images.0.id
#    image_type        = "common"
#
#    # optional param
#    root_disk_size_in_gb   = 100
#    root_disk_storage_type = "ssd"
#    admin_pass             = "baiduPasswd@123"
#
#    # please set your keypair id
#    #keypair_id            = "k-xxxxxxx"
#
#    cds_disks {
#      volume_type     = "sata"
#      disk_size_in_gb = 10
#    }
#
#    eip = {
#      bandwidth_in_mbps = 100
#      sub_product_type  = "netraffic"
#    }
#  }
#
#  master_config {
#    instance_type     = "10"
#    cpu               = 4
#    memory            = 8
#    image_type        = "common"
#    logical_zone      = "zoneA"
#    subnet_uuid       = baiducloud_subnet.defaultA.id
#    security_group_id = baiducloud_security_group.defualt.id
#    product_type      = "postpay"
#    image_id          = data.baiducloud_images.default.images.0.id
#    # please set your keypair id
#    #keypair_id       = "k-xxxxxxx"
#  }
#}

# get CCE Cluster kubectl.conf
data "baiducloud_cce_kubeconfig" "default" {
  cluster_uuid = baiducloud_cce_cluster.default_managed.id
  output_file = "${path.cwd}/kubectl.conf"
}

data "baiducloud_cce_cluster_nodes" "default" {
  cluster_uuid = baiducloud_cce_cluster.default_managed.id
}