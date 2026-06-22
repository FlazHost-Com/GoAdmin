# API GoAdmin (REST)

Base path: **`/api/v1`**. Tersedia di kedua varian (`full` & `api`).

## Autentikasi

JWT **HS256** via header `Authorization: Bearer <token>`. Token didapat dari login, dicabut (blacklist) saat logout. Endpoint ber-RBAC butuh permission spesifik; **Administrator bypass** semua permission.

## Bentuk respons

```jsonc
// sukses
{ "success": true,  "message": "OK", "data": { /* ... */ } }
// gagal (status sesuai AppError: 400/401/403/404/409/422/429/500)
{ "success": false, "message": "Email atau password salah", "errors": { /* opsional per-field */ } }
```

Paginasi (list): query `page`, `per_page`, `search`; `data` berisi `{ data: [...], meta: { total, per_page, current_page, last_page, from, to } }`.

---

## Auth

| Method | Path | Auth | Keterangan |
|---|---|---|---|
| POST | `/api/v1/auth/login` | publik | body `{ email, password }` → `{ token, user }` |
| POST | `/api/v1/auth/logout` | JWT | cabut token saat ini (blacklist) |
| GET | `/api/v1/auth/me` | JWT | profil user dari token |

## Users — RBAC `user.*`

| Method | Path | Permission | Body |
|---|---|---|---|
| GET | `/api/v1/users` | `user.view` | — (query: page/per_page/search) |
| GET | `/api/v1/users/:id` | `user.view` | — |
| POST | `/api/v1/users` | `user.create` | `{ name, email, phone?, password, status?, timezone?, role_ids?[] }` |
| PUT | `/api/v1/users/:id` | `user.update` | sama, `password` opsional |
| DELETE | `/api/v1/users/:id` | `user.delete` | — |

## Roles — RBAC `role.*`

| Method | Path | Permission | Body |
|---|---|---|---|
| GET | `/api/v1/roles` | `role.view` | — |
| GET | `/api/v1/roles/:id` | `role.view` | — |
| POST | `/api/v1/roles` | `role.create` | `{ name, permission_ids?[] }` |
| PUT | `/api/v1/roles/:id` | `role.update` | `{ name, permission_ids?[] }` |
| DELETE | `/api/v1/roles/:id` | `role.delete` | — (role Administrator ditolak) |

## Permissions — RBAC `permission.*`

| Method | Path | Permission | Body |
|---|---|---|---|
| GET | `/api/v1/permissions` | `permission.view` | — |
| GET | `/api/v1/permissions/:id` | `permission.view` | — |
| POST | `/api/v1/permissions` | `permission.create` | `{ name }` |
| PUT | `/api/v1/permissions/:id` | `permission.update` | `{ name }` |
| DELETE | `/api/v1/permissions/:id` | `permission.delete` | — |

## Setting — RBAC `setting.*`

| Method | Path | Permission | Body |
|---|---|---|---|
| GET | `/api/v1/setting` | `setting.view` | — → `{ setting, themes }` |
| PUT | `/api/v1/setting` | `setting.update` | `{ name?, initial?, description?, phone?, address?, email?, copyright?, theme?, fe_template? }` (parsial) |

## Profile — JWT (tanpa permission khusus)

| Method | Path | Auth | Body |
|---|---|---|---|
| GET | `/api/v1/profile` | JWT | — |
| PUT | `/api/v1/profile` | JWT | `{ name, email, phone?, timezone?, password?, password_confirmation? }` |

> Profil least-privilege: tak bisa mengubah status/role sendiri.

## Dashboard — JWT

| Method | Path | Auth | Keterangan |
|---|---|---|---|
| GET | `/api/v1/dashboard/stats` | JWT | `{ users, roles, permissions }` |

## Lain-lain

| Method | Path | Keterangan |
|---|---|---|
| GET | `/healthz` | health check (publik) → `{ status: "ok" }` |

---

### Contoh

```bash
# Login
curl -s -X POST http://localhost:3000/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@admin.com","password":"12345678"}'

# Pakai token
TOKEN=... # dari data.token di atas
curl -s http://localhost:3000/api/v1/users -H "Authorization: Bearer $TOKEN"
```
