# Provider Service - GoUserService

This is a **self-contained provider service** ready for deployment in a distributed environment.

## 🎯 Purpose

This service provides APIs consumed by other services and uses Pact to:
- Verify consumer contracts
- Ensure backward compatibility
- Prevent breaking changes to consumers

## 📁 Repository Structure

```
provider/
├── Makefile              # All Pact operations and CI/CD workflows
├── make/
│   └── config.mk        # Configuration (broker URL, credentials, etc.)
├── cmd/
│   └── usersvc/
│       └── main.go      # Service entry point
├── user_service.go      # HTTP handlers and routing
├── user_service_test.go # Pact provider verification tests
├── repository/
│   └── user_repository.go
└── model/
    └── user.go
```

## 🚀 Quick Start

### 1. Install Dependencies

```bash
# Install Pact CLI tools (one-time setup)
make install-cli

# Install Go dependencies
make install
```

### 2. Run Tests

```bash
# Run unit tests
make test-unit

# Run Pact verification (verifies ALL consumer contracts)
make test-pact
```

### 3. Run Service

```bash
# Start the provider service
make run
# Service available at http://localhost:8080
```

## 📋 Available Commands

### Infrastructure
```bash
make install        # Install Pact Go library
make install-cli    # Install Pact CLI tools
```

### Testing
```bash
make test-unit      # Run unit tests
make test-pact      # Run Pact provider verification tests
```

### Deployment Safety
```bash
make can-i-deploy ENV=production    # Check if safe to deploy
make can-i-deploy-dev               # Shortcut for dev
make can-i-deploy-prod              # Shortcut for production

make record-deploy ENV=production   # Record deployment after deploying
make record-deploy-dev              # Shortcut for dev
make record-deploy-prod             # Shortcut for production
```

### Service
```bash
make run            # Run the provider service locally
```

### Utilities
```bash
make broker-info    # Show Pact Broker connection info
```

### CI/CD
```bash
make ci             # Full CI pipeline: test-unit → test-pact
```

## 🔄 Development Workflow

### Daily Development

```bash
# 1. Make changes to provider code
vim user_service.go

# 2. Run unit tests
make test-unit

# 3. Verify ALL consumer contracts
make test-pact

# 4. Run service locally for manual testing
make run
```

### Adding a New Endpoint

1. **Add handler** to `user_service.go`:
   ```go
   func GetOrders(w http.ResponseWriter, r *http.Request) {
       // implementation
   }
   ```

2. **Register route**:
   ```go
   func GetHTTPHandler() *http.ServeMux {
       mux := http.NewServeMux()
       mux.HandleFunc("/orders/", GetOrders)
       return mux
   }
   ```

3. **Run verification** to ensure no consumer contracts are broken:
   ```bash
   make test-pact
   ```

## 🧪 Provider Verification

### How It Works

1. **Fetches contracts** from Pact Broker for ALL consumers
2. **Starts real provider** service on a dynamic port
3. **Replays requests** from contracts against the real service
4. **Verifies responses** match consumer expectations
5. **Publishes results** back to Pact Broker

### Provider States

Provider states set up data for specific test scenarios:

```go
// user_service_test.go
var stateHandlers = models.StateHandlers{
    "User sally exists": func(setup bool, s models.ProviderState) (models.ProviderStateResponse, error) {
        userRepository = sallyExists  // Set up test data
        return models.ProviderStateResponse{}, nil
    },
    "User sally does not exist": func(setup bool, s models.ProviderState) (models.ProviderStateResponse, error) {
        userRepository = sallyDoesNotExist
        return models.ProviderStateResponse{}, nil
    },
}
```

### Request Filters

Handle dynamic data like time-based tokens:

```go
// user_service_test.go
func fixBearerToken(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Replace consumer's static token with current time-based token
        if r.Header.Get("Authorization") != "" {
            r.Header.Set("Authorization", getAuthToken())
        }
        next.ServeHTTP(w, r)
    })
}
```

### Consumer Version Selectors

Configure which consumer versions to verify in `user_service_test.go`:

