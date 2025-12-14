package config

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
)

// ConfigServer - конфиг сервер
type ConfigServer struct {
	Address              string `env:"ADDRESS"`           // Адрес сервера
	IntervalSave         int    `env:"STORE_INTERVAL"`    // Интервал сохранения метрик
	FileStoragePath      string `env:"FILE_STORAGE_PATH"` // Имя и путь до файла сохранения метрика
	DateBaseDSN          string `env:"DATABASE_DSN"`      // Подключение к БД
	KEY                  string `env:"KEY"`               // Ключ используемый в SHA256
	Restore              bool   `env:"RESTORE"`           // Имя и путь до файла, читается при запуске сервер
	AuditFile            string `env:"AUDIT_FILE"`        // Путь к файлу, в который сохраняются логи аудита
	AuditURL             string `env:"AUDIT_URL"`         // Путь к url, в который отправляются логи аудита
	CryptoKeyPrivatePath string `env:"CRYPTO_KEY"`        // Путь до файла с приватным ключом
	ConfigPath           string `env:"CONFIG"`            // Путь до файла с json конфигом
	TrustedSubnet        string `env:"TRUSTED_SUBNET"`    // CIDR, example "192.168.1.0/24"

}

func NewConfigServer() *ConfigServer {
	var cfg ConfigServer
	_ = env.Parse(&cfg)

	homePath, _ := os.UserHomeDir()
	homePath = filepath.Join(homePath, "private.pem")

	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	address := serverFlags.String("a", "localhost:8080", "port server")
	intervalSave := serverFlags.Int("i", 300, "time interval save")
	fileStoragePath := serverFlags.String("f", "jsonBD.json", "the path to the file to save the data")
	// postgres://metrics:qwerty!23@localhost:5432/metrics_db?sslmode=disable&search_path=metrics_schema
	// go run main.go -d='postgres://metrics:qwerty!23@localhost:5432/metrics_db?sslmode=disable&search_path=metrics_schema'
	// alias migrate-up='migrate -database "postgres://metrics:qwerty!23@localhost:5432/metrics_db?sslmode=disable&search_path=metrics_schema" -path ./migrations up'
	dateBaseDSN := serverFlags.String("d", "", "host db PostgreSQL")
	restore := serverFlags.Bool("r", true, "read file to start server")
	key := serverFlags.String("k", "", "specify the encryption key")
	cryptoKeyPath := serverFlags.String("crypto-key", homePath, "path crypto key")

	auditFile := serverFlags.String("audit-file", "", "путь к файлу, в который сохраняются логи аудита")
	auditURL := serverFlags.String("audit-url", "", "путь к url, в который отправляются логи аудита")
	trustedSubnet := serverFlags.String("t", "", "CIDR, example 192.168.1.0/24")

	configPath := serverFlags.String("c", "", "path config file")

	serverFlags.Parse(os.Args[1:])
	if cfg.Address == "" {
		cfg.Address = *address
	}

	if cfg.IntervalSave == 0 {
		cfg.IntervalSave = *intervalSave
	}

	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = *fileStoragePath
	}

	if cfg.DateBaseDSN == "" {
		cfg.DateBaseDSN = *dateBaseDSN
	}

	if cfg.KEY == "" {
		cfg.KEY = *key
	}

	if cfg.AuditFile == "" {
		cfg.AuditFile = *auditFile
	}

	if cfg.AuditURL == "" {
		cfg.AuditURL = *auditURL
	}

	if !cfg.Restore {
		cfg.Restore = *restore
	}

	if cfg.CryptoKeyPrivatePath == "" {
		cfg.CryptoKeyPrivatePath = *cryptoKeyPath
	}

	if cfg.ConfigPath == "" {
		cfg.ConfigPath = *configPath
	}

	if cfg.TrustedSubnet == "" {
		cfg.TrustedSubnet = *trustedSubnet
	}

	configJSON(&cfg)

	return &cfg
}

// configJSON - читаем данные из configJSON
// данные параметры менее приоритетны чем из командной строки или env
func configJSON(c *ConfigServer) {
	if c.ConfigPath == "" {
		return
	}
	configJSON := NewConfigServerJSON(c.ConfigPath)

	if c.Address == "" {
		c.Address = configJSON.Address
	}

	if !c.Restore {
		c.Restore = configJSON.Restore
	}
	if c.IntervalSave == 0 {
		c.IntervalSave = configJSON.StoreInterval
	}
	if c.FileStoragePath == "" {
		c.FileStoragePath = configJSON.StoreFile
	}
	if c.CryptoKeyPrivatePath == "" {
		c.CryptoKeyPrivatePath = configJSON.CryptoKey
	}
	if c.DateBaseDSN == "" {
		c.DateBaseDSN = configJSON.Database
	}

	if c.TrustedSubnet == "" {
		c.TrustedSubnet = configJSON.TrustedSubnet

	}
}
