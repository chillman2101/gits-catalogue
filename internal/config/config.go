package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	AppPort       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	DBSSL         string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() *Config {
	return &Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPass:        getEnv("DB_PASSWORD", "password"),
		DBName:        getEnv("DB_NAME", "catalogue"),
		DBSSL:         getEnv("DB_SSLMODE", "disable"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getIntEnv("REDIS_DB", 0),
	}
}

func (c *Config) CreateDB() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBSSL,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	var exists bool
	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", c.DBName).Scan(&exists); err != nil {
		return err
	}
	if exists {
		return nil
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %q", c.DBName))
	return err
}

func (c *Config) ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSL,
	)
	logrus.WithFields(logrus.Fields{
		"db": c.DBName,
	}).Info("Successfully connected to DB")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		valInt, err := strconv.Atoi(val)
		if err != nil {
			return fallback
		}
		return valInt
	}
	return fallback
}
