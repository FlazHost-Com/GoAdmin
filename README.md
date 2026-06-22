# GoAdmin

Bootstrap **admin panel** dalam Go (Gin + GORM) — port idiomatik dari [NodeAdmin](../NodeAdmin). Modular per fitur, RBAC, auth ganda (sesi web + JWT API), theme switcher, multi-database, dengan **guardrail** (convention checker + generator modul) yang menjaga konsistensi saat dikembangkan.

> Satu basis kode, dua varian dipilih runtime via `APP_MODE`: **full** (UI + API) atau **api** (REST + JWT saja).

---

## Mulai cepat (SQLite, tanpa server DB)

```bash
cp .env.example .env          # default: APP_MODE=full, DB sqlite
make migrate                  # buat tabel + seed admin (idempoten)
make run                      # http://localhost:3000
```

Login admin default: **admin@admin.com / 12345678**.
Landing publik di `/`, login di `/auth/login`, dashboard di `/admin/v1/dashboard`.

> Ganti DB ke MySQL/Postgres cukup lewat `.env` (`DB_TYPE=...`) — tanpa ubah kode.

## Perintah (Makefile)

```
make migrate   # AutoMigrate semua modul + seed admin (dev)
make run       # jalankan server (cmd/server)
make check     # convention checker (gate pola/prinsip)
make verify    # check + vet + build + test  ← gate "selesai"
make test      # go test ./...
make module ARGS="--name product"   # scaffold modul baru lengkap
```

## Fitur

- **RBAC** — User · Role · Permission (many2many), per-route, Administrator bypass. CRUD penuh via web **dan** REST API.
- **Auth ganda** — sesi web (cookie, Redis-ready) + **JWT** (HS256 di-pin) dengan **blacklist saat logout** (TTL = sisa berlaku). **Reset password via OTP email** (crypto-random, ter-hash, expiry, rate-limited).
- **Keamanan web** — CSRF token (form), rate-limit login per-IP, security headers (helmet setara), bcrypt, secret fail-fast di production, error generik di production (HTML untuk web, JSON untuk API).
- **Upload gambar aman** — `internal/storage` driver **local/disk atau S3/OSS/MinIO** (via `STORAGE_DRIVER`), validasi **magic-byte** (bukan MIME klien) + whitelist + batas ukuran + **re-encode** (buang metadata/payload); dipakai logo (Setting) & avatar (Profil/User).
- **Setting global** — singleton ber-**cache** (TTL + invalidasi saat update) + **theme switcher** (palet di DB, ganti tanpa rebuild).
- **Profile** self-service (least-privilege: tak bisa ubah status/role sendiri).
- **Dashboard** statistik, **Components** showcase UI, **Home** landing publik data-driven ke Setting + **frontend template switcher** (`/admin/v1/appearance`): builtin (Go view) + katalog **eksternal opentailwind** (opsional `FE_TEMPLATE_REMOTE=true` → fetch daftar + unduh on-demand + cache + anti-SSRF; default off → katalog kurasi).
- **Email/SMTP** — `internal/mail` (SMTP via `net/smtp`, fallback log saat dev) tersedia di container untuk reset OTP/notifikasi.
- **Multi-database** (MySQL/Postgres/SQLite) dialect-agnostic; kode dijaga portabel (tipe abstrak, `ciLike`, tanpa raw SQL vendor) — di-enforce checker.
- **Guardrail** — `cmd/checkconventions` (CI gate) + generator `make module` + `AGENTS.md`.

## Environment

Semua konfigurasi via `.env` (lihat [`.env.example`](.env.example)), dibaca **hanya** lewat `internal/config`. Kunci utama: `APP_MODE`, `DB_TYPE`/`DB_*`, `REDIS_URL`, `SESSION_SECRET`, `JWT_SECRET`, `BCRYPT_ROUNDS`, `CORS_ORIGINS`. Secret **wajib** di production (app berhenti bila kosong).

## Struktur

```
cmd/
  server/           entry-point (full/api via APP_MODE), graceful shutdown
  migrate/          AutoMigrate + seed admin (dev)
  checkconventions/ convention checker (go/ast)
  make-module/      generator modul (text/template)
internal/
  app/              perakitan engine Gin (middleware global, mount web/api)
  config/           env tervalidasi (viper)
  container/        DI manual terpusat (composition root)
  database/         koneksi GORM dialect-agnostic + SQLite in-memory (test)
  errors/           AppError (status HTTP + pesan publik)
  helpers/          paginate, ciLike, id, response (DRY)
  middleware/       error handler, security headers, CSRF, rate-limit, flash
  router/           registry modul + named routes
  view/             RenderView (html/template) + inject currentUser/_csrf/flash
  modules/          fitur: access, setting, profile, dashboard, components, home
web/templates/      layout/partials (chrome admin) bersama
tests/              unit · integration · api
```

## Dokumentasi

- [`AGENTS.md`](AGENTS.md) — aturan pengembangan (sumber kebenaran) + DO/DON'T.
- [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) — arsitektur, lifecycle, DI, varian.
- [`docs/MODULE_GUIDE.md`](docs/MODULE_GUIDE.md) — cara membuat modul (manual & generator).
- [`docs/TESTING.md`](docs/TESTING.md) — lapisan test & cara menjalankan.
- [`docs/API.md`](docs/API.md) — daftar endpoint REST.

## Testing

```bash
make verify   # konvensi + vet + build + test (gate)
make test     # hanya test
```

Test berjalan di **SQLite in-memory** (cepat, membuktikan portabilitas). Lihat [`docs/TESTING.md`](docs/TESTING.md).

## Deployment (ringkas)

1. Set `NODE_ENV=production`, `DB_TYPE` + kredensial DB, `REDIS_URL`, dan **`SESSION_SECRET`/`JWT_SECRET`** (wajib).
2. Migrasi: untuk produksi gunakan tool migrasi versioned (golang-migrate); `make migrate` (AutoMigrate) untuk dev.
3. `go build -o goadmin ./cmd/server` lalu jalankan di belakang reverse proxy (TLS). Cookie otomatis `Secure` di production.
4. Stateless → siap horizontal scaling (sesi/blacklist di Redis).

## Lisensi

ISC.
