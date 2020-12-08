### POST <ur>

#### Description

#### Query parameters

None.

#### Payload

This is an example `xxx` payload:

```json
{
    "foo":"bar"
}
```

#### Response

This is an example `xxx` response:

```json
{
    "foo":"bar"
}
```

#### Errors

Standard HTTP response codes.

#### Example

```shell

curl -X POST -H 'Authorization: Bearer xoxb-1234-56789abcdefghijklmnop' \
-H 'Content-type: application/json' \
--data '{"channel":"C061EG9SL","text":"I hope the tour went well, Mr. Wonka.","attachments": [{"text":"Who wins the lifetime supply of chocolate?","fallback":"You could be telling the computer exactly what it can do with a lifetime supply of chocolate.","color":"#3AA3E3","attachment_type":"default","callback_id":"select_simple_1234","actions":[{"name":"winners_list","text":"Who should win?","type":"select","data_source":"users"}]}]}' \
https://slack.com/api/chat.postMessage

```
