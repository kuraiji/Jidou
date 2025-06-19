variable "cluster_name" {
  type = string
  default = "MyCluster"
  description = "The name of the ecs cluster."
}

variable "image_uri" {
  type = string
  default = ""
  description = "The uri of the image file."
}

variable "exposed_port" {
  type = number
  default = 80
  description = "The exposed port for the application."
}

variable "ssh_key_name" {
  type = string
  default = "MyKey"
  description = "The name of the ssh key you have created on AWS EC2 console."
}

variable "region" {
  type = string
  default = "us-west-1"
  description = "The region that will be set for the env variable."
}

variable "frontend_image_uri" {
  type = string
  default = ""
  description = "The uri of the front end image file."
}

variable "frontend_port" {
  type = number
  default = 80
  description = "The port that the frontend listens on."
}

variable "frontend_aki" {
  type = string
  default = ""
  description = "Permission restricted AKI."
}

variable "frontend_asak" {
  type = string
  default = ""
  description = "Permission restricted ASAK."
}

variable "zone_id" {
  type = string
  default = ""
  description = "Zone ID."
}