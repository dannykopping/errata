version = "0.1"

options {
  prefix   = "errata-"
  base_url = "https://dannykopping.github.io/errata/errata/"
  description = "Below is a set of errors that the `eish` program can return."
}

error "file-not-found" {
  message    = "File path %q is incorrect or inaccessible"
  categories = ["file"]
  guide      = "Ensure the given file exists and can be access by errata"
  args       = [
    arg("path", "string")
  ]
  labels     = {
    level : "warning",
  }
}
error "file-not-readable" {
  message    = "File %q is unreadable"
  categories = ["file"]
  guide      = "Ensure the given file can be read by errata"
  args       = [
    arg("path", "string")
  ]
  labels     = {
    level : "warning",
  }
}
error "invalid-definitions" {
  message    = "One or more definitions declared in %q are invalid"
  categories = ["definitions", "validation"]
  guide      = "Review the error(s) and try again"
  args       = [
    arg("path", "string")
  ]
  labels     = {
    level : "error",
  }
}

error "invalid-syntax" {
  message    = "File is malformed"
  categories = ["parsing"]
  guide      = "Check the given datasource file for syntax errors"
}

error "invalid-datasource" {
  message    = "Datasource is invalid"
  categories = ["datasource"]
  guide      = "Check the given datasource file for errors"
}

error "code-gen" {
  message    = "Code generation failed"
  categories = ["codegen"]
}
error "markdown-render" {
  message    = "Markdown rendering failed"
  categories = ["web-ui"]
}

error "template-not-found" {
  message    = "Template path is incorrect or inaccessible"
  categories = ["file"]
}
error "template-not-readable" {
  message    = "Template path is unreadable"
  categories = ["file"]
}
error "template-syntax" {
  message    = "Syntax error in template"
  categories = ["codegen"]
}
error "template-execution" {
  message    = "Error in template execution"
  cause      = "Possible use of missing or renamed field"
  categories = ["codegen"]
}
error "serve-web-ui" {
  message    = "Cannot serve web UI for datasource %q"
  args       = [
    arg("path", "string")
  ]
  categories = ["serve", "web-ui"]
}
error "serve-unknown-code" {
  message    = "Cannot find error definition for given code %q"
  args       = [
    arg("code", "string")
  ]
  categories = ["serve", "web-ui"]
}
error "serve-search-index" {
  message    = "Failed to build search index"
  categories = ["serve", "web-ui", "search"]
}
error "serve-search-missing-term" {
  message    = "Search request is missing a \"term\" query string parameter"
  categories = ["serve", "web-ui", "search"]
}