# BLedger ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

## Table of Contents

- [BLedger ](#bledger-)
  - [Table of Contents](#table-of-contents)
  - [Requirements](#requirements)
  - [Getting Started](#getting-started)
    - [Installation](#installation)
    - [Running](#running)
    - [Using](#using)
    - [Linting](#linting)
    - [Tests](#tests)
    - [Integration Tests](#integration-tests)
    - [Environment](#environment)
    - [Layout](#layout)
    - [Description](#description)
    - [Future considerations](#future-considerations)
    - [Prompt](#prompt)

## Requirements
- docker
- [go](https://go.dev/dl/)
- [air](https://github.com/cosmtrek/air) (if running without docker)
- git
- make
- node (integration tests)
- pnpm (integration tests)

## Getting Started

### Installation
| This repo mainly makes use of docker and docker compose

Docker
- `make docker-build` - Build the BLedger dockerfile

Go + Node
- `make update` -- installs requirements
- `make build` -- builds binary to `./bin/server`

### Running

> **Warning**
>
> When deciding to run through docker or go directly, keep in mind that you have to change the `DB_DSN` env variable to maintain a proper connection. Docker Container <->  Docker Container uses an internal network
>
> DOCKER:
> - DB_DSN=postgres://postgres:postgres@host.docker.internal:5432/postgres?sslmode=disable
> 
> GO :
> - DB_DSN=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

Docker Services
- Start all services with `docker-compose up` 
  - Note: the healthcheck on the docker-compose.yml waits for the postgres instance to be fully healthy before starting the BLedger instance. This will take a few seconds to all go live.
  - Note: `BLedger` is commented out but will work if uncommented. Fastest way to run is to `docker-compose up -d` & `make run-hot` to get a reloadable server calling a local postgres instance with a persistent volume without config changes.

API
- Go + Air (preferred)
  - `make run-hot`

- Go (bin)
  - `make build`
  - `make run-direct` -- runs the server without hot reloading

- Go (raw cmd/server.go)
  - `make run-raw`

### Using

*API*

Served at http://localhost:8080 

| HTTP Verb | Route                      | Handler                                                             |
|-----------|----------------------------|---------------------------------------------------------------------|
| GET       | /v1/transactions/:id       | github.com/$user/bledger/internal/router.(*Manager).GetTransaction         |
| POST      | /v1/transactions/          | github.com/$user/bledger/internal/router.(*Manager).CreatePendingTransaction |
| PUT       | /v1/transactions/:id       | github.com/$user/bledger/internal/router.(*Manager).ExecutePendingTransaction |
| POST      | /v1/transactions/immediate | github.com/$user/bledger/internal/router.(*Manager).CreateTransaction        |
| DELETE    | /v1/transactions/:id       | github.com/$user/bledger/internal/router.(*Manager).ReverseTransaction      |
| GET       | /v1/accounts/:id           | github.com/$user/bledger/internal/router.(*Manager).GetAccount              |
| POST      | /v1/accounts/              | github.com/$user/bledger/internal/router.(*Manager).CreateAccount            |
| GET       | /health_check              | github.com/$user/bledger/internal/router.(*Manager).InitRouter.func1           |
| GET       | /                          | github.com/$user/bledger/internal/router.(*Manager).InitRouter.func2           |

### Linting
- `make lint` -- runs linter and security checker

### Tests
- `make test` -- runs controller, and other misc tests

### Integration Tests
- `docker-compose up -d` -- to start local redis and postgres
- `make run-hot` -- to start a server instance
- `make integration_test` -- runs raw typescript integration tests (not using jest for sake of time)

### Environment
| Fill in your `.env` at your root with (this is currently included and not ignored in the .gitignore):
> If using all docker containers
```bash
ENVIRONMENT=local
PORT=8080
DB_DSN=postgres://postgres:postgres@host.docker.internal:5432/postgres?sslmode=disable 
CACHE_URI=redis://localhost:6379
CACHE_PASSWORD=eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81
```

> If using docker container for redis and postgres but with a raw go server
```bash
ENVIRONMENT=local
PORT=8080
DB_DSN=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
CACHE_URI=localhost:6379
CACHE_PASSWORD=eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81
```

### Layout
```bash
├── Dockerfile
├── Makefile
├── README.md
├── bin
├── cmd
│   └── server.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── integration_tests
│   ├── account.ts
│   ├── helpers.ts
│   ├── index.ts
│   ├── interfaces.ts
│   ├── package.json
│   ├── pnpm-lock.yaml
│   ├── transactions.ts
│   └── tsconfig.json
├── internal
│   ├── cache
│   │   ├── cache.go
│   │   └── redis
│   │       └── redis.go
│   ├── common
│   │   ├── constant.go
│   │   ├── error.go
│   │   ├── error_test.go
│   │   ├── helper.go
│   │   ├── helper_test.go
│   │   ├── idempotency.go
│   │   ├── server.go
│   │   └── server_test.go
│   ├── config
│   │   └── config.go
│   ├── controller
│   │   ├── account.go
│   │   ├── controller.go
│   │   └── transaction.go
│   ├── db
│   │   └── db.go
│   ├── middleware
│   │   ├── idempotency.go
│   │   └── logger.go
│   ├── model
│   │   ├── account.go
│   │   ├── environment.go
│   │   ├── response.go
│   │   ├── transaction.go
│   │   └── version.go
│   └── router
│       ├── account.go
│       ├── router.go
│       └── transaction.go
└── pkg
    └── version.go
```

### Description
The following project is a non-versioned transction ledger that allows for 1-step and 2-step transctions to be created on an account. You can create accounts, fetch their balances and create immediate and 2-step transactions using CREDIT or DEBIT.

Immediate transactions immediately move to a `COMPLETED` state if it passes balance and currency checks, while a 2-step transaction lets the consumer create a transaction that moves to `PENDING` state with a subsequent api call that moves it into `COMPLETED` if it passes the checks

The transactions and account balance update management rely on DB transaction atomicity and mutexes. The transaction controller sets account and transaction locks to prevent multiple writes to the same row of data that could cause data loss.

For instance, if we are executing a 2-step transaction for `Account A` we set a row level lock on the both the transaction and account rows to prevent corruption. If multiple other 1-step or 2-step transactions wanted to take place, they must wait for the row-level lock to end before accessing the database

There are no server-side mutexes or channel synchronizations because we are relying on the database layer and DB locks as our mutex. When using GORM and an ACID-compliant database, relying on database transactions should be sufficient to maintain consistency and performance in the transaction ledger.

There also exists, an idempotency middleware, that allows for an api consumer to prevent duplicate writes of the transaction. The idempotency middleware is commented out in the router, but can be simply uncommented and will work amongst all apis. The idempotency keys are set in redis for hot caching and faster duplicate-write prevention.

### Future considerations
- Properly version the the transactions. For time's sake, went with a single-version transaction with an updatable `state` rather then multiple versions of transactions with a shared primary key. Ideally you would have debits and credits have some sort of version and a primary key that is generated from some components of the transaction and display it's lifecycle. Today, we simply revert the transaction and update the state and subsequent account balance changes
- Allow for 2-sided transactions to take place. E.g. instead of only `DEBIT` or `CREDIT` on `Account A`, we might have a transction of `DEBIT AccountA` and `CREDIT AccountB` and their subsequent balances checks rather than making 2 different transactions using this api.
- More unit tests, currently relying heavy on integration tests for brevity to avoid creating golang mocks


### Prompt

```
Create a small application that can track monetary balances for multiple accounts in real time.

Requirements:
    [x] System should accept credits and debits, each starting in a pending state that may be moved to a completed state at some point in the future.
    [x] Credits and debits may fail before succeeding.
    [x] Credits and debits may be reversed after initially succeeding.
    [x] Does not allow debits for more than what’s available.
    [x] Multiple requests to create transactions for an individual account can be executed in parallel. The system should remain correct in such circumstances.
    [x] Provides a way to query an account’s balance at any point in time.

```