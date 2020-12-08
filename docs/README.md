# API specification

### POST /a/1/short

#### Description

Shortens a longform URL.

#### Payload

```json

"AssetRequest" : {
    "link":"Required. Link is the long form URL",
    "owner":"Optional. Owner identifies the owner of the asset",
    "parent":"Optional. ParentID is the id of the category the asset belongs to",
    "source":"Optional. Source identiefies the client who created the request",
}

```

#### Response

```json

"AssetResponse": {
    "link":"Link is the long form URL",
    "short_link":"ShortLink is the ID of the shortened link",
    "preview_link":"PreviewLink is not use for now. Defaults to the canonical short link for now",
    "owner":"Owner identifies the owner of the asset",
    "token":"AccessToken is used as a 'Secret' in order to claim or access the asset",
}

```

### GET /r/:short [optional query parameters]

#### Description

Takes a shortform URL (e.g. https://lnkk.host/r/12345) and redirects the browser to the matching longform URL.

#### Payload

* :short -> the Short ID of the asset

#### Query parameters

* mtu_source
* mtu_medium
* mtu_campaign
* mtu_content

Example: https://lnkk.host/r/12345&mtu_source=blog&mtu_campaign=hackernews
