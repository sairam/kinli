
## Installation

    go get github.com/sairam/kinli

## Using `kinli`

* `kinli.PathTemplate`
* `kinli.PathPartialTemplate`
* `kinli.CacheMode` - true or false based on production / dev
* `kinli.ViewFuncs` - additional list of ViewFuncs you'd like to access from views
* `kinli.InitTmpl()` - to start your template rendering
* `kinli.SessionStore` - this is a mandatory field to be initialized

## Example
See [`example1/`](https://github.com/sairam/kinli/tree/master/example1/) for a quick webpage

## What is `kinli`?
`kinli` is a code wrapper extracted for creating simple web pages quicker
* Helper functions come in as requirements come. All helper functions are taken from [`hugo`](https://github.com/spf13/hugo/) project
* repeating methods on top of sessions using the `gorilla/sessions`

**This code is extracted from [GitNotify](https://github.com/sairam/gitnotify)**

### LICENSE
MIT
