// Package config memuat konfigurasi environment terpusat & tervalidasi.
//
// Prinsip (sejajar NodeAdmin src/config/env.ts):
//   - Akses env HANYA lewat paket ini. Modul TIDAK boleh memanggil os.Getenv
//     langsung (di-enforce convention checker).
//   - Secret wajib (SESSION_SECRET, JWT_SECRET) → fail-fast bila kosong di
//     production, agar app tak pernah jalan dengan secret default yang bisa ditebak.
//   - Tipe sudah dikonversi (int/bool/duration), bukan string mentah.
//   - Sumber: file .env (via viper) + environment OS (override).
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// AppMode menentukan varian aplikasi yang dijalankan dari satu basis kode.
type AppMode string

const (
	// ModeFull memasang lapisan web (sesi, static, layout, route web) + REST API.
	ModeFull AppMode = "full"
	// ModeAPI hanya REST + JWT (stateless), melewati lapisan web.
	ModeAPI AppMode = "api"
)

// Config adalah konfigurasi tervalidasi seluruh aplikasi.
type Config struct {
	Env     string // development | production | test
	IsProd  bool
	IsTest  bool
	App     AppConfig
	DB      DBConfig
	Redis   RedisConfig
	Session SessionConfig
	JWT     JWTConfig
	Security   SecurityConfig
	SMTP       SMTPConfig
	FeTemplate FeTemplateConfig
	Storage    StorageConfig
}

type AppConfig struct {
	Host string
	Port int
	Name string
	Mode AppMode // 'full' (UI+API) atau 'api' (API saja)
}

type DBConfig struct {
	// Type: mysql | postgres | sqlite — dialect-agnostic (ganti cukup lewat env).
	Type            string
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	Logging         bool
	ConnMaxOpen     int
	ConnMaxIdle     int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	URL string
}

type SessionConfig struct {
	Secret string
	TTL    time.Duration
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
	Algorithm string // di-pin HS256
}

type SecurityConfig struct {
	BcryptRounds int
	CORSOrigins  []string
}

// SMTPConfig = pengiriman email. Host kosong → mailer fallback (log saja, dev).
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// StorageConfig = penyimpanan file upload (gambar). Driver "local" (disk) atau
// "s3" (S3/OSS/MinIO-compatible). Lokal disajikan static di URLBase.
type StorageConfig struct {
	Driver  string // local | s3
	Dir     string // (local) folder, mis. web/uploads
	URLBase string // (local) prefix URL publik, mis. /uploads
	// S3 (driver=s3)
	S3Endpoint  string
	S3Region    string
	S3Bucket    string
	S3AccessKey string
	S3SecretKey string
	S3UseSSL    bool
	S3PublicURL string // base URL publik objek (mis. https://cdn.example.com)
}

// FeTemplateConfig = frontend template switcher (katalog landing eksternal).
// Remote=true (DEFAULT, sejajar NodeAdmin yang selalu mem-fetch) → fetch daftar
// 640 landing dari TreeURL (lazy, sekali, di-cache) + unduh HTML on-demand dari
// RawBaseURL; gagal/offline → fallback ke katalog kurasi (15). Set false untuk
// memaksa offline (hanya kurasi) — mis. lingkungan air-gapped/CI tanpa jaringan.
type FeTemplateConfig struct {
	Remote     bool
	TreeURL    string
	RawBaseURL string
	CacheDir   string
}

