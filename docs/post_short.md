### POST /a/1/short

#### Description

Shortens a longform URL.

#### Query parameters

None.

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
