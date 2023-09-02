resource "aws_iam_user" "deploy_api" {
  name = "${local.prefix}-deploy-api"
}

resource "aws_iam_access_key" "deploy_api" {
  user = aws_iam_user.deploy_api.name
}

resource "aws_iam_user_policy" "deploy_api" {
  user = aws_iam_user.deploy_api.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "lambda:UpdateFunctionCode"
        ]
        Effect = "Allow"
        Resource = [
          module.lambda_api.arn
        ]
      }
    ]
  })
}
