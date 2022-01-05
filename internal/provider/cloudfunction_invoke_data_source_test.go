package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceCloudFunctionInvoke(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: getConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.google-cloudfunction-https-trigger_cloudfunction_invoke_data_source.foo", "cloud_function_url", regexp.MustCompile(os.Getenv("CLOUD_FUNCTION_URL"))),
				),
			},
		},
	})
}

func getConfig() string {
	return fmt.Sprintf(testAccDataSourceCloudFunctionInvoke, os.Getenv("CLOUD_FUNCTION_URL"))
}

const testAccDataSourceCloudFunctionInvoke = `
data "google-cloudfunction-https-trigger_cloudfunction_invoke_data_source" "foo" {
  cloud_function_url = "%s"
}
`
