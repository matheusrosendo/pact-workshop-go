# Pact Project Structure - Best Practices

This repository demonstrates production-ready Pact testing organization with separate Makefiles for consumer and provider services.

## Directory Structure

```
pact-workshop-go/
├── Makefile                    # Root: delegates to service Makefiles + infrastructure
├── consumer/
│   ├── Makefile               # Consumer-specific Pact operations
│   ├── pacts/                 # Generated contracts (one per provider)
│   │   └── GoAdminService-GoUserService.json
│   ├── client/
│   │   ├── client.go          # HTTP client implementation
│   │   ├── client_test.go     # Unit tests (httptest)
│   │   └── client_pact_test.go # Pact consumer tests
│   └── ...
├── provider/
│   ├── Makefile               # Provider-specific Pact operations
│   ├── user_service.go        # Provider implementation
│   ├── user_service_test.go   # Pact provider verification
│   └── ...
└── pact/
    └── bin/                   # Pact CLI tools
```

## Consumer Makefile (`consumer/Makefile`)

### Purpose
Manages all Pact operations for a **consumer service** that calls multiple providers.

### Key Targets

```bash
# From consumer/ directory:
make test-unit          # Run unit tests
make test-pact          # Run Pact tests → generates contracts
make publish            # Publish contracts to Pact Broker
make can-i-deploy-prod  # Check if safe to deploy to production
make record-deploy-prod # Record deployment after deploying
make ci                 # Full CI pipeline: test-unit → test-pact → publish
```

### Configuration
- **CONSUMER_NAME**: `GoAdminService`
- **Contracts generated**: One JSON file per provider in `pacts/`
- **Versioning**: Uses git commit SHA automatically

### Multiple Providers Example

If your consumer talks to 3 providers, organize tests like this:

```
consumer/
├── Makefile
├── client/
│   ├── user_service_client.go
│   ├── user_service_pact_test.go      # → GoAdminService-UserService.json
│   ├── order_service_client.go
│   ├── order_service_pact_test.go     # → GoAdminService-OrderService.json
│   ├── payment_service_client.go
│   └── payment_service_pact_test.go   # → GoAdminService-PaymentService.json
└── pacts/
    ├── GoAdminService-UserService.json
    ├── GoAdminService-OrderService.json
    └── GoAdminService-PaymentService.json
```

Each test file specifies a different `Provider` name:

```go
// user_service_pact_test.go
mockProvider, _ := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
    Consumer: "GoAdminService",
    Provider: "UserService",  // ← Different per provider
    PactDir:  os.Getenv("PACT_DIR"),
})
```

## Provider Makefile (`provider/Makefile`)

### Purpose
Manages Pact verification for a **provider service** that is called by multiple consumers.

### Key Targets

```bash
# From provider/ directory:
make test-unit          # Run unit tests
make test-pact          # Verify contracts from ALL consumers
make can-i-deploy-prod  # Check if safe to deploy (won't break consumers)
make record-deploy-prod # Record deployment after deploying
make run                # Run provider service locally
make ci                 # Full CI pipeline: test-unit → test-pact
```

### Configuration
- **PROVIDER_NAME**: `GoUserService`
- **Verifies**: All consumer contracts from Pact Broker
- **ConsumerVersionSelectors**: Configurable in `user_service_test.go`

### Multiple Consumers

The provider automatically verifies contracts from **all consumers**:

```
Pact Broker:
  ConsumerA → ProviderX (contract 1)
  ConsumerB → ProviderX (contract 2)
  ConsumerC → ProviderX (contract 3)

Provider verification:
  make test-pact → verifies all 3 contracts
```

## Root Makefile

### Purpose
- **Infrastructure setup** (install Pact, start broker)
- **Delegates** to service-specific Makefiles
- **Backward compatibility** with existing scripts

### Key Targets

```bash
# Infrastructure
make install            # Install Pact Go library
make install-cli        # Install Pact CLI tools
make broker             # Start Pact Broker (Docker)
make broker-stop        # Stop Pact Broker

# Consumer (delegates to consumer/Makefile)
make consumer-test-pact
make consumer-publish
make consumer-ci

# Provider (delegates to provider/Makefile)
make provider-test-pact
make provider-run
make provider-ci

# Full workflow
make full-ci            # Run complete CI for both services
```

## Typical Workflows

### Consumer Development Workflow

```bash
# 1. Make changes to consumer code
vim consumer/client/client.go

# 2. Update Pact tests
vim consumer/client/client_pact_test.go

# 3. Run tests locally
cd consumer
make test-unit
make test-pact          # Generates contracts

# 4. Publish to broker
make publish

# 5. Check if safe to deploy
make can-i-deploy-prod

# 6. Deploy (your deployment script)
./deploy.sh

# 7. Record deployment
make record-deploy-prod
```

### Provider Development Workflow

