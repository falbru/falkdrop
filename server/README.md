# FalkDrop Server

## Dependencies

- go version 1.26

## Running locally

1. Copy the `.env.example` to `.env` and fill out the values

```sh
cp .env.example .env
$EDITOR .env
```

2. Migrate the database

```sh
go run cmd/migrate/main.go up
```

3. Run the server

```sh
go run cmd/server/main.go
```
