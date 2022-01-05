data "cloudfunction_invoke_data_source" "my-data" {
  cloud_function_url = "https://europe-west2-project.cloudfunctions.net/my-function"
}