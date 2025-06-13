resource "aws_dynamodb_table" "table" {
  name = "JIDOU"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "Date"

  attribute {
    name = "Date"
    type = "S"
  }

  lifecycle {
    prevent_destroy = true
  }
}