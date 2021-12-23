# `demeris-backend-models`

[![codecov](https://codecov.io/gh/allinbits/demeris-backend-models/branch/main/graph/badge.svg?token=U1YDBROZJ3)](https://codecov.io/gh/allinbits/demeris-backend-models)
[![Build status](https://github.com/allinbits/demeris-backend-models/workflows/Build/badge.svg)](https://github.com/allinbits/demeris-backend-models/commits/main)
[![Tests status](https://github.com/allinbits/demeris-backend-models/workflows/Tests/badge.svg)](https://github.com/allinbits/demeris-backend-models/commits/main)
[![Lint](https://github.com/allinbits/demeris-backend-models/workflows/Lint/badge.svg?token)](https://github.com/allinbits/demeris-backend-models/commits/main)

This repository holds all the Demeris backend models, mostly made up of SQL table definitions and generated code created with [`sqlc`](https://github.com/kyleconroy/sqlc).

## Repository organization

This repository implies that for each piece of software, we create an associated directory.

In each directory, the `sql` directory contains all the SQL definitions, which `sqlc` will consume to generate code.

`sqlc` configuration will live in the upper level.

A `Makefile` is provided to automate the process of generation for each directory.

```
gsora@BFG ~/T/n/backend-models (main)> tree -C
.
├── cns
│   ├── cns.go
│   ├── ibc.go
│   ├── sql
│   └── sqlc.json
├── go.mod
├── go.sum
├── Makefile
└── tracelistener
    ├── sql
    ├── sqlc.json
    └── tracelistener.go
```
## Custom tags

The module defines the following struct tags

* `binding:"derivationpath"`  
  A field's value conforms to a [key derivation path](https://learnmeabitcoin.com/technical/derivation-paths). 
* `binding:"cosmosrpcurl"`  
  A field contains a valid Cosmos RPC URL (`https://host:port`).  
  The implementation extends the definition to allow for:
  * `http` for local/DEV
  * BASIC auth (e.g. `https://username:pwd@host:port`) for private RPCs
  * path info (e.g. `https://host:port/foo/bar`) for PRCs behind an API gateway 