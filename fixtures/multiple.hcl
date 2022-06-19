version = "0.1"

error "code-1" {
  message = "This is a basic error"
}
error "code-2" {
  message = "This is another basic error"
  args = [
    arg("first", "string"),
    arg("second", "bool"),
  ]
}
error "code-3" {
  message = "This one has a guide file"
  guide = file("fixtures/guide.md")
}
error "code-4" {
  message = "This one has a guide as a string"
  guide = "# Hello Errata"
}