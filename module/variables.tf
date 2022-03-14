variable "region" {
  default     = "eu-west-1"
  description = "The aws region."
  type        = string
}

variable "name" {
  type        = string
  default     = ""
  description = "The name of the bucket is build like that : name = var.name != empty ? {var.project}-{var.name}-{local.suffix} : {local.prefix}-{local.suffix}"
  validation {
    condition     = length(var.name) < 20
    error_message = "Your name is too long."
  }
}

variable "environment" {
  type        = string
  description = "localstack, dev, test, integration, prod, etc."
}

variable "project" {
  type        = string
  description = "The project name."
}

variable "extra_tags" {
  description = "If you need extra tags for billings etc, add it here."
  default     = {}
  type        = map(any)
}

## Versionning
variable "versioning_enabled" {
  type        = bool
  description = "The versioning state of the bucket. Valid values: Enabled or Suspended."
  default     = false
}

variable "versioning_mfa" {
  type        = string
  default     = "Disabled"
  description = "(Optional, Required if versioning_configuration mfa_delete is enabled) The concatenation of the authentication device's serial number, a space, and the value that is displayed on your authentication device."
}

variable "versioning_mfa_delete" {
  type        = string
  default     = "Disabled"
  description = "Specifies whether MFA delete is enabled in the bucket versioning configuration. Valid values: Enabled or Disabled."
}

