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