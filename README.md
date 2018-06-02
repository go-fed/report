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
