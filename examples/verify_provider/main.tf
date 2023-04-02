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
