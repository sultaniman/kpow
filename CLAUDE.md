# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is KPow

KPow is a self-hosted, privacy-focused contact form server written in Go. It encrypts messages using Age, PGP, or RSA public keys before delivering them via SMTP email or webhook. No third-party services required.

## Common Commands

```bash
just build          # Build binary: go build -o kpow
just test           # Run all tests: go test -v ./...
just check          # Security scan: gosec ./...
just fmt            # Format: gofmt -w .
just dev            # Live reload with air
just styles         # Build Tailwind CSS (requires bun)
just setup-tools    # Install gosec, air, mockgen
```

Run a single test:

```bash
go test -v ./server/enc -run TestEncryptAge
```

Tests require `TEST_KEYS_DIR` (auto-set by Justfile to `server/enc/testkeys`). When running tests manually:

```bash
TEST_KEYS_DIR=$(pwd)/server/enc/testkeys go test -v ./...
```

Generate mocks with `mockgen` (go.uber.org/mock).

## Architecture

**Request flow:** HTTP form POST → Echo handler → encrypt message → send via SMTP mailer and/or webhook → on failure, store in inbox → cron retries from inbox.

**Key packages:**

- `cmd/` — Cobra CLI commands (`start`, `verify`). The `start` command loads config, applies CLI flag overrides, and boots the server.
- `config/` — TOML config parsing + env var overrides. Priority: TOML file → env vars → CLI flags (last wins).
- `server/` — Echo server setup, handler, error pages, embedded templates/static assets via `embed.FS`.
- `server/enc/` — Encryption providers implementing `enc.KeyLike` interface (Age, PGP, RSA).
- `server/mailer/` — `mailer.Mailer` interface with SMTP and webhook implementations. Webhook sends encrypted JSON payload.
- `server/form/` — Form binding and validation.
- `server/cron/` — Scheduled inbox retry for failed message delivery.

**Interfaces:** The codebase uses interface-based DI — `enc.KeyLike` for encryption, `mailer.Mailer` for delivery. The `Handler` struct composes these.

**Config:** Supports TOML file, environment variables, and CLI flags. Either `mailer.dsn` or `webhook.url` must be configured (webhook-only mode is supported).

**Frontend:** HTML templates + Tailwind CSS. Styles compiled with `bunx @tailwindcss/cli`. Static assets embedded in binary.

## Workflow

1. Check `docs/` for relevant ADRs and plans before exploring broadly.
2. When a plan file exists, read the ENTIRE plan before starting.
3. Summarize your intended approach and wait for confirmation before implementing.
4. Do NOT deviate from a plan's specified approach unless explicitly asked.

## Architecture & Style

- **Architecture**: DDD + Clean Architecture. Manual DI, no frameworks. Business logic in application layer, not database.
- **Philosophy**: Pragmatic MVP — core functionality over complex features, essential security with practical implementation.
- **Dependencies**: Prefer well-maintained, widely-known libraries. Leverage existing solutions over building custom.
- **Naming**: Prefer longer but clear variable names (`workspace` over `ws`, `titlebar` over `tb`).
- **Comments**: Simple docstrings only. No decorative separators (`//-----`, `//=====`, `/* ── ... ── */`), no ASCII art.
- **Complex logic**: Use mermaid diagrams over lengthy text. Add a README in the module if needed.

## Commits

Use simple commits, should not start with capital case letter. No conventional commits. Simple short sentences.

## Rust

### Implementation

- Prefer the simplest approach. Do not add complexity unless explicitly requested.
- Before implementing, grep the codebase for similar patterns and match them exactly.
- Use existing project APIs and utilities (e.g. `render_doc`, `KdfPreset::Default`, settings.toml persistence). Do not reimplement what exists.
- If a plan specifies a migration path, follow it exactly. Check `darkroom-db` for existing migration directory structure.
- Keep types simple — use `String` over `Ulid` for error types and IDs unless there's a clear reason.

### Crypto & Security

When modifying crypto or security code:

1. Read relevant files in `darkroom-crypto/` and `darkroom-db/` to understand patterns for PEK, MEK, KDF, envelope contexts.
2. Read `docs/adr/` for security-related decisions.
3. Grep for all usages of the types/functions being changed.
4. Present analysis: current patterns, proposed approach, conventions you'll follow.
5. Wait for approval before implementing.

After changes verify: no keys logged/exposed in errors, PEK redacted in audit logs, envelope contexts consistent.

### Testing

After changes to Rust files, always run:

1. `cargo check` — compilation errors
2. `cargo test` — all tests must pass
3. `cargo clippy -- -D warnings` — for security-sensitive changes
