// Package saldo is a double-entry accounting ledger.
//
// # Architecture
//
// Directories:
//
//   - domain - Double entry bookkeeping domain, and interfaces (ports).
//     Repositories accept domain objects as arguments and return domain objects.
//     Does not depend on anything.
//
//   - app - Application services, requests and responses.
//     Requests and Response (DTOs) live here.
//     Each file is one use case.
//     UnitOfWork interface is defined here and it allows transactions across aggregate or entity boundaries.
//     Uses UnitOfWork for all db access.
//     Defines Error type with Code which wraps domain error.
//     Depends on domain.
//
//   - infrastructure - Implementation of interfaces (adapters) from domain and app.
//     UnitOfWork is implemented here, and it is DB transaction boundary ie. there are no DB transactions in repositories, only simple methods.
//     Translates postgres driver errors into app errors (boundary).
//     Depends on domain and app.
//
//   - cmd - Runners.
//
// Runners execute application.
// They read config, wire up dependency tree, create requests from input, run application services against
// those requests to produce response, and then render those responses back to caller.
// Dependencies wired up through constructor injection, manual dependency tree build up.
// Runners also translate errors (Code field is used for this) returned from app layer into error codes suitable for caller: http error codes for http runner, exit code (either as process exit code, or explicit) for CLI.
// Depend on app and infrastructure.
//
//   - cmd/saldod - HTTP runner.
//     Defines Config struct cmd/saldod/config.go which is used only by this runner.
//     Config is read from environment.
//
//   - cmd/saldo - CLI runner.
//     Defines Config struct cmd/saldo/config.go which is used only by this runner.
//     Config is read from environment or passed through CLI arguments.
package saldo
