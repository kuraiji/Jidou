module "s3" {
  source = "./modules/s3"
}

resource "random_uuid" "param_uuid" {}

module "smm" {
  source = "./modules/smm"
  param_name = var.param_name
  param_value = random_uuid.param_uuid.id
}

module "dsql" {
  source = "./modules/dsql"
  cluster_name = "JIDOU"
  is_deletion_protected = false
}