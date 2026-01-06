terraform {
  required_providers {
    dsfhub = {
      source = "imperva/dsfhub"
    }
  }
}

variable "dsfhub_host" {}
variable "dsfhub_token" {}

provider "dsfhub" {
  dsfhub_host  = var.dsfhub_host
  dsfhub_token = var.dsfhub_token
}

# Run #1 with is_load_balancer = "true"
resource "dsfhub_ciphertrust" "this" {
  cm_name                  = "test-cm-tf"
  hostname                 = "1.2.3.4"
  port                     = 443
  is_load_balancer         = "true"
  auth_method              = "password"
  username                 = "myuser"
  password                 = "mypassword"
  ddc_enabled              = true
  ddc_active_node_hostname = "ddc-active-node-hostname"
  ddc_active_node_port     = 443
}

# Run #2 with is_load_balancer = "false"
# resource "dsfhub_ciphertrust" "this" {
#   cm_name                  = "test-cm-tf"
#   hostname                 = "1.2.3.4"
#   port                     = 443
#   is_load_balancer         = "false"
#   auth_method              = "password"
#   username                 = "myuser"
#   password                 = "mypassword"
# }
