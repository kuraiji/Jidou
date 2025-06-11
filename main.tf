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
      version = "~> 4.16"
    }
  }
  required_version = ">= 1.2.0"
}

provider "aws" {
  region = var.instance_region
}

resource "aws_s3_bucket" "bucket" {
  bucket = var.bucket_name

  tags = {
    Name        = var.instance_name
    Environment = var.instance_environment
  }
}

data "local_file" "index_file" {
  filename = "./index.html"
}

resource "aws_s3_object" "index" {
  bucket       = aws_s3_bucket.bucket.id
  content      = data.local_file.index_file.content
  key          = var.index_document
  content_type = "text/html"
}

data "local_file" "error_file" {
  filename = "./error.html"
}

resource "aws_s3_object" "error" {
  bucket       = aws_s3_bucket.bucket.id
  key          = var.error_document
  content      = data.local_file.error_file.content
  content_type = "text/html"
}

resource "aws_s3_bucket_website_configuration" "website" {
  bucket = aws_s3_bucket.bucket.id

  index_document {
    suffix = var.index_document
  }
  error_document {
    key = var.error_document
  }
}

resource "aws_s3_bucket_public_access_block" "unblock" {
  bucket = aws_s3_bucket.bucket.id
}

data "aws_iam_policy_document" "policy" {
  statement {
    principals {
      type        = "*"
      identifiers = ["*"]
    }
    sid    = "PublicReadGetObject"
    effect = "Allow"
    actions = [
      "s3:GetObject"
    ]
    resources = [
      "${aws_s3_bucket.bucket.arn}/*",
    ]
  }
}

resource "aws_s3_bucket_policy" "bucket_policy" {
  bucket = aws_s3_bucket.bucket.id
  policy = data.aws_iam_policy_document.policy.json
}