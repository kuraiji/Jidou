variable "instance_region" {
  type        = string
  default     = "us-east-1"
  description = "The region of the project"
}

variable "PARAM_NAME" {
  type = string
  default = "MyParam"
  description = "The name of the param."
}

variable "IMAGE_URI" {
  type = string
  default = ""
  description = "The uri of the image file."
}

variable "backend_port" {
  type = number
  description = "The port that the backend listens on."
  default = 1323
}

variable "FRONTEND_IMAGE_URI" {
  type = string
  default = ""
  description = "The uri of the image file."
}

variable FRONTEND_AWS_ACCESS_KEY_ID {
  type = string
  default = ""
  description = "Permission restricted AKI."
}

variable FRONTEND_AWS_SECRET_ACCESS_KEY {
  type = string
  default = ""
  description = "Permission restricted ASAK."
}