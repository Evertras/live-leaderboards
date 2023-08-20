terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.13.1"
    }
  }

  # This is set up in another repository's terraform
  backend "s3" {
    bucket  = "evertras-home-terraform"
    key     = "global/s3/leaderboards.state"
    region  = "ap-northeast-1"
    encrypt = true
    profile = "admin"

    dynamodb_table = "evertras-home-terraform-locks"
  }

  required_version = ">= 1.5.5"
}

provider "aws" {
  region  = "ap-northeast-1"
  profile = "admin"

  default_tags {
    tags = {
      Deployment = "${local.prefix}-leaderboards"
    }
  }
}
