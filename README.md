# report

This is the small stub server used to generate the implementation reports hosted
at [test.activitypub.rocks](https://test.activitypub.rocks/).

This is **NOT** a reference implementation. This should **NOT** be used as a
reference for using the [go-fed/activity](https://github.com/go-fed/activity)
library. In fact, it should **NOT** be used as a reference for *any*thing. It is
**NOT** meant for any sort of production use. It is highly and purposefully
**insecure**.

The **only** reason to look at this library is to confirm that it does not 
implement any additional ActivityPub features beyond what
[go-fed/activity](https://github.com/go-fed/activity) provides out of the box
and therefore confirm that the implementation report is an accurate
representation of the library.

## Automatic Test Cases

When using the [ActivityPub Rocks Test Suite](),
the following parameters are used (when hosted on the domain `$HOST`, and the
server is being run with:

```
./repsrv -cert $CERTPATH/fullchain.pem -key $KEYPATH/privkey.pem -https -host $HOST
```

```
Actor id (a uri): https://$HOST/actor
Auth token: doNotDoThisInRealImplementations
```

A sample command being run for the automatic test cases:

```
curl -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Content-Type: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Create",
  "object": {
    "type": "Note",
    "content": "Up for some root beer floats?",
    "attributedTo": "https://'"$HOST"'/actor",
    "id": "http://tsyesika.co.uk/chat/sup-yo/"
  },
  "actor": "https://'"$HOST"'/actor",
  "id": "http://tsyesika.co.uk/act/foo-id-here/"
}' \
     -v https://$HOST/actor/outbox
```

## Federation Manual Test Cases

To test the non-automatic common test cases, the following commands are used.
Note that `$HOST` must be set.

### Outbox

Using `$TESTACCOUNT` as a test account IRI, the following test the recipient
fields (`to`, `bto`, `cc`, `bcc`):

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Create",
  "object": {
    "type": "Note",
    "content": "This is a test note."
  },
  "actor": "https://'"$HOST"'/actor",
  "to": "'"$TESTACCOUNT"'"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Create",
  "object": {
    "type": "Note",
    "content": "This is a test note."
  },
  "actor": "https://'"$HOST"'/actor",
  "bto": "'"$TESTACCOUNT"'"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Create",
  "object": {
    "type": "Note",
    "content": "This is a test note."
  },
  "actor": "https://'"$HOST"'/actor",
  "cc": "'"$TESTACCOUNT"'"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Create",
  "object": {
    "type": "Note",
    "content": "This is a test note."
  },
  "actor": "https://'"$HOST"'/actor",
  "bcc": "'"$TESTACCOUNT"'"
}' \
     -v https://$HOST/actor/outbox
```

Test `object` requirement:

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Create",
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Update",
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Delete",
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Follow",
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Add",
  "target": {
    "id": "https://'"$HOST"'/new/1",
    "type": "OrderedCollection"
  },
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Remove",
  "target": {
    "id": "https://'"$HOST"'/new/1",
    "type": "OrderedCollection"
  },
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Like",
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Block",
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Undo",
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
```

Test `target` requirement:

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Add",
  "object": {
    "id": "https://'"$HOST"'/new/1",
    "type": "Note"
  },
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Remove",
  "object": {
    "id": "https://'"$HOST"'/new/1",
    "type": "Note"
  },
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
```

Test deduplication requirement:

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Create",
  "object": {
    "type": "Note",
    "content": "This is a test note."
  },
  "actor": "https://'"$HOST"'/actor",
  "to": [
    "'"$TESTACCOUNT"'",
    "'"$TESTACCOUNT"'"
  ]
}' \
     -v https://$HOST/actor/outbox
```

Test no duplication upon receipt:

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Create",
  "id": "https://'"$HOST"'/new/1",
  "object": {
    "type": "Note",
    "id": "https://'"$HOST"'/new/2",
    "content": "This is a test note."
  },
  "actor": "https://'"$HOST"'/actor",
  "to": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
```

Test no block delivery:

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Block",
  "object": "'"$TESTACCOUNT"'",
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
```

### Inbox

Testing deduplicating same received activities, run this twice:

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Create",
  "id": "https://example.com/new/1",
  "object": {
    "type": "Note",
    "id": "https://example.com/new/2",
    "content": "This is a test note."
  },
  "actor": "https://example.com/actor",
  "to": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/inbox
```

Testing `Update` activity:

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Create",
  "id": "https://example.com/new/1",
  "object": {
    "type": "Note",
    "id": "https://example.com/new/2",
    "content": "This is a test note."
  },
  "actor": "https://example.com/actor",
  "to": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/inbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Update",
  "id": "https://example.com/new/3",
  "object": {
    "type": "Note",
    "id": "https://example.com/new/2",
    "content": "Completely new test note."
  },
  "actor": "https://example.com/actor",
  "to": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/inbox
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Update",
  "id": "https://bad.example.com/new/4",
  "object": {
    "type": "Note",
    "id": "https://example.com/new/2",
    "content": "Not allowed to update the note."
  },
  "actor": "https://bad.example.com/actor",
  "to": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/inbox
```

Testing `Delete` activity, using setup from above:

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Delete",
  "id": "https://example.com/new/8",
  "object": {
    "type": "Note",
    "id": "https://example.com/new/2"
  },
  "actor": "https://example.com/actor",
  "to": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/inbox
```

Testing `Follow` activity:

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Follow",
  "id": "https://example.com/new/9",
  "object": "https://'"$HOST"'/actor",
  "actor": "https://example.com/actor",
  "to": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/inbox
```

### Common Tests

Test fetching `inbox`:

```
curl -H "Accept: application/activity+json" https://$HOST/actor/inbox
```

Test fetching activities / deleted activities:

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Create",
  "object": {
    "type": "Note",
    "content": "This is a test note."
  },
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
curl -H "Accept: application/activity+json" https://$HOST/new/1
curl -H "Accept: application/ld+json;profile=\"https://www.w3.org/ns/activitystreams\"" https://$HOST/new/1
curl -H "Accept: application/ld+json;profile=\"https://www.w3.org/ns/activitystreams\"" -v https://$HOST/new/555
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     -H "Accept: application/ld+json;profile=https://www.w3.org/ns/activitystreams" \
     --data '{
  "type": "Delete",
  "object": "https://'"$HOST"'/new/1",
  "actor": "https://'"$HOST"'/actor"
}' \
     -v https://$HOST/actor/outbox
curl -H "Accept: application/activity+json" https://$HOST/new/1
curl -H "Accept: application/ld+json;profile=\"https://www.w3.org/ns/activitystreams\"" https://$HOST/new/1
```

```
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" \
     -H "Authorization: Bearer doNotDoThisInRealImplementations" \
     --data '{"@context":"https://www.w3.org/ns/activitystreams","actor":"https://'"$HOST"'/actor","type":"Delete","object":"https://'"$HOST"'/new/3"}' \
     -v https://$HOST/actor/outbox
curl -H "Accept: application/ld+json;profile=\"https://www.w3.org/ns/activitystreams\"" -v https://$HOST/new/3
```
