# Auth

Intended to be the single Auth service used with
Traefik's [ForwardAuth](https://docs.traefik.io/middlewares/forwardauth/) middleware.

See [this nice blog post](https://rogerwelin.github.io/traefik/api/go/auth/2019/08/19/build-external-api-with-trafik-go.html)
for more info.

## Architecture explanation

This auth service stores credential information as a projection of registration events.
**This is not currently a secure solution.**  This is just to push the idea of event driven
architecture to the limit to see where it breaks.  It may be possible to make this
actually secure with Kafka encryption and ACLs, but that is beyond the current scope.

Pieces:
* Main Kafka shared by system
* Redis to share cached credentials as a projection of registration events
* N auth server instances behind a load balancer

This auth service is built to horizontally scale primarily for reliability.  If one
instance goes down, there should be zero interruption.  Additionally, updates to the
service should require zero downtime.  Performance via scaling is a secondary concern
but may become relevant under heavy JWT validation load.

The auth instances are connected to Kafka with a randomly generated and shared
consumer group ID.  The consumer group ID is shared via Redis.  If Redis information
is lost, a new group ID will be generated and all Kafka events will be processed.

## Hashing mechanism

Passwords are **never** stored plaintext.  They are always hashed as soon as possible
using [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt) which is the generally
recommended way to hash passwords in Go.  The bcrypt package handles salting for us.

## Logins

A login attempt will make the auth service simply check an associated redis key and
try to match the hashed password.  If the password hashes match, a JWT is created
containing the user's basic information and an expiration time.  This JWT will be
sent by the client in subsequent requests in the `X-Auth-Token` header.

## Auth check

Actual authentication is handled via JWT.  A simple check is done to ensure the JWT
is valid, then a header is set for `X-User-ID` that will be passed on via Traefik's
[ForwardAuth](https://docs.traefik.io/middlewares/forwardauth/) middleware.  Any
other service behind Traefik should now be able to trust that `X-User-ID` is valid
and authenticated without having to know any mechanisms behind it.

The choice of JWT is primarily for this step; no database check is required, and because
the system is built to easily scale horizontally we are happy to trade database IO
for CPU here.

## Registration flow

An incoming registration request gets a simple filter and sanity check before it's
simply shipped off to Kafka.  Once it's on Kafka, it is then read back into the auth
service.  This way the event is recorded permanently in Kafka without requiring a
nasty "how to do atomic transactions in distributed systems" rabbit hole.

Once the registration event is read by an auth processor, the user is registered by
creating an entry in the Redis auth database.  Note that this database is not like
a traditional user database.  Consider it a read-only view of all Kafka events.

**Important note**: Because we're in distributed land at this point, the auth service
instance processing the registration request is not necessarily the same instance that
received the request from the user.

At this point, the original user request is still waiting to know if it succeeded.
The original auth server instance that handled the request will wait for the key
in Redis to be updated via a subscribe to
[keyspace notifications](https://redis.io/topics/notifications).  Once the key is
updated, the original request returns successful with a login token.

This system is somewhat complex due to the distributed and event-driven nature of
the system as a whole.  However, it does give some potential advantages.

Firstly it allows other components to be notified of user registrations or login events
and react without putting any burden of knowledge of those systems onto the auth system
itself.  For example, we could add a campaign system that checks user registration
dates and provides rewards to players without having to modify the auth system.

All data is kept as a single source of truth in Kafka.  If state is mutated incorrectly,
we can simply scrap the entire stack and let it rebuild off the Kafka event stream.
We don't have to worry about whether our local information in the database is in sync
with what's in Kafka.  Attempting to treat both pushing to the database and pushing to
Kafka as a single atomic operation is a recipe for enormous headaches.

Whether this is all worth the complexity is another question...

## Seriously this isn't secure

Redis is wide open for anyone to connect/modify.  The JWT sign key is randomly generated,
but its value is stored in the totally insecure Redis instance.  Kafka is unencrypted.  No ACLs
are in place.  Don't use this in prod for anything.

An additional concern is storing the hashed passwords in the event store where other services
can potentially read them.  How dangerous is this, actually?  If the hash is properly done,
do we even care?  An open question for debate!

