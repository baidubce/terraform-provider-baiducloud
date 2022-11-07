provider "baiducloud" {
  # option config, you can use assume role as the operation account
  #assume_role {
  #  account_id = "your-account-id"
  #  role_name = "your-role-name"
  #}
}


resource "baiducloud_sms_template" "default" {
  name	         = "My test template"
  content        = "Test content"
  sms_type       = "CommonNotice"
  country_type   = "GLOBAL"
  description    = "this is a test sms template"
}

data "baiducloud_sms_template" "default" {
	template_id = "${baiducloud_sms_template.default.id}"
}