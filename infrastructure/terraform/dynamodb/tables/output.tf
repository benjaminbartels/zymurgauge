output "table_arns" {
  value = [
    "${module.beers_table.table_arn}",
    "${module.chambers_table.table_arn}",
    "${module.fermentations_table.table_arn}"
  ]
}
