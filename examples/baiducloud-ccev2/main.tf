provider "baiducloud" {
  # option config, you can use assume role as the operation account
  # assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  # }
}

# To Recommend Cluster IP CIDR or Container CIDR
data "baiducloud_ccev2_container_cidr" "default" {
  vpc_id = baiducloud_vpc.default.id
  vpc_cidr = baiducloud_vpc.default.cidr
  cluster_max_node_num = 16
  max_pods_per_node = 32
  private_net_cidrs = ["192.168.0.0/16",]
  k8s_version = "1.16.8"
  ip_version = "ipv4"
  #   output_file = "${path.cwd}/recommendContainerCidr.txt"
}

data "baiducloud_ccev2_clusterip_cidr" "default" {
  vpc_cidr = baiducloud_vpc.default.cidr
  container_cidr = var.container_cidr
  cluster_max_service_num = 32
  private_net_cidrs = ["192.168.0.0/16",]
  ip_version = "ipv4"
  # output_file = "${path.cwd}/recommendClusterIPCidr.txt"
}

# ====To delete a specific instance of a instance group====
# A instance with a lower "cce_instance_priority" will be deleted firstly when you shrink a instance group.
# The default value of "cce_instance_priority" is 5.
# Step 1:
# If you wish delete a specific instance of a instance group, set the "cce_instance_priority" to a lower value, for example, 1.
# Tips: Apply a new "baiducloud_ccev2_instance" resource will not create a new instance.
#       It is just a bind to a existed remote instance.
#       If you want to create more instances, please use resource "baiducloud_ccev2_instance_group"
resource "baiducloud_ccev2_instance" "default" {
  cluster_id        = baiducloud_ccev2_cluster.default_custom.id
  instance_id       = data.baiducloud_ccev2_instance_group_instances.default.instance_list.0.instance_spec.0.cce_instance_id
  spec {
    cce_instance_priority = 1 # Set this value lower than the default value, fox example, 1.
  }
}
# Step 2: type "terraform apply" to update "cce_instance_priority" value of the instance.

# YOU MUST APPLY THE CHANGE IN STEP_2 BEFORE STEP_3
# Step 3: set the value "replicas" of the instance group from "N" to "N-1" and then apply the change.
variable "instance_group_replica_3" {
  # default 4 => 3
  default = 3
}
#===========================================================

# Steps to create an CCEv2 cluster:
# For custom cluster,  follow steps 1, 2, 3, 4, 5.1, 6
# For managed cluster, follow steps 1, 2, 3, 4, 5.2, 6

# 1.Create a vpc for the cluster
resource "baiducloud_vpc" "default" {
  name        = "test-vpc-tf-auto"
  description = "test-vpc-tf-auto"
  cidr        = "192.168.0.0/16"
}

# 2.Get available zone name
data "baiducloud_zones" "defaultA" {
  name_regex = ".*a$"
}

data "baiducloud_zones" "defaultB" {
  name_regex = ".*b$"
}

data "baiducloud_zones" "defaultC" {
  name_regex = ".*c$"
}

# 3.Create subnet in different available zone
resource "baiducloud_subnet" "defaultA" {
  name        = "test-subnet-tf-auto-1"
  zone_name   = data.baiducloud_zones.defaultA.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "test-subnet-tf-auto-1"
}

resource "baiducloud_subnet" "defaultB" {
  name        = "test-subnet-tf-auto-2"
  zone_name   = data.baiducloud_zones.defaultB.zones.0.zone_name
  cidr        = "192.168.2.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "test-subnet-tf-auto-2"
}

resource "baiducloud_subnet" "defaultC" {
  name = "test-subnet-tf-auto-3"
  zone_name = data.baiducloud_zones.defaultC.zones.0.zone_name
  cidr = "192.168.3.0/24"
  vpc_id = baiducloud_vpc.default.id
  description = "test-subnet-tf-auto-3"
}

# 4.Create security group and security group rules
resource "baiducloud_security_group" "default" {
  name   = "test-security-group-tf-auto"
  vpc_id = baiducloud_vpc.default.id
}

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

# 5.Create a cluster
data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}

# 5.1 Create a user custom cluster with masters spread in different AvailableZone.
resource "baiducloud_ccev2_cluster" "default_custom" {
  cluster_spec  {
    cluster_name = var.cluster_name
    cluster_type = "normal"
    k8s_version = "1.16.8"
    runtime_type = "docker"
    vpc_id = baiducloud_vpc.default.id
    plugins = ["core-dns", "kube-proxy"]
    master_config {
      master_type = "custom"
      exposed_public = false
      cluster_blb_vpc_subnet_id = baiducloud_subnet.defaultA.id
    }
    container_network_config  {
      mode = "kubenet"
      lb_service_vpc_subnet_id = baiducloud_subnet.defaultA.id
      node_port_range_min = 30000
      node_port_range_max = 32767
      max_pods_per_node = 64
      cluster_pod_cidr = var.cluster_pod_cidr
      cluster_ip_service_cidr = var.cluster_ip_service_cidr
      ip_version = "ipv4"
      kube_proxy_mode = "iptables"
    }
    cluster_delete_option {
      delete_resource = true
      delete_cds_snapshot = true
    }
  }
  # Tips: If you wish to create master nodes spread in different available zone
  #       please add [master_spces] with different [master_spec.vpc_config.vpc_subnet_id] and
  #       [master_spec.vpc_config.available_zone].
  #       Notes that the total number of masters can only be 1,3 or 5
  dynamic "master_specs" {
    for_each = [
      {
        master_name: "tf_instance_1",
        subnet_id: baiducloud_subnet.defaultA.id,
        zone_name: "zoneA"
        count: 1
      },
      {
        master_name: "tf_instance_2",
        subnet_id: baiducloud_subnet.defaultB.id,
        zone_name: "zoneB"
        count: 0
      },
      {
        master_name: "tf_instance_3",
        subnet_id: baiducloud_subnet.defaultC.id,
        zone_name: "zoneC"
        count: 0
      },
    ]
    content {
      count = master_specs.value.count
      master_spec {
        cce_instance_id = ""
        instance_name = master_specs.value.master_name
        cluster_role = "master"
        existed = false
        instance_type = "N3"
        vpc_config {
          vpc_id = baiducloud_vpc.default.id
          vpc_subnet_id = master_specs.value.subnet_id
          security_group_id = baiducloud_security_group.default.id
          available_zone = master_specs.value.zone_name
        }
        instance_resource {
          cpu = 4
          mem = 8
          root_disk_size = 40
          local_disk_size = 0
        }
        image_id = data.baiducloud_images.default.images.0.id
        instance_os {
          image_type = "System"
        }
        need_eip = false
        admin_password = "test123!YT"
        ssh_key_id = ""
        instance_charging_type = "Postpaid"
        runtime_type = "docker"
      }
    }
  }
}

