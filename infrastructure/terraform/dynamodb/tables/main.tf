
provider "aws" {
  region = "${var.aws_region}"
}

terraform {
  backend "s3" {}
  required_version = ">= 0.12.0"
}

module "beers_table" {
  source         = "git::git@github.com:benjaminbartels/terraform-modules.git//dynamodb"
  name           = "beers"
  read_capacity  = 5
  write_capacity = 5
  aws_region     = "${var.aws_region}"
}

module "chambers_table" {
  source         = "git::git@github.com:benjaminbartels/terraform-modules.git//dynamodb"
  name           = "chambers"
  read_capacity  = 5
  write_capacity = 5
  aws_region     = "${var.aws_region}"
}

module "fermentations_table" {
  source         = "git::git@github.com:benjaminbartels/terraform-modules.git//dynamodb"
  name           = "fermentations"
  read_capacity  = 5
  write_capacity = 5
  aws_region     = "${var.aws_region}"
}
