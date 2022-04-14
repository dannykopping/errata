# Errata

## Background

Modern applications typically involve interactions between a client and a server. As these applications become more complex, so to
do their failure modes. The server is usually delegated as the "source of truth" in the interaction, and the server has to respond
with error codes and messages under failure conditions.

Application clients come in broadly two type:

// TODO improve terminology, too stuffy
1. **Relaying** clients relay server responses to the user ~verbatim, or with minimal interpretation
2. **Conditionalising** clients interpret responses and conditionalise their subsequent behaviour

## Problems

### Type 1 vs Type 2 Clients

HTTP-based APIs typically utilise the set of [status codes](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes) to describe
the responses they generate. Some of these status codes (like `500 Internal Server Error`) are so broad as to be effectively
meaningless. APIs that return an HTTP status code and error message _only_ run the risk of complicating client-side handling if
the client needs to be a **Type 2** client.

Applications of sufficient complexity cannot simply relay (**Type 1** client) error responses to the user, for various reasons:

// TODO reasons

### Error Tracking

...unique ID per response
...log details (conditionally to avoid DDOS?)

### Localisation

### Self-service

...provide details of the error for users to understand what has gone wrong
...interaction with support

### Proxies

...mangling response code

### Security Considerations

...exposing system internals by relaying errors verbatim

## Solution

Modern APIs must take the same approach that almost all programming languages and operating systems take, which is to return
specific _error codes_ that describe as precisely as possible what has gone wrong.

...still use HTTP status code but provide additional semantic layer
...define error code to divorce message from meaning
  ...localisation
  ...client-generated messaging
...error code enumeration in file
  ...common error definition shared by both client and server
    ...code generation
      ...versioning through hash of contents
      ...add to CI
  ...publish error definitions along with API spec
    ...errors have the same semantic importance as regular responses
...logging
  ...unique transaction ID per entry
  ...searchable, debuggable
...obfuscation
  ...security, leaking system internals
...UI
  ...link to errors (RFC 7807 has nice ideas here)
  ...hide sensitive error codes
...HTTP status cost in response
  ...proxies can mangle response (RFC 7807)
...templates are configurable
...makes no demands over response format
  ...recommend using X-Errata-Code header
...schema should be relatively rigid
  ...additional fields can be defined BUT this introduces possible issues with client
  ...validatable schemas for metadata (too far?)
  ...versions for error codes (too far?)
    ...probably not worth putting in countermeasures for these, users must protect themselves

## See also

- [RFC 7807 - Problem Details for HTTP APIs](https://datatracker.ietf.org/doc/html/rfc7807)