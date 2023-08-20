resource "aws_dynamodb_table" "events" {
  name = "${local.prefix}-events"

  billing_mode = "PAY_PER_REQUEST"

  stream_enabled = true
  stream_view_type = "NEW_IMAGE"

  hash_key = "pk"
  range_key = "sk"

  attribute {
    name = "pk"
    type = "S"
  }

  attribute {
    name = "sk"
    type = "S"
  }
}
