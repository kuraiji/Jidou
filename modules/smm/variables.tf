variable "param_name" {
  type = string
  default = "MyParam"
  description = "The name of the param."
}

variable "param_value" {
  type = string
  description = "The password of the instance"
  default = ""
}

variable "backend_ip" {
  type = string
  description = "The ip of the backend server"
  default = "127.0.0.1"
}

variable "backend_port" {
  type = string
  description = "The port of the backend server app"
  default = ""
}