output "smm_param_name" {
  value     = module.smm.param_name
  sensitive = true
}

output "smm_param_value" {
  value     = module.smm.param_value
  sensitive = true
}

output "ec2_url" {
  value     = module.ecs.ec2_url
  sensitive = true
}

output "ec2_dns" {
  value     = module.ecs.ec2_dns
  sensitive = true
}

output "dns_name" {
  value = module.ecs.dns_name
  sensitive = false
}

/*output "dsql_endpoint" {
  value = "${module.dsql.dsql_id}.dsql.${var.instance_region}.on.aws"
  sensitive = true
}*/