locals {
  name = "${var.project}-${var.name}-${var.environment}-${var.region}"
}

resource "aws_s3_bucket" "default" {
  bucket = local.name

  tags = merge({
    Name        = local.name
    Project     = var.project
    region      = var.region
    Environment = var.environment
    Account     = data.aws_caller_identity.current.account_id
    CreatedBy   = "Terraform"
  }, var.extra_tags)
}

resource "aws_s3_bucket_policy" "bucket_policy" {
  bucket = aws_s3_bucket.default.id
  policy = data.aws_iam_policy_document.policy.json
}

data "aws_iam_policy_document" "policy" {

  ## Policy statement to prevent unsecure access to the bucket
  statement {
    sid    = "Deny-non-secure-transport"
    effect = "Deny"
    actions = [
      "s3:*"
    ]
    resources = [
      "${aws_s3_bucket.default.arn}/*",
    ]
    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
    condition {
      test     = "Bool"
      variable = "aws:SecureTransport"
      values = [
        "false"
      ]
    }
  }
}