# 5.2 Create a managed cluster. Please uncomment the following code and comment the code in step-5.1 before you apply.
# resource "baiducloud_ccev2_cluster" "default_managed" {
#   cluster_spec  {
#    cluster_name = var.cluster_name
#    cluster_type = "normal"
#    k8s_version = "1.16.8"
#    runtime_type = "docker"
#    vpc_id = baiducloud_vpc.default.id
#    plugins = ["core-dns", "kube-proxy"]
#   master_config {
      #     master_type = "managed"
      #     cluster_ha = 1
      #     exposed_public = false
      #     cluster_blb_vpc_subnet_id = baiducloud_subnet.defaultA.id
      #     managed_cluster_master_option {
        #        master_vpc_subnet_zone = "zoneA"
        #      }
      #    }
#    container_network_config  {
      #       mode = "kubenet"
      #       lb_service_vpc_subnet_id = baiducloud_subnet.defaultA.id
      #      node_port_range_min = 30000
      #     node_port_range_max = 32767
      #     max_pods_per_node = 64
      #     cluster_pod_cidr = var.cluster_pod_cidr
      #     cluster_ip_service_cidr = var.cluster_ip_service_cidr
      #      ip_version = "ipv4"
      #     kube_proxy_mode = "iptables"
      #    }
#    cluster_delete_option {
      #     delete_resource = true
      #       delete_cds_snapshot = true
      #     }
#  }
# }

# 6. Add worker nodes to cluster by using instance group.
# Tips: If you wish to create more worker nodes spread in different available zone
#       please add more [baiducloud_ccev2_instance_group] with different
#       spec.instance_template.vpc_config.vpc_subnet_id and spec.instance_template.vpc_config.available_zone.
resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_1" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default_custom.id
    replicas = var.instance_group_replica_1
    instance_group_name = "ig_1"
    instance_template {
      cce_instance_id = ""
      instance_name = "tf_ins_ig_1"
      cluster_role = "node"
      existed = false
      instance_type = "N3"

      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.defaultA.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneA"
      }
      deploy_custom_config {
        pre_user_script  = "ls"
        post_user_script = "date"
      }
      instance_resource {
        cpu = 1
        mem = 4
        root_disk_size = 40
        local_disk_size = 0
      }
      image_id = data.baiducloud_images.default.images.0.id
      instance_os {
        image_type = "System"
      }
      need_eip = false
      admin_password = "test123!YT"
      ssh_key_id = ""
      instance_charging_type = "Postpaid"
      runtime_type = "docker"
    }
  }
}

resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_2" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default_custom.id
    replicas = var.instance_group_replica_2
    instance_group_name = "ig_2"
    instance_template {
      cce_instance_id = ""
      instance_name = "tf_ins_ig_2"
      cluster_role = "node"
      existed = false
      instance_type = "N3"
      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.defaultB.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneB"
      }
      instance_resource {
        cpu = 1
        mem = 4
        root_disk_size = 40
        local_disk_size = 0
      }
      image_id = data.baiducloud_images.default.images.0.id
      instance_os {
        image_type = "System"
      }
      need_eip = false
      admin_password = "test123!YT"
      ssh_key_id = ""
      instance_charging_type = "Postpaid"
      runtime_type = "docker"
    }
  }
}

resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_3" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default_custom.id
    replicas = var.instance_group_replica_3
    instance_group_name = "ig_3"
    instance_template {
      cce_instance_id = ""
      instance_name = "tf_ins_ig_3"
      cluster_role = "node"
      existed = false
      instance_type = "N3"
      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.defaultC.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneC"
      }
      instance_resource {
        cpu = 1
        mem = 4
        root_disk_size = 40
        local_disk_size = 0
      }
      image_id = data.baiducloud_images.default.images.0.id
      instance_os {
        image_type = "System"
      }
      need_eip = false
      admin_password = "test123!YT"
      ssh_key_id = ""
      instance_charging_type = "Postpaid"
      runtime_type = "docker"
    }
  }
}

# Get the master-node list and follower-node list of the cluster
# Tips: If the type of cluster is "managed", "master_list" will be empty
data "baiducloud_ccev2_cluster_instances" "default" {
  cluster_id = baiducloud_ccev2_cluster.default_custom.id
}

# Get the instance list of the instance group
data "baiducloud_ccev2_instance_group_instances" "default" {
  cluster_id = baiducloud_ccev2_cluster.default_custom.id
  instance_group_id = baiducloud_ccev2_instance_group.ccev2_instance_group_3.id
}




