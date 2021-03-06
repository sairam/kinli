[![GoDoc](https://godoc.org/github.com/sairam/kinli?status.svg)](https://godoc.org/github.com/sairam/kinli)

## Go Versions
Tested on Go versions > 1.6

## Installation

    go get github.com/sairam/kinli

or

    dep ensure -add github.com/sairam/kinli

## Using `kinli`

* `kinli.PathTemplate`
* `kinli.PathPartialTemplate`
* `kinli.CacheMode` - true or false based on production / dev
* `kinli.ViewFuncs` - additional list of ViewFuncs you'd like to access from views
* `kinli.InitTmpl()` - to start your template rendering
* `kinli.SessionStore` - this is a mandatory field to be initialized

## Sending Emails

```
var smtpConfig = &kinli.EmailSMTPConfig{config.SMTPHost, config.SMTPPort, config.SMTPUser, config.SMTPPass}
kinli.InitMailer(smtpConfig)

e := &kinli.EmailCtx{} // fill up the fields from,to, subject, TextBody, HTMLBody, optional headers
e.SendEmail()
```

## Example
See [`example1/`](https://github.com/sairam/kinli/tree/master/example1/) for a quick webpage

## What is `kinli`?
`kinli` is a code wrapper extracted for creating simple web pages quicker
* Helper functions come in as requirements come. All helper functions are taken from [`hugo`](https://github.com/spf13/hugo/) project
* repeating methods on top of sessions using the `gorilla/sessions`

**This code is extracted from [GitNotify](https://github.com/sairam/gitnotify)**

### LICENSE
MIT
