# `demeris-backend-models`

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
  A field contains a valid Cosmos RPC URL (`https://host:port`)