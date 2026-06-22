// Command migrate menyiapkan skema + data awal (dev). Menjalankan AutoMigrate
// seluruh modul ber-tabel + seed admin/RBAC. Untuk PRODUKSI gunakan
// golang-migrate (.up/.down.sql portabel) — AutoMigrate hanya iterasi cepat.
//
// Pakai (sqlite dev):
//
//	APP_MODE=full DB_TYPE=sqlite DB_DATABASE=goadmin.db go run ./cmd/migrate
//	# override admin: go run ./cmd/migrate --email a@b.com --password rahasia
package main

import (
	"flag"
	"log"

	"goadmin/internal/bootstrap"
	"goadmin/internal/config"
	"goadmin/internal/database"
)

func main() {
	email := flag.String("email", "admin@admin.com", "email admin default")
	password := flag.String("password", "12345678", "password admin default")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("FATAL config: %v", err)
	}
	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("FATAL database: %v", err)
	}

	if err := bootstrap.MigrateAndSeed(db, *email, *password, cfg.Security.BcryptRounds); err != nil {
		log.Fatalf("FATAL migrate+seed: %v", err)
	}
	log.Printf("migrate + seed selesai (admin: %s)", *email)
}
