output "api_gateway_prod_url" {
  description = "Base URL for API Gateway stage."

  value = aws_apigatewayv2_stage.prod.invoke_url
}

output "lambda_server_function_name" {
  value = module.lambda_api.function_name
}

output "api_gw_id" {
  value = aws_apigatewayv2_api.api.id
}

output "api_gw_stage" {
  value = aws_apigatewayv2_stage.prod.id
}

output "api_update_key" {
  value = aws_iam_access_key.deploy_api.id
}

output "api_update_secret" {
  value     = aws_iam_access_key.deploy_api.secret
  sensitive = true
}
