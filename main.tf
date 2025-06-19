/* "s3" {
  source = "./modules/s3"
}*/

resource "random_uuid" "param_uuid" {}

module "smm" {
  source = "./modules/smm"
  param_name = var.PARAM_NAME
  param_value = random_uuid.param_uuid.id
}

/*module "ecs" {
  source = "./modules/ecs"
  cluster_name = "jidou"
  image_uri = var.IMAGE_URI
  exposed_port = var.backend_port
  ssh_key_name = "jidou_key"
  region = var.instance_region
  frontend_image_uri = var.FRONTEND_IMAGE_URI
  frontend_port = 3000
  frontend_aki = var.FRONTEND_AWS_ACCESS_KEY_ID
  frontend_asak = var.FRONTEND_AWS_SECRET_ACCESS_KEY
}*/

/*module "dsql" {
  source = "./modules/dsql"
  cluster_name = "JIDOU"
  is_deletion_protected = false
}*/