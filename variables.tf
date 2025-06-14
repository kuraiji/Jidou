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