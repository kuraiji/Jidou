module "s3" {
  source = "./modules/s3"
}

module "dynamodb" {
  source = "./modules/dynamodb"
}