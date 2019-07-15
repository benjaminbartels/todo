remote_state {
  backend = "s3"

  config = {
    encrypt        = true
    bucket         = "todo-terraform-remote-state"
    key            = "${path_relative_to_include()}/terraform.tfstate"
    region         = "us-west-2"
    dynamodb_table = "terraform_locks"
  }
}

terraform {
  extra_arguments "retry_lock" {
    commands = "${get_terraform_commands_that_need_locking()}"

    arguments = [
      "-lock-timeout=10m",
    ]
  }

  extra_arguments "auto_approve" {
    commands = [
      "apply",
    ]

    arguments = [
      "-auto-approve",
    ]
  }
}
