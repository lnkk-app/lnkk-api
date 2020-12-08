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

* POST /a/1/short
* GET /a/1/short/:short
* PUT /a/1/short/:short

### Redirecting

* GET /r/:short [optional query parameters]

---

### POST /a/1/short

#### Description

Shortens a longform URL.

#### Payload

This is an example `AssetRequest` payload:

```json
{
    "link":"Required. Link is the long form URL",
    "owner":"Optional. Owner identifies the owner of the asset",
    "parent":"Optional. ParentID is the id of the category the asset belongs to",
    "source":"Optional. Source identiefies the client who created the request",
}
```

#### Response

This is an example `AssetResponse` response:

```json
{
    "link":"Link is the long form URL",
    "short_link":"ShortLink is the ID of the shortened link",
    "preview_link":"PreviewLink is not use for now. Defaults to the canonical short link for now",
    "owner":"Owner identifies the owner of the asset",
    "token":"AccessToken is used as a 'Secret' in order to claim or access the asset",
}
```

---

### GET /a/1/short/:short

#### Description

Retrieves metadata of an asset. Depending on the ownership, either the full data or only the longform URL is returned.

#### Payload

* :short -> the Short ID of the asset

#### Query parameters

* t - the access token that proofs ownership of the asset. 

Note: If the short ID references an asset that belongs to the owner of the bearer token, the access token query parameter is optional.

#### Response

This is an example `Asset` response:

```json
{
    "long_link":"https://lnkk.host",
    "short_link":"6fb841d5b0e2",
    "preview_link":"https://lnkk.host/r/6fb841d5b0e2",
    "owner":"me",
    "token":"bd5a48639e82",
    "tags":"not_used",
    "parent":"not_used",
    "title":"Welcome - lnkk.host",
    "description":"not_used",
    "state":2,
    "last_access":1601484658,
    "source":"not_used",
    "created":1601484658,
    "modified":1601484658
}
```

All timestamps are in seconds since 1.1.1970

The following asset states are used:

* active (2) - the asset is active and redirection is in place.
* archived (3) - the owner deactivated the asset.
* expired (4) - the asset was not redirected in the last <n> days and was therefor deactivated.
* broken (5) - the longform link is no longer accessible.

---

---

### PUT /a/1/short/:short

#### Description

Updates the metadata of an asset. Not all attributes can be updated though.

#### Payload

* :short -> the Short ID of the asset

This is an example `Asset` response:

```json
{
    "tags":"tag1,tag2",
    "parent":"just a category",
    "title":"Welcome - lnkk.host",
    "description":"not_used",
    "state":2,
    "source":"not_used",
}
```

Only the above attributes are accepted.

In order to change the state of asset, e.g. in order to deactivate it, change its `state` attribute.

#### Query parameters

No query parameters.

#### Response

This is an example `StandardResponse` response:

```json
{
}
```

---

### GET /r/:short [optional query parameters]

#### Description

Takes a shortform URL (e.g. https://lnkk.host/r/6fb841d5b0e2) and redirects the browser to the matching longform URL.

#### Payload

* :short -> the Short ID of the asset

#### Query parameters

* mtu_source
* mtu_medium
* mtu_campaign
* mtu_content

Example: https://lnkk.host/r/12345&mtu_source=blog&mtu_campaign=hackernews

---
