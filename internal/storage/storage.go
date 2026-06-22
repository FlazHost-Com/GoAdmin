// Package storage menyediakan penyimpanan file upload yang dapat ditukar
// (interface), dengan implementasi lokal (disk). Validasi gambar berbasis
// MAGIC-BYTE (bukan MIME klien) — lihat image.go.
package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"goadmin/internal/config"
	"goadmin/internal/helpers"
)

// Storage = kontrak penyimpanan (mudah diganti S3/OSS kemudian).
type Storage interface {
	// SaveImage menyimpan byte gambar (sudah tervalidasi) → mengembalikan URL publik.
	SaveImage(ctx context.Context, data []byte, ext string) (string, error)
}

// Local menyimpan ke folder disk yang disajikan sebagai static di URLBase.
type Local struct {
	dir     string
	urlBase string
}

// New memilih implementasi berdasar Driver: "s3" → S3, selain itu → Local.
// Bila init S3 gagal (config buruk), fallback ke Local agar app tetap jalan.
func New(cfg config.StorageConfig) Storage {
	if strings.EqualFold(cfg.Driver, "s3") {
		if s, err := NewS3(cfg); err == nil {
			return s
		}
	}
	return NewLocal(cfg)
}

// NewLocal merakit storage lokal dari config.
func NewLocal(cfg config.StorageConfig) *Local {
	return &Local{dir: cfg.Dir, urlBase: cfg.URLBase}
}

// ValidateAndSave membaca reader (mis. file upload), MEMVALIDASI + RE-ENCODE
// gambar (magic-byte + sanitasi), lalu menyimpan → URL publik. Helper DRY.
func ValidateAndSave(ctx context.Context, store Storage, r io.Reader) (string, error) {
	data, err := io.ReadAll(io.LimitReader(r, MaxImageBytes+1))
	if err != nil {
		return "", err
	}
	clean, ext, verr := SanitizeImage(data, MaxImageBytes)
	if verr != nil {
		return "", verr
	}
	return store.SaveImage(ctx, clean, ext)
}

// SaveImage menulis file bernama acak (UUID + ext) lalu mengembalikan URL publik.
func (s *Local) SaveImage(_ context.Context, data []byte, ext string) (string, error) {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return "", err
	}
	name := helpers.NewID() + ext
	if err := os.WriteFile(filepath.Join(s.dir, name), data, 0o644); err != nil {
		return "", err
	}
	return strings.TrimRight(s.urlBase, "/") + "/" + name, nil
}
