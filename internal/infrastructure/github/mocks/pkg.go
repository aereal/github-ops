//go:generate go tool mockgen -typed -destination ./secret_gen.go -package mocks github.com/aereal/github-ops/internal/infrastructure/github GHSecretScanningService
package mocks
