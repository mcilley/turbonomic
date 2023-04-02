package turbonomic

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// var testAccProviders map[string]terraform.ResourceProvider
// var testAccProvider *schema.Provider

// func init() {
// 	testAccProvider = Provider().(*schema.Provider)
// 	testAccProviders = map[string]terraform.ResourceProvider{
// 		"turbonomic": testAccProvider,
// 	}
// }

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

// func testAccPreCheck(t *testing.T) {
// 	if v := os.Getenv("SERVER_HOSTNAME"); v == "" {
// 		t.Fatal("SERVER_HOSTNAME must be set for acceptance tests")
// 	}
// 	if v := os.Getenv("CLIENT_USERNAME"); v == "" {
// 		t.Fatal("CLIENT_USERNAME must be set for acceptance tests")
// 	}
// 	if v := os.Getenv("CLIENT_PASSWORD"); v == "" {
// 		t.Fatal("CLIENT_PASSWORD must be set for acceptance tests")
// 	}
// }

// func skipIfEnvNotSet(t *testing.T, envs ...string) {
// 	for _, k := range envs {
// 		if os.Getenv(k) == "" {
// 			t.Skipf("Environment variable %s is not set", k)
// 		}
// 	}
// }
