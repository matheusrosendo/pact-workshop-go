# Consumer Service - GoAdminService

This is a **self-contained consumer service** ready for deployment in a distributed environment.

## 🎯 Purpose

This service consumes APIs from provider services and uses Pact to:
- Define consumer-driven contracts
- Generate and publish contracts to the Pact Broker
- Verify safe deployments before going to production

## 📁 Repository Structure

```
consumer/
├── Makefile              # All Pact operations and CI/CD workflows
├── make/
│   └── config.mk        # Configuration (broker URL, credentials, etc.)
├── pacts/               # Generated contracts (one per provider)
│   └── GoAdminService-GoUserService.json
├── client/
│   ├── client.go        # HTTP client implementation
│   ├── client_test.go   # Unit tests
│   └── client_pact_test.go  # Pact consumer tests
└── model/
    └── user.go          # Data models
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
# Run unit tests (fast, no Pact)
make test-unit

# Run Pact tests (generates contracts)
make test-pact
```

### 3. Publish Contracts

```bash
# Publish to Pact Broker
make publish
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
make test-pact      # Run Pact consumer tests (generates contracts)
```

### Publishing
```bash
make publish        # Publish contracts to Pact Broker
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

### Utilities
```bash
make clean          # Clean generated pact files
make broker-info    # Show Pact Broker connection info
```

### CI/CD
```bash
make ci             # Full CI pipeline: test-unit → test-pact → publish
```

## 🔄 Development Workflow

### Adding a New Provider

1. **Create client code** for the new provider:
   ```go
   // client/order_service_client.go
   func (c *Client) GetOrder(id int) (*Order, error) {
       // implementation
   }
   ```

2. **Create Pact test** for the new provider:
   ```go
   // client/order_service_pact_test.go
   func TestOrderServicePact(t *testing.T) {
       mockProvider, _ := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
           Consumer: "GoAdminService",
           Provider: "OrderService",  // ← New provider name
           PactDir:  os.Getenv("PACT_DIR"),
       })
       // ... define interactions
   }
   ```

3. **Run tests** to generate contract:
   ```bash
   make test-pact
   # Creates: pacts/GoAdminService-OrderService.json
   ```

4. **Publish** the new contract:
   ```bash
   make publish
   ```

### Daily Development

```bash
# 1. Make changes to consumer code
vim client/client.go

# 2. Update Pact tests
vim client/client_pact_test.go

# 3. Run tests
make test-unit
make test-pact

# 4. Publish contracts
make publish
```

## 🚢 Deployment Workflow

### Before Deploying to Production

```bash
# 1. Check if safe to deploy
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
make publish
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

### Consumer Name

Edit `make/config.mk`:

```makefile
export CONSUMER_NAME = GoAdminService
```

## 🧪 Testing Strategy

### Unit Tests (`client_test.go`)
- Fast, isolated tests using `httptest`
- No Pact involvement
- Test business logic and error handling

### Pact Tests (`client_pact_test.go`)
- Define consumer expectations
- Generate contracts
- Run against Pact mock server
- Slower but verify contract compliance

**Best Practice:** Write unit tests for all scenarios, Pact tests for contract verification.

## 📊 CI/CD Integration

### GitHub Actions Example

```yaml
name: Consumer CI
on: [push]
jobs:
  test-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Install Pact CLI
        run: make install-cli
        working-directory: consumer
      
      - name: Run CI Pipeline
        run: make ci
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

## 🔍 Troubleshooting

### "No pacts found"
- Run `make test-pact` first to generate contracts
- Check that `pacts/` directory contains JSON files

### "Can-i-deploy" fails
- Provider hasn't verified the contract yet
- Check Pact Broker UI for verification status
- Ensure provider is running their verification tests

### "Environment does not exist"
- Create environment in Pact Broker:
  ```bash
  pact/bin/pact-broker create-environment --name production --broker-base-url http://localhost:8079
  ```

## 📚 Additional Resources

- [Pact Documentation](https://docs.pact.io/)
- [Pact Go](https://github.com/pact-foundation/pact-go)
- [Consumer-Driven Contracts](https://martinfowler.com/articles/consumerDrivenContracts.html)

## 🎓 Key Concepts

### Consumer-Driven Contracts
The **consumer** defines what it expects from the provider. This contract is then verified by the provider.

### Pact Broker
Central repository for:
- Storing contracts
- Tracking verifications
- Managing deployments
- Preventing breaking changes

### Can-I-Deploy
Safety check that answers: "Can I deploy this version without breaking consumers/providers?"

---

**This repository is production-ready and can be deployed independently.**
