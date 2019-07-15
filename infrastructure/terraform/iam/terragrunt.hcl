terraform {
  source = "git::git@github.com:benjaminbartels/terraform-modules.git//lambda-dynamodb-iam-role"
}

inputs = {
  app_name   = "todo"
  aws_region = "us-west-2"
}

include {
  path = "${find_in_parent_folders()}"
}
