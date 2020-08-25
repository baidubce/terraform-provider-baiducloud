terraform {
  required_providers {
    baiducloud = {
      versions = ["1.0.1"]
      source = "example.com/myorg/baiducloud"
    }
  }
}

provider "baiducloud" {

}

data "baiducloud_ccev2_container_cidr" "default" {
  vpc_id = baiducloud_vpc.default.id
  vpc_cidr = baiducloud_vpc.default.cidr
  cluster_max_node_num = 16
  max_pods_per_node = 32
  private_net_cidrs = ["192.168.0.0/16",]
  k8s_version = "1.16.8"
  ip_version = "ipv4"
  output_file = "${path.cwd}/recommendContainerCidr.txt"
}

data "baiducloud_ccev2_clusterip_cidr" "default" {
  vpc_cidr = baiducloud_vpc.default.cidr
  container_cidr = var.container_cidr
  cluster_max_service_num = 32
  private_net_cidrs = ["192.168.0.0/16",]
  ip_version = "ipv4"
  output_file = "${path.cwd}/recommendClusterIPCidr.txt"
}

//Create a vpc for cluster
resource "baiducloud_vpc" "default" {
  name        = "test-vpc-tf-auto"
  description = "test-vpc-tf-auto"
  cidr        = "192.168.0.0/16"
}

//Get available zone name
data "baiducloud_zones" "defaultA" {
  name_regex = ".*a$"
}

data "baiducloud_zones" "defaultB" {
  name_regex = ".*b$"
}

data "baiducloud_zones" "defaultC" {
  name_regex = ".*c$"
}

//Create subnet in different available zone
resource "baiducloud_subnet" "defaultA" {
  name        = "test-subnet-tf-auto-1"
  zone_name   = data.baiducloud_zones.defaultA.zones.0.zone_name
  cidr        = "192.168.1.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "test-subnet-tf-auto-1"
  tags = {
    "testKey"  = "testValue"
    "testKey2" = "testValue2"
  }
}

resource "baiducloud_subnet" "defaultB" {
  name        = "test-subnet-tf-auto-2"
  zone_name   = data.baiducloud_zones.defaultB.zones.0.zone_name
  cidr        = "192.168.2.0/24"
  vpc_id      = baiducloud_vpc.default.id
  description = "test-subnet-tf-auto-2"

  tags = {
    "testKey"  = "testValue"
    "testKey2" = "testValue2"
  }
}

resource "baiducloud_subnet" "defaultC" {
  name = "test-subnet-tf-auto-3"
  zone_name = data.baiducloud_zones.defaultC.zones.0.zone_name
  cidr = "192.168.3.0/24"
  vpc_id = baiducloud_vpc.default.id
  description = "test-subnet-tf-auto-3"

  tags = {
    "testKey" = "testValue"
    "testKey2" = "testValue2"
  }
}

//Crate security group and rules
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

//Create cluster with masters in different availavle zone.
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
//      cluster_ha = 1
//      exposed_public = false
      cluster_blb_vpc_subnet_id = baiducloud_subnet.defaultA.id
//      managed_cluster_master_option {
//        master_vpc_subnet_zone = "zoneA"
//      }
    }
    container_network_config  {
      mode = "kubenet"
      lb_service_vpc_subnet_id = baiducloud_subnet.defaultA.id
      node_port_range_min = 30000
      node_port_range_max = 32767
      max_pods_per_node = 64
      cluster_pod_cidr = var.cluster_pod_cidr
      cluster_ip_service_cidr = var.cluster_ip_service_cidr
    }
    cluster_delete_option {
      delete_resource = true
      delete_cds_snapshot = true
    }
  }
  //If you with to create different master nodes in different available zone
  //  please append master_spces with different vpc_subnet_id and available_zone.
  master_specs {
    count = 1
    master_spec {
      cce_instance_id = ""
      instance_name = "ccev2_test_instance_master-1"
      cluster_role = "master"
      existed = false
      machine_type = "BCC"
      instance_type = "N3"
      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.defaultA.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneA"
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
  master_specs {
    count = 1
    master_spec {
      cce_instance_id = ""
      instance_name = "ccev2_test_instance_master-2"
      cluster_role = "master"
      existed = false
      machine_type = "BCC"
      instance_type = "N3"
      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.defaultB.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneB"
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
  master_specs {
    count = 1
    master_spec {
      cce_instance_id = ""
      instance_name = "ccev2_test_instance_master-3"
      cluster_role = "master"
      existed = false
      machine_type = "BCC"
      instance_type = "N3"
      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.defaultC.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneC"
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

data "baiducloud_images" "default" {
  image_type = "System"
  name_regex = "7.5.*"
  os_name    = "CentOS"
}

//If you wish to create more instances in different available zone
//  please append more instance group with different vpc_subnet_id and available_zone.
resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_1" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default_custom.id
    replicas = var.instance_group_replica_1
    instance_group_name = "ccev2_instance_group_1"
    instance_template {
      cce_instance_id = ""
      instance_name = "ccev2_test_instance_1"
      cluster_role = "node"
      existed = false
      machine_type = "BCC"
      instance_type = "N3"
      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.defaultA.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneA"
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

resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_2" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default_custom.id
    replicas = var.instance_group_replica_1
    instance_group_name = "ccev2_instance_group_2"
    instance_template {
      cce_instance_id = ""
      instance_name = "ccev2_test_instance_2"
      cluster_role = "node"
      existed = false
      machine_type = "BCC"
      instance_type = "N3"
      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.defaultB.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneB"
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

resource "baiducloud_ccev2_instance_group" "ccev2_instance_group_3" {
  spec {
    cluster_id = baiducloud_ccev2_cluster.default_custom.id
    replicas = var.instance_group_replica_1
    instance_group_name = "ccev2_instance_group_3"
    instance_template {
      cce_instance_id = ""
      instance_name = "ccev2_test_instance_3"
      cluster_role = "node"
      existed = false
      machine_type = "BCC"
      instance_type = "N3"
      vpc_config {
        vpc_id = baiducloud_vpc.default.id
        vpc_subnet_id = baiducloud_subnet.defaultC.id
        security_group_id = baiducloud_security_group.default.id
        available_zone = "zoneC"
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




