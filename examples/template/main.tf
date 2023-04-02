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

data "turbonomic_template" "example" {
  display_name   = "MAX:VMs_MA2\\c05.esx.ma2"
  vcenter_server = "vcenter.host.foo.foo.com"
}

data "turbonomic_deployment_profile" "example" {
  display_name = "DEP-PCKR20181023141134_CENTOS751804_BO1_4.0"
}

/*
resource "turbonomic_template" "example" {

  class_name   = "VirtualMachine"
  display_name = "terraform_test_template"
  description  = "Test VM Template created with Terraform"

  //deployment_profile_id = "${data.turbonomic_deployment_profile.example.id}"

  compute_resource {
    name = "numOfCpu"
    value = 8
  }
  compute_resource {
    name = "cpuSpeed"
    units = "MHz"
    value = 3000
  }
  compute_resource {
    name = "memorySize"
    units = "MB"
    value = 4096
  }

  // NOTE(ALL): These appear to be added by default?

  compute_resource {
    name  = "networkThroughput"
    units = "MB/s"
    value = 0
  }
  compute_resource {
    name  = "memoryConsumedFactor"
    units = "%"
    value = 75
  }
  compute_resource {
    name  = "cpuConsumedFactor"
    units = "%"
    value = 50
  }
  compute_resource {
    name  = "ioThroughput"
    units = "MB/s"
    value = 0
  }

  storage_resource {
    name  = "diskConsumedFactor"
    units = "%"
    value = 100
  }
  storage_resource {
    name  = "diskIopsConsumed"
    value = 0
  }
  storage_resource {
    name  = "diskSize"
    units = "GB"
    value = 0
  }

}
*/
