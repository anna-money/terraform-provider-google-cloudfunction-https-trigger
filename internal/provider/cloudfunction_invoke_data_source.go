package provider

import (
	"context"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
	htransport "google.golang.org/api/transport/http"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudFunctionInvoke() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Data source that request a cloud function and fetch a response body and headers",

		ReadContext: dataSourceCloudFunctionInvokeRead,

		Schema: map[string]*schema.Schema{
			"cloud_function_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "URL of the Cloud Function",
			},
			"body": {
				Type:     schema.TypeString,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Returned body of a request",
			},
			"response_headers": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "HTTP headers from a response",
			},
		},
	}
}

func dataSourceCloudFunctionInvokeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	authentication := meta.(Authentication)
	targetURL := d.Get("cloud_function_url").(string)

	client, err := newClient(ctx, targetURL, authentication.credentials)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	resp, err := client.Get(targetURL)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return append(diags, diag.Errorf("HTTP request error. Response code: %d", resp.StatusCode)...)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	responseHeaders := make(map[string]string)
	for k, v := range resp.Header {
		responseHeaders[k] = strings.Join(v, ", ")
	}

	err = d.Set("body", string(bytes))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}

	if err = d.Set("response_headers", responseHeaders); err != nil {
		return append(diags, diag.Errorf("Error setting HTTP response headers: %s", err)...)
	}

	d.SetId(targetURL)

	return diags
}

func newClient(ctx context.Context, audience string, opts ...option.ClientOption) (*http.Client, error) {
	ts, err := idtoken.NewTokenSource(ctx, audience, opts...)
	if err != nil {
		return nil, err
	}
	opts = []option.ClientOption{option.WithTokenSource(ts)}
	t, err := htransport.NewTransport(ctx, http.DefaultTransport, opts...)
	if err != nil {
		return nil, err
	}
	return &http.Client{Transport: t}, nil
}
