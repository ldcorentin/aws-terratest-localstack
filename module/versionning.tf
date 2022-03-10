resource "aws_s3_bucket_versioning" "default" {
  count = var.versioning_enabled ? 1 : 0

  bucket                = aws_s3_bucket.default.id
  expected_bucket_owner = data.aws_caller_identity.current.account_id
  mfa                   = var.versioning_mfa

  versioning_configuration {
    status     = var.versioning_enabled ? "Enabled" : "Suspended"
    mfa_delete = var.versioning_mfa_delete
  }
}
