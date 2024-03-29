package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"google-cloudfunction-https-trigger": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TF_ACC"); v == "" {
		t.Skip("Acceptance tests and bootstrapping skipped unless env 'TF_ACC' set")
		return
	}

	if v := os.Getenv("CREDENTIALS_JSON"); v == "" {
		t.Fatal("CREDENTIALS_JSON must be set for acceptance tests")
	}

	if v := os.Getenv("CLOUD_FUNCTION_URL"); v == "" {
		t.Fatal("CLOUD_FUNCTION_URL must be set for acceptance tests")
	}
}
