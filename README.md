# LangZ

## Backend commands

The Go module lives in `backend/`, so run Go commands from that directory:

```bash
cd backend
go run ./cmd/api
go build ./cmd/api
go test ./... -v

swag init -g cmd/api/main.go
grep -o '"/[^"]*"' docs/swagger.json | sort -u
```

The database compose file lives at the repo root:

```bash
docker compose up -d
```

`golangci-lint run ./...` also needs to be run from `backend/`, after `golangci-lint` is installed on your machine.
