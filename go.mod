module github.com/aereal/github-ops

go 1.24.0

require (
	github.com/google/go-cmp v0.7.0
	github.com/google/go-github/v72 v72.0.0
	github.com/google/wire v0.6.0
	github.com/hashicorp/go-set/v3 v3.0.0
	go.uber.org/mock v0.5.0
	golang.org/x/crypto v0.35.0
	golang.org/x/sync v0.7.0
)

require (
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/subcommands v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
)

tool go.uber.org/mock/mockgen

tool github.com/google/wire/cmd/wire
