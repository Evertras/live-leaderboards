module "lambda_api" {
  source = "./modules/lambda"

  name        = "api"
  binary_name = "leaderboard-api-lambda"
  prefix      = local.prefix

  environment_vars = {
    EVERTRAS_LEADERBOARD_TABLE = aws_dynamodb_table.events.name
  }
}
