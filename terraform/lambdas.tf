module "lambda_api" {
  source = "./modules/lambda"

  name        = "api"
  binary_name = "leaderboard-api-lambda"
  prefix      = local.prefix

  environment_vars = {
    EVERTRAS_LEADERBOARD_TABLE = aws_dynamodb_table.events.name
  }
}

resource "aws_iam_policy" "api_db_policy" {
  name        = "${local.prefix}-api-db-policy"
  description = "API access to write/query database"
  policy      = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": [
       "dynamodb:Query",
       "dynamodb:PutItem",
       "dynamodb:GetItem"
     ],
     "Resource": "${aws_dynamodb_table.events.arn}",
     "Effect": "Allow"
   }
 ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "attach_api_db_policy" {
  role       = module.lambda_api.role_name
  policy_arn = aws_iam_policy.api_db_policy.arn
}
