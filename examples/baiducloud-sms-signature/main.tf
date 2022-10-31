provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}


resource "baiducloud_sms_signature" "default" {
  content      = "baidu"
  content_type = "Enterprise"
  description  = "terraform test"
  country_type = "DOMESTIC"

}

data "baiducloud_sms_signature" "default" {
	signature_id = "${baiducloud_sms_signature.default.id}"
}