```go
ConsumerVersionSelectors: []provider.Selector{
    &provider.ConsumerVersionSelector{Branch: os.Getenv("VERSION_BRANCH")},  // Current branch
    &provider.ConsumerVersionSelector{MainBranch: true},                      // Main branch
    &provider.ConsumerVersionSelector{DeployedOrReleased: true},              // All deployed versions
},
```

## 🚢 Deployment Workflow

### Before Deploying to Production

```bash
# 1. Verify all consumer contracts
make test-pact

# 2. Check if safe to deploy
make can-i-deploy-prod

# Output example:
# Computer says yes \o/
# CONSUMER       | C.VERSION | PROVIDER      | P.VERSION | SUCCESS?
# GoAdminService | abc123    | GoUserService | def456    | true
```

### After Deploying

```bash
# Record the deployment
make record-deploy-prod
```

### Complete Deployment Example

```bash
# Development
make test-pact
make can-i-deploy-dev
# ... deploy to dev ...
make record-deploy-dev

# Staging
make can-i-deploy ENV=staging
# ... deploy to staging ...
make record-deploy ENV=staging

# Production
make can-i-deploy-prod
# ... deploy to production ...
make record-deploy-prod
```

## ⚙️ Configuration

### Pact Broker

Edit `make/config.mk`:

```makefile
export PACT_BROKER_PROTO = https
export PACT_BROKER_URL = pact-broker.your-company.com
export PACT_BROKER_USERNAME = your-username
export PACT_BROKER_PASSWORD = your-password
```

Or use environment variables:

```bash
export PACT_BROKER_URL=pact-broker.company.com
export PACT_BROKER_USERNAME=user
export PACT_BROKER_PASSWORD=pass
```

### Provider Name

Edit `make/config.mk`:

```makefile
export PROVIDER_NAME = GoUserService
```

## 📊 CI/CD Integration

### GitHub Actions Example

```yaml
name: Provider CI
on: [push]
jobs:
  test-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Install Pact CLI
        run: make install-cli
        working-directory: provider
      
      - name: Run CI Pipeline
        run: make ci
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

## 🔍 Troubleshooting

### "No pacts found"
- No consumers have published contracts yet
- Check Pact Broker UI to see available contracts
- Ensure consumer has run `make publish`

### Verification fails
- Provider implementation doesn't match consumer expectations
- Check the verification output for specific failures
- Update provider code or negotiate contract changes with consumer

### "Can-i-deploy" fails
- Consumer version not deployed to target environment yet
- Or provider hasn't verified the latest consumer contract
- Check Pact Broker UI for deployment and verification status

### Provider state not found
- Add missing state handler to `stateHandlers` in `user_service_test.go`
- Or the state doesn't require setup (e.g., authentication failures)

## 🎯 Best Practices

### 1. Verify on Every Build
Always run `make test-pact` in CI to catch breaking changes early.

### 2. Use Provider States
Set up realistic test data for each scenario:
```go
"User with orders exists": func(setup bool, s models.ProviderState) {
    // Create user with orders in test database
}
```

### 3. Handle Dynamic Data
Use request filters for time-sensitive or random data:
```go
func fixDynamicData(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Normalize dynamic headers/tokens
        next.ServeHTTP(w, r)
    })
}
```

### 4. Version Selectors
Verify contracts from:
- Current branch (feature development)
- Main branch (integration)
- Deployed environments (production safety)

### 5. Can-I-Deploy
**Always** check before deploying to production:
```bash
make can-i-deploy-prod && ./deploy.sh && make record-deploy-prod
```

## 📚 Additional Resources

- [Pact Documentation](https://docs.pact.io/)
- [Pact Go](https://github.com/pact-foundation/pact-go)
- [Provider Verification](https://docs.pact.io/implementation_guides/go/readme#provider-verification)

## 🎓 Key Concepts

### Provider Verification
The provider verifies that it can satisfy all consumer expectations defined in contracts.

### Provider States
Setup functions that prepare the provider with specific data for each test scenario.

### Request Filters
Middleware that normalizes dynamic data (tokens, timestamps) before verification.

### Can-I-Deploy
Safety check that prevents deploying a provider version that would break existing consumers.

### Multiple Consumers
One provider can verify contracts from many consumers in a single test run.

---

**This repository is production-ready and can be deployed independently.**
