## Simple
output "simple_id" {
  value = module.test.id
}
output "simple_arn" {
  value = module.test.arn
}

# Versionning bucket
output "versionning_id" {
  value = module.test_versionning.id
}
output "versionning_arn" {
  value = module.test_versionning.arn
}
