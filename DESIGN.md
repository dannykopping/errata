# Errata

## Background

Users interact with software via interfaces; these interfaces can be graphical, network-based, function calls, etc. When
we design software systems, we not only need to make the software _work_ but also handle error conditions that are expected,
and even unexpected. The clarity and completeness of error handling in a system is what can delight and enable a user,
or deeply frustrate or impede them if implemented poorly.

Users don't expect software to work all the time, and _when_ it breaks they need to be communicated with in a way that
explains what **went wrong and what they can do about it**.

---

Let's consider an example:

We have an eCommerce store for musical instruments. A user, _Anathi_ (based in South Africa and speaks [isiXhosa](https://en.wikipedia.org/wiki/Xhosa_language)), would like to buy a pair of 
drumsticks. In this store's database the particular drumsticks _Anathi_ is interested in was input incorrectly, which is
causing the "add to cart" action to fail.

Here's how a typical interaction occurs in modern software:

1. _Anathi_ clicks "Add to cart" on her chosen product
2. the system generates a `NullPointerException` due to a `null` value for this product's `shipping_notes` field
3. the exception is caught, and the API server responds with a `500 Internal Server Error` error
4. _Anathi_ sees the following message:
  ```
  500 Internal Server Error
  
  java.lang.NullPointerException: Cannot read field "shipping" because "this.product.notes" is null
    at com.store.product.Cart.addItem(Cart.java:37)
  ```
5. _Anathi_ has no idea what this means, and calls _Customer Support_
6. _Customer Support_ is also unsure what this error means, and contacts one of the engineers
7. the engineer cannot track down this particular error in the logs, and gives up
8. _Customer Support_ informs _Anathi_ that there's nothing they can do and apologises
9. _Anathi_ is frustrated, and goes to another online store for the sticks

_Anathi_ really wanted this product, and we don't want to lose the sale; our values are completely aligned. Error handling
can make or break a great user experience.

**As software engineers, we can and <u>should</u> do better than this for our users and business.**

## Solution

**errata**'s goal is to make errors easier to define, raise, handle, and most importantly understand.

Let's take the scenario defined above in the _Background_, and see how we might improve it:

1. _Anathi_ clicks "Add to cart" on her chosen product
2. the system generates a `NullPointerException` due to a `null` value for this product's `shipping_notes` field
3. the exception is caught, and a new exception `AddToCartException` is thrown with this product's [SKU](https://en.wikipedia.org/wiki/Stock_keeping_unit) and wraps the original exception
4. the exception is logged, along with a UUID reference "beafdad1-3770-48a1-bba7-887126ccf504":
  ```
  com.store.AddToCartException: Failed to add product "drumstick-1": ref "beafdad1-3770-48a1-bba7-887126ccf504"
      at Program.main(Program.java:9)
  Caused by: java.lang.NullPointerException: Cannot read field "shipping" because "this.product.notes" is null
      at com.store.product.Cart.addItem(Cart.java:37)
  ```
5. the API server responds with a `500 Internal Server Error` error, and a friendly error message in the locale of her browser
6. _Anathi_ sees the following message:
  ```
  Code: add_to_cart_error (500 Internal Server Error)
  Reference: beafdad1-3770-48a1-bba7-887126ccf504
  
  Message (English):
  There was a problem adding your item to your cart
  Please contact Customer Support and provide this error code: "beafdad1-3770-48a1-bba7-887126ccf504"
  
  Message (isiXhosa):
  Bekukho ingxaki ekufakeni into yakho kwinqwelo yakho
  Nceda uqhagamshelane neNkxaso yoMthengi kwaye unikeze le khowudi yempazamo: "beafdad1-3770-48a1-bba7-887126ccf504"
  ```
7. _Anathi_ contacts _Customer Support_ with the error code
8. _Customer Support_ looks up this error using the **errata Web UI**, and determines that they need to hand this off to an engineer
9. an engineer locates the error log, and knows exactly which product has an issue
10. an engineer fixes the product and notifies _Customer Support_
11. _Customer Support_ notifies _Anathi_, who is then able to buy the product

We saved the sale and delighted the customer! Errors are inevitable; providing a clear means to solving these errors is
crucial in the software we build.

Let's dive deeper into the core ideas described above, as well as some additional features of **errata**:
- error definition and usage
- error handling
- error tracking
- internalisation (i18n)
- search & self-service

### Error Definition & Usage

**errata** makes use of a YAML file with an application's list of various errors defined. The format is very straightforward,
and makes very few assumptions about or demands your application.

**errata** comes with a CLI tool called **eish** (_**e**rrata **i**nteractive **sh**ell_,
[pronounced](http://ipa-reader.xyz/?text=e%C9%AA%CA%83) "eɪʃ") which generates code based on the error definitions.

There is a [sample application](sample/) available which demonstrates this.

The errors are defined in [`errata.yml`](sample/errata.yml) file, and [here](sample/errata/errors.go) is what's
generated by **eish**. Let's look at a simple example, the definition of `invalid_email`:

```yaml
...
errors:
  ...
  invalid_email:
    message: Please provide a valid email address
    cause: Given email address is invalid
    categories: [ login ]
    labels:
      http_response_code: 400
      shell_exit_code: 1
```

The code that gets generated is:

```go
...

const (
	...
	ErrInvalidEmail        = "invalid_email"
	...
)

var list = map[string]Error{
	...
	
	ErrInvalidEmail: {
		Code:       ErrInvalidEmail,
		Message:    "Please provide a valid email address",
		Cause:      "Given email address is invalid",
		Solution:   "",
		Categories: []string{"login"},
		Labels: map[string]string{
			"http_response_code": "400",
			"shell_exit_code":    "1",
		},

		translations: map[string]Error{},
	},
	
	...
}

func NewInvalidEmailErr(wrapped error) Error {
	return NewFromCode(ErrInvalidEmail, wrapped)
}

...
```

### Error Handling

When an application error occurs, a lot of valuable context is generated (file, line number, stack trace, etc) which should
not be discarded as it's essential for debugging. All **errata** errors allow for the original error to be "wrapped", providing
access to it for logging or inspection.

Here's an example from the sample application [`sample/http/server.go`](sample/http/server.go):

```go
	r, err := json.Marshal(&s)
	if err != nil {
		return "", errata.NewResponseFormattingErr(err)
	}
```

Here is the corresponding definition of the error:

```yaml
  response_formatting:
    message: Failed to format response body
    categories: [ internal ]
    labels:
      http_response_code: 500
```

Notice the `http_response_code` label? We can add arbitrary values in the `labels` field, and use them in our application.

### Error Tracking

Errors may need to be tracked for effective debugging. In a large multi-tenant application, it may become difficult or
even impractical to determine the cause for an error without a unique code associated to each error occurrence.

The generated code contains a `UUID()` function on the `Error` type, which can be used to log with the error itself for
tracking purposes.

### Internationalisation (i18n)

Applications that need to service many languages and cultures will degrade their UX by only displaying errors in a single
language. **errata** supports i18n by allowing you to configure per-locale overrides of the `message`, `cause` and `solution`
fields of error definitions.

### Search & Self-service

**errata** provides a Web UI which can be used to search errors and explore additional details of each error.

This is crucially important for users' self-service, and many of the world's largest APIs already provide this:

- [Twilio API (one of the best IMHO)](https://www.twilio.com/docs/api/errors)
- [AWS S3 API](https://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html#ErrorCodeList)
- [Azure Blob API](https://docs.microsoft.com/en-us/rest/api/storageservices/blob-service-error-codes)
- [GCP Storage API](https://cloud.google.com/storage/docs/json_api/v1/status-codes)
- [Twitter API](https://developer.twitter.com/en/support/twitter-api/error-troubleshooting)
- [Stripe API](https://stripe.com/docs/error-codes)

## See also

- [RFC 7807 - Problem Details for HTTP APIs](https://datatracker.ietf.org/doc/html/rfc7807)