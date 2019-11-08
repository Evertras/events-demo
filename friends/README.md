# Friends service

This service handles friend lists.

## Neo4j

Because ENTERPRISE, let's try a graph database!  We'll use Neo4j
because it's the most mature and has the largest community.

[The go driver](https://github.com/neo4j/neo4j-go-driver) requires
cgo wrapping around their C client, which is a bit of a pain.  
The complicated bits have been handled by a Docker image in
[Dockerfile.dev](./Dockerfile.dev) which is referenced by the
`dev` task.

```bash
task dev

# ...after Telepresence finishes loading, we can use
# "--tags seabolt_static" with go commands to just run
# seabolt statically.
go run --tags seabolt_static cmd/server/main.go
go build --tags seabolt_static cmd/server/main.go
# etc.
```

An ingress has been added to access the Neo4j UI.  You can access it at
[http://friends-db.localhost](http://friends-db.localhost).  To connect
to the database, you'll need to run the following in a terminal:

```bash
task db-ui-forward
```

This will expose the Neo4j instance to localhost:7687.  Use this address
to connect in the UI.  Note that auth has been disabled for simplicity,
so do not enter a user or password.

## Actions

The following actions can be done:

### Invite

A player can invite another player to be their friend.  The friend
is not added yet; the other player must accept first.

Requires the user ID of the other player.

### Get Pending Invites

A player can get all pending invites from other players.

### Accept

A player can accept a pending invitation from another player.  Once
the invitation is accepted, the two players are considered friends.

### Get Friend List

A player should be able to get their friend list at any time.

### Unfriend

A player can remove a friend.  If a friend is removed, the action
is mutual; both players will be removed from each others' lists.
Friendship is always a two way street!

## Future stuff

Blocks, friend-of-friend recommendations.

