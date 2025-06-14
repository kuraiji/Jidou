variable "cluster_name" {
  type = string
  default = "MyCluster"
  description = "The name of the db cluster."
}

variable "is_deletion_protected" {
  type = bool
  default = false
  description = "Sets if the cluster is deletion protected."
}