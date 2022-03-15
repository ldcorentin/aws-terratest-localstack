module "test" {
  source = "../../module"

  project     = var.project
  name        = "s3-simple"
  environment = var.environment
}

module "test_versioning" {
  source = "../../module"

  project            = var.project
  name               = "s3-versioning"
  environment        = var.environment
  versioning_enabled = true
}
