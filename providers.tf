terraform {
  cloud {
    organization = "Kuraiji"
    workspaces {
      name = "Jidou"
    }
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.100.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "0.12.1"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5"
    }
  }
  required_version = ">= 1.2.0"
}

provider "aws" {
  region = var.instance_region
}

provider "cloudflare" {}