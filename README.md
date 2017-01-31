[![GoDoc](https://godoc.org/github.com/sebkinne/go-recaptcha?status.svg)](https://godoc.org/github.com/sebkinne/go-recaptcha)

# recaptcha
Package recaptcha provides support for [reCaptcha 2.0](https://www.google.com/recaptcha) user response verification. It allows the use of a custom http.Client, and will fall back to http.DefaultClient if none is supplied.

## Todo
- Testing: While this package has been tested, proper tests should be written in recaptcha_test.go.
- Optional 'remoteip' parameter. This can be automatically retrieved from the http.request, or if we are behind a proxy, via a header.