// Load membaca konfigurasi dari .env + environment dan memvalidasinya.
// Mengembalikan error (bukan panic) agar pemanggil (main) yang memutuskan
// fail-fast — memudahkan test memuat config tanpa mematikan proses.
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	// Abaikan bila .env tidak ada (mis. di container/CI yang inject env langsung).
	_ = v.ReadInConfig()

	setDefaults(v)

	env := strings.ToLower(v.GetString("NODE_ENV"))
	if env == "" {
		env = "development"
	}
	isProd := env == "production"
	isTest := env == "test"

	mode := ModeFull
	if strings.ToLower(v.GetString("APP_MODE")) == "api" {
		mode = ModeAPI
	}

	cfg := &Config{
		Env:    env,
		IsProd: isProd,
		IsTest: isTest,
		App: AppConfig{
			Host: v.GetString("APP_HOST"),
			Port: v.GetInt("APP_PORT"),
			Name: v.GetString("APP_NAME"),
			Mode: mode,
		},
		DB: DBConfig{
			Type:            strings.ToLower(v.GetString("DB_TYPE")),
			Host:            v.GetString("DB_HOST"),
			Port:            v.GetInt("DB_PORT"),
			Username:        v.GetString("DB_USERNAME"),
			Password:        v.GetString("DB_PASSWORD"),
			Database:        v.GetString("DB_DATABASE"),
			Logging:         v.GetBool("DB_LOGGING"),
			ConnMaxOpen:     v.GetInt("DB_CONNECTION_LIMIT"),
			ConnMaxIdle:     v.GetInt("DB_CONNECTION_IDLE"),
			ConnMaxLifetime: time.Duration(v.GetInt("DB_CONNECTION_LIFETIME_MIN")) * time.Minute,
		},
		Redis: RedisConfig{
			URL: v.GetString("REDIS_URL"),
		},
		Session: SessionConfig{
			Secret: v.GetString("SESSION_SECRET"),
			TTL:    time.Duration(v.GetInt("SESSION_TTL_HOURS")) * time.Hour,
		},
		JWT: JWTConfig{
			Secret:    v.GetString("JWT_SECRET"),
			ExpiresIn: time.Duration(v.GetInt("JWT_EXPIRES_IN_MIN")) * time.Minute,
			Algorithm: "HS256",
		},
		Security: SecurityConfig{
			BcryptRounds: v.GetInt("BCRYPT_ROUNDS"),
			CORSOrigins:  splitAndTrim(v.GetString("CORS_ORIGINS")),
		},
		SMTP: SMTPConfig{
			Host:     v.GetString("SMTP_HOST"),
			Port:     v.GetInt("SMTP_PORT"),
			Username: v.GetString("SMTP_USERNAME"),
			Password: v.GetString("SMTP_PASSWORD"),
			From:     v.GetString("SMTP_FROM"),
		},
		FeTemplate: FeTemplateConfig{
			Remote:     v.GetBool("FE_TEMPLATE_REMOTE"),
			TreeURL:    v.GetString("FE_TEMPLATE_TREE_URL"),
			RawBaseURL: v.GetString("FE_TEMPLATE_RAW_URL"),
			CacheDir:   v.GetString("FE_TEMPLATE_CACHE_DIR"),
		},
		Storage: StorageConfig{
			Driver:      v.GetString("STORAGE_DRIVER"),
			Dir:         v.GetString("STORAGE_DIR"),
			URLBase:     v.GetString("STORAGE_URL"),
			S3Endpoint:  v.GetString("S3_ENDPOINT"),
			S3Region:    v.GetString("S3_REGION"),
			S3Bucket:    v.GetString("S3_BUCKET"),
			S3AccessKey: v.GetString("S3_ACCESS_KEY"),
			S3SecretKey: v.GetString("S3_SECRET_KEY"),
			S3UseSSL:    v.GetBool("S3_USE_SSL"),
			S3PublicURL: v.GetString("S3_PUBLIC_URL"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_HOST", "http://localhost")
	v.SetDefault("APP_PORT", 3000)
	v.SetDefault("APP_NAME", "Go Admin")
	v.SetDefault("DB_TYPE", "mysql")
	v.SetDefault("DB_PORT", 3306)
	v.SetDefault("DB_LOGGING", false)
	v.SetDefault("DB_CONNECTION_LIMIT", 10)
	v.SetDefault("DB_CONNECTION_IDLE", 5)
	v.SetDefault("DB_CONNECTION_LIFETIME_MIN", 60)
	v.SetDefault("REDIS_URL", "redis://127.0.0.1:6379")
	v.SetDefault("SESSION_TTL_HOURS", 6)
	v.SetDefault("JWT_EXPIRES_IN_MIN", 60)
	v.SetDefault("BCRYPT_ROUNDS", 10)
	v.SetDefault("CORS_ORIGINS", "")
	v.SetDefault("SMTP_PORT", 587)
	v.SetDefault("SMTP_FROM", "no-reply@goadmin.local")
	v.SetDefault("FE_TEMPLATE_REMOTE", true)
	v.SetDefault("FE_TEMPLATE_TREE_URL", "https://api.github.com/repos/lindoai/opentailwind/git/trees/master?recursive=1")
	v.SetDefault("FE_TEMPLATE_RAW_URL", "https://raw.githubusercontent.com/lindoai/opentailwind/master/landings")
	v.SetDefault("FE_TEMPLATE_CACHE_DIR", "web/cache/fetemplates")
	v.SetDefault("STORAGE_DRIVER", "local")
	v.SetDefault("STORAGE_DIR", "web/uploads")
	v.SetDefault("STORAGE_URL", "/uploads")
	v.SetDefault("S3_USE_SSL", true)
}

// validate menerapkan aturan fail-fast: secret wajib di production.
func (c *Config) validate() error {
	if c.IsProd {
		var missing []string
		if c.Session.Secret == "" {
			missing = append(missing, "SESSION_SECRET")
		}
		if c.JWT.Secret == "" {
			missing = append(missing, "JWT_SECRET")
		}
		if len(missing) > 0 {
			return fmt.Errorf("config: secret wajib di production kosong: %s", strings.Join(missing, ", "))
		}
	}
	// Di luar production, beri secret dev agar app bisa jalan lokal/test.
	if c.Session.Secret == "" {
		c.Session.Secret = "dev-session-secret-change-me"
	}
	if c.JWT.Secret == "" {
		c.JWT.Secret = "dev-jwt-secret-change-me"
	}
	switch c.DB.Type {
	case "mysql", "postgres", "sqlite":
	default:
		return fmt.Errorf("config: DB_TYPE '%s' tak didukung (mysql|postgres|sqlite)", c.DB.Type)
	}
	return nil
}

func splitAndTrim(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		// Buang trailing slash agar match origin tepat (pelajaran NodeAdmin CORS).
		t := strings.TrimRight(strings.TrimSpace(p), "/")
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
