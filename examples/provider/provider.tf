variable "credentials" {}

provider "google-cloudfunction-https-trigger" {
  credentials_json = var.credentials
}