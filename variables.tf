variable "instance_name" {
  type        = string
  default     = "Jidou"
  description = "The name of the project"
}

variable "instance_region" {
  type        = string
  default     = "us-west-1"
  description = "The region of the project"
}

variable "bucket_name" {
  type        = string
  default     = "jidou"
  description = "The name of the bucket"
}

variable "instance_environment" {
  type        = string
  default     = "Dev"
  description = "The name of the environment"
}

variable "index_document" {
  type    = string
  default = "index.html"
}
variable "error_document" {
  type    = string
  default = "error.html"
}