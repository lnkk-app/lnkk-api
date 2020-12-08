### GET /r/:short

#### Description

Takes a shortform URL (e.g. https://lnkk.host/r/6fb841d5b0e2) and redirects the browser to the matching longform URL.

#### Query parameters

* :short -> the Short ID of the asset
* mtu_source
* mtu_medium
* mtu_campaign
* mtu_content

#### Example

https://lnkk.host/r/12345&mtu_source=blog&mtu_campaign=hackernews
