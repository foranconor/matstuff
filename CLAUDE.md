# matstuff

Internal tool for a staircase/joinery business to manage materials, suppliers, contacts, and handrail brackets/profiles.

## Architecture

Full-stack Go + React SPA:

- **Backend**: Go stdlib `net/http`, no framework. Entry point: `cmd/main.go`. Handlers in `handlers/`, domain types in `domain/`, database layer in `db/`.
- **Frontend**: React 19 + Vite, in `web/`. Built output served by the Go server from `web/dist/`.
- **Database**: PostgreSQL with schemas: `materials`, `handrails`, `common`.

## Running

- Backend config from `.env` (DB creds, PORT, JWT_KEY). Never read `.env` or `web/.env`.
- Build and run backend: `go run ./cmd`
- Frontend dev server: `cd web && npm run dev` (proxies `/api` to backend)
- Build frontend: `cd web && npm run build`

## Development rules

- **No ORM** — raw SQL only in the `db/` layer.
- After any frontend file edits, run `npm run build` in `web/` before testing.
- Delete unused files — don't comment them out or leave empty shells.

## Frontend style conventions

- **CSS**: single file only — `index.css`. No per-component CSS files.
- **Braces**: always use curly braces, even for single-line `if` statements and one-liner arrow functions.
- **No ternary operators** — use `if/else` or early returns instead.

## Key patterns

### Handler flow

Every handler follows this sequence:

```go
func (h *Handler) Example(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    body, ok := receiveJSON[SomeType](w, r)
    if !ok { return }
    tx, ok := beginTx(h.DB, w)
    if !ok { return }
    defer tx.Rollback()
    result, err := db.SomeFunc(tx, ...)
    if err != nil { dbErr(w, err); return }
    if !commitTx(tx, w) { return }
    sendJSON(result, w)
}
```

Read-only handlers skip `receiveJSON` but still use a transaction. DELETE handlers typically respond with `w.WriteHeader(http.StatusNoContent)` instead of `sendJSON`.

Error helpers: `dbErr(w, err)` for DB errors (500), `sendError(w, msg, status, err)` for others.

### Database layer

All `db/` functions take `*sql.Tx` as first param. No exceptions.

### JSON field naming

Go structs use exported PascalCase fields. JSON marshaling uses defaults — no `json:` tags. Frontend receives/sends PascalCase: `material.Name`, `mwu.Material.ID`.

### Frontend

- All API calls go through `web/src/api.js` helpers: `get`, `put`, `post`, `del`.
- JWT stored in localStorage (`access_token`). API helper adds `Authorization: Bearer` header automatically. 401/403 redirects to `/login`.
- List endpoints accept `?all=true` to include archived/unpublished items.
- Detail pages use tabs for sub-resources (e.g. MaterialDetail has Specs, Uses, Suppliers, Profiles, Notes tabs).
- Field editing uses `InlineField` (text/number, saves on blur) and `InlineSelect` (saves on change).
- Creation uses `Modal` component with a form.

## Authentication

Uses a challenge/response flow against a separate database (`AUTH_DB_NAME` in .env, typically `portal`):

1. `POST /api/auth/nonce` with `{ Email }` → returns a nonce
2. Client computes proof: `SHA512(nonce + SHA512(password + salt) + nonce)`
3. `POST /api/auth/tokens` with `{ Email, Proof }` → returns JWT

Auth routes are excluded from `AuthMiddleware`. All other `/api/*` routes require a valid JWT.

## Middleware (applied in order)

1. **RouteLogger** — logs method + path + X-Real-IP
2. **Auth** — validates JWT, adds user to request context; skips non-`/api/` and `/api/auth/*`
3. **CORS** — allows all origins, standard methods + headers
4. **Recovery** — catches panics, returns 500

## Routes

### Auth
| Method | Path | Handler |
|--------|------|---------|
| POST | `/api/auth/nonce` | NonceRoute |
| POST | `/api/auth/tokens` | TokensRoute |

### Materials
| Method | Path | Handler |
|--------|------|---------|
| GET | `/api/materials` | ListMaterials (`?all=true` includes archived) |
| POST | `/api/materials` | CreateMaterial |
| GET | `/api/materials/{id}` | GetMaterial (returns MaterialWithUses) |
| PUT | `/api/materials/{id}` | UpdateMaterial |
| PUT | `/api/materials/{id}/uses` | UpdateUses |
| GET | `/api/materials/{id}/suppliers` | ListMaterialSuppliers |
| POST | `/api/materials/{id}/suppliers` | AddSupplierToMaterial |
| GET | `/api/materials/{id}/profiles` | ListMaterialProfiles |
| POST | `/api/materials/{id}/profiles` | AddProfileToMaterial |
| GET | `/api/materials/{id}/notes` | ListMaterialNotes |
| POST | `/api/materials/{id}/notes` | CreateMaterialNote |
| PUT | `/api/materials/{id}/notes/{nid}` | UpdateMaterialNote |
| DELETE | `/api/materials/{id}/notes/{nid}` | DeleteMaterialNote |

