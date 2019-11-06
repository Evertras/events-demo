# Friends service

This service handles friend lists.

## Neo4j

Because ENTERPRISE, let's try a graph database.  We'll use Neo4j.
[The go driver](https://github.com/neo4j/neo4j-go-driver) requires
cgo wrapping around their C client, which is a bit of a pain.

The complicated bits have been handled by a Docker image in
[Dockerfile.dev](./Dockerfile.dev) which is referenced by the `dev`
task.

```bash
task dev

# ...after Telepresence finishes loading, we can use
# "--tags seabolt_static" with go commands to just run
# seabolt statically.
go run --tags seabolt_static cmd/server/main.go
go build --tags seabolt_static cmd/server/main.go
# etc.
```

## Actions

The following actions can be done:

### Invite

A player can invite another player to be their friend.  The friend
is not added yet; the other player must accept first.

Requires the user ID of the other player.

### Accept

A player can accept a pending invitation from another player.  Once
the invitation is accepted, the two players are considered friends.

### Get

A player should be able to get their friend list at any time.

### Unfriend

A player can remove a friend.  If a friend is removed, the action
is mutual; both players will be removed from each others' lists.
Friendship is always a two way street!

## Tech

Let's see what graph databases are all about!  TBD

## Future stuff

Blocks, friend-of-friend recommendations.

