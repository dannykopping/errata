version = "0.1"

options {
  base_url = "https://dannykopping.github.io/errata/"
  prefix = "err-"
  imports = [
    "fmt",
    "github.com/hashicorp/hcl/v2",
  ]
  description = "This is a description"
}

error "code-1" {
  message = "This is a basic error"
}