### Suppliers
| Method | Path | Handler |
|--------|------|---------|
| GET | `/api/suppliers` | ListSuppliers |
| POST | `/api/suppliers` | CreateSupplier |
| GET | `/api/suppliers/{id}` | GetSupplier (returns SupplierWithContact) |
| PUT | `/api/suppliers/{id}` | UpdateSupplier |
| GET | `/api/suppliers/{id}/materials` | ListSupplierMaterials |
| POST | `/api/suppliers/{id}/materials` | AddMaterialToSupplier |
| PUT | `/api/supplier-materials/{id}` | UpdateMaterialSupplier |
| DELETE | `/api/supplier-materials/{id}` | DeleteMaterialSupplier |

### Profiles
| Method | Path | Handler |
|--------|------|---------|
| GET | `/api/profiles` | ListProfiles |
| POST | `/api/profiles` | CreateProfile |
| GET | `/api/profiles/{id}` | GetProfile |
| PUT | `/api/profiles/{id}` | UpdateProfile |
| DELETE | `/api/profiles/{id}` | DeleteProfile |
| GET | `/api/profiles/{id}/materials` | ListProfileMaterials |
| POST | `/api/profiles/{id}/materials` | AddMaterialToProfile |
| PUT | `/api/material-profiles/{id}` | UpdateMaterialProfile |
| DELETE | `/api/material-profiles/{id}` | DeleteMaterialProfile |

### Contacts
| Method | Path | Handler |
|--------|------|---------|
| GET | `/api/contacts` | ListContacts |
| POST | `/api/contacts` | CreateContact |
| GET | `/api/contacts/{id}` | GetContact |
| PUT | `/api/contacts/{id}` | UpdateContact |

### Brackets
| Method | Path | Handler |
|--------|------|---------|
| GET | `/api/brackets` | ListBrackets (`?all=true` includes archived) |
| POST | `/api/brackets` | CreateBracket |
| GET | `/api/brackets/{id}` | GetBracket |
| PUT | `/api/brackets/{id}` | UpdateBracket |
| GET | `/api/brackets/{id}/notes` | ListBracketNotes |
| POST | `/api/brackets/{id}/notes` | CreateBracketNote |
| PUT | `/api/brackets/{id}/notes/{nid}` | UpdateBracketNote |
| DELETE | `/api/brackets/{id}/notes/{nid}` | DeleteBracketNote |

## Schema

| Schema | Tables |
|--------|--------|
| `materials` | `materials`, `uses`, `suppliers`, `suppliers_to_materials`, `materials_notes` |
| `handrails` | `brackets`, `profiles`, `material_profiles`, `brackets_notes` |
| `common` | `contacts`, `notes` |

Key relationships:
- Each material has exactly one `uses` row (1:1, created together)
- Materials ↔ Suppliers: M:M via `suppliers_to_materials` (has priority, price, lead_time)
- Materials ↔ Profiles: M:M via `material_profiles` (has price)
- Materials and Brackets each have M:M notes via junction tables (`materials_notes`, `brackets_notes`)
- Suppliers reference a Contact (optional)
- Brackets reference a Supplier (optional)

## Domain types (key fields)

- **Material**: ID, Name, Nickname, Treatment, Blurb, Units, Length, Width, Thickness, MaxSpan, Density, MaxOverhang, Radius, Color, Published, Archived
- **Uses**: ID, MaterialID, Exterior, Interior, EarlyAccess, Stringers, Risers, Treads, Timber, Panel, Handrail, Published, Archived
- **MaterialWithUses**: `Material` + `Uses` (returned by GetMaterial)
- **MaterialSupplier**: ID, Priority, SupplierID, MaterialID, Price, LeadTime, SupplierName, MaterialName
- **Supplier**: ID, Name, Website, ContactID, ContactName
- **SupplierWithContact**: `Supplier` + `Contact`
- **Contact**: ID, Name, Phone, Email
- **Bracket**: ID, Name, Nickname, SupplierID (pointer), Price, Published, Archived, SupplierName
- **Profile**: ID, Name
- **MaterialProfile**: ID, MaterialID, ProfileID, Price, MaterialName, ProfileName
- **Note**: ID, Content, Created, Modified
