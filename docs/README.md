# Specification

## Description

TBD

## Authentication

The API supports authentication using a `bearer tokens`.

Example:

```
POST https://lnkk.host/a/1/shorten
Content-type: application/json
Authorization: Bearer lnkk-your-token
{
    "foo":"bar"
}
```

## Methods

The API provides the following public endpoints:

### Asset management

* [POST /a/1/short](post_short.md)
* [GET /a/1/short/:short](get_short.md)
* [PUT /a/1/short/:short](put_short.md)

### Redirecting

* [GET /r/:short](redirect.md)
