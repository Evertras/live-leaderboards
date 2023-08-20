locals {
  prefix = terraform.workspace == "default" ? "evertras-leaderboards" : "evertras-leaderboards-${terraform.workspace}"
}
