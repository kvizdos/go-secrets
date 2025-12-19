# go-secrets

An access layer for reading values from external providers, with consistent semantics and composition.

## Out-of-the-box Providers

go-secrets provides out-of-the-box support for the following providers:

- [x] Environment Variables
- [x] AWS Secrets Manager
- [x] AWS Systems Manager (SSM)
- [ ] Hashicorp Vault (open a PR, plz!)
- [ ] BitWarden Secrets Manager (open a PR, plz!)

## Out-of-the-box Preflights

Preflights are used to run logic before a secret is retrieved. This can be useful for caching secrets or ensuring that a secret is only retrieved once:

- [x] **Cache TTL**: A preflight that caches the result of a secret lookup for a specified duration.
- [x] **Single Flight**: A preflight that ensures that a secret is only retrieved once per instance.

```go
wrappedProvider := go_secrets_providers.WithPreflights(
		baseProvider,
		go_secrets_providers.WithCacheTTL(time.Minute),
		go_secrets_providers.WithSingleFlight(),
	)

```

## Why use `go-secrets`

### 1. Normalization

Every provider responds differently. Ensuring env, AWS Secrets Manager, and Hashicorp Vault all return a consistent interface is crucial for maintaining a unified experience across different environments. This library normalizes the following into stable, testable error semantics:

- not found
- access denied
- lookup failure
- empty values
- execution safety

### 2. Composition

Detailed composition of secrets providers is possible, allowing for a flexible and modular approach to managing secrets. This library provides a set of interfaces and utilities that enable the creation of custom providers and the composition of multiple providers into a single, unified interface:

```go
wrappedProvider := go_secrets_providers.WithPreflights(
		baseProvider,
		go_secrets_providers.WithCacheTTL(time.Minute),
		go_secrets_providers.WithSingleFlight(),
	)

```

### 3. Multi-provider Support

Each go-secret service can contain multiple `Channel`s. By default, `SECRETS` and `CONFIG` are defined, however you can define your own channels as needed. Each channel can have its own set of providers, allowing for a flexible and modular approach to managing secrets:

```go
secretProvider := go_secrets_aws.NewSecretsManager(awsCfg)
configProvider := go_secrets_providers.NewEnvProvider()

secrets := go_secrets.New(
	go_secrets.WithSecretProvider(secretProvider),
	go_secrets.WithConfigProvider(configProvider),
)
```

Since each provider is designed to be independent and composable, you can easily create custom providers or combine multiple providers to suit your specific needs.

This style also makes testing easier, as you can use a drop-in testing provider:

```go
func TestSomething(t *testing.T) {
	demoEnv := func(key string) string {
		return "Hello World"
	}
	
	demoProvider := go_secrets_providers.NewTestingProvider(demoEnv)
	secrets := go_secrets.New(
		go_secrets.WithSecretProvider(demoProvider),
	)
	// Run somethign that depends on a provider / svc
}
```

### 4. Fully Tested

This library is fully tested using a combination of unit tests and integration tests. The unit tests cover the core functionality of the library, while the integration tests ensure that the library works correctly with various providers and configurations (via TestContainers).
