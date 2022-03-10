module "test" {
  source = "../../module"

  project     = var.project
  name        = "s3-simple"
  environment = var.environment
}

module "test_versionning" {
  source = "../../module"

  project            = var.project
  name               = "s3-versionning"
  environment        = var.environment
  versioning_enabled = true
  # versioning_mfa        = ""
  # versioning_mfa_delete = "Disabled"
}
