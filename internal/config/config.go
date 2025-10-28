package config

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string `env:"DB_HOST,required"`
	DBPort     int    `env:"DB_PORT,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`
	ServerPort string `env:"SERVER_PORT,required"`

	JWTSecret  string `env:"JWT_SECRET,required"`
	TokenTTL   int    `env:"TOKEN_TTL,required"` // minutes
}

func Load() *Config {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(filename), "..", "..")
	envPath := filepath.Join(rootDir, ".env")
	
	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Предупреждение: не удалось загрузить .env файл: %v. Попытка использовать системные переменные окружения.", err)
	}

	cfg := Config{}
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию из переменных окружения: %v", err)
	}

	return &cfg
}
