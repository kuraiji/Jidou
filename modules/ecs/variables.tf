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