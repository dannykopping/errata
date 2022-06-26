version = "0.1"

options {
  base_url = "https://errata.codes/errata/sample/"
}

error "missing-values" {
  message    = "The %q field is missing from the request"
  categories = ["login"]
  args       = [
    arg("missingField", "string")
  ]
  labels     = {
    http_response_code : 400
    shell_exit_code : 1
  }
}

error "invalid-request" {
  message    = "One or more values are missing"
  categories = ["login"]
  labels     = {
    http_response_code : 400
    shell_exit_code : 1
  }
}

error "invalid-email" {
  message    = "Given email %q is invalid"
  guide      = "Ensure the email address has a _username_, an `@` symbol, and a _domain_ name."
  categories = ["login"]
  args       = [
    arg("email", "string")
  ]
  labels     = {
    http_response_code : 400
    shell_exit_code : 1
  }
}

# unsuccessful login errors
error "incorrect-credentials" {
  message    = "Given credentials are incorrect"
  categories = ["login"]
  labels     = {
    http_response_code : 403
    shell_exit_code : 2
  }
}

# blocked errors
error "account-blocked-spam" {
  message    = "Account is blocked because of spam"
  categories = ["login"]
  labels     = {
    http_response_code : 403
    shell_exit_code : 3
  }
}
error "account-blocked-abuse" {
  message    = "Account is blocked because of abuse"
  categories = ["login"]
  labels     = {
    http_response_code : 403
    shell_exit_code : 3
  }
}

# response errors
error "response-formatting" {
  message    = "Failed to format response body"
  categories = ["internal"]
  labels     = {
    http_response_code : 500
  }
}
