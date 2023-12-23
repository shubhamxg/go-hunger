package models

import (
	"log"
	"os"
)

type envs int

const (
	JWT_SECRET envs = iota
	DB_HOST
	DB_PORT
	DB_USER
	DB_PASSWORD
	DB_DATABASE
	REDIS_HOST
	REDIS_PORT
	REDIS_PASSWORD
)

func Env(env_var envs) string {
	switch env_var {
	case JWT_SECRET:
		return os.Getenv("JWT_SECRET")
	case DB_HOST:
		return os.Getenv("DB_HOST")
	case DB_PORT:
		return os.Getenv("DB_PORT")
	case DB_USER:
		return os.Getenv("DB_USER")
	case DB_PASSWORD:
		return os.Getenv("DB_PASSWORD")
	case DB_DATABASE:
		return os.Getenv("DB_DATABASE")
	case REDIS_HOST:
		return os.Getenv("REDIS_HOST")
	case REDIS_PORT:
		return os.Getenv("REDIS_PORT")
	case REDIS_PASSWORD:
		return os.Getenv("REDIS_PASSWORD")
	default:
		log.Panicf("Env variable not found: %v", env_var)
		return ""
	}
}
