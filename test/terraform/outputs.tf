##┬áSimple
output "simple_id" {
  value = module.test.id
}
output "simple_arn" {
  value = module.test.arn
}

#┬áversioning bucket
output "versioning_id" {
  value = module.test_versioning.id
}
output "versioning_arn" {
  value = module.test_versioning.arn
}
