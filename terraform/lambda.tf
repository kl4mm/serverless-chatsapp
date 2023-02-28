module "connect" {
  source        = "terraform-aws-modules/lambda/aws"
  function_name = "${local.app}-connect"
  handler       = "main"
  runtime       = "go1.x"

  attach_policy_json = true
  policy_json        = data.aws_iam_policy_document.connect.json

  create_package         = false
  local_existing_package = "../cmd/connect/connect.zip"

  create_current_version_allowed_triggers = false
  allowed_triggers = {
    "APIGateway" = {
      service    = "apigateway"
      source_arn = "${aws_apigatewayv2_stage.websocket.execution_arn}/$connect"
    }
  }

  environment_variables = {
    "table" = "${local.app}-connections"
  }
}

module "disconnect" {
  source        = "terraform-aws-modules/lambda/aws"
  function_name = "${local.app}-disconnect"
  handler       = "main"
  runtime       = "go1.x"

  attach_policy_json = true
  policy_json        = data.aws_iam_policy_document.connect.json

  create_package         = false
  local_existing_package = "../cmd/disconnect/disconnect.zip"

  create_current_version_allowed_triggers = false
  allowed_triggers = {
    "APIGateway" = {
      service    = "apigateway"
      source_arn = "${aws_apigatewayv2_stage.websocket.execution_arn}/$disconnect"
    }
  }

  environment_variables = {
    "table" = "${local.app}-connections"
  }
}

data "aws_iam_policy_document" "connect" {
  statement {

    actions = [
      "dynamodb:BatchWriteItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem"
    ]

    resources = [
      module.connections.dynamodb_table_arn
    ]
  }
}

module "default" {
  source        = "terraform-aws-modules/lambda/aws"
  function_name = "${local.app}-default"
  handler       = "main"
  runtime       = "go1.x"

  attach_policy_json = true
  policy_json        = data.aws_iam_policy_document.default.json

  create_package         = false
  local_existing_package = "../cmd/default/default.zip"

  create_current_version_allowed_triggers = false
  allowed_triggers = {
    "APIGateway" = {
      service    = "apigateway"
      source_arn = "${aws_apigatewayv2_stage.websocket.execution_arn}/$default"
    }
  }

  environment_variables = {
    "table" = "${local.app}-connections"
  }
}

data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "default" {
  statement {

    actions = [
      "dynamodb:BatchWriteItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem",
      "dynamodb:DeleteItem"
    ]

    resources = [
      module.connections.dynamodb_table_arn
    ]
  }

  statement {
    actions = [
      "execute-api:ManageConnections"
    ]

    resources = [
      "arn:aws:execute-api:eu-west-2:${data.aws_caller_identity.current.account_id}:*/*/POST/@connections/*",
      "arn:aws:execute-api:eu-west-2:${data.aws_caller_identity.current.account_id}:*/*/GET/@connections/*"
    ]
  }
}

module "sendmessage" {
  source        = "terraform-aws-modules/lambda/aws"
  function_name = "${local.app}-send-message"
  handler       = "main"
  runtime       = "go1.x"

  attach_policy_json = true
  policy_json        = data.aws_iam_policy_document.sendmessage.json

  create_package         = false
  local_existing_package = "../cmd/sendmessage/sendmessage.zip"

  create_current_version_allowed_triggers = false
  allowed_triggers = {
    "APIGateway" = {
      service    = "apigateway"
      source_arn = "${aws_apigatewayv2_stage.websocket.execution_arn}/sendmessage"
    }
  }

  environment_variables = {
    "table" = "${local.app}-connections"
  }
}

data "aws_iam_policy_document" "sendmessage" {
  statement {

    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:GetRecords",
      "dynamodb:GetShardIterator",
      "dynamodb:Query",
      "dynamodb:GetItem",
      "dynamodb:Scan",
      "dynamodb:ConditionCheckItem"
    ]

    resources = [
      module.connections.dynamodb_table_arn
    ]
  }
}
