module "s3" {
  source = "./modules/s3"
}

resource "random_uuid" "param_uuid" {}

module "smm" {
  source = "./modules/smm"
  param_name = var.PARAM_NAME
  param_value = random_uuid.param_uuid.id
}

module "ecs" {
  source = "./modules/ecs"
  cluster_name = "jidou"
  image_uri = var.IMAGE_URI
}

/*module "dsql" {
  source = "./modules/dsql"
  cluster_name = "JIDOU"
  is_deletion_protected = false
}*/