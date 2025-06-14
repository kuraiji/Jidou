resource "aws_dsql_cluster" "dsql" {
  deletion_protection_enabled = var.is_deletion_protected

  tags = {
    Name = var.cluster_name
  }

  lifecycle {
    prevent_destroy = true
  }
}