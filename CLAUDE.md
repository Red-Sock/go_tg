# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`go_tg` is a Go library (module `github.com/Red-Sock/go_tg`) that wraps `go-telegram-bot-api/v5` to provide a higher-level framework for building Telegram bots. It is a dependency, not a standalone binary — there is no `main` package.

## Commands

```bash
# Build / type-check
go build ./...

# Run tests
go test ./...

# Run a single test
go test ./path/to/package -run TestName

# Lint (if golangci-lint is available)
golangci-lint run
```

## Architecture

### Entry point — `Bot` (`client.go`, package `go_tg`)

`Bot` is the central struct. Callers create one with `NewBot(token, ...opts)`, register handlers via `AddCommandHandler` / `MustAddCommandHandler`, then call `Start()` / `Stop()`.

- `handleInComing` — long-polls Telegram updates, dispatches `Message` and `CallbackQuery` events.
- `handleMessage` — resolves which `Handler` to call based on the parsed command string. Falls back to `defaultHandler` (set with `SetDefaultCommandHandler`).
- `handleOutgoing` — sends a `MessageOut` via the underlying bot API; tolerates two known benign Telegram API quirks.
- `send.SetSender` is called at `Start()` to wire the global `send.Send` helper to `handleOutgoing`.

### Interfaces (`interfaces/`)

- `Handler` — `Handle(in *model.MessageIn, out Chat) error`
- `CommandHandler` — extends `Handler` with `GetCommand() string`; optionally implements `Description` for bot-command autocomplete.
- `Chat` — what a handler receives as `out`: `SendMessage` and `DeleteIncomingMessage`.
- `MessageOut` — the outbound message contract (`SetChatIdIfZero`, `GetMessage`, etc.).
- `ExternalContext` — a `func(*model.MessageIn) context.Context` assigned to `Bot`; called per-message to inject user metadata into `MessageIn.Ctx`.

### Model (`model/`)

- `MessageIn` — thin wrapper around `tgbotapi.Message` adding `Command`, `Args`, `Ctx`, and `IsCallback`.
- `model/response/` — builder pattern for outgoing messages: `response.New()` returns a `Builder`; call `SetText`, `SetKeyboard`, `EditMessage(id)`, or `Delete(id)`, then `Build()`.
- `model/keyboard/` — `Keyboard` interface plus `GridKeyboard`, `FloatingKeyboard`, `ReplyKeyboard` implementations and a `Builder`.

### Internal (`internal/`)

- `Chat` (`responser.go`) — implements `interfaces.Chat`; sets the chat ID on outbound messages before forwarding to `COut`.
- `DefaultHandler` — logs unrecognised commands.

### Options (`options.go`)

- `WithLogger` — swap the default `logrus.New()` logger.
- `WithOnlyDirectCalls` — in group chats, only respond to `/command@botname` style calls.

### CI

- Branch pushes matching `RSI-*` automatically open a PR via `RedSockActions/create_pr`.
- Release workflow in `.github/workflows/release.yaml`.