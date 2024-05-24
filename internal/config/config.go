package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	"gihub.com/gmohmad/wb_l0/internal/utils"
)

type Config struct {
	Service
	DB
	Nats
}

type Service struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Host           string
	Port           string
	User           string
	Password       string
	DBName         string
	SSLMode        string
	MigrationsPath string
}

type Nats struct {
	ClusterId string
	ClientId  string
	Host      string
	Port      string
}

func MustLoad() *Config {

	if err := godotenv.Load(); err != nil {
		log.Fatal("can't load the env file")
	}

	configPath := utils.GetEnvOrFatal("CONFIG_PATH")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("%s no such file or derictory", configPath)
	}

	var service Service

	if err := cleanenv.ReadConfig(configPath, &service); err != nil {
		log.Fatal("can't read config")
	}

	dbHost := utils.GetEnvOrFatal("POSTGRES_HOST")
	dbPort := utils.GetEnvOrFatal("POSTGRES_PORT")
	dbUser := utils.GetEnvOrFatal("POSTGRES_USER")
	dbPassword := utils.GetEnvOrFatal("POSTGRES_PASSWORD")
	dbName := utils.GetEnvOrFatal("POSTGRES_DB")
	dbSslMode := utils.GetEnvOrFatal("SSL_MODE")
	migrationsPath := utils.GetEnvOrFatal("MIGRATIONS_PATH")

	natsClusterId := utils.GetEnvOrFatal("NATS_CLUSTER_ID")
	natsClientId := utils.GetEnvOrFatal("NATS_CLIENT_ID")
	natsHost := utils.GetEnvOrFatal("NATS_HOST")
	natsPort := utils.GetEnvOrFatal("NATS_PORT")

	db := DB{
		Host:           dbHost,
		Port:           dbPort,
		User:           dbUser,
		Password:       dbPassword,
		DBName:         dbName,
		SSLMode:        dbSslMode,
		MigrationsPath: migrationsPath,
	}

	nats := Nats{
		ClusterId: natsClusterId,
		ClientId:  natsClientId,
		Host:      natsHost,
		Port:      natsPort,
	}

	cfg := Config{
		Service: service,
		DB:      db,
		Nats:    nats,
	}

	return &cfg
}
