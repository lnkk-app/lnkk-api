### PUT /a/1/short/:short

#### Description

Updates the metadata of an asset. Not all attributes can be updated though.

#### Query parameters

* :short -> the Short ID of the asset

#### Payload

This is an example `Asset` response:

```json
{
    "token":"bd5a48639e82",
    "tags":"tag1,tag2",
    "parent":"just a category",
    "title":"Welcome - lnkk.host",
    "description":"not_used",
    "state":2,
    "source":"not_used",
}
```

The following restrictions apply:

* Only the above attributes are accepted.
* `token` can not be changed. It has to be provided to verify ownership of the asset.
* In order to change the state of an asset (e.g. in order to deactivate it), change its `state` attribute.

#### Response

This is an example `StandardResponse` response:

```json
TBD
```
