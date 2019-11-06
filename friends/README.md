# Friends service

This service handles friend lists.

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

