# properties-svc

Service holds `properties` and their `values` for users,
`properties` organized into `bundles`, each bundle holds
one or more `values`.

Bundles can be tagged, and organized into hierarchy by
parent-child relationships.

The main purpose of service is to answer what values was set for
user in given time (or actual state for requests without time).

The main goal to achieve is a maximum performance in following limitations:
- mysql as storage
- separated databases (we cant access some tables in single request)

# api

##### `GET`

- `/tags` - retuns list of available tag names.
- `/bundles` - retuns list of available bundles.
- `/settings` - returns list of available settings names.
- `/settings/{user:int}[?when=RFC3339:string]` - returns list of `{"name": "...", "value": "..."}`
objects, where `name` is a setting name and `value` is a setting value for user in given time.

##### `POST`

- `/set-tag` - sets bundles for user by tag (by one at time).
- `/set-bundles` - sets bundles for user by bundle names.
- `/unset-tag` - un-sets bundles for user by tag.
- `/unset-bundles` - un-sets bundles for user by bundle names.

All `POST`-endpoints consumes following object for simplicity:
```
{
  "user_id": {int},
  "items": [{string},],
  "expire": "RFC3339:string"
}
```

# examples

set 'jun'-tagged bundles to user with id 1
```
curl -d '{"user_id": 1, "items": ["jun"]}' http://localhost:8080/set-tag
```

set 'mid'-tagged bundles to user with id 1
```
curl -d '{"user_id": 1, "items": ["mid"]}' http://localhost:8080/set-tag
```

check actual settings
```
curl http://localhost:8080/settings/1
```

set bundle 'deals-sen' for user 1
```
curl -d '{"user_id": 1, "items": ["deals-sen"]}' http://localhost:8080/set-bundles
```

un-set bundle 'deals-sen' for user 1
```
curl -d '{"user_id": 1, "items": ["deals-sen"]}' http://localhost:8080/unset-bundles
```

check history
```
curl "http://localhost:8080/settings/1?when=3020-01-01T00:01:02Z"
```
