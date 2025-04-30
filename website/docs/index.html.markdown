---
layout: "baiducloud"
page_title: "Provider: baiducloud"
sidebar_current: "docs-baiducloud-index"
description: |-
  The BaiduCloud provider is used to interact with many resources supported by BaiduCloud. The provider needs to be configured with proper credentials before being used.
---

# BaiduCloud Provider

The BaiduCloud provider is used to interact with
many resources supported by [BaiduCloud](https://cloud.baidu.com). The provider needs to be configured
with the proper credentials before being used.

Use the navigation on the left to read about available resources.

## Example Usage

```hcl
# Configure the BaiduCloud Provider
provider "baiducloud" {
  access_key  = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}

# Create a web server
resource "baiducloud_instance" "my-server"{
  image_id = "m-DpgNg8lO"
  name = "from-terraform"
  availability_zone = "cn-bj-c"
}

# Create a security group
resource "baiducloud_security_group" "my-sg"{
  name = "from-terraform"
  description = "baiducloud security group created by terraform"
}

# Create an eip
resource "baiducloud_eip" "my-eip"{
  name        = "from-terraform"
  bandwidth_in_mbps = 100
  payment_timing = "Postpaid"
  billing_method = "ByTraffic"
}

# Create a VPC
resource "baiducloud_vpc" "default" {
    name = "my-vpc"
    description = "baiducloud vpc created by terraform"
	cidr = "192.168.0.0/24"
}

# Create a VPC subnet
resource "baiducloud_subnet" "default" {
  name = "my-subnet"
  zone_name = "cn-bj-a"
  cidr = "192.168.3.0/24"
  vpc_id = "${baiducloud_vpc.default.id}"
}

resource "baiducloud_appblb" "my-appblb" {
  name        = "${var.name}"
  description = "${var.description}"
  vpc_id      = "${var.vpc_id}"
  subnet_id   = "${var.subnet_id}"
}

resource "baiducloud_appblb_server_group" "my-appblb-sg" {
  name        = "${var.name}"
  description = "${var.description}"
  blb_id      = "${var.blb_id}"
}
```

## Authentication

The BaiduCloud provider offers a flexible means of providing credentials for authentication.
The following methods are supported, and explained below in this order:

- Static credentials
- Environment variables
- AssumeRole credentials

### Static credentials

Static credentials can be provided by adding `access_key` `secret_key` and `region` in-line in the
baiducloud provider block:

Usage:

```hcl
provider "baiducloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"
}
```

### Environment variables

You can provide your credentials via `BAIDUCLOUD_ACCESS_KEY` and `BAIDUCLOUD_SECRET_KEY`,
environment variables, representing your BaiduCloud Access Key and Secret Key, respectively.
`BAIDUCLOUD_REGION` is also used, if applicable:

```hcl
provider "baiducloud" {}
```

Usage:

```shell
$ export BAIDUCLOUD_ACCESS_KEY="your_fancy_accesskey"
$ export BAIDUCLOUD_SECRET_KEY="your_fancy_secretkey"
$ export BAIDUCLOUD_REGION="bj"
$ terraform plan
```

### AssumeRole credentials

You can use `assume_role` as your credential role:

Usage:

```hcl
provider "baiducloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "${var.region}"

  assume_role {
    account_id = "your-account-id"
    role_name = "your-role-name"
  }
}
```

## Endpoints

Endpoints can be provided by adding an `endpoints` in-line in the baiducloud provider block:

Usage:

```hcl
provider "baiducloud" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  endpoints [
      bcc = "your_fancy_bcc_custom_endpoint"
      vpc = "your_fancy_vpc_custom_endpoint"
      eip = "your_fancy_eip_custom_endpoint"
      appblb = "your_fancy_blb_custom_endpoint"
    ]
}
```

## Argument Reference

The following arguments are supported:

* `access_key` - (Optional) This is the BaiduCloud access key. It must be provided, but
  it can also be sourced from the `BAIDUCLOUD_ACCESS_KEY` environment variable.

* `secret_key` - (Optional) This is the BaiduCloud secret key. It must be provided, but
  it can also be sourced from the `BAIDUCLOUD_SECRET_KEY` environment variable.

* `session_token` - (Optional) This is the BaiduCloud session token. It must be provided when 
   using a temporary access key, it can also be sourced from the `BAIDUCLOUD_SESSION_TOKEN` environment variable.

* `region` - (Required) This is the BaiduCloud region. It must be provided, but
  it can also be sourced from the `BAIDUCLOUD_REGION` environment variables.
  The default input value is bj. Available value is [bj, bd, gz, su, fsh, fwh, hkg, sin]

* `endpoints` - (Optional) An `endpoints` block (documented below) to support custom endpoints.

* `assume_role` - (Optional) An `assume_role` block (documented below) to support assume role credentials. Assume role configurations, for more information, please refer to [STS Service](https://cloud.baidu.com/doc/IAM/s/Qjwvyc8ov).

Nested `endpoints` block supports the following:

* `bcc` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BCC endpoints.

* `vpc` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom VPC endpoints.

* `eip` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom EIP endpoints.

* `appblb` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BLB endpoints.

* `bos` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom BOS endpoints.

* `cfc` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CFC endpoints.

* `scs` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom SCS endpoints.

* `cce` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CCE endpoints.

* `ccev2` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom CCEv2 endpoints.

* `rds` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom RDS endpoints.

* `dts` - (Optional) Use this to override the default endpoint URL constructed from the `region`. It's typically used to connect to custom DTS endpoints.

Nested `assume_role` block supports the following:

* `role_name` - (Required) The role name for assume role.

* `account_id` - (Required) The main account id for assume role account.

* `user_id` - (Optional) The user id for assume role.

* `acl` - (Optional) The acl for this assume role.


## Testing

Credentials must be provided via the `BAIDUCLOUD_ACCESS_KEY`, and `BAIDUCLOUD_SECRET_KEY` environment variables in order to run acceptance tests.
