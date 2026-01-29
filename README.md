# Pact Go workshop

## Introduction

This workshop is aimed at demonstrating core features and benefits of contract testing with Pact.

Whilst contract testing can be applied retrospectively to systems, we will follow the [consumer driven contracts](https://martinfowler.com/articles/consumerDrivenContracts.html) approach in this workshop - where a new consumer and provider are created in parallel to evolve a service over time, especially where there is some uncertainty with what is to be built.

This workshop should take from 1 to 2 hours, depending on how deep you want to go into each topic.

**Workshop outline**:

- [step 1: **create consumer**](//github.com/pact-foundation/pact-workshop-go/tree/step1): Create our consumer before the Provider API even exists
- [step 2: **unit test**](//github.com/pact-foundation/pact-workshop-go/tree/step2): Write a unit test for our consumer
- [step 3: **pact test**](//github.com/pact-foundation/pact-workshop-go/tree/step3): Write a Pact test for our consumer
- [step 4: **pact verification**](//github.com/pact-foundation/pact-workshop-go/tree/step4): Verify the consumer pact with the Provider API
- [step 5: **fix consumer**](//github.com/pact-foundation/pact-workshop-go/tree/step5): Fix the consumer's bad assumptions about the Provider
- [step 6: **pact test**](//github.com/pact-foundation/pact-workshop-go/tree/step6): Write a pact test for `404` (missing User) in consumer
- [step 7: **provider states**](//github.com/pact-foundation/pact-workshop-go/tree/step7): Update API to handle `404` case
- [step 8: **pact test**](//github.com/pact-foundation/pact-workshop-go/tree/step8): Write a pact test for the `401` case
- [step 9: **pact test**](//github.com/pact-foundation/pact-workshop-go/tree/step9): Update API to handle `401` case
- [step 10: **request filters**](//github.com/pact-foundation/pact-workshop-go/tree/step10): Fix the provider to support the `401` case
- [step 11: **pact broker**](//github.com/pact-foundation/pact-workshop-go/tree/step11): Implement a broker workflow for integration with CI/CD

_NOTE: Each step is tied to, and must be run within, a git branch, allowing you to progress through each stage incrementally. For example, to move to step 2 run the following: `git checkout step2`_

## Learning objectives

If running this as a team workshop format, you may want to take a look through the [learning objectives](./LEARNING.md).

## Scenario

There are two components in scope for our workshop.

1. Admin Service (Consumer). Does Admin-y things, and often needs to communicate to the User service. But really, it's just a placeholder for a more useful consumer (e.g. a website or another microservice) - it doesn't do much!
1. User Service (Provider). Provides useful things about a user, such as listing all users and getting the details of individuals.

For the purposes of this workshop, we won't implement any functionality of the Admin Service, except the bits that require User information.

**Project Structure - Distributed Environment Simulation**

⚠️ **IMPORTANT**: This workshop has been restructured to simulate a **real distributed environment** where consumer and provider are in **separate repositories**.

Each service is **completely self-contained** and ready to be extracted:

```sh
├── consumer/              # SELF-CONTAINED Consumer Repository
│   ├── Makefile          # All Pact operations (test, publish, deploy)
│   ├── make/config.mk    # Configuration (broker, credentials)
│   ├── pact/bin/         # Pact CLI tools (installed per service)
│   ├── pacts/            # Generated contracts
│   ├── client/           # HTTP client + Pact tests
│   ├── model/            # Data models
│   └── README.md         # Complete consumer documentation
│
├── provider/              # SELF-CONTAINED Provider Repository
│   ├── Makefile          # All Pact operations (verify, deploy)
│   ├── make/config.mk    # Configuration (broker, credentials)
│   ├── pact/bin/         # Pact CLI tools (installed per service)
│   ├── cmd/usersvc/      # Service entry point
│   ├── user_service.go   # HTTP handlers
│   ├── user_service_test.go  # Pact verification tests
│   ├── repository/       # Data layer
│   ├── model/            # Data models
│   └── README.md         # Complete provider documentation
│
├── docker-compose.yml     # Shared Pact Broker (workshop only)
└── PACT_STRUCTURE.md      # Detailed architecture documentation
```

**Key Differences from Traditional Monorepo:**

- ✅ No root-level `go.mod`, `go.sum`, or `Makefile`
- ✅ Each service has its own Pact CLI installation
- ✅ Each service has its own configuration
- ✅ Services can be tested and deployed independently
- ✅ Ready to split into separate Git repositories

**Working with Services:**

```bash
# Consumer (in separate terminal/repo)
cd consumer
make help
make install-cli
make test-pact
make publish

# Provider (in separate terminal/repo)
cd provider
make help
make install-cli
make test-pact
```

See `consumer/README.md` and `provider/README.md` for complete documentation.

## Step 1 - Simple Consumer calling Provider

We need to first create an HTTP client to make the calls to our provider service:

![Simple Consumer](diagrams/workshop_step1.png)

_NOTE_: even if the API client had been been graciously provided for us by our Provider Team, it doesn't mean that we shouldn't write contract tests - because the version of the client we have may not always be in sync with the deployed API - and also because we will write tests on the output appropriate to our specific needs.
