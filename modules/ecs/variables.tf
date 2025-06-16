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