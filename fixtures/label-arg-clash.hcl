version = "0.1"

error "code-1" {
  message = "This is a basic error"
  args = [
    arg("key", "string")
  ]
  labels = {
    "key": "value"
  }
}