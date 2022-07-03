version = "0.1"

options {
  base_url = "https://errata.codes/errata/sample/"
}

error "missing-command" {
  message    = "Command was not specified"
  categories = ["validation"]
  labels     = {
    http_response_code : 400
  }
}

error "script-not-found" {
  message    = "Given script was not found"
  categories = ["validation"]
  args       = [
    arg("path", "string")
  ]
  labels = {
    http_response_code : 404
  }
}

error "script-execution-failed" {
  message    = "Given script returned an error"
  categories = ["execution"]
  labels = {
    http_response_code : 500
  }
}