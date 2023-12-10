provider "aws" {
  region = "eu-central-1" 
}

resource "aws_dynamodb_table" "table" {
  name           = "ShortLinkDB" 
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "id"

  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_lambda_function" "generate_link" {
  function_name = "GenerateLinkFunction"
  handler       = "bootstrap"
  runtime       = "provided.al2"
  memory_size   = 128
  timeout       = 5

  environment {
    variables = {
      TABLE = aws_dynamodb_table.table.name
    }
  }

  tracing_config {
    mode = "Active"
  }

  filename         = "functions/generate_link/deployment-package.zip"
  source_code_hash = filebase64sha256("functions/generate_link/deployment-package.zip")

  role = aws_iam_role.lambda_exec.arn
}

resource "aws_iam_role" "lambda_exec" {
  name = "lambda_exec_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "lambda.amazonaws.com"
        },
      },
    ],
  })
}

resource "aws_iam_policy_attachment" "lambda_basic_execution" {
  name       = "lambda_basic_execution"
  roles      = [aws_iam_role.lambda_exec.name]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}


resource "aws_iam_policy" "lambda_policy" {
  name        = "lambda_policy"
  path        = "/"
  description = "IAM policy for Lambda to access DynamoDB"

  policy = jsonencode({
    Version   = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Action    = "dynamodb:PutItem",
        Resource  = aws_dynamodb_table.table.arn,
      },
    ],
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy_attach" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_policy.arn
}

resource "aws_apigatewayv2_api" "http_api" {
  name          = "MyHttpApi"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_deployment" "http_api_deployment" {
  api_id = aws_apigatewayv2_api.http_api.id

  depends_on = [
    aws_apigatewayv2_route.http_api_route,
  ]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_apigatewayv2_stage" "http_api_stage" {
  api_id        = aws_apigatewayv2_api.http_api.id
  name          = "$default" 
  deployment_id = aws_apigatewayv2_deployment.http_api_deployment.id
}

resource "aws_apigatewayv2_integration" "http_api_integration" {
  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.generate_link.invoke_arn
  payload_format_version = "2.0"
}


resource "aws_apigatewayv2_route" "http_api_route" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "PUT /links"
  target    = "integrations/${aws_apigatewayv2_integration.http_api_integration.id}"
}

resource "aws_lambda_permission" "api_gateway_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.generate_link.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
  depends_on = [aws_apigatewayv2_deployment.http_api_deployment]
}

output "api_url" {
  description = "API Gateway endpoint URL for the GenerateLink function"
  value       = "${aws_apigatewayv2_api.http_api.api_endpoint}/links"
}