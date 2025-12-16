package environment

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port      string `env:"PORT" envDefault:"8000"`
	AdminPort string `env:"ADMIN_PORT" envDefault:"8001"`
	AppEnv    string `env:"APP_ENV" envDefault:"local"`

	// Database
	DBHost     string `env:"POSTGRES_DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"POSTGRES_DB_PORT" envDefault:"5432"`
	DBUser     string `env:"POSTGRES_DB_USER" envDefault:"goodtodo"`
	DBPassword string `env:"POSTGRES_DB_PASSWORD" envDefault:"secret"`
	DBName     string `env:"POSTGRES_DB_NAME" envDefault:"goodtodo_dev"`

	// JWT
	JWTSecret           string `env:"JWT_SECRET" envDefault:"your-super-secret-key"`
	JWTExpiresIn        int    `env:"JWT_EXPIRES_IN" envDefault:"3600"`
	JWTRefreshExpiresIn int    `env:"JWT_REFRESH_EXPIRES_IN" envDefault:"604800"`

	// SMTP
	SMTPHost     string `env:"SMTP_HOST" envDefault:"localhost"`
	SMTPPort     string `env:"SMTP_PORT" envDefault:"1025"`
	SMTPUser     string `env:"SMTP_USER" envDefault:""`
	SMTPPassword string `env:"SMTP_PASSWORD" envDefault:""`
	SMTPFrom     string `env:"SMTP_FROM" envDefault:"noreply@goodtodo.local"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
