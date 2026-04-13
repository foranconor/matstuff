# matstuff

Internal tool for a staircase/joinery business to manage materials, suppliers, contacts, and handrail brackets.

## Architecture

Full-stack Go + React SPA:

- **Backend**: Go stdlib `net/http`, no framework. Entry point: `cmd/main.go`. Handlers in `handlers/`, domain types in `domain/`, database layer in `db/`.
- **Frontend**: React 19 + Vite, in `web/`. Built output served by the Go server from `web/dist/`.
- **Database**: PostgreSQL with multiple schemas (`materials`, `handrails`, `common`, and others not yet implemented: `customers`, `jobs`, `stairways`).

## Running

- Backend config comes from `.env` (DB credentials, PORT).
- **TODO**: A bash script will build the Go binary and restart a systemd unit. Not yet written.
- Frontend dev server: `cd web && npm run dev`

## Development rules

- **No ORM** — raw SQL only in the `db/` layer.
- After any frontend file edits, run `npm run build` in `web/` before testing.
- Never read `.env` or `web/.env`.
- Delete unused files — don't comment them out or leave empty shells.

## Frontend style conventions

- **CSS**: single file only — `index.css`. No per-component CSS files.
- **Braces**: always use curly braces, even for single-line `if` statements and one-liner arrow functions.
- **No ternary operators** — use `if/else` or early returns instead.

## Key patterns

- All DB access goes through `*sql.Tx` (handlers call `beginTx`, then pass the tx to `db/` functions, then `commitTx`). Read-only handlers still use a transaction for consistency.
- HTTP handlers follow: parse path params → `recieveJSON` → `beginTx` → db call → `commitTx` → `sendJSON`.
- Frontend field editing uses `InlineField` / `InlineSelect` components — click to edit, save on blur/change.
- API calls all go through `web/src/api.js` helpers (`get`, `put`, `post`, `del`).
- JSON field names between Go and React use Go's default export casing (PascalCase), e.g. `material.Name`, `mwu.Material.ID`.

## Schema overview

| Schema | Tables |
|---|---|
| `materials` | `materials`, `uses`, `suppliers`, `suppliers_to_materials`, `materials_notes` |
| `handrails` | `brackets` |
| `common` | `contacts`, `notes` |
| `customers` | `customers` (referenced, not yet implemented) |
| `jobs` | `jobs` (referenced, not yet implemented) |
| `stairways` | `stairways` (referenced, not yet implemented) |
