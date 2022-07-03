# errata

**errata** is a general purpose, language-agnostic toolkit for _error code enumeration (ECE)_.

ECE is the process of defining all the ways in which your software can fail. Think of it as negative documentation (describing how your system _fails_ as opposed to how it _works_).

Consider a web application. _Aren't HTTP status codes enough?_ No. They are sometimes useful (`404 Not Found` is fairly clear), but others are so vague as to be pretty meaningless (`500 Internal Server Error`). Often an additional semantic layer is required to communicate what exactly went wrong, and what the caller (be they a human, an API client, etc) can do in response. HTTP status codes are perfect for describing the _category_ of problem (4xx client error, 5xx server error), but insufficient for the complex software of today. 

Major API vendors such as [Twilio](https://www.twilio.com/docs/api/errors), [Instagram](https://developers.facebook.com/docs/instagram-api/reference/error-codes/), and [Google Cloud](https://cloud.google.com/resource-manager/docs/core_errors) recognise this and use ECE in some form. **errata** aims to provide a mechanism for any software to deliver world-class error handling.

## Basic Concepts

**errata**'s philosophy is that errors should have _at least_ a static error `code` and a human-readable `message`.

- the `code` is searchable, and because it's static it becomes _more easily_ searchable
- the `message` is displayed to the user alongside the code, to provide immediate context, and ideally to give an insight into what went wrong

Besides the basic `code` and `message`, including other valuable metadata like a unique reference (particularly useful in SaaS applications), labels, user guides, etc can be included.

### Definitions

**errata** uses the [HCL](https://github.com/hashicorp/hcl) structured configuration language, used primarily by [Terraform](https://www.terraform.io/). It's extensible, simple to read and write, and frankly - _fuck YAML_.

```hcl
version = "0.1"

error "file-not-found" {
  message    = "File path is incorrect or inaccessible"
  categories = ["file"]
  guide      = "Ensure the given file exists and can be accessed"
  args       = [
    arg("path", "string")
  ]
  labels     = {
    severity = "warn"
  }
}

...
```

The above example defines the `code` (**file-not-found**) and the `message`, along with some other useful metadata (more on this below).

So, what can this definitions file be used for?

### Code Generation

**errata** provides a language-agnostic mechanism for generating code based on these definitions using the [Pongo2](https://github.com/flosch/pongo2) templating engine.

**errata** comes with a CLI tool called **eish** (_**e**rrata **i**nteractive **sh**ell_,
[pronounced "eɪʃ"](http://ipa-reader.xyz/?text=e%C9%AA%CA%83)) which generates code based on given **errata** definitions.

```bash
$ eish generate --source=errata.hcl --template=golang --package=errors
```

This will generate a single file with all error definitions. See the [sample application](sample/) which uses **errata** definitions (and rather recursively, the **errata** library also uses [**errata** definitions](errata.hcl)).

### Web UI

`eish` also provides a simple web UI, allowing your **errata** definitions to be viewed and searched.

```bash
$ eish serve --source=errata.hcl
```

The web UI by default runs on port `37707`.

![Web UI](web-ui.png)

### Supported Languages

_If your language of choice is not yet available, consider contributing a template!_

> _**errata** uses the [Pongo2](https://github.com/flosch/pongo2) templating engine_

- [Golang](templates/golang.tmpl) (reference implementation)

The code produced by **errata** aims to be:

- idiomatic
- easy to use
- using native error/exception types as much as possible