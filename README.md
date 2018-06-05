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

## Manual Test Cases

To test the non-automatic common test cases, the following commands are used.
Note that `$HOST` must be set.

```
curl -H "Accept: application/activity+json" https://$HOST/actor/inbox
curl -H "Accept: application/activity+json" https://$HOST/new/3
curl -H "Accept: application/ld+json;profile=\"https://www.w3.org/ns/activitystreams\"" https://$HOST/new/3
curl -H "Accept: application/ld+json;profile=\"https://www.w3.org/ns/activitystreams\"" -v https://$HOST/new/555
curl -H "Content-Type: application/ld+json; profile=\"https://www.w3.org/ns/activitystreams\"" -H "Authorization: Bearer doNotDoThisInRealImplementations" --data '{"@context":"https://www.w3.org/ns/activitystreams","actor":"https://$HOST/actor","type":"Delete","object":"https://$HOST/new/3"}' -v https://$HOST/actor/outbox
curl -H "Accept: application/ld+json;profile=\"https://www.w3.org/ns/activitystreams\"" -v https://$HOST/new/3
```
