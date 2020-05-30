# redi-shop

## Running
Builds the binary and starts the server.
```
go build
./redi-shop
```
The server now needs a postgres  or redis database to connect to, to do this, start a docker container with one of the two options.

### Postgres
```
docker run --rm --name redi_postgres -e POSTGRES_DB=redi -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres:alpine
```
The database that is created in the docker container needs the `uuid-ossp` extension to be able to generate uuid's. Enable the extension with this command:
```
docker exec -d redi_postgres psql -U postgres -h localhost -d redi -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp"'
```

### Redis
```
docker run --rm --name redi_redis -p 6379:6379 -d redis:5.0.9-alpine
```

## Testing

This command runs the `_test.go` files to verify the behavior.
```
go test
```
