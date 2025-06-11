resource "aws_dynamodb_table" "table" {
  name = "JIDOU"
  billing_mode = "PAY_PER_REQUEST"
  hash_key = "Date"
  range_key = "Name"

  attribute {
    name = "Date"
    type = "S"
  }

  attribute {
    name = "Name"
    type = "S"
  }

  lifecycle {
    prevent_destroy = true
  }
}