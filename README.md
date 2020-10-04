# archai
[EventStore](https://geteventstore.com/) replacement, leveraging Cassandra for storage.

---
Archai - plural of [Arche](https://en.wikipedia.org/wiki/Arche). It
> designates the source, origin or root of things that exist.

## Testing

Requires `docker` and `docker-compose`.

```
$ make test
```

if you encounter error

```
gocql: unable to create session: unable to discover protocol version: EOF
CreateSession failed
```

it means Cassandra is not yet up. Wait a couple more seconds and see with

```
$ docker-compose logs -f --tail=20
```

if any errors were encountered.
