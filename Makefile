run:
	go run ./cmd/api/main.go
air:
	air -c .air.toml

air-build-debug:
	CGO_ENABLED=0 go build -gcflags=all="-N -l" -o ./tmp/main ./cmd/api/main.go