```bash
# 1. Make changes to provider code
vim provider/user_service.go

# 2. Run tests locally
cd provider
make test-unit
make test-pact          # Verifies all consumer contracts

# 3. Check if safe to deploy
make can-i-deploy-prod

# 4. Deploy (your deployment script)
./deploy.sh

# 5. Record deployment
make record-deploy-prod
```

### Full CI Pipeline

```bash
# From root directory
make full-ci

# This runs:
# 1. Consumer: test-unit → test-pact → publish
# 2. Provider: test-unit → test-pact
```

## Environment Management

Both Makefiles support multiple environments:

```bash
# Consumer
make can-i-deploy ENV=dev
make can-i-deploy ENV=staging
make can-i-deploy ENV=production
make record-deploy ENV=dev

# Provider
make can-i-deploy ENV=dev
make can-i-deploy ENV=staging
make can-i-deploy ENV=production
make record-deploy ENV=dev
```

Or use shortcuts:
```bash
make can-i-deploy-dev
make can-i-deploy-prod
make record-deploy-dev
make record-deploy-prod
```

## CI/CD Integration

### Consumer Pipeline (GitHub Actions)

```yaml
name: Consumer CI
on: [push]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run Consumer CI
        run: make consumer-ci
        working-directory: consumer
      - name: Can I Deploy?
        run: make can-i-deploy-prod
        working-directory: consumer
      - name: Deploy
        if: success()
        run: ./deploy.sh
      - name: Record Deployment
        run: make record-deploy-prod
        working-directory: consumer
```

### Provider Pipeline (GitHub Actions)

```yaml
name: Provider CI
on: [push]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run Provider CI
        run: make provider-ci
        working-directory: provider
      - name: Can I Deploy?
        run: make can-i-deploy-prod
        working-directory: provider
      - name: Deploy
        if: success()
        run: ./deploy.sh
      - name: Record Deployment
        run: make record-deploy-prod
        working-directory: provider
```

## Real-World Scenario: Separate Repositories

In production, consumer and provider are typically in **separate repositories**:

### Consumer Repository
```
consumer-api/
├── Makefile           # Consumer Makefile (from this example)
├── client/
│   ├── user_client_pact_test.go
│   ├── order_client_pact_test.go
│   └── payment_client_pact_test.go
└── pacts/
    ├── ConsumerAPI-UserService.json
    ├── ConsumerAPI-OrderService.json
    └── ConsumerAPI-PaymentService.json
```

### Provider Repository
```
user-service/
├── Makefile           # Provider Makefile (from this example)
├── service.go
└── service_test.go    # Verifies all consumer contracts
```

Each repository has its own CI/CD pipeline that:
1. Runs Pact tests
2. Publishes/verifies contracts
3. Checks `can-i-deploy`
4. Deploys if safe
5. Records deployment

## Key Benefits

✅ **Clear separation** - Consumer and provider concerns are isolated
✅ **Scalable** - Easy to add more providers (consumer) or handle more consumers (provider)
✅ **Environment-aware** - Support for dev, staging, production
✅ **CI/CD ready** - Each service has its own pipeline
✅ **Safe deployments** - `can-i-deploy` prevents breaking changes
✅ **Audit trail** - All deployments recorded in Pact Broker

## Configuration

### Pact Broker
Update in each Makefile:
```makefile
PACT_BROKER_PROTO = https
PACT_BROKER_URL = pact-broker.your-company.com
PACT_BROKER_USERNAME = your-username
PACT_BROKER_PASSWORD = your-password
```

Or use environment variables:
```bash
export PACT_BROKER_URL=pact-broker.company.com
export PACT_BROKER_USERNAME=user
export PACT_BROKER_PASSWORD=pass
```

### Consumer Version Selectors (Provider)

In `provider/user_service_test.go`, configure which consumer versions to verify:

```go
ConsumerVersionSelectors: []provider.Selector{
    &provider.ConsumerVersionSelector{Branch: os.Getenv("VERSION_BRANCH")},  // Current branch
    &provider.ConsumerVersionSelector{MainBranch: true},                      // Main branch
    &provider.ConsumerVersionSelector{DeployedOrReleased: true},              // All deployed versions
},
```

## Troubleshooting

### "No pacts found"
- Consumer hasn't published contracts yet
- Run `make consumer-publish` first

### "Can-i-deploy" fails
- Provider hasn't verified the contract
- Or consumer version not deployed to target environment
- Check Pact Broker UI for verification status

### "Environment does not exist"
- Create environment first:
  ```bash
  pact/bin/pact-broker create-environment --name dev --broker-base-url http://localhost:8079
  ```

## Summary

This structure provides a **production-ready** Pact setup that:
- Scales to multiple providers/consumers
- Integrates seamlessly with CI/CD
- Prevents breaking changes in production
- Maintains clear separation of concerns

Each service is self-contained with its own Makefile, making it easy to work with microservices across multiple repositories.
