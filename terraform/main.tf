terraform {
  backend "s3" {
    bucket     = "terraform-state-storage-586877430255"
    dynamodb_table = "terraform-state-lock-586877430255"
    region     = "us-west-2"

    // THIS MUST BE UNIQUE
    key = "av-control-api.tfstate"
  }
}

provider "aws" {
  region = "us-west-2"
}

data "aws_ssm_parameter" "eks_cluster_endpoint" {
  name = "/eks/av-cluster-endpoint"
}

provider "kubernetes" {
  host = data.aws_ssm_parameter.eks_cluster_endpoint.value
}

// pull all env vars out of ssm

data "aws_ssm_parameter" "atlona_username" {
  name = "/atlona/username"
}

data "aws_ssm_parameter" "atlona_password" {
  name = "/atlona/password"
}

data "aws_ssm_parameter" "endpoint_auth_url" {
  name = "/av-api/endpoint_auth_url"
}

data "aws_ssm_parameter" "messenger_hub_address" {
  name = "/env/hub-address"
}

data "aws_ssm_parameter" "prd_db_address" {
  name = "/env/couch-address"
}

data "aws_ssm_parameter" "prd_db_password" {
  name = "/env/couch-password"
}

data "aws_ssm_parameter" "prd_db_username" {
  name = "/env/couch-username"
}

data "aws_ssm_parameter" "sonyrest_psk" {
  name = "/sonyrest/psk"
}

module "av_api_prd" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "av-api-prd"
  image          = "byuoitav/av-api"
  image_version  = "development"
  container_port = 8000
  repo_url       = "https://github.com/byuoitav/av-control-api"

  // optional
  public_urls = ["api.av.byu.edu"]
  container_env = {
    DB_ADDRESS                 = data.aws_ssm_parameter.prd_db_address.value
    DB_PASSWORD                = data.aws_ssm_parameter.prd_db_password.value
    DB_USERNAME                = data.aws_ssm_parameter.prd_db_username.value
    ENDPOINT_AUTHORIZATION_URL = data.aws_ssm_parameter.endpoint_auth_url.value
    HUB_ADDRESS                = data.aws_ssm_parameter.messenger_hub_address.value
    STOP_REPLICATION           = "true"
    SYSTEM_ID                  = "aws-avapi-prd"
  }
}

module "atlona_driver_prd" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "atlona-driver-prd"
  image          = "docker.pkg.github.com/byuoitav/av-control-api/atlona-driver"
  image_version  = "v0.2.2"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/av-control-api"

  // optional
  image_pull_secret = "github-docker-registry"
  container_args = [
    "--username", data.aws_ssm_parameter.atlona_username.value, // Atlona device username
    "--password", data.aws_ssm_parameter.atlona_password.value  // Atlona device password
  ]
}

module "justaddpower_driver_prd" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "justaddpower-driver-prd"
  image          = "docker.pkg.github.com/byuoitav/av-control-api/justaddpower-driver"
  image_version  = "v0.1.1"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/av-control-api"

  // optional
  image_pull_secret = "github-docker-registry"
  container_args = [
    "--port", 8080
  ]
}

module "sonyrest_driver_prd" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "sonyrest-driver-prd"
  image          = "docker.pkg.github.com/byuoitav/av-control-api/sonyrest-driver-dev"
  image_version  = "26ede10"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/av-control-api"

  // optional
  image_pull_secret = "github-docker-registry"
  container_args = [
    "--psk", data.aws_ssm_parameter.sonyrest_psk.value // Sony device psk
  ]
}

module "adcp_driver_prd" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "adcp-driver-prd"
  image          = "docker.pkg.github.com/byuoitav/av-control-api/adcp-driver-dev"
  image_version  = "1639c73"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/av-control-api"

  // optional
  image_pull_secret = "github-docker-registry"
}
