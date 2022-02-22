terraform {
  backend "remote" {
    hostname     = "app.terraform.io"
    organization = "takescoop-oss"

    workspaces {
      name = "remote-state-action-e2e-test"
    }
  }
}