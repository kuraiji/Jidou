output "smm_param_name" {
  value = module.smm.param_name
  sensitive = true
}

output "smm_param_value" {
  value = module.smm.param_value
  sensitive = true
}

output "dsql_endpoint" {
  value = "${module.dsql.dsql_id}.dsql.${var.instance_region}.on.aws"
  sensitive = true
}