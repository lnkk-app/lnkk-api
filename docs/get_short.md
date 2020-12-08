### GET /a/1/short/:short

#### Description

Retrieves metadata of an asset. Depending on the ownership, either the full data or only the longform URL is returned.

#### Query parameters

* t - the access token that proofs ownership of the asset.
* :short -> the Short ID of the asset

Note: If the short ID references an asset that belongs to the owner of the bearer token, the access token query parameter is optional.

#### Payload

None.

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
