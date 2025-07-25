### Test
`ENV_PATH=../.env go test ./...`

Testing with coverage:
`ENV_PATH=../.env go test -cover ./...`

Testing with vendor:
`ENV_PATH=../.env go test -v ./...`

Testing with vendor and coverage:
`ENV_PATH=../.env go test -v -cover ./...`