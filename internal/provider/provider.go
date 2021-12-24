package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/option"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"credentials_json": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("CREDENTIALS_JSON", nil),
					Description: "Service account JSON key",
					Sensitive:   true,
				},
			},

			DataSourcesMap: map[string]*schema.Resource{
				"cloudfunction_invoke_data_source": dataSourceCloudFunctionInvoke(),
			},
			ResourcesMap: map[string]*schema.Resource{},

			ConfigureContextFunc: providerConfigure,
		}

		return p
	}
}

type Authentication struct {
	credentials option.ClientOption
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (
	client interface{},
	diags diag.Diagnostics,
) {
	credentials := d.Get("credentials_json").(string)
	client = Authentication{
		credentials: option.WithCredentialsJSON([]byte(credentials)),
	}

	return client, diags
}
