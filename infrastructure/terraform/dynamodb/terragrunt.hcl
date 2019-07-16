terraform {
  source = "github.com/benjaminbartels/terraform-modules.git//dynamodb"
}

inputs = {
  name           = "todos"
  read_capacity  = 5
  write_capacity = 5
  aws_region     = "us-west-2"
}

include {
  path = "${find_in_parent_folders()}"
}
