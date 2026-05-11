package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	RedisURL     string
	JWTSecret    string
	SMSCLogin    string
	SMSCPassword string
	Port         string
	Env          string
	Clinic       ClinicInfo
}

// ClinicInfo holds the issuing-organization fields required by Russian tax certificate form КНД 1151156.
type ClinicInfo struct {
	Name              string
	INN               string
	KPP               string
	LicenseNumber     string
	LicenseIssuedAt   string // YYYY-MM-DD
	LicenseIssuedBy   string
	SignatoryName     string
	SignatoryPosition string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, reading from environment")
	}
	return &Config{
		DatabaseURL:  mustEnv("DATABASE_URL"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:    mustEnv("JWT_SECRET"),
		SMSCLogin:    getEnv("SMSC_LOGIN", ""),
		SMSCPassword: getEnv("SMSC_PASSWORD", ""),
		Port:         getEnv("PORT", "8080"),
		Env:          getEnv("ENV", "development"),
		Clinic: ClinicInfo{
			Name:              getEnv("CLINIC_NAME", "ХС «БьютиМед»"),
			INN:               getEnv("CLINIC_INN", "000000000"),
			KPP:               getEnv("CLINIC_KPP", "000000000"),
			LicenseNumber:     getEnv("CLINIC_LICENSE_NUMBER", "ЛЦ-000000"),
			LicenseIssuedAt:   getEnv("CLINIC_LICENSE_ISSUED_AT", "2020-01-01"),
			LicenseIssuedBy:   getEnv("CLINIC_LICENSE_ISSUED_BY", "Министерство здравоохранения и медицинской промышленности Туркменистана"),
			SignatoryName:     getEnv("CLINIC_SIGNATORY_NAME", "Гулиев А."),
			SignatoryPosition: getEnv("CLINIC_SIGNATORY_POSITION", "Главный врач"),
		},
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required env var %s is not set", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
