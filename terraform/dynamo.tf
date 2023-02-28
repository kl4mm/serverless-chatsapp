module "connections" {
  source = "terraform-aws-modules/dynamodb-table/aws"

  name     = "${local.app}-connections"
  hash_key = "connectionId"

  attributes = [
    {
      name = "connectionId"
      type = "S"
    }
  ]
}
