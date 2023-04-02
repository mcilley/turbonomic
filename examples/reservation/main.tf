// -----------------------------------------------------------------------------
// On Prem Lab / Dev Instance
// -----------------------------------------------------------------------------

// If you have not set the provider environment variables, uncomment
// these lines:
//   TURBO_CLIENT_USERNAME
//   TURBO_CLIENT_PASSWORD
//   TURBO_SERVER_HOSTNAME
//variable "client_username" {}
//variable "client_password" {}
//variable "server_hostname" {}

// -----------------------------------------------------------------------------

provider "turbonomic" {

  // If you have not set the provider environment variables, uncomment
  // these lines:
  //   TURBO_CLIENT_USERNAME
  //   TURBO_CLIENT_PASSWORD
  //   TURBO_SERVER_HOSTNAME
  //client_username = "${var.client_username}"
  //client_password = "${var.client_password}"
  //server_hostname = "${var.server_hostname}"

  client_tls_insecure = "true"

  server_protocol = "https"

}

// -----------------------------------------------------------------------------



/*
resource "turbonomic_reservation" "reservation" {
  count                    = "${var.vsphere_vm_count}"
  action                   = "RESERVATION"
  constraint_ids           = ["${data.turbonomic_market_policy.policy.id}"]
  deployment_profile_id    = "${data.turbonomic_deployment_profile.profile.id}"
  entity_name              = "${format("%s%d", var.vsphere_vm_name, count.index+1)}"
  template_id              = "${turbonomic_template.template.id}"
  reservation_reserve_time = "${timeadd( local.reservation_base_time, format("%dm", count.index + 5))}"
  reservation_expire_time  = "${timeadd( local.reservation_base_time, format("%dm", var.vsphere_vm_count + 5))}"
  depends_on               = ["turbonomic_template.template"]

  lifecycle {
    ignore_changes = [
      "reservation_reserve_time",
      "reservation_expire_time",
    ]
  }
}
*/
