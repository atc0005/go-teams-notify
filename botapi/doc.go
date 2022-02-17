/*
Package botapi is intended to provide limited support for adding user mention
functionality to messages sent to a Microsoft Teams channel.

This package is currently a work-in-progress; when complete, this package will
provide support for generating a message equivalent to the example below,
contributed by @ghokun via
https://github.com/atc0005/go-teams-notify/issues/127.

curl -X POST -H "Content-type: application/json" -d '{
    "type": "message",
    "text": "Hey <at>Some User</at> check out this message",
    "entities": [
        {
            "type":"mention",
            "mentioned":{
                "id":"some.user@company.com",
                "name":"Some User"
            },
            "text": "<at>Some User</at>"
        }
    ]
}' <webhook_url>

*/
package botapi
