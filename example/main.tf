module "lambda_archive_rds_query" {
  source = "rojopolis/lambda-python-archive/aws"
  src_dir = "../lambda_function"
  output_path = "${path.module}/lambda_query.zip"
}

resource "aws_lambda_function" "query_rds" {
  filename         = "${module.lambda_archive_rds_query.archive_path}"
  source_code_hash = "${module.lambda_archive_rds_query.source_code_hash}"
  function_name = "lambda_query"
  description = "serverless db query function"
  runtime          = "python3.6"
  role             = "role_arn"
  timeout          = 60
  memory_size      = 128
  handler          = "lambda_function.handler"

  vpc_config {
    subnet_ids = ["your_primary_subnet_id", "your_secondary_subnet_id"]
    security_group_ids = ["your_security_group_id"]
  }

  environment {
    variables = {
      DB_HOST = "db host address" # e.g. aws_db_instance.main.address
      DB_NAME = "database name"
      DB_USER = "username"
      DB_PASS = "database password"
      # SSM_DB_PASS_NAME = aws_ssm_parameter.mysql_master_password.name # you can use SSM parameter as password store
      ENV = "development"
    }
  }
}
