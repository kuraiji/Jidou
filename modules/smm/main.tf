resource "aws_ssm_parameter" "param" {
  name = var.param_name
  type = "String"
  value = var.param_value
}

resource "aws_ssm_parameter" "backend_ip" {
  name = "/JIDOU-API/BACKEND_IP"
  type = "String"
  value = var.backend_ip
}

resource "aws_ssm_parameter" "backend_port" {
  name = "/JIDOU-API/BACKEND_PORT"
  type = "String"
  value = var.backend_port
}