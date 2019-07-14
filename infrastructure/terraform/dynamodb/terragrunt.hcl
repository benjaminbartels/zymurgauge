terraform {
  source = ".//tables"
}

include {
  path = "${find_in_parent_folders()}"
}
