# Shared libraries

Contains some shared libraries for convenience.

***DO NOT EVER REQUIRE ANY SHARED LIBRARY TO BE UPDATED IN LOCKSTEP BETWEEN MULTIPLE SERVICES.***

For example, if Service A uses this lib, and Service B uses a library, the library
SHOULD NOT EVER require that Service A and Service B are updated at the same time.
This breaks the entire point of splitting out Service A and Service B in the first
place!

***This is a hard and fast rule that must be followed:*** shared libraries cannot make
breaking changes to themselves that would require updating multiple services at once.
For example, changing the event ID header in Kafka *must* be implemented in a backwards
compatible way so that services can update at any time.

This is from painful past experience.  Don't let it happen to you!

Seriously this is for convenience only; consider these libraries as boilerplate reduction,
NOT for enforcing logic.

