package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	jwt struct {
		secret string
	}
}

func newConfig() *config {
	var cfg config

	flag.IntVar(&cfg.port, "port", envToInt("GUILD_HALL_PORT", 8080), "server port")
	flag.StringVar(&cfg.env, "env", envStringWithFallback("GUILD_HALL_ENV", "development"), "server environment")

	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("GUILD_HALL_DB_DSN"), "database DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", envToInt("GUILD_HALL_DB_MAX_OPEN_CONNS", 25), "max open connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", envToDuration("GUILD_HALL_DB_MAX_IDLE_TIME", 15*time.Minute), "max idle connections")

	flag.StringVar(&cfg.jwt.secret, "jwt-secret", os.Getenv("GUILD_HALL_JWT_SECRET"), "JWT secret")

	flag.Parse()

	return &cfg
}

func envStringWithFallback(env string, fallback string) string {
	str := os.Getenv(env)
	if str == "" {
		return fallback
	}
	return str
}

func envToInt(env string, fallback int) int {
	num := os.Getenv(env)
	if num == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(num)
	if err != nil {
		log.Printf("env %s is not a number", env)
		return fallback
	}

	return parsed
}

func envToBool(env string, fallback bool) bool {
	str := os.Getenv(env)
	if str == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(str)
	if err != nil {
		return fallback
	}

	return parsed
}

func envToDuration(env string, fallback time.Duration) time.Duration {
	str := os.Getenv(env)
	if str == "" {
		return fallback
	}

	duration, err := time.ParseDuration(str)
	if err != nil {
		return fallback
	}

	return duration
}
