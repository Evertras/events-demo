# Auth

Intended to be the single Auth service used with
Traefik's [ForwardAuth](https://docs.traefik.io/middlewares/forwardauth/) middleware.

See [this nice blog post](https://rogerwelin.github.io/traefik/api/go/auth/2019/08/19/build-external-api-with-trafik-go.html)
for more info.

This auth service stores user information in a simple postgres database.  A basic hash
is done because I can't bring myself to store things in plaintext, but this is **NOT**
a secure solution.  The point of this is to provide a simple example of the concept
of how an auth service might fit into the overall architecture, not a foundation for
a real auth solution.

