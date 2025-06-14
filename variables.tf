variable "instance_region" {
  type        = string
  default     = "us-east-1"
  description = "The region of the project"
}

variable "param_name" {
  type = string
  default = "MyParam"
  description = "The name of the param."
}