terraform {
  source = "git::git@github.com:benjaminbartels/terraform-modules.git//lambda-dynamodb-iam-role"
}

inputs = {
  aws_region = "us-west-2"
  app_name   = "zymurgauge"
}

include {
  path = "${find_in_parent_folders()}"
}
